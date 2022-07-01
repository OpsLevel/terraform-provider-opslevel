package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func datasourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceGroupsRead),
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
	var groups []opslevel.Group
	var err error

	groups, err = client.ListGroups()
	if err != nil {
		return err
	}

	count := len(groups)
	aliases := make([]string, count)
	ids := make([]string, count)
	names := make([]string, count)
	for i, item := range groups {
		aliases[i] = item.Alias
		ids[i] = item.Id.(string)
		names[i] = item.Name
	}

	d.SetId(timeID())
	d.Set("aliases", aliases)
	d.Set("ids", ids)
	d.Set("names", names)

	return nil
}
