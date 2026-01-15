package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opslevel/opslevel-go/v2026"
)

type TerraformSource interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}

func read[T any](ctx context.Context, d *diag.Diagnostics, state TerraformSource) T {
	var data T
	d.Append(state.Get(ctx, &data)...)
	return data
}

func GetTeamID(d *diag.Diagnostics, client *opslevel.Client, identifier string) *opslevel.Nullable[opslevel.ID] {
	if opslevel.IsID(identifier) {
		return opslevel.RefOf(*opslevel.NewID(identifier))
	}
	team, err := client.GetTeamWithAlias(identifier)
	if err != nil {
		d.AddError("opslevel error", fmt.Sprintf("failed to find team with alias '%s': %s", identifier, err))
		return opslevel.RefOf(*opslevel.NewID())
	}
	return opslevel.RefOf(team.Id)
}

// Because the opslevel.RefOf changed to be a Nullable[T] we need a helper in here for backwards compatibility for things needed plain old *T
func refOf[T any](value T) *T {
	return &value
}

func nullable[T comparable](s *T) *opslevel.Nullable[T] {
	if s == nil {
		return &opslevel.Nullable[T]{
			SetNull: true,
		}
	}
	return opslevel.RefOf[T](*s)
}

func nullableID(s *string) *opslevel.Nullable[opslevel.ID] {
	if s == nil {
		return &opslevel.Nullable[opslevel.ID]{
			SetNull: true,
		}
	}
	if *s == "" {
		return nil
	}
	return opslevel.RefOf(opslevel.ID(*s))
}

func asEnum[T ~string](s *string) *T {
	if s == nil {
		return nil
	}
	value := T(*s)
	return &value
}
