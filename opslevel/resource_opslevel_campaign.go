package opslevel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
	"github.com/relvacode/iso8601"
)

var _ resource.ResourceWithConfigure = &CampaignResource{}

var _ resource.ResourceWithImportState = &CampaignResource{}

func NewCampaignResource() resource.Resource {
	return &CampaignResource{}
}

type CampaignResource struct {
	CommonResourceClient
}

type CampaignResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	OwnerId      types.String `tfsdk:"owner_id"`
	FilterId     types.String `tfsdk:"filter_id"`
	ProjectBrief types.String `tfsdk:"project_brief"`
	StartDate    types.String `tfsdk:"start_date"`
	TargetDate   types.String `tfsdk:"target_date"`
	Status       types.String `tfsdk:"status"`
	HtmlUrl      types.String `tfsdk:"html_url"`
}

func NewCampaignResourceModel(campaign opslevel.Campaign, givenModel CampaignResourceModel) CampaignResourceModel {
	model := CampaignResourceModel{
		Id:           ComputedStringValue(string(campaign.Id)),
		Name:         RequiredStringValue(campaign.Name),
		OwnerId:      OptionalStringValue(string(campaign.Owner.Id)),
		FilterId:     OptionalStringValue(string(campaign.Filter.Id)),
		ProjectBrief: StringValueFromResourceAndModelField(campaign.RawProjectBrief, givenModel.ProjectBrief),
		Status:       ComputedStringValue(string(campaign.Status)),
		HtmlUrl:      ComputedStringValue(campaign.HtmlUrl),
	}

	if !campaign.StartDate.IsZero() {
		model.StartDate = types.StringValue(campaign.StartDate.Format("2006-01-02"))
	} else if givenModel.StartDate.IsNull() || givenModel.StartDate.IsUnknown() {
		model.StartDate = types.StringNull()
	} else {
		model.StartDate = types.StringNull()
	}

	if !campaign.TargetDate.IsZero() {
		model.TargetDate = types.StringValue(campaign.TargetDate.Format("2006-01-02"))
	} else if givenModel.TargetDate.IsNull() || givenModel.TargetDate.IsUnknown() {
		model.TargetDate = types.StringNull()
	} else {
		model.TargetDate = types.StringNull()
	}

	return model
}

func (r *CampaignResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaign"
}

func (r *CampaignResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Campaign Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the campaign.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the campaign.",
				Required:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The ID of the team that owns this campaign.",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"filter_id": schema.StringAttribute{
				Description: "The ID of the filter applied to this campaign.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"project_brief": schema.StringAttribute{
				Description: "The project brief of the campaign (Markdown).",
				Optional:    true,
			},
			"start_date": schema.StringAttribute{
				Description: "The start date of the campaign (YYYY-MM-DD). Setting both start_date and target_date schedules the campaign.",
				Optional:    true,
			},
			"target_date": schema.StringAttribute{
				Description: "The target end date of the campaign (YYYY-MM-DD). Setting both start_date and target_date schedules the campaign.",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "The current status of the campaign (draft, scheduled, in_progress, delayed, ended).",
				Computed:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "The URL to the campaign in the OpsLevel UI.",
				Computed:    true,
			},
		},
	}
}

func (r *CampaignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CampaignResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CampaignCreateInput{
		Name:    planModel.Name.ValueString(),
		OwnerId: opslevel.ID(planModel.OwnerId.ValueString()),
	}
	if !planModel.FilterId.IsNull() {
		input.FilterId = opslevel.RefOf(opslevel.ID(planModel.FilterId.ValueString()))
	}
	if !planModel.ProjectBrief.IsNull() {
		brief := planModel.ProjectBrief.ValueString()
		input.ProjectBrief = &brief
	}

	campaign, err := r.client.CreateCampaign(input)
	if err != nil || campaign == nil {
		title, detail := formatOpslevelError("create campaign", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	if !planModel.StartDate.IsNull() && !planModel.TargetDate.IsNull() {
		startDate, sdErr := iso8601.ParseString(planModel.StartDate.ValueString() + "T00:00:00Z")
		targetDate, tdErr := iso8601.ParseString(planModel.TargetDate.ValueString() + "T00:00:00Z")
		if sdErr != nil || tdErr != nil {
			resp.Diagnostics.AddError("invalid date", "start_date and target_date must be valid dates (YYYY-MM-DD)")
			return
		}
		scheduled, err := r.client.ScheduleCampaign(opslevel.CampaignScheduleUpdateInput{
			Id:         campaign.Id,
			StartDate:  iso8601.Time{Time: startDate},
			TargetDate: iso8601.Time{Time: targetDate},
		})
		if err != nil {
			title, detail := formatOpslevelError("schedule campaign", err)
			resp.Diagnostics.AddError(title, detail)
			return
		}
		campaign = scheduled
	}

	createdModel := NewCampaignResourceModel(*campaign, planModel)
	tflog.Trace(ctx, "created a campaign resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdModel)...)
}

func (r *CampaignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CampaignResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	campaign, err := r.client.GetCampaign(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil || campaign == nil {
		if (campaign == nil || campaign.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		title, detail := formatOpslevelError("read campaign", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	readModel := NewCampaignResourceModel(*campaign, stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

func (r *CampaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CampaignResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[CampaignResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	campaignId := opslevel.ID(stateModel.Id.ValueString())

	updateInput := opslevel.CampaignUpdateInput{
		Id: campaignId,
	}

	nameVal := planModel.Name.ValueString()
	updateInput.Name = &nameVal

	updateInput.OwnerId = opslevel.RefOf(opslevel.ID(planModel.OwnerId.ValueString()))

	if !planModel.FilterId.IsNull() {
		updateInput.FilterId = opslevel.RefOf(opslevel.ID(planModel.FilterId.ValueString()))
	} else if !stateModel.FilterId.IsNull() {
		updateInput.FilterId = opslevel.RefOf(opslevel.ID(""))
	}

	if !planModel.ProjectBrief.IsNull() {
		brief := planModel.ProjectBrief.ValueString()
		updateInput.ProjectBrief = &brief
	}

	campaign, err := r.client.UpdateCampaign(updateInput)
	if err != nil {
		title, detail := formatOpslevelError("update campaign", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	planHasDates := !planModel.StartDate.IsNull() && !planModel.TargetDate.IsNull()
	stateHasDates := !stateModel.StartDate.IsNull() && !stateModel.TargetDate.IsNull()

	if planHasDates {
		startDate, sdErr := iso8601.ParseString(planModel.StartDate.ValueString() + "T00:00:00Z")
		targetDate, tdErr := iso8601.ParseString(planModel.TargetDate.ValueString() + "T00:00:00Z")
		if sdErr != nil || tdErr != nil {
			resp.Diagnostics.AddError("invalid date", "start_date and target_date must be valid dates (YYYY-MM-DD)")
			return
		}
		scheduled, err := r.client.ScheduleCampaign(opslevel.CampaignScheduleUpdateInput{
			Id:         campaignId,
			StartDate:  iso8601.Time{Time: startDate},
			TargetDate: iso8601.Time{Time: targetDate},
		})
		if err != nil {
			title, detail := formatOpslevelError("schedule campaign", err)
			resp.Diagnostics.AddError(title, detail)
			return
		}
		campaign = scheduled
	} else if stateHasDates && !planHasDates {
		unscheduled, err := r.client.UnscheduleCampaign(campaignId)
		if err != nil {
			title, detail := formatOpslevelError("unschedule campaign", err)
			resp.Diagnostics.AddError(title, detail)
			return
		}
		campaign = unscheduled
	}

	updatedModel := NewCampaignResourceModel(*campaign, planModel)
	tflog.Trace(ctx, "updated a campaign resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func (r *CampaignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CampaignResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCampaign(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		title, detail := formatOpslevelError("delete campaign", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}
	tflog.Trace(ctx, "deleted a campaign resource")
}

func (r *CampaignResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
