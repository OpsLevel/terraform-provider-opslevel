package opslevel_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	opslevelgo "github.com/opslevel/opslevel-go/v2026"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

const (
	componentTypeID      = "Z2lkOi8vb3BzbGV2ZWwvQ29tcG9uZW50VHlwZS8x" // gid://opslevel/ComponentType/1
	otherComponentTypeID = "Z2lkOi8vb3BzbGV2ZWwvQ29tcG9uZW50VHlwZS8y" // gid://opslevel/ComponentType/2
)

func relationshipDefinition() opslevelgo.RelationshipDefinitionType {
	return opslevelgo.RelationshipDefinitionType{
		Id:            opslevelgo.ID("Z2lkOi8vb3BzbGV2ZWwvUmVsYXRpb25zaGlwRGVmaW5pdGlvbi8x"),
		Alias:         "infrastructure",
		Name:          "Infrastructure",
		ComponentType: opslevelgo.ComponentTypeId{Id: opslevelgo.ID(componentTypeID), Aliases: []string{"business_unit"}},
	}
}

func componentTypeReadWith(stateValue types.String) types.String {
	model := opsleveltf.RelationshipDefinitionResourceModel{ComponentType: stateValue}
	return opsleveltf.NewRelationshipDefinitionResourceModel(relationshipDefinition(), model).ComponentType
}

func TestNewRelationshipDefinitionResourceModel_ImportTakesComponentTypeFromAPI(t *testing.T) {
	// ImportState sets only the ID, so Read starts from a model with no component_type.
	if got := componentTypeReadWith(types.StringNull()); !got.Equal(types.StringValue("business_unit")) {
		t.Errorf("expected an imported definition to read component_type from the API, got %s", got)
	}
}

func TestNewRelationshipDefinitionResourceModel_ReadsTheCurrentAliasFromTheAPI(t *testing.T) {
	// An alias in state that no longer matches, for example after the component type's alias was renamed.
	if got := componentTypeReadWith(types.StringValue("old_business_unit")); !got.Equal(types.StringValue("business_unit")) {
		t.Errorf("expected component_type to read the API's current alias, got %s", got)
	}
}

func TestNewRelationshipDefinitionResourceModel_KeepsTheConfiguredSpellingOfTheSameAlias(t *testing.T) {
	// OpsLevel resolves aliases case-insensitively, so this names the same component type as "business_unit".
	if got := componentTypeReadWith(types.StringValue("Business_Unit")); !got.Equal(types.StringValue("Business_Unit")) {
		t.Errorf("expected the configured spelling of the same alias to be kept, got %s", got)
	}
}

func TestNewRelationshipDefinitionResourceModel_ReadsAnIDFromTheAPI(t *testing.T) {
	for _, stateValue := range []string{componentTypeID, otherComponentTypeID} {
		if got := componentTypeReadWith(types.StringValue(stateValue)); !got.Equal(types.StringValue(componentTypeID)) {
			t.Errorf("expected an ID in state (%s) to read the API's ID, got %s", stateValue, got)
		}
	}
}

const relationshipDefinitionID = "Z2lkOi8vb3BzbGV2ZWwvUmVsYXRpb25zaGlwRGVmaW5pdGlvbi8x"

// readRelationshipDefinition runs Read with only the ID in state, against a fake API whose
// relationshipDefinition field returns definitionJSON.
func readRelationshipDefinition(t *testing.T, definitionJSON string) resource.ReadResponse {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		if !strings.Contains(string(body), "relationshipDefinition(") {
			t.Errorf("unexpected request: %s", body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"account":{"relationshipDefinition":` + definitionJSON + `}}}`))
	}))
	t.Cleanup(server.Close)

	ctx := context.Background()
	r := opsleveltf.NewRelationshipDefinitionResource().(*opsleveltf.RelationshipDefinitionResource)
	client := opslevelgo.NewGQLClient(
		opslevelgo.SetURL(server.URL+"/LOCAL_TESTING/test"),
		opslevelgo.SetAPIToken("test-token"),
		opslevelgo.SetMaxRetries(0),
	)
	var configureResp resource.ConfigureResponse
	r.Configure(ctx, resource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("unexpected configure errors: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	attributes := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	attributes["id"] = tftypes.NewValue(tftypes.String, relationshipDefinitionID)
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, attributes)}

	resp := resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &resp)
	return resp
}

func TestRelationshipDefinitionResource_ReadRemovesADeletedDefinition(t *testing.T) {
	// The API returns null, with no error, for a definition deleted outside Terraform.
	resp := readRelationshipDefinition(t, `null`)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected the deleted definition to be removed from state, got %s", resp.State.Raw)
	}
}

func TestRelationshipDefinitionResource_ReadKeepsAnExistingDefinition(t *testing.T) {
	resp := readRelationshipDefinition(t, `{"id":"`+relationshipDefinitionID+`","alias":"infrastructure","name":"Infrastructure","description":"",`+
		`"metadata":{"allowedCategories":[],"allowedTypes":["team"],"maxItems":0,"minItems":0},`+
		`"componentType":{"id":"`+componentTypeID+`","aliases":["business_unit"]},"managementRules":[]}`)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("expected the existing definition to stay in state")
	}
	var id types.String
	if diags := resp.State.GetAttribute(context.Background(), path.Root("id"), &id); diags.HasError() {
		t.Fatalf("unable to read id from state: %v", diags.Errors())
	}
	if !id.Equal(types.StringValue(relationshipDefinitionID)) {
		t.Errorf("expected id %s in state, got %s", relationshipDefinitionID, id)
	}
}
