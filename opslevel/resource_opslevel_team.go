package opslevel

import (
	"context"
	"fmt"
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

func NewTeamResourceModel(ctx context.Context, team opslevel.Team) (TeamResourceModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, team.Aliases)
	if diags != nil && diags.HasError() {
		return TeamResourceModel{}, diags
	}
	teamMembers := make([]TeamMember, 0)
	if team.Memberships != nil {
		for _, mem := range team.Memberships.Nodes {
			teamMembers = append(teamMembers, convertTeamMember(mem))
		}
	}

	return TeamResourceModel{
		Aliases:          aliases,
		Id:               types.StringValue(string(team.Id)),
		Member:           teamMembers,
		Name:             types.StringValue(team.Name),
		Parent:           types.StringValue(team.ParentTeam.Alias),
		Responsibilities: types.StringValue(team.Responsibilities),
	}, diags
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
				Optional: true,
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
							Description: "The email address of the team member.",
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
	var data TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := getMembers(data.Member)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("unable to read members, got error: %s", err))
		return
	}
	team, err := teamResource.client.CreateTeam(opslevel.TeamCreateInput{
		Members:          &members,
		Name:             data.Name.ValueString(),
		ParentTeam:       opslevel.NewIdentifier(data.Parent.ValueString()),
		Responsibilities: data.Responsibilities.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team, got error: %s", err))
		return
	}

	createdTeamResourceModel, diags := NewTeamResourceModel(ctx, *team)
	resp.Diagnostics.Append(diags...)
	createdTeamResourceModel.LastUpdated = timeLastUpdated()
	tflog.Trace(ctx, "created a team resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdTeamResourceModel)...)
}

func (teamResource *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := teamResource.client.GetTeam(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read team, got error: %s", err))
		return
	}

	readTeamResourceModel, diags := NewTeamResourceModel(ctx, *team)
	resp.Diagnostics.Append(diags...)
	readTeamResourceModel.Aliases = data.Aliases
	resp.Diagnostics.Append(resp.State.Set(ctx, &readTeamResourceModel)...)
}

func (teamResource *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := getMembers(data.Member)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("unable to read members, got error: %s", err))
		return
	}
	updatedTeam, err := teamResource.client.UpdateTeam(opslevel.TeamUpdateInput{
		Id:               opslevel.NewID(data.Id.ValueString()),
		Members:          &members,
		Name:             data.Name.ValueStringPointer(),
		ParentTeam:       opslevel.NewIdentifier(data.Parent.ValueString()),
		Responsibilities: data.Responsibilities.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team, got error: %s", err))
		return
	}

	updatedTeamResourceModel, diags := NewTeamResourceModel(ctx, *updatedTeam)
	resp.Diagnostics.Append(diags...)
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
