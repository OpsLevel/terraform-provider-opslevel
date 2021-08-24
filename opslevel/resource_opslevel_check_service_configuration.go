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
	input := opslevel.CheckServiceConfigurationCreateInput{}
	setCheckCreateInput(d, &input)

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

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckServiceConfigurationUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceConfigurationUpdateInput{}
	setCheckUpdateInput(d, &input)

	_, err := client.UpdateCheckServiceConfiguration(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServiceConfigurationRead(d, client)
}
