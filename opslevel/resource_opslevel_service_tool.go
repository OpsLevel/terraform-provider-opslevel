package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
)

func resourceServiceTool() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service tool",
		Create:      wrap(resourceServiceToolCreate),
		Read:        wrap(resourceServiceToolRead),
		Update:      wrap(resourceServiceToolUpdate),
		Delete:      wrap(resourceServiceToolDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service": {
				Type:        schema.TypeString,
				Description: "The id of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"service_alias": {
				Type:        schema.TypeString,
				Description: "The alias of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the tool.",
				ForceNew:    false,
				Required:    true,
			},
			"category": {
				Type:         schema.TypeString,
				Description:  "The category that the tool belongs to.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllToolCategory(), false),
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The URL of the tool.",
				ForceNew:    false,
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Description: "The environment that the tool belongs to.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceServiceToolCreate(d *schema.ResourceData, client *opslevel.Client) error {
	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}

	input := opslevel.ToolCreateInput{
		ServiceId: service.Id,

		DisplayName: d.Get("name").(string),
		Category:    opslevel.ToolCategory(d.Get("category").(string)),
		Url:         d.Get("url").(string),
	}
	if env := d.Get("environment"); env != nil {
		input.Environment = env.(string)
	}
	resource, err := client.CreateTool(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("category", string(resource.Category)); err != nil {
		return err
	}
	if err := d.Set("url", resource.Url); err != nil {
		return err
	}
	if err := d.Set("environment", resource.Environment); err != nil {
		return err
	}

	return nil
}

func resourceServiceToolRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	// Handle Import by spliting the ID into the 2 parts
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		d.Set("service", parts[0])
		id = parts[1]
		d.SetId(id)
	}

	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}

	var resource *opslevel.Tool
	for _, t := range service.Tools.Nodes {
		if t.Id == id {
			resource = &t
			break
		}
	}
	if resource == nil {
		return fmt.Errorf("unable to find tool with id '%s' on service '%s'", id, service.Aliases[0])
	}

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("category", string(resource.Category)); err != nil {
		return err
	}
	if err := d.Set("url", resource.Url); err != nil {
		return err
	}
	if err := d.Set("environment", resource.Environment); err != nil {
		return err
	}

	return nil
}

func resourceServiceToolUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ToolUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.DisplayName = d.Get("name").(string)
	}
	if d.HasChange("category") {
		input.Category = opslevel.ToolCategory(d.Get("category").(string))
	}
	if d.HasChange("url") {
		input.Url = d.Get("url").(string)
	}
	if d.HasChange("environment") {
		input.Environment = d.Get("environment").(string)
	}

	resource, err := client.UpdateTool(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("category", string(resource.Category)); err != nil {
		return err
	}
	if err := d.Set("url", resource.Url); err != nil {
		return err
	}
	if err := d.Set("environment", resource.Environment); err != nil {
		return err
	}
	return nil
}

func resourceServiceToolDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTool(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
