package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure      = &CheckTagDefinedResource{}
	_ resource.ResourceWithImportState    = &CheckTagDefinedResource{}
	_ resource.ResourceWithValidateConfig = &CheckTagDefinedResource{}
)

func NewCheckTagDefinedResource() resource.Resource {
	return &CheckTagDefinedResource{}
}

// CheckTagDefinedResource defines the resource implementation.
type CheckTagDefinedResource struct {
	CommonResourceClient
}

type CheckTagDefinedResourceModel struct {
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

	TagKey       types.String `tfsdk:"tag_key"`
	TagPredicate types.Object `tfsdk:"tag_predicate"`
}

func NewCheckTagDefinedResourceModel(ctx context.Context, check opslevel.Check, planModel CheckTagDefinedResourceModel) CheckTagDefinedResourceModel {
	var stateModel CheckTagDefinedResourceModel

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

	stateModel.TagKey = RequiredStringValue(check.TagKey)

	if check.TagPredicate == nil {
		stateModel.TagPredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *check.TagPredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.TagPredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}

	return stateModel
}

func (r *CheckTagDefinedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_tag_defined"
}

func (r *CheckTagDefinedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Tag Defined Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"tag_key": schema.StringAttribute{
				Description: "The tag key where the tag predicate should be applied.",
				Required:    true,
			},
			"tag_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckTagDefinedResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: CheckBaseAttributes(map[string]schema.Attribute{
					"tag_key": schema.StringAttribute{
						Description: "The tag key where the tag predicate should be applied.",
						Required:    true,
					},
					"tag_predicate": PredicateSchema(),
				}),
			},
			StateUpgrader: CheckUpgradeFunc[CheckTagDefinedResourceModel](),
		},
	}
}

func (r *CheckTagDefinedResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configModel CheckTagDefinedResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicateModel, diags := PredicateObjectToModel(ctx, configModel.TagPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		return
	}
	if err := predicateModel.Validate(); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("tag_predicate"), "Invalid Attribute Configuration", err.Error())
	}
}

func (r *CheckTagDefinedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckTagDefinedResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckTagDefinedCreateInput{
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

	input.TagKey = planModel.TagKey.ValueString()

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.TagPredicate)
	resp.Diagnostics.Append(diags...)
	if !predicateModel.Type.IsUnknown() && !predicateModel.Type.IsNull() {
		if err := predicateModel.Validate(); err == nil {
			input.TagPredicate = predicateModel.ToCreateInput()
		} else {
			resp.Diagnostics.AddAttributeError(path.Root("tag_predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateCheckTagDefined(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_tag_defined, got error: %s", err))
		return
	}

	stateModel := NewCheckTagDefinedResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check tag defined resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckTagDefinedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckTagDefinedResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check tag defined, got error: %s", err))
		return
	}
	stateModel := NewCheckTagDefinedResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckTagDefinedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckTagDefinedResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckTagDefinedUpdateInput{
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

	input.TagKey = planModel.TagKey.ValueStringPointer()

	// convert environment_predicate object to model from plan
	predicateModel, diags := PredicateObjectToModel(ctx, planModel.TagPredicate)
	resp.Diagnostics.Append(diags...)
	if predicateModel.Type.IsUnknown() || predicateModel.Type.IsNull() {
		input.TagPredicate = &opslevel.PredicateUpdateInput{}
	} else if err := predicateModel.Validate(); err == nil {
		input.TagPredicate = predicateModel.ToUpdateInput()
	} else {
		resp.Diagnostics.AddAttributeError(path.Root("tag_predicate"), "Invalid Attribute Configuration", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdateCheckTagDefined(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_tag_defined, got error: %s", err))
		return
	}

	stateModel := NewCheckTagDefinedResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check tag defined resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckTagDefinedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckTagDefinedResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check tag defined, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check tag defined resource")
}

func (r *CheckTagDefinedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
