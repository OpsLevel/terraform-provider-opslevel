package opslevel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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

var (
	_ resource.ResourceWithConfigure   = &integrationGoogleCloudResource{}
	_ resource.ResourceWithImportState = &integrationGoogleCloudResource{}
)

func NewIntegrationGoogleCloudResource() resource.Resource {
	return &integrationGoogleCloudResource{}
}

type integrationGoogleCloudResource struct {
	CommonResourceClient
}

type integrationGoogleCloudResourceModel struct {
	Aliases               types.List                      `tfsdk:"aliases"`
	ClientEmail           types.String                    `tfsdk:"client_email"`
	CreatedAt             types.String                    `tfsdk:"created_at"`
	Id                    types.String                    `tfsdk:"id"`
	InstalledAt           types.String                    `tfsdk:"installed_at"`
	Name                  types.String                    `tfsdk:"name"`
	OwnershipTagKeys      types.Set                       `tfsdk:"ownership_tag_keys"`
	PrivateKey            types.String                    `tfsdk:"private_key"`
	Projects              googleCloudProjectResourceModel `tfsdk:"projects"`
	TagsOverrideOwnership types.Bool                      `tfsdk:"ownership_tag_overrides"`
}

type googleCloudProjectResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func newIntegrationGoogleCloudResourceModel(GoogleCloudIntegration opslevel.Integration, givenModel integrationGoogleCloudResourceModel) integrationGoogleCloudResourceModel {
	resourceModel := integrationGoogleCloudResourceModel{
		Aliases:     OptionalStringListValue(GoogleCloudIntegration.GoogleCloudIntegrationFragment.Aliases),
		ClientEmail: givenModel.ClientEmail,
		CreatedAt:   ComputedStringValue(GoogleCloudIntegration.CreatedAt.Local().Format(time.RFC850)),
		Id:          ComputedStringValue(string(GoogleCloudIntegration.Id)),
		InstalledAt: ComputedStringValue(GoogleCloudIntegration.InstalledAt.Local().Format(time.RFC850)),
		Name:        RequiredStringValue(GoogleCloudIntegration.Name),
		PrivateKey:  givenModel.PrivateKey,
		Projects:    givenModel.Projects,
	}
	if givenModel.OwnershipTagKeys.IsNull() {
		resourceModel.OwnershipTagKeys = types.SetNull(types.StringType)
	} else {
		resourceModel.OwnershipTagKeys = StringSliceToSetValue(GoogleCloudIntegration.GoogleCloudIntegrationFragment.OwnershipTagKeys)
	}
	if givenModel.TagsOverrideOwnership.IsNull() {
		resourceModel.TagsOverrideOwnership = types.BoolNull()
	} else {
		resourceModel.TagsOverrideOwnership = types.BoolValue(GoogleCloudIntegration.GoogleCloudIntegrationFragment.TagsOverrideOwnership)
	}

	return resourceModel
}

func (r *integrationGoogleCloudResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_google_cloud"
}

func (r *integrationGoogleCloudResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Google Cloud Integration resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				Description: "All of the aliases attached to the integration.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"client_email": schema.StringAttribute{
				Description: "The service account email OpsLevel uses to access the Google Cloud account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Description: "The private key for the service account that OpsLevel uses to access the Google Cloud account.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"ownership_tag_keys": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "An Array of tag keys used to associate ownership from an integration. Max 5",
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 5),
				},
			},
			"ownership_tag_overrides": schema.BoolAttribute{
				Description: "Allow tags imported from Google Cloud to override ownership set in OpsLevel directly.",
				Optional:    true,
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of the Google Cloud projects that were imported by the integration.",
				Required:    true,
				Computed:    true,
			},
		},
	}
}

func (r *integrationGoogleCloudResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel integrationGoogleCloudResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipTagKeys, diags := SetValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.GoogleCloudIntegrationInput{
		ClientEmail:           planModel.ClientEmail.ValueStringPointer(),
		Name:                  planModel.Name.ValueStringPointer(),
		PrivateKey:            planModel.PrivateKey.ValueStringPointer(),
		TagsOverrideOwnership: planModel.TagsOverrideOwnership.ValueBoolPointer(),
	}
	if len(ownershipTagKeys) > 0 {
		input.OwnershipTagKeys = &ownershipTagKeys
	}

	createdIntegration, err := r.client.CreateIntegrationGCP(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Google Cloud integration, got error: '%s'", err))
		return
	}

	stateModel := newIntegrationGoogleCloudResourceModel(*createdIntegration, planModel)

	tflog.Trace(ctx, "created a Google Cloud integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *integrationGoogleCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel integrationGoogleCloudResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readIntegration, err := r.client.GetIntegration(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Google Cloud integration, got error: '%s'", err))
		return
	}

	verifiedStateModel := newIntegrationGoogleCloudResourceModel(*readIntegration, stateModel)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a Google Cloud integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *integrationGoogleCloudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel integrationGoogleCloudResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.GoogleCloudIntegrationInput{
		ClientEmail:           planModel.ClientEmail.ValueStringPointer(),
		Name:                  planModel.Name.ValueStringPointer(),
		PrivateKey:            planModel.PrivateKey.ValueStringPointer(),
		TagsOverrideOwnership: planModel.TagsOverrideOwnership.ValueBoolPointer(),
	}
	ownershipTagKeys, diags := SetValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// schema requires at least one ownership tag,
	if planModel.OwnershipTagKeys.IsNull() {
		input.OwnershipTagKeys = new([]string)
	} else if len(ownershipTagKeys) > 0 {
		input.OwnershipTagKeys = &ownershipTagKeys
	}

	updatedIntegration, err := r.client.UpdateIntegrationGCP(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Google Cloud integration, got error: '%s'", err))
		return
	}

	stateModel := newIntegrationGoogleCloudResourceModel(*updatedIntegration, planModel)

	tflog.Trace(ctx, "updated a Google Cloud integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *integrationGoogleCloudResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationGoogleCloudResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteIntegration(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Google Cloud integration, got error: '%s'", err))
		return
	}
	tflog.Trace(ctx, "deleted a Google Cloud integration")
}

func (r *integrationGoogleCloudResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
