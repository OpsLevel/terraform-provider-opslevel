package opslevel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure DomainDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &DomainDataSource{}

func NewDomainDataSource() datasource.DataSource {
	return &DomainDataSource{}
}

// DomainDataSource manages a Domain data source.
type DomainDataSource struct {
	CommonClient
}

// DomainModel is the model for a domain object.
type DomainModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Owner       types.String `tfsdk:"owner"`
}

func NewDomainModel(ctx context.Context, domain opslevel.Domain) (DomainModel, diag.Diagnostics) {
	var domainDataSourceModel DomainModel

	domainDataSourceModel.Id = types.StringValue(string(domain.Id))
	domainAliases, diags := types.ListValueFrom(ctx, types.StringType, domain.Aliases)

	domainDataSourceModel.Aliases = domainAliases
	domainDataSourceModel.Name = types.StringValue(string(domain.Name))
	domainDataSourceModel.Description = types.StringValue(string(domain.Description))
	domainDataSourceModel.Owner = types.StringValue(string(domain.Owner.Id()))
	return domainDataSourceModel, diags
}

func (d *DomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *DomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Domain data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Terraform specific identifier.",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the domain to find.",
				Optional:    true,
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
				Description: "The id of the domain owner - could be a group or team.",
				Computed:    true,
			},
		},
	}
}

func (d *DomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var identifier string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("identifier"), &identifier)...)
	var data DomainModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := d.client.GetDomain(identifier)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}
	model, diags := NewDomainModel(ctx, *domain)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Domain data source")
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), identifier)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
