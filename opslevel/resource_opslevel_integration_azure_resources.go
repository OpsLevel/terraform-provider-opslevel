package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

const AzureIdRegexPattern = `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`

var (
	_ resource.ResourceWithConfigure   = &IntegrationAzureResourcesResource{}
	_ resource.ResourceWithImportState = &IntegrationAzureResourcesResource{}
)

func NewIntegrationAzureResourcesResource() resource.Resource {
	return &IntegrationAzureResourcesResource{}
}

type IntegrationAzureResourcesResource struct {
	CommonResourceClient
}

type IntegrationAzureResourcesResourceModel struct {
	Aliases               types.List   `tfsdk:"aliases"`
	ClientId              types.String `tfsdk:"client_id"`
	ClientSecret          types.String `tfsdk:"client_secret"`
	CreatedAt             types.String `tfsdk:"created_at"`
	Id                    types.String `tfsdk:"id"`
	InstalledAt           types.String `tfsdk:"installed_at"`
	Name                  types.String `tfsdk:"name"`
	OwnershipTagKeys      types.List   `tfsdk:"ownership_tag_keys"`
	SubscriptionId        types.String `tfsdk:"subscription_id"`
	TagsOverrideOwnership types.Bool   `tfsdk:"ownership_tag_overrides"`
	TenantId              types.String `tfsdk:"tenant_id"`
}

func NewIntegrationAzureResourcesResourceModel(ctx context.Context, azureResourcesIntegration opslevel.Integration, givenModel IntegrationAzureResourcesResourceModel) IntegrationAzureResourcesResourceModel {
	resourceModel := IntegrationAzureResourcesResourceModel{
		Aliases:          OptionalStringListValue(azureResourcesIntegration.AzureResourcesIntegrationFragment.Aliases),
		ClientId:         givenModel.ClientId,
		ClientSecret:     givenModel.ClientSecret,
		CreatedAt:        ComputedStringValue(azureResourcesIntegration.CreatedAt.Local().Format(time.RFC850)),
		Id:               ComputedStringValue(string(azureResourcesIntegration.Id)),
		InstalledAt:      ComputedStringValue(azureResourcesIntegration.InstalledAt.Local().Format(time.RFC850)),
		Name:             RequiredStringValue(azureResourcesIntegration.Name),
		OwnershipTagKeys: OptionalStringListValue(azureResourcesIntegration.AzureResourcesIntegrationFragment.OwnershipTagKeys),
		SubscriptionId:   RequiredStringValue(azureResourcesIntegration.SubscriptionId),
		TenantId:         RequiredStringValue(azureResourcesIntegration.TenantId),
	}
	if givenModel.TagsOverrideOwnership.IsNull() {
		resourceModel.TagsOverrideOwnership = types.BoolNull()
	} else {
		resourceModel.TagsOverrideOwnership = types.BoolValue(azureResourcesIntegration.AzureResourcesIntegrationFragment.TagsOverrideOwnership)
	}

	return resourceModel
}

func (r *IntegrationAzureResourcesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_azure_resources"
}

func (r *IntegrationAzureResourcesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Azure Resources Integration resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				Description: "All of the aliases attached to the integration.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"client_id": schema.StringAttribute{
				Description: "The client id OpsLevel uses to access the Azure account.",
				Required:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret OpsLevel uses to access the Azure account.",
				Required:    true,
				Sensitive:   true,
			},
			"created_at": schema.StringAttribute{
				Description: "The time this integration was created.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"installed_at": schema.StringAttribute{
				Description: "The time that this integration was successfully installed, if null, this indicates the integration was not completely installed.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
			"ownership_tag_keys": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "An Array of tag keys used to associate ownership from an integration. Max 5 (default = [\"owner\"])",
				Optional:    true,
				Computed:    true,
				// current API default below
				Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("owner")})),
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.SizeBetween(1, 5),
				},
			},
			"ownership_tag_overrides": schema.BoolAttribute{
				Description: "Allow tags imported from Azure to override ownership set in OpsLevel directly.",
				Optional:    true,
				Computed:    true,
				// current API default is true
				Default: booldefault.StaticBool(true),
			},
			"subscription_id": schema.StringAttribute{
				MarkdownDescription: "The subscription OpsLevel uses to access the Azure account. [Microsoft's docs on regex pattern for ID](https://learn.microsoft.com/en-us/rest/api/defenderforcloud/tasks/get-subscription-level-task?view=rest-defenderforcloud-2015-06-01-preview&tabs=HTTP#uri-parameters)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(AzureIdRegexPattern),
						fmt.Sprintf("expected ID matching regex pattern: ' %s '", AzureIdRegexPattern),
					),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant OpsLevel uses to access the Azure account. [Microsoft's docs on regex pattern for ID](https://learn.microsoft.com/en-us/rest/api/defenderforcloud/tasks/get-subscription-level-task?view=rest-defenderforcloud-2015-06-01-preview&tabs=HTTP#uri-parameters)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(AzureIdRegexPattern),
						fmt.Sprintf("expected ID matching regex pattern: ' %s '", AzureIdRegexPattern),
					),
				},
			},
		},
	}
}

func (r *IntegrationAzureResourcesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[IntegrationAzureResourcesResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipTagKeys, diags := ListValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.AzureResourcesIntegrationInput{
		ClientId:              nullable(planModel.ClientId.ValueStringPointer()),
		ClientSecret:          nullable(planModel.ClientSecret.ValueStringPointer()),
		Name:                  nullable(planModel.Name.ValueStringPointer()),
		OwnershipTagKeys:      &opslevel.Nullable[[]string]{Value: ownershipTagKeys}, // TODO: why does this need to be nullable?
		SubscriptionId:        nullable(planModel.SubscriptionId.ValueStringPointer()),
		TagsOverrideOwnership: nullable(planModel.TagsOverrideOwnership.ValueBoolPointer()),
		TenantId:              nullable(planModel.TenantId.ValueStringPointer()),
	}

	azureResourcesIntegration, err := r.client.CreateIntegrationAzureResources(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Azure Resources integration, got error: '%s'", err))
		return
	}

	stateModel := NewIntegrationAzureResourcesResourceModel(ctx, *azureResourcesIntegration, planModel)

	tflog.Trace(ctx, "created an Azure Resources integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationAzureResourcesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[IntegrationAzureResourcesResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	azureResourcesIntegration, err := r.client.GetIntegration(asID(stateModel.Id))
	if err != nil {
		if (azureResourcesIntegration == nil || azureResourcesIntegration.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Azure Resources integration, got error: '%s'", err))
		return
	}

	verifiedStateModel := NewIntegrationAzureResourcesResourceModel(ctx, *azureResourcesIntegration, stateModel)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read an Azure Resources integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationAzureResourcesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[IntegrationAzureResourcesResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}
	ownershipTagKeys, diags := ListValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.AzureResourcesIntegrationInput{
		ClientId:              nullable(planModel.ClientId.ValueStringPointer()),
		ClientSecret:          nullable(planModel.ClientSecret.ValueStringPointer()),
		Name:                  nullable(planModel.Name.ValueStringPointer()),
		OwnershipTagKeys:      &opslevel.Nullable[[]string]{Value: ownershipTagKeys}, // TODO: why does this need to be nullable?
		SubscriptionId:        nullable(planModel.SubscriptionId.ValueStringPointer()),
		TagsOverrideOwnership: nullable(planModel.TagsOverrideOwnership.ValueBoolPointer()),
		TenantId:              nullable(planModel.TenantId.ValueStringPointer()),
	}

	azureResourcesIntegration, err := r.client.UpdateIntegrationAzureResources(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Azure Resources integration, got error: '%s'", err))
		return
	}

	stateModel := NewIntegrationAzureResourcesResourceModel(ctx, *azureResourcesIntegration, planModel)

	tflog.Trace(ctx, "updated an Azure Resources integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationAzureResourcesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[IntegrationAzureResourcesResourceModel](ctx, &resp.Diagnostics, req.State)
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
