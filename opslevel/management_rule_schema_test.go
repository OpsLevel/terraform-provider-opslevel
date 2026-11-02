package opslevel_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

// ruleAttributeRejects runs the validators of one management rule attribute against value.
func ruleAttributeRejects(t *testing.T, name string, value types.String) bool {
	t.Helper()
	attribute, ok := opsleveltf.ManagementRulesResourceAttribute().NestedObject.Attributes[name].(resschema.StringAttribute)
	if !ok {
		t.Fatalf("management rule attribute %q is not a string attribute", name)
	}
	for _, v := range attribute.Validators {
		resp := &validator.StringResponse{}
		v.ValidateString(context.Background(), validator.StringRequest{Path: path.Root(name), ConfigValue: value}, resp)
		if resp.Diagnostics.HasError() {
			return true
		}
	}
	return false
}

func TestManagementRuleTagOperation_AcceptsEqualsAndStartsWith(t *testing.T) {
	for _, name := range []string{"source_tag_operation", "target_tag_operation"} {
		for _, value := range []string{"equals", "starts_with"} {
			if ruleAttributeRejects(t, name, types.StringValue(value)) {
				t.Errorf("expected %s to accept %q", name, value)
			}
		}
	}
}

func TestManagementRuleTagOperation_RejectsOtherValues(t *testing.T) {
	for _, name := range []string{"source_tag_operation", "target_tag_operation"} {
		for _, value := range []string{"EQUALS", "eq", "startswith", ""} {
			if !ruleAttributeRejects(t, name, types.StringValue(value)) {
				t.Errorf("expected %s to reject %q, which the API would not store as written", name, value)
			}
		}
	}
}

func TestManagementRuleTagOperation_AllowsUnsetValue(t *testing.T) {
	for _, name := range []string{"source_tag_operation", "target_tag_operation"} {
		if ruleAttributeRejects(t, name, types.StringNull()) {
			t.Errorf("expected %s to allow an unset value, which rules that don't use 'tag' rely on", name)
		}
	}
}

// ruleConfig returns a rule matching the source's alias to the target's business_unit tag, with the given
// attributes overridden.
func ruleConfig(overrides map[string]attr.Value) map[string]attr.Value {
	rule := map[string]attr.Value{
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
		rule[name] = value
	}
	return rule
}

// planRuleAttribute runs the schema's plan modifiers for one attribute of a configured rule and returns the
// planned value. Like Terraform, it starts from the configured value, or from unknown when an optional computed
// attribute is unset.
func planRuleAttribute(t *testing.T, name string, rule map[string]attr.Value) types.String {
	t.Helper()
	ctx := context.Background()
	attribute, ok := opsleveltf.ManagementRulesResourceAttribute().NestedObject.Attributes[name].(resschema.StringAttribute)
	if !ok {
		t.Fatalf("management rule attribute %q is not a string attribute", name)
	}

	s := resschema.Schema{Attributes: map[string]resschema.Attribute{"management_rules": opsleveltf.ManagementRulesResourceAttribute()}}
	plan := tfsdk.Plan{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	ruleType := types.ObjectType{AttrTypes: opsleveltf.ManagementRuleModelAttrs()}
	rules := types.ListValueMust(ruleType, []attr.Value{types.ObjectValueMust(opsleveltf.ManagementRuleModelAttrs(), rule)})
	if diags := plan.Set(ctx, &struct {
		ManagementRules types.List `tfsdk:"management_rules"`
	}{ManagementRules: rules}); diags.HasError() {
		t.Fatalf("unable to build the config: %v", diags)
	}

	initial := rule[name].(types.String)
	if initial.IsNull() {
		initial = types.StringUnknown()
	}
	resp := &planmodifier.StringResponse{PlanValue: initial}
	for _, m := range attribute.PlanModifiers {
		m.PlanModifyString(ctx, planmodifier.StringRequest{
			Path:        path.Root("management_rules").AtListIndex(0).AtName(name),
			Config:      tfsdk.Config{Schema: s, Raw: plan.Raw},
			ConfigValue: rule[name].(types.String),
			Plan:        plan,
			PlanValue:   resp.PlanValue,
		}, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("plan modifier for %s failed: %v", name, resp.Diagnostics)
		}
	}
	return resp.PlanValue
}

func TestManagementRuleTagOperation_IsComputed(t *testing.T) {
	// Without Computed, Terraform rejects a planned value for an attribute the config leaves unset.
	for _, name := range []string{"source_tag_operation", "target_tag_operation"} {
		attribute := opsleveltf.ManagementRulesResourceAttribute().NestedObject.Attributes[name].(resschema.StringAttribute)
		if !attribute.Computed {
			t.Errorf("expected %s to be computed", name)
		}
	}
}

func TestManagementRuleTagOperation_PlansEqualsForUnsetTagRule(t *testing.T) {
	if got := planRuleAttribute(t, "target_tag_operation", ruleConfig(nil)); !got.Equal(types.StringValue("equals")) {
		t.Errorf("expected an unset target_tag_operation on a tag rule to plan \"equals\", got %s", got)
	}

	sourceTagRule := ruleConfig(map[string]attr.Value{
		"source_property":      types.StringValue("tag"),
		"source_tag_key":       types.StringValue("business_unit"),
		"target_property":      types.StringValue("alias"),
		"target_tag_key":       types.StringNull(),
		"target_tag_operation": types.StringNull(),
	})
	if got := planRuleAttribute(t, "source_tag_operation", sourceTagRule); !got.Equal(types.StringValue("equals")) {
		t.Errorf("expected an unset source_tag_operation on a tag rule to plan \"equals\", got %s", got)
	}
}

func TestManagementRuleTagOperation_PlansNullForRuleWithoutTag(t *testing.T) {
	aliasRule := ruleConfig(map[string]attr.Value{
		"target_property": types.StringValue("alias"),
		"target_tag_key":  types.StringNull(),
	})
	if got := planRuleAttribute(t, "target_tag_operation", aliasRule); !got.IsNull() {
		t.Errorf("expected no target_tag_operation for a rule whose target_property is not 'tag', got %s", got)
	}
	if got := planRuleAttribute(t, "source_tag_operation", aliasRule); !got.IsNull() {
		t.Errorf("expected no source_tag_operation for a rule whose source_property is not 'tag', got %s", got)
	}
}

func TestManagementRuleTagOperation_KeepsConfiguredOperation(t *testing.T) {
	startsWithRule := ruleConfig(map[string]attr.Value{"target_tag_operation": types.StringValue("starts_with")})
	if got := planRuleAttribute(t, "target_tag_operation", startsWithRule); !got.Equal(types.StringValue("starts_with")) {
		t.Errorf("expected a configured target_tag_operation to be kept, got %s", got)
	}
}
