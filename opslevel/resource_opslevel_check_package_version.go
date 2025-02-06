package opslevel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.ResourceWithConfigure   = &CheckPackageVersionResource{}
	_ resource.ResourceWithImportState = &CheckPackageVersionResource{}
)

func NewCheckPackageVersionResource() resource.Resource {
	return &CheckPackageVersionResource{}
}

// CheckPackageVersionResource defines the resource implementation.
type CheckPackageVersionResource struct {
	CommonResourceClient
}

type CheckPackageVersionResourceModel struct {
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

	MissingPackageResult       types.String `tfsdk:"missing_package_result"`
	PackageConstraint          types.String `tfsdk:"package_constraint"`
	PackageManager             types.String `tfsdk:"package_manager"`
	PackageName                types.String `tfsdk:"package_name"`
	PackageNameIsRegex         types.Bool   `tfsdk:"package_name_is_regex"`
	VersionConstraintPredicate types.Object `tfsdk:"version_constraint_predicate"`
}

func ParsePredicate(predicate *opslevel.Predicate) types.Object {
	if predicate == nil {
		return types.ObjectNull(predicateType)
	}
	predicateAttrValues := map[string]attr.Value{
		"type":  types.StringValue(string(predicate.Type)),
		"value": OptionalStringValue(predicate.Value),
	}
	return types.ObjectValueMust(predicateType, predicateAttrValues)
}

func NewCheckPackageVersionResourceModel(ctx context.Context, check opslevel.Check, planModel CheckPackageVersionResourceModel) CheckPackageVersionResourceModel {
	var stateModel CheckPackageVersionResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
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

	if check.MissingPackageResult != nil {
		stateModel.MissingPackageResult = OptionalStringValue(string(*check.MissingPackageResult))
	}
	stateModel.PackageConstraint = RequiredStringValue(string(check.PackageConstraint))
	stateModel.PackageManager = RequiredStringValue(string(check.PackageManager))
	stateModel.PackageName = RequiredStringValue(check.PackageName)
	if !planModel.PackageNameIsRegex.IsNull() {
		stateModel.PackageNameIsRegex = OptionalBoolValue(&check.PackageNameIsRegex)
	}
	stateModel.VersionConstraintPredicate = ParsePredicate(check.VersionConstraintPredicate)

	return stateModel
}

func (r *CheckPackageVersionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_package_version"
}

func (r *CheckPackageVersionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check PackageVersion Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"missing_package_result": schema.StringAttribute{
				Description: "The check result if the package isn't being used by a service. (Optional.)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllCheckResultStatusEnum...),
				},
			},
			"package_constraint": schema.StringAttribute{
				Description: "The package constraint the service is to be checked for. (Required.)",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllPackageConstraintEnum...),
				},
			},
			"package_manager": schema.StringAttribute{
				Description: "The package manager (ecosystem) this package relates to. (Required.)",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllPackageManagerEnum...),
				},
			},
			"package_name": schema.StringAttribute{
				Description: "The name of the package to be checked. (Required.)",
				Required:    true,
			},
			"package_name_is_regex": schema.BoolAttribute{
				Description: "Whether or not the value in the package name field is a regular expression. (Optional.)",
				Optional:    true,
			},
			"version_constraint_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckPackageVersionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	packageVersionPossiblePredicateTypes := []opslevel.PredicateTypeEnum{
		opslevel.PredicateTypeEnumDoesNotMatchRegex,
		opslevel.PredicateTypeEnumMatchesRegex,
		opslevel.PredicateTypeEnumSatisfiesVersionConstraint,
	}

	configModel := read[CheckPackageVersionResourceModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	if configModel.PackageConstraint.ValueString() == string(opslevel.PackageConstraintEnumMatchesVersion) {
		if configModel.MissingPackageResult.IsNull() && !configModel.MissingPackageResult.IsUnknown() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("missing_package_result"),
				"Invalid Configuration",
				"missing_package_result is required when package_constraint is 'matches_version'",
			)
		}
		if configModel.VersionConstraintPredicate.IsNull() && !configModel.VersionConstraintPredicate.IsUnknown() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("version_constraint_predicate"),
				"Invalid Configuration",
				"version_constraint_predicate is required when package_constraint is 'matches_version'",
			)
		}
		if !configModel.VersionConstraintPredicate.IsNull() {
			predicateModel, diags := PredicateObjectToModel(ctx, configModel.VersionConstraintPredicate)
			resp.Diagnostics.Append(diags...)
			if !slices.Contains(packageVersionPossiblePredicateTypes, opslevel.PredicateTypeEnum(predicateModel.Type.ValueString())) {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("version_constraint_predicate"),
					"Invalid Configuration",
					fmt.Sprintf("version_constraint_predicate must be one of %v", packageVersionPossiblePredicateTypes),
				)
			} else {
				if err := predicateModel.Validate(); err != nil {
					resp.Diagnostics.AddAttributeWarning(path.Root("version_constraint_predicate"), "Invalid Configuration", err.Error())
				}
			}
		}
	} else {
		if !configModel.MissingPackageResult.IsUnknown() && !configModel.MissingPackageResult.IsNull() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("missing_package_result"),
				"Invalid Configuration",
				"missing_package_result is only valid when package_constraint is 'matches_version'",
			)
		}
		if !configModel.VersionConstraintPredicate.IsUnknown() && !configModel.VersionConstraintPredicate.IsNull() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("version_constraint_predicate"),
				"Invalid Configuration",
				"version_constraint_predicate is only valid when package_constraint is 'matches_version'",
			)
		}
	}
}

func (r *CheckPackageVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckPackageVersionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckPackageVersionCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),

		PackageConstraint: opslevel.PackageConstraintEnum(planModel.PackageConstraint.ValueString()),
		PackageManager:    opslevel.PackageManagerEnum(planModel.PackageManager.ValueString()),
		PackageName:       planModel.PackageName.ValueString(),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	if !planModel.MissingPackageResult.IsNull() {
		input.MissingPackageResult = asEnum[opslevel.CheckResultStatusEnum](planModel.MissingPackageResult.ValueString())
	}
	if !planModel.PackageNameIsRegex.IsNull() {
		input.PackageNameIsRegex = nullable(planModel.PackageNameIsRegex.ValueBoolPointer())
	}
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.VersionConstraintPredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.VersionConstraintPredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("version_constraint_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}

	data, err := r.client.CreateCheckPackageVersion(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check package version, got error: %s", err))
		return
	}

	stateModel := NewCheckPackageVersionResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check package_version resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckPackageVersionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckPackageVersionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check package_version, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckPackageVersionResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform stateModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckPackageVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckPackageVersionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[CheckPackageVersionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckPackageVersionUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Id:         asID(planModel.Id),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),

		PackageConstraint: asEnum[opslevel.PackageConstraintEnum](planModel.PackageConstraint.ValueString()),
		PackageManager:    asEnum[opslevel.PackageManagerEnum](planModel.PackageManager.ValueString()),
		PackageName:       opslevel.RefOf(planModel.PackageName.ValueString()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	if planModel.MissingPackageResult.IsNull() {
		input.MissingPackageResult = (*opslevel.CheckResultStatusEnum)(nil)
	} else {
		input.MissingPackageResult = asEnum[opslevel.CheckResultStatusEnum](planModel.MissingPackageResult.ValueString())
	}

	if !planModel.PackageNameIsRegex.IsNull() {
		input.PackageNameIsRegex = nullable(planModel.PackageNameIsRegex.ValueBoolPointer())
	} else if !stateModel.PackageNameIsRegex.IsNull() { // Then Unset
		input.PackageNameIsRegex = opslevel.RefOf(false)
	}

	predicateModel, diags := PredicateObjectToModel(ctx, planModel.VersionConstraintPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.VersionConstraintPredicate = &opslevel.PredicateUpdateInput{}
	} else if err := predicateModel.Validate(); err == nil {
		input.VersionConstraintPredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("version_constraint_predicate"), "Invalid Attribute Configuration", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdateCheckPackageVersion(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check package version, got error: %s", err))
		return
	}

	validatedModel := NewCheckPackageVersionResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check package_version resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &validatedModel)...)
}

func (r *CheckPackageVersionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckPackageVersionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check package_version, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check package_version resource")
}

func (r *CheckPackageVersionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
