package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	Name             types.String `tfsdk:"name"`
	Role             types.String `tfsdk:"role"`
	SendInvite       types.Bool   `tfsdk:"send_invite"`
	SkipWelcomeEmail types.Bool   `tfsdk:"skip_welcome_email"`
}

func NewUserResourceModel(user opslevel.User, model UserResourceModel) UserResourceModel {
	return UserResourceModel{
		Email:            RequiredStringValue(user.Email),
		Id:               ComputedStringValue(string(user.Id)),
		Name:             RequiredStringValue(user.Name),
		Role:             OptionalStringValue(string(user.Role)),
		SendInvite:       model.SendInvite,
		SkipWelcomeEmail: model.SkipWelcomeEmail,
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
			"name": schema.StringAttribute{
				Description: "The name of the user.",
				Required:    true,
			},
			"role": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The access role of the user. One of `%s`",
					strings.Join(opslevel.AllUserRole, "`, `"),
				),
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllUserRole...),
				},
			},
			"send_invite": schema.BoolAttribute{
				MarkdownDescription: "Send an invite email even if notifications are disabled for the account. **(default: false)**",
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				Optional:            true,
			},
			"skip_welcome_email": schema.BoolAttribute{
				MarkdownDescription: "Don't send an email welcoming the user to OpsLevel. **(default: true)**",
				Default:             booldefault.StaticBool(true),
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[UserResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	userInput := opslevel.UserInput{
		Name:             nullable(planModel.Name.ValueStringPointer()),
		Role:             asEnum[opslevel.UserRole](planModel.Role.ValueStringPointer()),
		SkipWelcomeEmail: nullable(planModel.SkipWelcomeEmail.ValueBoolPointer()),
	}
	user, err := r.client.InviteUser(planModel.Email.ValueString(), userInput, planModel.SendInvite.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}
	stateModel := NewUserResourceModel(*user, planModel)

	tflog.Trace(ctx, "created a user resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[UserResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(stateModel.Id.ValueString())
	if err != nil {
		if (user == nil || user.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	updatedStateModel := NewUserResourceModel(*user, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedStateModel)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[UserResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[UserResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if planModel.SendInvite != stateModel.SendInvite {
		resp.Diagnostics.AddWarning("opslevel_user update no-op", "Modifying the send_invite attribute has no effect on an existing user.")
	}
	input := opslevel.UserInput{
		Name:             unsetStringHelper(planModel.Name, stateModel.Name),
		SkipWelcomeEmail: nullable(planModel.SkipWelcomeEmail.ValueBoolPointer()),
	}
	if !planModel.Role.IsNull() {
		input.Role = asEnum[opslevel.UserRole](planModel.Role.ValueStringPointer())
	} // We don't check state here because its ambigous what to do if the end user removes the role configuration

	res, err := r.client.UpdateUser(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update user, got error: %s", err))
		return
	}
	stateModel = NewUserResourceModel(*res, planModel)

	tflog.Trace(ctx, "updated a user resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[UserResourceModel](ctx, &resp.Diagnostics, req.State)
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
