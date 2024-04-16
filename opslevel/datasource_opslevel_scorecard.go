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

func NewScorecardDataSource() datasource.DataSource {
	return &ScorecardDataSource{}
}

// ScorecardDataSource manages a Scorecard data source.
type ScorecardDataSource struct {
	CommonDataSourceClient
}

// ScorecardDataSourceModel describes the data source data model.
type ScorecardDataSourceModel struct {
	AffectsOverallServiceLevels types.Bool   `tfsdk:"affects_overall_service_levels"`
	Aliases                     types.List   `tfsdk:"aliases"`
	Description                 types.String `tfsdk:"description"`
	FilterId                    types.String `tfsdk:"filter_id"`
	Id                          types.String `tfsdk:"id"`
	Identifier                  types.String `tfsdk:"identifier"`
	Name                        types.String `tfsdk:"name"`
	OwnerId                     types.String `tfsdk:"owner_id"`
	PassingChecks               types.Int64  `tfsdk:"passing_checks"`
	ServiceCount                types.Int64  `tfsdk:"service_count"`
	TotalChecks                 types.Int64  `tfsdk:"total_checks"`
}

func NewScorecardDataSourceModel(ctx context.Context, scorecard opslevel.Scorecard, identifier string) (ScorecardDataSourceModel, diag.Diagnostics) {
	scorecardDataSourceModel := ScorecardDataSourceModel{
		AffectsOverallServiceLevels: types.BoolValue(scorecard.AffectsOverallServiceLevels),
		Description:                 types.StringValue(scorecard.Description),
		FilterId:                    types.StringValue(string(scorecard.Filter.Id)),
		Id:                          types.StringValue(string(scorecard.Id)),
		Identifier:                  types.StringValue(identifier),
		Name:                        types.StringValue(scorecard.Name),
		OwnerId:                     types.StringValue(string(scorecard.Owner.Id())),
		PassingChecks:               types.Int64Value(int64(scorecard.PassingChecks)),
		ServiceCount:                types.Int64Value(int64(scorecard.ServiceCount)),
		TotalChecks:                 types.Int64Value(int64(scorecard.ChecksCount)),
	}

	scorecardAliases, diags := types.ListValueFrom(ctx, types.StringType, scorecard.Aliases)
	scorecardDataSourceModel.Aliases = scorecardAliases

	return scorecardDataSourceModel, diags
}

func (d *ScorecardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard"
}

func (d *ScorecardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Scorecard data source",

		Attributes: map[string]schema.Attribute{
			"affects_overall_service_levels": schema.BoolAttribute{
				Description: "Specifies whether the checks on this scorecard affect services' overall maturity level.",
				Computed:    true,
			},
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The scorecard's aliases.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The scorecard's description.",
				Computed:    true,
			},
			"filter_id": schema.StringAttribute{
				MarkdownDescription: "The scorecard's filter.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the scorecard to find.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The scorecard's name.",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The scorecard's owner id.",
				Computed:    true,
			},
			"passing_checks": schema.Int64Attribute{
				Description: "The scorecard's number of checks that are passing.",
				Computed:    true,
			},
			"service_count": schema.Int64Attribute{
				Description: "The scorecard's number of services matched.",
				Computed:    true,
			},
			"total_checks": schema.Int64Attribute{
				Description: "The scorecard's total number of checks.",
				Computed:    true,
			},
		},
	}
}

func (d *ScorecardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ScorecardDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := d.client.GetScorecard(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scorecard datasource, got error: %s", err))
		return
	}
	scorecardDataModel, diags := NewScorecardDataSourceModel(ctx, *scorecard, data.Identifier.ValueString())
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Scorecard data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &scorecardDataModel)...)
}
