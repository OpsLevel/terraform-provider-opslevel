package opslevel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
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

	model.Category = RequiredStringValue(string(check.Category.Id))
	model.Enabled = OptionalBoolValue(&check.Enabled)
	model.EnableOn = OptionalStringValue(check.EnableOn.Time.Format(time.RFC3339))
	model.Filter = OptionalStringValue(string(check.Filter.Id))
	model.Id = ComputedStringValue(string(check.Id))
	model.Level = RequiredStringValue(string(check.Level.Id))
	model.Name = RequiredStringValue(check.Name)
	model.Notes = OptionalStringValue(check.Notes)
	model.Owner = OptionalStringValue(string(check.Owner.Team.Id))

	model.UpdateRequiresComment = RequiredBoolValue(check.UpdateRequiresComment)
	model.UpdateFrequency = CheckUpdateFrequency{
		StartingDate: RequiredStringValue(check.UpdateFrequency.StartingDate.Time.Format(time.RFC3339)),
		TimeScale:    RequiredStringValue(string(check.UpdateFrequency.FrequencyTimeScale)),
		Value:        RequiredIntValue(check.UpdateFrequency.FrequencyValue),
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
	var planModel CheckManualResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckManualCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	input.UpdateRequiresComment = planModel.UpdateRequiresComment.ValueBool()
	input.UpdateFrequency = opslevel.NewManualCheckFrequencyInput(
		planModel.UpdateFrequency.StartingDate.ValueString(),
		opslevel.FrequencyTimeScale(planModel.UpdateFrequency.TimeScale.ValueString()),
		int(planModel.UpdateFrequency.Value.ValueInt64()),
	)

	data, err := r.client.CreateCheckManual(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_manual, got error: %s", err))
		return
	}

	stateModel := NewCheckManualResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.UpdateFrequency.StartingDate = planModel.UpdateFrequency.StartingDate
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a check manual resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckManualResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckManualResourceModel

	// Read Terraform prior stateModel data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check manual, got error: %s", err))
		return
	}
	stateModel := NewCheckManualResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.UpdateFrequency.StartingDate = planModel.UpdateFrequency.StartingDate

	// Save updated data into Terraform stateModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckManualResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckManualResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
		return
	}
	input := opslevel.CheckManualUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Id:         asID(planModel.Id),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	input.UpdateRequiresComment = planModel.UpdateRequiresComment.ValueBoolPointer()
	input.UpdateFrequency = opslevel.NewManualCheckFrequencyUpdateInput(
		planModel.UpdateFrequency.StartingDate.ValueString(),
		opslevel.FrequencyTimeScale(planModel.UpdateFrequency.TimeScale.ValueString()),
		int(planModel.UpdateFrequency.Value.ValueInt64()),
	)

	data, err := r.client.UpdateCheckManual(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_manual, got error: %s", err))
		return
	}

	stateModel := NewCheckManualResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.UpdateFrequency.StartingDate = planModel.UpdateFrequency.StartingDate
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a check manual resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckManualResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckManualResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check manual, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check manual resource")
}

func (r *CheckManualResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
