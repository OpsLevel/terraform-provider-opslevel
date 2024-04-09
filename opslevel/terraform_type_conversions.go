package opslevel

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
	"golang.org/x/net/context"
)

// Returns value wrapped in a types.StringValue, even if blank string given
func RequiredStringValue(value string) basetypes.StringValue {
	return types.StringValue(unquote(value))
}

// Returns value wrapped in a types.StringValue, or types.StringNull if blank
func OptionalStringValue(value string) basetypes.StringValue {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(unquote(value))
}

// Syntactic sugar for OptionalStringValue
func ComputedStringValue(value string) basetypes.StringValue {
	return OptionalStringValue(value)
}

// Returns value wrapped in a types.StringValue, or types.ListNull if blank
//
// NOTE: an empty list is not the same as 'null'
func OptionalStringListValue(ctx context.Context, value []string) (basetypes.ListValue, diag.Diagnostics) {
	if value == nil {
		return types.ListNull(types.StringType), diag.Diagnostics{}
	}
	return types.ListValueFrom(ctx, types.StringType, value)
}

// unquotes unwanted quotes from strings in maps, returns original value in most cases
func unquote(value string) string {
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		if unquotedValue, err := strconv.Unquote(value); err == nil {
			return unquotedValue
		}
	}
	return value
}

// Converts a basetypes.ListValue to a []string
func ListValueToStringSlice(ctx context.Context, listValue basetypes.ListValue) ([]string, diag.Diagnostics) {
	dataAsSlice := []string{}
	if listValue.IsNull() {
		return dataAsSlice, nil
	}
	diags := listValue.ElementsAs(ctx, &dataAsSlice, true)
	return dataAsSlice, diags
}

// Converts a basetypes.MapValue to an opslevel.JSON
func MapValueToOpslevelJson(ctx context.Context, mapValue basetypes.MapValue) (opslevel.JSON, diag.Diagnostics) {
	mapAsJson := opslevel.JSON{}
	stringMap := map[string]string{}

	diags := mapValue.ElementsAs(ctx, &stringMap, false)
	for k, v := range stringMap {
		mapAsJson[k] = v
	}
	return mapAsJson, diags
}

// asID converts a types.String to an opslevel.ID
func asID(input types.String) opslevel.ID {
	return *opslevel.NewID(input.ValueString())
}

// asISO8601 convert timetypes.RFC3339 to opslevel go's iso8601 time
func asISO8601(input timetypes.RFC3339) (*iso8601.Time, diag.Diagnostics) {
	t, diags := input.ValueRFC3339Time()
	return &iso8601.Time{Time: t}, diags
}
