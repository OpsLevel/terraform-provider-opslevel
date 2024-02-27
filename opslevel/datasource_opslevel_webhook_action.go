package opslevel

// import (
// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceWebhookAction() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceWebhookActionRead),
// 		Schema: map[string]*schema.Schema{
// 			"identifier": {
// 				Type:        schema.TypeString,
// 				Description: "The id or alias of the webhook action to find.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The name of the Webhook Action.",
// 				Computed:    true,
// 			},
// 			"description": {
// 				Type:        schema.TypeString,
// 				Description: "The description of the Webhook Action.",
// 				Computed:    true,
// 			},
// 			"payload": {
// 				Type:        schema.TypeString,
// 				Description: "Template that can be used to generate a webhook payload.",
// 				Computed:    true,
// 			},
// 			"url": {
// 				Type:        schema.TypeString,
// 				Description: "The URL of the Webhook Action.",
// 				Computed:    true,
// 			},
// 			"method": {
// 				Type:        schema.TypeString,
// 				Description: "The http method used to call the Webhook Action.",
// 				Computed:    true,
// 			},
// 			"headers": {
// 				Type:        schema.TypeMap,
// 				Description: "HTTP headers to be passed along with your webhook when triggered.",
// 				Computed:    true,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 		},
// 	}
// }

// func datasourceWebhookActionRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	identifier := d.Get("identifier").(string)
// 	resource, err := client.GetCustomAction(identifier)
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId(string(resource.Id))
// 	d.Set("name", resource.Name)
// 	d.Set("description", resource.Description)
// 	d.Set("payload", resource.LiquidTemplate)
// 	d.Set("url", resource.WebhookURL)
// 	d.Set("method", resource.HTTPMethod)
// 	d.Set("headers", resource.Headers)

// 	return nil
// }
