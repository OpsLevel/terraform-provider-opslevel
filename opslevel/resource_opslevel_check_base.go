package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2025"
)

type CheckCodeBaseResourceModel struct {
	Category    types.String `tfsdk:"category"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	EnableOn    types.String `tfsdk:"enable_on"`
	Filter      types.String `tfsdk:"filter"`
	Id          types.String `tfsdk:"id"`
	Level       types.String `tfsdk:"level"`
	Name        types.String `tfsdk:"name"`
	Notes       types.String `tfsdk:"notes"`
	Owner       types.String `tfsdk:"owner"`
}

func NewCheckCodeBaseResourceModel(check opslevel.Check, givenModel CheckCodeBaseResourceModel) CheckCodeBaseResourceModel {
	var stateModel CheckCodeBaseResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
	stateModel.Description = ComputedStringValue(check.Description)
	if givenModel.Enabled.IsNull() {
		stateModel.Enabled = types.BoolValue(false)
	} else {
		stateModel.Enabled = OptionalBoolValue(&check.Enabled)
	}
	if givenModel.EnableOn.IsNull() {
		stateModel.EnableOn = types.StringNull()
	} else {
		// We pass through the plan value because of time formatting issue to ensure the state gets the exact value the customer specified
		stateModel.EnableOn = givenModel.EnableOn
	}
	stateModel.Filter = OptionalStringValue(string(check.Filter.Id))
	stateModel.Id = ComputedStringValue(string(check.Id))
	stateModel.Level = RequiredStringValue(string(check.Level.Id))
	stateModel.Name = RequiredStringValue(check.Name)
	stateModel.Notes = OptionalStringValue(check.Notes)
	stateModel.Owner = OptionalStringValue(string(check.Owner.Team.Id))

	return stateModel
}

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

var predicateType = map[string]attr.Type{
	"type":  types.StringType,
	"value": types.StringType,
}

func (p PredicateModel) Validate() error {
	// Skip validation for now, when input variables or for loops are used
	if p.Type.IsNull() || p.Type.IsUnknown() || p.Value.IsUnknown() {
		return nil
	}
	predicate := opslevel.Predicate{
		Type:  opslevel.PredicateTypeEnum(p.Type.ValueString()),
		Value: p.Value.ValueString(),
	}
	return predicate.Validate()
}

func (p PredicateModel) ToCreateInput() *opslevel.PredicateInput {
	return &opslevel.PredicateInput{
		Type:  opslevel.PredicateTypeEnum(p.Type.ValueString()),
		Value: opslevel.RefOf(p.Value.ValueString()),
	}
}

func (s PredicateModel) ToUpdateInput() *opslevel.PredicateUpdateInput {
	return &opslevel.PredicateUpdateInput{
		Type:  asEnum[opslevel.PredicateTypeEnum](s.Type.ValueStringPointer()),
		Value: opslevel.RefOf(s.Value.ValueString()),
	}
}

func PredicateSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: fmt.Sprintf(
			"A condition that should be satisfied. One of `%s`",
			strings.Join(opslevel.AllPredicateTypeEnum, "`, `"),
		),
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "A condition that should be satisfied.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllPredicateTypeEnum...)},
			},
			"value": schema.StringAttribute{
				Description: "The condition value used by the predicate.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
		},
	}
}

// pre-v1.x base schema fields for checks used by StateUpgraders
func getCheckBaseSchemaV0(extras map[string]schema.Attribute) map[string]schema.Attribute {
	output := map[string]schema.Attribute{
		"last_updated": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"name": schema.StringAttribute{
			Description: "The display name of the check.",
			Required:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: `Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'.`,
			Required:    true,
		},
		"enable_on": schema.StringAttribute{
			Description: `The date when the check will be automatically enabled.
If you use this field you should add both 'enabled' and 'enable_on' to the lifecycle ignore_changes settings.
See example in opslevel_check_manual for proper configuration.
`,
			Optional: true,
		},
		"category": schema.StringAttribute{
			Description: "The id of the category the check belongs to.",
			Required:    true,
		},
		"level": schema.StringAttribute{
			Description: "The id of the level the check belongs to.",
			Required:    true,
		},
		"owner": schema.StringAttribute{
			Description: "The id of the team that owns the check.",
			Optional:    true,
		},
		"filter": schema.StringAttribute{
			Description: "The id of the filter of the check.",
			Optional:    true,
		},
		"notes": schema.StringAttribute{
			Description: "Additional information about the check.",
			Optional:    true,
		},
	}
	for k, v := range extras {
		output[k] = v
	}
	return output
}

var predicateSchemaV0 = schema.NestedBlockObject{
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
