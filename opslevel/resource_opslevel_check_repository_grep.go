package opslevel

// TODO: register resource to provider
// TODO: write tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckRepositoryGrepResource{}
	_ resource.ResourceWithImportState = &CheckRepositoryGrepResource{}
)

func NewCheckRepositoryGrepResource() resource.Resource {
	return &CheckRepositoryGrepResource{}
}

// CheckRepositoryGrepResource defines the resource implementation.
type CheckRepositoryGrepResource struct {
	CommonResourceClient
}

type CheckRepositoryGrepResourceModel struct {
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

	DirectorySearch       types.Bool      `tfsdk:"directory_search"`
	Filepaths             types.List      `tfsdk:"filepaths"`
	FileContentsPredicate *PredicateModel `tfsdk:"file_contents_predicate"`
}

func NewCheckRepositoryGrepResourceModel(ctx context.Context, check opslevel.Check) (CheckRepositoryGrepResourceModel, diag.Diagnostics) {
	var model CheckRepositoryGrepResourceModel

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

	model.DirectorySearch = types.BoolValue(check.RepositoryGrepCheckFragment.DirectorySearch)
	data, diags := types.ListValueFrom(ctx, types.StringType, check.RepositoryGrepCheckFragment.Filepaths)
	model.Filepaths = data
	model.FileContentsPredicate = NewPredicateModel(*check.RepositoryGrepCheckFragment.FileContentsPredicate)

	return model, diags
}

func (r *CheckRepositoryGrepResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_repository_grep"
}

func (r *CheckRepositoryGrepResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Repository Grep Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"directory_search": schema.BoolAttribute{
				Description: "Whether the check looks for the existence of a directory instead of a file.",
				Required:    true,
			},
			"filepaths": schema.ListAttribute{
				Description: "Restrict the search to certain file paths.",
				Required:    true,
				ElementType: types.StringType,
			},
			"file_contents_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckRepositoryGrepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckRepositoryGrepResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckRepositoryGrepCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	input.DirectorySearch = planModel.DirectorySearch.ValueBoolPointer()
	resp.Diagnostics.Append(planModel.Filepaths.ElementsAs(ctx, &input.FilePaths, false)...)
	if planModel.FileContentsPredicate != nil {
		input.FileContentsPredicate = opslevel.PredicateInput{
			Type:  opslevel.PredicateTypeEnum(planModel.FileContentsPredicate.Type.String()),
			Value: opslevel.RefOf(planModel.FileContentsPredicate.Value.String()),
		}
	}

	data, err := r.client.CreateCheckRepositoryGrep(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_repository_grep, got error: %s", err))
		return
	}

	stateModel, diags := NewCheckRepositoryGrepResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	// Custom Prop Overrides Here
	stateModel.LastUpdated = timeLastUpdated()
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created a check repositorygrep resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositoryGrepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckRepositoryGrepResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check repositorygrep, got error: %s", err))
		return
	}
	stateModel, diags := NewCheckRepositoryGrepResourceModel(ctx, *data)
	resp.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositoryGrepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckRepositoryGrepResourceModel

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
	input := opslevel.CheckRepositoryGrepUpdateInput{
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

	input.DirectorySearch = planModel.DirectorySearch.ValueBoolPointer()
	resp.Diagnostics.Append(planModel.Filepaths.ElementsAs(ctx, &input.FilePaths, false)...)
	if planModel.FileContentsPredicate != nil {
		input.FileContentsPredicate = &opslevel.PredicateUpdateInput{
			Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(planModel.FileContentsPredicate.Type.String())),
			Value: opslevel.RefOf(planModel.FileContentsPredicate.Value.String()),
		}
	}

	data, err := r.client.UpdateCheckRepositoryGrep(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_repository_grep, got error: %s", err))
		return
	}

	stateModel, diags := NewCheckRepositoryGrepResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.LastUpdated = timeLastUpdated()
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "updated a check repositorygrep resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositoryGrepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckRepositoryGrepResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check repositorygrep, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check repositorygrep resource")
}

func (r *CheckRepositoryGrepResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
