package opslevel

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var checkBaseAttributes = map[string]schema.Attribute{
	"category": schema.StringAttribute{
		Description: "The id of the category the check belongs to.",
		Required:    true,
		Validators:  []validator.String{IdStringValidator()},
	},
	"description": schema.StringAttribute{
		Description: "The description the check.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'.",
		Optional:    true,
		Computed:    true,
		Validators:  []validator.Bool{boolvalidator.ConflictsWith(path.MatchRoot("enable_on"))},
	},
	"enable_on": schema.StringAttribute{
		Description: `The date when the check will be automatically enabled.
 If you use this field you should add both 'enabled' and 'enable_on' to the lifecycle ignore_changes settings.
 See example in opslevel_check_manual for proper configuration.
 `,
		Optional:   true,
		Validators: []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("enabled"))},
	},
	"filter": schema.StringAttribute{
		Description: "The id of the filter of the check.",
		Optional:    true,
		Validators:  []validator.String{IdStringValidator()},
	},
	"id": schema.StringAttribute{
		Description: "The id of the check.",
		Computed:    true,
	},
	"level": schema.StringAttribute{
		Description: "The id of the level the check belongs to.",
		Required:    true,
		Validators:  []validator.String{IdStringValidator()},
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
		Validators:  []validator.String{IdStringValidator()},
	},
	"last_updated": schema.StringAttribute{
		Optional: true,
		Computed: true,
	},
}

func CheckBaseAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range checkBaseAttributes {
		attrs[key] = value
	}
	return attrs
}
