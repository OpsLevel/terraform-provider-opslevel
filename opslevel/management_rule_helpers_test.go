package opslevel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	opslevelgo "github.com/opslevel/opslevel-go/v2026"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func tagManagementRule(sourceProperty, targetProperty string) opslevelgo.RelationshipDefinitionManagementRule {
	return opslevelgo.RelationshipDefinitionManagementRule{
		Operator:              opslevelgo.RelationshipDefinitionManagementRuleOperatorEquals,
		SourceProperty:        sourceProperty,
		SourcePropertyBuiltin: true,
		TargetCategory:        opslevelgo.NewNullableFrom("infrastructure"),
		TargetProperty:        targetProperty,
		TargetPropertyBuiltin: true,
	}
}

// managementRuleModelList returns a one-rule list as Terraform would plan it, with the given attributes
// overriding a rule that matches the alias of the source to the target's business_unit tag.
func managementRuleModelList(overrides map[string]attr.Value) basetypes.ListValue {
	attrs := map[string]attr.Value{
		"operator":             types.StringValue("EQUALS"),
		"source_property":      types.StringValue("alias"),
		"source_tag_key":       types.StringNull(),
		"source_tag_operation": types.StringNull(),
		"target_category":      types.StringValue("infrastructure"),
		"target_property":      types.StringValue("tag"),
		"target_tag_key":       types.StringValue("business_unit"),
		"target_tag_operation": types.StringNull(),
		"target_type":          types.StringNull(),
	}
	for name, value := range overrides {
		attrs[name] = value
	}
	objectType := types.ObjectType{AttrTypes: opsleveltf.ManagementRuleModelAttrs()}
	return types.ListValueMust(objectType, []attr.Value{types.ObjectValueMust(opsleveltf.ManagementRuleModelAttrs(), attrs)})
}

func ruleAttribute(list basetypes.ListValue, index int, name string) attr.Value {
	return list.Elements()[index].(basetypes.ObjectValue).Attributes()[name]
}

func TestManagementRuleListValueFromResourceAndModel_UnsetTargetTagOperationMatchesPlan(t *testing.T) {
	planRules := managementRuleModelList(nil)

	stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
		[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("alias", "tag_key_eq:business_unit")},
		planRules,
	)

	if !stateRules.Equal(planRules) {
		t.Errorf("expected state to match the plan when target_tag_operation is unset, got %s", stateRules)
	}
}

func TestManagementRuleListValueFromResourceAndModel_UnsetSourceTagOperationMatchesPlan(t *testing.T) {
	planRules := managementRuleModelList(map[string]attr.Value{
		"source_property": types.StringValue("tag"),
		"source_tag_key":  types.StringValue("business_unit"),
	})

	stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
		[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("tag_key_eq:business_unit", "tag_key_eq:business_unit")},
		planRules,
	)

	if !stateRules.Equal(planRules) {
		t.Errorf("expected state to match the plan when source_tag_operation is unset, got %s", stateRules)
	}
}

func TestManagementRuleListValueFromResourceAndModel_ExplicitTagOperationsMatchPlan(t *testing.T) {
	for _, operation := range []struct{ configured, apiProperty string }{
		{"equals", "tag_key_eq:business_unit"},
		{"starts_with", "tag_key_starts_with:business_unit"},
	} {
		planRules := managementRuleModelList(map[string]attr.Value{"target_tag_operation": types.StringValue(operation.configured)})

		stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
			[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("alias", operation.apiProperty)},
			planRules,
		)

		if !stateRules.Equal(planRules) {
			t.Errorf("expected state to match the plan for target_tag_operation %q, got %s", operation.configured, stateRules)
		}
	}
}

func TestManagementRuleListValueFromResourceAndModel_UnsetTagOperationStillShowsStartsWithDrift(t *testing.T) {
	stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
		[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("alias", "tag_key_starts_with:business_unit")},
		managementRuleModelList(nil),
	)

	if got := ruleAttribute(stateRules, 0, "target_tag_operation"); !got.Equal(types.StringValue("starts_with")) {
		t.Errorf("expected a starts_with rule changed outside Terraform to read back as starts_with, got %s", got)
	}
}

func TestManagementRuleListValueFromResourceAndModel_NoModelRuleReadsEquals(t *testing.T) {
	objectType := types.ObjectType{AttrTypes: opsleveltf.ManagementRuleModelAttrs()}

	importedRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
		[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("alias", "tag_key_eq:business_unit")},
		types.ListNull(objectType),
	)

	if got := ruleAttribute(importedRules, 0, "target_tag_operation"); !got.Equal(types.StringValue("equals")) {
		t.Errorf("expected an imported rule to read back target_tag_operation as equals, got %s", got)
	}
}

func TestManagementRuleListValueFromResourceAndModel_RuleBeyondModelReadsEquals(t *testing.T) {
	stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
		[]opslevelgo.RelationshipDefinitionManagementRule{
			tagManagementRule("alias", "tag_key_eq:business_unit"),
			tagManagementRule("alias", "tag_key_eq:team"),
		},
		managementRuleModelList(nil),
	)

	if got := ruleAttribute(stateRules, 0, "target_tag_operation"); !got.IsNull() {
		t.Errorf("expected the planned rule to keep target_tag_operation unset, got %s", got)
	}
	if got := ruleAttribute(stateRules, 1, "target_tag_operation"); !got.Equal(types.StringValue("equals")) {
		t.Errorf("expected a rule added outside Terraform to read back target_tag_operation as equals, got %s", got)
	}
}

func TestNewManagementRuleValue_ReadsEquals(t *testing.T) {
	rule := opsleveltf.NewManagementRuleValue(tagManagementRule("alias", "tag_key_eq:business_unit")).(basetypes.ObjectValue)

	if got := rule.Attributes()["target_tag_operation"]; !got.Equal(types.StringValue("equals")) {
		t.Errorf("expected data sources to report target_tag_operation as equals, got %s", got)
	}
}
