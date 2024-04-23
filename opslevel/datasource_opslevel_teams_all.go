package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
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

func newTeamDataSourcesAllModel(teams []opslevel.Team) teamDataSourcesAllModel {
	teamModels := make([]teamDataSourceModel, 0)
	for _, team := range teams {
		teamModel := newTeamDataSourceModel(team)
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
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							MarkdownDescription: "The alias attached to the Team.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description: "The ID of this Team.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the Team.",
							Computed:    true,
						},
						"members": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: memberNestedSchemaAttrs,
							},
							Description: "List of team members on the team with email address and role.",
							Computed:    true,
						},
						"parent_alias": schema.StringAttribute{
							Description: "The alias of the parent team.",
							Computed:    true,
						},
						"parent_id": schema.StringAttribute{
							Description: "The id of the parent team.",
							Computed:    true,
						},
					},
				},
				Description: "List of team data sources",
				Computed:    true,
			},
		},
	}
}

func (d *TeamDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel teamDataSourcesAllModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teams, err := d.client.ListTeams(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list teams, got error: %s", err))
		return
	}
	stateModel = newTeamDataSourcesAllModel(teams.Nodes)

	tflog.Trace(ctx, "listed all OpsLevel Team data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
