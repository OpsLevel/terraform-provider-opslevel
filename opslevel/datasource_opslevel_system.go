package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceSystem() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceSystemRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the system to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "The aliases of the system.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the system.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the system.",
				Computed:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The id of the team that owns the system.",
				Computed:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The id of the domain this system is child to.",
				Computed:    true,
			},
		},
	}
}

func datasourceSystemRead(d *schema.ResourceData, client *opslevel.Client) error {
	identifier := d.Get("identifier").(string)
	resource, err := client.GetSystem(identifier)
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("aliases", resource.Aliases)
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("owner", resource.Owner.Id())
	d.Set("domain", resource.Parent.Id)

	return nil
}
