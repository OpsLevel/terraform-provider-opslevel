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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

	planModel := ComponentTypeModel{
		Relationships: managedByRelationship(),
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel)

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

// TestReconcileRelationships_RemoveOneKeepOne verifies that when the plan keeps
// one relationship and drops another, the dropped one is deleted and the kept
// one is updated.
func TestReconcileRelationships_RemoveOneKeepOne(t *testing.T) {
	server, operations := testAPIServerTwoRels(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	planModel := ComponentTypeModel{
		Relationships: managedByRelationship(),
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel)

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
// a relationship not on the API, it gets created and the existing one is updated.
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

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel)

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

// TestReconcileRelationships_EmptyPlanDeletesAll verifies that when the plan
// passes an empty (but non-nil) relationships map -- i.e., the user explicitly
// set `relationships = {}` -- every relationship on the component type is
// deleted. This is the contract the guard in Update() protects: reconcile only
// runs when the user is managing relationships via this resource.
func TestReconcileRelationships_EmptyPlanDeletesAll(t *testing.T) {
	server, operations := testAPIServerTwoRels(t)
	defer server.Close()

	client := newTestClient(server.URL)
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: client},
	}

	ctx := context.Background()
	resp := &resource.UpdateResponse{}

	planModel := ComponentTypeModel{
		Relationships: map[string]RelationshipModel{},
	}

	res.reconcileRelationships(ctx, nil, "Z2lkOi8vY3Q", resp, planModel)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	ops := *operations
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 2 {
		t.Errorf("expected 2 deletes (both removed from plan), got %d. Operations: %v", deleteCount, ops)
	}
}

// testComponentTypeServer answers the requests ComponentTypeResource.Read and Update make:
// the component type, its update, its properties, and its relationship definitions (one,
// "managed_by"), including deleting one.
func testComponentTypeServer(t *testing.T) (*httptest.Server, *[]string) {
	t.Helper()
	var mu sync.Mutex
	var operations []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}

		var gqlReq graphqlRequest
		if err := json.Unmarshal(body, &gqlReq); err != nil {
			t.Errorf("failed to unmarshal request: %v", err)
		}

		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		// The properties and relationship definition queries also mention componentType, so match them first.
		switch query := gqlReq.Query; {
		case strings.Contains(query, "componentTypeUpdate("):
			operations = append(operations, "componentTypeUpdate")
			w.Write([]byte(`{"data":{"componentTypeUpdate":{"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"],"description":"","href":"","icon":{"color":"#000000","name":"PhBird"},"isDefault":false,"name":"Test Type","ownerRelationship":{"managementRules":[]},"systemRelationship":{"managementRules":[]},"timestamps":{"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}},"errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitionDelete("):
			operations = append(operations, "relationshipDefinitionDelete")
			w.Write([]byte(`{"data":{"relationshipDefinitionDelete":{"deletedId":"Z2lkOi8vMTIz","errors":[]}}}`))

		case strings.Contains(query, "relationshipDefinitions("):
			operations = append(operations, "relationshipDefinitionsList")
			w.Write([]byte(`{"data":{"account":{"relationshipDefinitions":{"nodes":[{"id":"Z2lkOi8vMTIz","alias":"managed_by","name":"Managed By","description":"","metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"]},"managementRules":[]}],"pageInfo":{"hasNextPage":false,"hasPreviousPage":false,"startCursor":"","endCursor":""}}}}}`))

		case strings.Contains(query, "properties("):
			operations = append(operations, "componentTypePropertyList")
			w.Write([]byte(`{"data":{"account":{"componentType":{"properties":{"nodes":[],"pageInfo":{"hasNextPage":false,"hasPreviousPage":false,"startCursor":"","endCursor":""}}}}}}`))

		case strings.Contains(query, "componentType("):
			operations = append(operations, "componentTypeGet")
			w.Write([]byte(`{"data":{"account":{"componentType":{"id":"Z2lkOi8vY3Q","aliases":["test_type"],"description":"","href":"","icon":{"color":"#000000","name":"PhBird"},"isDefault":false,"name":"Test Type","ownerRelationship":{"managementRules":[]},"systemRelationship":{"managementRules":[]},"timestamps":{"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}}}}}`))

		default:
			t.Errorf("unhandled query: %s", gqlReq.Query)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	return server, &operations
}

// readComponentType runs Read with only the ID in state, plus `relationships = {}` when emptyRelationships is true.
func readComponentType(t *testing.T, emptyRelationships bool) (resource.ReadResponse, []string) {
	t.Helper()
	server, operations := testComponentTypeServer(t)
	t.Cleanup(server.Close)

	ctx := context.Background()
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: newTestClient(server.URL)},
	}

	var schemaResp resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	attributes := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	attributes["id"] = tftypes.NewValue(tftypes.String, "Z2lkOi8vY3Q")
	if emptyRelationships {
		attributes["relationships"] = tftypes.NewValue(objectType.AttributeTypes["relationships"], map[string]tftypes.Value{})
	}
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, attributes)}

	resp := resource.ReadResponse{State: state}
	res.Read(ctx, resource.ReadRequest{State: state}, &resp)
	return resp, *operations
}

// TestComponentTypeRead_EmptyRelationshipsReadsEveryDefinition verifies that an explicitly
// empty `relationships` map reads every definition on the component type into state, so
// the plan shows the definitions that `relationships = {}` would delete.
func TestComponentTypeRead_EmptyRelationshipsReadsEveryDefinition(t *testing.T) {
	resp, ops := readComponentType(t, true)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	var relationships map[string]RelationshipModel
	if diags := resp.State.GetAttribute(context.Background(), path.Root("relationships"), &relationships); diags.HasError() {
		t.Fatalf("unable to read relationships from state: %v", diags.Errors())
	}
	if _, ok := relationships["managed_by"]; !ok || len(relationships) != 1 {
		t.Errorf("expected managed_by in state, got %v. Operations: %v", relationships, ops)
	}
}

// TestComponentTypeRead_UnsetRelationshipsStayUnset verifies that a component type that
// doesn't set `relationships` never reads relationship definitions into state.
func TestComponentTypeRead_UnsetRelationshipsStayUnset(t *testing.T) {
	resp, ops := readComponentType(t, false)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	var relationships map[string]RelationshipModel
	if diags := resp.State.GetAttribute(context.Background(), path.Root("relationships"), &relationships); diags.HasError() {
		t.Fatalf("unable to read relationships from state: %v", diags.Errors())
	}
	if relationships != nil {
		t.Errorf("expected relationships to stay unset, got %v", relationships)
	}
	if listCount := countOps(ops, "relationshipDefinitionsList"); listCount != 0 {
		t.Errorf("expected no relationship definitions query, got %d. Operations: %v", listCount, ops)
	}
}

// updateComponentType runs Update with the given `relationships` in the plan and in the prior state.
func updateComponentType(t *testing.T, planRelationships, stateRelationships map[string]RelationshipModel) (resource.UpdateResponse, []string) {
	t.Helper()
	server, operations := testComponentTypeServer(t)
	t.Cleanup(server.Close)

	ctx := context.Background()
	res := ComponentTypeResource{
		CommonResourceClient: CommonResourceClient{client: newTestClient(server.URL)},
	}

	var schemaResp resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	nullObject := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)
	model := func(relationships map[string]RelationshipModel) ComponentTypeModel {
		return ComponentTypeModel{
			Id:            types.StringValue("Z2lkOi8vY3Q"),
			Name:          types.StringValue("Test Type"),
			Alias:         types.StringValue("test_type"),
			Description:   types.StringValue(""),
			Icon:          &ComponentTypeIconModel{Color: types.StringValue("#000000"), Name: types.StringValue("PhBird")},
			Properties:    map[string]PropertyModel{},
			Relationships: relationships,
		}
	}

	plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: nullObject}
	if diags := plan.Set(ctx, model(planRelationships)); diags.HasError() {
		t.Fatalf("unable to build plan: %v", diags.Errors())
	}
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: nullObject}
	if diags := state.Set(ctx, model(stateRelationships)); diags.HasError() {
		t.Fatalf("unable to build state: %v", diags.Errors())
	}

	resp := resource.UpdateResponse{State: state}
	res.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)
	return resp, *operations
}

// TestComponentTypeUpdate_RemovedRelationshipsStopsManaging verifies that removing
// `relationships` stops managing the definitions without deleting them.
func TestComponentTypeUpdate_RemovedRelationshipsStopsManaging(t *testing.T) {
	resp, ops := updateComponentType(t, nil, managedByRelationship())
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	if listCount := countOps(ops, "relationshipDefinitionsList"); listCount != 0 {
		t.Errorf("expected no reconcile, got %d relationship definition lists. Operations: %v", listCount, ops)
	}
	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 0 {
		t.Errorf("expected 0 deletes, got %d. Operations: %v", deleteCount, ops)
	}

	var relationships map[string]RelationshipModel
	if diags := resp.State.GetAttribute(context.Background(), path.Root("relationships"), &relationships); diags.HasError() {
		t.Fatalf("unable to read relationships from state: %v", diags.Errors())
	}
	if relationships != nil {
		t.Errorf("expected relationships to be unset in state, got %v", relationships)
	}
}

// TestComponentTypeUpdate_EmptyRelationshipsDeletesEveryDefinition verifies that
// `relationships = {}` means the component type has no relationship definitions, so
// Update deletes every definition on it.
func TestComponentTypeUpdate_EmptyRelationshipsDeletesEveryDefinition(t *testing.T) {
	resp, ops := updateComponentType(t, map[string]RelationshipModel{}, map[string]RelationshipModel{})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	if deleteCount := countOps(ops, "relationshipDefinitionDelete"); deleteCount != 1 {
		t.Errorf("expected 1 delete (managed_by), got %d. Operations: %v", deleteCount, ops)
	}
}
