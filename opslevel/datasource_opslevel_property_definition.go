package opslevel

import (
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourcePropertyDefinition() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourcePropertyDefinitionRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the property definition to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the property definition.",
				Computed:    true,
			},
			"schema": {
				Type:        schema.TypeString,
				Description: "The schema of the property definition.",
				Computed:    true,
			},
		},
	}
}

func datasourcePropertyDefinitionRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Get("id").(string)
	resource, err := client.GetPropertyDefinition(id)
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("schema", resource.Schema.ToJSON()); err != nil {
		return err
	}

	return nil
}
