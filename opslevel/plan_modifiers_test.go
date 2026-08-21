package opslevel

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
