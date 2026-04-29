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
		state             *ServiceResourceModel
		expectedUpdate    bool
		expectedDocPath   string
		expectedDocSource *opslevelgo.ApiDocumentSourceEnum
	}{
		{
			name: "create ignores unmanaged settings",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringNull(),
			},
			expectedUpdate: false,
		},
		{
			name: "create updates push source without path",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringValue(string(opslevelgo.ApiDocumentSourceEnumPush)),
			},
			expectedUpdate:    true,
			expectedDocPath:   "",
			expectedDocSource: &opslevelgo.ApiDocumentSourceEnumPush,
		},
		{
			name: "create updates pull source without path",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringValue(string(opslevelgo.ApiDocumentSourceEnumPull)),
			},
			expectedUpdate:    true,
			expectedDocPath:   "",
			expectedDocSource: &opslevelgo.ApiDocumentSourceEnumPull,
		},
		{
			name: "create updates path without source",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringValue("openapi.yaml"),
				PreferredApiDocumentSource: types.StringNull(),
			},
			expectedUpdate:  true,
			expectedDocPath: "openapi.yaml",
		},
		{
			name: "update clears managed source",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringNull(),
			},
			state: &ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringValue(string(opslevelgo.ApiDocumentSourceEnumPull)),
			},
			expectedUpdate:  true,
			expectedDocPath: "",
		},
		{
			name: "update clears managed path",
			plan: ServiceResourceModel{
				ApiDocumentPath:            types.StringNull(),
				PreferredApiDocumentSource: types.StringNull(),
			},
			state: &ServiceResourceModel{
				ApiDocumentPath:            types.StringValue("openapi.yaml"),
				PreferredApiDocumentSource: types.StringNull(),
			},
			expectedUpdate:  true,
			expectedDocPath: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			shouldUpdate, docPath, docSource := serviceApiDocSettingsUpdateInput(testCase.plan, testCase.state)

			if shouldUpdate != testCase.expectedUpdate {
				t.Fatalf("expected update %t, got %t", testCase.expectedUpdate, shouldUpdate)
			}
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
