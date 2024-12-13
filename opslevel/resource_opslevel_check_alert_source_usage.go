package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.ResourceWithConfigure      = &CheckAlertSourceUsageResource{}
	_ resource.ResourceWithImportState    = &CheckAlertSourceUsageResource{}
	_ resource.ResourceWithValidateConfig = &CheckAlertSourceUsageResource{}
)

func NewCheckAlertSourceUsageResource() resource.Resource {
	return &CheckAlertSourceUsageResource{}
}

// CheckAlertSourceUsageResource defines the resource implementation.
type CheckAlertSourceUsageResource struct {
	CommonResourceClient
}

type CheckAlertSourceUsageResourceModel struct {
	CheckCodeBaseResourceModel

	AlertType          types.String `tfsdk:"alert_type"`
	AlertNamePredicate types.Object `tfsdk:"alert_name_predicate"`
}

func NewCheckAlertSourceUsageResourceModel(ctx context.Context, check opslevel.Check, planModel CheckAlertSourceUsageResourceModel) CheckAlertSourceUsageResourceModel {
	var stateModel CheckAlertSourceUsageResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
	stateModel.Description = ComputedStringValue(check.Description)
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
	stateModel.AlertType = types.StringValue(string(check.AlertSourceType))

	if check.AlertSourceNamePredicate == nil {
		stateModel.AlertNamePredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *check.AlertSourceNamePredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.AlertNamePredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}

	return stateModel
}

func (r *CheckAlertSourceUsageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_alert_source_usage"
}

func (r *CheckAlertSourceUsageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Alert Source Usage Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"alert_type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The type of the alert source. One of `%s`",
					strings.Join(opslevel.AllAlertSourceTypeEnum, "`, `"),
				),
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf(opslevel.AllAlertSourceTypeEnum...)},
			},
			"alert_name_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckAlertSourceUsageResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Description: "Check Alert Source Usage Resource",
				Attributes: getCheckBaseSchemaV0(map[string]schema.Attribute{
					"alert_type": schema.StringAttribute{
						Description: "The type of the alert source.",
						Required:    true,
					},
					"id": schema.StringAttribute{
						Description: "The ID of this resource.",
						Computed:    true,
					},
				}),
				Blocks: map[string]schema.Block{
					"alert_name_predicate": schema.ListNestedBlock{
						NestedObject: predicateSchemaV0,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var diags diag.Diagnostics
				upgradedStateModel := CheckAlertSourceUsageResourceModel{}
				alertNamePredicateList := types.ListNull(types.ObjectType{AttrTypes: predicateType})

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

				// alert source specific attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("alert_type"), &upgradedStateModel.AlertType)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("alert_name_predicate"), &alertNamePredicateList)...)
				if len(alertNamePredicateList.Elements()) == 1 {
					alertNamePredicate := alertNamePredicateList.Elements()[0]
					upgradedStateModel.AlertNamePredicate, diags = types.ObjectValueFrom(ctx, predicateType, alertNamePredicate)
					resp.Diagnostics.Append(diags...)
				} else {
					upgradedStateModel.AlertNamePredicate = types.ObjectNull(predicateType)
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateModel)...)
			},
		},
	}
}

func (r *CheckAlertSourceUsageResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	alertNamePredicate := types.ObjectNull(predicateType)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("alert_name_predicate"), &alertNamePredicate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicateModel, diags := PredicateObjectToModel(ctx, alertNamePredicate)
	resp.Diagnostics.Append(diags...)
	if err := predicateModel.Validate(); err != nil {
		resp.Diagnostics.AddAttributeWarning(path.Root("alert_name_predicate"), "Invalid Attribute Configuration", err.Error())
	}
}

func (r *CheckAlertSourceUsageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckAlertSourceUsageResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckAlertSourceUsageCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = &iso8601.Time{Time: enabledOn}
	}
	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(planModel.AlertType.ValueString()))

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.AlertNamePredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.AlertSourceNamePredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("alert_name_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateCheckAlertSourceUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_alert_source_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check alert source usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckAlertSourceUsageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckAlertSourceUsageResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check alert source usage, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckAlertSourceUsageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckAlertSourceUsageResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckAlertSourceUsageUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		Id:         asID(planModel.Id),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.RefOf(planModel.Notes.ValueString()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = &iso8601.Time{Time: enabledOn}
	}
	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(planModel.AlertType.ValueString()))

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.AlertNamePredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.AlertSourceNamePredicate = &opslevel.PredicateUpdateInput{}
	} else if err := predicateModel.Validate(); err == nil {
		input.AlertSourceNamePredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("alert_name_predicate"), "Invalid Attribute Configuration", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdateCheckAlertSourceUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_alert_source_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckAlertSourceUsageResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check alert source usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckAlertSourceUsageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckAlertSourceUsageResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check alert source usage, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check alert source usage resource")
}

func (r *CheckAlertSourceUsageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
