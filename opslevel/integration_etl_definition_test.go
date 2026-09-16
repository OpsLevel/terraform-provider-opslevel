package opslevel

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUseStateForEquivalentYAML(t *testing.T) {
	t.Parallel()

	stateValue := types.StringValue("---\nextractors:\n- external_kind: widget\n  external_id: \".id\"\n")

	testCases := map[string]struct {
		planValue types.String
		expected  types.String
	}{
		"equivalent-yaml": {
			planValue: types.StringValue("extractors:\n  - external_kind: widget\n    external_id: \".id\"\n"),
			expected:  stateValue,
		},
		"different-yaml": {
			planValue: types.StringValue("extractors:\n  - external_kind: widget\n    external_id: \".other\"\n"),
			expected:  types.StringValue("extractors:\n  - external_kind: widget\n    external_id: \".other\"\n"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				StateValue: stateValue,
				PlanValue:  testCase.planValue,
			}
			resp := &planmodifier.StringResponse{
				PlanValue: req.PlanValue,
			}

			UseStateForEquivalentYAML().PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected) {
				t.Fatalf("expected plan value %q, got %q", testCase.expected.ValueString(), resp.PlanValue.ValueString())
			}
		})
	}
}

func TestNonEmptyYAML(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value       string
		expectError bool
	}{
		"empty-string":     {value: "", expectError: true},
		"single-space":     {value: " ", expectError: true},
		"whitespace-only":  {value: "   \n\t\n", expectError: true},
		"document-marker":  {value: "---\n", expectError: true},
		"explicit-null":    {value: "null", expectError: true},
		"tilde-null":       {value: "~", expectError: true},
		"comment-only":     {value: "# nothing here\n", expectError: true},
		"valid-definition": {value: "extractors:\n- external_kind: widget\n", expectError: false},
		"malformed-yaml":   {value: "extractors:\n\t- bad tab indent\n", expectError: true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:        path.Root("etl_definition").AtName("extract_definition"),
				ConfigValue: types.StringValue(testCase.value),
			}
			resp := &validator.StringResponse{}

			NonEmptyYAML().ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() != testCase.expectError {
				t.Fatalf("expected error=%v for %q, got error=%v (%v)",
					testCase.expectError, testCase.value, resp.Diagnostics.HasError(), resp.Diagnostics)
			}
		})
	}
}

// A custom integration created without an etl_definition block stores an empty
// pair in state. Terraform proposes that state as the plan on the next update, so
// this path must omit the fields rather than error, or the resource can never be
// updated again.
func TestIntegrationEtlDefinitionInputOmitsEmptyPair(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	diags := &diag.Diagnostics{}
	state := reconcileIntegrationEtlDefinition(ctx, "", "", types.ObjectNull(integrationEtlDefinitionAttrs()), diags)

	if state.IsNull() || state.IsUnknown() {
		t.Fatalf("expected a known object in state, got null=%v unknown=%v", state.IsNull(), state.IsUnknown())
	}

	inputDiags := &diag.Diagnostics{}
	extract, transform := integrationEtlDefinitionInput(ctx, state, inputDiags)

	if inputDiags.HasError() {
		t.Fatalf("an empty pair from state must not error, got: %v", inputDiags.Errors())
	}
	if extract != nil || transform != nil {
		t.Fatalf("expected both definitions omitted, got extract=%v transform=%v", extract, transform)
	}
}
