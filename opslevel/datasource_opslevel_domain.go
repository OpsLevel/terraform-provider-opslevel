package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

// Ensure DomainDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &DomainDataSource{}

func NewDomainDataSource() datasource.DataSource {
	return &DomainDataSource{}
}

// DomainDataSource manages a Domain data source.
type DomainDataSource struct {
	CommonDataSourceClient
}

var domainDatasourceSchemaAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this Doamin",
		Computed:            true,
	},
	"aliases": schema.ListAttribute{
		ElementType: types.StringType,
		Description: "The aliases of the domain.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the domain.",
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "The description of the domain.",
		Computed:    true,
	},
	"owner": schema.StringAttribute{
		Description: "The id of the domain owner (team)",
		Computed:    true,
	},
}

func DomainAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range domainDatasourceSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

// domainDataSourceModel describes the data source data model.
type domainDataSourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Owner       types.String `tfsdk:"owner"`
}

// domainDataSourceModelWithIdentifier needed for a single Domain
type domainDataSourceModelWithIdentifier struct {
	Identifier types.String `tfsdk:"identifier"`
	domainDataSourceModel
}

// newDomainDataSourceModelWithIdentifier used for a single Domain
func newDomainDataSourceModelWithIdentifier(domain opslevel.Domain, identifier types.String) domainDataSourceModelWithIdentifier {
	domainDataSourceModel := newDomainDataSourceModel(domain)
	domainDataSourceModelWithIdentifier := domainDataSourceModelWithIdentifier{
		domainDataSourceModel: domainDataSourceModel,
		Identifier:            identifier,
	}
	return domainDataSourceModelWithIdentifier
}

func newDomainDataSourceModel(domain opslevel.Domain) domainDataSourceModel {
	domainAliases := OptionalStringListValue(domain.Aliases)
	domainDataSourceModel := domainDataSourceModel{
		Aliases:     domainAliases,
		Description: ComputedStringValue(domain.Description),
		Id:          ComputedStringValue(string(domain.Id)),
		Name:        ComputedStringValue(domain.Name),
		Owner:       ComputedStringValue(string(domain.Owner.Id())),
	}
	return domainDataSourceModel
}

func (d *DomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *DomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Domain data source",

		Attributes: DomainAttributes(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the domain to find.",
				Required:    true,
			},
		}),
	}
}

func (d *DomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := read[domainDataSourceModelWithIdentifier](ctx, &resp.Diagnostics, req.Config)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := d.client.GetDomain(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read domain, got error: %s", err))
		return
	}
	domainDataModel := newDomainDataSourceModelWithIdentifier(*domain, data.Identifier)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Domain data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &domainDataModel)...)
}
