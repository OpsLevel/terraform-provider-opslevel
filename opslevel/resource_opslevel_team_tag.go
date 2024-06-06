package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Tag key names are stored in OpsLevel as lowercase so we need to ensure the configuration is written as lowercase
var (
	TagKeyRegex    = regexp.MustCompile(`\A[a-z][0-9a-z_\.\/\\-]*\z`)
	TagKeyErrorMsg = "tag key name must start with a letter and be only lowercase alphanumerics, underscores, hyphens, periods, and slashes."
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
	Id        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Team      types.String `tfsdk:"team"`
	TeamAlias types.String `tfsdk:"team_alias"`
	Value     types.String `tfsdk:"value"`
}

func NewTeamTagResourceModel(team opslevel.Team, teamTag opslevel.Tag) TeamTagResourceModel {
	return TeamTagResourceModel{
		Key:       RequiredStringValue(teamTag.Key),
		Value:     RequiredStringValue(teamTag.Value),
		Id:        ComputedStringValue(fmt.Sprintf("%s:%s", string(team.Id), string(teamTag.Id))),
		Team:      OptionalStringValue(string(team.Id)),
		TeamAlias: OptionalStringValue(team.Alias),
	}
}

func (teamTagResource *TeamTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_tag"
}

func (teamTagResource *TeamTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Team Tag Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource. Formatted as <team-id>:<tag-id>.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The tag's key.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(TagKeyRegex, TagKeyErrorMsg),
				},
			},
			"value": schema.StringAttribute{
				Description: "The tag's value.",
				Required:    true,
			},
			"team": schema.StringAttribute{
				Description: "The id of the team that this will be added to.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					IdStringValidator(),
					stringvalidator.AtLeastOneOf(path.MatchRoot("team"),
						path.MatchRoot("team_alias")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team_alias": schema.StringAttribute{
				Description: "The alias of the team that this will be added to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (teamTagResource *TeamTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel TeamTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, getTeamDiag := getTeamWithIdOrAlias(*teamTagResource.client, planModel.Team.ValueString(), planModel.TeamAlias.ValueString())
	resp.Diagnostics.Append(getTeamDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagCreateInput := opslevel.TagCreateInput{
		Id:    &team.Id,
		Type:  opslevel.RefOf(opslevel.TaggableResourceTeam),
		Key:   planModel.Key.ValueString(),
		Value: planModel.Value.ValueString(),
	}

	teamTag, err := teamTagResource.client.CreateTag(tagCreateInput)
	if err != nil || teamTag == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team (%s) tag (with key '%s'), got error: %s", team.Id, planModel.Key.ValueString(), err))
		return
	}

	createdTeamTagResourceModel := NewTeamTagResourceModel(*team, *teamTag)
	tflog.Trace(ctx, "created a team tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdTeamTagResourceModel)...)
}

func (teamTagResource *TeamTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel TeamTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	teamAndTagId := stateModel.Id.ValueString()
	if !isTagValid(teamAndTagId) {
		resp.Diagnostics.AddError(
			"State Error - invalid Id format",
			fmt.Sprintf("Id expected to be formatted as '<team-id>:<tag-id>'. State has '%s'", teamAndTagId),
		)
	}
	ids := strings.Split(teamAndTagId, ":")
	teamTagId := ids[1]

	// use either the team ID or alias based on what is used in the config
	team, getTeamDiag := getTeamWithIdOrAlias(*teamTagResource.client, ids[0], stateModel.TeamAlias.ValueString())
	resp.Diagnostics.Append(getTeamDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := team.GetTags(teamTagResource.client, nil)
	if err != nil || team.Tags == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read tags on team (%s), got error: %s", team.Id, err))
	}
	var teamTag *opslevel.Tag
	for _, readTag := range team.Tags.Nodes {
		if teamTagId == string(readTag.Id) || readTag.Key == stateModel.Key.ValueString() {
			teamTag = &readTag
			break
		}
	}
	if teamTag == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("team tag (with key '%s') not found on team (%s)", stateModel.Key.ValueString(), team.Id))
		return
	}

	readTeamResourceModel := NewTeamTagResourceModel(*team, *teamTag)
	// use either the team ID or alias based on what is used in the config
	resp.Diagnostics.Append(resp.State.Set(ctx, &readTeamResourceModel)...)
}

func (teamTagResource *TeamTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel TeamTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, getTeamDiag := getTeamWithIdOrAlias(*teamTagResource.client, planModel.Team.ValueString(), planModel.TeamAlias.ValueString())
	resp.Diagnostics.Append(getTeamDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := strings.Split(planModel.Id.ValueString(), ":")
	teamTagId := ids[1]
	tagUpdateInput := opslevel.TagUpdateInput{
		Id:    opslevel.ID(teamTagId),
		Key:   planModel.Key.ValueStringPointer(),
		Value: planModel.Value.ValueStringPointer(),
	}

	teamTag, err := teamTagResource.client.UpdateTag(tagUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update team tag (with id '%s'), got error: %s", teamTagId, err))
		return
	}

	updatedTeamTagResourceModel := NewTeamTagResourceModel(*team, *teamTag)
	tflog.Trace(ctx, "updated a team tag")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedTeamTagResourceModel)...)
}

func (teamTagResource *TeamTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := strings.Split(data.Id.ValueString(), ":")
	teamTagId := ids[1]
	err := teamTagResource.client.DeleteTag(opslevel.ID(teamTagId))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to delete team tag (with id '%s'), got error: %s", data.Id.ValueString(), err))
		return
	}
	tflog.Trace(ctx, "deleted a team tag resource")
}

func (teamTagResource *TeamTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !isTagValid(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<team-id>:<tag-id>'. Given '%s'", req.ID),
		)
	}
	// put these into req - this is passed to Read() - it works
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getTeamWithIdOrAlias(client opslevel.Client, teamId, teamAlias string) (*opslevel.Team, diag.Diagnostics) {
	var diag diag.Diagnostics
	var err error
	var team *opslevel.Team

	if opslevel.IsID(teamId) {
		team, err = client.GetTeam(opslevel.ID(teamId))
	} else {
		team, err = client.GetTeamWithAlias(teamAlias)
	}
	if err != nil {
		diag.AddError(
			"opslevel client error",
			fmt.Sprintf("unable to get team with id: '%s' or alias: '%s', got error: %s",
				teamId,
				teamAlias,
				err),
		)
	}
	return team, diag
}
