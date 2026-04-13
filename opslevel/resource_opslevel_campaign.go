package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

var _ resource.ResourceWithValidateConfig = &CampaignResource{}

func NewCampaignResource() resource.Resource {
	return &CampaignResource{}
}

type CampaignResource struct {
	CommonResourceClient
}

type CampaignResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OwnerId        types.String `tfsdk:"owner_id"`
	FilterId       types.String `tfsdk:"filter_id"`
	ProjectBrief   types.String `tfsdk:"project_brief"`
	CheckIds       types.List   `tfsdk:"check_ids"`
	StartDate      types.String `tfsdk:"start_date"`
	TargetDate     types.String `tfsdk:"target_date"`
	Status         types.String `tfsdk:"status"`
	HtmlUrl        types.String `tfsdk:"html_url"`
}

func NewCampaignResourceModel(campaign opslevel.Campaign, givenModel CampaignResourceModel) CampaignResourceModel {
	model := CampaignResourceModel{
		Id:           ComputedStringValue(string(campaign.Id)),
		Name:         RequiredStringValue(campaign.Name),
		OwnerId:      OptionalStringValue(string(campaign.Owner.Id)),
		FilterId:     OptionalStringValue(string(campaign.Filter.Id)),
		ProjectBrief: StringValueFromResourceAndModelField(campaign.RawProjectBrief, givenModel.ProjectBrief),
		CheckIds:     types.ListNull(types.StringType),
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
			"check_ids": schema.ListAttribute{
				Description: "List of rubric check IDs to associate with this campaign. On create, checks are copied into the campaign. On update, checks are added or removed to match the desired set.",
				Optional:    true,
				ElementType: types.StringType,
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

func (r *CampaignResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CampaignResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasStart := !config.StartDate.IsNull() && !config.StartDate.IsUnknown()
	hasTarget := !config.TargetDate.IsNull() && !config.TargetDate.IsUnknown()
	if hasStart != hasTarget {
		resp.Diagnostics.AddError(
			"Invalid Campaign Schedule",
			"Both start_date and target_date must be set together to schedule a campaign, or both must be omitted.",
		)
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

	if !planModel.CheckIds.IsNull() && !planModel.CheckIds.IsUnknown() {
		checkIds := extractCheckIds(ctx, &resp.Diagnostics, planModel.CheckIds)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(checkIds) > 0 {
			updated, err := r.client.CopyChecksToCampaign(opslevel.ChecksCopyToCampaignInput{
				CampaignId: campaign.Id,
				CheckIds:   checkIds,
			})
			if err != nil {
				title, detail := formatOpslevelError("copy checks to campaign", err)
				resp.Diagnostics.AddError(title, detail)
				return
			}
			campaign = updated
		}
	}

	createdModel := NewCampaignResourceModel(*campaign, planModel)
	createdModel.CheckIds = planModel.CheckIds
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
	readModel.CheckIds = r.readCampaignCheckIds(ctx, &resp.Diagnostics, campaign.Id, stateModel.CheckIds)
	if resp.Diagnostics.HasError() {
		return
	}
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

	r.reconcileCampaignChecks(ctx, &resp.Diagnostics, campaignId, stateModel, planModel)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedModel := NewCampaignResourceModel(*campaign, planModel)
	updatedModel.CheckIds = planModel.CheckIds
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

// readCampaignCheckIds queries the campaign's actual checks from the API and
// returns only the rubric check IDs (from priorCheckIds) whose corresponding
// campaign check still exists. This enables drift detection when checks are
// removed outside Terraform.
func (r *CampaignResource) readCampaignCheckIds(
	ctx context.Context,
	diags *diag.Diagnostics,
	campaignId opslevel.ID,
	priorCheckIds types.List,
) types.List {
	if priorCheckIds.IsNull() || priorCheckIds.IsUnknown() {
		return types.ListNull(types.StringType)
	}

	var priorIds []string
	diags.Append(priorCheckIds.ElementsAs(ctx, &priorIds, false)...)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	if len(priorIds) == 0 {
		return priorCheckIds
	}

	campaignChecks, err := r.client.ListCampaignChecks(campaignId)
	if err != nil {
		title, detail := formatOpslevelError("list campaign checks for read", err)
		diags.AddError(title, detail)
		return types.ListNull(types.StringType)
	}

	campaignCheckNames := make(map[string]bool, len(campaignChecks))
	for _, cc := range campaignChecks {
		campaignCheckNames[cc.Name] = true
	}

	var verified []string
	for _, rubricID := range priorIds {
		check, err := r.client.GetCheck(opslevel.ID(rubricID))
		if err != nil {
			tflog.Warn(ctx, "could not look up rubric check during read, keeping in state",
				map[string]any{"rubric_check_id": rubricID, "error": err.Error()})
			verified = append(verified, rubricID)
			continue
		}
		if campaignCheckNames[check.Name] {
			verified = append(verified, rubricID)
		} else {
			tflog.Info(ctx, "rubric check no longer present in campaign, removing from state",
				map[string]any{"rubric_check_id": rubricID, "check_name": check.Name})
		}
	}

	if len(verified) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}

	vals := make([]attr.Value, len(verified))
	for i, id := range verified {
		vals[i] = types.StringValue(id)
	}
	return types.ListValueMust(types.StringType, vals)
}

// DiffCheckIds computes which IDs to add and remove given two sets of IDs.
func DiffCheckIds(stateIds, planIds map[string]bool) (toAdd []string, toRemove []string) {
	for id := range planIds {
		if !stateIds[id] {
			toAdd = append(toAdd, id)
		}
	}
	for id := range stateIds {
		if !planIds[id] {
			toRemove = append(toRemove, id)
		}
	}
	return toAdd, toRemove
}

func (r *CampaignResource) reconcileCampaignChecks(
	ctx context.Context,
	diags *diag.Diagnostics,
	campaignId opslevel.ID,
	stateModel CampaignResourceModel,
	planModel CampaignResourceModel,
) {
	stateIds := extractCheckIdSet(ctx, diags, stateModel.CheckIds)
	planIds := extractCheckIdSet(ctx, diags, planModel.CheckIds)
	if diags.HasError() {
		return
	}

	added, toRemove := DiffCheckIds(stateIds, planIds)

	toAdd := make([]opslevel.ID, len(added))
	for i, id := range added {
		toAdd[i] = opslevel.ID(id)
	}

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return
	}

	if len(toRemove) > 0 {
		rubricNamesByID := make(map[string]string, len(toRemove))
		for _, rubricID := range toRemove {
			check, err := r.client.GetCheck(opslevel.ID(rubricID))
			if err != nil {
				diags.AddWarning(
					"could not look up rubric check",
					fmt.Sprintf("Could not fetch rubric check %s to match for removal: %s", rubricID, err),
				)
				continue
			}
			rubricNamesByID[rubricID] = check.Name
		}

		campaignChecks, err := r.client.ListCampaignChecks(campaignId)
		if err != nil {
			title, detail := formatOpslevelError("list campaign checks", err)
			diags.AddError(title, detail)
			return
		}

		campaignCheckByName := make(map[string]opslevel.ID, len(campaignChecks))
		for _, cc := range campaignChecks {
			campaignCheckByName[cc.Name] = cc.Id
		}

		for _, name := range rubricNamesByID {
			ccID, ok := campaignCheckByName[name]
			if !ok {
				tflog.Warn(ctx, "campaign check not found for removal", map[string]any{"check_name": name})
				continue
			}
			if err := r.client.DeleteCheck(ccID); err != nil {
				title, detail := formatOpslevelError("delete campaign check", err)
				diags.AddError(title, detail)
				return
			}
			tflog.Info(ctx, "removed campaign check", map[string]any{"check_name": name, "campaign_check_id": string(ccID)})
		}
	}

	if len(toAdd) > 0 {
		_, err := r.client.CopyChecksToCampaign(opslevel.ChecksCopyToCampaignInput{
			CampaignId: campaignId,
			CheckIds:   toAdd,
		})
		if err != nil {
			title, detail := formatOpslevelError("copy checks to campaign", err)
			diags.AddError(title, detail)
			return
		}
		tflog.Info(ctx, "added checks to campaign", map[string]any{"count": len(toAdd)})
	}
}

func extractCheckIdSet(ctx context.Context, diags *diag.Diagnostics, list types.List) map[string]bool {
	if list.IsNull() || list.IsUnknown() {
		return map[string]bool{}
	}
	var ids []string
	diags.Append(list.ElementsAs(ctx, &ids, false)...)
	if diags.HasError() {
		return nil
	}
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set
}

func extractCheckIds(ctx context.Context, diags *diag.Diagnostics, list types.List) []opslevel.ID {
	var ids []string
	diags.Append(list.ElementsAs(ctx, &ids, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]opslevel.ID, len(ids))
	for i, id := range ids {
		result[i] = opslevel.ID(id)
	}
	return result
}
