package opslevel

import (
	"context"
	"fmt"

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
	_ resource.ResourceWithConfigure   = &CheckToolUsageResource{}
	_ resource.ResourceWithImportState = &CheckToolUsageResource{}
)

func NewCheckToolUsageResource() resource.Resource {
	return &CheckToolUsageResource{}
}

// CheckToolUsageResource defines the resource implementation.
type CheckToolUsageResource struct {
	CommonResourceClient
}

type CheckToolUsageResourceModel struct {
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

	ToolCategory         types.String    `tfsdk:"tool_category"`
	ToolNamePredicate    *PredicateModel `tfsdk:"tool_name_predicate"`
	ToolUrlPredicate     *PredicateModel `tfsdk:"tool_url_predicate"`
	EnvironmentPredicate *PredicateModel `tfsdk:"environment_predicate"`
}

func NewCheckToolUsageResourceModel(ctx context.Context, check opslevel.Check, planModel CheckToolUsageResourceModel) CheckToolUsageResourceModel {
	var stateModel CheckToolUsageResourceModel

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

	stateModel.ToolCategory = RequiredStringValue(string(check.ToolCategory))
	if check.ToolNamePredicate != nil {
		stateModel.ToolNamePredicate = NewPredicateModel(*check.ToolNamePredicate)
	}
	if check.ToolUrlPredicate != nil {
		stateModel.ToolUrlPredicate = NewPredicateModel(*check.ToolUrlPredicate)
	}
	if check.EnvironmentPredicate != nil {
		stateModel.EnvironmentPredicate = NewPredicateModel(*check.EnvironmentPredicate)
	}

	return stateModel
}

func (r *CheckToolUsageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_tool_usage"
}

func (r *CheckToolUsageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Tool Usage Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"tool_category": schema.StringAttribute{
				Description: "",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllToolCategory...)},
			},
			"tool_name_predicate":   PredicateSchema(),
			"tool_url_predicate":    PredicateSchema(),
			"environment_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckToolUsageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckToolUsageCreateInput{
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

	input.ToolCategory = opslevel.ToolCategory(planModel.ToolCategory.ValueString())
	if planModel.ToolNamePredicate != nil {
		input.ToolNamePredicate = planModel.ToolNamePredicate.ToCreateInput()
	}
	if planModel.ToolUrlPredicate != nil {
		input.ToolUrlPredicate = planModel.ToolUrlPredicate.ToCreateInput()
	}
	if planModel.EnvironmentPredicate != nil {
		input.EnvironmentPredicate = planModel.EnvironmentPredicate.ToCreateInput()
	}

	data, err := r.client.CreateCheckToolUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_tool_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckToolUsageResourceModel(ctx, *data, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a check tool usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckToolUsageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check tool usage, got error: %s", err))
		return
	}
	stateModel := NewCheckToolUsageResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckToolUsageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckToolUsageUpdateInput{
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

	input.ToolCategory = opslevel.RefOf(opslevel.ToolCategory(planModel.ToolCategory.ValueString()))
	if planModel.ToolNamePredicate != nil {
		input.ToolNamePredicate = planModel.ToolNamePredicate.ToUpdateInput()
	} else {
		input.ToolNamePredicate = &opslevel.PredicateUpdateInput{}
	}
	if planModel.ToolUrlPredicate != nil {
		input.ToolUrlPredicate = planModel.ToolUrlPredicate.ToUpdateInput()
	} else {
		input.ToolUrlPredicate = &opslevel.PredicateUpdateInput{}
	}
	if planModel.EnvironmentPredicate != nil {
		input.EnvironmentPredicate = planModel.EnvironmentPredicate.ToUpdateInput()
	} else {
		input.EnvironmentPredicate = &opslevel.PredicateUpdateInput{}
	}

	data, err := r.client.UpdateCheckToolUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_tool_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckToolUsageResourceModel(ctx, *data, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a check tool usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckToolUsageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check tool usage, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check tool usage resource")
}

func (r *CheckToolUsageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
