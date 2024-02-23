// package opslevel

// import (
// 	"fmt"

// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceIntegration() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceIntegrationRead),
// 		Schema: map[string]*schema.Schema{
// 			"filter": getDatasourceFilter(true, []string{"id", "name"}),
// 			"name": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 		},
// 	}
// }

// func filterIntegrations(data []opslevel.Integration, field string, value string) (*opslevel.Integration, error) {
// 	if value == "" {
// 		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
// 	}

// 	var output opslevel.Integration
// 	found := false
// 	for _, item := range data {
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
// 		return nil, fmt.Errorf("Unable to find integration with: %s==%s", field, value)
// 	}
// 	return &output, nil
// }

// func datasourceIntegrationRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resp, err := client.ListIntegrations(nil)
// 	if err != nil {
// 		return err
// 	}
// 	if resp == nil {
// 		return fmt.Errorf("unexpected: listing integrations returned nil")
// 	}
// 	results := resp.Nodes

// 	field := d.Get("filter.0.field").(string)
// 	value := d.Get("filter.0.value").(string)

// 	item, itemErr := filterIntegrations(results, field, value)
// 	if itemErr != nil {
// 		return itemErr
// 	}

// 	d.SetId(string(item.Id))
// 	d.Set("names", item.Name)

// 	return nil
// }
