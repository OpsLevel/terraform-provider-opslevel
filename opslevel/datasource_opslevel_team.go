package opslevel

// import (
// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceTeam() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceTeamRead),
// 		Schema: map[string]*schema.Schema{
// 			"alias": {
// 				Type:        schema.TypeString,
// 				Description: "An alias of the team to find by.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"id": {
// 				Type:        schema.TypeString,
// 				Description: "The id of the team to find.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The name of the team.",
// 				Computed:    true,
// 			},
// 			"members": {
// 				Type:        schema.TypeList,
// 				Description: "List of members in the team with email address and role.",
// 				Computed:    true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"email": {
// 							Type:        schema.TypeString,
// 							Description: "The email address of the team member.",
// 							Computed:    true,
// 						},
// 						"role": {
// 							Type:        schema.TypeString,
// 							Description: "The role of the team member.",
// 							Computed:    true,
// 						},
// 					},
// 				},
// 			},
// 			"parent_alias": {
// 				Type:        schema.TypeString,
// 				Description: "The alias of the parent team.",
// 				Computed:    true,
// 			},
// 			"parent_id": {
// 				Type:        schema.TypeString,
// 				Description: "The id of the parent team.",
// 				Computed:    true,
// 			},
// 		},
// 	}
// }

// func datasourceTeamRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resource, err := findTeam("alias", "id", d, client)
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId(string(resource.Id))
// 	d.Set("alias", resource.Alias)
// 	d.Set("name", resource.Name)

// 	if err := d.Set("members", mapMembershipsArray(resource.Memberships)); err != nil {
// 		return err
// 	}

// 	if err := d.Set("parent_alias", resource.ParentTeam.Alias); err != nil {
// 		return err
// 	}
// 	if err := d.Set("parent_id", resource.ParentTeam.Id); err != nil {
// 		return err
// 	}

// 	return nil
// }
