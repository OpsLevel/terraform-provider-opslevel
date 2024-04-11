package opslevel

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2024"
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
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'.",
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(false),
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
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
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

type PredicateModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

func NewPredicateModel(predicate opslevel.Predicate) *PredicateModel {
	return &PredicateModel{
		Type:  RequiredStringValue(string(predicate.Type)),
		Value: OptionalStringValue(predicate.Value),
	}
}

func (s PredicateModel) ToCreateInput() *opslevel.PredicateInput {
	return &opslevel.PredicateInput{
		Type:  opslevel.PredicateTypeEnum(s.Type.ValueString()),
		Value: opslevel.RefOf(s.Value.ValueString()),
	}
}

func (s PredicateModel) ToUpdateInput() *opslevel.PredicateUpdateInput {
	return &opslevel.PredicateUpdateInput{
		Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(s.Type.ValueString())),
		Value: opslevel.RefOf(s.Value.ValueString()),
	}
}

func PredicateSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "A condition that should be satisfied.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "A condition that should be satisfied.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllPredicateTypeEnum...)},
			},
			"value": schema.StringAttribute{
				Description: "The condition value used by the predicate.",
				Optional:    true,
			},
		},
	}
}
