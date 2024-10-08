package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.ResourceWithConfigure      = &CheckRepositoryFileResource{}
	_ resource.ResourceWithImportState    = &CheckRepositoryFileResource{}
	_ resource.ResourceWithValidateConfig = &CheckRepositoryFileResource{}
)

func NewCheckRepositoryFileResource() resource.Resource {
	return &CheckRepositoryFileResource{}
}

// CheckRepositoryFileResource defines the resource implementation.
type CheckRepositoryFileResource struct {
	CommonResourceClient
}

type CheckRepositoryFileResourceModel struct {
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

	DirectorySearch       types.Bool   `tfsdk:"directory_search"`
	Filepaths             types.List   `tfsdk:"filepaths"`
	FileContentsPredicate types.Object `tfsdk:"file_contents_predicate"`
	UseAbsoluteRoot       types.Bool   `tfsdk:"use_absolute_root"`
}

func NewCheckRepositoryFileResourceModel(ctx context.Context, check opslevel.Check, planModel CheckRepositoryFileResourceModel) CheckRepositoryFileResourceModel {
	var stateModel CheckRepositoryFileResourceModel

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

	stateModel.DirectorySearch = RequiredBoolValue(check.RepositoryFileCheckFragment.DirectorySearch)
	stateModel.Filepaths = OptionalStringListValue(check.RepositoryFileCheckFragment.Filepaths)

	if check.RepositoryFileCheckFragment.FileContentsPredicate == nil {
		stateModel.FileContentsPredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *check.RepositoryFileCheckFragment.FileContentsPredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.FileContentsPredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}
	stateModel.UseAbsoluteRoot = RequiredBoolValue(check.RepositoryFileCheckFragment.UseAbsoluteRoot)

	return stateModel
}

func (r *CheckRepositoryFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_repository_file"
}

func (r *CheckRepositoryFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Repository File Resource",

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
			"use_absolute_root": schema.BoolAttribute{
				Description: "Whether the checks looks at the absolute root of a repo or the relative root (the directory specified when attached a repo to a service).",
				Required:    true,
			},
		}),
	}
}

func (r *CheckRepositoryFileResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Description: "Check Repository File Resource",
				Attributes: getCheckBaseSchemaV0(map[string]schema.Attribute{
					"directory_search": schema.BoolAttribute{
						Description: "Whether the check looks for the existence of a directory instead of a file.",
						Required:    true,
					},
					"filepaths": schema.ListAttribute{
						Description: "Restrict the search to certain file paths.",
						ElementType: types.StringType,
						Required:    true,
					},
					"id": schema.StringAttribute{
						Description: "The ID of this resource.",
						Computed:    true,
					},
					"use_absolute_root": schema.BoolAttribute{
						Description: "Whether the checks looks at the absolute root of a repo or the relative root (the directory specified when attached a repo to a service).",
						Required:    true,
					},
				}),
				Blocks: map[string]schema.Block{
					"file_contents_predicate": schema.ListNestedBlock{
						NestedObject: predicateSchemaV0,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var diags diag.Diagnostics
				upgradedStateModel := CheckRepositoryFileResourceModel{}
				fileContentsPredicateList := types.ListNull(types.ObjectType{AttrTypes: predicateType})

				// base check attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("category"), &upgradedStateModel.Category)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("enable_on"), &upgradedStateModel.EnableOn)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("enabled"), &upgradedStateModel.Enabled)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("filter"), &upgradedStateModel.Filter)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &upgradedStateModel.Id)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("level"), &upgradedStateModel.Level)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &upgradedStateModel.Name)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("notes"), &upgradedStateModel.Notes)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("owner"), &upgradedStateModel.Owner)...)

				// repository file specific attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("directory_search"), &upgradedStateModel.DirectorySearch)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("filepaths"), &upgradedStateModel.Filepaths)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("use_absolute_root"), &upgradedStateModel.UseAbsoluteRoot)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("file_contents_predicate"), &fileContentsPredicateList)...)
				if len(fileContentsPredicateList.Elements()) == 1 {
					fileContentsPredicate := fileContentsPredicateList.Elements()[0]
					upgradedStateModel.FileContentsPredicate, diags = types.ObjectValueFrom(ctx, predicateType, fileContentsPredicate)
					resp.Diagnostics.Append(diags...)
				} else {
					upgradedStateModel.FileContentsPredicate = types.ObjectNull(predicateType)
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateModel)...)
			},
		},
	}
}

func (r *CheckRepositoryFileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	fileContentsPredicate := types.ObjectNull(predicateType)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("file_contents_predicate"), &fileContentsPredicate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicateModel, diags := PredicateObjectToModel(ctx, fileContentsPredicate)
	resp.Diagnostics.Append(diags...)
	if err := predicateModel.Validate(); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("file_contents_predicate"), "Invalid Attribute Configuration", err.Error())
	}
}

func (r *CheckRepositoryFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckRepositoryFileResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckRepositoryFileCreateInput{
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

	input.DirectorySearch = planModel.DirectorySearch.ValueBoolPointer()
	resp.Diagnostics.Append(planModel.Filepaths.ElementsAs(ctx, &input.FilePaths, false)...)

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.FileContentsPredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.FileContentsPredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("file_contents_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}
	input.UseAbsoluteRoot = planModel.UseAbsoluteRoot.ValueBoolPointer()

	data, err := r.client.CreateCheckRepositoryFile(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_repository_file, got error: %s", err))
		return
	}

	stateModel := NewCheckRepositoryFileResourceModel(ctx, *data, planModel)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created a check repository file resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositoryFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckRepositoryFileResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddWarning("State drift", stateResourceMissingMessage("opslevel_check_repository_file"))
		resp.State.RemoveResource(ctx)
		return
	}
	stateModel := NewCheckRepositoryFileResourceModel(ctx, *data, planModel)
	stateModel.EnableOn = planModel.EnableOn

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositoryFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckRepositoryFileResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckRepositoryFileUpdateInput{
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

	input.DirectorySearch = planModel.DirectorySearch.ValueBoolPointer()
	resp.Diagnostics.Append(planModel.Filepaths.ElementsAs(ctx, &input.FilePaths, false)...)

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.FileContentsPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.FileContentsPredicate = &opslevel.PredicateUpdateInput{}
	} else if err := predicateModel.Validate(); err == nil {
		input.FileContentsPredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("environment_predicate"), "Invalid Attribute Configuration", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	input.UseAbsoluteRoot = planModel.UseAbsoluteRoot.ValueBoolPointer()

	data, err := r.client.UpdateCheckRepositoryFile(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_repository_file, got error: %s", err))
		return
	}

	stateModel := NewCheckRepositoryFileResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check repository file resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositoryFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckRepositoryFileResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddWarning("State drift", stateResourceMissingMessage("opslevel_check_repository_file"))
		return
	}
	tflog.Trace(ctx, "deleted a check repository file resource")
}

func (r *CheckRepositoryFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
