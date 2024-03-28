package opslevel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &SecretResource{}

var _ resource.ResourceWithImportState = &SecretResource{}

func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

// SecretResource defines the resource implementation.
type SecretResource struct {
	CommonResourceClient
}

// SecretResourceModel describes the Secret managed resource.
type SecretResourceModel struct {
	Alias       types.String `tfsdk:"alias"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Id          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Owner       types.String `tfsdk:"owner"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Value       types.String `tfsdk:"value"`
}

func NewSecretResourceModel(secret opslevel.Secret, ownerIdentifier, sensitiveValue string) SecretResourceModel {
	return SecretResourceModel{
		Alias:     types.StringValue(secret.Alias),
		CreatedAt: types.StringValue(secret.Timestamps.CreatedAt.Local().Format(time.RFC850)),
		Id:        types.StringValue(string(secret.ID)),
		Owner:     types.StringValue(ownerIdentifier),
		UpdatedAt: types.StringValue(secret.Timestamps.UpdatedAt.Local().Format(time.RFC850)),
		Value:     types.StringValue(sensitiveValue),
	}
}

func (r *SecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *SecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Secret Resource",

		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "The alias for this secret. Can only be set at create time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp of time created at.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the secret.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Description: "The owner of this secret.",
				Required:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp of last update.",
				Computed:    true,
			},
			"value": schema.StringAttribute{
				Description: "A sensitive value",
				Sensitive:   true,
				Required:    true,
			},
		},
	}
}

func (r *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecretResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.CreateSecret(data.Alias.ValueString(), opslevel.SecretInput{
		Owner: opslevel.NewIdentifier(data.Owner.ValueString()),
		Value: opslevel.RefOf(data.Value.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create secret, got error: %s", err))
		return
	}
	createdSecretResourceModel := NewSecretResourceModel(*secret, data.Owner.ValueString(), data.Value.ValueString())
	createdSecretResourceModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "created a secret resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdSecretResourceModel)...)
}

func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.GetSecret(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read secret, got error: %s", err))
		return
	}
	readSecretResourceModel := NewSecretResourceModel(*secret, data.Owner.ValueString(), data.Value.ValueString())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readSecretResourceModel)...)
}

func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SecretResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedSecret, err := r.client.UpdateSecret(data.Id.ValueString(), opslevel.SecretInput{
		Owner: opslevel.NewIdentifier(data.Owner.ValueString()),
		Value: opslevel.RefOf(data.Value.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update secret, got error: %s", err))
		return
	}
	updatedSecretResourceModel := NewSecretResourceModel(*updatedSecret, data.Owner.ValueString(), data.Value.ValueString())
	updatedSecretResourceModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "updated a secret resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedSecretResourceModel)...)
}

func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSecret(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete secret, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a secret resource")
}

func (r *SecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
