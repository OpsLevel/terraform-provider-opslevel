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
			"pass_pending": {
				Type:        schema.TypeBool,
				Description: "True if this check should pass by default. Otherwise the default 'pending' state counts as a failure.",
				ForceNew:    false,
				Optional:    true,
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
	input := opslevel.CheckCustomEventCreateInput{}
	setCheckCreateInput(d, &input)

	input.Integration = getID(d, "integration")
	input.PassPending = d.Get("pass_pending").(bool)
	input.ServiceSelector = d.Get("service_selector").(string)
	input.SuccessCondition = d.Get("success_condition").(string)
	input.Message = d.Get("message").(string)

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

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	if err := d.Set("integration", resource.Integration.Id.(string)); err != nil {
		return err
	}
	if err := d.Set("pass_pending", resource.PassPending); err != nil {
		return err
	}
	if err := d.Set("service_selector", resource.ServiceSelector); err != nil {
		return err
	}
	if err := d.Set("success_condition", resource.SuccessCondition); err != nil {
		return err
	}
	if err := d.Set("message", resource.ResultMessage); err != nil {
		return err
	}

	return nil
}

func resourceCheckCustomEventUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckCustomEventUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("integration") {
		input.Integration = getID(d, "integration")
	}
	if d.HasChange("pass_pending") {
		input.PassPending = d.Get("pass_pending").(bool)
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
	if err := d.Set("last_updated", timeLastUpdated()); err != nil {
		return err
	}
	return resourceCheckCustomEventRead(d, client)
}
