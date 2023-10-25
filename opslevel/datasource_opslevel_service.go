package opslevel

import (
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceService() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceServiceRead),
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Description: "An alias of the service to find by.",
				ForceNew:    true,
				Optional:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the service to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the service.",
				Computed:    true,
			},
			"product": {
				Type:        schema.TypeString,
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A brief description of the service.",
				Computed:    true,
			},
			"language": {
				Type:        schema.TypeString,
				Description: "The primary programming language that the service is written in.",
				Computed:    true,
			},
			"framework": {
				Type:        schema.TypeString,
				Description: "The primary software development framework that the service uses.",
				Computed:    true,
			},
			"tier_alias": {
				Type:        schema.TypeString,
				Description: "The software tier that the service belongs to.",
				Computed:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The team that owns the service.",
				Computed:    true,
			},
			"owner_alias": {
				Type:        schema.TypeString,
				Description: "The team that owns the service.",
				Computed:    true,
				Deprecated:  "field 'owner_alias' on service is no longer supported please use the 'owner' field.",
			},
			"owner_id": {
				Type:        schema.TypeString,
				Description: "The team ID that owns the service.",
				Computed:    true,
			},
			"lifecycle_alias": {
				Type:        schema.TypeString,
				Description: "The lifecycle stage of the service.",
				Computed:    true,
			},
			"api_document_path": {
				Type:        schema.TypeString,
				Description: "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.",
				Computed:    true,
			},
			"preferred_api_document_source": {
				Type:        schema.TypeString,
				Description: "The API document source (PUSH or PULL) used to determine the displayed document. If null, we use the order push and then pull.",
				Computed:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "A list of human-friendly, unique identifiers for the service",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "A list of tags applied to the service.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceServiceRead(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := findService("alias", "id", d, client)
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("product", resource.Product); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("language", resource.Language); err != nil {
		return err
	}
	if err := d.Set("framework", resource.Framework); err != nil {
		return err
	}
	if err := d.Set("tier_alias", resource.Tier.Alias); err != nil {
		return err
	}
	if err := d.Set("owner", resource.Owner.Alias); err != nil {
		return err
	}
	if err := d.Set("owner_alias", resource.Owner.Alias); err != nil {
		return err
	}
	if err := d.Set("owner_id", resource.Owner.Id); err != nil {
		return err
	}
	if err := d.Set("lifecycle_alias", resource.Lifecycle.Alias); err != nil {
		return err
	}

	if err := d.Set("aliases", resource.Aliases); err != nil {
		return err
	}
	if err := d.Set("tags", flattenTagArray(resource.Tags.Nodes)); err != nil {
		return err
	}

	if err := d.Set("api_document_path", resource.ApiDocumentPath); err != nil {
		return err
	}
	if err := d.Set("preferred_api_document_source", resource.PreferredApiDocumentSource); err != nil {
		return err
	}

	return nil
}
