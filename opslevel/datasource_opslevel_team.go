package opslevel

import (
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceTeam() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceTeamRead),
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Description: "An alias of the team to find by.",
				ForceNew:    true,
				Optional:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the team to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_team_alias": {
				Type:        schema.TypeString,
				Description: "The alias of the parent team.",
				Computed:    true,
			},
			"parent_team_id": {
				Type:        schema.TypeString,
				Description: "The id of the parent team.",
				Computed:    true,
			},
			"group_alias": {
				Type:        schema.TypeString,
				Description: "The name of the group the team belongs to.",
				Computed:    true,
				Deprecated:  "field 'group' on team is no longer supported please use the 'parent_team' field.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Description: "The id of the group the team belongs to.",
				Computed:    true,
				Deprecated:  "field 'group' on team is no longer supported please use the 'parent_team' field.",
			},
		},
	}
}

func datasourceTeamRead(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := findTeam("alias", "id", d, client)
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("alias", resource.Alias)
	d.Set("name", resource.Name)
	if err := d.Set("group_alias", resource.Group.Alias); err != nil {
		return err
	}
	if err := d.Set("group_id", resource.Group.Id); err != nil {
		return err
	}
	if err := d.Set("parent_team_alias", resource.ParentTeam.Alias); err != nil {
		return err
	}
	if err := d.Set("parent_team_id", resource.ParentTeam.Id); err != nil {
		return err
	}

	return nil
}
