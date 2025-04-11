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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
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
	AccessControl          types.String    `tfsdk:"access_control"`
	Action                 types.String    `tfsdk:"action"`
	ApprovalRequired       types.Bool      `tfsdk:"approval_required"`
	// ApprovalTeams          types.List      `tfsdk:"approval_teams"`
	ApprovalUsers          types.List      `tfsdk:"approval_users"`
	Description            types.String    `tfsdk:"description"`
	EntityType             types.String    `tfsdk:"entity_type"`
	ExtendedTeamAccess     types.List      `tfsdk:"extended_team_access"`
	Filter                 types.String    `tfsdk:"filter"`
	Id                     types.String    `tfsdk:"id"`
	ManualInputsDefinition types.String    `tfsdk:"manual_inputs_definition"`
	Name                   types.String    `tfsdk:"name"`
	Owner                  types.String    `tfsdk:"owner"`
	ResponseTemplate       types.String    `tfsdk:"response_template"`
	Published              types.Bool      `tfsdk:"published"`
}

// func convertUser(user opslevel.UserIdentifierInput) string {
// 	return OptionalStringValue(user.Email)
// }

func NewTriggerDefinitionResourceModel(client *opslevel.Client, triggerDefinition opslevel.CustomActionsTriggerDefinition, givenModel TriggerDefinitionResourceModel) (TriggerDefinitionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err error

	triggerDefinitionResourceModel := TriggerDefinitionResourceModel{
		AccessControl:          RequiredStringValue(string(triggerDefinition.AccessControl)),
		Action:                 RequiredStringValue(string(triggerDefinition.Action.Id)),
		ApprovalRequired:       types.BoolValue(triggerDefinition.ApprovalConfig.ApprovalRequired),
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

	// if givenModel.ApprovalUsers.IsNull() || givenModel.ApprovalUsers.IsUnknown() {
	// 	triggerDefinitionResourceModel.ApprovalUsers = types.ListNull(types.StringType)
	// } else if len(givenModel.ApprovalUsers.Elements()) == 0 {
	// 	triggerDefinitionResourceModel.ApprovalUsers = types.ListValueMust(types.StringType, []attr.Value{})
	// } else {
	// 	:= OptionalStringListValue(flattenUsersArray(extendedTeams))
	// 	for _, user := range triggerDefinition.ApprovalConfig.Users.Nodes {
	// 		triggerDefinitionResourceModel.ApprovalUsers = append(triggerDefinitionResourceModel.ApprovalUsers, convertUser(user))
	// 	}
	// }

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
			"approval_required": schema.BoolAttribute{
				Description: "Flag indicating approval is required.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			// "approval_teams": schema.ListAttribute{
			// 	ElementType: types.StringType,
			// 	Description: "Teams that can approve this Trigger Definition.",
			// 	Optional:    true,
			// },
			"approval_users": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Users that can approve this Trigger Definition.",
				Optional:    true,
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
	planModel := read[TriggerDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	accessControl := opslevel.CustomActionsTriggerDefinitionAccessControlEnum(planModel.AccessControl.ValueString())
	triggerDefinitionInput := opslevel.CustomActionsTriggerDefinitionCreateInput{
		AccessControl:          &accessControl,
		ActionId:               nullableID(planModel.Action.ValueStringPointer()),       
		Name:                   planModel.Name.ValueString(),
		Description:            nullable(planModel.Description.ValueStringPointer()),
		OwnerId:                opslevel.ID(planModel.Owner.ValueString()),
		ManualInputsDefinition: nullable(planModel.ManualInputsDefinition.ValueStringPointer()),
		Published:              nullable(planModel.Published.ValueBoolPointer()),
		ResponseTemplate:       nullable(planModel.ResponseTemplate.ValueStringPointer()),
	}
	if !planModel.Filter.IsNull() {
		triggerDefinitionInput.FilterId = nullable(opslevel.NewID(planModel.Filter.ValueString()))
	}
	if !planModel.EntityType.IsNull() && !planModel.EntityType.IsUnknown() {
		entityType := opslevel.CustomActionsEntityTypeEnum(planModel.EntityType.ValueString())
		triggerDefinitionInput.EntityType = &entityType
	}

	if planModel.ApprovalRequired == types.BoolValue(true) {
		var approvalConfig opslevel.ApprovalConfigInput
		var err_string string
		approvalConfig, err_string = getApprovalConfig(ctx, planModel)
		if len(err_string) > 0 {
			resp.Diagnostics.AddError("opslevel client error", err_string)
			return
		}
		triggerDefinitionInput.ApprovalConfig = &approvalConfig
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
	stateModel := read[TriggerDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	triggerDefinition, err := r.client.GetTriggerDefinition(stateModel.Id.ValueString())
	if err != nil {
		if (triggerDefinition == nil || triggerDefinition.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read trigger definition, got error: %s", err))
		return
	}

	verifiedStateModel, diags := NewTriggerDefinitionResourceModel(r.client, *triggerDefinition, stateModel)
	resp.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *TriggerDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[TriggerDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[TriggerDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	accessControl := opslevel.CustomActionsTriggerDefinitionAccessControlEnum(planModel.AccessControl.ValueString())
	entityType := opslevel.CustomActionsEntityTypeEnum(planModel.EntityType.ValueString())
	triggerDefinitionInput := opslevel.CustomActionsTriggerDefinitionUpdateInput{
		AccessControl:          &accessControl,
		ActionId:               unsetIDHelper(planModel.Action, stateModel.Action),
		Name:                   unsetStringHelper(planModel.Name, stateModel.Name),
		Description:            opslevel.RefOf(planModel.Description.ValueString()),
		EntityType:             &entityType,
		Id:                     opslevel.ID(planModel.Id.ValueString()),
		OwnerId:                unsetIDHelper(planModel.Owner, stateModel.Owner),
		FilterId:               unsetIDHelper(planModel.Filter, stateModel.Filter),
		ManualInputsDefinition: unsetStringHelper(planModel.ManualInputsDefinition, stateModel.ManualInputsDefinition),
		Published:              opslevel.RefOf(planModel.Published.ValueBool()),
		ResponseTemplate:       unsetStringHelper(planModel.ResponseTemplate, stateModel.ResponseTemplate),
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

	stateModelFinal, diags := NewTriggerDefinitionResourceModel(r.client, *updatedTriggerDefinition, planModel)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "updated a trigger definition resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModelFinal)...)
}

func (r *TriggerDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[TriggerDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTriggerDefinition(stateModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete trigger definition, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a trigger definition resource")
}

func (r *TriggerDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getApprovalConfig(ctx context.Context, planModel TriggerDefinitionResourceModel) (opslevel.ApprovalConfigInput, string) {
	approvalConfig := opslevel.ApprovalConfigInput{
		ApprovalRequired: planModel.ApprovalRequired.ValueBoolPointer(),
	}
	usersStringSlice, diags := ListValueToStringSlice(ctx, planModel.ApprovalUsers)
	if diags.HasError() {
		return approvalConfig, "failed to convert 'approval_users' to string slice"
	}
	users, user_err := getUsers(usersStringSlice)
	if user_err != nil {
		return approvalConfig, fmt.Sprintf("unable to read members, got error: %s", user_err)
	}
	if len(users) > 0 {
		approvalConfig.Users = &users
	}
	return approvalConfig, ""
}

func getUsers(users []string) ([]opslevel.UserIdentifierInput, error) {
	userInputs := make([]opslevel.UserIdentifierInput, len(users))
	for i, user := range users  {
		userInputs[i] = *opslevel.NewUserIdentifier(user)
	}
	if len(userInputs) > 0 {
		return userInputs, nil
	}
	return nil, nil
}

func getExtendedTeamAccessListValue(client *opslevel.Client, triggerDefinition *opslevel.CustomActionsTriggerDefinition) (basetypes.ListValue, error) {
	extendedTeams, err := triggerDefinition.ExtendedTeamAccess(client, nil)
	if err != nil {
		return types.ListNull(types.StringType), err
	}
	extendedTeamsAccess := OptionalStringListValue(flattenTeamsArray(extendedTeams))
	return extendedTeamsAccess, nil
}
