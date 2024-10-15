package opslevel

// TODO: write tests

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

var (
	_ resource.ResourceWithConfigure   = &CheckGitBranchProtectionResource{}
	_ resource.ResourceWithImportState = &CheckGitBranchProtectionResource{}
)

func NewCheckGitBranchProtectionResource() resource.Resource {
	return &CheckGitBranchProtectionResource{}
}

// CheckGitBranchProtectionResource defines the resource implementation.
type CheckGitBranchProtectionResource struct {
	CommonResourceClient
}

type CheckGitBranchProtectionResourceModel struct {
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
}

func NewCheckGitBranchProtectionResourceModel(ctx context.Context, check opslevel.Check, planModel CheckGitBranchProtectionResourceModel) CheckGitBranchProtectionResourceModel {
	var stateModel CheckGitBranchProtectionResourceModel

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

	return stateModel
}

func (r *CheckGitBranchProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_git_branch_protection"
}

func (r *CheckGitBranchProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Git Branch Protection Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{}),
	}
}

func (r *CheckGitBranchProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckGitBranchProtectionResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckGitBranchProtectionCreateInput{
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

	data, err := r.client.CreateCheckGitBranchProtection(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_git_branch_protection, got error: %s", err))
		return
	}

	stateModel := NewCheckGitBranchProtectionResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check git branch protection resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckGitBranchProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckGitBranchProtectionResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check git branch protection, got error: %s", err))
		return
	}
	stateModel := NewCheckGitBranchProtectionResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckGitBranchProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckGitBranchProtectionResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckGitBranchProtectionUpdateInput{
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

	data, err := r.client.UpdateCheckGitBranchProtection(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_git_branch_protection, got error: %s", err))
		return
	}

	stateModel := NewCheckGitBranchProtectionResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check git branch protection resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckGitBranchProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckGitBranchProtectionResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check git branch protection, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check git branch protection resource")
}

func (r *CheckGitBranchProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
