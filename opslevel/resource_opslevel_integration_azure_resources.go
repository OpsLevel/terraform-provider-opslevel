package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &IntegrationAzureResourcesResource{}

var _ resource.ResourceWithImportState = &IntegrationAzureResourcesResource{}

func NewIntegrationAzureResourcesResource() resource.Resource {
	return &IntegrationAzureResourcesResource{}
}

type IntegrationAzureResourcesResource struct {
	CommonResourceClient
}

type IntegrationAzureResourcesResourceModel struct {
	Aliases        types.List   `tfsdk:"aliases"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	Id             types.String `tfsdk:"id"`
	LastSyncedAt   types.String `tfsdk:"last_synced_at"`
	Name           types.String `tfsdk:"name"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
}

func NewIntegrationAzureResourcesResourceModel(ctx context.Context, azureResourcesIntegration opslevel.Integration, givenModel IntegrationAzureResourcesResourceModel) (IntegrationAzureResourcesResourceModel, diag.Diagnostics) {
	resourceModel := IntegrationAzureResourcesResourceModel{
		ClientId:       givenModel.ClientId,
		ClientSecret:   givenModel.ClientSecret,
		Id:             ComputedStringValue(string(azureResourcesIntegration.Id)),
		Name:           RequiredStringValue(azureResourcesIntegration.Name),
		SubscriptionId: RequiredStringValue(azureResourcesIntegration.SubscriptionId),
		TenantId:       RequiredStringValue(azureResourcesIntegration.TenantId),
	}

	if azureResourcesIntegration.LastSyncedAt == nil {
		resourceModel.LastSyncedAt = types.StringNull()
	} else {
		resourceModel.LastSyncedAt = types.StringValue(azureResourcesIntegration.LastSyncedAt.String())
	}
	var diags diag.Diagnostics
	resourceModel.Aliases, diags = types.ListValueFrom(ctx, types.StringType, azureResourcesIntegration.Aliases)

	return resourceModel, diags
}

func (r *IntegrationAzureResourcesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_azure_resources"
}

func (r *IntegrationAzureResourcesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Azure Resources Integration resource",

		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The client id OpsLevel uses to access the Azure account.",
				Required:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret OpsLevel uses to access the Azure account.",
				Required:    true,
				Sensitive:   true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "The tenant OpsLevel uses to access the Azure account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subscription_id": schema.StringAttribute{
				Description: "The subscription OpsLevel uses to access the Azure account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_synced_at": schema.StringAttribute{
				Description: "The time the Integration last imported data from Azure.",
				Computed:    true,
			},
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "All of the aliases attached to the resource.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
		},
	}
}

func (r *IntegrationAzureResourcesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel IntegrationAzureResourcesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.AzureResourcesIntegrationInput{
		Name:           planModel.Name.ValueStringPointer(),
		TenantId:       planModel.TenantId.ValueStringPointer(),
		SubscriptionId: planModel.SubscriptionId.ValueStringPointer(),
		ClientId:       planModel.ClientId.ValueStringPointer(),
		ClientSecret:   planModel.ClientSecret.ValueStringPointer(),
	}

	azureResourcesIntegration, err := r.client.CreateIntegrationAzureResources(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Azure Resources integration, got error: '%s'", err))
		return
	}

	stateModel, diags := NewIntegrationAzureResourcesResourceModel(ctx, *azureResourcesIntegration, planModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Trace(ctx, "created an Azure Resources integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationAzureResourcesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel IntegrationAzureResourcesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	azureResourcesIntegration, err := r.client.GetIntegration(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Azure Resources integration, got error: '%s'", err))
		return
	}

	verifiedStateModel, diags := NewIntegrationAzureResourcesResourceModel(ctx, *azureResourcesIntegration, stateModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read an Azure Resources integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationAzureResourcesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel IntegrationAzureResourcesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.AzureResourcesIntegrationInput{
		Name:           planModel.Name.ValueStringPointer(),
		TenantId:       planModel.TenantId.ValueStringPointer(),
		SubscriptionId: planModel.SubscriptionId.ValueStringPointer(),
		ClientId:       planModel.ClientId.ValueStringPointer(),
		ClientSecret:   planModel.ClientSecret.ValueStringPointer(),
	}

	azureResourcesIntegration, err := r.client.UpdateIntegrationAzureResources(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Azure Resources integration, got error: '%s'", err))
		return
	}

	stateModel, diags := NewIntegrationAzureResourcesResourceModel(ctx, *azureResourcesIntegration, planModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Trace(ctx, "updated an Azure Resources integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationAzureResourcesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationAzureResourcesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteIntegration(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Azure Resources integration, got error: '%s'", err))
		return
	}
	tflog.Trace(ctx, "deleted an Azure Resources integration")
}

func (r *IntegrationAzureResourcesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
