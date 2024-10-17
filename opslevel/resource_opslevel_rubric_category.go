package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &RubricCategoryResource{}

var _ resource.ResourceWithImportState = &RubricCategoryResource{}

func NewRubricCategoryResource() resource.Resource {
	return &RubricCategoryResource{}
}

// RubricCategoryResource defines the resource implementation.
type RubricCategoryResource struct {
	CommonResourceClient
}

// RubricCategoryResourceModel describes the rubric category managed resource.
type RubricCategoryResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewRubricCategoryResourceModel(rubricCategory opslevel.Category) RubricCategoryResourceModel {
	return RubricCategoryResourceModel{
		Id:   ComputedStringValue(string(rubricCategory.Id)),
		Name: RequiredStringValue(rubricCategory.Name),
	}
}

func (r *RubricCategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_category"
}

func (r *RubricCategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Rubric Category Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the rubric category.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the rubric category.",
				Required:    true,
			},
		},
	}
}

func (r *RubricCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RubricCategoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rubricCategory, err := r.client.CreateCategory(opslevel.CategoryCreateInput{
		Name: data.Name.ValueString(),
	})
	if err != nil || rubricCategory == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create rubric category, got error: %s", err))
		return
	}

	createdRubricCategoryResourceModel := NewRubricCategoryResourceModel(*rubricCategory)

	tflog.Trace(ctx, "created a rubric category resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdRubricCategoryResourceModel)...)
}

func (r *RubricCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RubricCategoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rubricCategory, err := r.client.GetCategory(asID(data.Id))
	if err != nil {
		if (rubricCategory == nil || rubricCategory.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read rubric category, got error: %s", err))
		return
	}
	readRubricCategoryResourceModel := NewRubricCategoryResourceModel(*rubricCategory)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readRubricCategoryResourceModel)...)
}

func (r *RubricCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RubricCategoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedRubricCategory, err := r.client.UpdateCategory(opslevel.CategoryUpdateInput{
		Id:   opslevel.ID(data.Id.ValueString()),
		Name: data.Name.ValueStringPointer(),
	})
	if err != nil || updatedRubricCategory == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update rubric category, got error: %s", err))
		return
	}
	updatedRubricCategoryResourceModel := NewRubricCategoryResourceModel(*updatedRubricCategory)

	tflog.Trace(ctx, "updated a rubric category resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedRubricCategoryResourceModel)...)
}

func (r *RubricCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RubricCategoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCategory(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rubric category, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a rubric category resource")
}

func (r *RubricCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
