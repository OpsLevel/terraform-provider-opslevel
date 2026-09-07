package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
)

var _ datasource.DataSourceWithConfigure = &TeamDataSourcesAll{}

func NewTeamDataSourcesAll() datasource.DataSource {
	return &TeamDataSourcesAll{}
}

type TeamDataSourcesAll struct {
	CommonDataSourceClient
}

type teamDataSourcesAllModel struct {
	Teams []teamDataSourceModel `tfsdk:"teams"`
}

func newTeamDataSourcesAllModel(ctx context.Context, teams []opslevel.Team, diags *diag.Diagnostics) teamDataSourcesAllModel {
	teamModels := make([]teamDataSourceModel, 0)
	for i := range teams {
		teamModel := newTeamDataSourceModel(teams[i])

		if teams[i].Properties != nil {
			propertiesModel, propDiags := NewPropertiesAllModel(ctx, teams[i].Properties.Nodes)
			diags.Append(propDiags...)
			teamModel.Properties = propertiesModel
		}

		teamModels = append(teamModels, teamModel)
	}
	return teamDataSourcesAllModel{Teams: teamModels}
}

func (d *TeamDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *TeamDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List of all Team data sources",

		Attributes: map[string]schema.Attribute{
			"teams": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: teamAttributes(map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							MarkdownDescription: "The alias attached to the Team.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description: "The ID of this Team.",
							Computed:    true,
						},
					}),
				},
				Description: "List of team data sources",
				Computed:    true,
			},
		},
	}
}

func (d *TeamDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	teams, err := d.client.ListTeams(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list teams, got error: %s", err))
		return
	}
	stateModel := newTeamDataSourcesAllModel(ctx, teams.Nodes, &resp.Diagnostics)

	tflog.Trace(ctx, "listed all OpsLevel Team data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
