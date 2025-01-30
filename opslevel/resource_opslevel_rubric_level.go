package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &RubricLevelResource{}

var _ resource.ResourceWithImportState = &RubricLevelResource{}

func NewRubricLevelResource() resource.Resource {
	return &RubricLevelResource{}
}

// RubricLevelResource defines the resource implementation.
type RubricLevelResource struct {
	CommonResourceClient
}

// RubricLevelResourceModel describes the rubric level managed resource.
type RubricLevelResourceModel struct {
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Index       types.Int64  `tfsdk:"index"`
	Name        types.String `tfsdk:"name"`
}

func NewRubricLevelResourceModel(rubricLevel opslevel.Level, givenModel RubricLevelResourceModel) RubricLevelResourceModel {
	return RubricLevelResourceModel{
		Description: StringValueFromResourceAndModelField(rubricLevel.Description, givenModel.Description),
		Id:          ComputedStringValue(string(rubricLevel.Id)),
		Index:       types.Int64Value(int64(rubricLevel.Index)),
		Name:        RequiredStringValue(rubricLevel.Name),
	}
}

func (r *RubricLevelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_level"
}

func (r *RubricLevelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Rubric Level Resource",

		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "The description of the rubric level.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the rubric level.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"index": schema.Int64Attribute{
				Description: "An integer allowing this level to be inserted between others. Must be unique per rubric.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the rubric level.",
				Required:    true,
			},
		},
	}
}

func (r *RubricLevelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[RubricLevelResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	levelCreateInput := opslevel.LevelCreateInput{
		Name:        planModel.Name.ValueString(),
		Description: nullable(planModel.Description.ValueStringPointer()),
	}
	if !planModel.Index.IsNull() && !planModel.Index.IsUnknown() {
		index := int(planModel.Index.ValueInt64())
		levelCreateInput.Index = &index
	}
	rubricLevel, err := r.client.CreateLevel(levelCreateInput)
	if err != nil || rubricLevel == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create rubric level, got error: %s", err))
		return
	}

	createdRubricLevelResourceModel := NewRubricLevelResourceModel(*rubricLevel, planModel)

	tflog.Trace(ctx, "created a rubric level resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdRubricLevelResourceModel)...)
}

func (r *RubricLevelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[RubricLevelResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	rubricLevel, err := r.client.GetLevel(asID(stateModel.Id))
	if err != nil {
		if (rubricLevel == nil || rubricLevel.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read rubric level, got error: %s", err))
		return
	}
	readRubricLevelResourceModel := NewRubricLevelResourceModel(*rubricLevel, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readRubricLevelResourceModel)...)
}

func (r *RubricLevelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[RubricLevelResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedRubricLevel, err := r.client.UpdateLevel(opslevel.LevelUpdateInput{
		Description: opslevel.RefOf(planModel.Description.ValueString()),
		Id:          opslevel.ID(planModel.Id.ValueString()),
		Name:        nullable(planModel.Name.ValueStringPointer()),
	})
	if err != nil || updatedRubricLevel == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update rubric level, got error: %s", err))
		return
	}
	updatedRubricLevelResourceModel := NewRubricLevelResourceModel(*updatedRubricLevel, planModel)

	tflog.Trace(ctx, "updated a rubric level resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedRubricLevelResourceModel)...)
}

func (r *RubricLevelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[RubricLevelResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLevel(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rubric level, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a rubric level resource")
}

func (r *RubricLevelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
