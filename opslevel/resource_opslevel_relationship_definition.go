package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

var (
	_ resource.ResourceWithConfigure   = &RelationshipDefinitionResource{}
	_ resource.ResourceWithImportState = &RelationshipDefinitionResource{}
)

var (
	TEAM_BUILTIN_PROPERTIES      = []string{"name", "alias", "contact", "tag"}
	USER_BUILTIN_PROPERTIES      = []string{"name", "contact", "tag"}
	COMPONENT_BUILTIN_PROPERTIES = []string{"name", "alias", "tag"}
)

func NewRelationshipDefinitionResource() resource.Resource {
	return &RelationshipDefinitionResource{}
}

// RelationshipDefinitionResource defines the resource implementation.
type RelationshipDefinitionResource struct {
	CommonResourceClient
}

// RelationshipDefinitionResourceModel describes the Relationship Definition managed resource.
type RelationshipDefinitionResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Alias             types.String `tfsdk:"alias"`
	Description       types.String `tfsdk:"description"`
	ComponentType     types.String `tfsdk:"component_type"`
	AllowedCategories types.List   `tfsdk:"allowed_categories"`
	AllowedTypes      types.List   `tfsdk:"allowed_types"`
	ManagementRules   types.List   `tfsdk:"management_rules"`
}

type ManagementRuleModel struct {
	Operator           types.String `tfsdk:"operator"`
	SourceProperty     types.String `tfsdk:"source_property"`
	SourceTagKey       types.String `tfsdk:"source_tag_key"`
	SourceTagOperation types.String `tfsdk:"source_tag_operation"`
	TargetCategory     types.String `tfsdk:"target_category"`
	TargetProperty     types.String `tfsdk:"target_property"`
	TargetTagKey       types.String `tfsdk:"target_tag_key"`
	TargetTagOperation types.String `tfsdk:"target_tag_operation"`
	TargetType         types.String `tfsdk:"target_type"`
}

func ManagementRuleModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"operator":             types.StringType,
		"source_property":      types.StringType,
		"source_tag_key":       types.StringType,
		"source_tag_operation": types.StringType,
		"target_category":      types.StringType,
		"target_property":      types.StringType,
		"target_tag_key":       types.StringType,
		"target_tag_operation": types.StringType,
		"target_type":          types.StringType,
	}
}

func NewRelationshipDefinitionResourceModel(definition opslevel.RelationshipDefinitionType, givenModel RelationshipDefinitionResourceModel) RelationshipDefinitionResourceModel {
	model := RelationshipDefinitionResourceModel{
		Id:            ComputedStringValue(string(definition.Id)),
		Name:          RequiredStringValue(definition.Name),
		Alias:         RequiredStringValue(definition.Alias),
		Description:   StringValueFromResourceAndModelField(definition.Description, givenModel.Description),
		ComponentType: givenModel.ComponentType,
		AllowedCategories: types.ListValueMust(
			types.StringType,
			func() []attr.Value {
				values := make([]attr.Value, len(definition.Metadata.AllowedCategories))
				for i, v := range definition.Metadata.AllowedCategories {
					values[i] = types.StringValue(v)
				}
				return values
			}(),
		),
		AllowedTypes: types.ListValueMust(
			types.StringType,
			func() []attr.Value {
				values := make([]attr.Value, len(definition.Metadata.AllowedTypes))
				for i, v := range definition.Metadata.AllowedTypes {
					values[i] = types.StringValue(v)
				}
				return values
			}(),
		),
	}

	if len(definition.ManagementRules) > 0 {
		ruleValues := make([]attr.Value, len(definition.ManagementRules))
		for i, rule := range definition.ManagementRules {

			ruleValues[i] = newManagementRuleValue(rule)
		}

		model.ManagementRules = types.ListValueMust(
			types.ObjectType{AttrTypes: ManagementRuleModelAttrs()},
			ruleValues,
		)
	} else if !givenModel.ManagementRules.IsNull() {
		model.ManagementRules = givenModel.ManagementRules
	} else {
		model.ManagementRules = types.ListNull(types.ObjectType{AttrTypes: ManagementRuleModelAttrs()})
	}

	return model
}

func (r *RelationshipDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_definition"
}

func (r *RelationshipDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `A relationship definition in OpsLevel defines how different resources can be related to each other. It specifies:
- Which component type the relationship is on
- What types of resources can be related (allowed_types)

`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the relationship definition.",
				Required:    true,
			},
			"alias": schema.StringAttribute{
				Description: "The unique identifier of the relationship.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the relationship definition.",
				Optional:    true,
			},
			"component_type": schema.StringAttribute{
				Description: "The component type that the relationship belongs to. Must be a valid component type alias from your OpsLevel account or 'team'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allowed_categories": schema.ListAttribute{
				Description: "The categories of resources that can be selected for this relationship definition. Can include any component category alias on your account.",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"allowed_types": schema.ListAttribute{
				Description: "The types of resources that can be selected for this relationship definition. Can include any component type alias on your account or 'team' or 'user'.",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"management_rules": schema.ListNestedAttribute{
				Description: "Rules that automatically manage relationships based on property matching conditions.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operator": schema.StringAttribute{
							Description: "The condition operator for this rule. Either EQUALS or ARRAY_CONTAINS.",
							Required:    true,
						},
						"source_property": schema.StringAttribute{
							Description: "The property on the source component to evaluate.",
							Required:    true,
						},
						"source_tag_key": schema.StringAttribute{
							Description: "When source_property is 'tag', this specifies the tag key to match.",
							Optional:    true,
						},
						"source_tag_operation": schema.StringAttribute{
							Description: "When source_property is 'tag', this specifies the matching operation. Either 'equals' or 'starts_with'. Defaults to 'equals'.",
							Optional:    true,
						},
						"target_category": schema.StringAttribute{
							Description: "The category of the target resource. Either target_category or target_type must be specified, but not both.",
							Optional:    true,
						},
						"target_property": schema.StringAttribute{
							Description: "The property on the target resource to match against.",
							Required:    true,
						},
						"target_tag_key": schema.StringAttribute{
							Description: "When target_property is 'tag', this specifies the tag key to match.",
							Optional:    true,
						},
						"target_tag_operation": schema.StringAttribute{
							Description: "When target_property is 'tag', this specifies the matching operation. Either 'equals' or 'starts_with'. Defaults to 'equals'.",
							Optional:    true,
						},
						"target_type": schema.StringAttribute{
							Description: "The type of the target resource. Either target_category or target_type must be specified, but not both.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *RelationshipDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	allowedCategories := make([]string, 0)
	if err := planModel.AllowedCategories.ElementsAs(ctx, &allowedCategories, false); err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_categories: %s", err))
		return
	}

	allowedTypes := make([]string, 0)
	if err := planModel.AllowedTypes.ElementsAs(ctx, &allowedTypes, false); err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_types: %s", err))
		return
	}

	componentTypeAlias := r.GetComponentTypeAlias(planModel.ComponentType.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	managementRules := parseManagementRules(ctx, planModel.ManagementRules, componentTypeAlias, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.RelationshipDefinitionInput{
		Name:            planModel.Name.ValueStringPointer(),
		Alias:           planModel.Alias.ValueStringPointer(),
		Description:     nullable(planModel.Description.ValueStringPointer()),
		ComponentType:   opslevel.NewIdentifier(planModel.ComponentType.ValueString()),
		ManagementRules: &managementRules,
		Metadata: &opslevel.RelationshipDefinitionMetadataInput{
			AllowedTypes:      allowedTypes,
			AllowedCategories: allowedCategories,
		},
	}

	definition, err := r.client.CreateRelationshipDefinition(input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create relationship definition with name '%s', got error: %s", planModel.Name.ValueString(), err))
		return
	}

	stateModel := NewRelationshipDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("created a relationship definition resource with id '%s'", definition.Id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RelationshipDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	definition, err := r.client.GetRelationshipDefinition(id)
	if err != nil {
		if (definition == nil || definition.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read relationship definition with id '%s', got error: %s", id, err))
		return
	}

	verifiedStateModel := NewRelationshipDefinitionResourceModel(*definition, stateModel)
	tflog.Trace(ctx, fmt.Sprintf("read a relationship definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *RelationshipDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	allowedCategories := make([]string, 0)
	if err := planModel.AllowedCategories.ElementsAs(ctx, &allowedCategories, false); err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_categories: %s", err))
		return
	}

	allowedTypes := make([]string, 0)
	if err := planModel.AllowedTypes.ElementsAs(ctx, &allowedTypes, false); err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_types: %s", err))
		return
	}

	componentTypeAlias := r.GetComponentTypeAlias(planModel.ComponentType.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	managementRules := parseManagementRules(ctx, planModel.ManagementRules, componentTypeAlias, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	id := planModel.Id.ValueString()
	input := opslevel.RelationshipDefinitionInput{
		Name:            planModel.Name.ValueStringPointer(),
		Alias:           planModel.Alias.ValueStringPointer(),
		Description:     nullable(planModel.Description.ValueStringPointer()),
		ComponentType:   opslevel.NewIdentifier(planModel.ComponentType.ValueString()),
		ManagementRules: &managementRules,
		Metadata: &opslevel.RelationshipDefinitionMetadataInput{
			AllowedCategories: allowedCategories,
			AllowedTypes:      allowedTypes,
		},
	}

	definition, err := r.client.UpdateRelationshipDefinition(id, input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update relationship definition with id '%s', got error: %s", id, err))
		return
	}

	stateModel := NewRelationshipDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("updated a relationship definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RelationshipDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	_, err := r.client.DeleteRelationshipDefinition(id)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete relationship definition (%s), got error: %s", id, err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("deleted a relationship definition resource with id '%s'", id))
}

func (r *RelationshipDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RelationshipDefinitionResource) GetComponentTypeAlias(componentTypeValue string, diags *diag.Diagnostics) string {
	if componentTypeValue == "team" {
		return componentTypeValue
	}

	componentType, err := r.client.GetComponentType(componentTypeValue)
	if err != nil {
		diags.AddError("opslevel client error", fmt.Sprintf("unable to fetch component type with id '%s': %s", componentTypeValue, err))
		return ""
	}

	return componentType.Aliases[0]
}

func parseManagementRules(ctx context.Context, planRules types.List, componentTypeAlias string, diags *diag.Diagnostics) []opslevel.RelationshipDefinitionManagementRulesInput {
	if planRules.IsNull() || planRules.IsUnknown() {
		return nil
	}

	var planRulesModels []ManagementRuleModel
	if err := planRules.ElementsAs(ctx, &planRulesModels, false); err != nil {
		diags.AddError("config error", fmt.Sprintf("unable to parse management_rules: %s", err))
		return nil
	}

	managementRules := make([]opslevel.RelationshipDefinitionManagementRulesInput, len(planRulesModels))
	for i, rule := range planRulesModels {
		var targetTypeOrCategory string
		isType := false

		if !rule.TargetType.IsNull() && !rule.TargetType.IsUnknown() {
			targetTypeOrCategory = rule.TargetType.ValueString()
			isType = true
		} else if !rule.TargetCategory.IsNull() && !rule.TargetCategory.IsUnknown() {
			targetTypeOrCategory = rule.TargetCategory.ValueString()
			isType = false
		}

		sourcePropertyStr := buildPropertyString(
			rule.SourceProperty.ValueString(),
			rule.SourceTagKey.ValueString(),
			rule.SourceTagOperation.ValueString(),
		)

		targetPropertyStr := buildPropertyString(
			rule.TargetProperty.ValueString(),
			rule.TargetTagKey.ValueString(),
			rule.TargetTagOperation.ValueString(),
		)

		sourcePropertyBuiltin := isBuiltinProperty(componentTypeAlias, rule.SourceProperty.ValueString(), true)
		targetPropertyBuiltin := isBuiltinProperty(targetTypeOrCategory, rule.TargetProperty.ValueString(), isType)

		managementRules[i] = opslevel.RelationshipDefinitionManagementRulesInput{
			Operator:              opslevel.RelationshipOperatorEnum(rule.Operator.ValueString()),
			SourceProperty:        sourcePropertyStr,
			SourcePropertyBuiltin: sourcePropertyBuiltin,
			TargetProperty:        targetPropertyStr,
			TargetPropertyBuiltin: targetPropertyBuiltin,
		}

		if !rule.TargetCategory.IsNull() && !rule.TargetCategory.IsUnknown() {
			targetCategory := rule.TargetCategory.ValueString()
			managementRules[i].TargetCategory = nullable(&targetCategory)
		}

		if !rule.TargetType.IsNull() && !rule.TargetType.IsUnknown() {
			targetType := rule.TargetType.ValueString()
			managementRules[i].TargetType = nullable(&targetType)
		}
	}

	return managementRules
}

func newManagementRuleValue(rule opslevel.RelationshipDefinitionManagementRules) attr.Value {
	var targetCategory types.String
	if rule.TargetCategory != nil && !rule.TargetCategory.SetNull {
		targetCategory = types.StringValue(rule.TargetCategory.Value)
	} else {
		targetCategory = types.StringNull()
	}

	var targetType types.String
	if rule.TargetType != nil && !rule.TargetType.SetNull {
		targetType = types.StringValue(rule.TargetType.Value)
	} else {
		targetType = types.StringNull()
	}

	sourceProperty, sourceTagKey, sourceTagOp := parsePropertyString(rule.SourceProperty)
	targetProperty, targetTagKey, targetTagOp := parsePropertyString(rule.TargetProperty)

	return types.ObjectValueMust(
		ManagementRuleModelAttrs(),
		map[string]attr.Value{
			"operator":             types.StringValue(string(rule.Operator)),
			"source_property":      types.StringValue(sourceProperty),
			"source_tag_key":       OptionalStringValue(sourceTagKey),
			"source_tag_operation": OptionalStringValue(sourceTagOp),
			"target_category":      targetCategory,
			"target_property":      types.StringValue(targetProperty),
			"target_tag_key":       OptionalStringValue(targetTagKey),
			"target_tag_operation": OptionalStringValue(targetTagOp),
			"target_type":          targetType,
		},
	)
}

func isBuiltinProperty(targetTypeOrCategory string, propertyName string, isType bool) bool {
	var builtinProps []string

	if isType {
		if targetTypeOrCategory == "team" {
			builtinProps = TEAM_BUILTIN_PROPERTIES
		} else if targetTypeOrCategory == "user" {
			builtinProps = USER_BUILTIN_PROPERTIES
		} else {
			builtinProps = COMPONENT_BUILTIN_PROPERTIES
		}
	} else {
		if targetTypeOrCategory == "people" {
			builtinProps = TEAM_BUILTIN_PROPERTIES
		} else {
			builtinProps = COMPONENT_BUILTIN_PROPERTIES
		}
	}

	for _, prop := range builtinProps {
		if prop == propertyName {
			return true
		}
	}
	return false
}

func buildPropertyString(property, tagKey, tagOperation string) string {
	if property != "tag" {
		return property
	}

	operation := "eq"
	if tagOperation != "" {
		if tagOperation == "starts_with" {
			operation = "starts_with"
		}
	}

	return fmt.Sprintf("tag_key_%s:%s", operation, tagKey)
}
func parsePropertyString(propertyStr string) (property, tagKey, tagOperation string) {
	if !strings.HasPrefix(propertyStr, "tag_key_") {
		return propertyStr, "", ""
	}

	property = "tag"

	remainder := strings.TrimPrefix(propertyStr, "tag_key_")

	if strings.HasPrefix(remainder, "eq:") {
		tagOperation = "equals"
		tagKey = strings.TrimPrefix(remainder, "eq:")
	} else if strings.HasPrefix(remainder, "starts_with:") {
		tagOperation = "starts_with"
		tagKey = strings.TrimPrefix(remainder, "starts_with:")
	}

	return
}
