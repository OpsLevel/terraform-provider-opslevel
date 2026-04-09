package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
)

var _ datasource.DataSourceWithConfigure = &CampaignDataSourcesAll{}

func NewCampaignDataSourcesAll() datasource.DataSource {
	return &CampaignDataSourcesAll{}
}

type CampaignDataSourcesAll struct {
	CommonDataSourceClient
}

type campaignListItemModel struct {
	FilterId     types.String `tfsdk:"filter_id"`
	HtmlUrl      types.String `tfsdk:"html_url"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	OwnerId      types.String `tfsdk:"owner_id"`
	ProjectBrief types.String `tfsdk:"project_brief"`
	StartDate    types.String `tfsdk:"start_date"`
	Status       types.String `tfsdk:"status"`
	TargetDate   types.String `tfsdk:"target_date"`
}

type campaignDataSourcesAllModel struct {
	Status    types.String            `tfsdk:"status"`
	Campaigns []campaignListItemModel `tfsdk:"campaigns"`
}

func newCampaignListItemModels(campaigns []opslevel.Campaign) []campaignListItemModel {
	models := make([]campaignListItemModel, 0, len(campaigns))
	for _, c := range campaigns {
		m := campaignListItemModel{
			FilterId:     OptionalStringValue(string(c.Filter.Id)),
			HtmlUrl:      ComputedStringValue(c.HtmlUrl),
			Id:           ComputedStringValue(string(c.Id)),
			Name:         ComputedStringValue(c.Name),
			OwnerId:      OptionalStringValue(string(c.Owner.Id)),
			ProjectBrief: OptionalStringValue(c.RawProjectBrief),
			Status:       ComputedStringValue(string(c.Status)),
		}
		if !c.StartDate.IsZero() {
			m.StartDate = types.StringValue(c.StartDate.Format("2006-01-02"))
		} else {
			m.StartDate = types.StringNull()
		}
		if !c.TargetDate.IsZero() {
			m.TargetDate = types.StringValue(c.TargetDate.Format("2006-01-02"))
		} else {
			m.TargetDate = types.StringNull()
		}
		models = append(models, m)
	}
	return models
}

func (d *CampaignDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaigns"
}

func (d *CampaignDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Campaign data sources",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Description: "Filter campaigns by status (draft, scheduled, in_progress, delayed, ended). Defaults to in_progress.",
				Optional:    true,
			},
			"campaigns": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: campaignSchemaAttrs,
				},
				Description: "List of Campaign data sources",
				Computed:    true,
			},
		},
	}
}

func (d *CampaignDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var configModel campaignDataSourcesAllModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var listVars *opslevel.ListCampaignsVariables
	if !configModel.Status.IsNull() && !configModel.Status.IsUnknown() {
		status := opslevel.CampaignStatusEnum(configModel.Status.ValueString())
		listVars = &opslevel.ListCampaignsVariables{Status: &status}
	}

	campaigns, err := d.client.ListCampaigns(listVars)
	if err != nil || campaigns == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list campaigns datasource, got error: %s", err))
		return
	}

	stateModel := campaignDataSourcesAllModel{
		Status:    configModel.Status,
		Campaigns: newCampaignListItemModels(campaigns.Nodes),
	}

	tflog.Trace(ctx, "read OpsLevel Campaigns data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
