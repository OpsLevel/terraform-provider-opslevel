package opslevel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	opslevelgo "github.com/opslevel/opslevel-go/v2026"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func relationshipDefinition() opslevelgo.RelationshipDefinitionType {
	return opslevelgo.RelationshipDefinitionType{
		Id:            opslevelgo.ID("Z2lkOi8vb3BzbGV2ZWwvUmVsYXRpb25zaGlwRGVmaW5pdGlvbi8x"),
		Alias:         "infrastructure",
		Name:          "Infrastructure",
		ComponentType: opslevelgo.ComponentTypeId{Id: opslevelgo.ID("Z2lkOi8vb3BzbGV2ZWwvQ29tcG9uZW50VHlwZS8x"), Aliases: []string{"business_unit"}},
	}
}

func TestNewRelationshipDefinitionResourceModel_ImportTakesComponentTypeFromAPI(t *testing.T) {
	// ImportState sets only the ID, so Read starts from a model with no component_type.
	model := opsleveltf.NewRelationshipDefinitionResourceModel(relationshipDefinition(), opsleveltf.RelationshipDefinitionResourceModel{})

	if !model.ComponentType.Equal(types.StringValue("business_unit")) {
		t.Errorf("expected an imported definition to read component_type from the API, got %s", model.ComponentType)
	}
}

func TestNewRelationshipDefinitionResourceModel_KeepsConfiguredComponentType(t *testing.T) {
	for _, configured := range []string{"business_unit", "Z2lkOi8vb3BzbGV2ZWwvQ29tcG9uZW50VHlwZS8x"} {
		model := opsleveltf.NewRelationshipDefinitionResourceModel(
			relationshipDefinition(),
			opsleveltf.RelationshipDefinitionResourceModel{ComponentType: types.StringValue(configured)},
		)

		if !model.ComponentType.Equal(types.StringValue(configured)) {
			t.Errorf("expected the configured component_type %q to be kept, got %s", configured, model.ComponentType)
		}
	}
}
