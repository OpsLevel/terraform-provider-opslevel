package opslevel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/opslevel/opslevel-go/v2025"
)

// OpsLevel ID String Validator
type idStringValidator struct{}

func (v idStringValidator) Description(_ context.Context) string {
	return "id expected to be a string starting with 'Z2lkOi8v' (which is 'gid://' encoded in base64)"
}

func (v idStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v idStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if request.ConfigValue.ValueString() == "" {
		response.Diagnostics.AddAttributeError(request.Path, "Config error", "expected a valid id but given an empty string. please set to 'null' or an id starting with 'Z2lkOi8v'")
		return
	}

	value := request.ConfigValue.ValueString()
	if !opslevel.IsID(value) {
		response.Diagnostics.AddAttributeError(request.Path, "Config error", fmt.Sprintf("expected a valid id. id should start with 'Z2lkOi8v'. given '%s'", value))
	}
}

func IdStringValidator() validator.String {
	return idStringValidator{}
}

// jsonStringValidator accepts any valid JSON (does not have to be an object), but not null and unknown
type jsonStringValidator struct{}

func (v jsonStringValidator) Description(_ context.Context) string {
	return "field expected to be valid JSON"
}

func (v jsonStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !json.Valid([]byte(value)) {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected valid JSON. %s was set to `%s`", request.Path, value))
	}
}

func JsonStringValidator() validator.String {
	return jsonStringValidator{}
}

// jsonObjectValidator accepts any valid JSON object
type jsonObjectValidator struct{}

func (v jsonObjectValidator) Description(_ context.Context) string {
	return "field expected to be valid JSON object"
}

func (v jsonObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonObjectValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	result := make(map[string]any)
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected a valid JSON object. '%s' was set to `%s`", request.Path, value))
		return
	}
}

func JsonObjectValidator() validator.String {
	return jsonObjectValidator{}
}

// jsonHasNameKeyValidator accepts any valid JSON object with a 'name' key
type jsonHasNameKeyValidator struct{}

func (v jsonHasNameKeyValidator) Description(_ context.Context) string {
	return "field expected to be a valid JSON object with a 'name' key to some value"
}

func (v jsonHasNameKeyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonHasNameKeyValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	result := make(map[string]any)
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected a valid JSON object with 'name' key mapped to some value.\n'%s' was set to `%s`", request.Path, value))
		return
	}
	if _, ok := result["name"]; !ok {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected a valid JSON object with 'name' key mapped to some value.\n'%s' was set to `%s`", request.Path, value))
		return
	}
}

func JsonHasNameKeyValidator() validator.String {
	return jsonHasNameKeyValidator{}
}

var _ validator.Set = tagFormatValidator{}

// tagFormatValidator validates that list contains items with tag format.
type tagFormatValidator struct {
	max int
}

func (v tagFormatValidator) Description(_ context.Context) string {
	return "list must contain elements with 'key:value' format"
}

func (v tagFormatValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v tagFormatValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elems := req.ConfigValue.Elements()
	for _, elem := range elems {
		if elem.IsNull() || elem.IsUnknown() || hasTagFormat(unquote(elem.String())) {
			continue
		}

		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

func TagFormatValidator() validator.Set {
	return tagFormatValidator{}
}

// managementRuleTagValidator validates that tag_key and tag_operation are only set when property is 'tag'
type managementRuleTagValidator struct{}

func (v managementRuleTagValidator) Description(_ context.Context) string {
	return "ensures that tag_key and tag_operation are only set when property is 'tag'"
}

func (v managementRuleTagValidator) MarkdownDescription(_ context.Context) string {
	return "ensures that `tag_key` and `tag_operation` are only set when `property` is `'tag'`"
}

func (v managementRuleTagValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var rules []ManagementRuleModel
	diags := req.ConfigValue.ElementsAs(ctx, &rules, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, rule := range rules {
		sourceProperty := rule.SourceProperty.ValueString()
		hasSourceTagKey := !rule.SourceTagKey.IsNull() && !rule.SourceTagKey.IsUnknown()
		hasSourceTagOp := !rule.SourceTagOperation.IsNull() && !rule.SourceTagOperation.IsUnknown()

		if sourceProperty != "tag" && (hasSourceTagKey || hasSourceTagOp) {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i),
				"Invalid Management Rule Configuration",
				fmt.Sprintf("source_tag_key and source_tag_operation can only be set when source_property is 'tag', but source_property is '%s'", sourceProperty),
			)
		}

		if sourceProperty == "tag" && !hasSourceTagKey {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i),
				"Invalid Management Rule Configuration",
				"source_tag_key is required when source_property is 'tag'",
			)
		}

		targetProperty := rule.TargetProperty.ValueString()
		hasTargetTagKey := !rule.TargetTagKey.IsNull() && !rule.TargetTagKey.IsUnknown()
		hasTargetTagOp := !rule.TargetTagOperation.IsNull() && !rule.TargetTagOperation.IsUnknown()

		if targetProperty != "tag" && (hasTargetTagKey || hasTargetTagOp) {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i),
				"Invalid Management Rule Configuration",
				fmt.Sprintf("target_tag_key and target_tag_operation can only be set when target_property is 'tag', but target_property is '%s'", targetProperty),
			)
		}

		if targetProperty == "tag" && !hasTargetTagKey {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i),
				"Invalid Management Rule Configuration",
				"target_tag_key is required when target_property is 'tag'",
			)
		}
	}
}

func ManagementRuleTagValidator() validator.List {
	return managementRuleTagValidator{}
}
