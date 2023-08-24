package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceScorecard() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a scorecard",
		Create:      wrap(resourceScorecardCreate),
		Read:        wrap(resourceScorecardRead),
		Update:      wrap(resourceScorecardUpdate),
		Delete:      wrap(resourceScorecardDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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

func resourceScorecardCreate(d *schema.ResourceData, client *opslevel.Client) error {
	ownerId := opslevel.NewID(d.Get("owner_id").(string))
	filterId := opslevel.NewID(d.Get("filter_id").(string))

	input := opslevel.ScorecardInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(*string),
		OwnerId:     *ownerId,
		FilterId:    filterId,
	}

	resource, err := client.CreateScorecard(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceScorecardRead(d, client)
}

func resourceScorecardRead(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := client.GetScorecard(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("aliases", resource.Aliases); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("filter", resource.Filter); err != nil {
		return err
	}
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("owner", resource.Owner); err != nil {
		return err
	}
	if err := d.Set("passing_checks", resource.PassingChecks); err != nil {
		return err
	}
	if err := d.Set("service_count", resource.ServiceCount); err != nil {
		return err
	}
	if err := d.Set("total_checks", resource.ChecksCount); err != nil {
		return err
	}

	return nil
}

func resourceScorecardUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	ownerId := opslevel.NewID(d.Get("owner_id").(string))
	filterId := opslevel.NewID(d.Get("filter_id").(string))

	input := opslevel.ScorecardInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(*string),
		OwnerId:     *ownerId,
		FilterId:    filterId,
	}

	_, err := client.UpdateScorecard(d.Id(), input)
	if err != nil {
		return err
	}

	return resourceScorecardRead(d, client)
}

func resourceScorecardDelete(d *schema.ResourceData, client *opslevel.Client) error {
	_, err := client.DeleteScorecard(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
