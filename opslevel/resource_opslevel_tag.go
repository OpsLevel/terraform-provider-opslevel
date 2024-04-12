package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &TagResource{}

var _ resource.ResourceWithImportState = &TagResource{}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the resource implementation.
type TagResource struct {
	CommonResourceClient
}

// TagResourceModel describes the Domain managed resource.
type TagResourceModel struct {
	TargetResource types.String `tfsdk:"resource_identifier"`
	TargetType     types.String `tfsdk:"resource_type"`
	Key            types.String `tfsdk:"key"`
	Value          types.String `tfsdk:"value"`

	Id          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func NewTagResourceModel(ctx context.Context, tag opslevel.Tag, planModel TagResourceModel) (TagResourceModel, diag.Diagnostics) {
	var stateModel TagResourceModel

	stateModel.TargetResource = RequiredStringValue(planModel.TargetResource.ValueString())
	stateModel.TargetType = RequiredStringValue(planModel.TargetType.ValueString())
	stateModel.Key = RequiredStringValue(tag.Key)
	stateModel.Value = RequiredStringValue(tag.Value)
	stateModel.Id = types.StringValue(string(tag.Id))

	return stateModel, diag.Diagnostics{}
}

func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Domain Resource",

		Attributes: map[string]schema.Attribute{
			"resource_identifier": schema.StringAttribute{
				Description: "The id or human-friendly, unique identifier of the resource this tag belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Description: "The resource type that the tag applies to.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllTaggableResource...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The key of the tag.",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The value of the tag.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the tag created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel TagResourceModel

	// Read Terraform plan planModel into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := planModel.TargetResource.ValueString()
	resourceType := opslevel.TaggableResource(planModel.TargetType.ValueString())
	tagCreateInput := opslevel.TagCreateInput{
		Key:   planModel.Key.ValueString(),
		Value: planModel.Value.ValueString(),
		Type:  opslevel.RefOf(resourceType),
	}

	if opslevel.IsID(resourceId) {
		tagCreateInput.Id = opslevel.NewID(resourceId)
	} else {
		tagCreateInput.Alias = opslevel.RefOf(resourceId)
	}
	data, err := r.client.CreateTag(tagCreateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create tag, got error: %s", err))
		return
	}
	stateModel, diags := NewTagResourceModel(ctx, *data, planModel)
	resp.Diagnostics.Append(diags...)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel TagResourceModel

	// Read Terraform prior state planModel into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := planModel.TargetResource.ValueString()
	resourceType := opslevel.TaggableResource(planModel.TargetType.ValueString())
	data, err := r.client.GetTaggableResource(resourceType, resourceId)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read tag, got error: %s", err))
		return
	}
	tags, err := data.GetTags(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get tags from '%s' with id '%s'", resourceType, resourceId))
	}

	id := planModel.Id.ValueString()
	tag, err := tags.GetTagById(*opslevel.NewID(id))
	if err != nil || tag == nil {
		resp.Diagnostics.AddError("opslevel client error",
			fmt.Sprintf("Tag '%s' for type %s with id '%s' not found. %s",
				id,
				data.ResourceType(),
				data.ResourceId(),
				err,
			))
		return
	}

	stateModel, diags := NewTagResourceModel(ctx, *tag, planModel)
	resp.Diagnostics.Append(diags...)

	// Save updated planModel into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("terraform plugin error", "tag assignments should never be updated, only replaced.\nplease file a bug report including your .tf file at: github.com/OpsLevel/terraform-provider-opslevel")
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TagResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTag(asID(data.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tag, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a tag resource")
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
