package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func datasourceScorecards() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceScorecardsRead),
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owner_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// computed fields
			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeList},
			},
			"passing_checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"service_counts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func datasourceScorecardsRead(d *schema.ResourceData, client *opslevel.Client) error {
	resp, err := client.ListScorecards(nil)
	if err != nil {
		return err
	}

	count := len(resp.Nodes)
	ids := make([]string, count)
	names := make([]string, count)
	ownerIds := make([]string, count)
	descriptions := make([]string, count)
	filterIds := make([]string, count)
	aliases := make([]string, count)
	passingChecks := make([]int, count)
	serviceCounts := make([]int, count)
	totalChecks := make([]int, count)
	for i, item := range resp.Nodes {
		if len(item.Aliases) > 0 {
			aliases[i] = item.Aliases[0]
		}
		ids[i] = string(item.Id)
		names[i] = item.Name
		ownerIds[i] = string(item.Owner.Id())
		descriptions[i] = item.Description
		filterIds[i] = string(item.Filter.Id)
		passingChecks[i] = item.PassingChecks
		serviceCounts[i] = item.ServiceCount
		totalChecks[i] = item.ChecksCount
	}

	d.SetId(timeID())
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("owner_ids", ownerIds)
	d.Set("descriptions", descriptions)
	d.Set("filter_ids", filterIds)
	d.Set("aliases", aliases)
	d.Set("passing_checks", passingChecks)
	d.Set("service_counts", serviceCounts)
	d.Set("total_checks", totalChecks)

	return nil
}
