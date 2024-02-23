// package opslevel

// import (
// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourcePropertyDefinitions() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourcePropertyDefinitionsRead),
// 		Schema: map[string]*schema.Schema{
// 			"ids": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"names": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"descriptions": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"schemas": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"property_display_statuses": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 		},
// 	}
// }

// func datasourcePropertyDefinitionsRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resources, err := client.ListPropertyDefinitions(nil)
// 	if err != nil {
// 		return err
// 	}

// 	count := len(resources.Nodes)
// 	ids := make([]string, count)
// 	names := make([]string, count)
// 	descriptions := make([]string, count)
// 	schemas := make([]string, count)
// 	propertyDisplayStatuses := make([]string, count)
// 	for i, item := range resources.Nodes {
// 		ids[i] = string(item.Id)
// 		names[i] = item.Name
// 		descriptions[i] = item.Description
// 		schemas[i] = item.Schema.ToJSON()
// 		propertyDisplayStatuses[i] = string(item.PropertyDisplayStatus)
// 	}

// 	d.SetId(timeID())
// 	d.Set("ids", ids)
// 	d.Set("names", names)
// 	d.Set("descriptions", descriptions)
// 	d.Set("schemas", schemas)
// 	d.Set("property_display_statuses", propertyDisplayStatuses)

// 	return nil
// }
