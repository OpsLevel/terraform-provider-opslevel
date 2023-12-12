package opslevel

import (
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourcePropertyDefinition() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourcePropertyDefinitionRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the property definition to find.",
				Required:    true,
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
	identifier := d.Get("identifier").(string)
	resource, err := client.GetPropertyDefinition(identifier)
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
