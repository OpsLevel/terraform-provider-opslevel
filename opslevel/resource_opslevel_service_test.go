package opslevel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	opslevelgo "github.com/opslevel/opslevel-go/v2026"
)

func TestServiceApiDocSettingsUpdateInput(t *testing.T) {
	testCases := []struct {
		name              string
		plan              ServiceResourceModel
		expectedDocPath   string
		expectedDocSource *opslevelgo.ApiDocumentSourceEnum
	}{
		{
			name: "neither field set",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringNull(),
			},
			expectedDocPath: "",
		},
		{
			name: "push source without path",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringValue(string(opslevelgo.ApiDocumentSourceEnumPush)),
			},
			expectedDocPath:   "",
			expectedDocSource: &opslevelgo.ApiDocumentSourceEnumPush,
		},
		{
			name: "pull source without path",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringValue(string(opslevelgo.ApiDocumentSourceEnumPull)),
			},
			expectedDocPath:   "",
			expectedDocSource: &opslevelgo.ApiDocumentSourceEnumPull,
		},
		{
			name: "path without source",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringValue("openapi.yaml"),
				PreferredApiDocumentSource: types.StringNull(),
			},
			expectedDocPath: "openapi.yaml",
		},
		{
			name: "path with source",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringValue("openapi.yaml"),
				PreferredApiDocumentSource: types.StringValue(string(opslevelgo.ApiDocumentSourceEnumPush)),
			},
			expectedDocPath:   "openapi.yaml",
			expectedDocSource: &opslevelgo.ApiDocumentSourceEnumPush,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			docPath, docSource := serviceApiDocSettingsUpdateInput(testCase.plan)

			if docPath != testCase.expectedDocPath {
				t.Fatalf("expected doc path %q, got %q", testCase.expectedDocPath, docPath)
			}
			if testCase.expectedDocSource == nil {
				if docSource != nil {
					t.Fatalf("expected nil doc source, got %q", *docSource)
				}
				return
			}
			if docSource == nil {
				t.Fatalf("expected doc source %q, got nil", *testCase.expectedDocSource)
			}
			if *docSource != *testCase.expectedDocSource {
				t.Fatalf("expected doc source %q, got %q", *testCase.expectedDocSource, *docSource)
			}
		})
	}
}
