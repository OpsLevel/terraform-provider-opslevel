package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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

var _ resource.ResourceWithConfigure = &IntegrationAwsResource{}

var _ resource.ResourceWithImportState = &IntegrationAwsResource{}

func NewIntegrationAwsResource() resource.Resource {
	return &IntegrationAwsResource{}
}

// IntegrationAwsResource defines the resource implementation.
type IntegrationAwsResource struct {
	CommonResourceClient
}

// IntegrationAwsResourceModel describes the AWS Integraion managed resource.
type IntegrationAwsResourceModel struct {
	ExternalID            types.String `tfsdk:"external_id"`
	IamRole               types.String `tfsdk:"iam_role"`
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	OwnershipTagOverrides types.Bool   `tfsdk:"ownership_tag_overrides"`
	OwnershipTagKeys      types.List   `tfsdk:"ownership_tag_keys"`
}

func NewIntegrationAwsResourceModel(awsIntegration opslevel.Integration) IntegrationAwsResourceModel {
	integrationAwsResourceModel := IntegrationAwsResourceModel{
		ExternalID:            RequiredStringValue(awsIntegration.ExternalID),
		IamRole:               RequiredStringValue(awsIntegration.IAMRole),
		Id:                    ComputedStringValue(string(awsIntegration.Id)),
		Name:                  OptionalStringValue(awsIntegration.Name),
		OwnershipTagOverrides: types.BoolValue(awsIntegration.OwnershipTagOverride),
	}
	ownershipTagKeys := OptionalStringListValue(awsIntegration.OwnershipTagKeys)
	integrationAwsResourceModel.OwnershipTagKeys = ownershipTagKeys

	return integrationAwsResourceModel
}

func (r *IntegrationAwsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_aws"
}

func (r *IntegrationAwsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "AWS Integration resource",

		Attributes: map[string]schema.Attribute{
			"external_id": schema.StringAttribute{
				Description: "The External ID defined in the trust relationship to ensure OpsLevel is the only third party assuming this role (See https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html for more details).",
				Required:    true,
			},
			"iam_role": schema.StringAttribute{
				Description: "The IAM role OpsLevel uses in order to access the AWS account.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the AWS integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ownership_tag_keys": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "An Array of tag keys used to associate ownership from an integration. Max 5",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(5),
				},
			},
			"ownership_tag_overrides": schema.BoolAttribute{
				Description: "Allow tags imported from AWS to override ownership set in OpsLevel directly.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
		},
	}
}

func (r *IntegrationAwsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel IntegrationAwsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipTagKeys, diags := ListValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	input := opslevel.AWSIntegrationInput{
		Name:                 planModel.Name.ValueStringPointer(),
		IAMRole:              planModel.IamRole.ValueStringPointer(),
		ExternalID:           planModel.ExternalID.ValueStringPointer(),
		OwnershipTagOverride: planModel.OwnershipTagOverrides.ValueBoolPointer(),
		OwnershipTagKeys:     ownershipTagKeys,
	}

	awsIntegration, err := r.client.CreateIntegrationAWS(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create aws integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationAwsResourceModel(*awsIntegration)

	tflog.Trace(ctx, "created an AWS integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationAwsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel IntegrationAwsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	awsIntegration, err := r.client.GetIntegration(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read AWS integration, got error: %s", err))
		return
	}

	verifiedStateModel := NewIntegrationAwsResourceModel(*awsIntegration)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read an AWS integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationAwsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel IntegrationAwsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipTagKeys, diags := ListValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	input := opslevel.AWSIntegrationInput{
		Name:                 opslevel.RefOf(planModel.Name.ValueString()),
		IAMRole:              opslevel.RefOf(planModel.IamRole.ValueString()),
		ExternalID:           opslevel.RefOf(planModel.ExternalID.ValueString()),
		OwnershipTagOverride: opslevel.RefOf(planModel.OwnershipTagOverrides.ValueBool()),
		OwnershipTagKeys:     ownershipTagKeys,
	}

	awsIntegration, err := r.client.UpdateIntegrationAWS(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update AWS integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationAwsResourceModel(*awsIntegration)

	tflog.Trace(ctx, "updated an AWS integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationAwsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationAwsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteIntegration(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete AWS integration, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted an AWS integration resource")
}

func (r *IntegrationAwsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
