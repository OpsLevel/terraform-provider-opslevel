package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// campaignCheckAttrTypes describes one entry in a campaign's checks list.
var campaignCheckAttrTypes = map[string]attr.Type{
	"id":              types.StringType,
	"name":            types.StringType,
	"source_check_id": types.StringType,
}

// CampaignCheckListValue converts a campaign's checks into the data source's
// checks attribute.
//
// A check with no source reports source_check_id as null rather than an empty
// string: the source may have been deleted, or the copy may predate OpsLevel
// recording where copies came from. Null distinguishes "unknown" from a real id.
func CampaignCheckListValue(ctx context.Context, checks []opslevel.CampaignCheckNode) (types.List, diag.Diagnostics) {
	values := make([]attr.Value, 0, len(checks))
	for _, check := range checks {
		sourceCheckId := types.StringNull()
		if check.SourceCheck.Id != "" {
			sourceCheckId = types.StringValue(string(check.SourceCheck.Id))
		}
		obj, diags := types.ObjectValue(campaignCheckAttrTypes, map[string]attr.Value{
			"id":              types.StringValue(string(check.Id)),
			"name":            types.StringValue(check.Name),
			"source_check_id": sourceCheckId,
		})
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: campaignCheckAttrTypes}), diags
		}
		values = append(values, obj)
	}
	return types.ListValue(types.ObjectType{AttrTypes: campaignCheckAttrTypes}, values)
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
	Checks       types.List   `tfsdk:"checks"`
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
		Checks:       types.ListNull(types.ObjectType{AttrTypes: campaignCheckAttrTypes}),
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
			"checks": schema.ListNestedAttribute{
				Description: "The checks that are on this campaign. Use these ids to import existing " +
					"campaign checks as opslevel_campaign_check resources.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The id of the check on the campaign.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The display name of the check.",
							Computed:    true,
						},
						"source_check_id": schema.StringAttribute{
							Description: "The id of the rubric check this check was copied from. Null if " +
								"the source check has been deleted, or if the copy was made before " +
								"OpsLevel began recording this.",
							Computed: true,
						},
					},
				},
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

	campaignChecks, err := d.client.ListCampaignChecks(campaign.Id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read campaign checks, got error: %s", err))
		return
	}
	checks, diags := CampaignCheckListValue(ctx, campaignChecks)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateModel.Checks = checks

	tflog.Trace(ctx, "read an OpsLevel Campaign data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
