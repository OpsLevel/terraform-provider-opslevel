package opslevel

import (
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceFilter() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceFilterRead),
		Schema: map[string]*schema.Schema{
			"filter": getDatasourceFilter(true, []string{"id", "name"}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func filterFilters(data []opslevel.Filter, field string, value string) (*opslevel.Filter, error) {
	if value == "" {
		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
	}

	var output opslevel.Filter
	found := false
	for _, item := range data {
		switch field {
		case "id":
			if string(item.Id) == value {
				output = item
				found = true
			}
		case "name":
			if item.Name == value {
				output = item
				found = true
			}
		}
		if found {
			break
		}
	}

	if found == false {
		return nil, fmt.Errorf("Unable to find filter with: %s==%s", field, value)
	}
	return &output, nil
}

func datasourceFilterRead(d *schema.ResourceData, client *opslevel.Client) error {
	resp, err := client.ListFilters(nil)
	results := resp.Nodes
	if err != nil {
		return err
	}

	field := d.Get("filter.0.field").(string)
	value := d.Get("filter.0.value").(string)

	item, itemErr := filterFilters(results, field, value)
	if itemErr != nil {
		return itemErr
	}

	d.SetId(string(item.Id))
	d.Set("name", item.Name)

	return nil
}
