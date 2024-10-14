package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ datasource.DataSourceWithConfigure = &WebhookActionDataSourcesAll{}

func NewWebhookActionDataSourcesAll() datasource.DataSource {
	return &WebhookActionDataSourcesAll{}
}

type WebhookActionDataSourcesAll struct {
	CommonDataSourceClient
}

type webhookActionDataSourcesAllModel struct {
	WebhookActions []webhookActionDataSourceModel `tfsdk:"webhook_actions"`
}

func newWebhookActionDataSourcesAllModel(webhookActions []opslevel.CustomActionsExternalAction) webhookActionDataSourcesAllModel {
	webhookActionModels := make([]webhookActionDataSourceModel, 0)
	for _, webhookAction := range webhookActions {
		webhookActionModels = append(webhookActionModels, newWebhookActionDataSourceModel(webhookAction))
	}
	return webhookActionDataSourcesAllModel{WebhookActions: webhookActionModels}
}

func (d *WebhookActionDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_actions"
}

func (d *WebhookActionDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List of all WebhookAction data sources",

		Attributes: map[string]schema.Attribute{
			"webhook_actions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: webhookActionDatasourceSchemaAttrs,
				},
				Description: "List of webhook action data sources",
				Computed:    true,
			},
		},
	}
}

func (d *WebhookActionDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel webhookActionDataSourcesAllModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	webhookActions, err := d.client.ListCustomActions(nil)
	if opslevel.HasBadHttpStatus(err) {
		resp.Diagnostics.AddError("HTTP status error", fmt.Sprintf("Unable to list webhookActions, got error: %s", err))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list webhookActions, got error: %s", err))
		return
	}
	stateModel = newWebhookActionDataSourcesAllModel(webhookActions.Nodes)

	tflog.Trace(ctx, "listed all OpsLevel WebhookAction data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
