package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure TeamDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

// TeamDataSource manages a Team data source.
type TeamDataSource struct {
	CommonDataSourceClient
}

// TeamDataSourceModel describes the data source data model.
type TeamDataSourceModel struct {
	Alias       types.String   `tfsdk:"alias"`
	Id          types.String   `tfsdk:"id"`
	Members     types.ListType `tfsdk:"members"`
	Name        types.String   `tfsdk:"name"`
	ParentAlias types.String   `tfsdk:"parent_alias"`
	ParentId    types.String   `tfsdk:"parent_id"`
}

func NewTeamDataSourceModel(ctx context.Context, team opslevel.Team) TeamDataSourceModel {
	return TeamDataSourceModel{
		Alias: types.StringValue(team.Alias),
		Id:    types.StringValue(string(team.Id)),
		// TODO: members
		Name:        types.StringValue(team.Name),
		ParentAlias: types.StringValue(team.ParentTeam.Alias),
		ParentId:    types.StringValue(string(team.ParentTeam.Id)),
	}
}

func (teamDataSource *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (teamDataSource *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team data source",

		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				MarkdownDescription: "The alias attached to the Team.",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of this Team.",
				Computed:    true,
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Team.",
				Computed:    true,
			},
			// TODO: members
			"parent_alias": schema.StringAttribute{
				Description: "The alias of the parent team.",
				Computed:    true,
			},
			"parent_id": schema.StringAttribute{
				Description: "The id of the parent team.",
				Computed:    true,
			},
		},
	}
}

func (teamDataSource *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var team *opslevel.Team
	if data.Alias.ValueString() != "" {
		team, err = teamDataSource.client.GetTeamWithAlias(data.Alias.ValueString())
	} else if data.Id.ValueString() != "" {
		team, err = teamDataSource.client.GetTeam(opslevel.ID(data.Id.ValueString()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to read team datasource, got error: %s", err))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to read team, got error: %s", err))
		return
	}
	if team == nil || team.Id == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to find team with alias=`%s` or id=`%s`", data.Alias.ValueString(), data.Id.ValueString()))
		return
	}

	// TODO: members

	teamDataModel := NewTeamDataSourceModel(ctx, *team)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Team data source")
	resp.Diagnostics.Append(resp.Diagnostics...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &teamDataModel)...)
}
