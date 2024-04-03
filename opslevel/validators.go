package opslevel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/opslevel/opslevel-go/v2024"
)

// OpsLevel ID String Validator
type idStringValidator struct{}

// Description describes the validation in plain text formatting.
func (v idStringValidator) Description(_ context.Context) string {
	return "id expected to be a string starting with 'Z2lkOi8v'"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v idStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
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

// OpsLevel ID String Validator
type jsonStringValidator struct{}

// Description describes the validation in plain text formatting.
func (v jsonStringValidator) Description(_ context.Context) string {
	return "field expected to be valid JSON"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v jsonStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
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

// OpsLevel ID String Validator
type jsonHasNameKeyValidator struct{}

// Description describes the validation in plain text formatting.
func (v jsonHasNameKeyValidator) Description(_ context.Context) string {
	return "field expected to be valid JSON with a 'name' key to some value"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v jsonHasNameKeyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v jsonHasNameKeyValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	result := make(map[string]any)
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected valid JSON. '%s' was set to `%s`", request.Path, value))
		return
	}
	if _, ok := result["name"]; !ok {
		response.Diagnostics.AddError("Config error", fmt.Sprintf("expected JSON with 'name' key mapped to some value.\n'%s' was set to `%s`", request.Path, value))
	}
}

func JsonHasNameKeyValidator() validator.String {
	return jsonHasNameKeyValidator{}
}
