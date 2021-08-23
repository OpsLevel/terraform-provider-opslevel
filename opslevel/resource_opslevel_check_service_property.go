package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
)

func resourceCheckServiceProperty() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service property check.",
		Create:      wrap(resourceCheckServicePropertyCreate),
		Read:        wrap(resourceCheckServicePropertyRead),
		Update:      wrap(resourceCheckServicePropertyUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"property": {
				Type:         schema.TypeString,
				Description:  "The property of the service that the check will verify.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.GetServicePropertyTypes(), false),
			},
			"predicate": getPredicateInputSchema(false),
		}),
	}
}

func resourceCheckServicePropertyCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServicePropertyCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		Property:  opslevel.ServiceProperty(d.Get("property").(string)),
		Predicate: getPredicateInput(d, "predicate"),
	}
	resource, err := client.CreateCheckServiceProperty(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckServicePropertyRead(d, client)
}

func resourceCheckServicePropertyRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
	if err != nil {
		return err
	}

	if err := resourceCheckRead(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckServicePropertyUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServicePropertyUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("enabled") {
		value := d.Get("enabled").(bool)
		input.Enabled = &value
	}
	if d.HasChange("category") {
		input.Category = getID(d, "category")
	}
	if d.HasChange("level") {
		input.Level = getID(d, "level")
	}
	if d.HasChange("owner") {
		input.Owner = getID(d, "owner")
	}
	if d.HasChange("filter") {
		input.Filter = getID(d, "filter")
	}
	if d.HasChange("notes") {
		input.Notes = d.Get("notes").(string)
	}

	if d.HasChange("property") {
		input.Property = opslevel.ServiceProperty(d.Get("property").(string))
	}
	if d.HasChange("predicate") {
		input.Predicate = getPredicateInput(d, "predicate")
	}

	_, err := client.UpdateCheckServiceProperty(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServicePropertyRead(d, client)
}
