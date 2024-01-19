package opslevel

import (
	"encoding/json"
	"github.com/opslevel/opslevel-go/v2024"
	"reflect"
	"testing"
)

func compareKey(t *testing.T, m map[string]interface{}, key string, exp interface{}) {
	got := m[key]
	if !reflect.DeepEqual(exp, m[key]) {
		t.Errorf("on key '%s' expected value %#v got\n\t%#v", key, exp, got)
	}
}

func TestInterfacesMaps(t *testing.T) {
	predicateInputs := `[{"case_insensitive":false,"case_sensitive":true,"key":"tags","key_data":"image","type":"ends_with","value":"ID"},{"case_insensitive":true,"case_sensitive":false,"key":"name","key_data":"","type":"contains","value":"runner"},{"case_insensitive":false,"case_sensitive":false,"key":"repository_ids","key_data":"","type":"exists","value":""}]`
	var unread interface{}
	err := json.Unmarshal([]byte(predicateInputs), &unread)
	if err != nil {
		t.Error(err)
	}
	output := interfacesMaps(unread)
	if len(output) != 3 {
		t.Errorf("expected resulting interfaces map len to be 3 got %d", len(output))
	}
	if len(output[0]) != 6 {
		t.Errorf("expected len 6 got %d", len(output[0]))
	}
	compareKey(t, output[0], "case_sensitive", true)
	compareKey(t, output[0], "case_insensitive", false)
	compareKey(t, output[0], "key", "tags")
	compareKey(t, output[0], "key_data", "image")
	compareKey(t, output[0], "type", "ends_with")
	compareKey(t, output[0], "value", "ID")
	if len(output[1]) != 6 {
		t.Errorf("expected len 6 got %d", len(output[1]))
	}
	compareKey(t, output[1], "case_sensitive", false)
	compareKey(t, output[1], "case_insensitive", true)
	compareKey(t, output[1], "key", "name")
	compareKey(t, output[1], "key_data", "")
	compareKey(t, output[1], "type", "contains")
	compareKey(t, output[1], "value", "runner")
	if len(output[2]) != 6 {
		t.Errorf("expected len 6 got %d", len(output[2]))
	}
	compareKey(t, output[2], "case_sensitive", false)
	compareKey(t, output[2], "case_insensitive", false)
	compareKey(t, output[2], "key", "repository_ids")
	compareKey(t, output[2], "type", "exists")
	compareKey(t, output[2], "key_data", "")
	compareKey(t, output[2], "value", "")
}

func TestExpandFilterPredicates(t *testing.T) {
	predicateInputs := []map[string]interface{}{
		{
			"type":             "ends_with",
			"value":            "ID",
			"key":              "tags",
			"key_data":         "image",
			"case_sensitive":   true,
			"case_insensitive": false,
		},
		{
			"type":             "contains",
			"value":            "runner",
			"key":              "name",
			"case_sensitive":   false,
			"case_insensitive": true,
		},
		{
			"type": "exists",
			"key":  "repository_ids",
		},
	}
	expectedInputs := []opslevel.FilterPredicateInput{
		{
			Type:          opslevel.PredicateTypeEnumEndsWith,
			Value:         opslevel.RefTo("ID"),
			Key:           opslevel.PredicateKeyEnumTags,
			KeyData:       opslevel.RefOf("image"),
			CaseSensitive: opslevel.RefOf(true),
		},
		{
			Type:  opslevel.PredicateTypeEnumContains,
			Value: opslevel.RefTo("runner"),
			Key:   opslevel.PredicateKeyEnumName,
			// KeyData here should be read as nil
			CaseSensitive: opslevel.RefOf(false),
		},
		{
			Type: opslevel.PredicateTypeEnumExists,
			Key:  opslevel.PredicateKeyEnumRepositoryIDs,
			// CaseSensitive here should be read as nil
		},
	}

	outputs := *expandFilterPredicateInputs(predicateInputs)
	for i := range outputs {
		if !reflect.DeepEqual(outputs[i], expectedInputs[i]) {
			t.Errorf("expected %#v\n\tgot %#v", expectedInputs[i], outputs[i])
		}
	}
}

func TestFlattenFilterPredicates(t *testing.T) {
	predicates := []opslevel.FilterPredicate{
		{
			Type:          opslevel.PredicateTypeEnumEndsWith,
			Value:         "ID",
			Key:           opslevel.PredicateKeyEnumTags,
			KeyData:       "image",
			CaseSensitive: opslevel.RefOf(true),
		},
		{
			Type:          opslevel.PredicateTypeEnumContains,
			Value:         "runner",
			Key:           opslevel.PredicateKeyEnumName,
			CaseSensitive: opslevel.RefOf(false),
		},
		{
			Type: opslevel.PredicateTypeEnumExists,
			Key:  opslevel.PredicateKeyEnumRepositoryIDs,
		},
	}
	// tf provider version can't differentiate between nil vs zero val on inputs
	// means required fields will be parsed as zero (or nil if they are pointers)
	expectedInputs := []map[string]interface{}{
		{
			"type":             "ends_with",
			"value":            "ID",
			"key":              "tags",
			"key_data":         "image",
			"case_sensitive":   true,
			"case_insensitive": false,
		},
		{
			"type":             "contains",
			"value":            "runner",
			"key":              "name",
			"case_sensitive":   false,
			"case_insensitive": true,
			"key_data":         "",
		},
		{
			"type":             "exists",
			"key":              "repository_ids",
			"value":            "",
			"key_data":         "",
			"case_sensitive":   false,
			"case_insensitive": false,
		},
	}

	outputs := flattenFilterPredicates(predicates)
	for i := range outputs {
		if !reflect.DeepEqual(expectedInputs[i], outputs[i]) {
			t.Errorf("expected %#v\n\tgot %#v", expectedInputs[i], outputs[i])
		}
	}
}
