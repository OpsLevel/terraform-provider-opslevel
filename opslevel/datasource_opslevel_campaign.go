package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	"reminder": schema.SingleNestedAttribute{
		Description: "The recurring reminder configured for this campaign, if any.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"channels": schema.ListAttribute{
				Description: "The channels through which reminders are delivered.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"frequency": schema.Int64Attribute{
				Description: "The interval (in frequency_unit) at which reminders are delivered.",
				Computed:    true,
			},
			"frequency_unit": schema.StringAttribute{
				Description: "The unit of the frequency interval.",
				Computed:    true,
			},
			"time_of_day": schema.StringAttribute{
				Description: "The time of day reminders are delivered, in 24-hour \"HH:MM\" format.",
				Computed:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "The IANA timezone in which time_of_day is evaluated.",
				Computed:    true,
			},
			"days_of_week": schema.ListAttribute{
				Description: "The weekdays on which reminders are delivered (weekly cadence only).",
				ElementType: types.StringType,
				Computed:    true,
			},
			"message": schema.StringAttribute{
				Description: "The custom message included in the reminder.",
				Computed:    true,
			},
			"default_slack_channel": schema.StringAttribute{
				Description: "Slack channel notified when a team has no default Slack contact.",
				Computed:    true,
			},
			"default_microsoft_teams_channel": schema.StringAttribute{
				Description: "Microsoft Teams channel notified when a team has no default Teams contact.",
				Computed:    true,
			},
			"next_occurrence": schema.StringAttribute{
				Description: "The next time a reminder will be delivered, based on the current configuration.",
				Computed:    true,
			},
		},
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
	Reminder     types.Object `tfsdk:"reminder"`
	StartDate    types.String `tfsdk:"start_date"`
	Status       types.String `tfsdk:"status"`
	TargetDate   types.String `tfsdk:"target_date"`
}

func newCampaignDataSourceModel(ctx context.Context, diags *diag.Diagnostics, campaign opslevel.Campaign, identifier string) campaignDataSourceModel {
	model := campaignDataSourceModel{
		FilterId:     ComputedStringValue(string(campaign.Filter.Id)),
		HtmlUrl:      ComputedStringValue(campaign.HtmlUrl),
		Id:           ComputedStringValue(string(campaign.Id)),
		Identifier:   ComputedStringValue(identifier),
		Name:         ComputedStringValue(campaign.Name),
		OwnerId:      ComputedStringValue(string(campaign.Owner.Id)),
		ProjectBrief: ComputedStringValue(campaign.RawProjectBrief),
		Reminder:     campaignReminderToObject(ctx, diags, campaign.Reminder, types.ObjectNull(reminderAttrTypes())),
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

	stateModel := newCampaignDataSourceModel(ctx, &resp.Diagnostics, *campaign, configModel.Identifier.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "read an OpsLevel Campaign data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
