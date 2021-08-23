package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceCheckServiceOwnership() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service ownership check.",
		Create:      wrap(resourceCheckServiceOwnershipCreate),
		Read:        wrap(resourceCheckServiceOwnershipRead),
		Update:      wrap(resourceCheckServiceOwnershipUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(nil),
	}
}

func resourceCheckServiceOwnershipCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceOwnershipCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),
	}
	resource, err := client.CreateCheckServiceOwnership(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckServiceOwnershipRead(d, client)
}

func resourceCheckServiceOwnershipRead(d *schema.ResourceData, client *opslevel.Client) error {
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

func resourceCheckServiceOwnershipUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceOwnershipUpdateInput{
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

	_, err := client.UpdateCheckServiceOwnership(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServiceOwnershipRead(d, client)
}
