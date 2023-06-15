package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceInfrastructure() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an infrastructure resource",
		Create:      wrap(resourceInfrastructureCreate),
		Read:        wrap(resourceInfrastructureRead),
		Update:      wrap(resourceInfrastructureUpdate),
		Delete:      wrap(resourceInfrastructureDelete),
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the infrastructure resource in its provider.",
				Optional:    true,
				ForceNew:    false,
			},
			"schema": {
				Type:        schema.TypeString,
				Description: "The schema of the infrastructure resource that determines its specification.",
				Required:    true,
				ForceNew:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The id of the owner for the infrastructure resource.  Can be a team or group",
				ForceNew:    false,
				Optional:    true,
			},
			"provider_data": {
				Type:        schema.TypeList,
				Description: "The provider data for the infrastructure resource",
				ForceNew:    false,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Description: "The account name for the infrastructure resource.",
							ForceNew:    false,
							Optional:    true,
						},
						"external_url": {
							Type:        schema.TypeString,
							Description: "The external url for the infrastructure resource.",
							ForceNew:    false,
							Optional:    true,
						},
						"provider_name": {
							Type:        schema.TypeString,
							Description: "The provider name for the infrastructure resource.",
							ForceNew:    false,
							Optional:    true,
						},
					},
				},
			},
			"data": {
				Type:        schema.TypeString,
				Description: "The data of the infrastructure resource in JSON format.",
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "The aliases of the infrastructure resource.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func expandInfrastructureProviderData(d *schema.ResourceData) *opslevel.InfrastructureResourceProviderInput {
	output := opslevel.InfrastructureResourceProviderInput{}
	config := d.Get("provider_data").([]interface{})
	if len(config) > 0 {
		item := config[0].(map[string]interface{})
		output.AccountName = item["account_name"].(string)
		output.ExternalURL = item["external_url"].(string)
		output.ProviderName = item["provider_name"].(string)
	}
	return &output
}

func resourceInfrastructureCreate(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := client.CreateInfrastructure(opslevel.InfrastructureResourceInput{
		Type: GetString(d, "type"),
		Schema: &opslevel.InfrastructureResourceSchemaInput{
			Type: *GetString(d, "schema"),
		},
		ProviderData: expandInfrastructureProviderData(d),
		Owner:        opslevel.NewID(d.Get("owner").(string)),
		Data:         opslevel.NewJSON(*GetString(d, "data")),
	})
	if err != nil {
		return err
	}
	d.SetId(resource.Id)
	return resourceInfrastructureRead(d, client)
}

func resourceInfrastructureRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetInfrastructure(id)
	if err != nil {
		return err
	}

	if err := d.Set("aliases", resource.Aliases); err != nil {
		return err
	}
	if err := d.Set("type", resource.Type); err != nil {
		return err
	}
	// TODO: check if this is not needed - but i fear Terraform might need this and our API doesn expose it.
	//if err := d.Set("schema", resource.Schema.Type); err != nil {
	//	return err
	//}
	if err := d.Set("owner", resource.Owner.Id()); err != nil {
		return err
	}

	// TODO: this might not work properly - might need to use resource.Data.MarshalJSON()
	if err := d.Set("data", resource.Data); err != nil {
		return err
	}

	return nil
}

func resourceInfrastructureUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	// TODO: i'm using a different pattern here then in other update functions, but i think this is actually better
	// please make sure to test the update function heavily so we don't break anything
	// make sure that each individual field change changed and doesn't cause an infinite plan change
	input := opslevel.InfrastructureResourceInput{
		Type: GetString(d, "type"),
		Schema: &opslevel.InfrastructureResourceSchemaInput{
			Type: *GetString(d, "schema"),
		},
		ProviderData: expandInfrastructureProviderData(d),
		Owner:        opslevel.NewID(d.Get("owner").(string)),
		Data:         opslevel.NewJSON(*GetString(d, "data")),
	}

	_, err := client.UpdateInfrastructure(id, input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceSystemRead(d, client)
}

func resourceInfrastructureDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteInfrastructure(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
