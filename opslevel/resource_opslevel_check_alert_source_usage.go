package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckAlertSourceUsageResource{}
	_ resource.ResourceWithImportState = &CheckAlertSourceUsageResource{}
)

func NewCheckAlertSourceUsageResource() resource.Resource {
	return &CheckAlertSourceUsageResource{}
}

// CheckAlertSourceUsageResource defines the resource implementation.
type CheckAlertSourceUsageResource struct {
	CommonResourceClient
}

type CheckAlertSourceUsageResourceModel struct {
	Category    types.String `tfsdk:"category"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	EnableOn    types.String `tfsdk:"enable_on"`
	Filter      types.String `tfsdk:"filter"`
	Id          types.String `tfsdk:"id"`
	Level       types.String `tfsdk:"level"`
	Name        types.String `tfsdk:"name"`
	Notes       types.String `tfsdk:"notes"`
	Owner       types.String `tfsdk:"owner"`
	LastUpdated types.String `tfsdk:"last_updated"`

	AlertType          types.String    `tfsdk:"alert_type"`
	AlertNamePredicate *PredicateModel `tfsdk:"alert_name_predicate"`
}

func NewCheckAlertSourceUsageResourceModel(ctx context.Context, check opslevel.Check, planModel CheckAlertSourceUsageResourceModel) CheckAlertSourceUsageResourceModel {
	var stateModel CheckAlertSourceUsageResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
	stateModel.Description = ComputedStringValue(check.Description)
	if planModel.Enabled.IsNull() {
		stateModel.Enabled = types.BoolValue(false)
	} else {
		stateModel.Enabled = OptionalBoolValue(&check.Enabled)
	}
	if planModel.EnableOn.IsNull() {
		stateModel.EnableOn = types.StringNull()
	} else {
		// We pass through the plan value because of time formatting issue to ensure the state gets the exact value the customer specified
		stateModel.EnableOn = planModel.EnableOn
	}
	stateModel.Filter = OptionalStringValue(string(check.Filter.Id))
	stateModel.Id = ComputedStringValue(string(check.Id))
	stateModel.Level = RequiredStringValue(string(check.Level.Id))
	stateModel.Name = RequiredStringValue(check.Name)
	stateModel.Notes = OptionalStringValue(check.Notes)
	stateModel.Owner = OptionalStringValue(string(check.Owner.Team.Id))

	stateModel.AlertType = types.StringValue(string(check.AlertSourceType))
	if planModel.AlertNamePredicate != nil && check.AlertSourceNamePredicate != nil {
		stateModel.AlertNamePredicate = NewPredicateModel(*check.AlertSourceNamePredicate)
	}

	return stateModel
}

func (r *CheckAlertSourceUsageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_alert_source_usage"
}

func (r *CheckAlertSourceUsageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Alert Source Usage Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"alert_type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The type of the alert source. One of `%s`",
					strings.Join(opslevel.AllAlertSourceTypeEnum, "`, `"),
				),
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf(opslevel.AllAlertSourceTypeEnum...)},
			},
			"alert_name_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckAlertSourceUsageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckAlertSourceUsageResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckAlertSourceUsageCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = &iso8601.Time{Time: enabledOn}
	}

	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(planModel.AlertType.ValueString()))
	if planModel.AlertNamePredicate != nil {
		input.AlertSourceNamePredicate = planModel.AlertNamePredicate.ToCreateInput()
	}

	data, err := r.client.CreateCheckAlertSourceUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_alert_source_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a check alert source usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckAlertSourceUsageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckAlertSourceUsageResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check alert source usage, got error: %s", err))
		return
	}
	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckAlertSourceUsageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckAlertSourceUsageResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckAlertSourceUsageUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		Id:         asID(planModel.Id),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.RefOf(planModel.Notes.ValueString()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = &iso8601.Time{Time: enabledOn}
	}

	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(planModel.AlertType.ValueString()))
	if planModel.AlertNamePredicate != nil {
		input.AlertSourceNamePredicate = planModel.AlertNamePredicate.ToUpdateInput()
	}

	data, err := r.client.UpdateCheckAlertSourceUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_alert_source_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a check alert source usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckAlertSourceUsageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckAlertSourceUsageResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check alert source usage, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check alert source usage resource")
}

func (r *CheckAlertSourceUsageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
