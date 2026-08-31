package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
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

type useStateForEquivalentYAMLModifier struct{}

func (m useStateForEquivalentYAMLModifier) Description(ctx context.Context) string {
	return "Preserves prior state when planned YAML is semantically equivalent."
}

func (m useStateForEquivalentYAMLModifier) MarkdownDescription(ctx context.Context) string {
	return "Preserves prior state when planned YAML is semantically equivalent."
}

func (m useStateForEquivalentYAMLModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	if yamlEquivalent(req.PlanValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

func UseStateForEquivalentYAML() planmodifier.String {
	return useStateForEquivalentYAMLModifier{}
}

type nonEmptyYAMLValidator struct{}

func (v nonEmptyYAMLValidator) Description(ctx context.Context) string {
	return "must be a YAML document that does not decode to null"
}

func (v nonEmptyYAMLValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Rejects "", " ", "---", "null" and "~". They all decode to nil, which the API
// stores as a cleared definition rather than rejecting - so without this they
// silently wipe an existing definition and still produce a clean plan.
func (v nonEmptyYAMLValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var decoded any
	if err := yaml.Unmarshal([]byte(req.ConfigValue.ValueString()), &decoded); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid YAML",
			fmt.Sprintf("Expected a valid YAML document, got error: %s", err),
		)
		return
	}

	if decoded == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Empty YAML definition",
			"This decodes to null, which clears the definition stored in OpsLevel. "+
				"Omit the entire etl_definition block instead of setting an empty definition.",
		)
	}
}

func NonEmptyYAML() validator.String {
	return nonEmptyYAMLValidator{}
}

// ImmutableViolation reports whether moving a write-once attribute from
// stateValue to planValue is forbidden, and why.
//
// A write-once attribute records something that happened at create time and
// cannot be restated afterwards. Setting one where none exists is refused as
// well as changing one: an imported resource has no record of the value, and
// the action it would describe has already taken place.
//
// An unknown plan value cannot be judged, since the value it will take is not
// known until apply.
func ImmutableViolation(stateValue, planValue types.String) (bool, string) {
	if planValue.IsUnknown() {
		return false, ""
	}
	if planValue.Equal(stateValue) {
		return false, ""
	}
	if stateValue.IsNull() {
		return true, "it cannot be set on a resource that does not already have it"
	}
	if planValue.IsNull() {
		return true, fmt.Sprintf("it cannot be removed once set (currently %q)", stateValue.ValueString())
	}
	return true, fmt.Sprintf("it cannot be changed once set (currently %q)", stateValue.ValueString())
}

// immutableStringPlanModifier refuses a change to a write-once attribute at plan
// time, rather than planning a replacement.
//
// RequiresReplace would be wrong here: replacing a campaign check deletes the
// copy and makes a new one, discarding the check's results and history, which is
// too destructive to happen implicitly from an edited attribute.
type immutableStringPlanModifier struct{}

func (m immutableStringPlanModifier) Description(ctx context.Context) string {
	return "Refuses changes to a write-once attribute after it has been created"
}

func (m immutableStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m immutableStringPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Creating: there is no prior state to preserve. Destroying: the value is
	// going away with the resource.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	violated, reason := ImmutableViolation(req.StateValue, req.PlanValue)
	if !violated {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		fmt.Sprintf("Cannot Change %s", req.Path),
		fmt.Sprintf(
			"%s is set when the resource is created and %s. "+
				"Destroy and recreate the resource if you really intend this - "+
				"note that doing so deletes the campaign check and its results and history.",
			req.Path, reason,
		),
	)
}

func ImmutableStringPlanModifier() planmodifier.String {
	return immutableStringPlanModifier{}
}
