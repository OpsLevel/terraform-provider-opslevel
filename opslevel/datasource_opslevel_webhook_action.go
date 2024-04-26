package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure WebhookActionDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &WebhookActionDataSource{}

func NewWebhookActionDataSource() datasource.DataSource {
	return &WebhookActionDataSource{}
}

// WebhookActionDataSource manages a WebhookAction data source.
type WebhookActionDataSource struct {
	CommonDataSourceClient
}

type webhookActionWithIdentifierDataSourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Headers     types.Map    `tfsdk:"headers"`
	Id          types.String `tfsdk:"id"`
	Identifier  types.String `tfsdk:"identifier"`
	Method      types.String `tfsdk:"method"`
	Name        types.String `tfsdk:"name"`
	Payload     types.String `tfsdk:"payload"`
	Url         types.String `tfsdk:"url"`
}

type webhookActionDataSourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Headers     types.Map    `tfsdk:"headers"`
	Id          types.String `tfsdk:"id"`
	Method      types.String `tfsdk:"method"`
	Name        types.String `tfsdk:"name"`
	Payload     types.String `tfsdk:"payload"`
	Url         types.String `tfsdk:"url"`
}

func jsonToMap(json map[string]any) map[string]attr.Value {
	jsonAttrs := make(map[string]attr.Value)
	for k, v := range json {
		if value, ok := v.(string); ok {
			jsonAttrs[k] = types.StringValue(value)
		} else {
			jsonAttrs[k] = types.StringNull()
		}
	}
	return jsonAttrs
}

func newWebhookActionWithIdentifierDataSourceModel(ctx context.Context, webhookAction opslevel.CustomActionsExternalAction, identifier string) (webhookActionWithIdentifierDataSourceModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, webhookAction.Aliases)
	jsonAttrs := jsonToMap(webhookAction.Headers)
	action := webhookActionWithIdentifierDataSourceModel{
		Aliases:     aliases,
		Description: types.StringValue(webhookAction.Description),
		Headers:     types.MapValueMust(types.StringType, jsonAttrs),
		Id:          RequiredStringValue(string(webhookAction.Id)),
		Identifier:  types.StringValue(identifier),
		Method:      types.StringValue(string(webhookAction.CustomActionsWebhookAction.HTTPMethod)),
		Name:        types.StringValue(webhookAction.Name),
		Payload:     types.StringValue(webhookAction.LiquidTemplate),
		Url:         types.StringValue(webhookAction.CustomActionsWebhookAction.WebhookURL),
	}

	return action, diags
}

func newWebhookActionDataSourceModel(ctx context.Context, webhookAction opslevel.CustomActionsExternalAction) (webhookActionDataSourceModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, webhookAction.Aliases)
	jsonAttrs := jsonToMap(webhookAction.Headers)
	action := webhookActionDataSourceModel{
		Aliases:     aliases,
		Description: types.StringValue(webhookAction.Description),
		Headers:     types.MapValueMust(types.StringType, jsonAttrs),
		Id:          RequiredStringValue(string(webhookAction.Id)),
		Method:      types.StringValue(string(webhookAction.CustomActionsWebhookAction.HTTPMethod)),
		Name:        types.StringValue(webhookAction.Name),
		Payload:     types.StringValue(webhookAction.LiquidTemplate),
		Url:         types.StringValue(webhookAction.CustomActionsWebhookAction.WebhookURL),
	}

	return action, diags
}

var webhookActionDatasourceSchemaAttrs = map[string]schema.Attribute{
	"aliases": schema.ListAttribute{
		ElementType: types.StringType,
		Description: "The aliases of the Webhook Action.",
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "The description of the Webhook Action.",
		Computed:    true,
	},
	"headers": schema.MapAttribute{
		ElementType: types.StringType,
		Description: "HTTP headers to be passed along with your webhook when triggered.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the Webhook Action.",
		Computed:    true,
	},
	"method": schema.StringAttribute{
		Description: "The http method used to call the Webhook Action.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the Webhook Action.",
		Computed:    true,
	},
	"payload": schema.StringAttribute{
		Description: "Template that can be used to generate a webhook payload.",
		Computed:    true,
	},
	"url": schema.StringAttribute{
		Description: "The URL of the Webhook Action.",
		Computed:    true,
	},
}

func webhookActionAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range webhookActionDatasourceSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

func (d *WebhookActionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_action"
}

func (d *WebhookActionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "WebhookAction data source",

		Attributes: webhookActionAttributes(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the Webhook Action to find.",
				Required:    true,
			},
		}),
	}
}

func (d *WebhookActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data webhookActionWithIdentifierDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookAction, err := d.client.GetCustomAction(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read webhookAction datasource, got error: %s", err))
		return
	}
	webhookActionDataModel, diags := newWebhookActionWithIdentifierDataSourceModel(ctx, *webhookAction, data.Identifier.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel WebhookAction data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &webhookActionDataModel)...)
}
