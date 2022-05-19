package opslevel

import (
	"github.com/opslevel/opslevel-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceGroupRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the group to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGroupRead(d *schema.ResourceData, client *opslevel.Client) error {
	identifier := d.Get("identifier").(string)
	var resource *opslevel.Group
	var err error
	if opslevel.IsID(identifier) {
		resource, err = client.GetGroup(identifier)
	} else {
		resource, err = client.GetGroupWithAlias(identifier)
	}
	if err != nil {
		return err
	}

	parent := map[string]string{}
	if resource.Parent.Id != nil {
		parent = map[string]string{
			"alias": resource.Parent.Alias,
			"id":    resource.Parent.Id.(string),
		}
	}

	members, err := resource.Members(client)
	members_list := []map[string]string{}
	if err != nil {
		return err
	}
	if len(members) > 0 {
		for _, member := range members {
			members_list = append(members_list, map[string]string{
				"email": member.Email,
			})
		}
	}

	childTeams, err := resource.ChildTeams(client)
	teams := []map[string]string{}
	if err != nil {
		return err
	}
	if len(childTeams) > 0 {
		for _, team := range childTeams {
			teams = append(teams, map[string]string{
				"alias": team.Alias,
				"id":    team.Id.(string),
			})
		}
	}

	d.SetId(resource.Id.(string))
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("parent", parent)
	d.Set("members", members_list)
	d.Set("teams", teams)

	return nil
}
