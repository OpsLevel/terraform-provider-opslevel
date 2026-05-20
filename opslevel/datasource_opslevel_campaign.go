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

var _ datasource.DataSourceWithConfigure = &CampaignDataSource{}

func NewCampaignDataSource() datasource.DataSource {
	return &CampaignDataSource{}
}

type CampaignDataSource struct {
	CommonDataSourceClient
}

var campaignSchemaAttrs = map[string]schema.Attribute{
	"filter_id": schema.StringAttribute{
		Description: "The ID of the filter applied to this campaign.",
		Computed:    true,
	},
	"html_url": schema.StringAttribute{
		Description: "The URL to the campaign in the OpsLevel UI.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the campaign.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the campaign.",
		Computed:    true,
	},
	"owner_id": schema.StringAttribute{
		Description: "The ID of the team that owns this campaign.",
		Computed:    true,
	},
	"project_brief": schema.StringAttribute{
		Description: "The raw project brief of the campaign (Markdown).",
		Computed:    true,
	},
	"start_date": schema.StringAttribute{
		Description: "The start date of the campaign.",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "The current status of the campaign.",
		Computed:    true,
	},
	"target_date": schema.StringAttribute{
		Description: "The target end date of the campaign.",
		Computed:    true,
	},
}

func CampaignAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range campaignSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

type campaignDataSourceModel struct {
	FilterId     types.String `tfsdk:"filter_id"`
	HtmlUrl      types.String `tfsdk:"html_url"`
	Id           types.String `tfsdk:"id"`
	Identifier   types.String `tfsdk:"identifier"`
	Name         types.String `tfsdk:"name"`
	OwnerId      types.String `tfsdk:"owner_id"`
	ProjectBrief types.String `tfsdk:"project_brief"`
	StartDate    types.String `tfsdk:"start_date"`
	Status       types.String `tfsdk:"status"`
	TargetDate   types.String `tfsdk:"target_date"`
}

func newCampaignDataSourceModel(campaign opslevel.Campaign, identifier string) campaignDataSourceModel {
	model := campaignDataSourceModel{
		FilterId:     ComputedStringValue(string(campaign.Filter.Id)),
		HtmlUrl:      ComputedStringValue(campaign.HtmlUrl),
		Id:           ComputedStringValue(string(campaign.Id)),
		Identifier:   ComputedStringValue(identifier),
		Name:         ComputedStringValue(campaign.Name),
		OwnerId:      ComputedStringValue(string(campaign.Owner.Id)),
		ProjectBrief: ComputedStringValue(campaign.RawProjectBrief),
		Status:       ComputedStringValue(string(campaign.Status)),
	}
	if !campaign.StartDate.IsZero() {
		model.StartDate = types.StringValue(campaign.StartDate.Format("2006-01-02"))
	} else {
		model.StartDate = types.StringNull()
	}
	if !campaign.TargetDate.IsZero() {
		model.TargetDate = types.StringValue(campaign.TargetDate.Format("2006-01-02"))
	} else {
		model.TargetDate = types.StringNull()
	}
	return model
}

func (d *CampaignDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaign"
}

func (d *CampaignDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Campaign data source",
		Attributes: CampaignAttributes(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Description: "The id of the campaign to find.",
				Required:    true,
			},
		}),
	}
}

func (d *CampaignDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var configModel campaignDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	campaign, err := d.client.GetCampaign(opslevel.ID(configModel.Identifier.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read campaign datasource, got error: %s", err))
		return
	}

	stateModel := newCampaignDataSourceModel(*campaign, configModel.Identifier.ValueString())
	tflog.Trace(ctx, "read an OpsLevel Campaign data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
