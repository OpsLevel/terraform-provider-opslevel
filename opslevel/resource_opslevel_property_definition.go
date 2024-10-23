package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var (
	_ resource.ResourceWithConfigure   = &PropertyDefinitionResource{}
	_ resource.ResourceWithImportState = &PropertyDefinitionResource{}
)

type PropertyDefinitionResource struct {
	CommonResourceClient
}

func NewPropertyDefinitionResource() resource.Resource {
	return &PropertyDefinitionResource{}
}

// PropertyDefinitionResourceModel describes the Property Definition managed resource.
type PropertyDefinitionResourceModel struct {
	AllowedInConfigFiles  types.Bool   `tfsdk:"allowed_in_config_files"`
	Description           types.String `tfsdk:"description"`
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	PropertyDisplayStatus types.String `tfsdk:"property_display_status"`
	Schema                types.String `tfsdk:"schema"`
}

func NewPropertyDefinitionResourceModel(definition opslevel.PropertyDefinition, givenModel PropertyDefinitionResourceModel) PropertyDefinitionResourceModel {
	model := PropertyDefinitionResourceModel{
		AllowedInConfigFiles:  types.BoolValue(definition.AllowedInConfigFiles),
		Description:           StringValueFromResourceAndModelField(definition.Description, givenModel.Description),
		Id:                    ComputedStringValue(string(definition.Id)),
		Name:                  RequiredStringValue(definition.Name),
		PropertyDisplayStatus: RequiredStringValue(string(definition.PropertyDisplayStatus)),
		Schema:                RequiredStringValue(definition.Schema.ToJSON()),
	}

	return model
}

func (resource *PropertyDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_definition"
}

func (resource *PropertyDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property Definition Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_in_config_files": schema.BoolAttribute{
				Description: "Whether or not the property is allowed to be set in opslevel.yml config files.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the property definition.",
				Required:    true,
			},
			"schema": schema.StringAttribute{
				Description: "The schema of the property definition.",
				Required:    true,
				Validators: []validator.String{
					JsonObjectValidator(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the property definition.",
				Optional:    true,
			},
			"property_display_status": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The display status of a custom property on service pages. One of `%s`",
					strings.Join(opslevel.AllPropertyDisplayStatusEnum, "`, `"),
				),
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllPropertyDisplayStatusEnum...),
				},
			},
		},
	}
}

func (resource *PropertyDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[PropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	definitionSchema, err := opslevel.NewJSONSchema(planModel.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to use definition schema '%s', got error: %s", planModel.Schema.ValueString(), err))
		return
	}

	var propertyDisplayStatus *opslevel.PropertyDisplayStatusEnum
	if planModel.PropertyDisplayStatus.ValueString() != "" {
		propertyDisplayStatus = opslevel.RefOf(opslevel.PropertyDisplayStatusEnum(planModel.PropertyDisplayStatus.ValueString()))
	}
	input := opslevel.PropertyDefinitionInput{
		AllowedInConfigFiles:  planModel.AllowedInConfigFiles.ValueBoolPointer(),
		Description:           planModel.Description.ValueStringPointer(),
		Name:                  planModel.Name.ValueStringPointer(),
		PropertyDisplayStatus: propertyDisplayStatus,
		Schema:                definitionSchema,
	}
	definition, err := resource.client.CreatePropertyDefinition(input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create definition with name '%s', got error: %s", *input.Name, err))
		return
	}

	stateModel := NewPropertyDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("created a definition resource with id '%s'", definition.Id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (resource *PropertyDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[PropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	definition, err := resource.client.GetPropertyDefinition(id)
	if err != nil {
		if (definition == nil || definition.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read definition with id '%s', got error: %s", id, err))
		return
	}

	verifiedStateModel := NewPropertyDefinitionResourceModel(*definition, stateModel)
	tflog.Trace(ctx, fmt.Sprintf("read a definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (resource *PropertyDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[PropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	definitionSchema, err := opslevel.NewJSONSchema(planModel.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("config error", fmt.Sprintf("unable to use definition schema '%s', got error: %s", planModel.Schema.ValueString(), err))
		return
	}

	id := planModel.Id.ValueString()
	var propertyDisplayStatus *opslevel.PropertyDisplayStatusEnum
	if planModel.PropertyDisplayStatus.ValueString() != "" {
		propertyDisplayStatus = opslevel.RefOf(opslevel.PropertyDisplayStatusEnum(planModel.PropertyDisplayStatus.ValueString()))
	}
	input := opslevel.PropertyDefinitionInput{
		AllowedInConfigFiles:  planModel.AllowedInConfigFiles.ValueBoolPointer(),
		Description:           opslevel.RefOf(planModel.Description.ValueString()),
		Name:                  planModel.Name.ValueStringPointer(),
		PropertyDisplayStatus: propertyDisplayStatus,
		Schema:                definitionSchema,
	}
	definition, err := resource.client.UpdatePropertyDefinition(id, input)
	if err != nil || definition == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update definition with id '%s', got error: %s", id, err))
		return
	}

	stateModel := NewPropertyDefinitionResourceModel(*definition, planModel)
	tflog.Trace(ctx, fmt.Sprintf("updated a definition resource with id '%s'", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (resource *PropertyDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[PropertyDefinitionResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	id := stateModel.Id.ValueString()
	err := resource.client.DeletePropertyDefinition(id)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to delete definition (%s), got error: %s", id, err))
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("deleted a definition resource with id '%s'", id))
}

func (r *PropertyDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
