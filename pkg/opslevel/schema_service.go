package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/iancoleman/strcase"
	"log"
)

var serviceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Description: "The display name of the service.",
		ForceNew:    false,
		Required:    true,
	},
	"aliases": {
		Type:        schema.TypeList,
		Description: "A list of human-friendly, unique identifiers for the service",
		ForceNew:    false,
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"description": {
		Type:        schema.TypeString,
		Description: "A brief description of the service.",
		ForceNew:    false,
		Optional:    true,
	},
	"framework": {
		Type:        schema.TypeString,
		Description: "The primary software development framework that the service uses.",
		ForceNew:    false,
		Optional:    true,
	},
	"language": {
		Type:        schema.TypeString,
		Description: "The primary programming language that the service is written in.",
		ForceNew:    false,
		Optional:    true,
	},
	"owner_id": {
		Type:        schema.TypeString,
		Description: "The ID of the Team that owns this service.",
		ForceNew:    false,
		Optional:    true,
	},
	"product": {
		Type:        schema.TypeString,
		Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
		ForceNew:    false,
		Optional:    true,
	},
	"tags": {
		Type:        schema.TypeMap,
		Description: "A map of tags applied to the service",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"tier_id": {
		Type:        schema.TypeString,
		Description: "The ID of the software tier that the service belongs to.",
		ForceNew:    false,
		Optional:    true,
	},
}

func flattenService(svc *Service) map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = svc.Id.(string)
	m["name"] = string(svc.Name)
	m["description"] = string(svc.Description)
	m["framework"] = string(svc.Framework)
	m["language"] = string(svc.Language)

	if svc.Owner.Id != nil {
		m["owner_id"] = svc.Owner.Id.(string)
	}

	m["product"] = string(svc.Product)
	m["tier_id"] = string(svc.Tier.Id)

	aliases := []string{}
	for _, alias := range svc.Aliases {
		str := string(alias)
		if str == strcase.ToSnake(string(svc.Name)) {
			log.Printf("[DEBUG] ignoring alias `%s`", str)
			continue
		}
		aliases = append(aliases, str)
	}
	m["aliases"] = aliases

	tags := map[string]string{}
	for _, tag := range svc.Tags.Nodes {
		tags[string(tag.Key)] = string(tag.Value)
	}
	m["tags"] = tags

	return m
}

func expandServiceAliases(cfg []interface{}) []string {
	aliases := make([]string, len(cfg))
	if len(cfg) > 0 {
		for i, v := range cfg {
			aliases[i] = v.(string)
		}
	}

	return aliases
}
