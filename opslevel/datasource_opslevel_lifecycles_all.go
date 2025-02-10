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

// Ensure LifecycleDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &LifecycleDataSource{}

func NewLifecycleDataSourcesAll() datasource.DataSource {
	return &LifecycleDataSourcesAll{}
}

// LifecycleDataSource manages a Lifecycle data source.
type LifecycleDataSourcesAll struct {
	CommonDataSourceClient
}

// CategoryDataSourcesAll manages a Category data source.
type lifecycleDataSourcesAllModel struct {
	Lifecycles []lifecycleDataSourceModel `tfsdk:"lifecycles"`
}

// lifecycleDataSourceModel describes the data source data model.
type lifecycleDataSourceModel struct {
	Alias types.String `tfsdk:"alias"`
	Id    types.String `tfsdk:"id"`
	Index types.Int64  `tfsdk:"index"`
	Name  types.String `tfsdk:"name"`
}

func NewLifecycleDataSourcesAllModel(lifecycles []opslevel.Lifecycle) lifecycleDataSourcesAllModel {
	lifecyclesModel := []lifecycleDataSourceModel{}
	for _, lifecycle := range lifecycles {
		lifecycleModel := lifecycleDataSourceModel{
			Alias: ComputedStringValue(lifecycle.Alias),
			Id:    ComputedStringValue(string(lifecycle.Id)),
			Index: types.Int64Value(int64(lifecycle.Index)),
			Name:  ComputedStringValue(lifecycle.Name),
		}
		lifecyclesModel = append(lifecyclesModel, lifecycleModel)
	}
	return lifecycleDataSourcesAllModel{Lifecycles: lifecyclesModel}
}

func (d *LifecycleDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycles"
}

func (d *LifecycleDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Lifecycle data source",

		Attributes: map[string]schema.Attribute{
			"lifecycles": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: lifecycleSchemaAttrs,
				},
				Description: "List of Rubric Category data sources",
				Computed:    true,
			},
		},
	}
}

func (d *LifecycleDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	lifecycles, err := d.client.ListLifecycles()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to list lifecycles, got error: %s", err))
		return
	}

	stateModel := NewLifecycleDataSourcesAllModel(lifecycles)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all lifecycle data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
