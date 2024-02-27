package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func datasourceScorecard() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceScorecardRead),
// 		Schema: map[string]*schema.Schema{
// 			"identifier": {
// 				Type:        schema.TypeString,
// 				Description: "The id or alias of the scorecard to find.",
// 				Required:    true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The scorecard's name.",
// 				Computed:    true,
// 			},
// 			"owner_id": {
// 				Type:        schema.TypeString,
// 				Description: "The scorecard's owner.",
// 				Computed:    true,
// 			},
// 			"description": {
// 				Type:        schema.TypeString,
// 				Description: "The scorecard's description.",
// 				Computed:    true,
// 			},
// 			"filter_id": {
// 				Type:        schema.TypeString,
// 				Description: "The scorecard's filter.",
// 				Computed:    true,
// 			},
// 			"affects_overall_service_levels": {
// 				Type:        schema.TypeBool,
// 				Description: "Specifies whether the checks on this scorecard affect services' overall maturity level.",
// 				Computed:    true,
// 			},

// 			// computed fields
// 			"aliases": {
// 				Type:        schema.TypeList,
// 				Description: "The scorecard's aliases.",
// 				Computed:    true,
// 				Elem:        &schema.Schema{Type: schema.TypeString},
// 			},
// 			"passing_checks": {
// 				Type:        schema.TypeInt,
// 				Description: "The scorecard's number of checks that are passing.",
// 				Computed:    true,
// 			},
// 			"service_count": {
// 				Type:        schema.TypeInt,
// 				Description: "The scorecard's number of services matched.",
// 				Computed:    true,
// 			},
// 			"total_checks": {
// 				Type:        schema.TypeInt,
// 				Description: "The scorecard's total number of checks.",
// 				Computed:    true,
// 			},
// 		},
// 	}
// }

// func datasourceScorecardRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	identifier := d.Get("identifier").(string)
// 	resource, err := client.GetScorecard(identifier)
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId(string(resource.Id))
// 	d.Set("name", resource.Name)
// 	d.Set("owner_id", resource.Id)
// 	d.Set("description", resource.Description)
// 	d.Set("filter_id", resource.Filter.Id)
// 	d.Set("affects_overall_service_levels", resource.AffectsOverallServiceLevels)
// 	d.Set("aliases", resource.Aliases)
// 	d.Set("passing_checks", resource.PassingChecks)
// 	d.Set("service_count", resource.ServiceCount)
// 	d.Set("total_checks", resource.ChecksCount)

// 	return nil
// }
