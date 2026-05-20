package opslevel_test

import (
	"sort"
	"testing"

	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

func TestDiffCheckIds_NoChange(t *testing.T) {
	state := map[string]bool{"a": true, "b": true}
	plan := map[string]bool{"a": true, "b": true}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	if len(toAdd) != 0 {
		t.Errorf("expected no additions, got %v", toAdd)
	}
	if len(toRemove) != 0 {
		t.Errorf("expected no removals, got %v", toRemove)
	}
}

func TestDiffCheckIds_AddOnly(t *testing.T) {
	state := map[string]bool{"a": true}
	plan := map[string]bool{"a": true, "b": true, "c": true}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	sort.Strings(toAdd)
	if len(toAdd) != 2 || toAdd[0] != "b" || toAdd[1] != "c" {
		t.Errorf("expected [b c], got %v", toAdd)
	}
	if len(toRemove) != 0 {
		t.Errorf("expected no removals, got %v", toRemove)
	}
}

func TestDiffCheckIds_RemoveOnly(t *testing.T) {
	state := map[string]bool{"a": true, "b": true, "c": true}
	plan := map[string]bool{"a": true}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	sort.Strings(toRemove)
	if len(toAdd) != 0 {
		t.Errorf("expected no additions, got %v", toAdd)
	}
	if len(toRemove) != 2 || toRemove[0] != "b" || toRemove[1] != "c" {
		t.Errorf("expected [b c], got %v", toRemove)
	}
}

func TestDiffCheckIds_AddAndRemove(t *testing.T) {
	state := map[string]bool{"a": true, "b": true}
	plan := map[string]bool{"b": true, "c": true}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	if len(toAdd) != 1 || toAdd[0] != "c" {
		t.Errorf("expected [c], got %v", toAdd)
	}
	if len(toRemove) != 1 || toRemove[0] != "a" {
		t.Errorf("expected [a], got %v", toRemove)
	}
}

func TestDiffCheckIds_EmptyPlan(t *testing.T) {
	state := map[string]bool{"a": true, "b": true}
	plan := map[string]bool{}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	sort.Strings(toRemove)
	if len(toAdd) != 0 {
		t.Errorf("expected no additions, got %v", toAdd)
	}
	if len(toRemove) != 2 || toRemove[0] != "a" || toRemove[1] != "b" {
		t.Errorf("expected [a b], got %v", toRemove)
	}
}

func TestDiffCheckIds_EmptyState(t *testing.T) {
	state := map[string]bool{}
	plan := map[string]bool{"a": true, "b": true}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	sort.Strings(toAdd)
	if len(toAdd) != 2 || toAdd[0] != "a" || toAdd[1] != "b" {
		t.Errorf("expected [a b], got %v", toAdd)
	}
	if len(toRemove) != 0 {
		t.Errorf("expected no removals, got %v", toRemove)
	}
}

func TestDiffCheckIds_BothEmpty(t *testing.T) {
	state := map[string]bool{}
	plan := map[string]bool{}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	if len(toAdd) != 0 {
		t.Errorf("expected no additions, got %v", toAdd)
	}
	if len(toRemove) != 0 {
		t.Errorf("expected no removals, got %v", toRemove)
	}
}

func TestDiffCheckIds_NilState(t *testing.T) {
	var state map[string]bool
	plan := map[string]bool{"a": true}
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	if len(toAdd) != 1 || toAdd[0] != "a" {
		t.Errorf("expected [a], got %v", toAdd)
	}
	if len(toRemove) != 0 {
		t.Errorf("expected no removals, got %v", toRemove)
	}
}

func TestDiffCheckIds_NilPlan(t *testing.T) {
	state := map[string]bool{"a": true}
	var plan map[string]bool
	toAdd, toRemove := opsleveltf.DiffCheckIds(state, plan)
	if len(toAdd) != 0 {
		t.Errorf("expected no additions, got %v", toAdd)
	}
	if len(toRemove) != 1 || toRemove[0] != "a" {
		t.Errorf("expected [a], got %v", toRemove)
	}
}
