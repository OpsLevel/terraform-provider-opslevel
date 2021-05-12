package opslevel

import (
	"fmt"
	"log"
	"strings"

	"github.com/kr/pretty"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceOpsLevelService() *schema.Resource {
	serviceSchema := datasourceSchemaFromResourceSchema(resourceOpsLevelService().Schema)

	serviceSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The service ID",
		Computed:    true,
	}

	return &schema.Resource{
		Read: datasourceOpsLevelServiceRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeList,
				Description: "The services matching filter.",
				ForceNew:    true,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:        schema.TypeString,
							Description: "The service field to filter on.",
							ForceNew:    true,
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The field value to match.",
							ForceNew:    true,
							Optional:    true,
						},
					},
				},
			},
			"services": {
				Type:        schema.TypeList,
				Description: "The services matching filter.",
				ForceNew:    false,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: serviceSchema,
				},
			},
		},
	}
}

func datasourceOpsLevelServiceRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)

	services := []interface{}{}

	field := d.Get("filter.0.field").(string)
	value := d.Get("filter.0.value").(string)
	log.Printf("[DEBUG] filtering for services %s=%s", field, value)
	d.SetId(fmt.Sprintf("services(%s=%s)", field, value))

	if field == "alias" {
		svc, err := p.client.GetServiceWithAlias(value)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] got service: %v", svc)

		services = append(services, flattenService(svc))

	} else if field == "id" {
		svc, err := p.client.GetServiceWithId(value)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] got service: %v", svc)

		services = append(services, flattenService(svc))

	} else {
		allServices, err := p.client.ListServices()
		log.Printf("[DEBUG] got %d services", len(allServices))
		if err != nil {
			return err
		}
		for _, svc := range allServices {
			smap := flattenService(&svc)
			if strings.ToLower(smap[field].(string)) == strings.ToLower(value) {
				services = append(services, smap)
			}
		}
	}

	log.Printf("[DEBUG] services len=%d cap=%d %s", len(services), cap(services), pretty.Sprint(services))
	d.Set("services", services)

	return nil
}
