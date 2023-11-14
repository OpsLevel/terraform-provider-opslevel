package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceTriggerDefinition() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a webhook action",
		Create:      wrap(resourceTriggerDefinitionCreate),
		Read:        wrap(resourceTriggerDefinitionRead),
		Update:      wrap(resourceTriggerDefinitionUpdate),
		Delete:      wrap(resourceTriggerDefinitionDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Trigger Definition",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of what the Trigger Definition will do.",
				ForceNew:    false,
				Optional:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The owner of the Trigger Definition.",
				ForceNew:    false,
				Required:    true,
			},
			"action": {
				Type:        schema.TypeString,
				Description: "The action that will be triggered by the Trigger Definition.",
				ForceNew:    false,
				Required:    true,
			},
			"filter": {
				Type:        schema.TypeString,
				Description: "A filter defining which services this Trigger Definition applies to.",
				ForceNew:    false,
				Optional:    true,
			},
			"manual_inputs_definition": {
				Type:        schema.TypeString,
				Description: "The YAML definition of any custom inputs for this Trigger Definition.",
				ForceNew:    false,
				Optional:    true,
			},
			"published": {
				Type:        schema.TypeBool,
				Description: "The published state of the Custom Action; true if the Trigger Definition is ready for use; false if it is a draft. Defaults to false.",
				Default:     false,
				ForceNew:    false,
				Optional:    true,
			},
			"access_control": {
				Type:         schema.TypeString,
				Description:  "The set of users that should be able to use the Trigger Definition. Requires a value of `everyone`, `admins`, or `service_owners`.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllCustomActionsTriggerDefinitionAccessControlEnum, false),
			},
			"response_template": {
				Type:        schema.TypeString,
				Description: "The liquid template used to parse the response from the Webhook Action.",
				ForceNew:    false,
				Optional:    true,
			},
			"entity_type": {
				Type:         schema.TypeString,
				Description:  "The entity type to associate with the Trigger Definition.",
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllCustomActionsEntityTypeEnum, false),
			},
			"extended_team_access": {
				Type:        schema.TypeList,
				Description: "The set of additional teams who can invoke this Trigger Definition.",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTriggerDefinitionCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CustomActionsTriggerDefinitionCreateInput{
		Name:          d.Get("name").(string),
		Owner:         *opslevel.NewID(d.Get("owner").(string)),
		Action:        *opslevel.NewID(d.Get("action").(string)),
		AccessControl: opslevel.CustomActionsTriggerDefinitionAccessControlEnum(d.Get("access_control").(string)),
	}
	extended_teams := opslevel.NewIdentifierArray(getStringArray(d, "extended_team_access"))
	if len(extended_teams) > 0 {
		input.ExtendedTeamAccess = &extended_teams
	}

	if _, ok := d.GetOk("description"); ok {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}

	if _, ok := d.GetOk("filter"); ok {
		input.Filter = opslevel.NewID(d.Get("filter").(string))
	}

	if _, ok := d.GetOk("manual_inputs_definition"); ok {
		manualInputsDefinition := d.Get("manual_inputs_definition").(string)
		input.ManualInputsDefinition = manualInputsDefinition
	}

	if _, ok := d.GetOk("response_template"); ok {
		responseTemplate := d.Get("response_template").(string)
		input.ResponseTemplate = responseTemplate
	}

	if published, ok := d.GetOk("published"); ok {
		input.Published = opslevel.Bool(published.(bool))
	}

	if _, ok := d.GetOk("entity_type"); ok {
		entityType := d.Get("entity_type").(string)
		input.EntityType = opslevel.CustomActionsEntityTypeEnum(entityType)
	}

	resource, err := client.CreateTriggerDefinition(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceTriggerDefinitionRead(d, client)
}

func resourceTriggerDefinitionRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetTriggerDefinition(*opslevel.NewIdentifier(id))
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("owner", resource.Owner.Id); err != nil {
		return err
	}
	if err := d.Set("action", resource.Action.Id); err != nil {
		return err
	}
	if err := d.Set("filter", resource.Filter.Id); err != nil {
		return err
	}
	if err := d.Set("manual_inputs_definition", resource.ManualInputsDefinition); err != nil {
		return err
	}
	if err := d.Set("published", resource.Published); err != nil {
		return err
	}
	if err := d.Set("access_control", string(resource.AccessControl)); err != nil {
		return err
	}
	if err := d.Set("response_template", resource.ResponseTemplate); err != nil {
		return err
	}
	if _, ok := d.GetOk("entity_type"); ok {
		if err := d.Set("entity_type", resource.EntityType); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("extended_team_access"); ok {
		extendedTeamAccess, err := resource.ExtendedTeamAccess(client, nil)
		if err != nil {
			return err
		}
		teams := flattenTeamsArray(extendedTeamAccess)
		if err := d.Set("extended_team_access", teams); err != nil {
			return err
		}
	}

	return nil
}

func resourceTriggerDefinitionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	input := opslevel.CustomActionsTriggerDefinitionUpdateInput{
		Id: *opslevel.NewID(id),
	}

	if d.HasChange("name") {
		input.Name = opslevel.NewString(d.Get("name").(string))
	}
	if d.HasChange("description") {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}
	if d.HasChange("owner") {
		input.Owner = opslevel.NewID(d.Get("owner").(string))
	}
	if d.HasChange("action") {
		input.Action = opslevel.NewID(d.Get("action").(string))
	}
	if d.HasChange("manual_inputs_definition") {
		manualInputsDefinition := d.Get("manual_inputs_definition").(string)
		input.ManualInputsDefinition = &manualInputsDefinition
	}

	if d.HasChange("filter") {
		input.Filter = opslevel.NewID(d.Get("filter").(string))
	}

	if d.HasChange("published") {
		published, ok := d.GetOk("published")
		if ok {
			input.Published = opslevel.Bool(published.(bool))
		} else {
			input.Published = opslevel.Bool(false)
		}
	}

	if d.HasChange("access_control") {
		input.AccessControl = opslevel.CustomActionsTriggerDefinitionAccessControlEnum(d.Get("access_control").(string))
	}

	if d.HasChange("response_template") {
		responseTemplate := d.Get("response_template").(string)
		input.ResponseTemplate = &responseTemplate
	}

	if d.HasChange("entity_type") {
		entityType := d.Get("entity_type").(string)
		input.EntityType = opslevel.CustomActionsEntityTypeEnum(entityType)
	}

	extended_teams := opslevel.NewIdentifierArray(getStringArray(d, "extended_team_access"))
	if d.HasChange("extended_team_access") {
		input.ExtendedTeamAccess = &extended_teams
	}

	_, err := client.UpdateTriggerDefinition(input)
	if err != nil {
		return err
	}

	return resourceTriggerDefinitionRead(d, client)
}

func resourceTriggerDefinitionDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTriggerDefinition(*opslevel.NewIdentifier(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
