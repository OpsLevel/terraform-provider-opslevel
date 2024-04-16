package opslevel

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opslevel/opslevel-go/v2024"
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

// Returns value wrapped in a types.BoolValue
func RequiredBoolValue(value bool) basetypes.BoolValue {
	return types.BoolValue(value)
}

// Returns value wrapped in a types.BoolValue, or types.BoolNull if blank
func OptionalBoolValue(value *bool) basetypes.BoolValue {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

// Returns value wrapped in a types.Int64Value
func RequiredIntValue(value int) basetypes.Int64Value {
	return types.Int64Value(int64(value))
}

// Returns value wrapped in a types.StringValue, or types.ListNull if blank
func OptionalStringListValue(ctx context.Context, value []string) (basetypes.ListValue, diag.Diagnostics) {
	if len(value) == 0 {
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
	return opslevel.ID(input.ValueString())
}
