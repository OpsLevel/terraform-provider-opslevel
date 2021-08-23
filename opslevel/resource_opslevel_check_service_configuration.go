package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceCheckServiceConfiguration() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service configuration check.",
		Create:      wrap(resourceCheckServiceConfigurationCreate),
		Read:        wrap(resourceCheckServiceConfigurationRead),
		Update:      wrap(resourceCheckServiceConfigurationUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(nil),
	}
}

func resourceCheckServiceConfigurationCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceConfigurationCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),
	}
	resource, err := client.CreateCheckServiceConfiguration(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckServiceConfigurationRead(d, client)
}

func resourceCheckServiceConfigurationRead(d *schema.ResourceData, client *opslevel.Client) error {
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

func resourceCheckServiceConfigurationUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceConfigurationUpdateInput{
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

	_, err := client.UpdateCheckServiceConfiguration(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServiceConfigurationRead(d, client)
}
