package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ datasource.DataSourceWithConfigure = &TierDataSourcesAll{}

func NewTierDataSourcesAll() datasource.DataSource {
	return &TierDataSourcesAll{}
}

type TierDataSourcesAll struct {
	CommonDataSourceClient
}

type tierDataSourcesAllModel struct {
	Tiers []tierDataSourceModel `tfsdk:"tiers"`
}

func newTierDataSourcesAllModel(tiers []opslevel.Tier) tierDataSourcesAllModel {
	tierModels := make([]tierDataSourceModel, 0)
	for _, tier := range tiers {
		tierModel := newTierDataSourceModel(tier)
		tierModels = append(tierModels, tierModel)
	}
	return tierDataSourcesAllModel{Tiers: tierModels}
}

func (d *TierDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tiers"
}

func (d *TierDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List of all Tier data sources",

		Attributes: map[string]schema.Attribute{
			"tiers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: tierDatasourceSchemaAttrs,
				},
				Description: "List of tier data sources",
				Computed:    true,
			},
		},
	}
}

func (d *TierDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel tierDataSourcesAllModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tiers, err := d.client.ListTiers()
	if opslevel.HasBadHttpStatus(err) {
		resp.Diagnostics.AddError("HTTP status error", fmt.Sprintf("Unable to list tiers, got error: %s", err))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list tiers, got error: %s", err))
		return
	}
	stateModel = newTierDataSourcesAllModel(tiers)

	tflog.Trace(ctx, "listed all OpsLevel Tier data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
