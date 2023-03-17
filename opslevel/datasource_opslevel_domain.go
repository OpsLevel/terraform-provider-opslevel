package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceDomain() *schema.Resource {
	return &schema.Resource{
		Read:   wrap(datasourceDomainRead),
		Schema: map[string]*schema.Schema{},
	}
}

func datasourceDomainRead(d *schema.ResourceData, client *opslevel.Client) error {

	return nil
}
