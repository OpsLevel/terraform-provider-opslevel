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
	input := opslevel.CheckServiceOwnershipCreateInput{}
	setCheckCreateInput(d, &input)

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

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckServiceOwnershipUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceOwnershipUpdateInput{}
	setCheckUpdateInput(d, &input)

	_, err := client.UpdateCheckServiceOwnership(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServiceOwnershipRead(d, client)
}
