package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceGroups() *schema.Resource {
	return &schema.Resource{
		Read:               wrap(datasourceGroupsRead),
		DeprecationMessage: "Groups are being deprecated. Please replace Groups with Teams.",
		Schema: map[string]*schema.Schema{
			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
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

func datasourceGroupsRead(d *schema.ResourceData, client *opslevel.Client) error {
	groups, err := client.ListGroups(nil)
	if err != nil {
		return err
	}

	count := len(groups.Nodes)
	aliases := make([]string, count)
	ids := make([]string, count)
	names := make([]string, count)
	for i, item := range groups.Nodes {
		aliases[i] = item.Alias
		ids[i] = string(item.Id)
		names[i] = item.Name
	}

	d.SetId(timeID())
	d.Set("aliases", aliases)
	d.Set("ids", ids)
	d.Set("names", names)

	return nil
}
