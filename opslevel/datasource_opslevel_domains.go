package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceDomains() *schema.Resource {
	return &schema.Resource{
		Read:   wrap(datasourceDomainsRead),
		Schema: map[string]*schema.Schema{},
	}
}

func datasourceDomainsRead(d *schema.ResourceData, client *opslevel.Client) error {

	return nil
}
