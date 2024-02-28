package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2024"
)

func resourceCheckHasRecentDeploy() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a has recent deploy check",
		Create:      wrap(resourceCheckHasRecentDeployCreate),
		Read:        wrap(resourceCheckHasRecentDeployRead),
		Update:      wrap(resourceCheckHasRecentDeployUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"days": {
				Type:        schema.TypeInt,
				Description: "The number of days to check since the last deploy.",
				ForceNew:    false,
				Required:    true,
			},
		}),
	}
}

func resourceCheckHasRecentDeployCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckHasRecentDeployCreateInput](checkCreateInput)
	input.Days = d.Get("days").(int)

	resource, err := client.CreateCheckHasRecentDeploy(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckHasRecentDeployRead(d, client)
}

func resourceCheckHasRecentDeployRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("days", resource.HasRecentDeployCheckFragment.Days); err != nil {
		return err
	}

	return nil
}

func resourceCheckHasRecentDeployUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckHasRecentDeployUpdateInput](checkUpdateInput)
	if d.HasChange("days") {
		input.Days = opslevel.RefOf(d.Get("days").(int))
	}

	_, err := client.UpdateCheckHasRecentDeploy(*input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckHasRecentDeployRead(d, client)
}
