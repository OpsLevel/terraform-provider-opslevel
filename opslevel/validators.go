package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/opslevel/opslevel-go/v2024"
)

// Non Empty String Validator
type idStringValidator struct{}

// Description describes the validation in plain text formatting.
func (v idStringValidator) Description(_ context.Context) string {
	return "value expected to be a string that isn't an empty string"
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
		response.Diagnostics.AddError("expected a valid id", fmt.Sprintf("a valid id should start with Z2lkOi8v. %s was set to `%s`", request.Path, value))
	}
}

func IdStringValidator() validator.String {
	return idStringValidator{}
}
