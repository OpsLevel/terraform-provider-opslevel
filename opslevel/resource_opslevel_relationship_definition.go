package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

var (
	_ resource.ResourceWithConfigure   = &RelationshipDefinitionResource{}
	_ resource.ResourceWithImportState = &RelationshipDefinitionResource{}
)

func NewRelationshipDefinitionResource() resource.Resource {
	return &RelationshipDefinitionResource{}
}

// RelationshipDefinitionResource defines the resource implementation.
type RelationshipDefinitionResource struct {
	CommonResourceClient
}

// RelationshipDefinitionResourceModel describes the Relationship Definition managed resource.
type RelationshipDefinitionResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Alias         types.String `tfsdk:"alias"`
	Description   types.String `tfsdk:"description"`
	ComponentType types.String `tfsdk:"component_type"`
	Metadata      types.Object `tfsdk:"metadata"`
}

func NewRelationshipDefinitionResourceModel(definition opslevel.RelationshipDefinitionType, givenModel RelationshipDefinitionResourceModel) RelationshipDefinitionResourceModel {
	model := RelationshipDefinitionResourceModel{
		Id:            ComputedStringValue(string(definition.Id)),
		Name:          RequiredStringValue(definition.Name),
		Alias:         RequiredStringValue(definition.Alias),
		Description:   StringValueFromResourceAndModelField(definition.Description, givenModel.Description),
		ComponentType: RequiredStringValue(string(definition.ComponentType.Id)),
		Metadata: types.ObjectValueMust(
			map[string]attr.Type{
				"allowed_types": types.ListType{ElemType: types.StringType},
				"max_items":     types.Int64Type,
				"min_items":     types.Int64Type,
			},
			map[string]attr.Value{
				"allowed_types": types.ListValueMust(types.StringType, func() []attr.Value {
					values := make([]attr.Value, len(definition.Metadata.AllowedTypes))
					for i, v := range definition.Metadata.AllowedTypes {
						values[i] = types.StringValue(v)
					}
					return values
				}()),
				"max_items": types.Int64Value(int64(definition.Metadata.MaxItems)),
				"min_items": types.Int64Value(int64(definition.Metadata.MinItems)),
			},
		),
	}

	return model
}

func (r *RelationshipDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_definition"
}

func (r *RelationshipDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Relationship Definition Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the relationship definition.",
				Required:    true,
			},
			"alias": schema.StringAttribute{
				Description: "The unique identifier of the relationship.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the relationship definition.",
				Optional:    true,
			},
			"component_type": schema.StringAttribute{
				Description: "The component type that the relationship belongs to.",
				Required:    true,
			},
			"metadata": schema.ObjectAttribute{
				Description: "The metadata of the relationship.",
				Required:    true,
				AttributeTypes: map[string]attr.Type{
					"allowed_types": types.ListType{ElemType: types.StringType},
					"max_items":     types.Int64Type,
					"min_items":     types.Int64Type,
				},
			},
		},
	}
}

func (r *RelationshipDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	metadata := planModel.Metadata.Attributes()
	allowedTypes := make([]string, 0)
	if err := metadata["allowed_types"].(types.List).ElementsAs(ctx, &allowedTypes, false); err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_types: %s", err))
		return
	}

	maxItems := int(metadata["max_items"].(types.Int64).ValueInt64())
	minItems := int(metadata["min_items"].(types.Int64).ValueInt64())

	input := opslevel.RelationshipDefinitionInput{
		Name:          nullable(planModel.Name.ValueStringPointer()),
		Alias:         nullable(planModel.Alias.ValueStringPointer()),
		Description:   nullable(planModel.Description.ValueStringPointer()),
		ComponentType: opslevel.NewIdentifier(planModel.ComponentType.ValueString()),
		Metadata: &opslevel.RelationshipDefinitionMetadataInput{
			AllowedTypes: allowedTypes,
			MaxItems:     &maxItems,
			MinItems:     &minItems,
		},
	}

	definition, err := r.client.CreateRelationshipDefinition(input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create relationship definition with name '%s', got error: %s", input.Name.Value, err))
		return
	}

	stateModel := NewRelationshipDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("created a relationship definition resource with id '%s'", definition.Id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RelationshipDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	definition, err := r.client.GetRelationshipDefinition(id)
	if err != nil {
		if (definition == nil || definition.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read relationship definition with id '%s', got error: %s", id, err))
		return
	}

	verifiedStateModel := NewRelationshipDefinitionResourceModel(*definition, stateModel)
	tflog.Trace(ctx, fmt.Sprintf("read a relationship definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *RelationshipDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	metadata := planModel.Metadata.Attributes()
	allowedTypes := make([]string, 0)
	if err := metadata["allowed_types"].(types.List).ElementsAs(ctx, &allowedTypes, false); err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to parse allowed_types: %s", err))
		return
	}

	maxItems := int(metadata["max_items"].(types.Int64).ValueInt64())
	minItems := int(metadata["min_items"].(types.Int64).ValueInt64())

	id := planModel.Id.ValueString()
	input := opslevel.RelationshipDefinitionInput{
		Name:          nullable(planModel.Name.ValueStringPointer()),
		Alias:         nullable(planModel.Alias.ValueStringPointer()),
		Description:   nullable(planModel.Description.ValueStringPointer()),
		ComponentType: opslevel.NewIdentifier(planModel.ComponentType.ValueString()),
		Metadata: &opslevel.RelationshipDefinitionMetadataInput{
			AllowedTypes: allowedTypes,
			MaxItems:     &maxItems,
			MinItems:     &minItems,
		},
	}

	definition, err := r.client.UpdateRelationshipDefinition(id, input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update relationship definition with id '%s', got error: %s", id, err))
		return
	}

	stateModel := NewRelationshipDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("updated a relationship definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RelationshipDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[RelationshipDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	err := r.client.DeleteRelationshipDefinition(id)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete relationship definition (%s), got error: %s", id, err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("deleted a relationship definition resource with id '%s'", id))
}

func (r *RelationshipDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
} 