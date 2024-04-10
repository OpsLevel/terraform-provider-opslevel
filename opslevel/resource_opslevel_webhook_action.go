package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &WebhookActionResource{}

var _ resource.ResourceWithImportState = &WebhookActionResource{}

func NewWebhookActionResource() resource.Resource {
	return &WebhookActionResource{}
}

// WebhookActionResource defines the resource implementation.
type WebhookActionResource struct {
	CommonResourceClient
}

// WebhookActionResourceModel describes the Webhook Action managed resource.
type WebhookActionResourceModel struct {
	Description types.String `tfsdk:"description"`
	Headers     types.Map    `tfsdk:"headers"`
	Id          types.String `tfsdk:"id"`
	Method      types.String `tfsdk:"method"`
	Name        types.String `tfsdk:"name"`
	Payload     types.String `tfsdk:"payload"`
	Url         types.String `tfsdk:"url"`
}

func NewWebhookActionResourceModel(webhookAction opslevel.CustomActionsExternalAction) WebhookActionResourceModel {
	jsonAttrs := jsonToMap(webhookAction.Headers)
	return WebhookActionResourceModel{
		Description: types.StringValue(webhookAction.Description),
		Headers:     types.MapValueMust(types.StringType, jsonAttrs),
		Id:          types.StringValue(string(webhookAction.Id)),
		Method:      types.StringValue(string(webhookAction.HTTPMethod)),
		Name:        types.StringValue(webhookAction.Name),
		Payload:     types.StringValue(webhookAction.LiquidTemplate),
		Url:         types.StringValue(webhookAction.WebhookURL),
	}
}

func (r *WebhookActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_action"
}

func (r *WebhookActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "WebhookAction resource",

		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "The description of the Webhook Action.",
				Optional:    true,
			},
			"headers": schema.MapAttribute{
				ElementType: types.StringType,
				Description: "HTTP headers to be passed along with your webhook when triggered.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the Webhook Action.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"method": schema.StringAttribute{
				Description: "The http method used to call the Webhook Action.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllCustomActionsHttpMethodEnum...),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Webhook Action.",
				Required:    true,
			},
			"payload": schema.StringAttribute{
				Description: "Template that can be used to generate a webhook payload.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the Webhook Action.",
				Required:    true,
			},
		},
	}
}

func (r *WebhookActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WebhookActionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	headersAsJson, diags := MapValueToOpslevelJson(ctx, data.Headers)
	if diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to create opslevel.JSON from 'headers': '%s'", data.Headers))
		return
	}

	webhookActionInput := opslevel.CustomActionsWebhookActionCreateInput{
		Description:    data.Description.ValueStringPointer(),
		Headers:        &headersAsJson,
		HttpMethod:     opslevel.CustomActionsHttpMethodEnum(data.Method.ValueString()),
		LiquidTemplate: data.Payload.ValueStringPointer(),
		Name:           data.Name.ValueString(),
		WebhookUrl:     data.Url.ValueString(),
	}
	webhookAction, err := r.client.CreateWebhookAction(webhookActionInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create webhook action, got error: %s", err))
		return
	}
	createdWebhookActionResourceModel := NewWebhookActionResourceModel(*webhookAction)

	tflog.Trace(ctx, "created a webhook action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdWebhookActionResourceModel)...)
}

func (r *WebhookActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WebhookActionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	webhookAction, err := r.client.GetCustomAction(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read webhookAction, got error: %s", err))
		return
	}
	readWebhookActionResourceModel := NewWebhookActionResourceModel(*webhookAction)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a webhook action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &readWebhookActionResourceModel)...)
}

func (r *WebhookActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WebhookActionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	headersAsJson, diags := MapValueToOpslevelJson(ctx, data.Headers)
	if diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to create opslevel.JSON from 'headers': '%s'", data.Headers))
		return
	}

	httpMethod := opslevel.CustomActionsHttpMethodEnum(data.Method.ValueString())
	updateWebhookActionInput := opslevel.CustomActionsWebhookActionUpdateInput{
		Description:    data.Description.ValueStringPointer(),
		Headers:        &headersAsJson,
		HttpMethod:     &httpMethod,
		Id:             opslevel.ID(data.Id.ValueString()),
		LiquidTemplate: data.Payload.ValueStringPointer(),
		Name:           data.Name.ValueStringPointer(),
		WebhookUrl:     data.Url.ValueStringPointer(),
	}
	resource, err := r.client.UpdateWebhookAction(updateWebhookActionInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update webhookAction, got error: %s", err))
		return
	}
	updatedWebhookActionResourceModel := NewWebhookActionResourceModel(*resource)

	tflog.Trace(ctx, "updated a webhook action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedWebhookActionResourceModel)...)
}

func (r *WebhookActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WebhookActionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteWebhookAction(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete webhookAction, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a webhook action resource")
}

func (r *WebhookActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// // Handles conversion from Terraform's interface{} struct to OpsLevel's JSON struct
// func expandHeaders(headers interface{}) opslevel.JSON {
// 	output := opslevel.JSON{}
// 	for k, v := range headers.(map[string]interface{}) {
// 		output[k] = v.(string)
// 	}
// 	return output
// }

// func resourceWebhookAction() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a webhook action",
// 		Create:      wrap(resourceWebhookActionCreate),
// 		Read:        wrap(resourceWebhookActionRead),
// 		Update:      wrap(resourceWebhookActionUpdate),
// 		Delete:      wrap(resourceWebhookActionDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The name of the Webhook Action.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"description": {
// 				Type:        schema.TypeString,
// 				Description: "The description of the Webhook Action.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"payload": {
// 				Type:        schema.TypeString,
// 				Description: "Template that can be used to generate a webhook payload.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"url": {
// 				Type:        schema.TypeString,
// 				Description: "The URL of the Webhook Action.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"method": {
// 				Type:         schema.TypeString,
// 				Description:  "The http method used to call the Webhook Action.",
// 				ForceNew:     false,
// 				Required:     true,
// 				ValidateFunc: validation.StringInSlice(opslevel.AllCustomActionsHttpMethodEnum, false),
// 			},
// 			"headers": {
// 				Type:        schema.TypeMap,
// 				Description: "HTTP headers to be passed along with your webhook when triggered.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 		},
// 	}
// }

// func resourceWebhookActionCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CustomActionsWebhookActionCreateInput{
// 		Name:           d.Get("name").(string),
// 		LiquidTemplate: opslevel.RefOf(d.Get("payload").(string)),
// 		WebhookUrl:     d.Get("url").(string),
// 		HttpMethod:     opslevel.CustomActionsHttpMethodEnum(d.Get("method").(string)),
// 		Headers:        opslevel.RefOf(expandHeaders(d.Get("headers"))),
// 	}

// 	if _, ok := d.GetOk("description"); ok {
// 		input.Description = opslevel.RefOf(d.Get("description").(string))
// 	}

// 	resource, err := client.CreateWebhookAction(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceWebhookActionRead(d, client)
// }

// func resourceWebhookActionRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetCustomAction(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}

// 	if err := d.Set("description", resource.Description); err != nil {
// 		return err
// 	}

// 	if err := d.Set("payload", resource.LiquidTemplate); err != nil {
// 		return err
// 	}

// 	if err := d.Set("url", resource.WebhookURL); err != nil {
// 		return err
// 	}

// 	if err := d.Set("method", string(resource.HTTPMethod)); err != nil {
// 		return err
// 	}

// 	if err := d.Set("headers", resource.Headers); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceWebhookActionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CustomActionsWebhookActionUpdateInput{
// 		Id: *opslevel.NewID(d.Id()),
// 	}

// 	if d.HasChange("name") {
// 		input.Name = opslevel.RefOf(d.Get("name").(string))
// 	}
// 	if d.HasChange("description") {
// 		input.Description = opslevel.RefOf(d.Get("description").(string))
// 	}
// 	if d.HasChange("payload") {
// 		input.LiquidTemplate = opslevel.RefOf(d.Get("payload").(string))
// 	}
// 	if d.HasChange("url") {
// 		input.WebhookUrl = opslevel.RefOf(d.Get("url").(string))
// 	}
// 	if d.HasChange("method") {
// 		input.HttpMethod = opslevel.RefOf(opslevel.CustomActionsHttpMethodEnum(d.Get("method").(string)))
// 	}
// 	if d.HasChange("headers") {
// 		headers := expandHeaders(d.Get("headers"))
// 		input.Headers = &headers
// 	}

// 	_, err := client.UpdateWebhookAction(input)
// 	if err != nil {
// 		return err
// 	}

// 	return resourceWebhookActionRead(d, client)
// }

// func resourceWebhookActionDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteWebhookAction(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
