package opslevel

import (
	"github.com/opslevel/opslevel-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceRubricLevels() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceRubricLevelsRead),
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
			"indexes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceRubricLevelsRead(d *schema.ResourceData, client *opslevel.Client) error {

	result, err := client.ListLevels()
	if err != nil {
		return err
	}

	count := len(result) - 1
	aliases := make([]string, count)
	ids := make([]string, count)
	indexes := make([]int, count)
	names := make([]string, count)
	for _, item := range result {
		if item.Index == 0 {
			continue
		}
		i := item.Index - 1
		aliases[i] = item.Alias
		ids[i] = item.Id.(string)
		indexes[i] = item.Index
		names[i] = item.Name
	}

	d.SetId(timeID())
	d.Set("aliases", aliases)
	d.Set("ids", ids)
	d.Set("indexes", indexes)
	d.Set("names", names)

	return nil
}
