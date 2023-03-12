package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

// Handles conversion from Terraform's interface{} struct to OpsLevel's JSON struct
func expandHeaders(headers interface{}) opslevel.JSON {
	output := opslevel.JSON{}
	for k, v := range headers.(map[string]interface{}) {
		output[k] = v.(string)
	}
	return output
}

func resourceWebhookAction() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a webhook action",
		Create:      wrap(resourceWebhookActionCreate),
		Read:        wrap(resourceWebhookActionRead),
		Update:      wrap(resourceWebhookActionUpdate),
		Delete:      wrap(resourceWebhookActionDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Webhook Action.",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Webhook Action.",
				ForceNew:    false,
				Optional:    true,
			},
			"payload": {
				Type:        schema.TypeString,
				Description: "Template that can be used to generate a webhook payload.",
				ForceNew:    false,
				Required:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The URL of the Webhook Action.",
				ForceNew:    false,
				Required:    true,
			},
			"method": {
				Type:         schema.TypeString,
				Description:  "The http method used to call the Webhook Action.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllCustomActionsHttpMethodEnum, false),
			},
			"headers": {
				Type:        schema.TypeMap,
				Description: "HTTP headers to be passed along with your webhook when triggered.",
				ForceNew:    false,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceWebhookActionCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CustomActionsWebhookActionCreateInput{
		Name:           d.Get("name").(string),
		LiquidTemplate: d.Get("payload").(string),
		WebhookURL:     d.Get("url").(string),
		HTTPMethod:     opslevel.CustomActionsHttpMethodEnum(d.Get("method").(string)),
		Headers:        expandHeaders(d.Get("headers")),
	}

	if _, ok := d.GetOk("description"); ok {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}

	resource, err := client.CreateWebhookAction(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceWebhookActionRead(d, client)
}

func resourceWebhookActionRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCustomAction(*opslevel.NewIdentifier(id))
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}

	if err := d.Set("description", resource.Description); err != nil {
		return err
	}

	if err := d.Set("payload", resource.LiquidTemplate); err != nil {
		return err
	}

	if err := d.Set("url", resource.WebhookURL); err != nil {
		return err
	}

	if err := d.Set("method", string(resource.HTTPMethod)); err != nil {
		return err
	}

	if err := d.Set("headers", resource.Headers); err != nil {
		return err
	}

	return nil
}

func resourceWebhookActionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CustomActionsWebhookActionUpdateInput{
		Id: *opslevel.NewID(d.Id()),
	}

	if d.HasChange("name") {
		input.Name = opslevel.NewString(d.Get("name").(string))
	}
	if d.HasChange("description") {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}
	if d.HasChange("payload") {
		input.WebhookURL = opslevel.NewString(d.Get("payload").(string))
	}
	if d.HasChange("url") {
		input.WebhookURL = opslevel.NewString(d.Get("url").(string))
	}
	if d.HasChange("method") {
		input.HTTPMethod = opslevel.CustomActionsHttpMethodEnum(d.Get("method").(string))
	}
	if d.HasChange("headers") {
		headers := expandHeaders(d.Get("headers"))
		input.Headers = &headers
	}

	_, err := client.UpdateWebhookAction(input)
	if err != nil {
		return err
	}

	return resourceWebhookActionRead(d, client)
}

func resourceWebhookActionDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteWebhookAction(*opslevel.NewIdentifier(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
