package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &TriggerDefinitionResource{}

var _ resource.ResourceWithImportState = &TriggerDefinitionResource{}

func NewTriggerDefinitionResource() resource.Resource {
	return &TriggerDefinitionResource{}
}

// TriggerDefinitionResource defines the resource implementation.
type TriggerDefinitionResource struct {
	CommonResourceClient
}

// TriggerDefinitionResourceModel describes the trigger definition managed resource.
type TriggerDefinitionResourceModel struct {
	AccessControl          types.String `tfsdk:"access_control"`
	Action                 types.String `tfsdk:"action"`
	Description            types.String `tfsdk:"description"`
	EntityType             types.String `tfsdk:"entity_type"`
	ExtendedTeamAccess     types.List   `tfsdk:"extended_team_access"`
	Filter                 types.String `tfsdk:"filter"`
	Id                     types.String `tfsdk:"id"`
	ManualInputsDefinition types.String `tfsdk:"manual_inputs_definition"`
	Name                   types.String `tfsdk:"name"`
	Owner                  types.String `tfsdk:"owner"`
	ResponseTemplate       types.String `tfsdk:"response_template"`
	Published              types.Bool   `tfsdk:"published"`
}

func NewTriggerDefinitionResourceModel(client *opslevel.Client, triggerDefinition opslevel.CustomActionsTriggerDefinition, givenModel TriggerDefinitionResourceModel) (TriggerDefinitionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err error

	triggerDefinitionResourceModel := TriggerDefinitionResourceModel{
		AccessControl:          RequiredStringValue(string(triggerDefinition.AccessControl)),
		Action:                 RequiredStringValue(string(triggerDefinition.Action.Id)),
		Description:            OptionalStringValue(triggerDefinition.Description),
		EntityType:             OptionalStringValue(string(triggerDefinition.EntityType)),
		Filter:                 OptionalStringValue(string(triggerDefinition.Filter.Id)),
		Id:                     ComputedStringValue(string(triggerDefinition.Id)),
		ManualInputsDefinition: OptionalStringValue(triggerDefinition.ManualInputsDefinition),
		Name:                   RequiredStringValue(triggerDefinition.Name),
		Owner:                  RequiredStringValue(string(triggerDefinition.Owner.Id)),
		ResponseTemplate:       OptionalStringValue(triggerDefinition.ResponseTemplate),
		Published:              types.BoolValue(triggerDefinition.Published),
	}

	if givenModel.ExtendedTeamAccess.IsNull() || givenModel.ExtendedTeamAccess.IsUnknown() {
		triggerDefinitionResourceModel.ExtendedTeamAccess = types.ListNull(types.StringType)
	} else if len(givenModel.ExtendedTeamAccess.Elements()) == 0 {
		triggerDefinitionResourceModel.ExtendedTeamAccess = types.ListValueMust(types.StringType, []attr.Value{})
	} else {
		triggerDefinitionResourceModel.ExtendedTeamAccess, err = getExtendedTeamAccessListValue(client, &triggerDefinition)
		if err != nil {
			diags.AddError("opslevel client error", fmt.Sprintf("Unable to get teams for 'extended_team_access', got error: %s", err))
		}
	}
	return triggerDefinitionResourceModel, diags
}

func (r *TriggerDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger_definition"
}

func (r *TriggerDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Trigger Definition Resource",

		Attributes: map[string]schema.Attribute{
			"access_control": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The set of users that should be able to use the Trigger Definition. One of `%s`",
					strings.Join(opslevel.AllCustomActionsTriggerDefinitionAccessControlEnum, "`, `"),
				),
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllCustomActionsTriggerDefinitionAccessControlEnum...),
				},
			},
			"action": schema.StringAttribute{
				Description: "The action that will be triggered by the Trigger Definition.",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"description": schema.StringAttribute{
				Description: "The description of what the Trigger Definition will do.",
				Optional:    true,
			},
			"entity_type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The entity type to associate with the Trigger Definition. One of `%s`",
					strings.Join(opslevel.AllCustomActionsEntityTypeEnum, "`, `"),
				),
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("SERVICE"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllCustomActionsEntityTypeEnum...),
				},
			},
			"extended_team_access": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The set of additional teams who can invoke this Trigger Definition.",
				Optional:    true,
			},
			"filter": schema.StringAttribute{
				Description: "A filter defining which services this Trigger Definition applies to.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the trigger definition.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"manual_inputs_definition": schema.StringAttribute{
				Description: "The YAML definition of any custom inputs for this Trigger Definition.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Trigger Definition",
				Required:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The owner of the Trigger Definition.",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"response_template": schema.StringAttribute{
				Description: "The liquid template used to parse the response from the Webhook Action.",
				Optional:    true,
			},
			"published": schema.BoolAttribute{
				Description: "The published state of the Custom Action; true if the Trigger Definition is ready for use; false if it is a draft.",
				Required:    true,
			},
		},
	}
}

func (r *TriggerDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel TriggerDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessControl := opslevel.CustomActionsTriggerDefinitionAccessControlEnum(planModel.AccessControl.ValueString())
	triggerDefinitionInput := opslevel.CustomActionsTriggerDefinitionCreateInput{
		AccessControl:          &accessControl,
		ActionId:               opslevel.NewID(planModel.Action.ValueString()),
		Name:                   planModel.Name.ValueString(),
		Description:            planModel.Description.ValueStringPointer(),
		OwnerId:                opslevel.ID(planModel.Owner.ValueString()),
		FilterId:               opslevel.NewID(planModel.Filter.ValueString()),
		ManualInputsDefinition: planModel.ManualInputsDefinition.ValueStringPointer(),
		Published:              planModel.Published.ValueBoolPointer(),
		ResponseTemplate:       planModel.ResponseTemplate.ValueStringPointer(),
	}
	if !planModel.EntityType.IsNull() && !planModel.EntityType.IsUnknown() {
		entityType := opslevel.CustomActionsEntityTypeEnum(planModel.EntityType.ValueString())
		triggerDefinitionInput.EntityType = &entityType
	}

	extendedTeamsStringSlice, diags := ListValueToStringSlice(ctx, planModel.ExtendedTeamAccess)
	if diags.HasError() {
		resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extended_team_access' to string slice")
		return
	}
	extendedTeamAccess := opslevel.NewIdentifierArray(extendedTeamsStringSlice)
	triggerDefinitionInput.ExtendedTeamAccess = &extendedTeamAccess

	triggerDefinition, err := r.client.CreateTriggerDefinition(triggerDefinitionInput)
	if err != nil || triggerDefinition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create trigger definition, got error: %s", err))
		return
	}

	stateModel, diags := NewTriggerDefinitionResourceModel(r.client, *triggerDefinition, planModel)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created a trigger definition resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TriggerDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel TriggerDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	triggerDefinition, err := r.client.GetTriggerDefinition(planModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddWarning("State drift", stateResourceMissingMessage("opslevel_trigger_definition"))
		resp.State.RemoveResource(ctx)
		return
	}

	stateModel, diags := NewTriggerDefinitionResourceModel(r.client, *triggerDefinition, planModel)
	resp.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TriggerDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel TriggerDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessControl := opslevel.CustomActionsTriggerDefinitionAccessControlEnum(planModel.AccessControl.ValueString())
	entityType := opslevel.CustomActionsEntityTypeEnum(planModel.EntityType.ValueString())
	triggerDefinitionInput := opslevel.CustomActionsTriggerDefinitionUpdateInput{
		AccessControl:          &accessControl,
		ActionId:               opslevel.NewID(planModel.Action.ValueString()),
		Name:                   planModel.Name.ValueStringPointer(),
		Description:            opslevel.RefOf(planModel.Description.ValueString()),
		EntityType:             &entityType,
		Id:                     opslevel.ID(planModel.Id.ValueString()),
		OwnerId:                opslevel.NewID(planModel.Owner.ValueString()),
		FilterId:               opslevel.NewID(planModel.Filter.ValueString()),
		ManualInputsDefinition: opslevel.RefOf(planModel.ManualInputsDefinition.ValueString()),
		Published:              opslevel.RefOf(planModel.Published.ValueBool()),
		ResponseTemplate:       opslevel.RefOf(planModel.ResponseTemplate.ValueString()),
	}
	extendedTeamsStringSlice, diags := ListValueToStringSlice(ctx, planModel.ExtendedTeamAccess)
	if diags.HasError() {
		resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extended_team_access' to string slice")
		return
	}
	extendedTeamAccess := opslevel.NewIdentifierArray(extendedTeamsStringSlice)
	triggerDefinitionInput.ExtendedTeamAccess = &extendedTeamAccess

	updatedTriggerDefinition, err := r.client.UpdateTriggerDefinition(triggerDefinitionInput)
	if err != nil || updatedTriggerDefinition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update trigger definition, got error: %s", err))
		return
	}

	stateModel, diags := NewTriggerDefinitionResourceModel(r.client, *updatedTriggerDefinition, planModel)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "updated a trigger definition resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TriggerDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel TriggerDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTriggerDefinition(planModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddWarning("State drift", stateResourceMissingMessage("opslevel_trigger_definition"))
		return
	}
	tflog.Trace(ctx, "deleted a trigger definition resource")
}

func (r *TriggerDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getExtendedTeamAccessListValue(client *opslevel.Client, triggerDefinition *opslevel.CustomActionsTriggerDefinition) (basetypes.ListValue, error) {
	extendedTeams, err := triggerDefinition.ExtendedTeamAccess(client, nil)
	if err != nil {
		return types.ListNull(types.StringType), err
	}
	extendedTeamsAccess := OptionalStringListValue(flattenTeamsArray(extendedTeams))
	return extendedTeamsAccess, nil
}
