package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &TeamTagResource{}

var _ resource.ResourceWithImportState = &TeamTagResource{}

type TeamTagResource struct {
	CommonResourceClient
}

func NewTeamTagResource() resource.Resource {
	return &TeamTagResource{}
}

type TeamTagResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Key         types.String `tfsdk:"key"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Team        types.String `tfsdk:"team"`
	TeamAlias   types.String `tfsdk:"team_alias"`
	Value       types.String `tfsdk:"value"`
}

func NewTeamTagResourceModel(teamTag opslevel.Tag, teamIdentifier string) TeamTagResourceModel {
	teamResourceModel := TeamTagResourceModel{
		Key:   types.StringValue(teamTag.Key),
		Value: types.StringValue(teamTag.Value),
		Id:    types.StringValue(string(teamTag.Id)),
	}

	// use either the team ID or alias based on what is used in the config
	if opslevel.IsID(teamIdentifier) {
		teamResourceModel.Team = types.StringValue(teamIdentifier)
	} else {
		teamResourceModel.TeamAlias = types.StringValue(teamIdentifier)
	}

	return teamResourceModel
}

func (teamTagResource *TeamTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_tag"
}

func (teamTagResource *TeamTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Team Tag Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"key": schema.StringAttribute{
				Description: "The tag's key.",
				Required:    true,
				// TODO: can add validator
			},
			"value": schema.StringAttribute{
				Description: "The tag's value.",
				Required:    true,
			},
			"team": schema.StringAttribute{
				Description: "The id of the team that this will be added to.",
				Optional:    true,
			},
			"team_alias": schema.StringAttribute{
				Description: "The alias of the team that this will be added to.",
				Optional:    true,
			},
		},
	}
}

func (teamTagResource *TeamTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagCreateInput := opslevel.TagCreateInput{
		Type:  opslevel.RefOf(opslevel.TaggableResourceTeam),
		Key:   data.Key.ValueString(),
		Value: data.Value.ValueString(),
	}

	// use either the team ID or alias based on what is used in the config
	var teamIdentifier string
	if data.Team.ValueString() == "" && data.TeamAlias.ValueString() == "" {
		resp.Diagnostics.AddError("config error", "require at least one of: 'team', 'team_alias'")
		return
	}
	if data.Team.ValueString() != "" {
		teamIdentifier = data.Team.ValueString()
		tagCreateInput.Id = opslevel.NewID(teamIdentifier)
	} else {
		teamIdentifier = data.TeamAlias.ValueString()
		tagCreateInput.Alias = &teamIdentifier
	}

	team, err := teamTagResource.client.CreateTag(tagCreateInput)
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team (%s) tag (with key '%s'), got error: %s", teamIdentifier, data.Key.ValueString(), err))
		return
	}

	createdTeamTagResourceModel := NewTeamTagResourceModel(*team, teamIdentifier)
	createdTeamTagResourceModel.LastUpdated = timeLastUpdated()
	tflog.Trace(ctx, "created a team tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdTeamTagResourceModel)...)
}

func (teamTagResource *TeamTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// use either the team ID or alias based on what is used in the config
	var teamIdentifier string
	var team *opslevel.Team
	var err error
	if data.Team.ValueString() == "" && data.TeamAlias.ValueString() == "" {
		resp.Diagnostics.AddError("config error", "require at least one of: 'team', 'team_alias'")
		return
	}
	if data.Team.ValueString() != "" {
		teamIdentifier = data.Team.ValueString()
		team, err = teamTagResource.client.GetTeam(opslevel.ID(teamIdentifier))
	} else {
		teamIdentifier = data.TeamAlias.ValueString()
		team, err = teamTagResource.client.GetTeamWithAlias(teamIdentifier)
	}
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read team (%s), got error: %s", teamIdentifier, err))
		return
	}
	_, err = team.GetTags(teamTagResource.client, nil)
	if err != nil || team.Tags == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read tags on team (%s), got error: %s", teamIdentifier, err))
	}
	var teamTag *opslevel.Tag
	for _, readTag := range team.Tags.Nodes {
		if readTag.Key == data.Key.ValueString() {
			teamTag = &readTag
			break
		}
	}
	if teamTag == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("team tag (with key '%s') not found on team (%s)", data.Key.ValueString(), teamIdentifier))
		return
	}

	readTeamResourceModel := NewTeamTagResourceModel(*teamTag, teamIdentifier)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readTeamResourceModel)...)
}

func (teamTagResource *TeamTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// use either the team ID or alias based on what is used in the config
	var teamIdentifier string
	if data.Team.ValueString() == "" && data.TeamAlias.ValueString() == "" {
		resp.Diagnostics.AddError("config error", "require at least one of: 'team', 'team_alias'")
		return
	}
	if data.Team.ValueString() != "" {
		teamIdentifier = data.Team.ValueString()
	} else {
		teamIdentifier = data.TeamAlias.ValueString()
	}

	tagUpdateInput := opslevel.TagUpdateInput{
		Id:    opslevel.ID(data.Id.ValueString()),
		Key:   data.Key.ValueStringPointer(),
		Value: data.Value.ValueStringPointer(),
	}

	teamTag, err := teamTagResource.client.UpdateTag(tagUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update team tag (with id '%s'), got error: %s", data.Id.ValueString(), err))
		return
	}

	updatedTeamTagResourceModel := NewTeamTagResourceModel(*teamTag, teamIdentifier)
	updatedTeamTagResourceModel.LastUpdated = timeLastUpdated()
	tflog.Trace(ctx, "updated a team tag")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedTeamTagResourceModel)...)
}

func (teamTagResource *TeamTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := teamTagResource.client.DeleteTag(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to delete team tag (with id '%s'), got error: %s", data.Id.ValueString(), err))
		return
	}
	tflog.Trace(ctx, "deleted a team tag resource")
}

func (teamTagResource *TeamTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
