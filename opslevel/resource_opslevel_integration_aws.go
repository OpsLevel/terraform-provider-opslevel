package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceIntegrationAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a AWS Integration",
		Create:      wrap(resourceIntegrationAWSCreate),
		Read:        wrap(resourceIntegrationAWSRead),
		Update:      wrap(resourceIntegrationAWSUpdate),
		Delete:      wrap(resourceIntegrationAWSDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceIntegrationAWSCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.AWSIntegrationInput{}

	resource, err := client.CreateAWSIntegration(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceIntegrationAWSRead(d, client)
}

func resourceIntegrationAWSRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	_, err := client.GetIntegration(opslevel.ID(id))
	if err != nil {
		return err
	}

	return nil
}

func resourceIntegrationAWSUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.AWSIntegrationInput{}

	_, err := client.UpdateAWSIntegration(d.Id(), input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceIntegrationAWSRead(d, client)
}

func resourceIntegrationAWSDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteAWSIntegration(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
