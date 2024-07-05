package opslevel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	opslevelgo "github.com/opslevel/opslevel-go/v2024"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func TestOpslevelFilterPredicateToObjectValueNullObject(t *testing.T) {
	nullObject := opsleveltf.OpslevelFilterPredicateToObjectValue(nil, nil)
	if !nullObject.Equal(types.ObjectNull(opsleveltf.FilterPredicateType)) {
		t.Fatal("expected null object")
	}
}

func TestOpslevelFilterPredicateToObjectValueAllFieldsSet(t *testing.T) {
	// FilterPredicate.CaseSensitive omitted intentionally, types.BoolNull set in object
	filterPredicate := opslevelgo.FilterPredicate{
		Key:     opslevelgo.PredicateKeyEnumName,
		KeyData: "test data",
		Type:    opslevelgo.PredicateTypeEnumExists,
		Value:   "test value",
	}
	predicateObj := opsleveltf.OpslevelFilterPredicateToObjectValue(nil, &filterPredicate)
	attrs := predicateObj.Attributes()

	if !attrs["case_insensitive"].Equal(types.BoolNull()) {
		t.Logf("expected filter predicate 'case_insensitive' to be BoolNull")
		t.Fail()
	}
	if !attrs["case_sensitive"].Equal(types.BoolNull()) {
		t.Logf("expected filter predicate 'case_sensitive' to be BoolNull")
		t.Fail()
	}
	if !attrs["key"].Equal(types.StringValue(string(filterPredicate.Key))) {
		t.Logf("expected filter predicate 'key' to be StringValue of %s", string(filterPredicate.Key))
		t.Fail()
	}
	if !attrs["key_data"].Equal(types.StringValue(filterPredicate.KeyData)) {
		t.Logf("expected filter predicate 'key_data' to be StringValue of %s", string(filterPredicate.KeyData))
		t.Fail()
	}
	if !attrs["type"].Equal(types.StringValue(string(filterPredicate.Type))) {
		t.Logf("expected filter predicate 'type' to be StringValue of %s", string(filterPredicate.Type))
		t.Fail()
	}
	if !attrs["value"].Equal(types.StringValue(filterPredicate.Value)) {
		t.Logf("expected filter predicate 'value' to be StringValue of %s", string(filterPredicate.Value))
		t.Fail()
	}
}

func TestOpslevelFilterPredicateToObjectValueMinimal(t *testing.T) {
	// FilterPredicate.CaseSensitive omitted intentionally, types.BoolNull set in object
	filterPredicates := []opslevelgo.FilterPredicate{
		{
			Key:  opslevelgo.PredicateKeyEnumAliases,
			Type: opslevelgo.PredicateTypeEnumDoesNotExist,
		},
		{
			Key:     opslevelgo.PredicateKeyEnumName,
			KeyData: "",
			Type:    opslevelgo.PredicateTypeEnumExists,
			Value:   "",
		},
	}
	for _, filterPredicate := range filterPredicates {
		// Validate filter predicates to ensure realistic test data
		if err := filterPredicate.Validate(); err != nil {
			t.Errorf(err.Error())
		}
		predicateObj := opsleveltf.OpslevelFilterPredicateToObjectValue(nil, &filterPredicate)
		attrs := predicateObj.Attributes()

		if !attrs["case_insensitive"].Equal(types.BoolNull()) {
			t.Logf("expected filter predicate 'case_insensitive' to be BoolNull")
			t.Fail()
		}
		if !attrs["case_sensitive"].Equal(types.BoolNull()) {
			t.Logf("expected filter predicate 'case_sensitive' to be BoolNull")
			t.Fail()
		}
		if !attrs["key"].Equal(types.StringValue(string(filterPredicate.Key))) {
			t.Logf("expected filter predicate 'key' to be StringValue of %s", string(filterPredicate.Key))
			t.Fail()
		}
		if !attrs["key_data"].Equal(types.StringNull()) {
			t.Log("expected filter predicate 'key_data' to be StringNull")
			t.Fail()
		}
		if !attrs["type"].Equal(types.StringValue(string(filterPredicate.Type))) {
			t.Logf("expected filter predicate 'type' to be StringValue of '%s'", string(filterPredicate.Type))
			t.Fail()
		}
		if !attrs["value"].Equal(types.StringNull()) {
			t.Log("expected filter predicate 'value' to be StringNull")
			t.Fail()
		}
	}
}

func TestExtractFilterPredicateModel(t *testing.T) {
	var (
		blankModel = opsleveltf.FilterPredicateModel{}
		// The schema ensures empty strings are forbidden, should be treated same as null
		emptyValues = opsleveltf.FilterPredicateModel{
			Key:     types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
			KeyData: types.StringValue(""),
			Type:    types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
			Value:   types.StringValue(""),
		}
		nullValues = opsleveltf.FilterPredicateModel{
			Key:  types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
			Type: types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
		}
	)
	possibleFilterPredicateModels := []opsleveltf.FilterPredicateModel{emptyValues, nullValues}

	attrsWithEmptyString := map[string]attr.Value{
		"case_insensitive": types.BoolNull(),
		"case_sensitive":   types.BoolNull(),
		"key":              types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
		"key_data":         types.StringValue(""),
		"type":             types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
		"value":            types.StringValue(""),
	}
	attrsWithNullStrings := map[string]attr.Value{
		"case_insensitive": types.BoolNull(),
		"case_sensitive":   types.BoolNull(),
		"key":              types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
		"key_data":         types.StringNull(),
		"type":             types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
		"value":            types.StringNull(),
	}

	foundModel, diags := opsleveltf.ExtractFilterPredicateModel(nil, attrsWithEmptyString, possibleFilterPredicateModels)
	if !diags.HasError() && foundModel == blankModel {
		t.Errorf("unexpectedly found filter predicate model with empty string values")
	}

	foundModel, diags = opsleveltf.ExtractFilterPredicateModel(nil, attrsWithNullStrings, possibleFilterPredicateModels)
	if diags.HasError() && foundModel != blankModel {
		t.Errorf("failed to find filter predicate model with null string values")
	}
}
