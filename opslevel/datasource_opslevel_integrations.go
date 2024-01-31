package opslevel

import (
	"fmt"
	"github.com/opslevel/opslevel-go/v2024"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIntegrations() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceIntegrationsRead),
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceIntegrationsRead(d *schema.ResourceData, client *opslevel.Client) error {
	resp, err := client.ListIntegrations(nil)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("unexpected: listing integrations returned nil")
	}
	result := resp.Nodes

	count := len(result)
	ids := make([]string, count)
	names := make([]string, count)
	for i, item := range result {
		ids[i] = string(item.Id)
		names[i] = item.Name
	}

	d.SetId(timeID())
	d.Set("ids", ids)
	d.Set("names", names)

	return nil
}
