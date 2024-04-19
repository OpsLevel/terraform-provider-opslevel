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
var _ datasource.DataSourceWithConfigure = &FilterDataSourcesAll{}

func NewFilterDataSourcesAll() datasource.DataSource {
	return &FilterDataSourcesAll{}
}

// FilterDataSourcesAll manages a list of all Filter data sources.
type FilterDataSourcesAll struct {
	CommonDataSourceClient
}

// FilterDataSourcesAllModel describes the data source data model.
type FilterDataSourcesAllModel struct {
	Filters []filterDataSourceModel `tfsdk:"filters"`
}

func NewFilterDataSourcesAllModel(filters []opslevel.Filter) FilterDataSourcesAllModel {
	filterModels := []filterDataSourceModel{}
	for _, filter := range filters {
		filterModel := newFilterDataSourceModel(filter)
		filterModels = append(filterModels, filterModel)
	}
	return FilterDataSourcesAllModel{Filters: filterModels}
}

func (d *FilterDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filters"
}

func (d *FilterDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List of all Filter data sources",

		Attributes: map[string]schema.Attribute{
			"filters": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: filterDatasourceSchemaAttrs,
				},
				Description: "List of filter data sources",
				Computed:    true,
			},
		},
	}
}

func (d *FilterDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel FilterDataSourcesAllModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters, err := d.client.ListFilters(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read filter, got error: %s", err))
		return
	}
	stateModel = NewFilterDataSourcesAllModel(filters.Nodes)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all OpsLevel Filter data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
