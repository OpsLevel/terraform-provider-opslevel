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

// RubricCategoryResourceModel describes the RubricCategory managed resource.
type RubricCategoryResourceModel struct {
	Id          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
}

func NewRubricCategoryResourceModel(rubricCategory opslevel.Category) RubricCategoryResourceModel {
	return RubricCategoryResourceModel{
		Id:   types.StringValue(string(rubricCategory.Id)),
		Name: types.StringValue(rubricCategory.Name),
	}
}

func (r *RubricCategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_category"
}

func (r *RubricCategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "RubricCategory Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the rubricCategory.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the category.",
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
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create rubricCategory, got error: %s", err))
		return
	}

	createdRubricCategoryResourceModel := NewRubricCategoryResourceModel(*rubricCategory)
	createdRubricCategoryResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a rubricCategory resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdRubricCategoryResourceModel)...)
}

func (r *RubricCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RubricCategoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rubricCategory, err := r.client.GetCategory(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read rubricCategory, got error: %s", err))
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
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update rubricCategory, got error: %s", err))
		return
	}
	updatedRubricCategoryResourceModel := NewRubricCategoryResourceModel(*updatedRubricCategory)
	updatedRubricCategoryResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a rubricCategory resource")
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rubricCategory, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a rubricCategory resource")
}

func (r *RubricCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceRubricCategory() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a rubric category",
// 		Create:      wrap(resourceRubricCategoryCreate),
// 		Read:        wrap(resourceRubricCategoryRead),
// 		Update:      wrap(resourceRubricCategoryUpdate),
// 		Delete:      wrap(resourceRubricCategoryDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The display name of the category.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 		},
// 	}
// }

// func resourceRubricCategoryCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CategoryCreateInput{
// 		Name: d.Get("name").(string),
// 	}
// 	resource, err := client.CreateCategory(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceRubricCategoryRead(d, client)
// }

// func resourceRubricCategoryRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetCategory(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceRubricCategoryUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CategoryUpdateInput{
// 		Id: opslevel.ID(d.Id()),
// 	}

// 	if d.HasChange("name") {
// 		input.Name = opslevel.RefOf(d.Get("name").(string))
// 	}

// 	_, err := client.UpdateCategory(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceRubricCategoryRead(d, client)
// }

// func resourceRubricCategoryDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteCategory(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
