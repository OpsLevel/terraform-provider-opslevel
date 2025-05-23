package opslevel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
	"strings"
)

var (
	_ resource.ResourceWithConfigure   = &RelationshipAssignmentResource{}
	_ resource.ResourceWithImportState = &RelationshipAssignmentResource{}
)

func NewRelationshipAssignmentResource() resource.Resource {
	return &RelationshipAssignmentResource{}
}

// RelationshipAssignmentResource defines the resource implementation.
type RelationshipAssignmentResource struct {
	CommonResourceClient
}

// RelationshipAssignmentResourceModel describes the Relationship Assignment managed resource.
type RelationshipAssignmentResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Source     types.String `tfsdk:"source"`
	Target     types.String `tfsdk:"target"`
	Type       types.String `tfsdk:"type"`
	Definition types.String `tfsdk:"definition"`
}

func NewRelationshipAssignmentResourceModel(relationship *opslevel.RelationshipType, givenModel RelationshipAssignmentResourceModel) RelationshipAssignmentResourceModel {
	model := RelationshipAssignmentResourceModel{
		Id:         ComputedStringValue(string(relationship.Id)),
		Source:     givenModel.Source,
		Target:     givenModel.Target,
		Type:       RequiredStringValue(string(relationship.Type)),
		Definition: givenModel.Definition,
	}

	return model
}

func (r *RelationshipAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_assignment"
}

// relationshipDefinitionValidator validates that:
// - definition is required when type is "related_to"
// - definition cannot be set when type is "belongs_to" or "depends_on"
type relationshipDefinitionValidator struct{}

func (v relationshipDefinitionValidator) Description(ctx context.Context) string {
	return "definition is required for 'related_to' type and cannot be set for 'belongs_to' or 'depends_on' types"
}

func (v relationshipDefinitionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v relationshipDefinitionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Get the type value from the config
	var typeValue string
	diags := req.Config.GetAttribute(ctx, path.Root("type"), &typeValue)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Get the definition value
	value := req.ConfigValue.ValueString()

	// Validate based on type
	switch typeValue {
	case "related_to":
		if value == "" && !req.ConfigValue.IsUnknown() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"definition is required when type is 'related_to'",
			)
		}
	case "belongs_to", "depends_on":
		if value != "" {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"definition cannot be set when type is 'belongs_to' or 'depends_on'",
			)
		}
	}
}

func (r *RelationshipAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `A relationship assignment in OpsLevel defines a connection between two resources. It specifies:
- The source resource (where the relationship starts)
- The target resource (where the relationship points to)
- The type of relationship (belongs_to, depends_on, or related_to)
- Optionally, a user defined relationship definition (required for related_to type)

## Notes

- For ` + "`type`" + `:
  - ` + "`belongs_to`" + ` and ` + "`depends_on`" + ` are predefined relationship types that don't require a relationship definition
  - ` + "`related_to`" + ` requires a ` + "`definition`" + ` to be specified
- The relationship definition must be created using the ` + "`opslevel_relationship_definition`" + ` resource first.

`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Description: "The ID of the source resource (where the relationship starts).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target": schema.StringAttribute{
				Description: "The ID of the target resource (where the relationship points to).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of relationship. Must be one of: belongs_to, depends_on, related_to",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllRelationshipTypeEnum...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"definition": schema.StringAttribute{
				Description: "The ID of the relationship definition to use. Required when type is 'related_to'.",
				Optional:    true,
				Validators: []validator.String{
					relationshipDefinitionValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *RelationshipAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[RelationshipAssignmentResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	source := planModel.Source.ValueString()
	target := planModel.Target.ValueString()
	relationshipType := planModel.Type.ValueString()
	relationshipDefinition := planModel.Definition.ValueString()

	input := opslevel.RelationshipDefinition{
		Source: *opslevel.NewIdentifier(source),
		Target: *opslevel.NewIdentifier(target),
		Type:   opslevel.RelationshipTypeEnum(relationshipType),
	}

	if relationshipDefinition != "" {
		input.RelationshipDefinition = opslevel.NewIdentifier(relationshipDefinition)
	}

	relationship, err := r.client.CreateRelationship(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("failed to create relationship from '%s' to '%s', got error: %s", source, target, err))
		return
	}

	stateModel := NewRelationshipAssignmentResourceModel(relationship, planModel)
	tflog.Trace(ctx, fmt.Sprintf("created a relationship from '%s' to '%s'", source, target))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RelationshipAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Since there's no GetRelationship function, we'll just keep the state as is
	// The relationship will be deleted if it doesn't exist when trying to delete it
	stateModel := read[RelationshipAssignmentResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RelationshipAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("terraform plugin error", "relationship assignments cannot be updated, only replaced.\nplease file a bug report including your .tf file at: github.com/OpsLevel/terraform-provider-opslevel")
}

func (r *RelationshipAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[RelationshipAssignmentResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	_, err := r.client.DeleteRelationship(id)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist on this account") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("failed to delete relationship '%s', got error: %s", id, err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("deleted relationship '%s'", id))
}

func (r *RelationshipAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
