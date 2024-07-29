package opslevel

import (
	"context"
	"fmt"
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

	// TODO: Unique Fields Here
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
		stateModel.MissingPackageResult = RequiredStringValue(string(*check.MissingPackageResult))
	}
	stateModel.PackageConstraint = RequiredStringValue(string(check.PackageConstraint))
	stateModel.PackageManager = RequiredStringValue(string(check.PackageManager))
	stateModel.PackageName = RequiredStringValue(check.PackageName)
	stateModel.PackageNameIsRegex = OptionalBoolValue(&check.PackageNameIsRegex)
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

func (r *CheckPackageVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckPackageVersionResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckPackageVersionCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),

		PackageConstraint: opslevel.PackageConstraintEnum(planModel.PackageConstraint.ValueString()),
		PackageManager:    opslevel.PackageManagerEnum(planModel.PackageManager.ValueString()),
		PackageName:       planModel.PackageName.ValueString(),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = &iso8601.Time{Time: enabledOn}
	}

	if !planModel.MissingPackageResult.IsNull() {
		input.MissingPackageResult = opslevel.RefOf(opslevel.CheckResultStatusEnum(planModel.MissingPackageResult.ValueString()))
	}
	if !planModel.PackageNameIsRegex.IsNull() {
		input.PackageNameIsRegex = planModel.PackageNameIsRegex.ValueBoolPointer()
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
	var planModel CheckPackageVersionResourceModel

	// Read Terraform prior stateModel data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check package_version, got error: %s", err))
		return
	}
	stateModel := NewCheckPackageVersionResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform stateModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckPackageVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel, stateModel CheckPackageVersionResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckPackageVersionUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Id:         asID(planModel.Id),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.RefOf(planModel.Notes.ValueString()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),

		PackageConstraint: opslevel.RefOf(opslevel.PackageConstraintEnum(planModel.PackageConstraint.ValueString())),
		PackageManager:    opslevel.RefOf(opslevel.PackageManagerEnum(planModel.PackageManager.ValueString())),
		PackageName:       opslevel.RefOf(planModel.PackageName.ValueString()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = &iso8601.Time{Time: enabledOn}
	}

	if !planModel.MissingPackageResult.IsNull() {
		input.MissingPackageResult = opslevel.RefOf(opslevel.CheckResultStatusEnum(planModel.MissingPackageResult.ValueString()))
	} else if !stateModel.MissingPackageResult.IsNull() { // Then Unset
		input.MissingPackageResult = opslevel.RefOf(opslevel.CheckResultStatusEnum(""))
	}

	if !planModel.PackageNameIsRegex.IsNull() {
		input.PackageNameIsRegex = planModel.PackageNameIsRegex.ValueBoolPointer()
	} else if !stateModel.PackageNameIsRegex.IsNull() { // Then Unset
		var v *bool
		input.PackageNameIsRegex = v
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
	var planModel CheckPackageVersionResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check package_version, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check package_version resource")
}

func (r *CheckPackageVersionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
