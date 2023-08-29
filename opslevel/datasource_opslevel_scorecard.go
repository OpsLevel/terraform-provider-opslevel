package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceScorecard() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceScorecardRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the scorecard to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The scorecard's name.",
				ForceNew:    false,
				Required:    true,
			},
			"owner_id": {
				Type:        schema.TypeString,
				Description: "The scorecard's owner.",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The scorecard's description.",
				ForceNew:    false,
				Optional:    true,
			},
			"filter_id": {
				Type:        schema.TypeString,
				Description: "The scorecard's filter.",
				ForceNew:    false,
				Optional:    true,
			},

			// computed fields
			"aliases": {
				Type:        schema.TypeList,
				Description: "The scorecard's aliases.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"passing_checks": {
				Type:        schema.TypeInt,
				Description: "The scorecard's number of checks that are passing.",
				Computed:    true,
			},
			"service_count": {
				Type:        schema.TypeInt,
				Description: "The scorecard's number of services matched.",
				Computed:    true,
			},
			"total_checks": {
				Type:        schema.TypeInt,
				Description: "The scorecard's total number of checks.",
				Computed:    true,
			},
		},
	}
}

func datasourceScorecardRead(d *schema.ResourceData, client *opslevel.Client) error {
	identifier := d.Get("identifier").(string)
	resource, err := client.GetScorecard(identifier)
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("name", resource.Name)
	d.Set("owner_id", resource.Id)
	d.Set("description", resource.Description)
	d.Set("filter_id", resource.Filter.FilterId)
	d.Set("aliases", resource.Aliases)
	d.Set("passing_checks", resource.PassingChecks)
	d.Set("service_count", resource.ServiceCount)
	d.Set("total_checks", resource.ChecksCount)

	return nil
}
