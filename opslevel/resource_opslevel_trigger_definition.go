package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

func NewTriggerDefinitionResourceModel(triggerDefinition opslevel.CustomActionsTriggerDefinition, extendedTeams basetypes.ListValue) TriggerDefinitionResourceModel {
	return TriggerDefinitionResourceModel{
		AccessControl:          types.StringValue(string(triggerDefinition.AccessControl)),
		Action:                 types.StringValue(string(triggerDefinition.Action.Id)),
		Description:            types.StringValue(triggerDefinition.Description),
		EntityType:             types.StringValue(string(triggerDefinition.EntityType)),
		ExtendedTeamAccess:     extendedTeams,
		Filter:                 types.StringValue(string(triggerDefinition.Filter.Id)),
		Id:                     types.StringValue(string(triggerDefinition.Id)),
		ManualInputsDefinition: types.StringValue(triggerDefinition.ManualInputsDefinition),
		Name:                   types.StringValue(triggerDefinition.Name),
		Owner:                  types.StringValue(string(triggerDefinition.Owner.Id)),
		ResponseTemplate:       types.StringValue(string(triggerDefinition.ResponseTemplate)),
		Published:              types.BoolValue(triggerDefinition.Published),
	}
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
				Description: "The set of users that should be able to use the Trigger Definition. Requires a value of `everyone`, `admins`, or `service_owners`.",
				Required:    true,
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
				Description: "The entity type to associate with the Trigger Definition.",
				Optional:    true,
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
	entityType := opslevel.CustomActionsEntityTypeEnum(planModel.EntityType.ValueString())
	triggerDefinitionInput := opslevel.CustomActionsTriggerDefinitionCreateInput{
		AccessControl:          &accessControl,
		ActionId:               opslevel.NewID(planModel.Action.ValueString()),
		Name:                   planModel.Name.ValueString(),
		Description:            planModel.Name.ValueStringPointer(),
		OwnerId:                opslevel.ID(planModel.Owner.ValueString()),
		FilterId:               opslevel.NewID(planModel.Filter.ValueString()),
		ManualInputsDefinition: planModel.ManualInputsDefinition.ValueStringPointer(),
		Published:              planModel.Published.ValueBoolPointer(),
		ResponseTemplate:       planModel.ResponseTemplate.ValueStringPointer(),
		EntityType:             &entityType,
	}
	extendedTeamsStringSlice, diags := ListValueToStringSlice(ctx, planModel.ExtendedTeamAccess)
	if diags.HasError() {
		resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extended_team_access' to string slice")
		return
	}
	if len(extendedTeamsStringSlice) > 0 {
		extendedTeamAccess := opslevel.NewIdentifierArray(extendedTeamsStringSlice)
		triggerDefinitionInput.ExtendedTeamAccess = &extendedTeamAccess
	}

	triggerDefinition, err := r.client.CreateTriggerDefinition(triggerDefinitionInput)
	if err != nil || triggerDefinition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create trigger definition, got error: %s", err))
		return
	}

	extendedTeams, err := triggerDefinition.ExtendedTeamAccess(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get teams for 'extended_team_access', got error: %s", err))
		return
	}
	extendedTeamsAccess := types.ListNull(types.StringType)
	if extendedTeams != nil {
		extendedTeamsAccess, diags = types.ListValueFrom(ctx, types.StringType, flattenTeamsArray(extendedTeams))
		if diags.HasError() {
			resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extendedTeams' to 'basetypes.ListValue'")
			return
		}
	}
	stateModel := NewTriggerDefinitionResourceModel(*triggerDefinition, extendedTeamsAccess)

	tflog.Trace(ctx, "created a trigger definition resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TriggerDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel TriggerDefinitionResourceModel
	var diags diag.Diagnostics

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	triggerDefinition, err := r.client.GetTriggerDefinition(planModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read trigger definition, got error: %s", err))
		return
	}
	extendedTeams, err := triggerDefinition.ExtendedTeamAccess(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get teams for 'extended_team_access', got error: %s", err))
		return
	}
	extendedTeamsAccess := types.ListNull(types.StringType)
	if extendedTeams != nil {
		extendedTeamsAccess, diags = types.ListValueFrom(ctx, types.StringType, flattenTeamsArray(extendedTeams))
		if diags.HasError() {
			resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extendedTeams' to 'basetypes.ListValue'")
			return
		}
	}
	stateModel := NewTriggerDefinitionResourceModel(*triggerDefinition, extendedTeamsAccess)

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
		Description:            planModel.Name.ValueStringPointer(),
		OwnerId:                opslevel.NewID(planModel.Owner.ValueString()),
		FilterId:               opslevel.NewID(planModel.Filter.ValueString()),
		ManualInputsDefinition: planModel.ManualInputsDefinition.ValueStringPointer(),
		Published:              planModel.Published.ValueBoolPointer(),
		ResponseTemplate:       planModel.ResponseTemplate.ValueStringPointer(),
		EntityType:             &entityType,
	}
	extendedTeamsStringSlice, diags := ListValueToStringSlice(ctx, planModel.ExtendedTeamAccess)
	if diags.HasError() {
		resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extended_team_access' to string slice")
		return
	}
	if len(extendedTeamsStringSlice) > 0 {
		extendedTeamAccess := opslevel.NewIdentifierArray(extendedTeamsStringSlice)
		triggerDefinitionInput.ExtendedTeamAccess = &extendedTeamAccess
	}

	updatedTriggerDefinition, err := r.client.UpdateTriggerDefinition(triggerDefinitionInput)
	if err != nil || updatedTriggerDefinition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update trigger definition, got error: %s", err))
		return
	}

	extendedTeams, err := updatedTriggerDefinition.ExtendedTeamAccess(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get teams for 'extended_team_access', got error: %s", err))
		return
	}
	extendedTeamsAccess := types.ListNull(types.StringType)
	if extendedTeams != nil {
		extendedTeamsAccess, diags = types.ListValueFrom(ctx, types.StringType, flattenTeamsArray(extendedTeams))
		if diags.HasError() {
			resp.Diagnostics.AddError("opslevel client error", "failed to convert 'extendedTeams' to 'basetypes.ListValue'")
			return
		}
	}
	stateModel := NewTriggerDefinitionResourceModel(*updatedTriggerDefinition, extendedTeamsAccess)

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete trigger definition, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a trigger definition resource")
}

func (r *TriggerDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceTriggerDefinition() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a webhook action",
// 		Create:      wrap(resourceTriggerDefinitionCreate),
// 		Read:        wrap(resourceTriggerDefinitionRead),
// 		Update:      wrap(resourceTriggerDefinitionUpdate),
// 		Delete:      wrap(resourceTriggerDefinitionDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The name of the Trigger Definition",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"description": {
// 				Type:        schema.TypeString,
// 				Description: "The description of what the Trigger Definition will do.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"owner": {
// 				Type:        schema.TypeString,
// 				Description: "The owner of the Trigger Definition.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"action": {
// 				Type:        schema.TypeString,
// 				Description: "The action that will be triggered by the Trigger Definition.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"filter": {
// 				Type:        schema.TypeString,
// 				Description: "A filter defining which services this Trigger Definition applies to.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"manual_inputs_definition": {
// 				Type:        schema.TypeString,
// 				Description: "The YAML definition of any custom inputs for this Trigger Definition.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"published": {
// 				Type:        schema.TypeBool,
// 				Description: "The published state of the Custom Action; true if the Trigger Definition is ready for use; false if it is a draft.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"access_control": {
// 				Type:         schema.TypeString,
// 				Description:  "The set of users that should be able to use the Trigger Definition. Requires a value of `everyone`, `admins`, or `service_owners`.",
// 				ForceNew:     false,
// 				Required:     true,
// 				ValidateFunc: validation.StringInSlice(opslevel.AllCustomActionsTriggerDefinitionAccessControlEnum, false),
// 			},
// 			"response_template": {
// 				Type:        schema.TypeString,
// 				Description: "The liquid template used to parse the response from the Webhook Action.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"entity_type": {
// 				Type:         schema.TypeString,
// 				Description:  "The entity type to associate with the Trigger Definition.",
// 				ForceNew:     false,
// 				Optional:     true,
// 				ValidateFunc: validation.StringInSlice(opslevel.AllCustomActionsEntityTypeEnum, false),
// 			},
// 			"extended_team_access": {
// 				Type:        schema.TypeList,
// 				Description: "The set of additional teams who can invoke this Trigger Definition.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem:        &schema.Schema{Type: schema.TypeString},
// 			},
// 		},
// 	}
// }

// func resourceTriggerDefinitionCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CustomActionsTriggerDefinitionCreateInput{
// 		Name:          d.Get("name").(string),
// 		OwnerId:       *opslevel.NewID(d.Get("owner").(string)),
// 		ActionId:      opslevel.NewID(d.Get("action").(string)),
// 		AccessControl: opslevel.RefOf(opslevel.CustomActionsTriggerDefinitionAccessControlEnum(d.Get("access_control").(string))),
// 	}
// 	extended_teams := opslevel.NewIdentifierArray(getStringArray(d, "extended_team_access"))
// 	if len(extended_teams) > 0 {
// 		input.ExtendedTeamAccess = &extended_teams
// 	}

// 	if _, ok := d.GetOk("description"); ok {
// 		input.Description = opslevel.RefOf(d.Get("description").(string))
// 	}

// 	if _, ok := d.GetOk("filter"); ok {
// 		input.FilterId = opslevel.NewID(d.Get("filter").(string))
// 	}

// 	if _, ok := d.GetOk("manual_inputs_definition"); ok {
// 		manualInputsDefinition := d.Get("manual_inputs_definition").(string)
// 		input.ManualInputsDefinition = opslevel.RefOf(manualInputsDefinition)
// 	}

// 	if _, ok := d.GetOk("response_template"); ok {
// 		responseTemplate := d.Get("response_template").(string)
// 		input.ResponseTemplate = opslevel.RefOf(responseTemplate)
// 	}

// 	input.Published = opslevel.Bool(d.Get("published").(bool))

// 	if _, ok := d.GetOk("entity_type"); ok {
// 		entityType := d.Get("entity_type").(string)
// 		input.EntityType = opslevel.RefOf(opslevel.CustomActionsEntityTypeEnum(entityType))
// 	}

// 	resource, err := client.CreateTriggerDefinition(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceTriggerDefinitionRead(d, client)
// }

// func resourceTriggerDefinitionRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetTriggerDefinition(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}
// 	if err := d.Set("description", resource.Description); err != nil {
// 		return err
// 	}
// 	if err := d.Set("owner", resource.Owner.Id); err != nil {
// 		return err
// 	}
// 	if err := d.Set("action", resource.Action.Id); err != nil {
// 		return err
// 	}
// 	if err := d.Set("filter", resource.Filter.Id); err != nil {
// 		return err
// 	}
// 	if err := d.Set("manual_inputs_definition", resource.ManualInputsDefinition); err != nil {
// 		return err
// 	}
// 	if err := d.Set("published", resource.Published); err != nil {
// 		return err
// 	}
// 	if err := d.Set("access_control", string(resource.AccessControl)); err != nil {
// 		return err
// 	}
// 	if err := d.Set("response_template", resource.ResponseTemplate); err != nil {
// 		return err
// 	}
// 	if _, ok := d.GetOk("entity_type"); ok {
// 		if err := d.Set("entity_type", resource.EntityType); err != nil {
// 			return err
// 		}
// 	}
// 	if _, ok := d.GetOk("extended_team_access"); ok {
// 		extendedTeamAccess, err := resource.ExtendedTeamAccess(client, nil)
// 		if err != nil {
// 			return err
// 		}
// 		teams := flattenTeamsArray(extendedTeamAccess)
// 		if err := d.Set("extended_team_access", teams); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func resourceTriggerDefinitionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	input := opslevel.CustomActionsTriggerDefinitionUpdateInput{
// 		Id: *opslevel.NewID(id),
// 	}

// 	if d.HasChange("name") {
// 		input.Name = opslevel.RefOf(d.Get("name").(string))
// 	}
// 	if d.HasChange("description") {
// 		input.Description = opslevel.RefOf(d.Get("description").(string))
// 	}
// 	if d.HasChange("owner") {
// 		input.OwnerId = opslevel.NewID(d.Get("owner").(string))
// 	}
// 	if d.HasChange("action") {
// 		input.ActionId = opslevel.NewID(d.Get("action").(string))
// 	}
// 	if d.HasChange("manual_inputs_definition") {
// 		manualInputsDefinition := d.Get("manual_inputs_definition").(string)
// 		input.ManualInputsDefinition = &manualInputsDefinition
// 	}

// 	if d.HasChange("filter") {
// 		input.FilterId = opslevel.NewID(d.Get("filter").(string))
// 	}

// 	input.Published = opslevel.Bool(d.Get("published").(bool))

// 	if d.HasChange("access_control") {
// 		input.AccessControl = opslevel.RefOf(opslevel.CustomActionsTriggerDefinitionAccessControlEnum(d.Get("access_control").(string)))
// 	}

// 	if d.HasChange("response_template") {
// 		responseTemplate := d.Get("response_template").(string)
// 		input.ResponseTemplate = &responseTemplate
// 	}

// 	if d.HasChange("entity_type") {
// 		entityType := d.Get("entity_type").(string)
// 		input.EntityType = opslevel.RefOf(opslevel.CustomActionsEntityTypeEnum(entityType))
// 	}

// 	extended_teams := opslevel.NewIdentifierArray(getStringArray(d, "extended_team_access"))
// 	if d.HasChange("extended_team_access") {
// 		input.ExtendedTeamAccess = &extended_teams
// 	}

// 	_, err := client.UpdateTriggerDefinition(input)
// 	if err != nil {
// 		return err
// 	}

// 	return resourceTriggerDefinitionRead(d, client)
// }

// func resourceTriggerDefinitionDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteTriggerDefinition(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
