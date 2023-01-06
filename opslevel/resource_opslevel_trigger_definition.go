package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hasura/go-graphql-client"
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
				Description: "The description of what the trigger definition will do.",
				ForceNew:    false,
				Optional:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The owner of the Trigger Definition",
				ForceNew:    false,
				Required:    true,
			},
			"action": {
				Type:        schema.TypeString,
				Description: "The action that will be triggered by the Trigger Definition",
				ForceNew:    false,
				Required:    true,
			},
			"filter": {
				Type:        schema.TypeString,
				Description: "A filter defining which services this trigger definition applies to.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceTriggerDefinitionCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CustomActionsTriggerDefinitionCreateInput{
		Name:   d.Get("name").(string),
		Owner:  *opslevel.NewID(d.Get("owner").(string)),
		Action: opslevel.NewID(d.Get("action").(string)),
	}

	if _, ok := d.GetOk("description"); !ok {
		input.Description = nil
	} else {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}

	if _, ok := d.GetOk("filter"); !ok {
		input.Filter = nil
	} else {
		input.Filter = opslevel.NewID(d.Get("filter").(string))
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

	return nil
}

func resourceTriggerDefinitionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	input := opslevel.CustomActionsTriggerDefinitionUpdateInput{
		Id: graphql.ID(id),
	}

	if d.HasChange("name") {
		input.Name = opslevel.NewString(d.Get("name").(string))
	}
	if d.HasChange("description") {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}
	if d.HasChange("owner") {
		input.Owner = opslevel.NewString(d.Get("owner").(string))
	}
	if d.HasChange("action") {
		input.Action = opslevel.NewString(d.Get("action").(string))
	}

	if d.HasChange("filter") {
		filter, ok := d.GetOk("filter")
		if ok {
			input.Filter = opslevel.NewString(filter.(string))
		} else {
			input.Filter = opslevel.NullString()
		}
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
