package opslevel

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2024"
	"time"
)

type CheckBaseModel struct {
	Category    types.String      `tfsdk:"category"`
	Description types.String      `tfsdk:"description"`
	Enabled     types.Bool        `tfsdk:"enabled"`
	EnableOn    timetypes.RFC3339 `tfsdk:"enable_on"`
	Filter      types.String      `tfsdk:"filter"`
	Id          types.String      `tfsdk:"id"`
	Level       types.String      `tfsdk:"level"`
	Name        types.String      `tfsdk:"name"`
	Notes       types.String      `tfsdk:"notes"`
	Owner       types.String      `tfsdk:"owner"`
	LastUpdated timetypes.RFC3339 `tfsdk:"last_updated"`
}

func CheckBaseAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	output := map[string]schema.Attribute{
		"category": schema.StringAttribute{
			Description: "The id of the category the check belongs to.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description the check.",
			Optional:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'.",
			Optional:    true,
		},
		"enable_on": schema.StringAttribute{
			Description: `The date when the check will be automatically enabled.
 If you use this field you should add both 'enabled' and 'enable_on' to the lifecycle ignore_changes settings.
 See example in opslevel_check_manual for proper configuration.
 `,
			Optional: true,
		},
		"filter": schema.StringAttribute{
			Description: "The id of the filter of the check.",
			Optional:    true,
		},
		"id": schema.StringAttribute{
			Description: "The id of the check.",
			Computed:    true,
		},
		"level": schema.StringAttribute{
			Description: "The id of the level the check belongs to.",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "The display name of the check.",
			Required:    true,
		},
		"notes": schema.StringAttribute{
			Description: "Additional information to display to the service owner about the check.",
			Optional:    true,
		},
		"owner": schema.StringAttribute{
			Description: "The id of the team that owns the check.",
			Optional:    true,
		},
		"last_updated": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	}
	for key, value := range attrs {
		output[key] = value
	}
	return output
}

func NewCheckCreateInputFrom[T any](model CheckBaseModel) (*T, diag.Diagnostics) {
	enabledOn, diags := asISO8601(model.EnableOn)
	input := opslevel.CheckCreateInput{
		Category: *opslevel.NewID(model.Category.ValueString()),
		Enabled:  model.Enabled.ValueBoolPointer(),
		EnableOn: enabledOn,
		Filter:   opslevel.NewID(model.Filter.ValueString()),
		Level:    *opslevel.NewID(model.Level.ValueString()),
		Name:     model.Name.ValueString(),
		Notes:    model.Notes.ValueStringPointer(),
		Owner:    opslevel.NewID(model.Owner.ValueString()),
	}
	return opslevel.NewCheckCreateInputTypeOf[T](input), diags
}

func NewCheckUpdateInputFrom[T any](model CheckBaseModel) (*T, diag.Diagnostics) {
	enabledOn, diags := asISO8601(model.EnableOn)
	input := opslevel.CheckUpdateInput{
		Category: *opslevel.NewID(model.Category.ValueString()),
		Enabled:  model.Enabled.ValueBoolPointer(),
		EnableOn: enabledOn,
		Filter:   opslevel.NewID(model.Filter.ValueString()),
		Level:    *opslevel.NewID(model.Level.ValueString()),
		Id:       *opslevel.NewID(model.Id.ValueString()),
		Name:     model.Name.ValueString(),
		Notes:    model.Notes.ValueStringPointer(),
		Owner:    opslevel.NewID(model.Owner.ValueString()),
	}
	return opslevel.NewCheckUpdateInputTypeOf[T](input), diags
}

func ApplyCheckBaseModel(check opslevel.Check, model *CheckBaseModel) {
	model.Category = types.StringValue(string(check.Category.Id))
	model.Enabled = types.BoolValue(check.Enabled)
	model.EnableOn = timetypes.NewRFC3339TimeValue(check.EnableOn.Time)
	model.Filter = types.StringValue(string(check.Filter.Id))
	model.Id = types.StringValue(string(check.Id))
	model.Level = types.StringValue(string(check.Level.Id))
	model.Name = types.StringValue(check.Name)
	model.Notes = types.StringValue(check.Notes)
	model.Owner = types.StringValue(string(check.Owner.Team.Id))
	model.LastUpdated = timetypes.NewRFC3339TimeValue(time.Now())
}
