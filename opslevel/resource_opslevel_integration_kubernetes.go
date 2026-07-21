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
	"github.com/opslevel/opslevel-go/v2026"
)

var (
	_ resource.ResourceWithConfigure   = &IntegrationKubernetesResource{}
	_ resource.ResourceWithImportState = &IntegrationKubernetesResource{}
)

func NewIntegrationKubernetesResource() resource.Resource {
	return &IntegrationKubernetesResource{}
}

// IntegrationKubernetesResource defines the resource implementation.
type IntegrationKubernetesResource struct {
	CommonResourceClient
}

// IntegrationKubernetesResourceModel describes the Kubernetes Integration managed resource.
type IntegrationKubernetesResourceModel struct {
	ExtractDefinition   types.String `tfsdk:"extract_definition"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	TransformDefinition types.String `tfsdk:"transform_definition"`
}

// NewIntegrationKubernetesResourceModel builds the state model from an API response. The API
// parses and re-serializes the YAML definitions, so the returned strings never match the
// configured ones byte-for-byte - keep the given model's values when they are known and only
// fall back to the API values when unset (import, server-side defaults).
func NewIntegrationKubernetesResourceModel(kubernetesIntegration opslevel.Integration, givenModel IntegrationKubernetesResourceModel) IntegrationKubernetesResourceModel {
	extractDefinition := givenModel.ExtractDefinition
	if extractDefinition.IsNull() || extractDefinition.IsUnknown() {
		extractDefinition = OptionalStringValue(kubernetesIntegration.KubernetesIntegrationFragment.ExtractDefinition)
	}
	transformDefinition := givenModel.TransformDefinition
	if transformDefinition.IsNull() || transformDefinition.IsUnknown() {
		transformDefinition = OptionalStringValue(kubernetesIntegration.KubernetesIntegrationFragment.TransformDefinition)
	}
	return IntegrationKubernetesResourceModel{
		ExtractDefinition:   extractDefinition,
		Id:                  ComputedStringValue(string(kubernetesIntegration.Id)),
		Name:                RequiredStringValue(kubernetesIntegration.Name),
		TransformDefinition: transformDefinition,
	}
}

func (r *IntegrationKubernetesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_kubernetes"
}

func (r *IntegrationKubernetesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Kubernetes Integration resource",

		Attributes: map[string]schema.Attribute{
			"extract_definition": schema.StringAttribute{
				Description: "The YAML definition for extracting data from inbound payloads. If not set, OpsLevel's default extract definition is used. Note: OpsLevel normalizes the stored YAML, so the value read from the API may be formatted differently than the configured value. Removing this attribute keeps the last applied definition.",
				Optional:    true,
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the Kubernetes integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
			"transform_definition": schema.StringAttribute{
				Description: "The YAML definition for transforming extracted data to OpsLevel resources. If not set, OpsLevel's default transform definition is used. Note: OpsLevel normalizes the stored YAML, so the value read from the API may be formatted differently than the configured value. Removing this attribute keeps the last applied definition.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// newKubernetesIntegrationInput only sets the definition fields that have a known value.
// Unknown values (unset in config) must be omitted so the API keeps its defaults - sending
// an explicit null would clear the definitions on the server.
func newKubernetesIntegrationInput(planModel IntegrationKubernetesResourceModel) opslevel.KubernetesIntegrationInput {
	input := opslevel.KubernetesIntegrationInput{
		Name: nullable(planModel.Name.ValueStringPointer()),
	}
	if !planModel.ExtractDefinition.IsNull() && !planModel.ExtractDefinition.IsUnknown() {
		input.ExtractDefinition = refOf(opslevel.YAML(planModel.ExtractDefinition.ValueString()))
	}
	if !planModel.TransformDefinition.IsNull() && !planModel.TransformDefinition.IsUnknown() {
		input.TransformDefinition = refOf(opslevel.YAML(planModel.TransformDefinition.ValueString()))
	}
	return input
}

func (r *IntegrationKubernetesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[IntegrationKubernetesResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	kubernetesIntegration, err := r.client.CreateIntegrationKubernetes(newKubernetesIntegrationInput(planModel))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Kubernetes integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationKubernetesResourceModel(*kubernetesIntegration, planModel)

	tflog.Trace(ctx, "created a Kubernetes integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationKubernetesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[IntegrationKubernetesResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	kubernetesIntegration, err := r.client.GetIntegration(asID(stateModel.Id))
	if err != nil {
		if (kubernetesIntegration == nil || kubernetesIntegration.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Kubernetes integration, got error: %s", err))
		return
	}

	verifiedStateModel := NewIntegrationKubernetesResourceModel(*kubernetesIntegration, stateModel)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a Kubernetes integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationKubernetesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[IntegrationKubernetesResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	kubernetesIntegration, err := r.client.UpdateIntegrationKubernetes(planModel.Id.ValueString(), newKubernetesIntegrationInput(planModel))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Kubernetes integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationKubernetesResourceModel(*kubernetesIntegration, planModel)

	tflog.Trace(ctx, "updated a Kubernetes integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationKubernetesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[IntegrationKubernetesResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteIntegration(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Kubernetes integration, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a Kubernetes integration resource")
}

func (r *IntegrationKubernetesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
