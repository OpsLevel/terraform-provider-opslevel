package opslevel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	opslevelgo "github.com/opslevel/opslevel-go/v2025"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func TestExtractFilterPredicateModel(t *testing.T) {
	// The schema ensures empty strings are forbidden, should be treated same as null
	var (
		apiPredicateEmptyStrings = opslevelgo.FilterPredicate{
			Key:     opslevelgo.PredicateKeyEnumLanguage,
			Type:    opslevelgo.PredicateTypeEnumExists,
			Value:   "",
			KeyData: "",
		}
		apiPredicateNullValues = opslevelgo.FilterPredicate{
			Key:  opslevelgo.PredicateKeyEnumLanguage,
			Type: opslevelgo.PredicateTypeEnumExists,
		}
		tfPredicateEmptyStrings = opsleveltf.FilterPredicateModel{
			Key:     types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
			KeyData: types.StringValue(""),
			Type:    types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
			Value:   types.StringValue(""),
		}
		tfPredicateNoValues   = opsleveltf.FilterPredicateModel{}
		tfPredicateNullValues = opsleveltf.FilterPredicateModel{
			Key:  types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
			Type: types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
		}
	)
	possibleFilterPredicateModels := []opsleveltf.FilterPredicateModel{tfPredicateEmptyStrings, tfPredicateNullValues}

	foundModelEmptyStrings := opsleveltf.ExtractFilterPredicateModel(&apiPredicateEmptyStrings, possibleFilterPredicateModels)
	if foundModelEmptyStrings == tfPredicateNoValues {
		t.Errorf("unexpectedly found filter predicate model with no values")
	}
	if foundModelEmptyStrings == tfPredicateEmptyStrings {
		t.Errorf("unexpectedly found filter predicate model with no values")
	}

	foundNullModel := opsleveltf.ExtractFilterPredicateModel(&apiPredicateNullValues, possibleFilterPredicateModels)
	if foundNullModel != foundModelEmptyStrings {
		t.Errorf("unexpectedly found filter predicate model with null string values")
	}
}

func TestExtractFilterPredicateModelNoValues(t *testing.T) {
	var emptyPredicate opsleveltf.FilterPredicateModel
	var models []opsleveltf.FilterPredicateModel

	foundPredicate := opsleveltf.ExtractFilterPredicateModel(nil, models)
	if emptyPredicate != foundPredicate {
		t.Errorf("expected FilterPredicateModel from ExtractFilterPredicateModel to have no values")
	}
}
