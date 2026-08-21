package opslevel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

func integrationEtlDefinitionSchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "The ETL definitions used to import data from the integration. If omitted when the integration is created, OpsLevel's default definitions are used. The API manages the two definitions as a unit, so both must be set together.",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"extract_definition": schema.StringAttribute{
				Description: "The YAML definition for extracting data from inbound payloads.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					UseStateForEquivalentYAML(),
				},
			},
			"transform_definition": schema.StringAttribute{
				Description: "The YAML definition for transforming extracted data to OpsLevel resources.",
				Required:    true,
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

	if etlModel.ExtractDefinition.ValueString() == "" && etlModel.TransformDefinition.ValueString() == "" {
		return nil, nil
	}

	extractDefinition := opslevel.YAML(etlModel.ExtractDefinition.ValueString())
	transformDefinition := opslevel.YAML(etlModel.TransformDefinition.ValueString())
	return &extractDefinition, &transformDefinition
}
