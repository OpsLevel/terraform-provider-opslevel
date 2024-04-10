package opslevel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/relvacode/iso8601"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var (
	_ resource.ResourceWithConfigure   = &CheckManualResource{}
	_ resource.ResourceWithImportState = &CheckManualResource{}
)

func NewCheckManualResource() resource.Resource {
	return &CheckManualResource{}
}

// CheckManualResource defines the resource implementation.
type CheckManualResource struct {
	CommonResourceClient
}

type CheckUpdateFrequency struct {
	StartingDate types.String `tfsdk:"starting_date"`
	TimeScale    types.String `tfsdk:"time_scale"`
	Value        types.Int64  `tfsdk:"value"`
}

type CheckManualResourceModel struct {
	Category    types.String `tfsdk:"category"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	EnableOn    types.String `tfsdk:"enable_on"`
	Filter      types.String `tfsdk:"filter"`
	Id          types.String `tfsdk:"id"`
	Level       types.String `tfsdk:"level"`
	Name        types.String `tfsdk:"name"`
	Notes       types.String `tfsdk:"notes"`
	Owner       types.String `tfsdk:"owner"`
	LastUpdated types.String `tfsdk:"last_updated"`

	UpdateFrequency       CheckUpdateFrequency `tfsdk:"update_frequency"`
	UpdateRequiresComment types.Bool           `tfsdk:"update_requires_comment"`
}

func NewCheckManualResourceModel(ctx context.Context, check opslevel.Check) CheckManualResourceModel {
	var model CheckManualResourceModel

	model.Category = types.StringValue(string(check.Category.Id))
	model.Enabled = types.BoolValue(check.Enabled)
	model.EnableOn = types.StringValue(check.EnableOn.Time.Format(time.RFC3339))
	model.Filter = types.StringValue(string(check.Filter.Id))
	model.Id = types.StringValue(string(check.Id))
	model.Level = types.StringValue(string(check.Level.Id))
	model.Name = types.StringValue(check.Name)
	model.Notes = types.StringValue(check.Notes)
	model.Owner = types.StringValue(string(check.Owner.Team.Id))
	model.LastUpdated = timeLastUpdated()

	model.UpdateRequiresComment = types.BoolValue(check.UpdateRequiresComment)
	model.UpdateFrequency = CheckUpdateFrequency{
		StartingDate: types.StringValue(check.UpdateFrequency.StartingDate.Time.Format(time.RFC3339)),
		TimeScale:    types.StringValue(string(check.UpdateFrequency.FrequencyTimeScale)),
		Value:        types.Int64Value(int64(check.UpdateFrequency.FrequencyValue)),
	}

	return model
}

func (r *CheckManualResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_manual"
}

func (r *CheckManualResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Manual Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"update_requires_comment": schema.BoolAttribute{
				Description: "Whether the check requires a comment or not.",
				Required:    true,
			},
			"update_frequency": schema.SingleNestedAttribute{
				Description: "Defines the minimum frequency of the updates.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"starting_date": schema.StringAttribute{
						Description: "The date that the check will start to evaluate.",
						Required:    true,
					},
					"time_scale": schema.StringAttribute{
						Description: "The time scale type for the frequency.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(opslevel.AllFrequencyTimeScale...),
						},
					},
					"value": schema.Int64Attribute{
						Description: "The value to be used together with the frequency time_scale.",
						Required:    true,
					},
				},
			},
		}),
	}
}

func (r *CheckManualResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model CheckManualResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(model.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckManualCreateInput{
		CategoryId: asID(model.Category),
		Enabled:    model.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(model.Filter)),
		LevelId:    asID(model.Level),
		Name:       model.Name.ValueString(),
		Notes:      model.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(model.Owner)),
	}
	input.UpdateRequiresComment = model.UpdateRequiresComment.ValueBool()
	input.UpdateFrequency = opslevel.NewManualCheckFrequencyInput(
		timetypes.NewRFC3339ValueMust(model.UpdateFrequency.StartingDate.ValueString()).String(),
		opslevel.FrequencyTimeScale(model.UpdateFrequency.TimeScale.ValueString()),
		int(model.UpdateFrequency.Value.ValueInt64()),
	)

	data, err := r.client.CreateCheckManual(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_manual, got error: %s", err))
		return
	}

	state := NewCheckManualResourceModel(ctx, *data)

	tflog.Trace(ctx, "created a check manual resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CheckManualResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model CheckManualResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(model.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check manual, got error: %s", err))
		return
	}
	state := NewCheckManualResourceModel(ctx, *data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CheckManualResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model CheckManualResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(model.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
		return
	}
	input := opslevel.CheckManualUpdateInput{
		CategoryId: opslevel.RefOf(asID(model.Category)),
		Enabled:    model.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(model.Filter)),
		LevelId:    opslevel.RefOf(asID(model.Level)),
		Id:         asID(model.Id),
		Name:       opslevel.RefOf(model.Name.ValueString()),
		Notes:      model.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(model.Owner)),
	}
	input.UpdateRequiresComment = model.UpdateRequiresComment.ValueBoolPointer()
	// TODO: this is fucking ugly
	startingDate, err := iso8601.ParseString(model.UpdateFrequency.StartingDate.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
		return
	}
	timescale := opslevel.FrequencyTimeScale(model.UpdateFrequency.TimeScale.ValueString())
	value := int(model.UpdateFrequency.Value.ValueInt64())
	input.UpdateFrequency = &opslevel.ManualCheckFrequencyUpdateInput{
		StartingDate:       &iso8601.Time{Time: startingDate},
		FrequencyTimeScale: &timescale,
		FrequencyValue:     &value,
	}

	data, err := r.client.UpdateCheckManual(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_manual, got error: %s", err))
		return
	}

	state := NewCheckManualResourceModel(ctx, *data)

	tflog.Trace(ctx, "updated a check manual resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CheckManualResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model CheckManualResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(model.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check manual, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check manual resource")
}

func (r *CheckManualResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
