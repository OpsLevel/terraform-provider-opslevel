package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
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

// FilterDataSourceModel describes the data source data model.
type FilterDataSourceModel struct {
	Filter FilterModel  `tfsdk:"filter"`
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
}

func NewFilterDataSourceModel(ctx context.Context, opslevelFilter opslevel.Filter, filterModel FilterModel) FilterDataSourceModel {
	return FilterDataSourceModel{
		Name:   types.StringValue(opslevelFilter.Name),
		Id:     types.StringValue(string(opslevelFilter.Id)),
		Filter: filterModel,
	}
}

func (d *FilterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filter"
}

func (d *FilterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"id", "name"}
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Filter data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Terraform specific identifier.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the domain.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (d *FilterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FilterDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
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

	filterDataModel := NewFilterDataSourceModel(ctx, *opslevelFilter, data.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Filter data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &filterDataModel)...)
}

func filterOpsLevelFilters(opslevelFilters []opslevel.Filter, filter FilterModel) (*opslevel.Filter, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
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

	return nil, fmt.Errorf("Unable to find filter with: %s==%s", filter.Field, filter.Value)
}
