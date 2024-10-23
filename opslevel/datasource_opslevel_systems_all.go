package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ datasource.DataSourceWithConfigure = &SystemDataSourcesAll{}

func NewSystemDataSourcesAll() datasource.DataSource {
	return &SystemDataSourcesAll{}
}

type SystemDataSourcesAll struct {
	CommonDataSourceClient
}

type SystemDataSourcesAllModel struct {
	Systems []systemDataSourceModel `tfsdk:"systems"`
}

func NewSystemDataSourcesAllModel(ctx context.Context, systems []opslevel.System) (SystemDataSourcesAllModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	systemModels := make([]systemDataSourceModel, 0)
	for _, system := range systems {
		systemModels = append(systemModels, newSystemDataSourceModel(system))
	}
	return SystemDataSourcesAllModel{Systems: systemModels}, diags
}

func (d *SystemDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_systems"
}

func (d *SystemDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List of all System data sources",

		Attributes: map[string]schema.Attribute{
			"systems": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: systemDatasourceSchemaAttrs,
				},
				Description: "List of system data sources",
				Computed:    true,
			},
		},
	}
}

func (d *SystemDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	systems, err := d.client.ListSystems(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read system, got error: %s", err))
		return
	}
	stateModel, diags := NewSystemDataSourcesAllModel(ctx, systems.Nodes)

	tflog.Trace(ctx, "listed all OpsLevel System data sources")
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
