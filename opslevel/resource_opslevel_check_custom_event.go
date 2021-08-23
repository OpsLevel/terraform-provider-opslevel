package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceCheckCustomEvent() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a custom event check.",
		Create:      wrap(resourceCheckCustomEventCreate),
		Read:        wrap(resourceCheckCustomEventRead),
		Update:      wrap(resourceCheckCustomEventUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"integration": {
				Type:        schema.TypeString,
				Description: "The integration id this check will use.",
				ForceNew:    false,
				Required:    true,
			},
			"service_selector": {
				Type:        schema.TypeString,
				Description: "A jq expression that will be ran against your payload. This will parse out the service identifier.",
				ForceNew:    false,
				Required:    true,
			},
			"success_condition": {
				Type:        schema.TypeString,
				Description: "A jq expression that will be ran against your payload. A truthy value will result in the check passing.",
				ForceNew:    false,
				Required:    true,
			},
			"message": {
				Type:        schema.TypeString,
				Description: "The check result message template. It is compiled with Liquid and formatted in Markdown.",
				ForceNew:    false,
				Optional:    true,
			},
		}),
	}
}

func resourceCheckCustomEventCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckCustomEventCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		Integration:      getID(d, "integration"),
		ServiceSelector:  d.Get("service_selector").(string),
		SuccessCondition: d.Get("success_condition").(string),
		Message:          d.Get("message").(string),
	}
	resource, err := client.CreateCheckCustomEvent(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckCustomEventRead(d, client)
}

func resourceCheckCustomEventRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
	if err != nil {
		return err
	}

	if err := resourceCheckRead(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckCustomEventUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckCustomEventUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("enabled") {
		value := d.Get("enabled").(bool)
		input.Enabled = &value
	}
	if d.HasChange("category") {
		input.Category = getID(d, "category")
	}
	if d.HasChange("level") {
		input.Level = getID(d, "level")
	}
	if d.HasChange("owner") {
		input.Owner = getID(d, "owner")
	}
	if d.HasChange("filter") {
		input.Filter = getID(d, "filter")
	}
	if d.HasChange("notes") {
		input.Notes = d.Get("notes").(string)
	}

	if d.HasChange("integration") {
		input.Integration = getID(d, "integration")
	}
	if d.HasChange("service_selector") {
		input.ServiceSelector = d.Get("service_selector").(string)
	}
	if d.HasChange("success_condition") {
		input.SuccessCondition = d.Get("success_condition").(string)
	}
	if d.HasChange("message") {
		input.Message = d.Get("message").(string)
	}

	_, err := client.UpdateCheckCustomEvent(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckCustomEventRead(d, client)
}
