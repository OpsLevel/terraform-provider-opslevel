package opslevel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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

type ComponentTypeModel struct {
	Id          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	Alias       types.String             `tfsdk:"alias"`
	Description types.String             `tfsdk:"description"`
	Icon        ComponentTypeIconModel   `tfsdk:"icon"`
	Properties  map[string]PropertyModel `tfsdk:"properties"`
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
	stateModel.Icon = ComponentTypeIconModel{
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
	finalModel, err := s.NewModel(res, stateModel)
	if err != nil {
		resp.Diagnostics.AddError("error", fmt.Sprintf("unable to build resource, got error: %s", err))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (s ComponentTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[ComponentTypeModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	if err := s.client.DeleteComponentType(id); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete resource, got error: %s", err))
		return
	}
}

func (s ComponentTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
