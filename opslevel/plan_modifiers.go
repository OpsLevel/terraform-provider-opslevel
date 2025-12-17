package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// preventRemovalPlanModifier checks if list items are being removed and returns an error
type preventRemovalPlanModifier struct {
	fieldName string
}

func (m preventRemovalPlanModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Prevents removal of %s from relationship definitions", m.fieldName)
}

func (m preventRemovalPlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Prevents removal of %s from relationship definitions", m.fieldName)
}

func (m preventRemovalPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If the resource is being created or destroyed, no validation needed
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	var stateList, planList []string

	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		diags := req.StateValue.ElementsAs(ctx, &stateList, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
		diags := req.PlanValue.ElementsAs(ctx, &planList, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	planMap := make(map[string]bool)
	for _, item := range planList {
		planMap[item] = true
	}

	var removedItems []string
	for _, item := range stateList {
		if !planMap[item] {
			removedItems = append(removedItems, item)
		}
	}

	if len(removedItems) > 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			fmt.Sprintf("Cannot Remove %s", m.fieldName),
			fmt.Sprintf(
				"The OpsLevel API does not support removing %s from relationship definitions. "+
					"The following items cannot be removed: %v. "+
					"You can only add new items to the list.",
				m.fieldName,
				removedItems,
			),
		)
	}
}

func PreventRemovalPlanModifier(fieldName string) planmodifier.List {
	return preventRemovalPlanModifier{
		fieldName: fieldName,
	}
}
