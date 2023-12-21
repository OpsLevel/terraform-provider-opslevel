package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourcePropertyDefinition() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a property definition",
		Create:      wrap(resourcePropertyDefinitionCreate),
		Read:        wrap(resourcePropertyDefinitionRead),
		Update:      wrap(resourcePropertyDefinitionUpdate),
		Delete:      wrap(resourcePropertyDefinitionDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the property definition.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the property definition.",
				Optional:    true,
			},
			"schema": {
				Type:        schema.TypeString,
				Description: "The schema of the property definition.",
				Required:    true,
			},
			"property_display_status": {
				Type:        schema.TypeString,
				Description: "The display status of a custom property on service pages. (Options: 'visible' or 'hidden')",
				Default:     "visible",
				Optional:    true,
			},
		},
	}
}

func resourcePropertyDefinitionCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.PropertyDefinitionInput{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Schema:                opslevel.NewJSON(d.Get("schema").(string)),
		PropertyDisplayStatus: opslevel.PropertyDisplayStatusEnum(d.Get("property_display_status").(string)),
	}

	resource, err := client.CreatePropertyDefinition(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourcePropertyDefinitionRead(d, client)
}

func resourcePropertyDefinitionRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	resource, err := client.GetPropertyDefinition(id)
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("schema", resource.Schema.ToJSON()); err != nil {
		return err
	}
	if err := d.Set("property_display_status", string(resource.PropertyDisplayStatus)); err != nil {
		return err
	}

	return nil
}

func resourcePropertyDefinitionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	input := opslevel.PropertyDefinitionInput{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Schema:                opslevel.NewJSON(d.Get("schema").(string)),
		PropertyDisplayStatus: opslevel.PropertyDisplayStatusEnum(d.Get("property_display_status").(string)),
	}

	_, err := client.UpdatePropertyDefinition(id, input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourcePropertyDefinitionRead(d, client)
}

func resourcePropertyDefinitionDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeletePropertyDefinition(id)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
