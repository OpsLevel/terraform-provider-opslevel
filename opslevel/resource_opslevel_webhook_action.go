package opslevel

import (
	"context"
	"fmt"
	"strings"

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

func NewWebhookActionResourceModel(webhookAction opslevel.CustomActionsExternalAction, givenModel WebhookActionResourceModel) WebhookActionResourceModel {
	return WebhookActionResourceModel{
		Description: StringValueFromResourceAndModelField(webhookAction.Description, givenModel.Description),
		Headers:     jsonToMapValue(webhookAction.Headers),
		Id:          ComputedStringValue(string(webhookAction.Id)),
		Method:      RequiredStringValue(string(webhookAction.HTTPMethod)),
		Name:        RequiredStringValue(webhookAction.Name),
		Payload:     RequiredStringValue(webhookAction.LiquidTemplate),
		Url:         RequiredStringValue(webhookAction.WebhookURL),
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
				Description: fmt.Sprintf(
					"The http method used to call the Webhook Action. One of `%s`",
					strings.Join(opslevel.AllCustomActionsHttpMethodEnum, "`, `"),
				),
				Required: true,
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
	planModel := read[WebhookActionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	headersAsJson, diags := MapValueToOpslevelJson(ctx, planModel.Headers)
	if diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to create opslevel.JSON from 'headers': '%s'", planModel.Headers))
		return
	}

	webhookActionInput := opslevel.CustomActionsWebhookActionCreateInput{
		Description:    planModel.Description.ValueStringPointer(),
		Headers:        &headersAsJson,
		HttpMethod:     opslevel.CustomActionsHttpMethodEnum(planModel.Method.ValueString()),
		LiquidTemplate: planModel.Payload.ValueStringPointer(),
		Name:           planModel.Name.ValueString(),
		WebhookUrl:     planModel.Url.ValueString(),
	}
	webhookAction, err := r.client.CreateWebhookAction(webhookActionInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create webhook action, got error: %s", err))
		return
	}
	createdWebhookActionResourceModel := NewWebhookActionResourceModel(*webhookAction, planModel)

	tflog.Trace(ctx, "created a webhook action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdWebhookActionResourceModel)...)
}

func (r *WebhookActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[WebhookActionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookAction, err := r.client.GetCustomAction(stateModel.Id.ValueString())
	if err != nil {
		if (webhookAction == nil || webhookAction.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read webhookAction, got error: %s", err))
		return
	}
	verifiedStateModel := NewWebhookActionResourceModel(*webhookAction, stateModel)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a webhook action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *WebhookActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[WebhookActionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	headersAsJson, diags := MapValueToOpslevelJson(ctx, planModel.Headers)
	if diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to create opslevel.JSON from 'headers': '%s'", planModel.Headers))
		return
	}

	httpMethod := opslevel.CustomActionsHttpMethodEnum(planModel.Method.ValueString())
	updateWebhookActionInput := opslevel.CustomActionsWebhookActionUpdateInput{
		Description:    opslevel.RefOf(planModel.Description.ValueString()),
		Headers:        &headersAsJson,
		HttpMethod:     &httpMethod,
		Id:             opslevel.ID(planModel.Id.ValueString()),
		LiquidTemplate: opslevel.RefOf(planModel.Payload.ValueString()),
		Name:           opslevel.RefOf(planModel.Name.ValueString()),
		WebhookUrl:     opslevel.RefOf(planModel.Url.ValueString()),
	}
	resource, err := r.client.UpdateWebhookAction(updateWebhookActionInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update webhookAction, got error: %s", err))
		return
	}
	updatedWebhookActionResourceModel := NewWebhookActionResourceModel(*resource, planModel)

	tflog.Trace(ctx, "updated a webhook action resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedWebhookActionResourceModel)...)
}

func (r *WebhookActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[WebhookActionResourceModel](ctx, &resp.Diagnostics, req.State)
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
