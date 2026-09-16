package opslevel_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	opslevelgo "github.com/opslevel/opslevel-go/v2026"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

// boolPtr is a small helper for constructing *bool fields in test fixtures.
func boolPtr(b bool) *bool {
	return &b
}

// buildStatePredicateList builds a non-null types.List of one predicate suitable
// for use as the Predicate field on FilterResourceModel in tests.
func buildStatePredicateList(t *testing.T, p opsleveltf.FilterPredicateModel) types.List {
	t.Helper()
	list, diags := types.ListValueFrom(
		context.Background(),
		types.ObjectType{AttrTypes: opsleveltf.FilterPredicateType},
		[]opsleveltf.FilterPredicateModel{p},
	)
	if diags.HasError() {
		t.Fatalf("constructing predicate list: %v", diags)
	}
	return list
}

// extractFirstPredicate decodes the first predicate from a result model's List.
func extractFirstPredicate(t *testing.T, result opsleveltf.FilterResourceModel) opsleveltf.FilterPredicateModel {
	t.Helper()
	var predicates []opsleveltf.FilterPredicateModel
	diags := result.Predicate.ElementsAs(context.Background(), &predicates, false)
	if diags.HasError() {
		t.Fatalf("decoding result predicates: %v", diags)
	}
	if len(predicates) != 1 {
		t.Fatalf("expected 1 predicate, got %d", len(predicates))
	}
	return predicates[0]
}

// stubFilter wraps a single API-returned predicate in a minimal opslevel.Filter.
func stubFilter(pred opslevelgo.FilterPredicate) opslevelgo.Filter {
	return opslevelgo.Filter{
		Predicates: []opslevelgo.FilterPredicate{pred},
	}
}

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

// TestNewFilterResourceModel_NoCrashWhenStateHasCaseSensitiveButAPIReturnsNil
// is the regression test for the nil-pointer dereference that occurs when
// Terraform state holds a `case_sensitive` boolean for a predicate but the
// OpsLevel API returns nil for that predicate's CaseSensitive (which it does
// for `exists`, `matches` against an ID, and other predicate types that
// don't take a case-sensitivity argument).
//
// Pre-fix: this panics with "invalid memory address or nil pointer dereference"
// inside NewFilterResourceModel.
// Post-fix: the resulting predicate has case_sensitive=null.
func TestNewFilterResourceModel_NoCrashWhenStateHasCaseSensitiveButAPIReturnsNil(t *testing.T) {
	ctx := context.Background()

	statePred := opsleveltf.FilterPredicateModel{
		CaseSensitive: types.BoolValue(true),
		Key:           types.StringValue(string(opslevelgo.PredicateKeyEnumRepositoryIDs)),
		Type:          types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
	}
	given := opsleveltf.FilterResourceModel{
		Predicate: buildStatePredicateList(t, statePred),
	}

	apiFilter := stubFilter(opslevelgo.FilterPredicate{
		Key:           opslevelgo.PredicateKeyEnumRepositoryIDs,
		Type:          opslevelgo.PredicateTypeEnumExists,
		CaseSensitive: nil,
	})

	result, diags := opsleveltf.NewFilterResourceModel(ctx, apiFilter, given)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	got := extractFirstPredicate(t, result)
	if !got.CaseSensitive.IsNull() {
		t.Errorf("case_sensitive: expected null when API returns nil, got %v", got.CaseSensitive)
	}
}

// TestNewFilterResourceModel_NoCrashWhenStateHasCaseInsensitiveButAPIReturnsNil
// covers the symmetric bug on the case_insensitive branch — the same nil deref
// pattern existed on the second `if` block.
func TestNewFilterResourceModel_NoCrashWhenStateHasCaseInsensitiveButAPIReturnsNil(t *testing.T) {
	ctx := context.Background()

	statePred := opsleveltf.FilterPredicateModel{
		CaseInsensitive: types.BoolValue(true),
		Key:             types.StringValue(string(opslevelgo.PredicateKeyEnumRepositoryIDs)),
		Type:            types.StringValue(string(opslevelgo.PredicateTypeEnumExists)),
	}
	given := opsleveltf.FilterResourceModel{
		Predicate: buildStatePredicateList(t, statePred),
	}

	apiFilter := stubFilter(opslevelgo.FilterPredicate{
		Key:           opslevelgo.PredicateKeyEnumRepositoryIDs,
		Type:          opslevelgo.PredicateTypeEnumExists,
		CaseSensitive: nil,
	})

	result, diags := opsleveltf.NewFilterResourceModel(ctx, apiFilter, given)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	got := extractFirstPredicate(t, result)
	if !got.CaseInsensitive.IsNull() {
		t.Errorf("case_insensitive: expected null when API returns nil, got %v", got.CaseInsensitive)
	}
}

// TestNewFilterResourceModel_PreservesCaseSensitiveWhenAPISetsIt covers the
// happy path: state and API agree, value is preserved through the round-trip.
func TestNewFilterResourceModel_PreservesCaseSensitiveWhenAPISetsIt(t *testing.T) {
	ctx := context.Background()

	statePred := opsleveltf.FilterPredicateModel{
		CaseSensitive: types.BoolValue(true),
		Key:           types.StringValue(string(opslevelgo.PredicateKeyEnumLanguage)),
		Type:          types.StringValue(string(opslevelgo.PredicateTypeEnumEquals)),
		Value:         types.StringValue("Go"),
	}
	given := opsleveltf.FilterResourceModel{
		Predicate: buildStatePredicateList(t, statePred),
	}

	apiFilter := stubFilter(opslevelgo.FilterPredicate{
		Key:           opslevelgo.PredicateKeyEnumLanguage,
		Type:          opslevelgo.PredicateTypeEnumEquals,
		Value:         "Go",
		CaseSensitive: boolPtr(true),
	})

	result, diags := opsleveltf.NewFilterResourceModel(ctx, apiFilter, given)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	got := extractFirstPredicate(t, result)
	if got.CaseSensitive.IsNull() || !got.CaseSensitive.ValueBool() {
		t.Errorf("case_sensitive: expected true, got %v", got.CaseSensitive)
	}
}

// TestNewFilterResourceModel_EmptyStateWithNilAPICaseSensitive covers the
// fresh-import path: no state predicates, API returns nil CaseSensitive.
// This already worked pre-fix, but is worth pinning because it's the path
// the documented workaround (state-rm + re-import) exercises.
func TestNewFilterResourceModel_EmptyStateWithNilAPICaseSensitive(t *testing.T) {
	ctx := context.Background()

	given := opsleveltf.FilterResourceModel{
		Predicate: types.ListNull(types.ObjectType{AttrTypes: opsleveltf.FilterPredicateType}),
	}

	apiFilter := stubFilter(opslevelgo.FilterPredicate{
		Key:           opslevelgo.PredicateKeyEnumRepositoryIDs,
		Type:          opslevelgo.PredicateTypeEnumExists,
		CaseSensitive: nil,
	})

	result, diags := opsleveltf.NewFilterResourceModel(ctx, apiFilter, given)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	got := extractFirstPredicate(t, result)
	if !got.CaseSensitive.IsNull() {
		t.Errorf("case_sensitive: expected null on fresh import with nil API, got %v", got.CaseSensitive)
	}
}
