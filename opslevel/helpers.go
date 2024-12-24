package opslevel

import (
	"context"
	"fmt"
	"github.com/opslevel/opslevel-go/v2024"

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

func GetTeamID(d *diag.Diagnostics, identifier string) *opslevel.ID {
	if opslevel.IsID(identifier) {
		return opslevel.NewID(identifier)
	}
	team, ok := opslevel.Cache.TryGetTeam(identifier)
	if !ok {
		d.AddError("opslevel error", fmt.Sprintf("failed to find team '%s'", identifier))
		return opslevel.NewID()
	}
	return &team.Id
}
