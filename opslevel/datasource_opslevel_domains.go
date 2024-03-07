package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSourceWithConfigure = &DomainDataSources{}

func NewDomainDataSources() datasource.DataSource {
	return &DomainDataSources{}
}

// DomainDataSources defines the data source implementation.
type DomainDataSources struct {
	CommonClient
}

var domainObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"aliases":     types.ListType{ElemType: types.StringType},
		"description": types.StringType,
		"id":          types.StringType,
		"name":        types.StringType,
		"owner":       types.StringType,
	},
}

// DomainDataSourceModel describes the data source data model.
type DomainDataSourcesModel struct {
	Domains types.List `tfsdk:"domains"`
}

func (d *DomainDataSources) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *DomainDataSources) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List of all Domain data sources",

		Attributes: map[string]schema.Attribute{
			"domains": schema.ListAttribute{
				ElementType: domainObjectType,
				Description: "List of domain data sources",
				Computed:    true,
			},
		},
	}
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
	parsedDomains, diags := parseDomains(ctx, domains.Nodes)
	resp.Diagnostics.Append(diags...)

	data.Domains = *parsedDomains

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "listed all OpsLevel Domain data sources")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseDomains(ctx context.Context, opslevelDomains []opslevel.Domain) (*basetypes.ListValue, diag.Diagnostics) {
	domains := make([]attr.Value, len(opslevelDomains))

	for i, domain := range opslevelDomains {
		domainModel := make(map[string]attr.Value)

		if len(domain.Aliases) > 0 {
			aliases, diags := types.ListValueFrom(ctx, types.StringType, domain.Aliases)
			if diags.HasError() {
				return nil, diags
			}
			domainModel["aliases"] = aliases
		} else {
			domainModel["aliases"] = types.ListNull(types.StringType)
		}

		domainModel["description"] = types.StringValue(string(domain.Description))
		domainModel["id"] = types.StringValue(string(domain.Id))
		domainModel["name"] = types.StringValue(string(domain.Name))
		domainModel["owner"] = types.StringValue(string(domain.Owner.Id()))

		domainObject, diags := types.ObjectValue(domainObjectType.AttrTypes, domainModel)
		if diags.HasError() {
			return nil, diags
		}

		domains[i] = domainObject
	}

	result, diags := basetypes.NewListValue(
		domainObjectType,
		domains,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}
