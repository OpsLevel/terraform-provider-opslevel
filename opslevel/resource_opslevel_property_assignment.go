package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &PropertyAssignmentResource{}

type PropertyAssignmentResource struct {
	CommonResourceClient
}

func NewPropertyAssignmentResource() resource.Resource {
	return &PropertyAssignmentResource{}
}

type PropertyAssignmentResourceModel struct {
	Definition types.String `tfsdk:"definition"`
	Id         types.String `tfsdk:"id"`
	Locked     types.Bool   `tfsdk:"locked"`
	Owner      types.String `tfsdk:"owner"`
	Value      types.String `tfsdk:"value"`
}

func NewPropertyAssignmentResourceModel(assignment opslevel.Property) PropertyAssignmentResourceModel {
	model := PropertyAssignmentResourceModel{
		Locked: types.BoolValue(assignment.Locked),
		Value:  types.StringValue(string(*assignment.Value)),
	}
	// TODO: do we need to keep using this method of setting an ID in the new plugin version?
	// the API does not have unique ID's for property assignments, so what we did in the past was use <owner_id>:<definition_id>
	// keeping this for now just in case it's necessary for backwards compatability.
	model.Id = RequiredStringValue(fmt.Sprintf("%s:%s", assignment.Owner.Id(), assignment.Definition.Id))

	return model
}

func (resource *PropertyAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_assignment"
}

func (resource *PropertyAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property Assignment Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"locked": schema.BoolAttribute{
				Description: "If locked = true, the property has been set in opslevel.yml and cannot be modified in Terraform!",
				Computed:    true,
			},
			"definition": schema.StringAttribute{
				Description: "The custom property definition's ID or alias.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The ID or alias of the entity (currently only supports service) that the property has been assigned to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "The value of the custom property (must be a valid JSON value or null or object).",
				Optional:    true,
				Validators: []validator.String{
					JsonStringValidator(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (resource *PropertyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel PropertyAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition := planModel.Definition.ValueString()
	owner := planModel.Owner.ValueString()
	value := opslevel.JsonString(planModel.Value.ValueString())
	input := opslevel.PropertyInput{
		Definition: *opslevel.NewIdentifier(planModel.Definition.ValueString()),
		Owner:      *opslevel.NewIdentifier(planModel.Owner.ValueString()),
		Value:      value,
	}
	assignment, err := resource.client.PropertyAssign(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("failed to assign property (%s) on service (%s), got error: %s", definition, owner, err))
		return
	}

	stateModel := NewPropertyAssignmentResourceModel(*assignment)
	// user is free to use either alias or ID for 'owner' and 'definition' fields
	stateModel.Owner = planModel.Owner
	stateModel.Definition = planModel.Definition

	tflog.Trace(ctx, fmt.Sprintf("assigned property (%s) on service (%s) with value: '%s'", definition, owner, value))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (resource *PropertyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel PropertyAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition := planModel.Definition.ValueString()
	owner := planModel.Owner.ValueString()
	assignment, err := resource.client.GetProperty(owner, definition)
	if err != nil || assignment == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read property assignment (%s) on service (%s), got error: %s", definition, owner, err))
		return
	}
	value := *assignment.Value

	stateModel := NewPropertyAssignmentResourceModel(*assignment)
	// user is free to use either alias or ID for 'owner' and 'definition' fields
	stateModel.Owner = planModel.Owner
	stateModel.Definition = planModel.Definition

	tflog.Trace(ctx, fmt.Sprintf("read property assignment (%s) on service (%s) with value: '%s'", definition, owner, value))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (resource *PropertyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("terraform plugin error", "property assignments should never be updated, only replaced.\nplease file a bug report including your .tf file at: github.com/OpsLevel/terraform-provider-opslevel")
}

func (resource *PropertyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel PropertyAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition := planModel.Definition.ValueString()
	owner := planModel.Owner.ValueString()
	err := resource.client.PropertyUnassign(owner, definition)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("failed to unassign property (%s) on service (%s), got error: %s", definition, owner, err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("unassigned property (%s) on service (%s)", definition, owner))
}
