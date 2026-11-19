package opslevel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	businessUnitComponentTypeID = "Z2lkOi8vb3BzbGV2ZWwvQ29tcG9uZW50VHlwZS8x"
	otherComponentTypeID        = "Z2lkOi8vb3BzbGV2ZWwvQ29tcG9uZW50VHlwZS8y"
)

// componentTypeAPIServer answers componentType lookups by ID, or by alias case-insensitively like OpsLevel does, and
// counts them. Anything else gets a GraphQL error, including "team", which the real API can't look up either.
func componentTypeAPIServer(t *testing.T) (*httptest.Server, func() int) {
	t.Helper()
	componentTypes := map[string]string{businessUnitComponentTypeID: "business_unit", otherComponentTypeID: "other_type"}
	var mu sync.Mutex
	lookups := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Variables struct {
				Input struct {
					Id    string `json:"id"`
					Alias string `json:"alias"`
				} `json:"input"`
			} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("unable to read request: %v", err)
			return
		}
		mu.Lock()
		lookups++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		for id, alias := range componentTypes {
			if id == req.Variables.Input.Id || (req.Variables.Input.Alias != "" && strings.EqualFold(alias, req.Variables.Input.Alias)) {
				body, _ := json.Marshal(map[string]any{"data": map[string]any{"account": map[string]any{
					"componentType": map[string]any{"id": id, "aliases": []string{alias}},
				}}})
				w.Write(body)
				return
			}
		}
		w.Write([]byte(`{"data":{"account":{"componentType":null}},"errors":[{"message":"component type not found"}]}`))
	}))

	return server, func() int {
		mu.Lock()
		defer mu.Unlock()
		return lookups
	}
}

func relationshipDefinitionSchema() schema.Schema {
	var resp resource.SchemaResponse
	(&RelationshipDefinitionResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp.Schema
}

// relationshipDefinitionRaw returns a definition with the given component_type as Terraform holds it, or a null
// object (no resource: create or destroy) when componentType is nil.
func relationshipDefinitionRaw(t *testing.T, componentType *types.String) tftypes.Value {
	t.Helper()
	ctx := context.Background()
	s := relationshipDefinitionSchema()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	if componentType == nil {
		return state.Raw
	}

	model := RelationshipDefinitionResourceModel{
		Id:                types.StringValue("Z2lkOi8vb3BzbGV2ZWwvUmVsYXRpb25zaGlwRGVmaW5pdGlvbi8x"),
		Name:              types.StringValue("Infrastructure"),
		Alias:             types.StringValue("infrastructure"),
		Description:       types.StringNull(),
		ComponentType:     *componentType,
		AllowedCategories: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("infrastructure")}),
		AllowedTypes:      types.ListNull(types.StringType),
		ManagementRules:   types.ListNull(types.ObjectType{AttrTypes: ManagementRuleModelAttrs()}),
	}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("unable to build the definition: %v", diags)
	}
	return state.Raw
}

func runModifyPlan(t *testing.T, serverURL string, state, plan *types.String) *resource.ModifyPlanResponse {
	t.Helper()
	s := relationshipDefinitionSchema()
	res := &RelationshipDefinitionResource{CommonResourceClient: CommonResourceClient{client: newTestClient(serverURL)}}
	planRaw := relationshipDefinitionRaw(t, plan)
	resp := &resource.ModifyPlanResponse{Plan: tfsdk.Plan{Schema: s, Raw: planRaw}}
	res.ModifyPlan(context.Background(), resource.ModifyPlanRequest{
		State: tfsdk.State{Schema: s, Raw: relationshipDefinitionRaw(t, state)},
		Plan:  tfsdk.Plan{Schema: s, Raw: planRaw},
	}, resp)
	return resp
}

func stringPtr(value string) *types.String {
	v := types.StringValue(value)
	return &v
}

func TestModifyPlan_UnchangedComponentTypeMakesNoLookups(t *testing.T) {
	server, lookups := componentTypeAPIServer(t)
	defer server.Close()

	resp := runModifyPlan(t, server.URL, stringPtr("business_unit"), stringPtr("business_unit"))

	if resp.Diagnostics.HasError() || len(resp.RequiresReplace) > 0 {
		t.Errorf("expected no replacement and no errors, got replace %v and %v", resp.RequiresReplace, resp.Diagnostics)
	}
	if got := lookups(); got != 0 {
		t.Errorf("expected no component type lookups for an unchanged component_type, got %d", got)
	}
}

func TestModifyPlan_SameAliasInAnotherCaseDoesNotReplace(t *testing.T) {
	server, lookups := componentTypeAPIServer(t)
	defer server.Close()

	resp := runModifyPlan(t, server.URL, stringPtr("business_unit"), stringPtr("Business_Unit"))

	if resp.Diagnostics.HasError() || len(resp.RequiresReplace) > 0 {
		t.Errorf("expected the same alias in another case not to replace, got replace %v and %v", resp.RequiresReplace, resp.Diagnostics)
	}
	if got := lookups(); got != 0 {
		t.Errorf("expected no component type lookups for a case-only change, got %d", got)
	}
}

func TestModifyPlan_AliasAndIDOfTheSameComponentTypeDoNotReplace(t *testing.T) {
	// An import puts the alias in state; the config can name the same component type by its ID.
	server, _ := componentTypeAPIServer(t)
	defer server.Close()

	resp := runModifyPlan(t, server.URL, stringPtr("business_unit"), stringPtr(businessUnitComponentTypeID))

	if resp.Diagnostics.HasError() || len(resp.RequiresReplace) > 0 {
		t.Errorf("expected the alias and the ID of the same component type not to replace, got replace %v and %v", resp.RequiresReplace, resp.Diagnostics)
	}
}

func TestModifyPlan_DifferentComponentTypeReplaces(t *testing.T) {
	server, _ := componentTypeAPIServer(t)
	defer server.Close()

	// "team" can't be looked up through the API, so it is treated as a different component type.
	for _, planned := range []string{"other_type", otherComponentTypeID, "team"} {
		resp := runModifyPlan(t, server.URL, stringPtr("business_unit"), stringPtr(planned))

		if resp.Diagnostics.HasError() {
			t.Errorf("expected no errors for %q, got %v", planned, resp.Diagnostics)
		}
		if !resp.RequiresReplace.Contains(path.Root("component_type")) {
			t.Errorf("expected moving to a different component type (%q) to replace the definition", planned)
		}
	}
}

func TestModifyPlan_UnresolvableComponentTypeReplacesAsBefore(t *testing.T) {
	// Without a lookup there is no way to tell it's the same component type, so this keeps the old RequiresReplace
	// behavior rather than failing the plan, which would also block `terraform destroy`.
	server, _ := componentTypeAPIServer(t)
	defer server.Close()

	resp := runModifyPlan(t, server.URL, stringPtr("business_unit"), stringPtr("no_such_component_type"))

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors, got %v", resp.Diagnostics)
	}
	if !resp.RequiresReplace.Contains(path.Root("component_type")) {
		t.Error("expected an unresolvable component_type to replace the definition, as before")
	}
}

func TestModifyPlan_CreateAndDestroyAreLeftAlone(t *testing.T) {
	server, lookups := componentTypeAPIServer(t)
	defer server.Close()

	for name, resp := range map[string]*resource.ModifyPlanResponse{
		"create":  runModifyPlan(t, server.URL, nil, stringPtr("business_unit")),
		"destroy": runModifyPlan(t, server.URL, stringPtr("business_unit"), nil),
	} {
		if resp.Diagnostics.HasError() || len(resp.RequiresReplace) > 0 {
			t.Errorf("expected %s to be left alone, got replace %v and %v", name, resp.RequiresReplace, resp.Diagnostics)
		}
	}
	if got := lookups(); got != 0 {
		t.Errorf("expected no component type lookups on create or destroy, got %d", got)
	}
}
