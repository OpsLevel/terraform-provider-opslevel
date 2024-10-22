package opslevel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type TerraformSource interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}

func read[T any](ctx context.Context, d *diag.Diagnostics, state TerraformSource) T {
	var data T
	d.Append(state.Get(ctx, &data)...)
	return data
}
