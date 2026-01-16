package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
)

var (
	_ resource.ResourceWithConfigure   = &RelationshipDefinitionResource{}
	_ resource.ResourceWithImportState = &RelationshipDefinitionResource{}
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

func NewRelationshipDefinitionResourceModel(definition opslevel.RelationshipDefinitionType, givenModel RelationshipDefinitionResourceModel) RelationshipDefinitionResourceModel {
	model := RelationshipDefinitionResourceModel{
		Id:                ComputedStringValue(string(definition.Id)),
		Name:              RequiredStringValue(definition.Name),
		Alias:             RequiredStringValue(definition.Alias),
		Description:       StringValueFromResourceAndModelField(definition.Description, givenModel.Description),
		ComponentType:     givenModel.ComponentType,
		AllowedCategories: StringListValueFromResourceAndModelField(definition.Metadata.AllowedCategories, givenModel.AllowedCategories),
		AllowedTypes:      StringListValueFromResourceAndModelField(definition.Metadata.AllowedTypes, givenModel.AllowedTypes),
	}

	model.ManagementRules = ManagementRuleListValueFromResourceAndModel(definition.ManagementRules, givenModel.ManagementRules)

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
				Validators: []validator.List{
					listvalidator.AtLeastOneOf(path.MatchRoot("allowed_types")),
				},
				PlanModifiers: []planmodifier.List{
					PreventRemovalPlanModifier("allowed_categories"),
				},
			},
			"allowed_types": schema.ListAttribute{
				Description: "The types of resources that can be selected for this relationship definition. Can include any component type alias on your account or 'team' or 'user'.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.AtLeastOneOf(path.MatchRoot("allowed_categories")),
				},
				PlanModifiers: []planmodifier.List{
					PreventRemovalPlanModifier("allowed_types"),
				},
			},
			"management_rules": ManagementRulesResourceAttribute(),
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

	managementRules := ParseManagementRules(ctx, planModel.ManagementRules, componentTypeAlias, &resp.Diagnostics)
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

	managementRules := ParseManagementRules(ctx, planModel.ManagementRules, componentTypeAlias, &resp.Diagnostics)
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
