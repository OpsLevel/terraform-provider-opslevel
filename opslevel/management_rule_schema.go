// management_rule_schema.go
package opslevel

import (
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ManagementRulesResourceAttribute returns the management_rules schema for resources
func ManagementRulesResourceAttribute() resschema.ListNestedAttribute {
	return resschema.ListNestedAttribute{
		Description: "Rules that automatically manage relationships based on property matching conditions.",
		Optional:    true,
		Validators: []validator.List{
			ManagementRuleTagValidator(),
		},
		NestedObject: resschema.NestedAttributeObject{
			Attributes: managementRuleResourceAttributes(),
		},
	}
}

func managementRuleResourceAttributes() map[string]resschema.Attribute {
	return map[string]resschema.Attribute{
		"operator": resschema.StringAttribute{
			Description: "The condition operator for this rule. Either EQUALS or ARRAY_CONTAINS.",
			Required:    true,
		},
		"source_property": resschema.StringAttribute{
			Description: "The property on the source component to evaluate.",
			Required:    true,
		},
		"source_tag_key": resschema.StringAttribute{
			Description: "When source_property is 'tag', this specifies the tag key to match. Required if source_property is 'tag', must not be set otherwise.",
			Optional:    true,
		},
		"source_tag_operation": resschema.StringAttribute{
			Description: "When source_property is 'tag', this specifies the matching operation. Either 'equals' or 'starts_with'. Defaults to 'equals'. Required if source_property is 'tag', must not be set otherwise",
			Optional:    true,
		},
		"target_category": resschema.StringAttribute{
			Description: "The category of the target resource. Either target_category or target_type must be specified, but not both.",
			Optional:    true,
		},
		"target_property": resschema.StringAttribute{
			Description: "The property on the target resource to match against.",
			Required:    true,
		},
		"target_tag_key": resschema.StringAttribute{
			Description: "When target_property is 'tag', this specifies the tag key to match. Required if target_property is 'tag', must not be set otherwise.",
			Optional:    true,
		},
		"target_tag_operation": resschema.StringAttribute{
			Description: "When target_property is 'tag', this specifies the matching operation. Either 'equals' or 'starts_with'. Defaults to 'equals'. Required if target_property is 'tag', must not be set otherwise.",
			Optional:    true,
		},
		"target_type": resschema.StringAttribute{
			Description: "The type of the target resource. Either target_category or target_type must be specified, but not both.",
			Optional:    true,
		},
	}
}

// ManagementRulesDataSourceAttribute returns the management_rules schema for data sources
func ManagementRulesDataSourceAttribute() dsschema.ListNestedAttribute {
	return dsschema.ListNestedAttribute{
		Description: "Rules that automatically manage relationships based on property matching conditions.",
		Computed:    true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: managementRuleDataSourceAttributes(),
		},
	}
}

func managementRuleDataSourceAttributes() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"operator": dsschema.StringAttribute{
			Description: "The condition operator for this rule.",
			Computed:    true,
		},
		"source_property": dsschema.StringAttribute{
			Description: "The property on the source component to evaluate.",
			Computed:    true,
		},
		"source_tag_key": dsschema.StringAttribute{
			Description: "When source_property is 'tag', this specifies the tag key to match.",
			Computed:    true,
		},
		"source_tag_operation": dsschema.StringAttribute{
			Description: "When source_property is 'tag', this specifies the matching operation.",
			Computed:    true,
		},
		"target_category": dsschema.StringAttribute{
			Description: "The category of the target resource.",
			Computed:    true,
		},
		"target_property": dsschema.StringAttribute{
			Description: "The property on the target resource to match against.",
			Computed:    true,
		},
		"target_tag_key": dsschema.StringAttribute{
			Description: "When target_property is 'tag', this specifies the tag key to match.",
			Computed:    true,
		},
		"target_tag_operation": dsschema.StringAttribute{
			Description: "When target_property is 'tag', this specifies the matching operation.",
			Computed:    true,
		},
		"target_type": dsschema.StringAttribute{
			Description: "The type of the target resource.",
			Computed:    true,
		},
	}
}
