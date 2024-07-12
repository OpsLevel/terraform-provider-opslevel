package opslevel

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// Returns resourceValue wrapped in types.StringValue or types.StringNull if modelValue is null
func StringValueFromResourceAndModelField(resourceValue string, modelValue basetypes.StringValue) basetypes.StringValue {
	if resourceValue == "" && modelValue.IsNull() {
		return types.StringNull()
	}
	return types.StringValue(unquote(resourceValue))
}

// Returns value wrapped in a types.StringValue, or types.StringNull if blank
func OptionalStringValue(value string) basetypes.StringValue {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(unquote(value))
}

// Returns value from config as a string OR null if the value is not set/explicitly set to null (supports empty strings)
func NullableStringConfigValue(s types.String) *opslevel.Nullable[string] {
	if s.IsNull() {
		return opslevel.NewNull[string]()
	}
	return opslevel.NewNullableFrom(s.ValueString())
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
func OptionalStringListValue(values []string) basetypes.ListValue {
	if len(values) == 0 {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}

	return types.ListValueMust(types.StringType, elems)
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

// Converts a basetypes.SetValue to a []string
func SetValueToStringSlice(ctx context.Context, setValue basetypes.SetValue) ([]string, diag.Diagnostics) {
	dataAsSlice := []string{}
	if setValue.IsNull() {
		return dataAsSlice, nil
	}
	diags := setValue.ElementsAs(ctx, &dataAsSlice, true)
	return dataAsSlice, diags
}

// Converts a []string to a basetypes.SetValue
func StringSliceToSetValue(values []string) basetypes.SetValue {
	result := []attr.Value{}
	for _, value := range values {
		result = append(result, types.StringValue(value))
	}
	return types.SetValueMust(types.StringType, result)
}

func TagSetValueToTagSlice(ctx context.Context, setValue basetypes.SetValue) ([]opslevel.Tag, diag.Diagnostics) {
	tagSlice := []opslevel.Tag{}
	if setValue.IsNull() {
		return tagSlice, nil
	}
	tagsAsStringSlice, diags := SetValueToStringSlice(ctx, setValue)
	if diags.HasError() {
		return tagSlice, diags
	}
	for _, tag := range tagsAsStringSlice {
		if hasTagFormat(tag) {
			parts := strings.Split(tag, ":")
			tagSlice = append(tagSlice, opslevel.Tag{Key: parts[0], Value: parts[1]})
		} else {
			diags.AddWarning("Invalid tag format", tag)
		}
	}
	return tagSlice, diags
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

// Converts a basetypes.ObjectValue to a PredicateModel
func PredicateObjectToModel(ctx context.Context, predicateObj basetypes.ObjectValue) (PredicateModel, diag.Diagnostics) {
	var predicateModel PredicateModel

	objOptions := basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}
	diags := predicateObj.As(ctx, &predicateModel, objOptions)
	if predicateModel.Value.ValueString() == "" {
		predicateModel.Value = types.StringNull()
	}
	return predicateModel, diags
}

// asID converts a types.String to an opslevel.ID
func asID(input types.String) opslevel.ID {
	return opslevel.ID(input.ValueString())
}
