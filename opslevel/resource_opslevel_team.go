package opslevel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// teamResourceModel describes the Team managed resource.
type teamResourceModel struct {
	Aliases          types.List        `tfsdk:"aliases"`
	Id               types.String      `tfsdk:"id"`
	LastUpdated      types.String      `tfsdk:"last_updated"`
	Member           []teamMemberModel `tfsdk:"member"`
	Name             types.String      `tfsdk:"name"`
	Parent           types.String      `tfsdk:"parent"`
	Responsibilities types.String      `tfsdk:"responsibilities"`
}

func newTeamResourceModel(ctx context.Context, team opslevel.Team, parentIdentifier string) (teamResourceModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, team.ManagedAliases)
	if diags != nil && diags.HasError() {
		return teamResourceModel{}, diags
	}
	teamMembers := make([]teamMemberModel, 0)
	if team.Memberships != nil {
		for _, mem := range team.Memberships.Nodes {
			teamMembers = append(teamMembers, newTeamMemberModel(mem))
		}
	}
	model := teamResourceModel{
		Aliases:          aliases,
		Id:               ComputedStringValue(string(team.Id)),
		Member:           teamMembers,
		Name:             RequiredStringValue(team.Name),
		Responsibilities: OptionalStringValue(team.Responsibilities),
	}
	// set parent team
	if team.ParentTeam.Id == "" {
		model.Parent = types.StringNull()
	} else if opslevel.IsID(parentIdentifier) {
		model.Parent = types.StringValue(string(team.ParentTeam.Id))
	} else {
		// can use non-default aliases
		model.Parent = types.StringValue(parentIdentifier)
	}

	return model, diags
}

// removeNonTerraformManagedMembers mutates the input model to exclude team members not managed by terraform
func removeNonTerraformManagedMembers(ctx context.Context, model *teamResourceModel, cachedModel teamResourceModel) {
	tfManagedMembers := make([]teamMemberModel, 0)
	for _, memberModel := range model.Member {
		if slices.Contains(cachedModel.Member, memberModel) {
			tfManagedMembers = append(tfManagedMembers, memberModel)
		} else {
			tflog.Debug(ctx, fmt.Sprintf("not a terraform managed team member: '%s'", memberModel.Email.ValueString()))
		}
	}
	model.Member = tfManagedMembers
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
			"member": schema.SetNestedAttribute{
				Description: "Unordered list of members. Only manages team members that were defined in terraform.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Description: "The email address of the team member. Must be sorted by email address.",
							Required:    true,
						},
						"role": schema.StringAttribute{
							Description: "The role of the team member. Options: `contributor` or `manager`",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (teamResource *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel teamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := getMembers(planModel.Member)
	teamCreateInput := opslevel.TeamCreateInput{
		Members:          &members,
		Name:             planModel.Name.ValueString(),
		ParentTeam:       opslevel.NewIdentifier(),
		Responsibilities: planModel.Responsibilities.ValueStringPointer(),
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
	if err = teamResource.reconcileTeamAliases(team, planModel); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to reconcile aliases, got error: %s", err))
		return
	}

	newStateModel, diags := newTeamResourceModel(ctx, *team, planModel.Parent.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	removeNonTerraformManagedMembers(ctx, &newStateModel, planModel)
	newStateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a team resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (teamResource *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel teamResourceModel
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

	newStateModel, diags := newTeamResourceModel(ctx, *team, stateModel.Parent.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	removeNonTerraformManagedMembers(ctx, &newStateModel, stateModel)

	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (teamResource *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel teamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := getMembers(planModel.Member)
	_, err := teamResource.client.AddMemberships(&opslevel.TeamId{Id: opslevel.ID(planModel.Id.ValueString())}, members...)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to add members, got error: %s", err))
		return
	}
	teamUpdateInput := opslevel.TeamUpdateInput{
		Id:               opslevel.NewID(planModel.Id.ValueString()),
		Name:             planModel.Name.ValueStringPointer(),
		ParentTeam:       opslevel.NewIdentifier(),
		Responsibilities: planModel.Responsibilities.ValueStringPointer(),
	}
	if planModel.Parent.ValueString() != "" {
		teamUpdateInput.ParentTeam = opslevel.NewIdentifier(planModel.Parent.ValueString())
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
	if err = teamResource.reconcileTeamAliases(updatedTeam, planModel); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to reconcile aliases, got error: %s", err))
		return
	}

	newStateModel, diags := newTeamResourceModel(ctx, *updatedTeam, planModel.Parent.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	removeNonTerraformManagedMembers(ctx, &newStateModel, planModel)
	newStateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a team resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (teamResource *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data teamResourceModel
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

// TODO: works fine for when we need to change a member from a contributor to manager and vv but does not work for when we need to unset members
// TODO: can go 2 ways - use remove memberships call or GET team first, put everything in here, remove the member(s) removed from teh plan.
func getMembers(tfManagedMembers []teamMemberModel) []opslevel.TeamMembershipUserInput {
	memberInputs := make([]opslevel.TeamMembershipUserInput, len(tfManagedMembers))
	for i, mem := range tfManagedMembers {
		memberInputs[i] = opslevel.TeamMembershipUserInput{
			User: opslevel.NewUserIdentifier(mem.Email.ValueString()),
			Role: mem.Role.ValueStringPointer(),
		}
	}
	return memberInputs
}

func (teamResource *TeamResource) reconcileTeamAliases(team *opslevel.Team, data teamResourceModel) error {
	// get list of expected aliases from terraform
	tmp := data.Aliases.Elements()
	expectedAliases := make([]string, len(tmp))
	for i, alias := range tmp {
		expectedAliases[i] = unquote(alias.String())
	}
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
