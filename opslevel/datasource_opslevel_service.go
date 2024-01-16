package opslevel

import (
	"strings"

	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rs/zerolog/log"

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
			"repositories": {
				Type:        schema.TypeList,
				Description: "List of repositories connected to the service.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"properties": {
				Type:        schema.TypeList,
				Description: "Custom properties assigned to this service.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Description: "The custom property definition's ID.",
							Computed:    true,
						},
						"owner": {
							Type:        schema.TypeString,
							Description: "The ID of the entity that the property has been assigned to.",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the custom property.",
							Computed:    true,
						},
					},
				},
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

	if err := d.Set("aliases", resource.ManagedAliases); err != nil {
		return err
	}
	if err := d.Set("tags", flattenTagArray(resource.Tags.Nodes)); err != nil {
		return err
	}
	if err := d.Set("repositories", flattenServiceRepositoriesArray(resource.Repositories)); err != nil {
		return err
	}

	properties, err := resource.GetProperties(client, nil)
	if err != nil {
		return err
	}
	// log warnings for any validation errors rather than adding them to state
	for _, property := range properties.Nodes {
		for _, validationErr := range property.ValidationErrors {
			log.Warn().Msgf("service '%s' property '%s' has a validation error\n\tmessage=\"%s\" path=[%s]",
				string(resource.Id),
				string(property.Definition.Id),
				validationErr.Message,
				strings.Join(validationErr.Path, ","))
		}
	}
	props := mapServiceProperties(properties)
	if err := d.Set("properties", props); err != nil {
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
