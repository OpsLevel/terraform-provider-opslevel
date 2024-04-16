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

// Ensure IntegrationDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &IntegrationDataSource{}

func NewIntegrationDataSource() datasource.DataSource {
	return &IntegrationDataSource{}
}

// IntegrationDataSource manages an Integration data source.
type IntegrationDataSource struct {
	CommonDataSourceClient
}

// IntegrationDataSourceModel describes the data source data model.
type IntegrationDataSourceModel struct {
	Filter FilterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Name   types.String     `tfsdk:"name"`
}

func NewIntegrationDataSourceModel(ctx context.Context, integration opslevel.Integration, filter FilterBlockModel) IntegrationDataSourceModel {
	return IntegrationDataSourceModel{
		Filter: filter,
		Id:     types.StringValue(string(integration.Id)),
		Name:   types.StringValue(integration.Name),
	}
}

func (i *IntegrationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (i *IntegrationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"id", "name"}
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Integration data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this Integration.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Integration.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (i *IntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IntegrationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrations, err := i.client.ListIntegrations(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to list integrations, got error: %s", err))
		return
	}

	integration, err := filterIntegrations(integrations.Nodes, data.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to filter integrations, got error: %s", err))
		return
	}

	integrationDataModel := NewIntegrationDataSourceModel(ctx, *integration, data.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Integration data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &integrationDataModel)...)
}

func filterIntegrations(data []opslevel.Integration, filter FilterBlockModel) (*opslevel.Integration, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("please provide a non-empty value for filter's value")
	}
	for _, integration := range data {
		switch filter.Field.ValueString() {
		case "id":
			if filter.Value.Equal(types.StringValue(string(integration.Id))) {
				return &integration, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(integration.Name)) {
				return &integration, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find integration with: %s==%s", filter.Field, filter.Value)
}
