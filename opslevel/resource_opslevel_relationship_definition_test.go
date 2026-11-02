package opslevel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
