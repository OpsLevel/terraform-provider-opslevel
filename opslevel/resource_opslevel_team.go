package opslevel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &TeamResource{}

var _ resource.ResourceWithImportState = &TeamResource{}

type TeamResource struct {
	CommonResourceClient
}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

// TeamResourceModel describes the Team managed resource.
type TeamResourceModel struct {
	Aliases          types.List   `tfsdk:"aliases"`
	Id               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	Member           []TeamMember `tfsdk:"member"`
	Name             types.String `tfsdk:"name"`
	Parent           types.String `tfsdk:"parent"`
	Responsibilities types.String `tfsdk:"responsibilities"`
}

type TeamMember struct {
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func convertTeamMember(teamMember opslevel.TeamMembership) TeamMember {
	return TeamMember{
		Email: RequiredStringValue(teamMember.User.Email),
		Role:  RequiredStringValue(teamMember.Role),
	}
}

func NewTeamResourceModel(ctx context.Context, team opslevel.Team, givenModel TeamResourceModel) (TeamResourceModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, team.ManagedAliases)
	if diags != nil && diags.HasError() {
		return TeamResourceModel{}, diags
	}
	teamResourceModel := TeamResourceModel{
		Aliases:          aliases,
		Id:               ComputedStringValue(string(team.Id)),
		Name:             RequiredStringValue(team.Name),
		Responsibilities: OptionalStringValue(team.Responsibilities),
	}
	if ListValueStringsAreEqual(ctx, givenModel.Aliases, teamResourceModel.Aliases) {
		teamResourceModel.Aliases = givenModel.Aliases
	}

	if len(givenModel.Member) > 0 && team.Memberships != nil {
		for _, mem := range team.Memberships.Nodes {
			teamResourceModel.Member = append(teamResourceModel.Member, convertTeamMember(mem))
		}
	}

	return teamResourceModel, diags
}

func (teamResource *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (teamResource *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Team Resource",
		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of human-friendly, unique identifiers for the team.",
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listSortStringModifier{},
				},
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The team's display name.",
				Required:    true,
			},
			"parent": schema.StringAttribute{
				Description: "The id or alias of the parent team.",
				Optional:    true,
			},
			"responsibilities": schema.StringAttribute{
				Description: "A description of what the team is responsible for.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"member": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Description: "The email address of the team member. Must be sorted by email address.",
							Required:    true,
						},
						"role": schema.StringAttribute{
							Description: "The role of the team member.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (teamResource *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	teamCreateInput := opslevel.TeamCreateInput{
		Name:             planModel.Name.ValueString(),
		Responsibilities: planModel.Responsibilities.ValueStringPointer(),
	}

	members, err := getMembers(planModel.Member)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("unable to read members, got error: %s", err))
		return
	}
	if len(members) > 0 {
		teamCreateInput.Members = &members
	}

	if planModel.Parent.ValueString() != "" {
		teamCreateInput.ParentTeam = opslevel.NewIdentifier(planModel.Parent.ValueString())
	}
	team, err := teamResource.client.CreateTeam(teamCreateInput)
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team, got error: %s", err))
		return
	}
	err = team.Hydrate(teamResource.client)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to hydrate team, got error: %s", err))
		return
	}

	aliases, diags := ListValueToStringSlice(ctx, planModel.Aliases)
	resp.Diagnostics.Append(diags...)
	if err = teamResource.reconcileTeamAliases(team, aliases); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to reconcile aliases, got error: %s", err))
		return
	}

	createdTeamResourceModel, diags := NewTeamResourceModel(ctx, *team, planModel)
	resp.Diagnostics.Append(diags...)

	// if parent is set, use an ID or alias for this field based on what is currently in the state
	if opslevel.IsID(planModel.Parent.ValueString()) {
		createdTeamResourceModel.Parent = types.StringValue(string(team.ParentTeam.Id))
	} else {
		// TODO: error thrown if config has alias from the parent team that is not the default alias
		createdTeamResourceModel.Parent = OptionalStringValue(team.ParentTeam.Alias)
	}

	createdTeamResourceModel.LastUpdated = timeLastUpdated()
	tflog.Trace(ctx, "created a team resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdTeamResourceModel)...)
}

func (teamResource *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := teamResource.client.GetTeam(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read team, got error: %s", err))
		return
	}
	err = team.Hydrate(teamResource.client)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to hydrate team, got error: %s", err))
		return
	}

	readTeamResourceModel, diags := NewTeamResourceModel(ctx, *team, stateModel)
	resp.Diagnostics.Append(diags...)
	// if parent is set, use an ID or alias for this field based on what is currently in the state
	if opslevel.IsID(stateModel.Parent.ValueString()) {
		readTeamResourceModel.Parent = types.StringValue(string(team.ParentTeam.Id))
	} else {
		// TODO: error thrown if config has alias from the parent team that is not the default alias
		readTeamResourceModel.Parent = OptionalStringValue(team.ParentTeam.Alias)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &readTeamResourceModel)...)
}

func (teamResource *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel, stateModel TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read state to help determine if Team Members should be deleted
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	teamUpdateInput := opslevel.TeamUpdateInput{
		Id:               opslevel.NewID(planModel.Id.ValueString()),
		Name:             planModel.Name.ValueStringPointer(),
		Responsibilities: opslevel.RefOf(planModel.Responsibilities.ValueString()),
	}

	// Delete Team Members only if we were tracking them and they have been removed
	if len(stateModel.Member) > 0 && len(planModel.Member) == 0 {
		teamUpdateInput.Members = &[]opslevel.TeamMembershipUserInput{}
	}

	members, err := getMembers(planModel.Member)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("unable to read members, got error: %s", err))
		return
	}
	if len(members) > 0 {
		teamUpdateInput.Members = &members
	}

	if planModel.Parent.ValueString() != "" {
		teamUpdateInput.ParentTeam = opslevel.NewIdentifier(planModel.Parent.ValueString())
	} else {
		teamUpdateInput.ParentTeam = opslevel.NewIdentifier()
	}
	updatedTeam, err := teamResource.client.UpdateTeam(teamUpdateInput)
	if err != nil || updatedTeam == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team, got error: %s", err))
		return
	}
	err = updatedTeam.Hydrate(teamResource.client)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to hydrate team, got error: %s", err))
		return
	}

	aliases, diags := ListValueToStringSlice(ctx, planModel.Aliases)
	resp.Diagnostics.Append(diags...)
	if err = teamResource.reconcileTeamAliases(updatedTeam, aliases); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to reconcile aliases, got error: %s", err))
		return
	}

	updatedTeamResourceModel, diags := NewTeamResourceModel(ctx, *updatedTeam, planModel)
	resp.Diagnostics.Append(diags...)
	// if parent is set, use an ID or alias for this field based on what is currently in the state
	if opslevel.IsID(planModel.Parent.ValueString()) {
		updatedTeamResourceModel.Parent = types.StringValue(string(updatedTeam.ParentTeam.Id))
	} else {
		// TODO: error thrown if config has alias from the parent team that is not the default alias
		updatedTeamResourceModel.Parent = OptionalStringValue(updatedTeam.ParentTeam.Alias)
	}
	updatedTeamResourceModel.LastUpdated = timeLastUpdated()
	tflog.Trace(ctx, "updated a team resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedTeamResourceModel)...)
}

func (teamResource *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := teamResource.client.DeleteTeam(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to delete team, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a team resource")
}

func (teamResource *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getMembers(members []TeamMember) ([]opslevel.TeamMembershipUserInput, error) {
	memberInputs := make([]opslevel.TeamMembershipUserInput, len(members))
	for i, mem := range members {
		memberInputs[i] = opslevel.TeamMembershipUserInput{
			User: opslevel.NewUserIdentifier(mem.Email.ValueString()),
			Role: mem.Role.ValueStringPointer(),
		}
	}
	if len(memberInputs) > 0 {
		return memberInputs, nil
	}
	return nil, nil
}

func (teamResource *TeamResource) reconcileTeamAliases(team *opslevel.Team, expectedAliases []string) error {
	// get list of existing aliases from OpsLevel
	existingAliases := team.ManagedAliases

	// if an existing alias is not supposed to be there, delete it
	for _, existingAlias := range existingAliases {
		if !slices.Contains(expectedAliases, existingAlias) {
			err := teamResource.client.DeleteTeamAlias(existingAlias)
			if err != nil {
				return err
			}
		}
	}
	// if an alias does not exist but is supposed to, create it
	for _, expectedAlias := range expectedAliases {
		if !slices.Contains(existingAliases, expectedAlias) {
			_, err := teamResource.client.CreateAliases(team.Id, []string{expectedAlias})
			if err != nil {
				return err
			}
		}
	}
	team.ManagedAliases = expectedAliases
	return nil
}
