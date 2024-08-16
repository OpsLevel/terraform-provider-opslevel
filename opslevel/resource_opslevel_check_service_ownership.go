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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure      = &CheckServiceOwnershipResource{}
	_ resource.ResourceWithImportState    = &CheckServiceOwnershipResource{}
	_ resource.ResourceWithValidateConfig = &CheckServiceOwnershipResource{}
)

func NewCheckServiceOwnershipResource() resource.Resource {
	return &CheckServiceOwnershipResource{}
}

// CheckServiceOwnershipResource defines the resource implementation.
type CheckServiceOwnershipResource struct {
	CommonResourceClient
}

type CheckServiceOwnershipResourceModel struct {
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

	RequireContactMethod types.Bool   `tfsdk:"require_contact_method"`
	ContactMethod        types.String `tfsdk:"contact_method"`
	TagKey               types.String `tfsdk:"tag_key"`
	TagPredicate         types.Object `tfsdk:"tag_predicate"`
}

func NewCheckServiceOwnershipResourceModel(ctx context.Context, check opslevel.Check, planModel CheckServiceOwnershipResourceModel) CheckServiceOwnershipResourceModel {
	var stateModel CheckServiceOwnershipResourceModel

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
	stateModel.RequireContactMethod = OptionalBoolValue(check.ServiceOwnershipCheckFragment.RequireContactMethod)

	if check.ServiceOwnershipCheckFragment.ContactMethod != nil {
		contactMethod := string(*check.ServiceOwnershipCheckFragment.ContactMethod)
		if strings.ToLower(planModel.ContactMethod.ValueString()) == strings.ToLower(contactMethod) {
			stateModel.ContactMethod = planModel.ContactMethod
		} else {
			stateModel.ContactMethod = OptionalStringValue(contactMethod)
		}
	}
	stateModel.TagKey = OptionalStringValue(check.ServiceOwnershipCheckFragment.TeamTagKey)

	if check.ServiceOwnershipCheckFragment.TeamTagPredicate == nil {
		stateModel.TagPredicate = types.ObjectNull(predicateType)
	} else {
		predicate := *&check.ServiceOwnershipCheckFragment.TeamTagPredicate
		predicateAttrValues := map[string]attr.Value{
			"type":  types.StringValue(string(predicate.Type)),
			"value": OptionalStringValue(predicate.Value),
		}
		stateModel.TagPredicate = types.ObjectValueMust(predicateType, predicateAttrValues)
	}

	return stateModel
}

func (r *CheckServiceOwnershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_service_ownership"
}

func (r *CheckServiceOwnershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	enumAllContactTypes := append(opslevel.AllContactType, "any")
	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Service Ownership Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"require_contact_method": schema.BoolAttribute{
				Description: "True if a service's owner must have a contact method, False otherwise.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"contact_method": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The type of contact method that is required. One of `%s`",
					strings.Join(enumAllContactTypes, "`, `"),
				),
				Computed:   true,
				Optional:   true,
				Default:    stringdefault.StaticString("ANY"),
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(enumAllContactTypes...)},
			},
			"tag_key": schema.StringAttribute{
				Description: "The tag key where the tag predicate should be applied.",
				Optional:    true,
			},
			"tag_predicate": PredicateSchema(),
		}),
	}
}

func (r *CheckServiceOwnershipResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	enumAllContactTypes := append(opslevel.AllContactType, "any")
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: CheckBaseAttributes(map[string]schema.Attribute{
					"require_contact_method": schema.BoolAttribute{
						Description: "True if a service's owner must have a contact method, False otherwise.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"contact_method": schema.StringAttribute{
						Description: fmt.Sprintf(
							"The type of contact method that is required. One of `%s`",
							strings.Join(enumAllContactTypes, "`, `"),
						),
						Computed:   true,
						Optional:   true,
						Default:    stringdefault.StaticString("ANY"),
						Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(enumAllContactTypes...)},
					},
					"tag_key": schema.StringAttribute{
						Description: "The tag key where the tag predicate should be applied.",
						Optional:    true,
					},
					"tag_predicate": PredicateSchema(),
				}),
			},
			StateUpgrader: CheckUpgradeFunc[CheckServiceOwnershipResourceModel](),
		},
	}
}

func (r *CheckServiceOwnershipResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configModel CheckServiceOwnershipResourceModel
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

func (r *CheckServiceOwnershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckServiceOwnershipResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckServiceOwnershipCreateInput{
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

	input.RequireContactMethod = planModel.RequireContactMethod.ValueBoolPointer()
	if planModel.ContactMethod.ValueString() != "" {
		input.ContactMethod = opslevel.RefOf(strings.ToUpper(planModel.ContactMethod.ValueString()))
	}
	input.TagKey = planModel.TagKey.ValueStringPointer()

	// convert tool_name_predicate object to model from plan
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

	data, err := r.client.CreateCheckServiceOwnership(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_service_ownership, got error: %s", err))
		return
	}

	stateModel := NewCheckServiceOwnershipResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check service ownership resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckServiceOwnershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckServiceOwnershipResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check service ownership, got error: %s", err))
		return
	}
	stateModel := NewCheckServiceOwnershipResourceModel(ctx, *data, planModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckServiceOwnershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckServiceOwnershipResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckServiceOwnershipUpdateInput{
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

	input.RequireContactMethod = planModel.RequireContactMethod.ValueBoolPointer()
	input.ContactMethod = opslevel.RefOf(strings.ToUpper(planModel.ContactMethod.ValueString()))
	input.TagKey = planModel.TagKey.ValueStringPointer()

	// convert tool_name_predicate object to model from plan
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

	data, err := r.client.UpdateCheckServiceOwnership(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_service_ownership, got error: %s", err))
		return
	}

	stateModel := NewCheckServiceOwnershipResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check service ownership resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckServiceOwnershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckServiceOwnershipResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check service ownership, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check service ownership resource")
}

func (r *CheckServiceOwnershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
