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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure      = &CheckServiceOwnershipResource{}
	_ resource.ResourceWithImportState    = &CheckServiceOwnershipResource{}
	_ resource.ResourceWithUpgradeState   = &CheckServiceOwnershipResource{}
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
				Description: "Check Service Ownership Resource",
				Attributes: getCheckBaseSchemaV0(map[string]schema.Attribute{
					"contact_method": schema.StringAttribute{
						Description: "The type of contact method that is required.",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive(enumAllContactTypes...)},
					},
					"id": schema.StringAttribute{
						Description: "The ID of this resource.",
						Computed:    true,
					},
					"require_contact_method": schema.BoolAttribute{
						Description: "True if a service's owner must have a contact method, False otherwise.",
						Optional:    true,
					},
					"tag_key": schema.StringAttribute{
						Description: "The tag key where the tag predicate should be applied.",
						Optional:    true,
					},
				}),
				Blocks: map[string]schema.Block{
					"tag_predicate": schema.ListNestedBlock{
						NestedObject: predicateSchemaV0,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var diags diag.Diagnostics
				upgradedStateModel := CheckServiceOwnershipResourceModel{}
				tagPredicateList := types.ListNull(types.ObjectType{AttrTypes: predicateType})

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

				// service ownership specific attributes
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("contact_method"), &upgradedStateModel.ContactMethod)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("require_contact_method"), &upgradedStateModel.RequireContactMethod)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tag_key"), &upgradedStateModel.TagKey)...)
				resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tag_predicate"), &tagPredicateList)...)
				if len(tagPredicateList.Elements()) == 1 {
					tagPredicate := tagPredicateList.Elements()[0]
					upgradedStateModel.TagPredicate, diags = types.ObjectValueFrom(ctx, predicateType, tagPredicate)
					resp.Diagnostics.Append(diags...)
				} else {
					upgradedStateModel.TagPredicate = types.ObjectNull(predicateType)
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateModel)...)
			},
		},
	}
}

func (r *CheckServiceOwnershipResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	tagPredicate := types.ObjectNull(predicateType)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("tag_predicate"), &tagPredicate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicateModel, diags := PredicateObjectToModel(ctx, tagPredicate)
	resp.Diagnostics.Append(diags...)
	if err := predicateModel.Validate(); err != nil {
		resp.Diagnostics.AddAttributeWarning(path.Root("tag_predicate"), "Invalid Attribute Configuration", err.Error())
	}
}

func (r *CheckServiceOwnershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckServiceOwnershipResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckServiceOwnershipCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	input.RequireContactMethod = nullable(planModel.RequireContactMethod.ValueBoolPointer())
	if planModel.ContactMethod.ValueString() != "" {
		input.ContactMethod = opslevel.RefOf(strings.ToUpper(planModel.ContactMethod.ValueString()))
	}
	input.TagKey = nullable(planModel.TagKey.ValueStringPointer())

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
	stateModel := read[CheckServiceOwnershipResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check service ownership, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckServiceOwnershipResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckServiceOwnershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckServiceOwnershipResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckServiceOwnershipUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		Id:         asID(planModel.Id),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	input.RequireContactMethod = nullable(planModel.RequireContactMethod.ValueBoolPointer())
	input.ContactMethod = opslevel.RefOf(strings.ToUpper(planModel.ContactMethod.ValueString()))
	input.TagKey = nullable(planModel.TagKey.ValueStringPointer())

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
	stateModel := read[CheckServiceOwnershipResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check service ownership, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check service ownership resource")
}

func (r *CheckServiceOwnershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
