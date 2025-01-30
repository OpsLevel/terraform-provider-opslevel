package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure ScorecardDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &ScorecardDataSource{}

func NewScorecardDataSourcesAll() datasource.DataSource {
	return &ScorecardDataSourcesAll{}
}

// ScorecardDataSourcesAll manages a Scorecard data source.
type ScorecardDataSourcesAll struct {
	CommonDataSourceClient
}

type scorecardDataSourceModel struct {
	AffectsOverallServiceLevels types.Bool                `tfsdk:"affects_overall_service_levels"`
	Aliases                     types.List                `tfsdk:"aliases"`
	Categories                  []categoryDataSourceModel `tfsdk:"categories"`
	Description                 types.String              `tfsdk:"description"`
	FilterId                    types.String              `tfsdk:"filter_id"`
	Id                          types.String              `tfsdk:"id"`
	Name                        types.String              `tfsdk:"name"`
	OwnerId                     types.String              `tfsdk:"owner_id"`
	PassingChecks               types.Int64               `tfsdk:"passing_checks"`
	ServiceCount                types.Int64               `tfsdk:"service_count"`
	TotalChecks                 types.Int64               `tfsdk:"total_checks"`
}

// scorecardDataSourcesAllModel describes the data source data model.
type scorecardDataSourcesAllModel struct {
	Scorecards []scorecardDataSourceModel `tfsdk:"scorecards"`
}

func NewScorecardDataSourcesAllModel(ctx context.Context, client *opslevel.Client, scorecards []opslevel.Scorecard) (scorecardDataSourcesAllModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	scorecardModels := []scorecardDataSourceModel{}
	for _, scorecard := range scorecards {
		scorecardAliases := OptionalStringListValue(scorecard.Aliases)

		categoriesModel, categoriesDiags := getCategoriesModelFromScorecard(client, &scorecard)
		diags.Append(categoriesDiags...)
		if diags.HasError() {
			return scorecardDataSourcesAllModel{}, diags
		}

		scorecardModel := scorecardDataSourceModel{
			AffectsOverallServiceLevels: types.BoolValue(scorecard.AffectsOverallServiceLevels),
			Aliases:                     scorecardAliases,
			Categories:                  categoriesModel,
			Description:                 ComputedStringValue(scorecard.Description),
			FilterId:                    ComputedStringValue(string(scorecard.Filter.Id)),
			Id:                          ComputedStringValue(string(scorecard.Id)),
			Name:                        ComputedStringValue(scorecard.Name),
			OwnerId:                     ComputedStringValue(string(scorecard.Owner.Id())),
			PassingChecks:               types.Int64Value(int64(scorecard.PassingChecks)),
			ServiceCount:                types.Int64Value(int64(scorecard.ServiceCount)),
			TotalChecks:                 types.Int64Value(int64(scorecard.TotalChecks)),
		}
		scorecardModels = append(scorecardModels, scorecardModel)
	}
	return scorecardDataSourcesAllModel{Scorecards: scorecardModels}, diags
}

func (d *ScorecardDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecards"
}

func (d *ScorecardDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Scorecard data sources",

		Attributes: map[string]schema.Attribute{
			"scorecards": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: scorecardSchemaAttrs,
				},
				Description: "List of Scorecard data sources",
				Computed:    true,
			},
		},
	}
}

func (d *ScorecardDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	scorecards, err := d.client.ListScorecards(nil)
	if err != nil || scorecards == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scorecards datasource, got error: %s", err))
		return
	}
	stateModel, diags := NewScorecardDataSourcesAllModel(ctx, d.client, scorecards.Nodes)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Scorecard data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
