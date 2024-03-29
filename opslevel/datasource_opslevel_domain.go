package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2024"
)

func datasourceDomain() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceDomainRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the domain to find.",
				ForceNew:    true,
				Required:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "The aliases of the domain.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the domain.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the domain.",
				Computed:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The id of the team that owns the domain.",
				Computed:    true,
			},
		},
	}
}

func datasourceDomainRead(d *schema.ResourceData, client *opslevel.Client) error {
	identifier := d.Get("identifier").(string)
	resource, err := client.GetDomain(identifier)
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("aliases", resource.ManagedAliases)
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("owner", resource.Owner.Id())

	return nil
}
