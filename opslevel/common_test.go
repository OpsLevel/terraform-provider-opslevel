package opslevel

import (
	"encoding/json"
	"github.com/opslevel/opslevel-go/v2024"
	"reflect"
	"testing"
)

func compareKey(t *testing.T, m map[string]string, key string, exp string) {
	got := m[key]
	if exp != m[key] {
		t.Errorf("on key '%s' expected value %#v got\n\t%#v", key, exp, got)
	}
}

func TestInterfacesMap(t *testing.T) {
	predicateInputs := `[{"case_sensitive":"true","key":"tags","key_data":"image","type":"ends_with","value":"ID"},{"case_sensitive":"false","key":"name","key_data":"","type":"contains","value":"runner"},{"case_sensitive":"","key":"repository_ids","key_data":"","type":"exists","value":""}]`
	var unread interface{}
	err := json.Unmarshal([]byte(predicateInputs), &unread)
	if err != nil {
		t.Error(err)
	}
	output := interfacesMap(unread)
	if len(output) != 3 {
		t.Errorf("expected resulting interfaces map len to be 3 got %d", len(output))
	}
	if len(output[0]) != 5 {
		t.Errorf("expected len 5 got %d", len(output[0]))
	}
	compareKey(t, output[0], "case_sensitive", "true")
	compareKey(t, output[0], "key", "tags")
	compareKey(t, output[0], "key_data", "image")
	compareKey(t, output[0], "type", "ends_with")
	compareKey(t, output[0], "value", "ID")
	if len(output[1]) != 5 {
		t.Errorf("expected len 5 got %d", len(output[1]))
	}
	compareKey(t, output[1], "case_sensitive", "false")
	compareKey(t, output[1], "key", "name")
	compareKey(t, output[1], "key_data", "")
	compareKey(t, output[1], "type", "contains")
	compareKey(t, output[1], "value", "runner")
	if len(output[2]) != 5 {
		t.Errorf("expected len 5 got %d", len(output[2]))
	}
	compareKey(t, output[2], "key", "repository_ids")
	compareKey(t, output[2], "type", "exists")
	compareKey(t, output[2], "case_sensitive", "")
	compareKey(t, output[2], "key_data", "")
	compareKey(t, output[2], "value", "")
}

func TestExpandFilterPredicates(t *testing.T) {
	predicateInputs := []map[string]string{
		{
			"type":           "ends_with",
			"value":          "ID",
			"key":            "tags",
			"key_data":       "image",
			"case_sensitive": "true",
		},
		{
			"type":           "contains",
			"value":          "runner",
			"key":            "name",
			"case_sensitive": "false",
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
