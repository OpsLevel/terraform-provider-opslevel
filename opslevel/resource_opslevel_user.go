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
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &UserResource{}

var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	CommonResourceClient
}

// UserResourceModel describes the User managed resource.
type UserResourceModel struct {
	Email            types.String `tfsdk:"email"`
	Id               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	Name             types.String `tfsdk:"name"`
	Role             types.String `tfsdk:"role"`
	SkipWelcomeEmail types.Bool   `tfsdk:"skip_welcome_email"` // not usable but kept for backwards compatibility
}

func NewUserResourceModel(user opslevel.User) UserResourceModel {
	return UserResourceModel{
		Email:            types.StringValue(user.Email),
		Id:               types.StringValue(string(user.Id)),
		LastUpdated:      timeLastUpdated(),
		Name:             types.StringValue(user.Name),
		Role:             types.StringValue(string(user.Role)),
		SkipWelcomeEmail: types.BoolNull(),
	}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User Resource",

		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "The email address of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the user.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user.",
				Required:    true,
			},
			"role": schema.StringAttribute{
				Description: "The access role (e.g. user or admin) of the user.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllUserRole...),
				},
			},
			"skip_welcome_email": schema.BoolAttribute{
				DeprecationMessage: "The skip_welcome_email attribute is deprecated and only kept for backward compatibility.",
				Description:        "Don't send an email welcoming the user to OpsLevel. Applies during creation only, this value cannot be read or updated.",
				Optional:           true,
			},
		},
	}
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.InviteUser(data.Email.ValueString(), opslevel.UserInput{
		Name: opslevel.RefOf(data.Name.ValueString()),
		Role: opslevel.RefOf(opslevel.UserRole(data.Role.ValueString())),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}
	createdUserResourceModel := NewUserResourceModel(*user)

	tflog.Trace(ctx, "created a user resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdUserResourceModel)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	data.Email = types.StringValue(user.Email)
	data.Id = types.StringValue(string(user.Id))
	data.Name = types.StringValue(user.Name)
	data.Role = types.StringValue(string(user.Role))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.UpdateUser(data.Id.ValueString(), opslevel.UserInput{
		Name: opslevel.RefOf(data.Name.ValueString()),
		Role: opslevel.RefOf(opslevel.UserRole(data.Role.ValueString())),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update user, got error: %s", err))
		return
	}
	updatedUserResourceModel := NewUserResourceModel(*resource)

	tflog.Trace(ctx, "updated a user resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedUserResourceModel)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a user resource")
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
