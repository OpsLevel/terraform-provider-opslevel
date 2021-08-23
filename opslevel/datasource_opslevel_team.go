package opslevel

import (
	"github.com/opslevel/opslevel-go"

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
		},
	}
}

func datasourceTeamRead(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := findTeam("alias", "id", d, client)
	if err != nil {
		return err
	}

	d.SetId(resource.Id.(string))
	d.Set("alias", resource.Alias)
	d.Set("name", resource.Name)

	return nil
}
