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

// Ensure IntegrationDataSourcesAll implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &IntegrationDataSourcesAll{}

func NewIntegrationDataSourcesAll() datasource.DataSource {
	return &IntegrationDataSourcesAll{}
}

// IntegrationDataSourcesAll manages an Integration data source.
type IntegrationDataSourcesAll struct {
	CommonDataSourceClient
}

// integrationDataSourceModel describes the data source data model.
type integrationDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// IntegrationDataSourcesAllModel describes the data source data model.
type IntegrationDataSourcesAllModel struct {
	Integrations []integrationDataSourceModel `tfsdk:"integrations"`
}

func NewIntegrationDataSourcesAllModel(integrations []opslevel.Integration) IntegrationDataSourcesAllModel {
	integrationsModel := []integrationDataSourceModel{}
	for _, integration := range integrations {
		integrationModel := integrationDataSourceModel{
			Id:   ComputedStringValue(string(integration.Id)),
			Name: ComputedStringValue(integration.Name),
		}
		integrationsModel = append(integrationsModel, integrationModel)
	}
	return IntegrationDataSourcesAllModel{Integrations: integrationsModel}
}

func (i *IntegrationDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integrations"
}

func (i *IntegrationDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Integrations data source",

		Attributes: map[string]schema.Attribute{
			"integrations": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: integrationSchemaAttrs,
				},
				Description: "List of Integration data sources",
				Computed:    true,
			},
		},
	}
}

func (i *IntegrationDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel IntegrationDataSourcesAllModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrations, err := i.client.ListIntegrations(nil)
	if err != nil || integrations == nil || integrations.Nodes == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to list integrations, got error: %s", err))
		return
	}

	foundIntegrations := *integrations
	stateModel = NewIntegrationDataSourcesAllModel(foundIntegrations.Nodes)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all integrations data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
