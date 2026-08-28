package opslevel_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func TestImmutableViolation_UnchangedValueIsAllowed(t *testing.T) {
	violated, reason := opsleveltf.ImmutableViolation(
		types.StringValue("check-a"),
		types.StringValue("check-a"),
	)
	if violated {
		t.Errorf("expected no violation for an unchanged value, got %q", reason)
	}
}

func TestImmutableViolation_ChangingValueIsForbidden(t *testing.T) {
	violated, reason := opsleveltf.ImmutableViolation(
		types.StringValue("check-a"),
		types.StringValue("check-b"),
	)
	if !violated {
		t.Fatal("expected changing a write-once value to be a violation")
	}
	if reason == "" {
		t.Error("expected a reason explaining the violation")
	}
}

// An imported campaign check has no record of what it was copied from, and the
// copy already happened, so the user cannot supply one after the fact.
func TestImmutableViolation_SettingAValueWhereNoneExistsIsForbidden(t *testing.T) {
	violated, reason := opsleveltf.ImmutableViolation(
		types.StringNull(),
		types.StringValue("check-a"),
	)
	if !violated {
		t.Fatal("expected setting a previously-null write-once value to be a violation")
	}
	if reason == "" {
		t.Error("expected a reason explaining the violation")
	}
}

func TestImmutableViolation_RemovingAValueIsForbidden(t *testing.T) {
	violated, _ := opsleveltf.ImmutableViolation(
		types.StringValue("check-a"),
		types.StringNull(),
	)
	if !violated {
		t.Error("expected removing a write-once value to be a violation")
	}
}

// A plan value can be unknown when it references a resource that does not exist
// yet. There is nothing to compare, so it cannot be judged a violation.
func TestImmutableViolation_UnknownPlanValueIsAllowed(t *testing.T) {
	violated, _ := opsleveltf.ImmutableViolation(
		types.StringValue("check-a"),
		types.StringUnknown(),
	)
	if violated {
		t.Error("expected an unknown plan value to be allowed")
	}
}

func TestImmutableViolation_BothNullIsAllowed(t *testing.T) {
	violated, _ := opsleveltf.ImmutableViolation(
		types.StringNull(),
		types.StringNull(),
	)
	if violated {
		t.Error("expected two null values to be allowed")
	}
}

func objectValue(null bool) tftypes.Value {
	t := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
	if null {
		return tftypes.NewValue(t, nil)
	}
	return tftypes.NewValue(t, map[string]tftypes.Value{})
}

func runImmutableModifier(state, plan tftypes.Value, stateValue, planValue types.String) diag.Diagnostics {
	req := planmodifier.StringRequest{
		Path:       path.Root("copied_from_check_id"),
		State:      tfsdk.State{Raw: state},
		Plan:       tfsdk.Plan{Raw: plan},
		StateValue: stateValue,
		PlanValue:  planValue,
	}
	resp := &planmodifier.StringResponse{PlanValue: planValue}
	opsleveltf.ImmutableStringPlanModifier().PlanModifyString(context.Background(), req, resp)
	return resp.Diagnostics
}

// On create there is no prior state, so every write-once attribute looks like it
// is being "set where none exists". Creating must not error.
func TestImmutableStringPlanModifier_AllowsCreate(t *testing.T) {
	diags := runImmutableModifier(
		objectValue(true), objectValue(false),
		types.StringNull(), types.StringValue("check-a"),
	)
	if diags.HasError() {
		t.Errorf("expected create to be allowed, got %v", diags.Errors())
	}
}

func TestImmutableStringPlanModifier_AllowsDestroy(t *testing.T) {
	diags := runImmutableModifier(
		objectValue(false), objectValue(true),
		types.StringValue("check-a"), types.StringNull(),
	)
	if diags.HasError() {
		t.Errorf("expected destroy to be allowed, got %v", diags.Errors())
	}
}

func TestImmutableStringPlanModifier_ErrorsOnChange(t *testing.T) {
	diags := runImmutableModifier(
		objectValue(false), objectValue(false),
		types.StringValue("check-a"), types.StringValue("check-b"),
	)
	if !diags.HasError() {
		t.Fatal("expected changing a write-once attribute to error")
	}
}
