package opslevel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/opslevel/terraform-provider-opslevel/internal"
)

var ComponentTypeDataSourceSchema = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "The unique name of the component type.",
		Computed:            true,
	},
	"alias": schema.StringAttribute{
		MarkdownDescription: "The unique alias of the component type.",
		Computed:            true,
	},
	"description": schema.StringAttribute{
		MarkdownDescription: "The description of the component type.",
		Computed:            true,
	},
	"icon": schema.SingleNestedAttribute{
		MarkdownDescription: "The icon associated with the component type",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"color": schema.StringAttribute{
				Description: "The color, represented as a hexcode, for the icon.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the icon in Phosphor icons for Vue, e.g. `PhBird`. See https://phosphoricons.com/ for a full list.",
				Computed:    true,
			},
		},
	},
	"properties": schema.MapNestedAttribute{
		Description: "The properties of this component type.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "The name of the property definition.",
					Computed:    true,
				},
				"description": schema.StringAttribute{
					Description: "The description of the property definition.",
					Computed:    true,
				},
				"allowed_in_config_files": schema.BoolAttribute{
					Description: "Whether or not the property is allowed to be set in opslevel.yml config files.",
					Computed:    true,
				},
				"display_status": schema.StringAttribute{
					Description: "The display status of the property.",
					Computed:    true,
				},
				"locked_status": schema.StringAttribute{
					Description: "Restricts what sources are able to assign values to this property.",
					Computed:    true,
				},
				"schema": schema.StringAttribute{
					Description: "The schema of the property.",
					Computed:    true,
				},
			},
		},
	},
}

// ComponentTypeDataSourceModel is a simplified version of ComponentTypeModel for data sources
// It excludes relationships, because opslevel-go does not yet have Relationships on ComponentType
type ComponentTypeDataSourceModel struct {
	Identifier  types.String             `tfsdk:"identifier"`
	Id          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	Alias       types.String             `tfsdk:"alias"`
	Description types.String             `tfsdk:"description"`
	Icon        *ComponentTypeIconModel  `tfsdk:"icon"`
	Properties  map[string]PropertyModel `tfsdk:"properties"`
}

func NewComponentTypeDataSourceSingle() datasource.DataSource {
	return &internal.TFDataSourceSingle[opslevel.ComponentType, ComponentTypeDataSourceModel]{
		Name:        "component_type",
		Description: "Component Type Definition Resource",
		Attributes:  ComponentTypeDataSourceSchema,
		ReadFn: func(ctx context.Context, client *opslevel.Client, identifier string) (opslevel.ComponentType, error) {
			data, err := client.GetComponentType(identifier)
			if err != nil {
				return *data, err
			}
			conn, err := data.GetProperties(client, nil)
			if err != nil {
				return *data, err
			}
			(*data).Properties = conn
			return *data, err
		},
		ToModel: func(ctx context.Context, identifier string, data opslevel.ComponentType) (ComponentTypeDataSourceModel, error) {
			model := ComponentTypeDataSourceModel{
				Identifier:  types.StringValue(identifier),
				Id:          types.StringValue(string(data.Id)),
				Name:        types.StringValue(data.Name),
				Alias:       types.StringValue(data.Aliases[0]),
				Description: types.StringValue(data.Description),
				Icon: &ComponentTypeIconModel{
					Color: types.StringValue(data.Icon.Color),
					Name:  types.StringValue(string(data.Icon.Name)),
				},
				Properties: map[string]PropertyModel{},
			}
			for _, prop := range data.Properties.Nodes {
				model.Properties[prop.Aliases[0]] = PropertyModel{
					Name:                 types.StringValue(prop.Name),
					Description:          ComputedStringValue(prop.Description),
					AllowedInConfigFiles: types.BoolValue(prop.AllowedInConfigFiles),
					DisplayStatus:        types.StringValue(string(prop.PropertyDisplayStatus)),
					LockedStatus:         types.StringValue(string(prop.LockedStatus)),
					Schema:               types.StringValue(prop.Schema.AsString()),
				}
			}
			return model, nil
		},
	}
}

func NewComponentTypeDataSourceMulti() datasource.DataSource {
	return &internal.TFDataSourceMulti[opslevel.ComponentType, ComponentTypeDataSourceModel]{
		Name:        "component_types",
		Description: "Component Type Definition Resource",
		Attributes:  ComponentTypeDataSourceSchema,
		ReadFn: func(ctx context.Context, client *opslevel.Client) ([]opslevel.ComponentType, error) {
			resp, err := client.ListComponentTypes(nil)
			if resp == nil {
				return nil, err
			}
			for i, item := range resp.Nodes {
				conn, err := item.GetProperties(client, nil)
				if err != nil {
					continue
				}
				resp.Nodes[i].Properties = conn
			}
			return resp.Nodes, err
		},
		ToModel: func(ctx context.Context, data opslevel.ComponentType) (ComponentTypeDataSourceModel, error) {
			model := ComponentTypeDataSourceModel{
				Id:          types.StringValue(string(data.Id)),
				Name:        types.StringValue(data.Name),
				Alias:       types.StringValue(data.Aliases[0]),
				Description: types.StringValue(data.Description),
				Icon: &ComponentTypeIconModel{
					Color: types.StringValue(data.Icon.Color),
					Name:  types.StringValue(string(data.Icon.Name)),
				},
				Properties: map[string]PropertyModel{},
			}
			for _, prop := range data.Properties.Nodes {
				model.Properties[prop.Aliases[0]] = PropertyModel{
					Name:                 types.StringValue(prop.Name),
					Description:          ComputedStringValue(prop.Description),
					AllowedInConfigFiles: types.BoolValue(prop.AllowedInConfigFiles),
					DisplayStatus:        types.StringValue(string(prop.PropertyDisplayStatus)),
					LockedStatus:         types.StringValue(string(prop.LockedStatus)),
					Schema:               types.StringValue(prop.Schema.AsString()),
				}
			}
			return model, nil
		},
	}
}
