package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	OwnerId      types.String `tfsdk:"owner_id"`
	FilterId     types.String `tfsdk:"filter_id"`
	ProjectBrief types.String `tfsdk:"project_brief"`
	CheckIds     types.List   `tfsdk:"check_ids"`
	StartDate    types.String `tfsdk:"start_date"`
	TargetDate   types.String `tfsdk:"target_date"`
	Reminder     types.Object `tfsdk:"reminder"`
	Status       types.String `tfsdk:"status"`
	HtmlUrl      types.String `tfsdk:"html_url"`
}

// CampaignReminderModel describes the nested reminder configuration block.
type CampaignReminderModel struct {
	Channels                     types.List   `tfsdk:"channels"`
	Frequency                    types.Int64  `tfsdk:"frequency"`
	FrequencyUnit                types.String `tfsdk:"frequency_unit"`
	TimeOfDay                    types.String `tfsdk:"time_of_day"`
	Timezone                     types.String `tfsdk:"timezone"`
	DaysOfWeek                   types.List   `tfsdk:"days_of_week"`
	Message                      types.String `tfsdk:"message"`
	DefaultSlackChannel          types.String `tfsdk:"default_slack_channel"`
	DefaultMicrosoftTeamsChannel types.String `tfsdk:"default_microsoft_teams_channel"`
	NextOccurrence               types.String `tfsdk:"next_occurrence"`
}

func reminderAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"channels":                        types.ListType{ElemType: types.StringType},
		"frequency":                       types.Int64Type,
		"frequency_unit":                  types.StringType,
		"time_of_day":                     types.StringType,
		"timezone":                        types.StringType,
		"days_of_week":                    types.ListType{ElemType: types.StringType},
		"message":                         types.StringType,
		"default_slack_channel":           types.StringType,
		"default_microsoft_teams_channel": types.StringType,
		"next_occurrence":                 types.StringType,
	}
}

// buildCampaignReminderInput converts the Terraform reminder object into the SDK input.
// Returns nil when no reminder block is configured. days_of_week is only sent for
// weekly cadences (the API rejects it otherwise).
func buildCampaignReminderInput(ctx context.Context, diags *diag.Diagnostics, obj types.Object) *opslevel.CampaignReminderInput {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	var rm CampaignReminderModel
	diags.Append(obj.As(ctx, &rm, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	channelStrings, d := ListValueToStringSlice(ctx, rm.Channels)
	diags.Append(d...)
	channels := make([]opslevel.CampaignReminderChannelEnum, len(channelStrings))
	for i, c := range channelStrings {
		channels[i] = opslevel.CampaignReminderChannelEnum(c)
	}

	input := &opslevel.CampaignReminderInput{
		Channels:      channels,
		Frequency:     int(rm.Frequency.ValueInt64()),
		FrequencyUnit: opslevel.CampaignReminderFrequencyUnitEnum(rm.FrequencyUnit.ValueString()),
		TimeOfDay:     rm.TimeOfDay.ValueString(),
		Timezone:      rm.Timezone.ValueString(),
	}

	if rm.FrequencyUnit.ValueString() == string(opslevel.CampaignReminderFrequencyUnitEnumWeek) {
		dayStrings, dd := ListValueToStringSlice(ctx, rm.DaysOfWeek)
		diags.Append(dd...)
		for _, day := range dayStrings {
			input.DaysOfWeek = append(input.DaysOfWeek, opslevel.DayOfWeekEnum(day))
		}
	}

	if !rm.Message.IsNull() && !rm.Message.IsUnknown() {
		input.Message = rm.Message.ValueStringPointer()
	}
	if !rm.DefaultSlackChannel.IsNull() && !rm.DefaultSlackChannel.IsUnknown() {
		input.DefaultSlackChannel = rm.DefaultSlackChannel.ValueStringPointer()
	}
	if !rm.DefaultMicrosoftTeamsChannel.IsNull() && !rm.DefaultMicrosoftTeamsChannel.IsUnknown() {
		input.DefaultMicrosoftTeamsChannel = rm.DefaultMicrosoftTeamsChannel.ValueStringPointer()
	}
	return input
}

// preserveHashChannel keeps the user-configured channel value when it is
// semantically equal to the API value (the API prefixes a leading '#'), so
// "platform-eng" in config does not perpetually drift against "#platform-eng".
func preserveHashChannel(apiVal string, given types.String) types.String {
	if apiVal == "" {
		return types.StringNull()
	}
	if !given.IsNull() && !given.IsUnknown() {
		g := given.ValueString()
		if g == apiVal || "#"+strings.TrimPrefix(g, "#") == apiVal {
			return given
		}
	}
	return types.StringValue(apiVal)
}

// campaignReminderToObject maps an SDK reminder onto a Terraform object,
// preserving user-formatted values from the given (plan/state) object where
// the API would otherwise introduce noise. Returns a null object when no
// reminder is configured on the campaign.
func campaignReminderToObject(ctx context.Context, diags *diag.Diagnostics, reminder *opslevel.CampaignReminder, given types.Object) types.Object {
	if reminder == nil {
		return types.ObjectNull(reminderAttrTypes())
	}

	var prior CampaignReminderModel
	havePrior := !given.IsNull() && !given.IsUnknown()
	if havePrior {
		diags.Append(given.As(ctx, &prior, basetypes.ObjectAsOptions{})...)
	}

	channels := make([]attr.Value, len(reminder.Channels))
	for i, c := range reminder.Channels {
		channels[i] = types.StringValue(string(c))
	}

	daysList := types.ListNull(types.StringType)
	if len(reminder.DaysOfWeek) > 0 {
		days := make([]attr.Value, len(reminder.DaysOfWeek))
		for i, day := range reminder.DaysOfWeek {
			days[i] = types.StringValue(string(day))
		}
		daysList = types.ListValueMust(types.StringType, days)
	}

	nextOccurrence := types.StringNull()
	if !reminder.NextOccurrence.IsZero() {
		nextOccurrence = types.StringValue(reminder.NextOccurrence.Format(time.RFC3339))
	}

	rm := CampaignReminderModel{
		Channels:                     types.ListValueMust(types.StringType, channels),
		Frequency:                    types.Int64Value(int64(reminder.Frequency)),
		FrequencyUnit:                types.StringValue(string(reminder.FrequencyUnit)),
		TimeOfDay:                    types.StringValue(reminder.TimeOfDay),
		Timezone:                     types.StringValue(reminder.Timezone),
		DaysOfWeek:                   daysList,
		Message:                      StringValueFromResourceAndModelField(reminder.Message, prior.Message),
		DefaultSlackChannel:          preserveHashChannel(reminder.DefaultSlackChannel, prior.DefaultSlackChannel),
		DefaultMicrosoftTeamsChannel: preserveHashChannel(reminder.DefaultMicrosoftTeamsChannel, prior.DefaultMicrosoftTeamsChannel),
		NextOccurrence:               nextOccurrence,
	}

	obj, d := types.ObjectValueFrom(ctx, reminderAttrTypes(), rm)
	diags.Append(d...)
	return obj
}

func NewCampaignResourceModel(campaign opslevel.Campaign, givenModel CampaignResourceModel) CampaignResourceModel {
	model := CampaignResourceModel{
		Id:           ComputedStringValue(string(campaign.Id)),
		Name:         RequiredStringValue(campaign.Name),
		OwnerId:      OptionalStringValue(string(campaign.Owner.Id)),
		FilterId:     OptionalStringValue(string(campaign.Filter.Id)),
		ProjectBrief: StringValueFromResourceAndModelField(campaign.RawProjectBrief, givenModel.ProjectBrief),
		CheckIds:     types.ListNull(types.StringType),
		Reminder:     types.ObjectNull(reminderAttrTypes()),
		Status:       ComputedStringValue(string(campaign.Status)),
		HtmlUrl:      ComputedStringValue(campaign.HtmlUrl),
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
				Computed:    true,
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
			"reminder": schema.SingleNestedAttribute{
				Description: "Configuration for recurring campaign reminders sent to component owners via Slack, email, or Microsoft Teams. Reminders are only delivered while the campaign is in_progress or delayed.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"channels": schema.ListAttribute{
						Description: "The channels through which reminders are delivered. One or more of: " + strings.Join(opslevel.AllCampaignReminderChannelEnum, ", ") + ".",
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(stringvalidator.OneOf(opslevel.AllCampaignReminderChannelEnum...)),
						},
					},
					"frequency": schema.Int64Attribute{
						Description: "The interval (in frequency_unit) at which reminders are delivered. Must be at least 1.",
						Required:    true,
						Validators:  []validator.Int64{int64validator.AtLeast(1)},
					},
					"frequency_unit": schema.StringAttribute{
						Description: "The unit of the frequency interval. One of: " + strings.Join(opslevel.AllCampaignReminderFrequencyUnitEnum, ", ") + ".",
						Required:    true,
						Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllCampaignReminderFrequencyUnitEnum...)},
					},
					"time_of_day": schema.StringAttribute{
						Description: "The time of day reminders are delivered, in 24-hour \"HH:MM\" format.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`),
								"must be a 24-hour time in \"HH:MM\" format (e.g. \"09:30\")",
							),
						},
					},
					"timezone": schema.StringAttribute{
						Description: "The IANA timezone in which time_of_day is evaluated (e.g. \"America/Chicago\").",
						Required:    true,
					},
					"days_of_week": schema.ListAttribute{
						Description: "The weekdays on which reminders are delivered. Only supported (and required) when frequency_unit is \"week\".",
						Optional:    true,
						ElementType: types.StringType,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(stringvalidator.OneOf(opslevel.AllDayOfWeekEnum...)),
						},
					},
					"message": schema.StringAttribute{
						Description: "An optional custom message included in the reminder.",
						Optional:    true,
					},
					"default_slack_channel": schema.StringAttribute{
						Description: "Slack channel notified when a team has no default Slack contact. A leading '#' is added automatically.",
						Optional:    true,
					},
					"default_microsoft_teams_channel": schema.StringAttribute{
						Description: "Microsoft Teams channel notified when a team has no default Teams contact.",
						Optional:    true,
					},
					"next_occurrence": schema.StringAttribute{
						Description: "The next time a reminder will be delivered, based on the current configuration. Null until the campaign is in_progress or delayed.",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
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

	if !config.Reminder.IsNull() && !config.Reminder.IsUnknown() {
		var rm CampaignReminderModel
		resp.Diagnostics.Append(config.Reminder.As(ctx, &rm, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		isWeekly := rm.FrequencyUnit.ValueString() == string(opslevel.CampaignReminderFrequencyUnitEnumWeek)
		hasDays := !rm.DaysOfWeek.IsNull() && !rm.DaysOfWeek.IsUnknown() && len(rm.DaysOfWeek.Elements()) > 0

		if isWeekly && !hasDays {
			resp.Diagnostics.AddAttributeError(
				path.Root("reminder").AtName("days_of_week"),
				"Invalid Reminder Configuration",
				"days_of_week must contain at least one day when frequency_unit is \"week\".",
			)
		}
		if !isWeekly && hasDays {
			resp.Diagnostics.AddAttributeError(
				path.Root("reminder").AtName("days_of_week"),
				"Invalid Reminder Configuration",
				"days_of_week is only supported when frequency_unit is \"week\".",
			)
		}
	}
}

func (r *CampaignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CampaignResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CampaignCreateInput{
		Name:     planModel.Name.ValueString(),
		OwnerId:  opslevel.ID(planModel.OwnerId.ValueString()),
		FilterId: nullableID(planModel.FilterId.ValueStringPointer()),
	}
	if !planModel.ProjectBrief.IsNull() {
		brief := planModel.ProjectBrief.ValueString()
		input.ProjectBrief = &brief
	}
	input.Reminder = buildCampaignReminderInput(ctx, &resp.Diagnostics, planModel.Reminder)
	if resp.Diagnostics.HasError() {
		return
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
	if !planModel.CheckIds.IsUnknown() {
		createdModel.CheckIds = planModel.CheckIds
	}
	createdModel.Reminder = campaignReminderToObject(ctx, &resp.Diagnostics, campaign.Reminder, planModel.Reminder)
	if resp.Diagnostics.HasError() {
		return
	}
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
	readModel.Reminder = campaignReminderToObject(ctx, &resp.Diagnostics, campaign.Reminder, stateModel.Reminder)
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

	updateInput.FilterId = nullableID(planModel.FilterId.ValueStringPointer())

	brief := planModel.ProjectBrief.ValueString()
	updateInput.ProjectBrief = &brief

	plannedReminder := buildCampaignReminderInput(ctx, &resp.Diagnostics, planModel.Reminder)
	if resp.Diagnostics.HasError() {
		return
	}
	if plannedReminder != nil {
		updateInput.Reminder = opslevel.NewNullableFrom(*plannedReminder)
	} else if !stateModel.Reminder.IsNull() {
		// Reminder block removed from config -> explicitly clear it on the API.
		updateInput.Reminder = opslevel.NewNullOf[opslevel.CampaignReminderInput]()
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
	if !planModel.CheckIds.IsUnknown() {
		updatedModel.CheckIds = planModel.CheckIds
	}
	updatedModel.Reminder = campaignReminderToObject(ctx, &resp.Diagnostics, campaign.Reminder, planModel.Reminder)
	if resp.Diagnostics.HasError() {
		return
	}
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

	campaignCheckNames := make(map[string]int, len(campaignChecks))
	for _, cc := range campaignChecks {
		campaignCheckNames[cc.Name]++
	}
	for name, count := range campaignCheckNames {
		if count > 1 {
			tflog.Warn(ctx, "multiple campaign checks share the same name; matching may be unreliable",
				map[string]any{"check_name": name, "count": count})
		}
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
		if campaignCheckNames[check.Name] > 0 {
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
		seen := make(map[string]int, len(campaignChecks))
		for _, cc := range campaignChecks {
			seen[cc.Name]++
			campaignCheckByName[cc.Name] = cc.Id
		}
		for name, count := range seen {
			if count > 1 {
				tflog.Warn(ctx, "multiple campaign checks share the same name; deletion may target the wrong check",
					map[string]any{"check_name": name, "count": count})
			}
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
