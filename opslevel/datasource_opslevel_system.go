package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceSystem() *schema.Resource {
	return &schema.Resource{
		Read:   wrap(datasourceSystemRead),
		Schema: map[string]*schema.Schema{},
	}
}

func datasourceSystemRead(d *schema.ResourceData, client *opslevel.Client) error {

	return nil
}
