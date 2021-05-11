package opslevel

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kr/pretty"
	"github.com/zapier/opslevel-go"
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
				Computed:    true,
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
			"product": {
				Type:        schema.TypeString,
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				ForceNew:    false,
				Optional:    true,
			},
			"tier": {
				Type:        schema.TypeString,
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceOpsLevelServiceCreate(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)

	serviceData := opslevel.Service{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Framework:   d.Get("framework").(string),
		Language:    d.Get("language").(string),
		Product:     d.Get("product").(string),
	}

	resp, err := p.client.CreateService(context.Background(), serviceData)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] created OpsLevel service ID:", resp.Id)
	d.SetId(resp.Id)

	return resourceOpsLevelServiceRead(d, meta)
}

func resourceOpsLevelServiceRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)
	id := d.Id()

	log.Println("[DEBUG] querying OpsLevel for service with ID:", id)

	svc, err := p.client.GetServiceById(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no service was found") {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] query service error: %s", pretty.Sprint(err))

		return err
	}

	d.Set("name", svc.Name)
	d.Set("description", svc.Description)
	d.Set("framework", svc.Framework)
	d.Set("language", svc.Language)
	d.Set("owner_id", svc.Owner.Id)
	d.Set("product", svc.Product)
	d.Set("tier_id", svc.Tier.Id)

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

	_, err := p.client.DeleteServiceById(context.Background(), id)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
