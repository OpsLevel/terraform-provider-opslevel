package opslevel

import (
	"fmt"
	"log"
	"strings"

	"github.com/kr/pretty"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceOpsLevelServices() *schema.Resource {
	dsServiceSchema := datasourceSchemaFromResourceSchema(serviceSchema)

	dsServiceSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The service ID",
		Computed:    true,
	}

	return &schema.Resource{
		Read: datasourceOpsLevelServicesRead,
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
							Description: "The service field to filter on. Accepts `alias`, `id`, `framework`, `language`, `ownerAlias`, `tag`.",
							ForceNew:    true,
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The service field value to match.",
							ForceNew:    true,
							Optional:    true,
						},
					},
				},
			},
			"services": {
				Type:        schema.TypeList,
				Description: "The services matching specified filters.",
				ForceNew:    false,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dsServiceSchema,
				},
			},
		},
	}
}

func datasourceOpsLevelServicesRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)

	services := []interface{}{}
	var err error

	field := d.Get("filter.0.field").(string)
	value := d.Get("filter.0.value").(string)
	log.Printf("[DEBUG] filtering for services %s=%s", field, value)
	d.SetId(fmt.Sprintf("services(%s=%s)", field, value))

	switch field {
	case "alias":
		svc, err := GetServiceWithAlias(p.client, value)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] got service: %v", svc)

		services = append(services, flattenService(svc))
	case "id":
		svc, err := GetServiceWithId(p.client, value)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] got service: %v", svc)

		services = append(services, flattenService(svc))
	default:
		var response []Service
		switch field {
		case "framework":
			response, err = ListServicesByFramework(p.client, value)
		case "language":
			response, err = ListServicesByLanguage(p.client, value)
		case "ownerAlias":
			response, err = ListServicesByOwnerAlias(p.client, value)
		case "tag":
			tagKV := strings.Split(value, ":")
			if len(tagKV) != 2 {
				return fmt.Errorf("tag filter requires `value` in format 'key:value'")
			}
			response, err = ListServicesByTag(p.client, tagKV[0], tagKV[1])
		}

		if err != nil {
			return err
		}

		for _, svc := range response {
			services = append(services, flattenService(&svc))
		}
	}
	log.Printf("[DEBUG] services len=%d cap=%d %s", len(services), cap(services), pretty.Sprint(services))
	d.Set("services", services)

	return nil
}
