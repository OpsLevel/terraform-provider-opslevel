package opslevel

import (
	"log"

	"github.com/iancoleman/strcase"
	"github.com/shurcooL/graphql"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kr/pretty"
	"github.com/opslevel/opslevel-go"
)

func resourceOpsLevelService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsLevelServiceCreate,
		Read:   resourceOpsLevelServiceRead,
		Update: resourceOpsLevelServiceUpdate,
		Delete: resourceOpsLevelServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
		},
	}
}

func resourceOpsLevelServiceCreate(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)

	svcCreate := opslevel.ServiceCreateInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Framework:   d.Get("framework").(string),
		Language:    d.Get("language").(string),
		Product:     d.Get("product").(string),
	}

	svc, err := p.client.CreateService(svcCreate)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] created OpsLevel service ID:", svc.Id)
	d.SetId(svc.Id.(string))

	// add Aliases
	aliases := expandServiceAliases(d.Get("aliases").([]interface{}))
	if len(aliases) > 0 {
		p.client.CreateAliases(svc.Id, aliases)
	}

	// add tags
	if v, ok := d.Get("tags").(map[string]interface{}); ok && len(v) > 0 {
		p.client.AssignTagsForId(svc.Id, expandStringMap(d.Get("tags").(map[string]interface{})))
	}

	return resourceOpsLevelServiceRead(d, meta)
}

func resourceOpsLevelServiceRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)
	id := d.Id()

	log.Println("[DEBUG] querying OpsLevel for service with ID:", id)

	svc, err := p.client.GetServiceWithId(id)
	if err != nil {
		log.Printf("[DEBUG] query service error: %s", pretty.Sprint(err))

		return err
	}
	if svc.Id == nil || svc.Id.(string) == "" {
		log.Println("[DEBUG] service no longer exists")
		d.SetId("")
		return nil
	}

	flattened := flattenService(svc)
	for k, v := range flattened {
		d.Set(k, v)
	}

	return nil
}

func resourceOpsLevelServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	//p := meta.(provider)
	//id := d.Id()

	return resourceOpsLevelServiceRead(d, meta)
}

func resourceOpsLevelServiceDelete(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)
	id := d.Id()

	delInput := opslevel.ServiceDeleteInput{
		Id: graphql.ID(id),
	}

	err := p.client.DeleteService(delInput)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
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

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func flattenService(svc *opslevel.Service) map[string]interface{} {
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
