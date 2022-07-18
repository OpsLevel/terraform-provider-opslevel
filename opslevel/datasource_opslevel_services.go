package opslevel

import (
	"github.com/opslevel/opslevel-go/v2022"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceServices() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceServicesRead),
		Schema: map[string]*schema.Schema{
			"filter": getDatasourceFilter(false, []string{"framework", "language", "lifecycle", "owner", "product", "tag", "tier"}),
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
			"urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceServicesRead(d *schema.ResourceData, client *opslevel.Client) error {
	field := d.Get("filter.0.field").(string)
	value := d.Get("filter.0.value").(string)

	var services []opslevel.Service
	var err error

	switch field {
	case "framework":
		services, err = client.ListServicesWithFramework(value)
	case "language":
		services, err = client.ListServicesWithLanguage(value)
	case "lifecycle":
		services, err = client.ListServicesWithLifecycle(value)
	case "owner":
		services, err = client.ListServicesWithOwner(value)
	case "product":
		services, err = client.ListServicesWithProduct(value)
	case "tag":
		services, err = client.ListServicesWithTag(opslevel.NewTagArgs(value))
	case "tier":
		services, err = client.ListServicesWithTier(value)
	default:
		services, err = client.ListServices()
	}
	if err != nil {
		return err
	}

	count := len(services)
	ids := make([]string, count)
	names := make([]string, count)
	urls := make([]string, count)
	for i, item := range services {
		ids[i] = item.Id.(string)
		names[i] = item.Name
		urls[i] = item.HtmlURL
	}

	d.SetId(timeID())
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("urls", urls)

	return nil
}
