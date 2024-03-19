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

// TierDataSourceModel describes the data source data model.
type TierDataSourceModel struct {
	Alias  types.String     `tfsdk:"alias"`
	Filter FilterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Index  types.Int64      `tfsdk:"index"`
	Name   types.String     `tfsdk:"name"`
}

func NewTierDataSourceModel(ctx context.Context, tier opslevel.Tier, filter FilterBlockModel) TierDataSourceModel {
	return TierDataSourceModel{
		Alias:  types.StringValue(string(tier.Alias)),
		Filter: filter,
		Id:     types.StringValue(string(tier.Id)),
		Index:  types.Int64Value(int64(tier.Index)),
		Name:   types.StringValue(string(tier.Name)),
	}
}

func (d *TierDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tier"
}

func (d *TierDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"alias", "id", "index", "name"}
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tier data source",

		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				MarkdownDescription: "Terraform specific identifier.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Terraform specific identifier.",
				Computed:            true,
			},
			"index": schema.Int64Attribute{
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

func (d *TierDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TierDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
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

	tierDataModel := NewTierDataSourceModel(ctx, *tier, data.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Tier data source")
	// resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &tierDataModel)...)
}

// func filterTiers(levels []opslevel.Tier, field string, value string) (*opslevel.Tier, error) {
func filterTiers(tiers []opslevel.Tier, filter FilterBlockModel) (*opslevel.Tier, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
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
			index := strconv.Itoa(int(tier.Index))
			if filter.Value.Equal(types.StringValue(index)) {
				return &tier, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(tier.Name)) {
				return &tier, nil
			}
		}
	}

	return nil, fmt.Errorf("Unable to find tier with: %s==%s", filter.Field, filter.Value)
}
