package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DomainDataSource{}

const providerIssueUrl = "https://github.com/OpsLevel/terraform-provider-opslevel/issues"

func NewDomainDataSource() datasource.DataSource {
	return &DomainDataSource{}
}

// DomainDataSource defines the data source implementation.
type DomainDataSource struct {
	client *opslevel.Client
}

// DomainDataSourceModel describes the data source data model.
type DomainDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Identifier  types.String `tfsdk:"identifier"`
	Aliases     types.List   `tfsdk:"aliases"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Owner       types.String `tfsdk:"owner"`
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

func (d *DomainDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*opslevel.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *opslevel.Client, got: %T. Please report this issue to the provider developers at %s.", req.ProviderData, providerIssueUrl),
		)

		return
	}

	d.client = client
}

func (d *DomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DomainDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := d.client.GetDomain(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	data.Id = types.StringValue(string(domain.Id))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read an OpsLevel Domain data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
