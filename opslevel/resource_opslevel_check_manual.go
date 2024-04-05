package opslevel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/opslevel/terraform-provider-opslevel/opslevel/utils"
)

// CheckManualResource defines the resource implementation.
type CheckManualResource struct {
	CommonResourceClient
	utils.BaseResource
}

var _ resource.ResourceWithConfigure = &CheckManualResource{}
var _ resource.ResourceWithImportState = &CheckManualResource{}

func NewCheckManualResource() resource.Resource {
	helper := utils.CRUD[opslevel.Check, CheckManualResourceModel]{
		Title: "Check Manual",
		Schema: schema.Schema{
			// This description is used by the documentation generator and the language server.
			MarkdownDescription: "Check Manual Resource",

			Attributes: CheckBaseAttributes(map[string]schema.Attribute{
				"update_requires_comment": schema.BoolAttribute{
					Description: "Whether the check requires a comment or not.",
					Required:    true,
				},
			}),
			Blocks: map[string]schema.Block{
				"update_frequency": schema.SingleNestedBlock{
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
			},
		},
		BuildModel: NewCheckManualResourceModel,
		DoCreate: func(client *opslevel.Client, model CheckManualResourceModel) (opslevel.Check, diag.Diagnostics) {
			var resp diag.Diagnostics
			input, diags := NewCheckCreateInputFrom[opslevel.CheckManualCreateInput](model.CheckBaseModel)
			resp.Append(diags...)
			input.UpdateRequiresComment = model.UpdateRequiresComment.ValueBool()
			input.UpdateFrequency = opslevel.NewManualCheckFrequencyInput(
				model.UpdateFrequency.StartingDate.ValueString(),
				opslevel.FrequencyTimeScale(model.UpdateFrequency.TimeScale.ValueString()),
				int(model.UpdateFrequency.Value.ValueInt64()),
			)

			data, err := client.CreateCheckManual(*input)
			if err != nil {
				resp.AddError("opslevel client error", err.Error())
			}
			return *data, resp
		},
		DoRead: func(client *opslevel.Client, model CheckManualResourceModel) (opslevel.Check, diag.Diagnostics) {
			var resp diag.Diagnostics
			data, err := client.GetCheck(AsID(model.Id))
			if err != nil {
				resp.AddError("opslevel client error", err.Error())
			}
			return *data, resp
		},
		DoUpdate: func(client *opslevel.Client, model CheckManualResourceModel) (opslevel.Check, diag.Diagnostics) {
			var resp diag.Diagnostics
			input, diags := NewCheckUpdateInputFrom[opslevel.CheckManualUpdateInput](model.CheckBaseModel)
			resp.Append(diags...)
			input.UpdateRequiresComment = model.UpdateRequiresComment.ValueBoolPointer()
			// TODO: this is fucking ugly
			startingDate, diags := AsISO8601(model.UpdateFrequency.StartingDate)
			timescale := opslevel.FrequencyTimeScale(model.UpdateFrequency.TimeScale.ValueString())
			value := int(model.UpdateFrequency.Value.ValueInt64())
			input.UpdateFrequency = &opslevel.ManualCheckFrequencyUpdateInput{
				StartingDate:       startingDate,
				FrequencyTimeScale: &timescale,
				FrequencyValue:     &value,
			}
			resp.Append(diags...)

			data, err := client.UpdateCheckManual(*input)
			if err != nil {
				resp.AddError("opslevel client error", err.Error())
			}
			return *data, resp
		},
		DoDelete: func(client *opslevel.Client, model CheckManualResourceModel) error {
			return client.DeleteCheck(AsID(model.Id))
		},
	}
	output := &CheckManualResource{}
	return output.SetHelper(helper)

}

type CheckUpdateFrequency struct {
	StartingDate timetypes.RFC3339 `tfsdk:"starting_date"`
	TimeScale    types.String      `tfsdk:"time_scale"`
	Value        types.Int64       `tfsdk:"value"`
}

type CheckManualResourceModel struct {
	CheckBaseModel
	UpdateFrequency       CheckUpdateFrequency `tfsdk:"update_frequency"`
	UpdateRequiresComment types.Bool           `tfsdk:"update_requires_comment"`
}

func NewCheckManualResourceModel(ctx context.Context, check opslevel.Check) (CheckManualResourceModel, diag.Diagnostics) {
	var model CheckManualResourceModel

	ApplyCheckBaseModel(check, &model.CheckBaseModel)

	model.UpdateFrequency = CheckUpdateFrequency{
		StartingDate: timetypes.NewRFC3339TimeValue(check.UpdateFrequency.StartingDate.Time),
		TimeScale:    types.StringValue(string(check.UpdateFrequency.FrequencyTimeScale)),
		Value:        types.Int64Value(int64(check.UpdateFrequency.FrequencyValue)),
	}

	return model, diag.Diagnostics{}
}
