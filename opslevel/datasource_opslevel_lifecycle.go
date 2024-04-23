package opslevel

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure LifecycleDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &LifecycleDataSource{}

func NewLifecycleDataSource() datasource.DataSource {
	return &LifecycleDataSource{}
}

// LifecycleDataSource manages a Lifecycle data source.
type LifecycleDataSource struct {
	CommonDataSourceClient
}

// lifecycleDataSourceWithFilterModel describes the data source data model.
type lifecycleDataSourceWithFilterModel struct {
	Alias  types.String     `tfsdk:"alias"`
	Filter filterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Index  types.Int64      `tfsdk:"index"`
	Name   types.String     `tfsdk:"name"`
}

func NewLifecycleDataSourceModel(ctx context.Context, lifecycle opslevel.Lifecycle, filter filterBlockModel) lifecycleDataSourceWithFilterModel {
	return lifecycleDataSourceWithFilterModel{
		Alias:  ComputedStringValue(lifecycle.Alias),
		Filter: filter,
		Id:     ComputedStringValue(string(lifecycle.Id)),
		Index:  types.Int64Value(int64(lifecycle.Index)),
		Name:   ComputedStringValue(lifecycle.Name),
	}
}

var lifecycleSchemaAttrs = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		MarkdownDescription: "The alias attached to the Lifecycle.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		Description: "The unique identifier for the Lifecycle.",
		Computed:    true,
	},
	"index": schema.Int64Attribute{
		Description: "The numerical representation of the Lifecycle.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the Lifecycle.",
		Computed:    true,
	},
}

func (lifecycleDataSource *LifecycleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle"
}

func (lifecycleDataSource *LifecycleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"alias", "id", "index", "name"}
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Lifecycle data source",

		Attributes: lifecycleSchemaAttrs,
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (lifecycleDataSource *LifecycleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel lifecycleDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lifecycles, err := lifecycleDataSource.client.ListLifecycles()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to list lifecycles, got error: %s", err))
		return
	}

	lifecycle, err := filterLifecycles(lifecycles, planModel.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to filter lifecycle, got error: %s", err))
		return
	}

	stateModel = NewLifecycleDataSourceModel(ctx, *lifecycle, planModel.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Lifecycle data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func filterLifecycles(lifecycles []opslevel.Lifecycle, filter filterBlockModel) (*opslevel.Lifecycle, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("please provide a non-empty value for lifecycle's value")
	}
	for _, lifecycle := range lifecycles {
		switch filter.Field.ValueString() {
		case "alias":
			if filter.Value.Equal(types.StringValue(lifecycle.Alias)) {
				return &lifecycle, nil
			}
		case "id":
			if filter.Value.Equal(types.StringValue(string(lifecycle.Id))) {
				return &lifecycle, nil
			}
		case "index":
			index := strconv.Itoa(lifecycle.Index)
			if filter.Value.Equal(types.StringValue(index)) {
				return &lifecycle, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(lifecycle.Name)) {
				return &lifecycle, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find lifecycle with: %s==%s", filter.Field, filter.Value)
}
