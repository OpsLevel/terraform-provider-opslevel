package opslevel

import (
	"context"
	"fmt"
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

var _ resource.ResourceWithConfigure = &TeamContactResource{}

var _ resource.ResourceWithImportState = &TeamContactResource{}

type TeamContactResource struct {
	CommonResourceClient
}

func NewTeamContactResource() resource.Resource {
	return &TeamContactResource{}
}

type TeamContactResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Team  types.String `tfsdk:"team"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

func NewTeamContactResourceModel(teamContact opslevel.Contact) TeamContactResourceModel {
	teamResourceModel := TeamContactResourceModel{
		Id:    RequiredStringValue(string(teamContact.Id)),
		Name:  RequiredStringValue(teamContact.DisplayName),
		Type:  RequiredStringValue(string(teamContact.Type)),
		Value: RequiredStringValue(teamContact.Address),
	}
	return teamResourceModel
}

func (teamContactResource *TeamContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_contact"
}

func (teamContactResource *TeamContactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	enumAllContactTypes := append(opslevel.AllContactType, "any")
	resp.Schema = schema.Schema{
		MarkdownDescription: "Team Contact Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name shown in the UI for the contact.",
				Required:    true,
			},
			"team": schema.StringAttribute{
				Description: "The id or alias of the team the contact belongs to.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The method of contact. One of `%s`",
					strings.Join(enumAllContactTypes, "`, `"),
				),
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(enumAllContactTypes...),
				},
			},
			"value": schema.StringAttribute{
				Description: "The contact value. Examples: support@company.com for type email, https://opslevel.com for type web, #devs for type slack",
				Required:    true,
			},
		},
	}
}

func (teamContactResource *TeamContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamContactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactCreateInput := opslevel.ContactInput{
		Address:     data.Value.ValueString(),
		DisplayName: data.Name.ValueStringPointer(),
		Type:        opslevel.ContactType(data.Type.ValueString()),
	}

	teamIdentifier := data.Team.ValueString()

	contact, err := teamContactResource.client.AddContact(teamIdentifier, contactCreateInput)
	if err != nil || contact == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to add contact on team (%s), got error: %s", teamIdentifier, err))
		return
	}

	createdTeamContactModel := NewTeamContactResourceModel(*contact)
	createdTeamContactModel.Team = RequiredStringValue(teamIdentifier)
	tflog.Trace(ctx, "created a team contact resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdTeamContactModel)...)
}

func (teamContactResource *TeamContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	teamIdentifier := data.Team.ValueString()
	contactID := data.Id.ValueString()

	var team *opslevel.Team
	var err error
	if opslevel.IsID(teamIdentifier) {
		team, err = teamContactResource.client.GetTeam(opslevel.ID(teamIdentifier))
	} else {
		team, err = teamContactResource.client.GetTeamWithAlias(teamIdentifier)
	}
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read team (%s), got error: %s", teamIdentifier, err))
		return
	}
	err = team.Hydrate(teamContactResource.client)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to hydrate team (%s), got error: %s", teamIdentifier, err))
	}

	var teamContact *opslevel.Contact
	for _, readContact := range team.Contacts {
		if string(readContact.Id) == contactID {
			teamContact = &readContact
			break
		}
	}
	if teamContact == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("team contact (with ID '%s') not found on team (%s)", contactID, teamIdentifier))
		return
	}

	readTeamContactResourceModel := NewTeamContactResourceModel(*teamContact)
	readTeamContactResourceModel.Team = RequiredStringValue(teamIdentifier)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readTeamContactResourceModel)...)
}

func (teamContactResource *TeamContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamContactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactCreateInput := opslevel.ContactInput{
		Address:     data.Value.ValueString(),
		DisplayName: opslevel.RefOf(data.Name.ValueString()),
		Type:        opslevel.ContactType(data.Type.ValueString()),
	}

	teamIdentifier := data.Team.ValueString()
	contactID := opslevel.ID(data.Id.ValueString())

	contact, err := teamContactResource.client.UpdateContact(contactID, contactCreateInput)
	if err != nil || contact == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update contact (with ID '%s') on team (%s), got error: %s", contactID, teamIdentifier, err))
		return
	}

	updatedTeamContactResourceModel := NewTeamContactResourceModel(*contact)
	updatedTeamContactResourceModel.Team = RequiredStringValue(teamIdentifier)
	tflog.Trace(ctx, "updated a team contact resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedTeamContactResourceModel)...)
}

func (teamContactResource *TeamContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID := opslevel.ID(data.Id.ValueString())

	err := teamContactResource.client.RemoveContact(contactID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to remove team contact (with id '%s'), got error: %s", contactID, err))
		return
	}
	tflog.Trace(ctx, "deleted a team contact resource")
}

func (teamContactResource *TeamContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !isTagValid(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<team-id>:<contact-id>'. Given '%s'", req.ID),
		)
	}

	ids := strings.Split(req.ID, ":")
	teamId := ids[0]
	contactId := ids[1]

	team, err := teamContactResource.client.GetTeam(opslevel.ID(teamId))
	if err != nil || team == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to get team with id '%s', got error: %s", teamId, err))
		return
	}
	if err = team.Hydrate(teamContactResource.client); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to hydrate team with id '%s', got error: %s", teamId, err))
		return
	}
	teamContact := extractContactFromContacts(opslevel.ID(contactId), team.Contacts)
	if teamContact == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to find contact with id '%s' in team with id '%s'", contactId, teamId))
		return
	}

	idPath := path.Root("id")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idPath, string(teamContact.Id))...)

	keyPath := path.Root("name")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, keyPath, teamContact.DisplayName)...)

	teamPath := path.Root("team")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, teamPath, string(team.Id))...)

	typePath := path.Root("type")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, typePath, teamContact.Type)...)

	valuePath := path.Root("value")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, valuePath, teamContact.Address)...)
}
