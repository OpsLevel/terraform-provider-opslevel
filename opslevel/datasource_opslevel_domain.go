package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceDomain() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceDomainRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the domain to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "The aliases of the domain.",
				Computed:    true,
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
				Description: "The id of the domain owner - could be a group or team.",
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
	d.Set("aliases", resource.Aliases)
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	if resource.Owner.GroupId.Id == "" {
		d.Set("owner", resource.Owner.TeamId.Id)
	} else {
		d.Set("owner", resource.Owner.GroupId.Id)
	}

	return nil
}
