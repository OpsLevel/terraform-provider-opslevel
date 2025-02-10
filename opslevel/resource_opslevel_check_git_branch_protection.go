package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
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
	planModel := read[CheckCodeBaseResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckGitBranchProtectionCreateInput{
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

	data, err := r.client.CreateCheckGitBranchProtection(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_git_branch_protection, got error: %s", err))
		return
	}

	stateModel := NewCheckCodeBaseResourceModel(*data, planModel)

	tflog.Trace(ctx, "created a check git branch protection resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckGitBranchProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	planModel := read[CheckCodeBaseResourceModel](ctx, &resp.Diagnostics, req.State)
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
	verifiedStateModel := NewCheckCodeBaseResourceModel(*data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckGitBranchProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckCodeBaseResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckGitBranchProtectionUpdateInput{
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

	data, err := r.client.UpdateCheckGitBranchProtection(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_git_branch_protection, got error: %s", err))
		return
	}

	stateModel := NewCheckCodeBaseResourceModel(*data, planModel)

	tflog.Trace(ctx, "updated a check git branch protection resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckGitBranchProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckCodeBaseResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check git branch protection, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check git branch protection resource")
}

func (r *CheckGitBranchProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
