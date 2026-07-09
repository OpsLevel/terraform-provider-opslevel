package opslevel

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opslevel/opslevel-go/v2026"
	"github.com/relvacode/iso8601"
)

func stringList(values ...string) types.List {
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elems)
}

func newReminderObject(t *testing.T, rm CampaignReminderModel) types.Object {
	t.Helper()
	obj, d := types.ObjectValueFrom(context.Background(), reminderAttrTypes(), rm)
	if d.HasError() {
		t.Fatalf("building reminder object: %v", d)
	}
	return obj
}

func readReminderModel(t *testing.T, obj types.Object) CampaignReminderModel {
	t.Helper()
	var rm CampaignReminderModel
	if d := obj.As(context.Background(), &rm, basetypes.ObjectAsOptions{}); d.HasError() {
		t.Fatalf("reading reminder object: %v", d)
	}
	return rm
}

func TestBuildCampaignReminderInput_NullReturnsNil(t *testing.T) {
	var diags diag.Diagnostics
	got := buildCampaignReminderInput(context.Background(), &diags, types.ObjectNull(reminderAttrTypes()))
	if got != nil {
		t.Fatalf("expected nil for null object, got %+v", got)
	}
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
}

func TestBuildCampaignReminderInput_WeeklyIncludesDays(t *testing.T) {
	var diags diag.Diagnostics
	obj := newReminderObject(t, CampaignReminderModel{
		Channels:                     stringList("slack", "email"),
		Frequency:                    types.Int64Value(1),
		FrequencyUnit:                types.StringValue("week"),
		TimeOfDay:                    types.StringValue("09:30"),
		Timezone:                     types.StringValue("America/Chicago"),
		DaysOfWeek:                   stringList("monday", "thursday"),
		Message:                      types.StringValue("hello"),
		DefaultSlackChannel:          types.StringValue("#platform-eng"),
		DefaultMicrosoftTeamsChannel: types.StringNull(),
		NextOccurrence:               types.StringNull(),
	})

	got := buildCampaignReminderInput(context.Background(), &diags, obj)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if got == nil {
		t.Fatal("expected an input, got nil")
	}
	if got.FrequencyUnit != opslevel.CampaignReminderFrequencyUnitEnumWeek {
		t.Errorf("frequency_unit = %q, want week", got.FrequencyUnit)
	}
	if got.Frequency != 1 {
		t.Errorf("frequency = %d, want 1", got.Frequency)
	}
	if len(got.Channels) != 2 {
		t.Errorf("channels = %v, want 2", got.Channels)
	}
	if len(got.DaysOfWeek) != 2 {
		t.Errorf("days_of_week = %v, want 2", got.DaysOfWeek)
	}
	if got.Message == nil || *got.Message != "hello" {
		t.Errorf("message = %v, want hello", got.Message)
	}
	if got.DefaultSlackChannel == nil || *got.DefaultSlackChannel != "#platform-eng" {
		t.Errorf("default_slack_channel = %v, want #platform-eng", got.DefaultSlackChannel)
	}
	if got.DefaultMicrosoftTeamsChannel != nil {
		t.Errorf("default_microsoft_teams_channel = %v, want nil", got.DefaultMicrosoftTeamsChannel)
	}
}

func TestBuildCampaignReminderInput_NonWeeklyOmitsDays(t *testing.T) {
	var diags diag.Diagnostics
	obj := newReminderObject(t, CampaignReminderModel{
		Channels:                     stringList("slack"),
		Frequency:                    types.Int64Value(2),
		FrequencyUnit:                types.StringValue("day"),
		TimeOfDay:                    types.StringValue("14:00"),
		Timezone:                     types.StringValue("America/Chicago"),
		DaysOfWeek:                   stringList("monday"), // present but must be ignored
		Message:                      types.StringNull(),
		DefaultSlackChannel:          types.StringNull(),
		DefaultMicrosoftTeamsChannel: types.StringNull(),
		NextOccurrence:               types.StringNull(),
	})

	got := buildCampaignReminderInput(context.Background(), &diags, obj)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if got == nil {
		t.Fatal("expected an input, got nil")
	}
	if len(got.DaysOfWeek) != 0 {
		t.Errorf("days_of_week = %v, want empty for daily cadence", got.DaysOfWeek)
	}
	if got.Message != nil {
		t.Errorf("message = %v, want nil", got.Message)
	}
	if got.DefaultSlackChannel != nil {
		t.Errorf("default_slack_channel = %v, want nil", got.DefaultSlackChannel)
	}
}

func TestPreserveHashChannel(t *testing.T) {
	cases := []struct {
		name     string
		api      string
		given    types.String
		wantNull bool
		want     string
	}{
		{"empty api returns null", "", types.StringValue("anything"), true, ""},
		{"given null keeps api value", "#x", types.StringNull(), false, "#x"},
		{"equal ignoring hash keeps given", "#x", types.StringValue("x"), false, "x"},
		{"exact equal keeps given", "#x", types.StringValue("#x"), false, "#x"},
		{"different value uses api", "#x", types.StringValue("y"), false, "#x"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := preserveHashChannel(tc.api, tc.given)
			if tc.wantNull {
				if !got.IsNull() {
					t.Fatalf("expected null, got %q", got.ValueString())
				}
				return
			}
			if got.IsNull() {
				t.Fatal("unexpected null")
			}
			if got.ValueString() != tc.want {
				t.Fatalf("got %q, want %q", got.ValueString(), tc.want)
			}
		})
	}
}

func TestCampaignReminderToObject_NilReturnsNull(t *testing.T) {
	var diags diag.Diagnostics
	obj := campaignReminderToObject(context.Background(), &diags, nil, types.ObjectNull(reminderAttrTypes()))
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if !obj.IsNull() {
		t.Fatal("expected a null object for a nil reminder")
	}
}

func TestCampaignReminderToObject_PopulatedPreservesGivenSlackForm(t *testing.T) {
	var diags diag.Diagnostics
	reminder := &opslevel.CampaignReminder{
		Channels:            []opslevel.CampaignReminderChannelEnum{opslevel.CampaignReminderChannelEnumSlack},
		DaysOfWeek:          []opslevel.DayOfWeekEnum{opslevel.DayOfWeekEnumMonday},
		DefaultSlackChannel: "#platform-eng",
		Frequency:           1,
		FrequencyUnit:       opslevel.CampaignReminderFrequencyUnitEnumWeek,
		Message:             "hello",
		NextOccurrence:      iso8601.Time{Time: time.Date(2026, 7, 1, 9, 30, 0, 0, time.UTC)},
		TimeOfDay:           "09:30",
		Timezone:            "America/Chicago",
	}
	// User configured the channel without the leading '#'; it should be preserved.
	given := newReminderObject(t, CampaignReminderModel{
		Channels:                     stringList("slack"),
		Frequency:                    types.Int64Value(1),
		FrequencyUnit:                types.StringValue("week"),
		TimeOfDay:                    types.StringValue("09:30"),
		Timezone:                     types.StringValue("America/Chicago"),
		DaysOfWeek:                   stringList("monday"),
		Message:                      types.StringValue("hello"),
		DefaultSlackChannel:          types.StringValue("platform-eng"),
		DefaultMicrosoftTeamsChannel: types.StringNull(),
		NextOccurrence:               types.StringNull(),
	})

	obj := campaignReminderToObject(context.Background(), &diags, reminder, given)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	rm := readReminderModel(t, obj)

	if rm.FrequencyUnit.ValueString() != "week" {
		t.Errorf("frequency_unit = %q, want week", rm.FrequencyUnit.ValueString())
	}
	if got := rm.DefaultSlackChannel.ValueString(); got != "platform-eng" {
		t.Errorf("default_slack_channel = %q, want platform-eng (preserved given form)", got)
	}
	if rm.NextOccurrence.IsNull() {
		t.Fatal("expected next_occurrence to be set")
	}
	if got := rm.NextOccurrence.ValueString(); got != "2026-07-01T09:30:00Z" {
		t.Errorf("next_occurrence = %q, want 2026-07-01T09:30:00Z", got)
	}
	if len(rm.DaysOfWeek.Elements()) != 1 {
		t.Errorf("days_of_week length = %d, want 1", len(rm.DaysOfWeek.Elements()))
	}
}

func TestCampaignReminderToObject_NoPriorUsesApiValue(t *testing.T) {
	var diags diag.Diagnostics
	reminder := &opslevel.CampaignReminder{
		Channels:            []opslevel.CampaignReminderChannelEnum{opslevel.CampaignReminderChannelEnumEmail},
		DefaultSlackChannel: "#platform-eng",
		Frequency:           2,
		FrequencyUnit:       opslevel.CampaignReminderFrequencyUnitEnumDay,
		TimeOfDay:           "08:00",
		Timezone:            "UTC",
	}

	obj := campaignReminderToObject(context.Background(), &diags, reminder, types.ObjectNull(reminderAttrTypes()))
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	rm := readReminderModel(t, obj)
	if got := rm.DefaultSlackChannel.ValueString(); got != "#platform-eng" {
		t.Errorf("default_slack_channel = %q, want #platform-eng", got)
	}
	if !rm.DaysOfWeek.IsNull() {
		t.Errorf("days_of_week = %v, want null for daily cadence", rm.DaysOfWeek)
	}
	if !rm.NextOccurrence.IsNull() {
		t.Errorf("next_occurrence = %v, want null when unset", rm.NextOccurrence)
	}
}

func TestNewCampaignDataSourceModel_Reminder(t *testing.T) {
	var diags diag.Diagnostics

	withReminder := opslevel.Campaign{
		Name: "C",
		Reminder: &opslevel.CampaignReminder{
			Channels:      []opslevel.CampaignReminderChannelEnum{opslevel.CampaignReminderChannelEnumEmail},
			DaysOfWeek:    []opslevel.DayOfWeekEnum{opslevel.DayOfWeekEnumMonday},
			Frequency:     1,
			FrequencyUnit: opslevel.CampaignReminderFrequencyUnitEnumWeek,
			TimeOfDay:     "08:00",
			Timezone:      "UTC",
		},
	}
	model := newCampaignDataSourceModel(context.Background(), &diags, withReminder, "id")
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if model.Reminder.IsNull() {
		t.Error("expected reminder to be set on the data source model")
	}

	noReminder := newCampaignDataSourceModel(context.Background(), &diags, opslevel.Campaign{Name: "C"}, "id")
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if !noReminder.Reminder.IsNull() {
		t.Error("expected null reminder when campaign has none")
	}
}

func TestValidateReminderConfig(t *testing.T) {
	base := func(freqUnit string, days types.List) CampaignReminderModel {
		return CampaignReminderModel{
			Channels:                     stringList("slack"),
			Frequency:                    types.Int64Value(1),
			FrequencyUnit:                types.StringValue(freqUnit),
			TimeOfDay:                    types.StringValue("09:30"),
			Timezone:                     types.StringValue("UTC"),
			DaysOfWeek:                   days,
			Message:                      types.StringNull(),
			DefaultSlackChannel:          types.StringNull(),
			DefaultMicrosoftTeamsChannel: types.StringNull(),
			NextOccurrence:               types.StringNull(),
		}
	}

	t.Run("null reminder is valid", func(t *testing.T) {
		if d := validateReminderConfig(context.Background(), types.ObjectNull(reminderAttrTypes())); d.HasError() {
			t.Fatalf("expected no error, got %v", d)
		}
	})

	t.Run("weekly with days is valid", func(t *testing.T) {
		obj := newReminderObject(t, base("week", stringList("monday")))
		if d := validateReminderConfig(context.Background(), obj); d.HasError() {
			t.Fatalf("expected no error, got %v", d)
		}
	})

	t.Run("weekly without days errors", func(t *testing.T) {
		obj := newReminderObject(t, base("week", types.ListNull(types.StringType)))
		if d := validateReminderConfig(context.Background(), obj); !d.HasError() {
			t.Fatal("expected an error for weekly cadence without days_of_week")
		}
	})

	t.Run("daily with days errors", func(t *testing.T) {
		obj := newReminderObject(t, base("day", stringList("monday")))
		if d := validateReminderConfig(context.Background(), obj); !d.HasError() {
			t.Fatal("expected an error for daily cadence with days_of_week")
		}
	})

	t.Run("daily without days is valid", func(t *testing.T) {
		obj := newReminderObject(t, base("day", types.ListNull(types.StringType)))
		if d := validateReminderConfig(context.Background(), obj); d.HasError() {
			t.Fatalf("expected no error, got %v", d)
		}
	})
}

func TestNewCampaignListItemModels_Reminder(t *testing.T) {
	var diags diag.Diagnostics
	campaigns := []opslevel.Campaign{
		{
			Name: "a",
			Reminder: &opslevel.CampaignReminder{
				Channels:      []opslevel.CampaignReminderChannelEnum{opslevel.CampaignReminderChannelEnumSlack},
				Frequency:     1,
				FrequencyUnit: opslevel.CampaignReminderFrequencyUnitEnumDay,
				TimeOfDay:     "08:00",
				Timezone:      "UTC",
			},
		},
		{Name: "b"},
	}

	models := newCampaignListItemModels(context.Background(), &diags, campaigns)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if len(models) != 2 {
		t.Fatalf("models length = %d, want 2", len(models))
	}
	if models[0].Reminder.IsNull() {
		t.Error("expected reminder set on first campaign")
	}
	if !models[1].Reminder.IsNull() {
		t.Error("expected null reminder on second campaign")
	}
}
