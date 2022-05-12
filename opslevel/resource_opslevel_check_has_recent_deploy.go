package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
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
	input := opslevel.CheckHasRecentDeployCreateInput{}
	setCheckCreateInput(d, &input)

	input.Days = d.Get("days").(int)

	resource, err := client.CreateCheckHasRecentDeploy(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckHasRecentDeployRead(d, client)
}

func resourceCheckHasRecentDeployRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
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
	input := opslevel.CheckHasRecentDeployUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("days") {
		input.Days = opslevel.NewInt(d.Get("days").(int))
	}

	_, err := client.UpdateCheckHasRecentDeploy(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckHasRecentDeployRead(d, client)
}
