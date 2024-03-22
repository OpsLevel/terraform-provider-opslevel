package opslevel

import (
	"context"
	"fmt"

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

// WebhookActionDataSourceModel describes the data source data model.
type WebhookActionDataSourceModel struct {
	Description types.String `tfsdk:"description"`
	Headers     types.Map    `tfsdk:"headers"`
	Identifier  types.String `tfsdk:"identifier"`
	Method      types.String `tfsdk:"method"`
	Name        types.String `tfsdk:"name"`
	Payload     types.String `tfsdk:"payload"`
	Url         types.String `tfsdk:"url"`
}

func jsonToMap(json map[string]any) map[string]attr.Value {
	jsonAttrs := make(map[string]attr.Value)
	for k, v := range json {
		jsonAttrs[k] = types.StringValue(v.(string))
		if v == nil {
			jsonAttrs[k] = types.StringNull()
		}
	}
	return jsonAttrs
}

func NewWebhookActionDataSourceModel(ctx context.Context, webhookAction opslevel.CustomActionsExternalAction, identifier string) WebhookActionDataSourceModel {
	jsonAttrs := jsonToMap(webhookAction.Headers)
	webhookActionDataSourceModel := WebhookActionDataSourceModel{
		Description: types.StringValue(webhookAction.Description),
		Headers:     types.MapValueMust(types.StringType, jsonAttrs),
		Identifier:  types.StringValue(identifier),
		Method:      types.StringValue(string(webhookAction.CustomActionsWebhookAction.HTTPMethod)),
		Name:        types.StringValue(webhookAction.Name),
		Payload:     types.StringValue(webhookAction.LiquidTemplate),
		Url:         types.StringValue(webhookAction.CustomActionsWebhookAction.WebhookURL),
	}

	return webhookActionDataSourceModel
}

func (d *WebhookActionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_action"
}

func (d *WebhookActionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "WebhookAction data source",

		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "The description of the Webhook Action.",
				Computed:    true,
			},
			"headers": schema.MapAttribute{
				ElementType: types.StringType,
				Description: "HTTP headers to be passed along with your webhook when triggered.",
				Computed:    true,
			},
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the Webhook Action to find.",
				Required:    true,
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
		},
	}
}

func (d *WebhookActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WebhookActionDataSourceModel

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
	webhookActionDataModel := NewWebhookActionDataSourceModel(ctx, *webhookAction, data.Identifier.ValueString())

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel WebhookAction data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &webhookActionDataModel)...)
}
