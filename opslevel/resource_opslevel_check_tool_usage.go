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
	_ resource.ResourceWithConfigure      = &CheckToolUsageResource{}
	_ resource.ResourceWithImportState    = &CheckToolUsageResource{}
	_ resource.ResourceWithValidateConfig = &CheckToolUsageResource{}
)

func NewCheckToolUsageResource() resource.Resource {
	return &CheckToolUsageResource{}
}

// CheckToolUsageResource defines the resource implementation.
type CheckToolUsageResource struct {
	CommonResourceClient
}

type CheckToolUsageResourceModel struct {
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

	ToolCategory         types.String `tfsdk:"tool_category"`
	ToolNamePredicate    types.Object `tfsdk:"tool_name_predicate"`
	ToolUrlPredicate     types.Object `tfsdk:"tool_url_predicate"`
	EnvironmentPredicate types.Object `tfsdk:"environment_predicate"`
}

func NewCheckToolUsageResourceModel(ctx context.Context, check opslevel.Check, planModel CheckToolUsageResourceModel) CheckToolUsageResourceModel {
	var stateModel CheckToolUsageResourceModel

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
	stateModel.ToolCategory = RequiredStringValue(string(check.ToolCategory))

	if check.ToolNamePredicate == nil {
		stateModel.ToolNamePredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *check.ToolNamePredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.ToolNamePredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}

	if check.ToolUrlPredicate == nil {
		stateModel.ToolUrlPredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *check.ToolUrlPredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.ToolUrlPredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}

	if check.EnvironmentPredicate == nil {
		stateModel.EnvironmentPredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *check.EnvironmentPredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.EnvironmentPredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}

	return stateModel
}

func (r *CheckToolUsageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_tool_usage"
}

func (r *CheckToolUsageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Tool Usage Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"tool_category": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The category that the tool belongs to. One of `%s`",
					strings.Join(opslevel.AllToolCategory, "`, `"),
				),
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf(opslevel.AllToolCategory...)},
			},
			"environment_predicate": PredicateSchema(),
			"tool_name_predicate":   PredicateSchema(),
			"tool_url_predicate":    PredicateSchema(),
		}),
	}
}

func (r *CheckToolUsageResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Description: "Check Tool Usage Resource",
				Attributes: getCheckBaseSchemaV0(map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of this resource.",
						Computed:    true,
					},
					"tool_category": schema.StringAttribute{
						Description: "The category that the tool belongs to.",
						Required:    true,
					},
				}),
				Blocks: map[string]schema.Block{
					"environment_predicate": schema.ListNestedBlock{
						NestedObject: predicateSchemaV0,
					},
					"tool_name_predicate": schema.ListNestedBlock{
						NestedObject: predicateSchemaV0,
					},
					"tool_url_predicate": schema.ListNestedBlock{
						NestedObject: predicateSchemaV0,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var diags diag.Diagnostics
				upgradedStateModel := CheckToolUsageResourceModel{}
				environmentPredicateList := types.ListNull(types.ObjectType{AttrTypes: predicateType})
				toolNamePredicateList := types.ListNull(types.ObjectType{AttrTypes: predicateType})
				toolUrlPredicateList := types.ListNull(types.ObjectType{AttrTypes: predicateType})

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

				// tool usage specific attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tool_category"), &upgradedStateModel.ToolCategory)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_predicate"), &environmentPredicateList)...)
				if len(environmentPredicateList.Elements()) == 1 {
					environmentPredicate := environmentPredicateList.Elements()[0]
					upgradedStateModel.EnvironmentPredicate, diags = types.ObjectValueFrom(ctx, predicateType, environmentPredicate)
					resp.Diagnostics.Append(diags...)
				} else {
					upgradedStateModel.EnvironmentPredicate = types.ObjectNull(predicateType)
				}
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tool_name_predicate"), &toolNamePredicateList)...)
				if len(toolNamePredicateList.Elements()) == 1 {
					toolNamePredicate := toolNamePredicateList.Elements()[0]
					upgradedStateModel.ToolNamePredicate, diags = types.ObjectValueFrom(ctx, predicateType, toolNamePredicate)
					resp.Diagnostics.Append(diags...)
				} else {
					upgradedStateModel.ToolNamePredicate = types.ObjectNull(predicateType)
				}
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tool_url_predicate"), &toolUrlPredicateList)...)
				if len(toolUrlPredicateList.Elements()) == 1 {
					toolUrlPredicate := toolUrlPredicateList.Elements()[0]
					upgradedStateModel.ToolUrlPredicate, diags = types.ObjectValueFrom(ctx, predicateType, toolUrlPredicate)
					resp.Diagnostics.Append(diags...)
				} else {
					upgradedStateModel.ToolUrlPredicate = types.ObjectNull(predicateType)
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateModel)...)
			},
		},
	}
}

func (r *CheckToolUsageResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configModel CheckToolUsageResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkToolUsagePredicates := map[string]types.Object{
		"environment_predicate": configModel.EnvironmentPredicate,
		"tool_name_predicate":   configModel.ToolNamePredicate,
		"tool_url_predicate":    configModel.ToolUrlPredicate,
	}
	for predicateSchemaName, predicate := range checkToolUsagePredicates {
		if predicate.IsNull() || predicate.IsUnknown() {
			continue
		}
		predicateModel, diags := PredicateObjectToModel(ctx, predicate)
		resp.Diagnostics.Append(diags...)
		if err := predicateModel.Validate(); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(predicateSchemaName), "Invalid Attribute Configuration", err.Error())
		}
	}
}

func (r *CheckToolUsageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckToolUsageCreateInput{
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
	if resp.Diagnostics.HasError() {
		return
	}

	input.ToolCategory = opslevel.ToolCategory(planModel.ToolCategory.ValueString())

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.EnvironmentPredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.EnvironmentPredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("environment_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// convert tool_name_predicate object to model from plan
	predicateModel, diags = PredicateObjectToModel(ctx, planModel.ToolNamePredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.ToolNamePredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("tool_name_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// convert tool_url_predicate object to model from plan
	predicateModel, diags = PredicateObjectToModel(ctx, planModel.ToolUrlPredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.ToolUrlPredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("tool_url_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateCheckToolUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_tool_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckToolUsageResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check tool usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckToolUsageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check tool usage, got error: %s", err))
		return
	}
	stateModel := NewCheckToolUsageResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckToolUsageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckToolUsageUpdateInput{
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
	if resp.Diagnostics.HasError() {
		return
	}

	input.ToolCategory = opslevel.RefOf(opslevel.ToolCategory(planModel.ToolCategory.ValueString()))
	nullPredicateModel := PredicateModel{}

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.EnvironmentPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.EnvironmentPredicate = nullPredicateModel.ToUpdateInput()
	} else if err := predicateModel.Validate(); err == nil {
		input.EnvironmentPredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("environment_predicate"), "Invalid Attribute Configuration", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// convert tool_name_predicate object to model from plan
	predicateModel, diags = PredicateObjectToModel(ctx, planModel.ToolNamePredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.ToolNamePredicate = nullPredicateModel.ToUpdateInput()
	} else if err := predicateModel.Validate(); err == nil {
		input.ToolNamePredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("tool_name_predicate"), "Invalid Attribute Configuration", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// convert tool_url_predicate object to model from plan
	predicateModel, diags = PredicateObjectToModel(ctx, planModel.ToolUrlPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.ToolUrlPredicate = nullPredicateModel.ToUpdateInput()
	} else if err := predicateModel.Validate(); err == nil {
		input.ToolUrlPredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("tool_url_predicate"), "Invalid Attribute Configuration", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	input.ToolCategory = opslevel.RefOf(opslevel.ToolCategory(planModel.ToolCategory.ValueString()))

	data, err := r.client.UpdateCheckToolUsage(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_tool_usage, got error: %s", err))
		return
	}

	stateModel := NewCheckToolUsageResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check tool usage resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckToolUsageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckToolUsageResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check tool usage, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check tool usage resource")
}

func (r *CheckToolUsageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
