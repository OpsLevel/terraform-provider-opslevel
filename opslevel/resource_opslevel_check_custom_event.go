package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
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

	input.Integration = *opslevel.NewID(d.Get("integration").(string))
	if passPending, ok := d.GetOk("pass_pending"); ok {
		input.PassPending = opslevel.Bool(passPending.(bool))
	}
	input.ServiceSelector = d.Get("service_selector").(string)
	input.SuccessCondition = d.Get("success_condition").(string)
	input.Message = d.Get("message").(string)

	resource, err := client.CreateCheckCustomEvent(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckCustomEventRead(d, client)
}

func resourceCheckCustomEventRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	if err := d.Set("integration", resource.Integration.Id); err != nil {
		return err
	}

	if _, ok := d.GetOk("pass_pending"); ok {
		if err := d.Set("pass_pending", resource.PassPending); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("service_selector"); ok {
		if err := d.Set("service_selector", resource.ServiceSelector); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("success_condition"); ok {
		if err := d.Set("success_condition", resource.SuccessCondition); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("message"); ok {
		if err := d.Set("message", resource.ResultMessage); err != nil {
			return err
		}
	}

	return nil
}

func resourceCheckCustomEventUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckCustomEventUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("integration") {
		input.Integration = opslevel.NewID(d.Get("integration").(string))
	}
	if d.HasChange("pass_pending") {
		if passPending, ok := d.GetOk("pass_pending"); ok {
			input.PassPending = opslevel.Bool(passPending.(bool))
		}
	}
	if d.HasChange("service_selector") {
		serviceSelector := d.Get("service_selector").(string)
		input.ServiceSelector = &serviceSelector
	}
	if d.HasChange("success_condition") {
		successCondition := d.Get("success_condition").(string)
		input.SuccessCondition = &successCondition
	}
	if d.HasChange("message") {
		message := d.Get("message").(string)
		input.Message = &message
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
