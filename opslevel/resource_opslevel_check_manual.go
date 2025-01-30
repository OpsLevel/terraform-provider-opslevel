package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

var updateFrequencyTypeV0 = map[string]attr.Type{
	"starting_data": types.StringType,
	"time_value":    types.StringType,
	"value":         types.StringType,
}

type CheckManualResourceModel struct {
	CheckCodeBaseResourceModel

	UpdateFrequency       *CheckUpdateFrequency `tfsdk:"update_frequency"`
	UpdateRequiresComment types.Bool            `tfsdk:"update_requires_comment"`
}

func NewCheckManualResourceModel(ctx context.Context, check opslevel.Check, planModel CheckManualResourceModel) CheckManualResourceModel {
	var stateModel CheckManualResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
	if planModel.Enabled.IsNull() {
		stateModel.Enabled = types.BoolValue(false)
	} else {
		stateModel.Enabled = OptionalBoolValue(&check.Enabled)
	}
	if planModel.EnableOn.IsNull() {
		stateModel.EnableOn = types.StringNull()
	} else {
		// We pass through the plan value because of time formatting issue to ensure the state gets the exact value the customer specified
		stateModel.EnableOn = planModel.EnableOn
	}
	stateModel.Filter = OptionalStringValue(string(check.Filter.Id))
	stateModel.Id = ComputedStringValue(string(check.Id))
	stateModel.Level = RequiredStringValue(string(check.Level.Id))
	stateModel.Name = RequiredStringValue(check.Name)
	stateModel.Notes = OptionalStringValue(check.Notes)
	stateModel.Owner = OptionalStringValue(string(check.Owner.Team.Id))

	stateModel.UpdateRequiresComment = RequiredBoolValue(check.UpdateRequiresComment)
	if planModel.UpdateFrequency != nil {
		stateModel.UpdateFrequency = &CheckUpdateFrequency{
			StartingDate: planModel.UpdateFrequency.StartingDate,
			TimeScale:    RequiredStringValue(string(check.UpdateFrequency.FrequencyTimeScale)),
			Value:        RequiredIntValue(check.UpdateFrequency.FrequencyValue),
		}
	}

	return stateModel
}

func (r *CheckManualResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_manual"
}

func (r *CheckManualResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
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
						Description: fmt.Sprintf(
							"The time scale type for the frequency. One of `%s`",
							strings.Join(opslevel.AllFrequencyTimeScale, "`, `"),
						),
						Required: true,
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

func (r *CheckManualResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Description: "Check Repository File Resource",
				Attributes: getCheckBaseSchemaV0(map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of this resource.",
						Computed:    true,
					},
					"update_requires_comment": schema.BoolAttribute{
						Description: "Whether the check requires a comment or not.",
						Optional:    true,
					},
				}),
				Blocks: map[string]schema.Block{
					"update_frequency": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"starting_data": schema.StringAttribute{
									Description: "The date that the check will start to evaluate.",
									Required:    true,
								},
								"time_scale": schema.StringAttribute{
									Description: "The time scale type for the frequency.",
									Required:    true,
								},
								"value": schema.Int64Attribute{
									Description: "The value to be used together with the frequency time_scale.",
									Required:    true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				// var diags diag.Diagnostics
				upgradedStateModel := CheckManualResourceModel{}
				updateFrequencyList := types.ListNull(types.ObjectType{AttrTypes: updateFrequencyTypeV0})

				// base check attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("category"), &upgradedStateModel.Category)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("enable_on"), &upgradedStateModel.EnableOn)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("enabled"), &upgradedStateModel.Enabled)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("filter"), &upgradedStateModel.Filter)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &upgradedStateModel.Id)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("level"), &upgradedStateModel.Level)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &upgradedStateModel.Name)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("notes"), &upgradedStateModel.Notes)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("owner"), &upgradedStateModel.Owner)...)

				// repository file specific attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("update_requires_comment"), &upgradedStateModel.UpdateRequiresComment)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("update_frequency"), &updateFrequencyList)...)
				if len(updateFrequencyList.Elements()) == 1 {
					updateFrequency := updateFrequencyList.Elements()[0].(basetypes.ObjectValue)
					updateFrequencyAttrs := updateFrequency.Attributes()
					upgradedStateModel.UpdateFrequency = &CheckUpdateFrequency{
						StartingDate: updateFrequencyAttrs["starting_data"].(basetypes.StringValue),
						TimeScale:    updateFrequencyAttrs["time_scale"].(basetypes.StringValue),
						Value:        updateFrequencyAttrs["value"].(basetypes.Int64Value),
					}
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateModel)...)
			},
		},
	}
}

func (r *CheckManualResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckManualResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckManualCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      nullable(planModel.Notes.ValueStringPointer()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}
	input.UpdateRequiresComment = planModel.UpdateRequiresComment.ValueBool()
	if planModel.UpdateFrequency != nil {
		input.UpdateFrequency = opslevel.NewManualCheckFrequencyInput(
			planModel.UpdateFrequency.StartingDate.ValueString(),
			opslevel.FrequencyTimeScale(planModel.UpdateFrequency.TimeScale.ValueString()),
			int(planModel.UpdateFrequency.Value.ValueInt64()),
		)
	}

	data, err := r.client.CreateCheckManual(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_manual, got error: %s", err))
		return
	}

	stateModel := NewCheckManualResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check manual resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckManualResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckManualResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check manual, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckManualResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform stateModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckManualResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckManualResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckManualUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Id:         asID(planModel.Id),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      nullable(planModel.Notes.ValueStringPointer()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}
	input.UpdateRequiresComment = opslevel.RefOf(planModel.UpdateRequiresComment.ValueBool())
	if planModel.UpdateFrequency != nil {
		input.UpdateFrequency = opslevel.NewManualCheckFrequencyUpdateInput(
			planModel.UpdateFrequency.StartingDate.ValueString(),
			opslevel.FrequencyTimeScale(planModel.UpdateFrequency.TimeScale.ValueString()),
			int(planModel.UpdateFrequency.Value.ValueInt64()),
		)
	}

	data, err := r.client.UpdateCheckManual(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_manual, got error: %s", err))
		return
	}

	stateModel := NewCheckManualResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check manual resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckManualResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckManualResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check manual, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check manual resource")
}

func (r *CheckManualResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
