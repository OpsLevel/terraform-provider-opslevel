// package opslevel

// import (
// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceUsers() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceUsersRead),
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
// 			"emails": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"roles": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 		},
// 	}
// }

// func datasourceUsersRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resp, err := client.ListUsers(nil)
// 	if err != nil {
// 		return err
// 	}
// 	result := resp.Nodes
// 	count := len(result)
// 	ids := make([]string, count)
// 	names := make([]string, count)
// 	emails := make([]string, count)
// 	roles := make([]string, count)
// 	for i, item := range result {
// 		ids[i] = string(item.Id)
// 		names[i] = item.Name
// 		emails[i] = item.Email
// 		roles[i] = string(item.Role)
// 	}

// 	d.SetId(timeID())
// 	d.Set("ids", ids)
// 	d.Set("names", names)
// 	d.Set("emails", emails)
// 	d.Set("roles", roles)

// 	return nil
// }
