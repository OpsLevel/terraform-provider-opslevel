package opslevel

import (
	"github.com/opslevel/opslevel-go/v2023"

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
				Type:        schema.TypeString,
				Description: "The name of the group.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the group.",
				Computed:    true,
			},
			"alias": {
				Type:        schema.TypeString,
				Description: "The alias of the group.",
				Computed:    true,
			},
			"parent": {
				Type:        schema.TypeString,
				Description: "The parent alias of the group.",
				Computed:    true,
			},
			"members": {
				Type:        schema.TypeList,
				Description: "The users' emails who are members of the group.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"teams": {
				Type:        schema.TypeList,
				Description: "The team aliases where this group is the direct parent.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceGroupRead(d *schema.ResourceData, client *opslevel.Client) error {
	identifier := d.Get("identifier").(string)
	var resource *opslevel.Group
	var err error
	if opslevel.IsID(identifier) {
		resource, err = client.GetGroup(opslevel.ID(identifier))
	} else {
		resource, err = client.GetGroupWithAlias(identifier)
	}
	if err != nil {
		return err
	}

	var parent string
	if resource.Parent.Alias != "" {
		parent = resource.Parent.Alias
	}

	groupMembers, err := resource.Members(client, nil)
	if err != nil {
		return err
	}
	members := flattenMembersArray(groupMembers)

	childTeams, err := resource.ChildTeams(client, nil)
	if err != nil {
		return err
	}
	teams := flattenTeamsArray(childTeams)

	d.SetId(string(resource.Id))
	d.Set("name", resource.Name)
	d.Set("alias", resource.Alias)
	d.Set("description", resource.Description)
	d.Set("parent", parent)
	d.Set("members", members)
	d.Set("teams", teams)

	return nil
}
