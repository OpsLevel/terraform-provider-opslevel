package opslevel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/opslevel/terraform-provider-opslevel/internal"
)

var RelationshipDefinitionDataSourceSchema = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "The display name of the relationship definition.",
		Computed:            true,
	},
	"alias": schema.StringAttribute{
		MarkdownDescription: "The unique identifier of the relationship.",
		Computed:            true,
	},
	"description": schema.StringAttribute{
		MarkdownDescription: "The description of the relationship definition.",
		Computed:            true,
	},
	"component_type": schema.StringAttribute{
		MarkdownDescription: "The component type that the relationship belongs to.",
		Computed:            true,
	},
	"allowed_categories": schema.ListAttribute{
		MarkdownDescription: "The categories of resources that can be selected for this relationship definition. Can include any component category alias on your account.",
		Computed:            true,
		ElementType:         types.StringType,
	},
	"allowed_types": schema.ListAttribute{
		MarkdownDescription: "The types of resources that can be selected for this relationship definition. Can include any component type alias on your account or 'team'.",
		Computed:            true,
		ElementType:         types.StringType,
	},
	"management_rules": schema.ListNestedAttribute{
		MarkdownDescription: "Rules that automatically manage relationships based on property matching conditions.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"operator": schema.StringAttribute{
					MarkdownDescription: "The condition operator for this rule. Either EQUALS or ARRAY_CONTAINS",
					Computed:            true,
				},
				"source_property": schema.StringAttribute{
					MarkdownDescription: "The property on the source component to evaluate.",
					Computed:            true,
				},
				"target_category": schema.StringAttribute{
					MarkdownDescription: "The category of the target resource. Either target_category or target_type must be specified, but not both.",
					Computed:            true,
				},
				"target_property": schema.StringAttribute{
					MarkdownDescription: "The property on the target resource to match against.",
					Computed:            true,
				},
				"target_type": schema.StringAttribute{
					MarkdownDescription: "The type of the target resource. Either target_category or target_type must be specified, but not both.",
					Computed:            true,
				},
			},
		},
	},
}

type RelationshipDefinitionDataSourceModel struct {
	Identifier types.String `tfsdk:"identifier"`
	RelationshipDefinitionResourceModel
}

func NewRelationshipDefinitionDataSourceSingle() datasource.DataSource {
	return &internal.TFDataSourceSingle[opslevel.RelationshipDefinitionType, RelationshipDefinitionDataSourceModel]{
		Name:        "relationship_definition",
		Description: "Relationship Definition Data Source",
		Attributes:  RelationshipDefinitionDataSourceSchema,
		ReadFn: func(ctx context.Context, client *opslevel.Client, identifier string) (opslevel.RelationshipDefinitionType, error) {
			data, err := client.GetRelationshipDefinition(identifier)
			if err != nil {
				return opslevel.RelationshipDefinitionType{}, err
			}
			return *data, nil
		},
		ToModel: func(ctx context.Context, identifier string, data opslevel.RelationshipDefinitionType) (RelationshipDefinitionDataSourceModel, error) {
			allowedCategories := make([]attr.Value, len(data.Metadata.AllowedCategories))
			for i, t := range data.Metadata.AllowedCategories {
				allowedCategories[i] = types.StringValue(t)
			}

			allowedTypes := make([]attr.Value, len(data.Metadata.AllowedTypes))
			for i, t := range data.Metadata.AllowedTypes {
				allowedTypes[i] = types.StringValue(t)
			}

			var managementRules types.List
			if len(data.ManagementRules) > 0 {
				ruleValues := make([]attr.Value, len(data.ManagementRules))
				for i, rule := range data.ManagementRules {
					ruleValues[i] = newManagementRuleValue(rule)
				}
				managementRules = types.ListValueMust(
					types.ObjectType{AttrTypes: ManagementRuleModelAttrs()},
					ruleValues,
				)
			} else {
				managementRules = types.ListNull(types.ObjectType{AttrTypes: ManagementRuleModelAttrs()})
			}

			model := RelationshipDefinitionResourceModel{
				Id:                types.StringValue(string(data.Id)),
				Name:              types.StringValue(data.Name),
				Alias:             types.StringValue(data.Alias),
				Description:       types.StringValue(data.Description),
				ComponentType:     types.StringValue(string(data.ComponentType.Id)),
				AllowedCategories: types.ListValueMust(types.StringType, allowedCategories),
				AllowedTypes:      types.ListValueMust(types.StringType, allowedTypes),
				ManagementRules:   managementRules,
			}

			return RelationshipDefinitionDataSourceModel{
				Identifier:                          types.StringValue(identifier),
				RelationshipDefinitionResourceModel: model,
			}, nil
		},
	}
}

func NewRelationshipDefinitionDataSourceMulti() datasource.DataSource {
	return &internal.TFDataSourceMulti[opslevel.RelationshipDefinitionType, RelationshipDefinitionResourceModel]{
		Name:        "relationship_definitions",
		Description: "Relationship Definition Data Source",
		Attributes:  RelationshipDefinitionDataSourceSchema,
		ReadFn: func(ctx context.Context, client *opslevel.Client) ([]opslevel.RelationshipDefinitionType, error) {
			resp, err := client.ListRelationshipDefinitions(nil)
			if resp == nil {
				return nil, err
			}
			return resp.Nodes, err
		},
		ToModel: func(ctx context.Context, data opslevel.RelationshipDefinitionType) (RelationshipDefinitionResourceModel, error) {
			allowedCategories := make([]attr.Value, len(data.Metadata.AllowedCategories))
			for i, t := range data.Metadata.AllowedCategories {
				allowedCategories[i] = types.StringValue(t)
			}

			allowedTypes := make([]attr.Value, len(data.Metadata.AllowedTypes))
			for i, t := range data.Metadata.AllowedTypes {
				allowedTypes[i] = types.StringValue(t)
			}

			var managementRules types.List
			if len(data.ManagementRules) > 0 {
				ruleValues := make([]attr.Value, len(data.ManagementRules))
				for i, rule := range data.ManagementRules {
					ruleValues[i] = newManagementRuleValue(rule)
				}
				managementRules = types.ListValueMust(
					types.ObjectType{AttrTypes: ManagementRuleModelAttrs()},
					ruleValues,
				)
			} else {
				managementRules = types.ListNull(types.ObjectType{AttrTypes: ManagementRuleModelAttrs()})
			}

			return RelationshipDefinitionResourceModel{
				Id:                types.StringValue(string(data.Id)),
				Name:              types.StringValue(data.Name),
				Alias:             types.StringValue(data.Alias),
				Description:       types.StringValue(data.Description),
				ComponentType:     types.StringValue(string(data.ComponentType.Id)),
				AllowedCategories: types.ListValueMust(types.StringType, allowedCategories),
				AllowedTypes:      types.ListValueMust(types.StringType, allowedTypes),
				ManagementRules:   managementRules,
			}, nil
		},
	}
}
