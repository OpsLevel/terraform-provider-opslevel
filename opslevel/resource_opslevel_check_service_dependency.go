package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceCheckServiceDependency() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service dependency check",
		Create:      wrap(resourceCheckServiceDependencyCreate),
		Read:        wrap(resourceCheckServiceDependencyRead),
		Update:      wrap(resourceCheckServiceDependencyUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{}),
	}
}

func resourceCheckServiceDependencyCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckServiceDependencyCreateInput](checkCreateInput)

	resource, err := client.CreateCheckServiceDependency(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckServiceDependencyRead(d, client)
}

func resourceCheckServiceDependencyRead(d *schema.ResourceData, client *opslevel.Client) error {
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

func resourceCheckServiceDependencyUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckServiceDependencyUpdateInput](checkUpdateInput)

	if _, err := client.UpdateCheckServiceDependency(*input); err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServiceDependencyRead(d, client)
}
