package opslevel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"

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

type googleCloudProjectResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func googleCloudProjectAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"url":  types.StringType,
	}
}

type integrationGoogleCloudResourceModel struct {
	Aliases               types.List   `tfsdk:"aliases"`
	ClientEmail           types.String `tfsdk:"client_email"`
	CreatedAt             types.String `tfsdk:"created_at"`
	Id                    types.String `tfsdk:"id"`
	InstalledAt           types.String `tfsdk:"installed_at"`
	Name                  types.String `tfsdk:"name"`
	OwnershipTagKeys      types.List   `tfsdk:"ownership_tag_keys"`
	PrivateKey            types.String `tfsdk:"private_key"`
	Projects              types.List   `tfsdk:"projects"`
	TagsOverrideOwnership types.Bool   `tfsdk:"ownership_tag_overrides"`
}

func newIntegrationGoogleCloudResourceModel(ctx context.Context, googleCloudIntegration opslevel.Integration, givenModel integrationGoogleCloudResourceModel, diags *diag.Diagnostics) integrationGoogleCloudResourceModel {
	resourceModel := integrationGoogleCloudResourceModel{
		Aliases:               OptionalStringListValue(googleCloudIntegration.GoogleCloudIntegrationFragment.Aliases),
		ClientEmail:           givenModel.ClientEmail,
		CreatedAt:             ComputedStringValue(googleCloudIntegration.CreatedAt.UTC().Format(time.RFC3339)),
		Id:                    ComputedStringValue(string(googleCloudIntegration.Id)),
		InstalledAt:           ComputedStringValue(googleCloudIntegration.InstalledAt.UTC().Format(time.RFC3339)),
		Name:                  RequiredStringValue(googleCloudIntegration.Name),
		OwnershipTagKeys:      StringSliceToListValue(googleCloudIntegration.GoogleCloudIntegrationFragment.OwnershipTagKeys),
		PrivateKey:            givenModel.PrivateKey,
		TagsOverrideOwnership: types.BoolValue(googleCloudIntegration.GoogleCloudIntegrationFragment.TagsOverrideOwnership),
	}

	projects := make([]googleCloudProjectResourceModel, len(googleCloudIntegration.Projects))
	for i, project := range googleCloudIntegration.Projects {
		projects[i] = googleCloudProjectResourceModel{
			ID:   RequiredStringValue(project.ID),
			Name: RequiredStringValue(project.Name),
			URL:  RequiredStringValue(project.URL),
		}
	}
	projectsList, tmp := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: googleCloudProjectAttrs()}, projects)
	diags.Append(tmp...)
	if diags.HasError() {
		return integrationGoogleCloudResourceModel{}
	}
	resourceModel.Projects = projectsList

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
				Description: "An Array of tag keys used to associate ownership from an integration. Max 5",
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(5),
				},
			},
			"ownership_tag_overrides": schema.BoolAttribute{
				Description: "Allow tags imported from Google Cloud to override ownership set in OpsLevel directly.",
				Required:    true,
			},
			"projects": schema.ListAttribute{
				Description: "A list of the Google Cloud projects that were imported by the integration.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: googleCloudProjectAttrs(),
				},
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

	ownershipTagKeys, diags := ListValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.GoogleCloudIntegrationInput{
		ClientEmail:           planModel.ClientEmail.ValueStringPointer(),
		Name:                  planModel.Name.ValueStringPointer(),
		OwnershipTagKeys:      &ownershipTagKeys,
		PrivateKey:            planModel.PrivateKey.ValueStringPointer(),
		TagsOverrideOwnership: planModel.TagsOverrideOwnership.ValueBoolPointer(),
	}

	createdIntegration, err := r.client.CreateIntegrationGCP(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Google Cloud integration, got error: '%s'", err))
		return
	}

	stateModel := newIntegrationGoogleCloudResourceModel(ctx, *createdIntegration, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a Google Cloud integration")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *integrationGoogleCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel integrationGoogleCloudResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readIntegration, err := r.client.GetIntegration(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Google Cloud integration, got error: '%s'", err))
		return
	}

	verifiedStateModel := newIntegrationGoogleCloudResourceModel(ctx, *readIntegration, stateModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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

	ownershipTagKeys, diags := ListValueToStringSlice(ctx, planModel.OwnershipTagKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.GoogleCloudIntegrationInput{
		ClientEmail:           planModel.ClientEmail.ValueStringPointer(),
		Name:                  planModel.Name.ValueStringPointer(),
		OwnershipTagKeys:      &ownershipTagKeys,
		PrivateKey:            planModel.PrivateKey.ValueStringPointer(),
		TagsOverrideOwnership: planModel.TagsOverrideOwnership.ValueBoolPointer(),
	}

	updatedIntegration, err := r.client.UpdateIntegrationGCP(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Google Cloud integration, got error: '%s'", err))
		return
	}

	stateModel := newIntegrationGoogleCloudResourceModel(ctx, *updatedIntegration, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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
