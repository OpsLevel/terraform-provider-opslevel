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

// Ensure PropertyDefinitionDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &PropertyDefinitionDataSource{}

func NewPropertyDefinitionDataSource() datasource.DataSource {
	return &PropertyDefinitionDataSource{}
}

// PropertyDefinitionDataSource manages a PropertyDefinition data source.
type PropertyDefinitionDataSource struct {
	CommonDataSourceClient
}

func PropertyDefinitionAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range propertyDefinitionSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

// propertyDefinitionDataSourceWithFilterModel describes the data source data model.
type propertyDefinitionDataSourceWithFilterModel struct {
	AllowedInConfigFiles  types.Bool   `tfsdk:"allowed_in_config_files"`
	Description           types.String `tfsdk:"description"`
	Id                    types.String `tfsdk:"id"`
	Identifier            types.String `tfsdk:"identifier"`
	Name                  types.String `tfsdk:"name"`
	PropertyDisplayStatus types.String `tfsdk:"property_display_status"`
	Schema                types.String `tfsdk:"schema"`
}

func NewPropertyDefinitionDataSourceWithFilterModel(propertydefinition opslevel.PropertyDefinition, identifier string) propertyDefinitionDataSourceWithFilterModel {
	return propertyDefinitionDataSourceWithFilterModel{
		AllowedInConfigFiles:  types.BoolValue(propertydefinition.AllowedInConfigFiles),
		Description:           types.StringValue(propertydefinition.Description),
		Id:                    types.StringValue(string(propertydefinition.Id)),
		Identifier:            types.StringValue(identifier),
		Name:                  types.StringValue(propertydefinition.Name),
		PropertyDisplayStatus: types.StringValue(string(propertydefinition.PropertyDisplayStatus)),
		Schema:                types.StringValue(propertydefinition.Schema.ToJSON()),
	}
}

func (d *PropertyDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_definition"
}

func (d *PropertyDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Property Definitions",

		Attributes: PropertyDefinitionAttributes(map[string]schema.Attribute{
			"allowed_in_config_files": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the property is allowed to be set in opslevel.yml config files.",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The id or alias of the property definition",
				Required:            true,
			},
		}),
	}
}

func (d *PropertyDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel propertyDefinitionDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	propertydefinition, err := d.client.GetPropertyDefinition(planModel.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read property definition datasource, got error: %s", err))
		return
	}
	stateModel = NewPropertyDefinitionDataSourceWithFilterModel(*propertydefinition, planModel.Identifier.ValueString())

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel PropertyDefinition data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
