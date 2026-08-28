package opslevel_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2026"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func campaignCheckAt(t *testing.T, list types.List, i int) map[string]interface{} {
	t.Helper()
	elements := list.Elements()
	if i >= len(elements) {
		t.Fatalf("wanted element %d, list has %d", i, len(elements))
	}
	obj, ok := elements[i].(types.Object)
	if !ok {
		t.Fatalf("expected element %d to be an object, got %T", i, elements[i])
	}
	out := map[string]interface{}{}
	for k, v := range obj.Attributes() {
		out[k] = v
	}
	return out
}

func TestCampaignCheckListValue_ReportsIdNameAndSource(t *testing.T) {
	checks := []opslevel.CampaignCheckNode{
		{Id: "copy-z", Name: "Arize", SourceCheck: opslevel.CheckId{Id: "check-a", Name: "Arize"}},
	}

	list, diags := opsleveltf.CampaignCheckListValue(context.Background(), checks)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}

	got := campaignCheckAt(t, list, 0)
	if got["id"] != types.StringValue("copy-z") {
		t.Errorf("id: got %v", got["id"])
	}
	if got["name"] != types.StringValue("Arize") {
		t.Errorf("name: got %v", got["name"])
	}
	if got["source_check_id"] != types.StringValue("check-a") {
		t.Errorf("source_check_id: got %v", got["source_check_id"])
	}
}

// A copy whose source has been deleted, or which predates OpsLevel recording
// lineage, reports no source check. That must surface as null rather than an
// empty string, so a migrating user can tell "unknown" from a real id.
func TestCampaignCheckListValue_MissingSourceIsNull(t *testing.T) {
	checks := []opslevel.CampaignCheckNode{
		{Id: "copy-z", Name: "Arize"},
	}

	list, diags := opsleveltf.CampaignCheckListValue(context.Background(), checks)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}

	got := campaignCheckAt(t, list, 0)
	if !got["source_check_id"].(types.String).IsNull() {
		t.Errorf("expected source_check_id to be null, got %v", got["source_check_id"])
	}
}

func TestCampaignCheckListValue_NoChecksIsEmptyNotNull(t *testing.T) {
	list, diags := opsleveltf.CampaignCheckListValue(context.Background(), []opslevel.CampaignCheckNode{})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}
	if list.IsNull() {
		t.Error("expected an empty list, got null")
	}
	if len(list.Elements()) != 0 {
		t.Errorf("expected no elements, got %d", len(list.Elements()))
	}
}
