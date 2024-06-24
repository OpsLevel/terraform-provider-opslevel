package opslevel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/opslevel/opslevel-go/v2024"
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

	value := request.ConfigValue.ValueString()
	if !opslevel.IsID(value) {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected a valid id. id should start with Z2lkOi8v. '%s' was set to `%s`", request.Path, value))
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
		if hasTagFormat(unquote(elem.String())) {
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
