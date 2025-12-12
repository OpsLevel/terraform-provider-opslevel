package opslevel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2025"
	"golang.org/x/net/context"
)

var (
	_ resource.Resource                = &ComponentTypeResource{}
	_ resource.ResourceWithImportState = &ComponentTypeResource{}
)

type PropertyModel struct {
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	AllowedInConfigFiles types.Bool   `tfsdk:"allowed_in_config_files"`
	DisplayStatus        types.String `tfsdk:"display_status"`
	LockedStatus         types.String `tfsdk:"locked_status"`
	Schema               types.String `tfsdk:"schema"`
}

type ComponentTypeIconModel struct {
	Color types.String `tfsdk:"color"`
	Name  types.String `tfsdk:"name"`
}

type RelationshipModel struct {
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	AllowedCategories types.List   `tfsdk:"allowed_categories"`
	AllowedTypes      types.List   `tfsdk:"allowed_types"`
}

type ComponentTypeModel struct {
	Id                types.String                 `tfsdk:"id"`
	Name              types.String                 `tfsdk:"name"`
	Alias             types.String                 `tfsdk:"alias"`
	Description       types.String                 `tfsdk:"description"`
	Icon              *ComponentTypeIconModel      `tfsdk:"icon"`
	OwnerRelationship *OwnerRelationshipModel      `tfsdk:"owner_relationship"`
	Properties        map[string]PropertyModel     `tfsdk:"properties"`
	Relationships     map[string]RelationshipModel `tfsdk:"relationships"`
}

type OwnerRelationshipModel struct {
	ManagementRules types.List `tfsdk:"management_rules"`
}

type ComponentTypeResource struct {
	CommonResourceClient
}

func NewComponentTypeResource() resource.Resource {
	return &ComponentTypeResource{}
}

func (s ComponentTypeResource) NewModel(res *opslevel.ComponentType, stateModel ComponentTypeModel) (ComponentTypeModel, error) {
	stateModel.Id = types.StringValue(string(res.Id))
	stateModel.Name = types.StringValue(res.Name)
	stateModel.Alias = types.StringValue(res.Aliases[0])
	stateModel.Description = types.StringValue(res.Description)
	stateModel.Icon = &ComponentTypeIconModel{
		Color: types.StringValue(res.Icon.Color),
		Name:  types.StringValue(string(res.Icon.Name)),
	}
	conn, err := res.GetProperties(s.client, nil)
	if err != nil {
		return stateModel, err
	}
	stateModel.Properties = map[string]PropertyModel{}
	for _, prop := range conn.Nodes {
		stateModel.Properties[prop.Aliases[0]] = PropertyModel{
			Name:                 types.StringValue(prop.Name),
			Description:          ComputedStringValue(prop.Description),
			AllowedInConfigFiles: types.BoolValue(prop.AllowedInConfigFiles),
			DisplayStatus:        types.StringValue(string(prop.PropertyDisplayStatus)),
			LockedStatus:         types.StringValue(string(prop.LockedStatus)),
			Schema:               types.StringValue(prop.Schema.AsString()),
		}
	}
	return stateModel, nil
}

func NewPropertiesInput(model ComponentTypeModel) (*[]opslevel.ComponentTypePropertyDefinitionInput, error) {
	var properties []opslevel.ComponentTypePropertyDefinitionInput
	for alias, prop := range model.Properties {
		propertySchema, err := opslevel.NewJSONSchema(prop.Schema.ValueString())
		if err != nil {
			return nil, err
		}
		properties = append(properties, opslevel.ComponentTypePropertyDefinitionInput{
			Name:                  prop.Name.ValueString(),
			Alias:                 alias,
			Description:           prop.Description.ValueString(),
			AllowedInConfigFiles:  prop.AllowedInConfigFiles.ValueBool(),
			PropertyDisplayStatus: *asEnum[opslevel.PropertyDisplayStatusEnum](prop.DisplayStatus.ValueStringPointer()),
			LockedStatus:          asEnum[opslevel.PropertyLockedStatusEnum](prop.LockedStatus.ValueStringPointer()),
			Schema:                *propertySchema,
		})
	}
	return &properties, nil
}

func (s ComponentTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component_type"
}

func (s ComponentTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Component Type Definition Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique name of the component type.",
				Required:    true,
			},
			"alias": schema.StringAttribute{
				Description: "The unique alias of the component type.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the component type.",
				Optional:    true,
			},
			"icon": schema.SingleNestedAttribute{
				Description: "The icon associated with the component type",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"color": schema.StringAttribute{
						Description: "The color, represented as a hexcode, for the icon.",
						Required:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the icon in Phosphor icons for Vue, e.g. `PhBird`. See https://phosphoricons.com/ for a full list.",
						Required:    true,
						Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllComponentTypeIconEnum...)},
					},
				},
			},
			"owner_relationship": schema.SingleNestedAttribute{
				Description: "The owner relationship configuration for this component type.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"management_rules": schema.ListNestedAttribute{
						Description: "Rules that automatically determine ownership based on property matching conditions.",
						Optional:    true,
						Validators: []validator.List{
							ManagementRuleTagValidator(),
						},
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
									Description: "When source_property is 'tag', this specifies the tag key to match. Required if source_property is 'tag', must not be set otherwise.",
									Optional:    true,
								},
								"source_tag_operation": schema.StringAttribute{
									Description: "When source_property is 'tag', this specifies the matching operation. Either 'equals' or 'starts_with'. Defaults to 'equals'. Required if source_property is 'tag', must not be set otherwise",
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
									Description: "When target_property is 'tag', this specifies the tag key to match. Required if target_property is 'tag', must not be set otherwise.",
									Optional:    true,
								},
								"target_tag_operation": schema.StringAttribute{
									Description: "When target_property is 'tag', this specifies the matching operation. Either 'equals' or 'starts_with'. Defaults to 'equals'. Required if target_property is 'tag', must not be set otherwise.",
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
			},
			"properties": schema.MapNestedAttribute{
				Description: "The properties of this component type.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the property definition.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the property definition.",
							Optional:    true,
						},
						"allowed_in_config_files": schema.BoolAttribute{
							Description: "Whether or not the property is allowed to be set in opslevel.yml config files.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"display_status": schema.StringAttribute{
							Description: "The display status of the property.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(string(opslevel.PropertyDisplayStatusEnumVisible)),
							Validators: []validator.String{
								stringvalidator.OneOf(opslevel.AllPropertyDisplayStatusEnum...),
							},
						},
						"locked_status": schema.StringAttribute{
							Description: "Restricts what sources are able to assign values to this property.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(string(opslevel.PropertyLockedStatusEnumUILocked)),
							Validators: []validator.String{
								stringvalidator.OneOf(opslevel.AllPropertyLockedStatusEnum...),
							},
						},
						"schema": schema.StringAttribute{
							Description: "The schema of the property.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("{}"),
							Validators: []validator.String{
								JsonObjectValidator(),
							},
						},
					},
				},
			},
			"relationships": schema.MapNestedAttribute{
				Description: "The relationships that can be defined for this component type.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The display name of the relationship definition.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the relationship definition.",
							Optional:    true,
						},
						"allowed_categories": schema.ListAttribute{
							MarkdownDescription: "The categories of resources that can be selected for this relationship definition. Can include any component category alias on your account.",
							Optional:            true,
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("allowed_types")),
							},
							PlanModifiers: []planmodifier.List{
								listplanmodifier.RequiresReplace(),
							},
						},
						"allowed_types": schema.ListAttribute{
							Description: "The types of resources that can be selected for this relationship definition. Can include any component type alias on your account or 'team'.",
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.List{
								listvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("allowed_categories")),
							},
							PlanModifiers: []planmodifier.List{
								listplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

func (s ComponentTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[ComponentTypeModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	properties, err := NewPropertiesInput(planModel)
	if err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to use schema, got error: %s", err))
		return
	}

	// Create the component type first
	input := opslevel.ComponentTypeInput{
		Name:        nullable(planModel.Name.ValueStringPointer()),
		Alias:       nullable(planModel.Alias.ValueStringPointer()),
		Description: nullable(planModel.Description.ValueStringPointer()),
		Properties:  properties,
	}
	if !planModel.Icon.Color.IsNull() && !planModel.Icon.Name.IsNull() {
		input.Icon = &opslevel.ComponentTypeIconInput{
			Color: planModel.Icon.Color.ValueString(),
			Name:  opslevel.ComponentTypeIconEnum(planModel.Icon.Name.ValueString()),
		}
	}

	res, err := s.client.CreateComponentType(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create resource, got error: %s", err))
		return
	}

	// Create relationship definitions if any are specified
	if len(planModel.Relationships) > 0 {
		for alias, rel := range planModel.Relationships {
			allowedCategories := make([]string, 0)
			if err := rel.AllowedCategories.ElementsAs(ctx, &allowedCategories, false); err != nil {
				resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_categories for relationship '%s': %s", alias, err))
				continue
			}

			allowedTypes := make([]string, 0)
			if err := rel.AllowedTypes.ElementsAs(ctx, &allowedTypes, false); err != nil {
				resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_types for relationship '%s': %s", alias, err))
				continue
			}

			relInput := opslevel.RelationshipDefinitionInput{
				Name:          rel.Name.ValueStringPointer(),
				Alias:         &alias,
				Description:   nullable(rel.Description.ValueStringPointer()),
				ComponentType: opslevel.NewIdentifier(string(res.Id)),
				Metadata: &opslevel.RelationshipDefinitionMetadataInput{
					AllowedCategories: allowedCategories,
					AllowedTypes:      allowedTypes,
				},
			}

			_, err := s.client.CreateRelationshipDefinition(relInput)
			if err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create relationship definition '%s', got error: %s", alias, err))
				continue
			}
		}
	}

	finalModel, err := s.NewModel(res, planModel)
	if err != nil {
		resp.Diagnostics.AddError("error", fmt.Sprintf("unable to build resource, got error: %s", err))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (s ComponentTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[ComponentTypeModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	res, err := s.client.GetComponentType(id)
	if err != nil {
		if (res == nil || res.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to get resource with id '%s', got error: %s", id, err))
		return
	}

	// Get relationship definitions
	rels, err := s.client.ListRelationshipDefinitions(&opslevel.PayloadVariables{
		"componentType": opslevel.NewIdentifier(id),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to list relationship definitions, got error: %s", err))
		return
	}

	// Add relationships to the state model
	stateModel.Relationships = make(map[string]RelationshipModel)
	for _, rel := range rels.Nodes {
		allowedCategories := make([]attr.Value, len(rel.Metadata.AllowedCategories))
		for i, t := range rel.Metadata.AllowedCategories {
			allowedCategories[i] = types.StringValue(t)
		}

		allowedTypes := make([]attr.Value, len(rel.Metadata.AllowedTypes))
		for i, t := range rel.Metadata.AllowedTypes {
			allowedTypes[i] = types.StringValue(t)
		}

		stateModel.Relationships[rel.Alias] = RelationshipModel{
			Name:              types.StringValue(rel.Name),
			Description:       types.StringValue(rel.Description),
			AllowedCategories: types.ListValueMust(types.StringType, allowedCategories),
			AllowedTypes:      types.ListValueMust(types.StringType, allowedTypes),
		}
	}

	finalModel, err := s.NewModel(res, stateModel)
	if err != nil {
		resp.Diagnostics.AddError("error", fmt.Sprintf("unable to build resource, got error: %s", err))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (s ComponentTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[ComponentTypeModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[ComponentTypeModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	properties, err := NewPropertiesInput(planModel)
	if err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to use schema, got error: %s", err))
		return
	}

	// Update the component type first
	input := opslevel.ComponentTypeInput{
		Name:        nullable(planModel.Name.ValueStringPointer()),
		Alias:       nullable(planModel.Alias.ValueStringPointer()),
		Description: nullable(planModel.Description.ValueStringPointer()),
		Properties:  properties,
	}
	if !planModel.Icon.Color.IsNull() && !planModel.Icon.Name.IsNull() {
		input.Icon = &opslevel.ComponentTypeIconInput{
			Color: planModel.Icon.Color.ValueString(),
			Name:  opslevel.ComponentTypeIconEnum(planModel.Icon.Name.ValueString()),
		}
	}

	id := stateModel.Id.ValueString()
	res, err := s.client.UpdateComponentType(id, input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update resource, got error: %s", err))
		return
	}

	if s.reconcileRelationships(ctx, err, id, resp, planModel) {
		return
	}

	finalModel, err := s.NewModel(res, planModel)
	if err != nil {
		resp.Diagnostics.AddError("error", fmt.Sprintf("unable to build resource, got error: %s", err))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (s ComponentTypeResource) reconcileRelationships(ctx context.Context, err error, id string, resp *resource.UpdateResponse, planModel ComponentTypeModel) bool {
	// Handle relationship definitions
	// First, get existing relationship definitions
	existingRels, err := s.client.ListRelationshipDefinitions(&opslevel.PayloadVariables{
		"componentType": opslevel.NewIdentifier(id),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to list existing relationship definitions, got error: %s", err))
		return true
	}

	// Create a map of existing relationships by alias
	existingRelMap := make(map[string]opslevel.RelationshipDefinitionType)
	for _, rel := range existingRels.Nodes {
		existingRelMap[rel.Alias] = rel
	}

	// Create or update relationships from the plan
	for alias, rel := range planModel.Relationships {
		allowedCategories := make([]string, 0)
		if err := rel.AllowedCategories.ElementsAs(ctx, &allowedCategories, false); err != nil {
			resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_categories for relationship '%s': %s", alias, err))
			continue
		}

		allowedTypes := make([]string, 0)
		if err := rel.AllowedTypes.ElementsAs(ctx, &allowedTypes, false); err != nil {
			resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_types for relationship '%s': %s", alias, err))
			continue
		}

		relInput := opslevel.RelationshipDefinitionInput{
			Name:          rel.Name.ValueStringPointer(),
			Alias:         &alias,
			Description:   nullable(rel.Description.ValueStringPointer()),
			ComponentType: opslevel.NewIdentifier(id),
			Metadata: &opslevel.RelationshipDefinitionMetadataInput{
				AllowedCategories: allowedCategories,
				AllowedTypes:      allowedTypes,
			},
		}

		if existingRel, exists := existingRelMap[alias]; exists {
			// Update existing relationship
			_, err := s.client.UpdateRelationshipDefinition(string(existingRel.Id), relInput)
			if err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update relationship definition '%s', got error: %s", alias, err))
			}
			delete(existingRelMap, alias)
		} else {
			// Create new relationship
			_, err := s.client.CreateRelationshipDefinition(relInput)
			if err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create relationship definition '%s', got error: %s", alias, err))
			}
		}
	}

	// Delete any relationships that were removed
	for _, rel := range existingRelMap {
		_, err := s.client.DeleteRelationshipDefinition(string(rel.Id))
		if err != nil {
			resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete relationship definition '%s', got error: %s", rel.Alias, err))
		}
	}
	return false
}

func (s ComponentTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[ComponentTypeModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()

	// Delete the component type
	if err := s.client.DeleteComponentType(id); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete resource, got error: %s", err))
		return
	}
}

func (s ComponentTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
