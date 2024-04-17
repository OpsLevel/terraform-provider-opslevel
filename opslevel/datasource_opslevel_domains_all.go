package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSourceWithConfigure = &DomainDataSourcesAll{}

func NewDomainDataSourcesAll() datasource.DataSource {
	return &DomainDataSourcesAll{}
}

// DomainDataSourcesAll manages a list of all Domain data sources.
type DomainDataSourcesAll struct {
	CommonDataSourceClient
}

// DomainDataSourcesAllModel describes the data source data model.
type DomainDataSourcesAllModel struct {
	Domains []domainDataSourceModel `tfsdk:"domains"`
}

func NewDomainDataSourcesAllModel(ctx context.Context, domains []opslevel.Domain) DomainDataSourcesAllModel {
	domainModels := []domainDataSourceModel{}
	for _, domain := range domains {
		domainModel, diags := NewDomainDataSourceModel(ctx, domain)
		if diags != nil && diags.HasError() {
			continue
		}
		domainModels = append(domainModels, domainModel)
	}
	return DomainDataSourcesAllModel{Domains: domainModels}
}

func (d *DomainDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *DomainDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List of all Domain data sources",

		Attributes: map[string]schema.Attribute{
			"domains": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: domainDatasourceSchemaAttrs,
				},
				Description: "List of domain data sources",
				Computed:    true,
			},
		},
	}
}

func (d *DomainDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel DomainDataSourcesAllModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domains, err := d.client.ListDomains(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}
	stateModel = NewDomainDataSourcesAllModel(ctx, domains.Nodes)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all OpsLevel Domain data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
