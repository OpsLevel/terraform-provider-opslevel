package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

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

func newWebhookActionDataSourcesAllModel(ctx context.Context, webhookActions []opslevel.CustomActionsExternalAction) (webhookActionDataSourcesAllModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	webhookActionModels := make([]webhookActionDataSourceModel, 0)
	for _, webhookAction := range webhookActions {
		webhookActionModel, tmpDiags := newWebhookActionDataSourceModel(ctx, webhookAction)
		diags.Append(tmpDiags...)
		webhookActionModels = append(webhookActionModels, webhookActionModel)
	}
	return webhookActionDataSourcesAllModel{WebhookActions: webhookActionModels}, diags
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
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list webhookActions, got error: %s", err))
		return
	}
	stateModel, diags := newWebhookActionDataSourcesAllModel(ctx, webhookActions.Nodes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "listed all OpsLevel WebhookAction data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
