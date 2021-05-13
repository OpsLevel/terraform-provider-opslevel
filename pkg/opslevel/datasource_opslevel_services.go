package opslevel

import (
	"fmt"
	"log"
	"strings"

	"github.com/opslevel/kubectl-opslevel/opslevel"
	"github.com/shurcooL/graphql"

	"github.com/kr/pretty"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceOpsLevelServices() *schema.Resource {
	serviceSchema := datasourceSchemaFromResourceSchema(resourceOpsLevelService().Schema)

	serviceSchema["id"] = &schema.Schema{
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
					Schema: serviceSchema,
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
		svc, err := p.client.GetServiceWithAlias(value)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] got service: %v", svc)

		services = append(services, flattenService(svc))
	case "id":
		svc, err := p.client.GetServiceWithId(value)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] got service: %v", svc)

		services = append(services, flattenService(svc))
	default:
		var response []opslevel.Service
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

// By Framework

type ListServicesByFrameworkQuery struct {
	Account struct {
		Services struct {
			Nodes    []opslevel.Service
			PageInfo opslevel.PageInfo
		} `graphql:"services(framework: $framework, after: $after, first: $first)"`
	}
}

func (q *ListServicesByFrameworkQuery) Query(client *opslevel.Client, framework string) error {
	var subQ ListServicesByFrameworkQuery
	v := opslevel.PayloadVariables{
		"framework": graphql.String(framework),
		"after":     q.Account.Services.PageInfo.End,
		"first":     graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, framework)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByFramework(client *opslevel.Client, framework string) ([]opslevel.Service, error) {
	q := ListServicesByFrameworkQuery{}
	if err := q.Query(client, framework); err != nil {
		return []opslevel.Service{}, err
	}
	return q.Account.Services.Nodes, nil
}


// By Language

type ListServicesByLanguageQuery struct {
	Account struct {
		Services struct {
			Nodes    []opslevel.Service
			PageInfo opslevel.PageInfo
		} `graphql:"services(language: $language, after: $after, first: $first)"`
	}
}

func (q *ListServicesByLanguageQuery) Query(client *opslevel.Client, language string) error {
	var subQ ListServicesByLanguageQuery
	v := opslevel.PayloadVariables{
		"language": graphql.String(language),
		"after":     q.Account.Services.PageInfo.End,
		"first":     graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, language)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByLanguage(client *opslevel.Client, framework string) ([]opslevel.Service, error) {
	q := ListServicesByLanguageQuery{}
	if err := q.Query(client, framework); err != nil {
		return []opslevel.Service{}, err
	}
	return q.Account.Services.Nodes, nil
}

// By OwnerAlias

type ListServicesByOwnerAliasQuery struct {
	Account struct {
		Services struct {
			Nodes    []opslevel.Service
			PageInfo opslevel.PageInfo
		} `graphql:"services(ownerAlias: $ownerAlias, after: $after, first: $first)"`
	}
}

func (q *ListServicesByOwnerAliasQuery) Query(client *opslevel.Client, ownerAlias string) error {
	var subQ ListServicesByOwnerAliasQuery
	v := opslevel.PayloadVariables{
		"ownerAlias": graphql.String(ownerAlias),
		"after":     q.Account.Services.PageInfo.End,
		"first":     graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, ownerAlias)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByOwnerAlias(client *opslevel.Client, framework string) ([]opslevel.Service, error) {
	q := ListServicesByOwnerAliasQuery{}
	if err := q.Query(client, framework); err != nil {
		return []opslevel.Service{}, err
	}
	return q.Account.Services.Nodes, nil
}

// By Tag

type ListServicesByTagQuery struct {
	Account struct {
		Services struct {
			Nodes    []opslevel.Service
			PageInfo opslevel.PageInfo
		} `graphql:"services(tag: {key:$key, value:$value}, after: $after, first: $first)"`
	}
}

func (q *ListServicesByTagQuery) Query(client *opslevel.Client, key, value string) error {
	var subQ ListServicesByTagQuery
	v := opslevel.PayloadVariables{
		"key":   graphql.String(key),
		"value": graphql.String(value),
		"after": q.Account.Services.PageInfo.End,
		"first": graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, key, value)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByTag(client *opslevel.Client, key, value string) ([]opslevel.Service, error) {
	q := ListServicesByTagQuery{}
	if err := q.Query(client, key, value); err != nil {
		return []opslevel.Service{}, err
	}
	return q.Account.Services.Nodes, nil
}
