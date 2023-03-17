package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a domain",
		Create:      wrap(resourceDomainCreate),
		Read:        wrap(resourceDomainRead),
		Update:      wrap(resourceDomainUpdate),
		Delete:      wrap(resourceDomainDelete),
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

func resourceDomainCreate(d *schema.ResourceData, client *opslevel.Client) error {

	return resourceDomainRead(d, client)
}

func resourceDomainRead(d *schema.ResourceData, client *opslevel.Client) error {

	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, client *opslevel.Client) error {

	d.Set("last_updated", timeLastUpdated())
	return resourceDomainRead(d, client)
}

func resourceDomainDelete(d *schema.ResourceData, client *opslevel.Client) error {

	return nil
}
