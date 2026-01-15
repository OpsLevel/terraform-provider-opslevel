package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
)

// Ensure FilterDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &FilterDataSource{}

func NewFilterDataSource() datasource.DataSource {
	return &FilterDataSource{}
}

// FilterDataSource manages a Filter data source.
type FilterDataSource struct {
	CommonDataSourceClient
}

// filterDataSourceModel describes the data source data model.
type filterDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// filterDataSourceWithFilterModel contains filterDataSourceModel fields and a filterBlockModel
type filterDataSourceWithFilterModel struct {
	Filter filterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Name   types.String     `tfsdk:"name"`
}

func newFilterDataSourceModel(opslevelFilter opslevel.Filter) filterDataSourceModel {
	return filterDataSourceModel{
		Id:   ComputedStringValue(string(opslevelFilter.Id)),
		Name: ComputedStringValue(opslevelFilter.Name),
	}
}

func newFilterDataSourceWithFilterModel(opslevelFilter opslevel.Filter, filterModel filterBlockModel) filterDataSourceWithFilterModel {
	return filterDataSourceWithFilterModel{
		Filter: filterModel,
		Id:     ComputedStringValue(string(opslevelFilter.Id)),
		Name:   ComputedStringValue(opslevelFilter.Name),
	}
}

func (d *FilterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filter"
}

var filterDatasourceSchemaAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this filter.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the filter.",
		Computed:    true,
	},
}

func (d *FilterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"id", "name"}
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Filter data source",

		Attributes: filterDatasourceSchemaAttrs,
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (d *FilterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := read[filterDataSourceWithFilterModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	opslevelFilters, err := d.client.ListFilters(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read filter datasource, got error: %s", err))
		return
	}

	opslevelFilter, err := filterOpsLevelFilters(opslevelFilters.Nodes, data.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read filter datasource, got error: %s", err))
		return
	}

	filterDataModel := newFilterDataSourceWithFilterModel(*opslevelFilter, data.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Filter data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &filterDataModel)...)
}

func filterOpsLevelFilters(opslevelFilters []opslevel.Filter, filter filterBlockModel) (*opslevel.Filter, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("please provide a non-empty value for filter's value")
	}
	for _, opslevelFilter := range opslevelFilters {
		switch filter.Field.ValueString() {
		case "id":
			if filter.Value.Equal(types.StringValue(string(opslevelFilter.Id))) {
				return &opslevelFilter, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(opslevelFilter.Name)) {
				return &opslevelFilter, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find filter with: %s==%s", filter.Field, filter.Value)
}
