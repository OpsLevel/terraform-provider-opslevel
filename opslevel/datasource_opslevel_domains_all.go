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
var _ datasource.DataSourceWithConfigure = &DomainDataSourcesAll{}

func NewDomainDataSourcesAll() datasource.DataSource {
	return &DomainDataSourcesAll{}
}

// DomainDataSourcesAll manages a list of all Domain data sources.
type DomainDataSourcesAll struct {
	CommonClient
}

// domainObjectType is derived from DomainDataSourceModel, needed for lists
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

func (d *DomainDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *DomainDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *DomainDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
	parsedDomains, diags := parseAllDomains(ctx, domains.Nodes)
	resp.Diagnostics.Append(diags...)

	data.Domains = *parsedDomains

	tflog.Trace(ctx, "listed all OpsLevel Domain data sources")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func parseAllDomains(ctx context.Context, opslevelDomains []opslevel.Domain) (*basetypes.ListValue, diag.Diagnostics) {
	domains := make([]attr.Value, len(opslevelDomains))

	for i, domain := range opslevelDomains {
		domainObject, diags := domainToObject(ctx, domain)
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

// domainToObject converts an opslevel.Domain to a basetypes.ObjectValue
func domainToObject(ctx context.Context, opslevelDomain opslevel.Domain) (basetypes.ObjectValue, diag.Diagnostics) {
	domainObject, diags := NewDomainDataSourceModel(ctx, opslevelDomain)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	domainModel := make(map[string]attr.Value)
	domainModel["aliases"] = domainObject.Aliases
	domainModel["description"] = domainObject.Description
	domainModel["id"] = domainObject.Id
	domainModel["name"] = domainObject.Name
	domainModel["owner"] = domainObject.Owner

	parsedDomain, diags := types.ObjectValue(domainObjectType.AttrTypes, domainModel)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	return parsedDomain, nil
}
