package opslevel

import (
	"fmt"
	"strconv"

	"github.com/opslevel/opslevel-go/v2022"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceRubricLevel() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceRubricLevelRead),
		Schema: map[string]*schema.Schema{
			"filter": getDatasourceFilter(true, []string{"alias", "id", "index", "name"}),
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func filterRubricLevels(levels []opslevel.Level, field string, value string) (*opslevel.Level, error) {
	if value == "" {
		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
	}

	var output opslevel.Level
	found := false
	for _, item := range levels {
		switch field {
		case "alias":
			if item.Alias == value {
				output = item
				found = true
			}
		case "id":
			if item.Id.(string) == value {
				output = item
				found = true
			}
		case "index":
			if v, err := strconv.Atoi(value); err == nil && item.Index == v {
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
		return nil, fmt.Errorf("Unable to find level with: %s==%s", field, value)
	}
	return &output, nil
}

func datasourceRubricLevelRead(d *schema.ResourceData, client *opslevel.Client) error {

	results, err := client.ListLevels()
	if err != nil {
		return err
	}

	field := d.Get("filter.0.field").(string)
	value := d.Get("filter.0.value").(string)

	item, itemErr := filterRubricLevels(results, field, value)
	if itemErr != nil {
		return itemErr
	}

	d.SetId(item.Id.(string))
	d.Set("alias", item.Alias)
	d.Set("index", item.Index)
	d.Set("name", item.Name)

	return nil
}
