package opslevel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opslevel/opslevel-go/v2026"
)

type IntegrationEtlDefinitionModel struct {
	ExtractDefinition   types.String `tfsdk:"extract_definition"`
	TransformDefinition types.String `tfsdk:"transform_definition"`
}

func integrationEtlDefinitionAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"extract_definition":   types.StringType,
		"transform_definition": types.StringType,
	}
}

// omittedBehaviour describes what happens when the block is omitted at create
// time. It differs by integration: Kubernetes seeds default definitions, Custom
// seeds nothing, so the two resources cannot share one description.
func integrationEtlDefinitionSchemaAttribute(omittedBehaviour string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "The ETL definitions used to import data from the integration. " + omittedBehaviour + " The API manages the two definitions as a unit, so both must be set together.",
		Optional:    true,
		Computed:    true,
		// UseNonNullStateForUnknown, not UseStateForUnknown: copying a null prior
		// state in would leave a known-null plan that Update then contradicts with
		// a known object, tripping "inconsistent result after apply".
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseNonNullStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"extract_definition": schema.StringAttribute{
				Description: "The YAML definition for extracting data from inbound payloads.",
				Required:    true,
				Validators: []validator.String{
					NonEmptyYAML(),
				},
				PlanModifiers: []planmodifier.String{
					UseStateForEquivalentYAML(),
				},
			},
			"transform_definition": schema.StringAttribute{
				Description: "The YAML definition for transforming extracted data to OpsLevel resources.",
				Required:    true,
				Validators: []validator.String{
					NonEmptyYAML(),
				},
				PlanModifiers: []planmodifier.String{
					UseStateForEquivalentYAML(),
				},
			},
		},
	}
}

func reconcileIntegrationEtlDefinition(ctx context.Context, apiExtractDefinition, apiTransformDefinition string, givenDefinition types.Object, diags *diag.Diagnostics) types.Object {
	etlModel := IntegrationEtlDefinitionModel{
		ExtractDefinition:   RequiredStringValue(apiExtractDefinition),
		TransformDefinition: RequiredStringValue(apiTransformDefinition),
	}
	if !givenDefinition.IsNull() && !givenDefinition.IsUnknown() {
		var givenEtlModel IntegrationEtlDefinitionModel
		diags.Append(givenDefinition.As(ctx, &givenEtlModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
		if yamlEquivalent(givenEtlModel.ExtractDefinition.ValueString(), apiExtractDefinition) {
			etlModel.ExtractDefinition = givenEtlModel.ExtractDefinition
		}
		if yamlEquivalent(givenEtlModel.TransformDefinition.ValueString(), apiTransformDefinition) {
			etlModel.TransformDefinition = givenEtlModel.TransformDefinition
		}
	}

	etlObject, objDiags := types.ObjectValueFrom(ctx, integrationEtlDefinitionAttrs(), etlModel)
	diags.Append(objDiags...)

	return etlObject
}

func integrationEtlDefinitionInput(ctx context.Context, etlDefinition types.Object, diags *diag.Diagnostics) (*opslevel.YAML, *opslevel.YAML) {
	if etlDefinition.IsNull() || etlDefinition.IsUnknown() {
		return nil, nil
	}

	var etlModel IntegrationEtlDefinitionModel
	diags.Append(etlDefinition.As(ctx, &etlModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	if diags.HasError() {
		return nil, nil
	}

	// Omit, never error. This runs against the plan, and for an Optional+Computed
	// attribute with a null config Terraform proposes the prior state - so an
	// integration created without definitions legitimately arrives here as an empty
	// pair, and erroring would permanently block every later update. Config is
	// already covered by the NonEmptyYAML validators. Omitting also protects the
	// half-empty case: sending "" would clear the stored definition, because
	// YamlType coerces it to nil and the interactors assign whenever provided.
	if etlModel.ExtractDefinition.ValueString() == "" || etlModel.TransformDefinition.ValueString() == "" {
		return nil, nil
	}

	extractDefinition := opslevel.YAML(etlModel.ExtractDefinition.ValueString())
	transformDefinition := opslevel.YAML(etlModel.TransformDefinition.ValueString())
	return &extractDefinition, &transformDefinition
}
