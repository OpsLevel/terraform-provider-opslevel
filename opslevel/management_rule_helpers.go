package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opslevel/opslevel-go/v2025"
)

var (
	TEAM_BUILTIN_PROPERTIES      = []string{"name", "alias", "contact", "tag"}
	USER_BUILTIN_PROPERTIES      = []string{"name", "contact", "tag"}
	COMPONENT_BUILTIN_PROPERTIES = []string{"name", "alias", "tag"}
)

type ManagementRuleModel struct {
	Operator           types.String `tfsdk:"operator"`
	SourceProperty     types.String `tfsdk:"source_property"`
	SourceTagKey       types.String `tfsdk:"source_tag_key"`
	SourceTagOperation types.String `tfsdk:"source_tag_operation"`
	TargetCategory     types.String `tfsdk:"target_category"`
	TargetProperty     types.String `tfsdk:"target_property"`
	TargetTagKey       types.String `tfsdk:"target_tag_key"`
	TargetTagOperation types.String `tfsdk:"target_tag_operation"`
	TargetType         types.String `tfsdk:"target_type"`
}

func ManagementRuleModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"operator":             types.StringType,
		"source_property":      types.StringType,
		"source_tag_key":       types.StringType,
		"source_tag_operation": types.StringType,
		"target_category":      types.StringType,
		"target_property":      types.StringType,
		"target_tag_key":       types.StringType,
		"target_tag_operation": types.StringType,
		"target_type":          types.StringType,
	}
}

func NewManagementRuleValue(rule opslevel.RelationshipDefinitionManagementRule) attr.Value {
	var targetCategory types.String
	if rule.TargetCategory != nil && !rule.TargetCategory.SetNull {
		targetCategory = types.StringValue(rule.TargetCategory.Value)
	} else {
		targetCategory = types.StringNull()
	}

	var targetType types.String
	if rule.TargetType != nil && !rule.TargetType.SetNull {
		targetType = types.StringValue(rule.TargetType.Value)
	} else {
		targetType = types.StringNull()
	}

	sourceProperty, sourceTagKey, sourceTagOp := ParsePropertyString(rule.SourceProperty)
	targetProperty, targetTagKey, targetTagOp := ParsePropertyString(rule.TargetProperty)

	return types.ObjectValueMust(
		ManagementRuleModelAttrs(),
		map[string]attr.Value{
			"operator":             types.StringValue(string(rule.Operator)),
			"source_property":      types.StringValue(sourceProperty),
			"source_tag_key":       OptionalStringValue(sourceTagKey),
			"source_tag_operation": OptionalStringValue(sourceTagOp),
			"target_category":      targetCategory,
			"target_property":      types.StringValue(targetProperty),
			"target_tag_key":       OptionalStringValue(targetTagKey),
			"target_tag_operation": OptionalStringValue(targetTagOp),
			"target_type":          targetType,
		},
	)
}

func ParseManagementRules(ctx context.Context, planRules types.List, componentTypeAlias string, diags *diag.Diagnostics) []opslevel.ManagementRuleInput {
	if planRules.IsNull() || planRules.IsUnknown() {
		return nil
	}

	var planRulesModels []ManagementRuleModel
	if err := planRules.ElementsAs(ctx, &planRulesModels, false); err != nil {
		diags.AddError("config error", fmt.Sprintf("unable to parse management_rules: %s", err))
		return nil
	}

	managementRules := make([]opslevel.ManagementRuleInput, len(planRulesModels))
	for i, rule := range planRulesModels {
		var targetTypeOrCategory string
		isType := false

		if !rule.TargetType.IsNull() && !rule.TargetType.IsUnknown() {
			targetTypeOrCategory = rule.TargetType.ValueString()
			isType = true
		} else if !rule.TargetCategory.IsNull() && !rule.TargetCategory.IsUnknown() {
			targetTypeOrCategory = rule.TargetCategory.ValueString()
			isType = false
		}

		sourcePropertyStr := BuildPropertyString(
			rule.SourceProperty.ValueString(),
			rule.SourceTagKey.ValueString(),
			rule.SourceTagOperation.ValueString(),
		)

		targetPropertyStr := BuildPropertyString(
			rule.TargetProperty.ValueString(),
			rule.TargetTagKey.ValueString(),
			rule.TargetTagOperation.ValueString(),
		)

		sourcePropertyBuiltin := IsBuiltinProperty(componentTypeAlias, rule.SourceProperty.ValueString(), true)
		targetPropertyBuiltin := IsBuiltinProperty(targetTypeOrCategory, rule.TargetProperty.ValueString(), isType)

		managementRules[i] = opslevel.ManagementRuleInput{
			Operator:              opslevel.RelationshipDefinitionManagementRuleOperator(rule.Operator.ValueString()),
			SourceProperty:        sourcePropertyStr,
			SourcePropertyBuiltin: sourcePropertyBuiltin,
			TargetProperty:        targetPropertyStr,
			TargetPropertyBuiltin: targetPropertyBuiltin,
		}

		if !rule.TargetCategory.IsNull() && !rule.TargetCategory.IsUnknown() {
			targetCategory := rule.TargetCategory.ValueString()
			managementRules[i].TargetCategory = nullable(&targetCategory)
		}

		if !rule.TargetType.IsNull() && !rule.TargetType.IsUnknown() {
			targetType := rule.TargetType.ValueString()
			managementRules[i].TargetType = nullable(&targetType)
		}
	}

	return managementRules
}

func BuildPropertyString(property, tagKey, tagOperation string) string {
	if property != "tag" {
		return property
	}

	operation := "eq"
	if tagOperation != "" {
		if tagOperation == "starts_with" {
			operation = "starts_with"
		}
	}

	return fmt.Sprintf("tag_key_%s:%s", operation, tagKey)
}

func ParsePropertyString(propertyStr string) (property, tagKey, tagOperation string) {
	if !strings.HasPrefix(propertyStr, "tag_key_") {
		return propertyStr, "", ""
	}

	property = "tag"

	remainder := strings.TrimPrefix(propertyStr, "tag_key_")

	if strings.HasPrefix(remainder, "eq:") {
		tagOperation = "equals"
		tagKey = strings.TrimPrefix(remainder, "eq:")
	} else if strings.HasPrefix(remainder, "starts_with:") {
		tagOperation = "starts_with"
		tagKey = strings.TrimPrefix(remainder, "starts_with:")
	}

	return
}

func IsBuiltinProperty(targetTypeOrCategory string, propertyName string, isType bool) bool {
	var builtinProps []string

	if isType {
		if targetTypeOrCategory == "team" {
			builtinProps = TEAM_BUILTIN_PROPERTIES
		} else if targetTypeOrCategory == "user" {
			builtinProps = USER_BUILTIN_PROPERTIES
		} else {
			builtinProps = COMPONENT_BUILTIN_PROPERTIES
		}
	} else {
		if targetTypeOrCategory == "people" {
			builtinProps = TEAM_BUILTIN_PROPERTIES
		} else {
			builtinProps = COMPONENT_BUILTIN_PROPERTIES
		}
	}

	for _, prop := range builtinProps {
		if prop == propertyName {
			return true
		}
	}
	return false
}

func ManagementRuleListValueFromResourceAndModel(resourceRules []opslevel.RelationshipDefinitionManagementRule, modelValue basetypes.ListValue) basetypes.ListValue {
	ruleValues := make([]attr.Value, len(resourceRules))
	for i, rule := range resourceRules {
		ruleValues[i] = NewManagementRuleValue(rule)
	}

	objectType := types.ObjectType{AttrTypes: ManagementRuleModelAttrs()}

	if len(resourceRules) > 0 || !modelValue.IsNull() {
		return types.ListValueMust(objectType, ruleValues)
	}

	return types.ListNull(objectType)
}
