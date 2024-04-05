package opslevel

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CheckBaseModel struct {
	Category    types.String `tfsdk:"category"`
	Description types.String `tfsdk:"description"`
	EnableOn    types.String `tfsdk:"enable_on"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Filter      types.String `tfsdk:"filter"`
	Level       types.String `tfsdk:"level"`
	Name        types.String `tfsdk:"name"`
	Notes       types.String `tfsdk:"notes"`
	Owner       types.String `tfsdk:"owner"`
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
		"enable_on": schema.StringAttribute{
			Description: `The date when the check will be automatically enabled.
 If you use this field you should add both 'enabled' and 'enable_on' to the lifecycle ignore_changes settings.
 See example in opslevel_check_manual for proper configuration.
 `,
			Optional: true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'.",
			Optional:    true,
		},
		"filter": schema.StringAttribute{
			Description: "The id of the filter of the check.",
			Optional:    true,
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
	}
	for key, value := range attrs {
		output[key] = value
	}
	return output
}
