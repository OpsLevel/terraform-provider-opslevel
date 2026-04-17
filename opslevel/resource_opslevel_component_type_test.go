package opslevel

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2026"
	"golang.org/x/net/context"
)

// graphqlRequest represents an incoming GraphQL request
type graphqlRequest struct {
	Query     string          `json:"query"`
	Variables json.RawMessage `json:"variables"`
}

// testAPIServer creates an httptest server that tracks which GraphQL operations
// were called. Returns the server and a pointer to the list of operation names.
func testAPIServer(t *testing.T) (*httptest.Server, *[]string) {
	t.Helper()
	var mu sync.Mutex
	var operations []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var gqlReq graphqlRequest
		if err := json.Unmarshal(body, &gqlReq); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		mu.Lock()
		defer mu.Unlock()

		query := gqlReq.Query

		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.Contains(query, "relationshipDefinitionDelete"):
			operations = append(operations, "relationshipDefinitionDelete")
			w.Write([]byte(`{"data":{"relationshipDefinitionDelete":{"deletedId":"Z2lkOi8vMTIz","errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitionCreate"):
			operations = append(operations, "relationshipDefinitionCreate")
			w.Write([]byte(`{"data":{"relationshipDefinitionCreate":{"definition":{"id":"Z2lkOi8vbmV3","alias":"new_rel","name":"New","description":"","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]},"errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitionUpdate"):
			operations = append(operations, "relationshipDefinitionUpdate")
			w.Write([]byte(`{"data":{"relationshipDefinitionUpdate":{"definition":{"id":"Z2lkOi8vMTIz","alias":"managed_by","name":"Managed By","description":"","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]},"errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitions"):
			operations = append(operations, "relationshipDefinitionsList")
			// Return one existing relationship definition "managed_by"
			w.Write([]byte(`{"data":{"account":{"relationshipDefinitions":{"nodes":[{"id":"Z2lkOi8vMTIz","alias":"managed_by","name":"Managed By","description":"A managed relationship","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]}],"pageInfo":{"hasNextPage":false,"hasPreviousPage":false,"startCursor":"","endCursor":""}}}}}`))

		default:
			t.Logf("unhandled query: %s", query)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	return server, &operations
}

// testAPIServerTwoRels is like testAPIServer but returns two existing
// relationships: "managed_by" and "depends_on".
func testAPIServerTwoRels(t *testing.T) (*httptest.Server, *[]string) {
	t.Helper()
	var mu sync.Mutex
	var operations []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var gqlReq graphqlRequest
		if err := json.Unmarshal(body, &gqlReq); err != nil {
			t.Fatalf("failed to unmarshal request: %v", err)
		}

		mu.Lock()
		defer mu.Unlock()

		query := gqlReq.Query

		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.Contains(query, "relationshipDefinitionDelete"):
			operations = append(operations, "relationshipDefinitionDelete")
			w.Write([]byte(`{"data":{"relationshipDefinitionDelete":{"deletedId":"Z2lkOi8vMTIz","errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitionCreate"):
			operations = append(operations, "relationshipDefinitionCreate")
			w.Write([]byte(`{"data":{"relationshipDefinitionCreate":{"definition":{"id":"Z2lkOi8vbmV3","alias":"new_rel","name":"New","description":"","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]},"errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitionUpdate"):
			operations = append(operations, "relationshipDefinitionUpdate")
			w.Write([]byte(`{"data":{"relationshipDefinitionUpdate":{"definition":{"id":"Z2lkOi8vMTIz","alias":"managed_by","name":"Managed By","description":"","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]},"errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitions"):
			operations = append(operations, "relationshipDefinitionsList")
			w.Write([]byte(`{"data":{"account":{"relationshipDefinitions":{"nodes":[` +
				`{"id":"Z2lkOi8vMTIz","alias":"managed_by","name":"Managed By","description":"","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]},` +
				`{"id":"Z2lkOi8vNDU2","alias":"depends_on","name":"Depends On","description":"","metadata":{"allowedCategories":[],"allowedTypes":["service"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]}` +
				`],"pageInfo":{"hasNextPage":false,"hasPreviousPage":false,"startCursor":"","endCursor":""}}}}}`))

		default:
			t.Logf("unhandled query: %s", query)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	return server, &operations
}

func newTestClient(serverURL string) *opslevel.Client {
	return opslevel.NewGQLClient(
		opslevel.SetURL(serverURL+"/LOCAL_TESTING/test"),
		opslevel.SetAPIToken("test-token"),
		opslevel.SetMaxRetries(0),
	)
}

func managedByRelationship() map[string]RelationshipModel {
	return map[string]RelationshipModel{
		"managed_by": {
			Name:              types.StringValue("Managed By"),
			AllowedCategories: types.ListNull(types.StringType),
			AllowedTypes:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("team")}),
		},
	}
}

func dependsOnRelationship() map[string]RelationshipModel {
	return map[string]RelationshipModel{
		"depends_on": {
			Name:              types.StringValue("Depends On"),
			AllowedCategories: types.ListNull(types.StringType),
			AllowedTypes:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("service")}),
		},
	}
}

func bothRelationships() map[string]RelationshipModel {
	return map[string]RelationshipModel{
		"managed_by": {
			Name:              types.StringValue("Managed By"),
			AllowedCategories: types.ListNull(types.StringType),
			AllowedTypes:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("team")}),
		},
		"depends_on": {
			Name:              types.StringValue("Depends On"),
			AllowedCategories: types.ListNull(types.StringType),
			AllowedTypes:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("service")}),
		},
	}
}

func countOps(ops []string, name string) int {
	count := 0
	for _, op := range ops {
		if op == name {
			count++
		}
	}
	return count
}

// TestReconcileRelationships_NilPlanNilState verifies that when neither plan
// nor state has relationships, no API calls are made to delete existing
// API-side relationships (e.g., ones created via the UI).
// Note: the guard in Update() skips calling reconcileRelationships entirely
// in this case. This test calls reconcileRelationships directly to verify
// it is also safe if called.
func TestReconcileRelationships_NilPlanNilState(t *testing.T) {
	server, operations := testAPIServer(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	planModel := ComponentTypeModel{
		Relationships: nil,
	}
	stateModel := ComponentTypeModel{
		Relationships: nil,
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes when plan and state are both nil, got %d. Operations: %v", deleteCount, ops)
	}
}

// TestReconcileRelationships_NilPlanWithState verifies that when relationships
// were in state but the plan drops them (user removed the block), only the
// state-managed relationships are deleted.
func TestReconcileRelationships_NilPlanWithState(t *testing.T) {
	server, operations := testAPIServer(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	planModel := ComponentTypeModel{
		Relationships: nil, // user removed the relationships block
	}
	stateModel := ComponentTypeModel{
		Relationships: managedByRelationship(), // was previously in state
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 1 {
		t.Errorf("expected 1 delete when state had relationship but plan removed it, got %d. Operations: %v", deleteCount, ops)
	}
}

// TestReconcileRelationships_PlanMatchesAPI verifies that when the plan includes
// the same relationships as the API, nothing is deleted -- only updated.
func TestReconcileRelationships_PlanMatchesAPI(t *testing.T) {
	server, operations := testAPIServer(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	rels := managedByRelationship()
	planModel := ComponentTypeModel{
		Relationships: rels,
	}
	stateModel := ComponentTypeModel{
		Relationships: rels,
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes when plan matches API, got %d. Operations: %v", deleteCount, ops)
	}
	if updateCount := countOps(ops, "relationshipDefinitionUpdate"); updateCount != 1 {
		t.Errorf("expected 1 update when plan matches API, got %d. Operations: %v", updateCount, ops)
	}
}

// TestReconcileRelationships_EmptyPlanNilState verifies that an empty map
// in plan with nil state does not delete API-side relationships.
func TestReconcileRelationships_EmptyPlanNilState(t *testing.T) {
	server, operations := testAPIServer(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	planModel := ComponentTypeModel{
		Relationships: map[string]RelationshipModel{}, // empty, not nil
	}
	stateModel := ComponentTypeModel{
		Relationships: nil,
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes with empty plan and nil state, got %d. Operations: %v", deleteCount, ops)
	}
}

// TestReconcileRelationships_APIHasExtraNotInState verifies that relationships
// existing on the API but never managed by Terraform (not in state) are left alone.
func TestReconcileRelationships_APIHasExtraNotInState(t *testing.T) {
	server, operations := testAPIServer(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	// Plan and state are both empty -- user never managed relationships via TF.
	// But the API has "managed_by" (created via UI).
	planModel := ComponentTypeModel{
		Relationships: map[string]RelationshipModel{},
	}
	stateModel := ComponentTypeModel{
		Relationships: map[string]RelationshipModel{},
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes for API-only relationships not in state, got %d. Operations: %v", deleteCount, ops)
	}
}

// TestReconcileRelationships_RemoveOneKeepOne verifies that when the user removes
// one relationship from config but keeps another, only the removed one is deleted.
func TestReconcileRelationships_RemoveOneKeepOne(t *testing.T) {
	server, operations := testAPIServerTwoRels(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	// Plan keeps "managed_by" but drops "depends_on"
	planModel := ComponentTypeModel{
		Relationships: managedByRelationship(),
	}
	// State had both
	stateModel := ComponentTypeModel{
		Relationships: bothRelationships(),
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 1 {
		t.Errorf("expected 1 delete (depends_on removed), got %d. Operations: %v", deleteCount, ops)
	}
	if updateCount := countOps(ops, "relationshipDefinitionUpdate"); updateCount != 1 {
		t.Errorf("expected 1 update (managed_by kept), got %d. Operations: %v", updateCount, ops)
	}
}

// TestReconcileRelationships_AddNewRelationship verifies that when the plan adds
// a relationship not on the API, it gets created.
func TestReconcileRelationships_AddNewRelationship(t *testing.T) {
	server, operations := testAPIServer(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	// Plan has existing "managed_by" plus a new "depends_on"
	planModel := ComponentTypeModel{
		Relationships: bothRelationships(),
	}
	// State only had "managed_by"
	stateModel := ComponentTypeModel{
		Relationships: managedByRelationship(),
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if createCount := countOps(ops, "relationshipDefinitionCreate"); createCount != 1 {
		t.Errorf("expected 1 create (depends_on new), got %d. Operations: %v", createCount, ops)
	}
	if updateCount := countOps(ops, "relationshipDefinitionUpdate"); updateCount != 1 {
		t.Errorf("expected 1 update (managed_by existing), got %d. Operations: %v", updateCount, ops)
	}
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes, got %d. Operations: %v", deleteCount, ops)
	}
}

// TestReconcileRelationships_RemoveAllFromState verifies that when the user
// explicitly removes all relationships (had them in state, plan is now empty),
// all state-managed relationships are deleted.
func TestReconcileRelationships_RemoveAllFromState(t *testing.T) {
	server, operations := testAPIServerTwoRels(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	// Plan has empty relationships (user removed the block)
	planModel := ComponentTypeModel{
		Relationships: map[string]RelationshipModel{},
	}
	// State had both relationships
	stateModel := ComponentTypeModel{
		Relationships: bothRelationships(),
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 2 {
		t.Errorf("expected 2 deletes (both removed from plan), got %d. Operations: %v", deleteCount, ops)
	}
}

// TestReconcileRelationships_StateHasOneAPIHasExtraKeepBoth verifies that when
// state manages one relationship and the API has an extra one (created via UI),
// the extra one is not deleted even when plan matches state.
func TestReconcileRelationships_StateHasOneAPIHasExtraKeepBoth(t *testing.T) {
	server, operations := testAPIServerTwoRels(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	// Plan and state both have only "managed_by".
	// API has both "managed_by" and "depends_on" (depends_on created via UI).
	planModel := ComponentTypeModel{
		Relationships: managedByRelationship(),
	}
	stateModel := ComponentTypeModel{
		Relationships: managedByRelationship(),
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel, stateModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes (depends_on is API-only, not in state), got %d. Operations: %v", deleteCount, ops)
	}
	if updateCount := countOps(ops, "relationshipDefinitionUpdate"); updateCount != 1 {
		t.Errorf("expected 1 update (managed_by), got %d. Operations: %v", updateCount, ops)
	}
}
