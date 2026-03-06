package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
)

var (
	_ resource.ResourceWithConfigure   = &TeamPropertyDefinitionResource{}
	_ resource.ResourceWithImportState = &TeamPropertyDefinitionResource{}
)

type TeamPropertyDefinitionResource struct {
	CommonResourceClient
}

func NewTeamPropertyDefinitionResource() resource.Resource {
	return &TeamPropertyDefinitionResource{}
}

type TeamPropertyDefinitionResourceModel struct {
	Alias        types.String `tfsdk:"alias"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	LockedStatus types.String `tfsdk:"locked_status"`
	Name         types.String `tfsdk:"name"`
	Schema       types.String `tfsdk:"schema"`
}

func NewTeamPropertyDefinitionResourceModel(definition opslevel.TeamPropertyDefinition, givenModel TeamPropertyDefinitionResourceModel) TeamPropertyDefinitionResourceModel {
	return TeamPropertyDefinitionResourceModel{
		Alias:        RequiredStringValue(definition.Alias),
		Description:  StringValueFromResourceAndModelField(definition.Description, givenModel.Description),
		Id:           ComputedStringValue(string(definition.Id)),
		LockedStatus: RequiredStringValue(string(definition.LockedStatus)),
		Name:         RequiredStringValue(definition.Name),
		Schema:       RequiredStringValue(definition.Schema.AsString()),
	}
}

func (r *TeamPropertyDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_property_definition"
}

func (r *TeamPropertyDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Team Property Definition Resource",
		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "The human-friendly, unique identifier of the team property definition.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the team property definition.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"locked_status": schema.StringAttribute{
				Description: fmt.Sprintf(
					"Restricts what sources are able to assign values to this property. One of `%s`",
					strings.Join(opslevel.AllPropertyLockedStatusEnum, "`, `"),
				),
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllPropertyLockedStatusEnum...),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the team property definition.",
				Required:    true,
			},
			"schema": schema.StringAttribute{
				Description: "The schema of the team property definition.",
				Required:    true,
				Validators: []validator.String{
					JsonObjectValidator(),
				},
			},
		},
	}
}

func (r *TeamPropertyDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[TeamPropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	definitionSchema, err := opslevel.NewJSONSchema(planModel.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to use definition schema '%s', got error: %s", planModel.Schema.ValueString(), err))
		return
	}

	input := opslevel.TeamPropertyDefinitionInput{
		Alias:        planModel.Alias.ValueString(),
		Description:  planModel.Description.ValueString(),
		LockedStatus: asEnum[opslevel.PropertyLockedStatusEnum](planModel.LockedStatus.ValueStringPointer()),
		Name:         planModel.Name.ValueString(),
		Schema:       *definitionSchema,
	}
	definition, err := r.client.CreateTeamPropertyDefinition(input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create team property definition '%s', got error: %s", input.Name, err))
		return
	}

	stateModel := NewTeamPropertyDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("created a team property definition resource with id '%s'", definition.Id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TeamPropertyDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[TeamPropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	definition, err := r.client.GetTeamPropertyDefinition(id)
	if err != nil {
		if (definition == nil || definition.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read team property definition with id '%s', got error: %s", id, err))
		return
	}

	verifiedStateModel := NewTeamPropertyDefinitionResourceModel(*definition, stateModel)
	tflog.Trace(ctx, fmt.Sprintf("read a team property definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *TeamPropertyDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[TeamPropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	definitionSchema, err := opslevel.NewJSONSchema(planModel.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to use definition schema '%s', got error: %s", planModel.Schema.ValueString(), err))
		return
	}

	id := planModel.Id.ValueString()
	input := opslevel.TeamPropertyDefinitionInput{
		Alias:        planModel.Alias.ValueString(),
		Description:  planModel.Description.ValueString(),
		LockedStatus: asEnum[opslevel.PropertyLockedStatusEnum](planModel.LockedStatus.ValueStringPointer()),
		Name:         planModel.Name.ValueString(),
		Schema:       *definitionSchema,
	}
	definition, err := r.client.UpdateTeamPropertyDefinition(id, input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update team property definition with id '%s', got error: %s", id, err))
		return
	}

	stateModel := NewTeamPropertyDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("updated a team property definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TeamPropertyDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[TeamPropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	// No individual delete mutation exists. List all, filter out this one, bulk reassign.
	defs, err := r.client.ListTeamPropertyDefinitions(nil)
	if err != nil || defs == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to list team property definitions for delete of '%s', got error: %s", stateModel.Id.ValueString(), err))
		return
	}

	remaining := make([]opslevel.TeamPropertyDefinitionInput, 0, len(defs.Nodes))
	for _, d := range defs.Nodes {
		if string(d.Id) == stateModel.Id.ValueString() {
			continue
		}
		lockedStatus := d.LockedStatus
		remaining = append(remaining, opslevel.TeamPropertyDefinitionInput{
			Alias:        d.Alias,
			Description:  d.Description,
			LockedStatus: &lockedStatus,
			Name:         d.Name,
			Schema:       d.Schema,
		})
	}

	_, err = r.client.AssignTeamPropertyDefinitions(opslevel.TeamPropertyDefinitionsAssignInput{
		Properties: remaining,
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete team property definition '%s', got error: %s", stateModel.Id.ValueString(), err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("deleted team property definition with id '%s'", stateModel.Id.ValueString()))
}

func (r *TeamPropertyDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
