package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2024"
)

func resourceCheckRepositoryIntegrated() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a repository integrated check.",
		Create:      wrap(resourceCheckRepositoryIntegratedCreate),
		Read:        wrap(resourceCheckRepositoryIntegratedRead),
		Update:      wrap(resourceCheckRepositoryIntegratedUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(nil),
	}
}

func resourceCheckRepositoryIntegratedCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckRepositoryIntegratedCreateInput](checkCreateInput)

	resource, err := client.CreateCheckRepositoryIntegrated(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckRepositoryIntegratedRead(d, client)
}

func resourceCheckRepositoryIntegratedRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckRepositoryIntegratedUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckRepositoryIntegratedUpdateInput](checkUpdateInput)

	if _, err := client.UpdateCheckRepositoryIntegrated(*input); err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositoryIntegratedRead(d, client)
}
