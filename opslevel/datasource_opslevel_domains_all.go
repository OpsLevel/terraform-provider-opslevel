package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
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

func NewDomainDataSourcesAllModel(ctx context.Context, domains []opslevel.Domain) (DomainDataSourcesAllModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	domainModels := []domainDataSourceModel{}
	for _, domain := range domains {
		domainModels = append(domainModels, newDomainDataSourceModel(domain))
	}
	return DomainDataSourcesAllModel{Domains: domainModels}, diags
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
	domains, err := d.client.ListDomains(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read domain, got error: %s", err))
		return
	}
	stateModel, diags := NewDomainDataSourcesAllModel(ctx, domains.Nodes)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all OpsLevel Domain data sources")
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
