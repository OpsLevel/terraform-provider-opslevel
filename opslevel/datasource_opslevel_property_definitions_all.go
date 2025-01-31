package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure PropertyDefinitionDataSourcesAll implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &PropertyDefinitionDataSourcesAll{}

func NewPropertyDefinitionDataSourcesAll() datasource.DataSource {
	return &PropertyDefinitionDataSourcesAll{}
}

// PropertyDefinitionDataSourcesAll manages a PropertyDefinition data source.
type PropertyDefinitionDataSourcesAll struct {
	CommonDataSourceClient
}

// propertyDefinitionDataSourcesAllModel describes the data source data model.
type propertyDefinitionDataSourcesAllModel struct {
	PropertyDefinitions []propertyDefinitionDataSourceModel `tfsdk:"property_definitions"`
}

func NewPropertyDefinitionDataSourcesAllModel(propertyDefinitions []opslevel.PropertyDefinition) propertyDefinitionDataSourcesAllModel {
	propDefinitionsModel := []propertyDefinitionDataSourceModel{}
	for _, propertyDefinition := range propertyDefinitions {
		propDefinitionModel := propertyDefinitionDataSourceModel{
			Description:           ComputedStringValue(propertyDefinition.Description),
			Id:                    ComputedStringValue(string(propertyDefinition.Id)),
			Name:                  ComputedStringValue(propertyDefinition.Name),
			PropertyDisplayStatus: ComputedStringValue(string(propertyDefinition.PropertyDisplayStatus)),
			LockedStatus:          ComputedStringValue(string(propertyDefinition.LockedStatus)),
			Schema:                ComputedStringValue(propertyDefinition.Schema.AsString()),
		}
		propDefinitionsModel = append(propDefinitionsModel, propDefinitionModel)
	}
	return propertyDefinitionDataSourcesAllModel{PropertyDefinitions: propDefinitionsModel}
}

type propertyDefinitionDataSourceModel struct {
	Description           types.String `tfsdk:"description"`
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	PropertyDisplayStatus types.String `tfsdk:"property_display_status"`
	LockedStatus          types.String `tfsdk:"locked_status"`
	Schema                types.String `tfsdk:"schema"`
}

var propertyDefinitionSchemaAttrs = map[string]schema.Attribute{
	"description": schema.StringAttribute{
		MarkdownDescription: "The description of the property definition.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "The display name of the property definition.",
		Computed:            true,
	},
	"property_display_status": schema.StringAttribute{
		MarkdownDescription: "The display status of a custom property on service pages. (Options: 'visible' or 'hidden')",
		Computed:            true,
	},
	"locked_status": schema.StringAttribute{
		MarkdownDescription: "Restricts what sources are able to assign values to this property. (Options: 'unlocked' or 'ui_locked')",
		Computed:            true,
	},
	"schema": schema.StringAttribute{
		MarkdownDescription: "The schema of the property definition.",
		Computed:            true,
	},
}

func (d *PropertyDefinitionDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_definitions"
}

func (d *PropertyDefinitionDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Property Definitions",

		Attributes: map[string]schema.Attribute{
			"property_definitions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: propertyDefinitionSchemaAttrs,
				},
				Description: "List of Property Definition data sources",
				Computed:    true,
			},
		},
	}
}

func (d *PropertyDefinitionDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	propertyDefinitions, err := d.client.ListPropertyDefinitions(nil)
	if err != nil || propertyDefinitions == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read property definition datasource, got error: %s", err))
		return
	}
	stateModel := NewPropertyDefinitionDataSourcesAllModel(propertyDefinitions.Nodes)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel PropertyDefinition data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
