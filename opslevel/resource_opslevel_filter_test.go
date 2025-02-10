package opslevel_test

import (
	"testing"

	opslevelgo "github.com/opslevel/opslevel-go/v2025"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func TestNewFilterPredicateModel(t *testing.T) {
	apiPredicateEmptyStrings := opslevelgo.FilterPredicate{
		Key:     opslevelgo.PredicateKeyEnumLanguage,
		Type:    opslevelgo.PredicateTypeEnumExists,
		Value:   "",
		KeyData: "",
	}
	apiPredicateNullValues := opslevelgo.FilterPredicate{
		Key:  opslevelgo.PredicateKeyEnumLanguage,
		Type: opslevelgo.PredicateTypeEnumExists,
	}
	predicateEmptyStrings := opsleveltf.NewFilterPredicateModel(&apiPredicateEmptyStrings)
	predicateNullValues := opsleveltf.NewFilterPredicateModel(&apiPredicateNullValues)
	if predicateEmptyStrings != predicateNullValues {
		t.Errorf("expected new FilterPredicateModels with empty strings and null values to be equal")
	}
}
