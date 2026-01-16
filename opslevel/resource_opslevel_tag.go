package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
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

	Id types.String `tfsdk:"id"`
}

func NewTagResourceModel(ctx context.Context, tag opslevel.Tag, planModel TagResourceModel) TagResourceModel {
	var stateModel TagResourceModel

	stateModel.TargetResource = RequiredStringValue(planModel.TargetResource.ValueString())
	stateModel.TargetType = RequiredStringValue(planModel.TargetType.ValueString())
	stateModel.Key = RequiredStringValue(tag.Key)
	stateModel.Value = RequiredStringValue(tag.Value)
	stateModel.Id = types.StringValue(string(tag.Id))

	return stateModel
}

func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tag Resource",

		Attributes: map[string]schema.Attribute{
			"resource_identifier": schema.StringAttribute{
				Description: "The id or human-friendly, unique identifier of the resource this tag belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The resource type that the tag applies to. One of `%s`",
					strings.Join(opslevel.AllTaggableResource, "`, `"),
				),
				Required: true,
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
		},
	}
}

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[TagResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := planModel.TargetResource.ValueString()
	resourceType := opslevel.TaggableResource(planModel.TargetType.ValueString())
	tagCreateInput := opslevel.TagCreateInput{
		Key:   planModel.Key.ValueString(),
		Value: planModel.Value.ValueString(),
		Type:  &resourceType,
	}

	if opslevel.IsID(resourceId) {
		tagCreateInput.Id = opslevel.NewID(resourceId)
	} else {
		tagCreateInput.Alias = &resourceId
	}
	data, err := r.client.CreateTag(tagCreateInput)
	if err != nil {
		title, detail := formatOpslevelError("create tag", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}
	stateModel := NewTagResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[TagResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := stateModel.TargetResource.ValueString()
	resourceType := opslevel.TaggableResource(stateModel.TargetType.ValueString())
	data, err := r.client.GetTaggableResource(resourceType, resourceId)
	if err != nil {
		// If the parent resource is gone, remove the tag from state
		if opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		title, detail := formatOpslevelError("read tag", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}
	tags, err := data.GetTags(r.client, nil)
	if err != nil {
		title, detail := formatOpslevelError(fmt.Sprintf("get tags from '%s' with id '%s'", resourceType, resourceId), err)
		resp.Diagnostics.AddError(title, detail)
		return
	}
	if tags == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	id := stateModel.Id.ValueString()
	tag, err := tags.GetTagById(*opslevel.NewID(id))
	if err != nil {
		if tag == nil || tag.Id == "" {
			resp.State.RemoveResource(ctx)
			return
		}
		title, detail := formatOpslevelError(
			fmt.Sprintf("find tag '%s' for type %s with id '%s'", id, data.ResourceType(), data.ResourceId()),
			err,
		)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	verifiedStateModel := NewTagResourceModel(ctx, *tag, stateModel)

	// Save updated planModel into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[TagResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	tagUpdateInput := opslevel.TagUpdateInput{
		Id:    asID(planModel.Id),
		Key:   planModel.Key.ValueStringPointer(),
		Value: planModel.Value.ValueStringPointer(),
	}
	updatedTag, err := r.client.UpdateTag(tagUpdateInput)
	if err != nil {
		title, detail := formatOpslevelError(fmt.Sprintf("update tag with id '%s'", planModel.Id.ValueString()), err)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	stateModel := NewTagResourceModel(ctx, *updatedTag, planModel)

	tflog.Trace(ctx, "updated a tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[TagResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTag(asID(data.Id))
	if err != nil {
		title, detail := formatOpslevelError("delete tag", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}
	tflog.Trace(ctx, "deleted a tag resource")
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
