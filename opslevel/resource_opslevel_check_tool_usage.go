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

func NewCheckToolUsageResourceModel(ctx context.Context, check opslevel.Check) CheckToolUsageResourceModel {
	var model CheckToolUsageResourceModel

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

	model.ToolCategory = types.StringValue(string(check.ToolCategory))
	model.ToolNamePredicate = NewPredicateModel(*check.ToolNamePredicate)
	model.ToolUrlPredicate = NewPredicateModel(*check.ToolUrlPredicate)
	model.EnvironmentPredicate = NewPredicateModel(*check.EnvironmentPredicate)

	return model
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

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckToolUsageCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	input.ToolCategory = opslevel.ToolCategory(planModel.ToolCategory.ValueString())
	if planModel.ToolNamePredicate != nil {
		input.ToolNamePredicate = &opslevel.PredicateInput{
			Type:  opslevel.PredicateTypeEnum(planModel.ToolNamePredicate.Type.String()),
			Value: opslevel.RefOf(planModel.ToolNamePredicate.Value.String()),
		}
	}
	if planModel.ToolUrlPredicate != nil {
		input.ToolUrlPredicate = &opslevel.PredicateInput{
			Type:  opslevel.PredicateTypeEnum(planModel.ToolUrlPredicate.Type.String()),
			Value: opslevel.RefOf(planModel.ToolUrlPredicate.Value.String()),
		}
	}
	if planModel.EnvironmentPredicate != nil {
		input.EnvironmentPredicate = &opslevel.PredicateInput{
			Type:  opslevel.PredicateTypeEnum(planModel.EnvironmentPredicate.Type.String()),
			Value: opslevel.RefOf(planModel.EnvironmentPredicate.Value.String()),
		}
	}

	data, err := r.client.CreateCheckToolUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_tool_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckToolUsageResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
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
	stateModel := NewCheckToolUsageResourceModel(ctx, *data)

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

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
		return
	}
	input := opslevel.CheckToolUsageUpdateInput{
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

	input.ToolCategory = opslevel.RefOf(opslevel.ToolCategory(planModel.ToolCategory.ValueString()))
	if planModel.ToolNamePredicate != nil {
		input.ToolNamePredicate = &opslevel.PredicateUpdateInput{
			Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(planModel.ToolNamePredicate.Type.String())),
			Value: opslevel.RefOf(planModel.ToolNamePredicate.Value.String()),
		}
	}
	if planModel.ToolUrlPredicate != nil {
		input.ToolUrlPredicate = &opslevel.PredicateUpdateInput{
			Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(planModel.ToolUrlPredicate.Type.String())),
			Value: opslevel.RefOf(planModel.ToolUrlPredicate.Value.String()),
		}
	}
	if planModel.EnvironmentPredicate != nil {
		input.EnvironmentPredicate = &opslevel.PredicateUpdateInput{
			Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(planModel.EnvironmentPredicate.Type.String())),
			Value: opslevel.RefOf(planModel.EnvironmentPredicate.Value.String()),
		}
	}

	data, err := r.client.UpdateCheckToolUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_tool_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckToolUsageResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
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
