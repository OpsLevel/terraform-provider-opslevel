package opslevel_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	resschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
