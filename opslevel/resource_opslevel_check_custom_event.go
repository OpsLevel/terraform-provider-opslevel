package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckCustomEventResource{}
	_ resource.ResourceWithImportState = &CheckCustomEventResource{}
)

func NewCheckCustomEventResource() resource.Resource {
	return &CheckCustomEventResource{}
}

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
	planModel := read[CheckCustomEventResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckCustomEventCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	input.IntegrationId = asID(planModel.Integration)
	input.PassPending = opslevel.RefOf(planModel.PassPending.ValueBool())
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
	stateModel := read[CheckCustomEventResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check custom event, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckCustomEventResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckCustomEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckCustomEventResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckCustomEventUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		Id:         asID(planModel.Id),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	input.IntegrationId = opslevel.RefOf(asID(planModel.Integration))
	input.PassPending = opslevel.RefOf(planModel.PassPending.ValueBool())
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
	stateModel := read[CheckCustomEventResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check custom event, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check custom event resource")
}

func (r *CheckCustomEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
