package opslevel

import (
	"github.com/kr/pretty"
	"github.com/opslevel/opslevel-go/v2024"
	"reflect"
	"testing"
)

func TestExpandFilterPredicatesReadsBool(t *testing.T) {
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
			Type:          opslevel.PredicateTypeEnumContains,
			Value:         opslevel.RefTo("runner"),
			Key:           opslevel.PredicateKeyEnumName,
			CaseSensitive: opslevel.RefOf(false),
		},
		{
			Type: opslevel.PredicateTypeEnumExists,
			Key:  opslevel.PredicateKeyEnumRepositoryIDs,
		},
	}

	outputs := *expandFilterPredicateInputs(predicateInputs)
	for i := range outputs {
		if !reflect.DeepEqual(outputs[i], expectedInputs[i]) {
			pretty.Print(outputs[i])
			pretty.Print(expectedInputs[i])
			t.Error("unexpected predicate input")
		}
	}
}
