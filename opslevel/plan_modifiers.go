package opslevel

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	planStrings, diags := ListValueToStringSlice(ctx, req.PlanValue)
	resp.Diagnostics.Append(diags...)

	stateStrings, diags := ListValueToStringSlice(ctx, req.StateValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	slices.Sort(planStrings)
	slices.Sort(stateStrings)
	if slices.Equal(planStrings, stateStrings) {
		resp.PlanValue = req.StateValue
	}
}
