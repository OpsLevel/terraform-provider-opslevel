package opslevel

import (
	"context"
	"fmt"
	"time"

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

func NewCheckAlertSourceUsageResourceModel(ctx context.Context, check opslevel.Check) CheckAlertSourceUsageResourceModel {
	var model CheckAlertSourceUsageResourceModel

	model.Category = types.StringValue(string(check.Category.Id))
	model.Enabled = types.BoolValue(check.Enabled)
	model.EnableOn = types.StringValue(check.EnableOn.Time.Format(time.RFC3339))
	model.Filter = types.StringValue(string(check.Filter.Id))
	model.Id = types.StringValue(string(check.Id))
	model.Level = types.StringValue(string(check.Level.Id))
	model.Name = types.StringValue(check.Name)
	model.Notes = types.StringValue(check.Notes)
	model.Owner = types.StringValue(string(check.Owner.Team.Id))
	model.LastUpdated = timeLastUpdated()

	model.AlertType = types.StringValue(string(check.AlertSourceType))
	model.AlertNamePredicate = NewPredicateModel(check.AlertSourceNamePredicate)

	return model
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
				Description: "",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllAlertSourceTypeEnum...)},
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

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckAlertSourceUsageCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(planModel.AlertType.ValueString()))
	if planModel.AlertNamePredicate != nil {
		input.AlertSourceNamePredicate = &opslevel.PredicateInput{
			Type:  opslevel.PredicateTypeEnum(planModel.AlertNamePredicate.Type.String()),
			Value: opslevel.RefOf(planModel.AlertNamePredicate.Value.String()),
		}
	}

	data, err := r.client.CreateCheckAlertSourceUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_alert_source_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
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
	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data)

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

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
		return
	}
	input := opslevel.CheckAlertSourceUsageUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Id:         asID(planModel.Id),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(planModel.AlertType.ValueString()))
	if planModel.AlertNamePredicate != nil {
		input.AlertSourceNamePredicate = &opslevel.PredicateUpdateInput{
			Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(planModel.AlertNamePredicate.Type.String())),
			Value: opslevel.RefOf(planModel.AlertNamePredicate.Value.String()),
		}
	}

	data, err := r.client.UpdateCheckAlertSourceUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_alert_source_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
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
