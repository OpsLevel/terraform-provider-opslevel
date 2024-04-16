package opslevel

// import (
// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceTeams() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceTeamsRead),
// 		Schema: map[string]*schema.Schema{
// 			"filter": getDatasourceFilter(false, []string{"manager-email"}),
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
// 		},
// 	}
// }

// func datasourceTeamsRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	field := d.Get("filter.0.field").(string)
// 	value := d.Get("filter.0.value").(string)

// 	var teams *opslevel.TeamConnection
// 	var err error
// 	switch field {
// 	case "manager-email":
// 		teams, err = client.ListTeamsWithManager(value, nil)
// 	default:
// 		teams, err = client.ListTeams(nil)
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	count := len(teams.Nodes)
// 	aliases := make([]string, count)
// 	ids := make([]string, count)
// 	names := make([]string, count)
// 	for i, item := range teams.Nodes {
// 		aliases[i] = item.Alias
// 		ids[i] = string(item.Id)
// 		names[i] = item.Name
// 	}

// 	d.SetId(timeID())
// 	d.Set("aliases", aliases)
// 	d.Set("ids", ids)
// 	d.Set("names", names)

// 	return nil
// }
