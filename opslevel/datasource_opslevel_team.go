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

// teamDataSourceModel describes the data source data model.
type teamDataSourceModel struct {
	Alias       types.String      `tfsdk:"alias"`
	Id          types.String      `tfsdk:"id"`
	Members     []teamMemberModel `tfsdk:"members"`
	Name        types.String      `tfsdk:"name"`
	ParentAlias types.String      `tfsdk:"parent_alias"`
	ParentId    types.String      `tfsdk:"parent_id"`
}

var teamSchemaAttrs = map[string]schema.Attribute{
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
}

var memberNestedSchemaAttrs = map[string]schema.Attribute{
	"email": schema.StringAttribute{
		MarkdownDescription: "The email address of the team member.",
		Computed:            true,
	},
	"role": schema.StringAttribute{
		MarkdownDescription: "The role of the team member.",
		Computed:            true,
	},
}

type teamMemberModel struct {
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func newTeamMemberModel(member opslevel.TeamMembership) teamMemberModel {
	return teamMemberModel{
		Email: ComputedStringValue(member.User.Email),
		Role:  ComputedStringValue(member.Role),
	}
}

func newTeamMembersAllModel(members []opslevel.TeamMembership) []teamMemberModel {
	membersModel := make([]teamMemberModel, 0)
	for _, member := range members {
		membersModel = append(membersModel, newTeamMemberModel(member))
	}
	return membersModel
}

func newTeamDataSourceModel(team opslevel.Team) teamDataSourceModel {
	teamDataSourceModel := teamDataSourceModel{
		Alias:       ComputedStringValue(team.Alias),
		Id:          ComputedStringValue(string(team.Id)),
		Name:        ComputedStringValue(team.Name),
		ParentAlias: ComputedStringValue(team.ParentTeam.Alias),
		ParentId:    ComputedStringValue(string(team.ParentTeam.Id)),
	}
	if team.Memberships != nil {
		teamDataSourceModel.Members = newTeamMembersAllModel(team.Memberships.Nodes)
	}
	return teamDataSourceModel
}

func (teamDataSource *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (teamDataSource *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team data source",

		Attributes: teamSchemaAttrs,
	}
}

func (teamDataSource *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data teamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var team *opslevel.Team
	if data.Alias.ValueString() != "" {
		team, err = teamDataSource.client.GetTeamWithAlias(data.Alias.ValueString())
	} else if opslevel.IsID(data.Id.ValueString()) {
		team, err = teamDataSource.client.GetTeam(opslevel.ID(data.Id.ValueString()))
	} else {
		resp.Diagnostics.AddError("Config Error", "'alias' or 'id' for opslevel_team datasource must be set")
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

	teamDataModel := newTeamDataSourceModel(*team)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Team data source")
	resp.Diagnostics.Append(resp.Diagnostics...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &teamDataModel)...)
}
