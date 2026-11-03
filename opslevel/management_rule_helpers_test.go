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

func TestTagOperation_PlannedValueSurvivesTheAPIRoundTrip(t *testing.T) {
	// Terraform requires the value read back after apply to match the plan.
	for _, rule := range []struct {
		property   string
		configured types.String
	}{
		{"tag", types.StringNull()},
		{"tag", types.StringValue("equals")},
		{"tag", types.StringValue("starts_with")},
		{"alias", types.StringNull()},
	} {
		planned := opsleveltf.PlannedTagOperation(rule.configured, types.StringValue(rule.property))
		sent := opsleveltf.BuildPropertyString(rule.property, "business_unit", planned.ValueString())
		_, _, readBack := opsleveltf.ParsePropertyString(sent)

		if readBack != planned.ValueString() {
			t.Errorf("property %q configured %s: planned %s, sent %q, read back %q", rule.property, rule.configured, planned, sent, readBack)
		}
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

func TestManagementRuleListValueFromResourceAndModel_ReportsTheAPIOperation(t *testing.T) {
	objectType := types.ObjectType{AttrTypes: opsleveltf.ManagementRuleModelAttrs()}
	for _, operation := range []struct{ apiProperty, expected string }{
		{"tag_key_eq:business_unit", "equals"},
		{"tag_key_starts_with:business_unit", "starts_with"},
	} {
		// A null model is what an import starts from.
		for _, model := range []basetypes.ListValue{managementRuleModelList(nil), types.ListNull(objectType)} {
			stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
				[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("alias", operation.apiProperty)},
				model,
			)

			if got := ruleAttribute(stateRules, 0, "target_tag_operation"); !got.Equal(types.StringValue(operation.expected)) {
				t.Errorf("expected %q to read back as %q, got %s", operation.apiProperty, operation.expected, got)
			}
		}
	}
}

func TestManagementRuleListValueFromResourceAndModel_RuleWithoutTagHasNoOperation(t *testing.T) {
	stateRules := opsleveltf.ManagementRuleListValueFromResourceAndModel(
		[]opslevelgo.RelationshipDefinitionManagementRule{tagManagementRule("alias", "name")},
		managementRuleModelList(nil),
	)

	if got := ruleAttribute(stateRules, 0, "target_tag_operation"); !got.IsNull() {
		t.Errorf("expected no target_tag_operation for a rule whose target_property is not 'tag', got %s", got)
	}
}

func TestNewManagementRuleValue_ReadsEquals(t *testing.T) {
	rule := opsleveltf.NewManagementRuleValue(tagManagementRule("alias", "tag_key_eq:business_unit")).(basetypes.ObjectValue)

	if got := rule.Attributes()["target_tag_operation"]; !got.Equal(types.StringValue("equals")) {
		t.Errorf("expected data sources to report target_tag_operation as equals, got %s", got)
	}
}
