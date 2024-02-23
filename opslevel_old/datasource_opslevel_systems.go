// package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func datasourceSystems() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceSystemsRead),
// 		Schema: map[string]*schema.Schema{
// 			"aliases": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
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
// 			"owners": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"domains": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 		},
// 	}
// }

// func datasourceSystemsRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resp, err := client.ListSystems(nil)
// 	if err != nil {
// 		return err
// 	}

// 	count := len(resp.Nodes)
// 	aliases := make([]string, count)
// 	ids := make([]string, count)
// 	names := make([]string, count)
// 	descriptions := make([]string, count)
// 	owners := make([]string, count)
// 	domains := make([]string, count)
// 	for i, item := range resp.Nodes {
// 		if len(item.Aliases) > 0 {
// 			aliases[i] = item.Aliases[0]
// 		}
// 		ids[i] = string(item.Id)
// 		names[i] = item.Name
// 		descriptions[i] = item.Description
// 		owners[i] = string(item.Owner.Id())
// 		domains[i] = string(item.Parent.Id)
// 	}

// 	d.SetId(timeID())
// 	d.Set("aliases", aliases)
// 	d.Set("ids", ids)
// 	d.Set("names", names)
// 	d.Set("descriptions", descriptions)
// 	d.Set("owners", owners)
// 	d.Set("domains", domains)

// 	return nil
// }
