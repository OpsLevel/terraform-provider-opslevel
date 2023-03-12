package opslevel

import (
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceRubricCategories() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceRubricCategoriesRead),
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

func datasourceRubricCategoriesRead(d *schema.ResourceData, client *opslevel.Client) error {

	result, err := client.ListCategories(nil)
	if err != nil {
		return err
	}

	count := result.TotalCount
	aliases := make([]string, count)
	ids := make([]string, count)
	indexes := make([]int, count)
	names := make([]string, count)
	for i, item := range result.Nodes {
		ids[i] = string(item.Id)
		names[i] = item.Name
	}

	d.SetId(timeID())
	d.Set("aliases", aliases)
	d.Set("ids", ids)
	d.Set("indexes", indexes)
	d.Set("names", names)

	return nil
}
