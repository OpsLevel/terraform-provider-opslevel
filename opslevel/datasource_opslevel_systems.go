package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceSystems() *schema.Resource {
	return &schema.Resource{
		Read:   wrap(datasourceSystemsRead),
		Schema: map[string]*schema.Schema{},
	}
}

func datasourceSystemsRead(d *schema.ResourceData, client *opslevel.Client) error {

	return nil
}
