package opslevel

import (
	"context"
	"fmt"

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

var _ resource.ResourceWithConfigure = &SystemResource{}

var _ resource.ResourceWithImportState = &SystemResource{}

func NewSystemResource() resource.Resource {
	return &SystemResource{}
}

// SystemResource defines the resource implementation.
type SystemResource struct {
	CommonResourceClient
}

// SystemResourceModel describes the System managed resource.
type SystemResourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Note        types.String `tfsdk:"note"`
	Owner       types.String `tfsdk:"owner"`
}

func NewSystemResourceModel(system opslevel.System, givenModel SystemResourceModel) SystemResourceModel {
	aliases := OptionalStringListValue(system.Aliases)
	systemDataSourceModel := SystemResourceModel{
		Aliases:     aliases,
		Description: StringValueFromResourceAndModelField(system.Description, givenModel.Description),
		Domain:      OptionalStringValue(string(system.Parent.Id)),
		Id:          ComputedStringValue(string(system.Id)),
		Name:        RequiredStringValue(system.Name),
		Note:        StringValueFromResourceAndModelField(system.Note, givenModel.Note),
		Owner:       OptionalStringValue(string(system.Owner.Id())),
	}
	return systemDataSourceModel
}

func (r *SystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (r *SystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "System Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases of the system.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description for the system.",
				Optional:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The id of the parent domain this system is a child for.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the system.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name for the system.",
				Required:    true,
			},
			"note": schema.StringAttribute{
				Description: "Additional information about the system.",
				Optional:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The id or alias of the team that owns the system.",
				Optional:    true,
			},
		},
	}
}

func (r *SystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[SystemResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.SystemInput{
		Name:        nullable(planModel.Name.ValueStringPointer()),
		Description: nullable(planModel.Description.ValueStringPointer()),
		Note:        nullable(planModel.Note.ValueStringPointer()),
	}

	teamIdentifier := planModel.Owner.ValueStringPointer()
	if teamIdentifier != nil && !opslevel.IsID(*teamIdentifier) {
		team, err := r.client.GetTeamWithAlias(*teamIdentifier)
		if err != nil {
			resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read team, got error: %s", err))
			return
		}
		*teamIdentifier = string(team.Id)
	}
	input.OwnerId = nullableID(teamIdentifier)

	if planModel.Domain.ValueString() != "" {
		input.Parent = opslevel.NewIdentifier(planModel.Domain.ValueString())
	}
	system, err := r.client.CreateSystem(input)
	if err != nil || system == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create system, got error: %s", err))
		return
	}
	stateModel := NewSystemResourceModel(*system, planModel)

	tflog.Trace(ctx, "created a system resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *SystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[SystemResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	readSystem, err := r.client.GetSystem(stateModel.Id.ValueString())
	if err != nil {
		if (readSystem == nil || readSystem.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read system, got error: %s", err))
		return
	}
	verifiedStateModel := NewSystemResourceModel(*readSystem, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *SystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[SystemResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[SystemResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.SystemInput{
		Name:        opslevel.RefOf(planModel.Name.ValueString()),
		Description: opslevel.RefOf(planModel.Description.ValueString()),
		Note:        opslevel.RefOf(planModel.Note.ValueString()),
	}

	teamIdentifier := planModel.Owner.ValueStringPointer()
	if teamIdentifier == nil {
		if !stateModel.Owner.IsNull() {
			input.OwnerId = nullableID(teamIdentifier)
		}
	} else {
		if !opslevel.IsID(*teamIdentifier) {
			team, err := r.client.GetTeamWithAlias(*teamIdentifier)
			if err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read team, got error: %s", err))
				return
			}
			*teamIdentifier = string(team.Id)
		}
		input.OwnerId = nullableID(teamIdentifier)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	if planModel.Domain.IsNull() {
		input.Parent = opslevel.NewIdentifier()
	} else {
		input.Parent = opslevel.NewIdentifier(planModel.Domain.ValueString())
	}

	system, err := r.client.UpdateSystem(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update system, got error: %s", err))
		return
	}
	finalModel := NewSystemResourceModel(*system, planModel)

	tflog.Trace(ctx, "updated a system resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalModel)...)
}

func (r *SystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[SystemResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSystem(stateModel.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete system, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a system resource")
}

func (r *SystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
