package opslevel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func newSecretResourceModel(secret opslevel.Secret, sensitiveValue string) SecretResourceModel {
	secretResourceModel := SecretResourceModel{
		Alias:     RequiredStringValue(secret.Alias),
		CreatedAt: ComputedStringValue(secret.Timestamps.CreatedAt.Local().Format(time.RFC850)),
		Id:        ComputedStringValue(string(secret.ID)),
		Value:     RequiredStringValue(sensitiveValue),
	}
	return secretResourceModel
}

func updateSecretResourceModelWithPlan(secret opslevel.Secret, secretResourceModel *SecretResourceModel, planModel SecretResourceModel, client *opslevel.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	// properly set owner to team alias OR id - based on what is in the secret
	owner, err := getValidOwner(client, &secret, planModel.Owner.ValueString())
	if err != nil {
		diags.AddError("opslevel client error", err.Error())
		return diags
	}
	if owner == types.StringNull() {
		diags.AddError("terraform state error", "unexpected: got owner = null")
		return diags
	}
	secretResourceModel.Owner = owner

	return diags
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
	var planModel SecretResourceModel

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.CreateSecret(planModel.Alias.ValueString(), opslevel.SecretInput{
		Owner: opslevel.NewIdentifier(planModel.Owner.ValueString()),
		Value: opslevel.RefOf(planModel.Value.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create secret, got error: %s", err))
		return
	}
	newStateModel := newSecretResourceModel(*secret, planModel.Value.ValueString())
	diags := updateSecretResourceModelWithPlan(*secret, &newStateModel, planModel, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newStateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a secret resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel SecretResourceModel

	// Read Terraform prior state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.GetSecret(stateModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read secret, got error: %s", err))
		return
	}
	newStateModel := newSecretResourceModel(*secret, stateModel.Value.ValueString())
	// for owner - use alias or ID based on what was previously in the state
	newStateModel.Owner = stateModel.Owner

	// Save updated into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel SecretResourceModel

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedSecret, err := r.client.UpdateSecret(planModel.Id.ValueString(), opslevel.SecretInput{
		Owner: opslevel.NewIdentifier(planModel.Owner.ValueString()),
		Value: opslevel.RefOf(planModel.Value.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update secret, got error: %s", err))
		return
	}
	newStateModel := newSecretResourceModel(*updatedSecret, planModel.Value.ValueString())
	diags := updateSecretResourceModelWithPlan(*updatedSecret, &newStateModel, planModel, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newStateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a secret resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateModel SecretResourceModel

	// Read Terraform prior state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSecret(stateModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete secret, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a secret resource")
}

func (r *SecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
