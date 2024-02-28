package opslevel

import (
	"context"
	"fmt"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DomainDataSources{}

func NewDomainDataSources() datasource.DataSource {
	return &DomainDataSources{}
}

// DomainDataSources defines the data source implementation.
type DomainDataSources struct {
	client *opslevel.Client
}

// DomainDataSourceModel describes the data source data model.
type DomainDataSourcesModel struct {
	Ids          types.List `tfsdk:"ids"`
	Aliases      types.List `tfsdk:"aliases"`
	Names        types.List `tfsdk:"names"`
	Descriptions types.List `tfsdk:"descriptions"`
	Owners       types.List `tfsdk:"owners"`
}

func (d *DomainDataSources) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *DomainDataSources) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Domain data sources",

		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Terraform specific identifier.",
				Computed:            true,
			},
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases of the domain.",
				Computed:    true,
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The name of the domain.",
				Computed:    true,
			},
			"descriptions": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The description of the domain.",
				Computed:    true,
			},
			"owners": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The id of the domain owner - could be a group or team.",
				Computed:    true,
			},
		},
	}
}

func (d *DomainDataSources) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DomainDataSourcesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domains, err := d.client.ListDomains(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	count := len(domains.Nodes)
	aliases := make([]string, count)
	ids := make([]string, count)
	names := make([]string, count)
	descriptions := make([]string, count)
	owners := make([]string, count)
	for i, domain := range domains.Nodes {
		if len(domain.Aliases) > 0 {
			aliases[i] = domain.Aliases[0]
		}
		ids[i] = string(domain.Id)
		names[i] = domain.Name
		descriptions[i] = domain.Description
		owners[i] = string(domain.Owner.Id())
	}

	domainAliases, diag := types.ListValueFrom(ctx, types.StringType, aliases)
	data.Aliases = domainAliases
	resp.Diagnostics.Append(diag...)

	domainIds, diag := types.ListValueFrom(ctx, types.StringType, ids)
	data.Ids = domainIds
	resp.Diagnostics.Append(diag...)

	domainNames, diag := types.ListValueFrom(ctx, types.StringType, names)
	data.Names = domainNames
	resp.Diagnostics.Append(diag...)

	domainDescriptions, diag := types.ListValueFrom(ctx, types.StringType, descriptions)
	data.Descriptions = domainDescriptions
	resp.Diagnostics.Append(diag...)

	domainOwners, diag := types.ListValueFrom(ctx, types.StringType, owners)
	data.Owners = domainOwners
	resp.Diagnostics.Append(diag...)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "listed all OpsLevel Domain data sources")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func datasourceDomains() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceDomainsRead),
// 		Schema: map[string]*schema.Schema{
// 			"aliases": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"ids": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"names": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"descriptions": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"owners": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 		},
// 	}
// }

// func datasourceDomainsRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resp, err := client.ListDomains(nil)
// 	if err != nil {
// 		return err
// 	}

// 	count := len(resp.Nodes)
// 	aliases := make([]string, count)
// 	ids := make([]string, count)
// 	names := make([]string, count)
// 	descriptions := make([]string, count)
// 	owners := make([]string, count)
// 	for i, item := range resp.Nodes {
// 		if len(item.Aliases) > 0 {
// 			aliases[i] = item.Aliases[0]
// 		}
// 		ids[i] = string(item.Id)
// 		names[i] = item.Name
// 		descriptions[i] = item.Description
// 		owners[i] = string(item.Owner.Id())
// 	}

// 	d.SetId(timeID())
// 	d.Set("aliases", aliases)
// 	d.Set("ids", ids)
// 	d.Set("names", names)
// 	d.Set("descriptions", descriptions)
// 	d.Set("owners", owners)

// 	return nil
// }
