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

// Ensure TierDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &TierDataSource{}

func NewTierDataSource() datasource.DataSource {
	return &TierDataSource{}
}

// TierDataSource manages a Tier data source.
type TierDataSource struct {
	CommonDataSourceClient
}

type tierDataSourceModel struct {
	Alias types.String `tfsdk:"alias"`
	Id    types.String `tfsdk:"id"`
	Index types.Int64  `tfsdk:"index"`
	Name  types.String `tfsdk:"name"`
}

func newTierDataSourceModel(tier opslevel.Tier) tierDataSourceModel {
	return tierDataSourceModel{
		Alias: types.StringValue(tier.Alias),
		Id:    types.StringValue(string(tier.Id)),
		Index: types.Int64Value(int64(tier.Index)),
		Name:  types.StringValue(tier.Name),
	}
}

type tierDataSourceModelWithFilter struct {
	Alias  types.String     `tfsdk:"alias"`
	Filter filterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Index  types.Int64      `tfsdk:"index"`
	Name   types.String     `tfsdk:"name"`
}

func newTierDataSourceModelWithFilter(tier opslevel.Tier, filter filterBlockModel) tierDataSourceModelWithFilter {
	return tierDataSourceModelWithFilter{
		Alias:  types.StringValue(tier.Alias),
		Filter: filter,
		Id:     types.StringValue(string(tier.Id)),
		Index:  types.Int64Value(int64(tier.Index)),
		Name:   types.StringValue(tier.Name),
	}
}

var tierDatasourceSchemaAttrs = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		MarkdownDescription: "The human-friendly, unique identifier for the tier.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		MarkdownDescription: "The unique identifier for the tier.",
		Computed:            true,
	},
	"index": schema.Int64Attribute{
		MarkdownDescription: "The numerical representation of the tier.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description: "The display name of the tier.",
		Computed:    true,
	},
}

func (d *TierDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tier"
}

func (d *TierDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"alias", "id", "index", "name"}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Tier data source",

		Attributes: tierDatasourceSchemaAttrs,
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (d *TierDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := read[tierDataSourceModelWithFilter](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	tiers, err := d.client.ListTiers()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tier datasource, got error: %s", err))
		return
	}

	tier, err := filterTiers(tiers, data.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to filter tier datasource, got error: %s", err))
		return
	}

	tierDataModel := newTierDataSourceModelWithFilter(*tier, data.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Tier data source")
	// resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &tierDataModel)...)
}

func filterTiers(tiers []opslevel.Tier, filter filterBlockModel) (*opslevel.Tier, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("please provide a non-empty value for filter's value")
	}
	for _, tier := range tiers {
		switch filter.Field.ValueString() {
		case "alias":
			if filter.Value.Equal(types.StringValue(tier.Alias)) {
				return &tier, nil
			}
		case "id":
			if filter.Value.Equal(types.StringValue(string(tier.Id))) {
				return &tier, nil
			}
		case "index":
			index := strconv.Itoa(tier.Index)
			if filter.Value.Equal(types.StringValue(index)) {
				return &tier, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(tier.Name)) {
				return &tier, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find tier with: %s==%s", filter.Field, filter.Value)
}
