package opslevel

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ planmodifier.List = listSortStringModifier{}

type listSortStringModifier struct{}

// Description returns a plain text description of the validator's behavior,
// suitable for a practitioner to understand its impact.
func (m listSortStringModifier) Description(_ context.Context) string {
	return "Sorts a given list of strings"
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (m listSortStringModifier) MarkdownDescription(_ context.Context) string {
	return "Sorts list of strings when returned order of strings may vary"
}

// PlanModifyList updates the planned value with the default if its not null.
func (m listSortStringModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if ListValueStringsAreEqual(ctx, req.PlanValue, req.StateValue) {
		resp.PlanValue = req.StateValue
	}
}

func ListValueStringsAreEqual(ctx context.Context, listValue1, listValue2 basetypes.ListValue) bool {
	stringValues1, diags := ListValueToStringSlice(ctx, listValue1)
	if diags.HasError() {
		return false
	}
	stringValues2, diags := ListValueToStringSlice(ctx, listValue2)
	if diags.HasError() {
		return false
	}
	slices.Sort(stringValues1)
	slices.Sort(stringValues2)
	return slices.Equal(stringValues1, stringValues2)
}
