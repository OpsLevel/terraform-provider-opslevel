package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckRelationshipResource{}
	_ resource.ResourceWithImportState = &CheckRelationshipResource{}
)

func NewCheckRelationshipResource() resource.Resource {
	return &CheckRelationshipResource{}
}

// CheckRelationshipResource defines the resource implementation.
type CheckRelationshipResource struct {
	CommonResourceClient
}

type CheckRelationshipResourceModel struct {
	Category                   types.String `tfsdk:"category"`
	Description                types.String `tfsdk:"description"`
	Enabled                    types.Bool   `tfsdk:"enabled"`
	EnableOn                   types.String `tfsdk:"enable_on"`
	Filter                     types.String `tfsdk:"filter"`
	Id                         types.String `tfsdk:"id"`
	Level                      types.String `tfsdk:"level"`
	Name                       types.String `tfsdk:"name"`
	Notes                      types.String `tfsdk:"notes"`
	Owner                      types.String `tfsdk:"owner"`
	RelationshipCountPredicate types.Object `tfsdk:"relationship_count_predicate"`
	RelationshipDefinitionId   types.String `tfsdk:"relationship_definition_id"`
}

func NewCheckRelationshipResourceModel(ctx context.Context, check opslevel.Check, planModel CheckRelationshipResourceModel) CheckRelationshipResourceModel {
	var stateModel CheckRelationshipResourceModel

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

	// Handle relationship count predicate
	if check.RelationshipCheckFragment.RelationshipCountPredicate != nil {
		predicate := check.RelationshipCheckFragment.RelationshipCountPredicate
		predicateModel := PredicateModel{
			Type:  types.StringValue(string(predicate.Type)),
			Value: OptionalStringValue(predicate.Value),
		}
		stateModel.RelationshipCountPredicate = types.ObjectValueMust(
			predicateType,
			map[string]attr.Value{
				"type":  predicateModel.Type,
				"value": predicateModel.Value,
			},
		)
	} else {
		stateModel.RelationshipCountPredicate = types.ObjectNull(predicateType)
	}

	stateModel.RelationshipDefinitionId = RequiredStringValue(string(check.RelationshipCheckFragment.RelationshipDefinition.Id))

	return stateModel
}

func (r *CheckRelationshipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_relationship"
}

func (r *CheckRelationshipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Relationship Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"relationship_count_predicate": PredicateSchema(),
			"relationship_definition_id": schema.StringAttribute{
				Description: "Count relationships of a specific relationship definition.",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
		}),
	}
}

func (r *CheckRelationshipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckRelationshipResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse relationship count predicate
	var predicateInput opslevel.PredicateInput
	if !planModel.RelationshipCountPredicate.IsNull() {
		predicateModel, diags := PredicateObjectToModel(ctx, planModel.RelationshipCountPredicate)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		predicateInput = *predicateModel.ToCreateInput()
	}

	input := opslevel.CheckRelationshipCreateInput{
		CategoryId:                 asID(planModel.Category),
		Enabled:                    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:                   nullableID(planModel.Filter.ValueStringPointer()),
		LevelId:                    asID(planModel.Level),
		Name:                       planModel.Name.ValueString(),
		Notes:                      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:                    nullableID(planModel.Owner.ValueStringPointer()),
		RelationshipCountPredicate: predicateInput,
		RelationshipDefinitionId:   asID(planModel.RelationshipDefinitionId),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	data, err := r.client.CreateCheckRelationship(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_relationship, got error: %s", err))
		return
	}

	stateModel := NewCheckRelationshipResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check relationship resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRelationshipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckRelationshipResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check relationship, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckRelationshipResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckRelationshipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckRelationshipResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse relationship count predicate
	var predicateInput *opslevel.PredicateInput
	if !planModel.RelationshipCountPredicate.IsNull() {
		predicateModel, diags := PredicateObjectToModel(ctx, planModel.RelationshipCountPredicate)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		predicateInput = predicateModel.ToCreateInput()
	}

	input := opslevel.CheckRelationshipUpdateInput{
		CategoryId:                 opslevel.RefOf(asID(planModel.Category)),
		Enabled:                    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:                   nullableID(planModel.Filter.ValueStringPointer()),
		Id:                         asID(planModel.Id),
		LevelId:                    opslevel.RefOf(asID(planModel.Level)),
		Name:                       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:                      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:                    nullableID(planModel.Owner.ValueStringPointer()),
		RelationshipCountPredicate: predicateInput,
		RelationshipDefinitionId:   opslevel.NewID(planModel.RelationshipDefinitionId.ValueString()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	data, err := r.client.UpdateCheckRelationship(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_relationship, got error: %s", err))
		return
	}

	stateModel := NewCheckRelationshipResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check relationship resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckRelationshipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckRelationshipResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check relationship, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check relationship resource")
}

func (r *CheckRelationshipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
