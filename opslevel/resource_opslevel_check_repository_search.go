package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure      = &CheckRepositorySearchResource{}
	_ resource.ResourceWithImportState    = &CheckRepositorySearchResource{}
	_ resource.ResourceWithValidateConfig = &CheckRepositorySearchResource{}
)

func NewCheckRepositorySearchResource() resource.Resource {
	return &CheckRepositorySearchResource{}
}

// CheckRepositorySearchResource defines the resource implementation.
type CheckRepositorySearchResource struct {
	CommonResourceClient
}

type CheckRepositorySearchResourceModel struct {
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

	FileExtensions        types.Set    `tfsdk:"file_extensions"`
	FileContentsPredicate types.Object `tfsdk:"file_contents_predicate"`
}

func NewCheckRepositorySearchResourceModel(ctx context.Context, check opslevel.Check, planModel CheckRepositorySearchResourceModel) (CheckRepositorySearchResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var stateModel CheckRepositorySearchResourceModel

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

	if planModel.FileExtensions.IsNull() {
		stateModel.FileExtensions = types.SetNull(types.StringType)
	} else {
		stateModel.FileExtensions, diags = types.SetValueFrom(ctx, types.StringType, check.RepositorySearchCheckFragment.FileExtensions)
		if diags != nil && diags.HasError() {
			return CheckRepositorySearchResourceModel{}, diags
		}
	}

	predicate := check.RepositorySearchCheckFragment.FileContentsPredicate
	predicateAttrValues := map[string]attr.Value{
		"type":  types.StringValue(string(predicate.Type)),
		"value": OptionalStringValue(predicate.Value),
	}
	stateModel.FileContentsPredicate = types.ObjectValueMust(predicateType, predicateAttrValues)

	return stateModel, diags
}

func (r *CheckRepositorySearchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_repository_search"
}

func (r *CheckRepositorySearchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	predicateSchema := PredicateSchema()
	predicateSchema.Optional = false
	predicateSchema.Required = true

	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Repository Search Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"file_extensions": schema.SetAttribute{
				Description: "Restrict the search to files of given extensions. Extensions should contain only letters and numbers. For example: [\"py\", \"rb\"].",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"file_contents_predicate": predicateSchema,
		}),
	}
}

func (r *CheckRepositorySearchResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Description: "Check Repository File Resource",
				Attributes: getCheckBaseSchemaV0(map[string]schema.Attribute{
					"file_extensions": schema.ListAttribute{
						Description: "Restrict the search to certain file paths.",
						ElementType: types.StringType,
						Optional:    true,
					},
					"id": schema.StringAttribute{
						Description: "The ID of this resource.",
						Computed:    true,
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
				upgradedStateModel := CheckRepositorySearchResourceModel{}
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
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("file_extensions"), &upgradedStateModel.FileExtensions)...)
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

func (r *CheckRepositorySearchResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	fileContentsPredicate := types.ObjectNull(predicateType)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("file_contents_predicate"), &fileContentsPredicate)...)
	if resp.Diagnostics.HasError() || fileContentsPredicate.IsNull() || fileContentsPredicate.IsUnknown() {
		return
	}

	predicateModel, diags := PredicateObjectToModel(ctx, fileContentsPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.ValueString() == "exists" || predicateModel.Type.ValueString() == "does_not_exist" {
		resp.Diagnostics.AddAttributeWarning(path.Root("file_contents_predicate"), "Invalid Attribute Configuration", "type must not be 'exists' or 'does_not_exist'")
	}
	if err := predicateModel.Validate(); err != nil {
		resp.Diagnostics.AddAttributeWarning(path.Root("file_contents_predicate"), "Invalid Attribute Configuration", err.Error())
	}
}

func (r *CheckRepositorySearchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckRepositorySearchResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	predicateModel, diags := PredicateObjectToModel(ctx, planModel.FileContentsPredicate)
	resp.Diagnostics.Append(diags...)
	input := opslevel.CheckRepositorySearchCreateInput{
		CategoryId:            asID(planModel.Category),
		Enabled:               nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:              nullableID(planModel.Filter.ValueStringPointer()),
		LevelId:               asID(planModel.Level),
		Name:                  planModel.Name.ValueString(),
		Notes:                 opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:               nullableID(planModel.Owner.ValueStringPointer()),
		FileContentsPredicate: *predicateModel.ToCreateInput(),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	fileExtensions, diags := SetValueToStringSlice(ctx, planModel.FileExtensions)
	resp.Diagnostics.Append(diags...)
	input.FileExtensions = opslevel.RefOf(fileExtensions)

	data, err := r.client.CreateCheckRepositorySearch(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_repository_search, got error: %s", err))
		return
	}

	stateModel, diags := NewCheckRepositorySearchResourceModel(ctx, *data, planModel)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created a check repository search resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositorySearchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckRepositorySearchResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check repository search, got error: %s", err))
		return
	}
	verifiedStateModel, diags := NewCheckRepositorySearchResourceModel(ctx, *data, stateModel)
	resp.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckRepositorySearchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckRepositorySearchResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckRepositorySearchUpdateInput{
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

	fileExtensions, diags := SetValueToStringSlice(ctx, planModel.FileExtensions)
	resp.Diagnostics.Append(diags...)
	input.FileExtensions = opslevel.RefOf(fileExtensions)

	predicateModel, diags := PredicateObjectToModel(ctx, planModel.FileContentsPredicate)
	resp.Diagnostics.Append(diags...)
	input.FileContentsPredicate = predicateModel.ToUpdateInput()

	data, err := r.client.UpdateCheckRepositorySearch(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_repository_search, got error: %s", err))
		return
	}

	stateModel, diags := NewCheckRepositorySearchResourceModel(ctx, *data, planModel)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "updated a check repository search resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRepositorySearchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckRepositorySearchResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check repository search, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check repository search resource")
}

func (r *CheckRepositorySearchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
