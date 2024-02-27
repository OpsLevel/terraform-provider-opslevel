package opslevel

// import (
// 	"fmt"

// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceRubricCategory() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceRubricCategoryRead),
// 		Schema: map[string]*schema.Schema{
// 			"filter": getDatasourceFilter(true, []string{"id", "name"}),
// 			"name": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 		},
// 	}
// }

// func filterRubricCategories(levels *opslevel.CategoryConnection, field string, value string) (*opslevel.Category, error) {
// 	if value == "" {
// 		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
// 	}

// 	var output opslevel.Category
// 	found := false
// 	for _, item := range levels.Nodes {
// 		switch field {
// 		case "id":
// 			if string(item.Id) == value {
// 				output = item
// 				found = true
// 			}
// 		case "name":
// 			if item.Name == value {
// 				output = item
// 				found = true
// 			}
// 		}
// 		if found {
// 			break
// 		}
// 	}

// 	if !found {
// 		return nil, fmt.Errorf("Unable to find category with: %s==%s", field, value)
// 	}
// 	return &output, nil
// }

// func datasourceRubricCategoryRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	results, err := client.ListCategories(nil)
// 	if err != nil {
// 		return err
// 	}

// 	field := d.Get("filter.0.field").(string)
// 	value := d.Get("filter.0.value").(string)

// 	item, itemErr := filterRubricCategories(results, field, value)
// 	if itemErr != nil {
// 		return itemErr
// 	}

// 	d.SetId(string(item.Id))
// 	d.Set("name", item.Name)

// 	return nil
// }
