package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

// CheckCustomEventResource defines the resource implementation.
type CheckCustomEventResource struct {
	CommonResourceClient
}

type CheckCustomEventResourceModel struct {
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

	Integration      types.String `tfsdk:"integration"`
	PassPending      types.Bool   `tfsdk:"pass_pending"`
	ServiceSelector  types.String `tfsdk:"service_selector"`
	SuccessCondition types.String `tfsdk:"success_condition"`
	Message          types.String `tfsdk:"message"`
}

func NewCheckCustomEventResourceModel(ctx context.Context, check opslevel.Check, planModel CheckCustomEventResourceModel) CheckCustomEventResourceModel {
	var stateModel CheckCustomEventResourceModel

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

	stateModel.Integration = RequiredStringValue(string(check.CustomEventCheckFragment.Integration.Id))
	stateModel.PassPending = RequiredBoolValue(check.CustomEventCheckFragment.PassPending)
	stateModel.ServiceSelector = RequiredStringValue(check.CustomEventCheckFragment.ServiceSelector)
	stateModel.SuccessCondition = RequiredStringValue(check.CustomEventCheckFragment.SuccessCondition)
	stateModel.Message = OptionalStringValue(check.CustomEventCheckFragment.ResultMessage)

	return stateModel
}

func (r *CheckCustomEventResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_custom_event"
}

func (r *CheckCustomEventResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Custom Event Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"integration": schema.StringAttribute{
				Description: "The integration id this check will use.",
				Required:    true,
			},
			"pass_pending": schema.BoolAttribute{
				Description: "True if this check should pass by default. Otherwise the default 'pending' state counts as a failure.",
				Required:    true,
			},
			"service_selector": schema.StringAttribute{
				Description: "A jq expression that will be ran against your payload. This will parse out the service identifier.",
				Required:    true,
			},
			"success_condition": schema.StringAttribute{
				Description: "A jq expression that will be ran against your payload. A truthy value will result in the check passing.",
				Required:    true,
			},
			"message": schema.StringAttribute{
				Description: "The check result message template. It is compiled with Liquid and formatted in Markdown.",
				Optional:    true,
			},
		}),
	}
}

func (r *CheckCustomEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckCustomEventResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckCustomEventCreateInput{
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

	input.IntegrationId = asID(planModel.Integration)
	input.PassPending = planModel.PassPending.ValueBoolPointer()
	input.ServiceSelector = planModel.ServiceSelector.ValueString()
	input.SuccessCondition = planModel.SuccessCondition.ValueString()
	input.ResultMessage = opslevel.RefOf(planModel.Message.ValueString())

	data, err := r.client.CreateCheckCustomEvent(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_custom_event, got error: %s", err))
		return
	}

	stateModel := NewCheckCustomEventResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check custom event resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckCustomEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckCustomEventResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check custom event, got error: %s", err))
		return
	}
	stateModel := NewCheckCustomEventResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckCustomEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckCustomEventResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckCustomEventUpdateInput{
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

	input.IntegrationId = opslevel.RefOf(asID(planModel.Integration))
	input.PassPending = planModel.PassPending.ValueBoolPointer()
	input.ServiceSelector = opslevel.RefOf(planModel.ServiceSelector.ValueString())
	input.SuccessCondition = opslevel.RefOf(planModel.SuccessCondition.ValueString())
	input.ResultMessage = opslevel.RefOf(planModel.Message.ValueString())

	data, err := r.client.UpdateCheckCustomEvent(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_custom_event, got error: %s", err))
		return
	}

	stateModel := NewCheckCustomEventResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check custom event resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckCustomEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckCustomEventResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check custom event, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check custom event resource")
}

func (r *CheckCustomEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
