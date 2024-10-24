package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

func NewTeamTagResourceModel(teamTag opslevel.Tag) TeamTagResourceModel {
	teamResourceModel := TeamTagResourceModel{
		Key:   RequiredStringValue(teamTag.Key),
		Value: RequiredStringValue(teamTag.Value),
		Id:    ComputedStringValue(string(teamTag.Id)),
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
				Optional:    true,
				Validators: []validator.String{
					IdStringValidator(),
					stringvalidator.ExactlyOneOf(path.MatchRoot("team"),
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
	data := read[TeamTagResourceModel](ctx, &resp.Diagnostics, req.Plan)
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

	createdTeamTagResourceModel := NewTeamTagResourceModel(*team)
	// use either the team ID or alias based on what is used in the config
	if opslevel.IsID(teamIdentifier) {
		createdTeamTagResourceModel.Team = OptionalStringValue(teamIdentifier)
	} else {
		createdTeamTagResourceModel.TeamAlias = OptionalStringValue(teamIdentifier)
	}
	tflog.Trace(ctx, "created a team tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdTeamTagResourceModel)...)
}

func (teamTagResource *TeamTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := read[TeamTagResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	// use either the team ID or alias based on what is used in the config
	var teamIdentifier string
	var team *opslevel.Team
	var err error
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
	if teamTag == nil || teamTag.Id == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	readTeamResourceModel := NewTeamTagResourceModel(*teamTag)
	// use either the team ID or alias based on what is used in the config
	if opslevel.IsID(teamIdentifier) {
		readTeamResourceModel.Team = OptionalStringValue(teamIdentifier)
	} else {
		readTeamResourceModel.TeamAlias = OptionalStringValue(teamIdentifier)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &readTeamResourceModel)...)
}

func (teamTagResource *TeamTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := read[TeamTagResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// use either the team ID or alias based on what is used in the config
	var teamIdentifier string
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

	updatedTeamTagResourceModel := NewTeamTagResourceModel(*teamTag)
	// use either the team ID or alias based on what is used in the config
	if opslevel.IsID(teamIdentifier) {
		updatedTeamTagResourceModel.Team = OptionalStringValue(teamIdentifier)
	} else {
		updatedTeamTagResourceModel.TeamAlias = OptionalStringValue(teamIdentifier)
	}
	tflog.Trace(ctx, "updated a team tag")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedTeamTagResourceModel)...)
}

func (teamTagResource *TeamTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[TeamTagResourceModel](ctx, &resp.Diagnostics, req.State)
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
	if !isTagValid(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<team-id>:<tag-id>'. Given '%s'", req.ID),
		)
		return
	}

	ids := strings.Split(req.ID, ":")
	teamId := ids[0]
	tagId := ids[1]

	team, err := teamTagResource.client.GetTaggableResource(opslevel.TaggableResourceTeam, teamId)
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read team (%s), got error: %s", teamId, err))
		return
	}
	tags, diags := getTagsFromResource(teamTagResource.client, team)
	resp.Diagnostics.Append(diags...)
	if tags == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to get tags from team with id '%s'", teamId))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	teamTag := extractTagFromTags(opslevel.ID(tagId), tags.Nodes)
	if teamTag == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to find tag with id '%s' in team with id '%s'", tagId, teamId))
		return
	}

	idPath := path.Root("id")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idPath, string(teamTag.Id))...)

	keyPath := path.Root("key")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, keyPath, teamTag.Key)...)

	teamPath := path.Root("team")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, teamPath, string(team.ResourceId()))...)

	valuePath := path.Root("value")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, valuePath, teamTag.Value)...)
}
