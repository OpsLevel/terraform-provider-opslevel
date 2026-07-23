package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

type IntegrationKubernetesEtlDefinitionModel struct {
	ExtractDefinition   types.String `tfsdk:"extract_definition"`
	TransformDefinition types.String `tfsdk:"transform_definition"`
}

func integrationKubernetesEtlDefinitionAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"extract_definition":   types.StringType,
		"transform_definition": types.StringType,
	}
}

// IntegrationKubernetesResourceModel describes the Kubernetes Integration managed resource.
type IntegrationKubernetesResourceModel struct {
	EtlDefinition types.Object `tfsdk:"etl_definition"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
}

func NewIntegrationKubernetesResourceModel(ctx context.Context, kubernetesIntegration opslevel.Integration, givenModel IntegrationKubernetesResourceModel, diags *diag.Diagnostics) IntegrationKubernetesResourceModel {
	etlModel := IntegrationKubernetesEtlDefinitionModel{
		ExtractDefinition:   RequiredStringValue(kubernetesIntegration.KubernetesIntegrationFragment.ExtractDefinition),
		TransformDefinition: RequiredStringValue(kubernetesIntegration.KubernetesIntegrationFragment.TransformDefinition),
	}
	if !givenModel.EtlDefinition.IsNull() && !givenModel.EtlDefinition.IsUnknown() {
		var givenEtlModel IntegrationKubernetesEtlDefinitionModel
		diags.Append(givenModel.EtlDefinition.As(ctx, &givenEtlModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
		if givenEtlModel.ExtractDefinition.ValueString() == "" && givenEtlModel.TransformDefinition.ValueString() == "" {
			etlModel = givenEtlModel
		} else {
			if yamlEquivalent(givenEtlModel.ExtractDefinition.ValueString(), kubernetesIntegration.KubernetesIntegrationFragment.ExtractDefinition) {
				etlModel.ExtractDefinition = givenEtlModel.ExtractDefinition
			}
			if yamlEquivalent(givenEtlModel.TransformDefinition.ValueString(), kubernetesIntegration.KubernetesIntegrationFragment.TransformDefinition) {
				etlModel.TransformDefinition = givenEtlModel.TransformDefinition
			}
		}
	}
	etlObject, objDiags := types.ObjectValueFrom(ctx, integrationKubernetesEtlDefinitionAttrs(), etlModel)
	diags.Append(objDiags...)

	return IntegrationKubernetesResourceModel{
		EtlDefinition: etlObject,
		Id:            ComputedStringValue(string(kubernetesIntegration.Id)),
		Name:          RequiredStringValue(kubernetesIntegration.Name),
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
			"etl_definition": schema.SingleNestedAttribute{
				Description: "The ETL definitions used to import data from the integration. If not set (or removed), OpsLevel's default definitions are used. The API manages the two definitions as a unit, so both must be set together.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"extract_definition": schema.StringAttribute{
						Description: "The YAML definition for extracting data from inbound payloads.",
						Required:    true,
					},
					"transform_definition": schema.StringAttribute{
						Description: "The YAML definition for transforming extracted data to OpsLevel resources.",
						Required:    true,
					},
				},
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					integrationKubernetesEtlDefinitionAttrs(),
					map[string]attr.Value{
						"extract_definition":   types.StringValue(""),
						"transform_definition": types.StringValue(""),
					},
				)),
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
		},
	}
}

// Both definitions are omitted when unset - the API rejects an empty string for YAML.
func newKubernetesIntegrationInput(ctx context.Context, planModel IntegrationKubernetesResourceModel, diags *diag.Diagnostics) opslevel.KubernetesIntegrationInput {
	input := opslevel.KubernetesIntegrationInput{
		Name: nullable(planModel.Name.ValueStringPointer()),
	}
	if !planModel.EtlDefinition.IsNull() && !planModel.EtlDefinition.IsUnknown() {
		var etlModel IntegrationKubernetesEtlDefinitionModel
		diags.Append(planModel.EtlDefinition.As(ctx, &etlModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
		if etlModel.ExtractDefinition.ValueString() != "" || etlModel.TransformDefinition.ValueString() != "" {
			input.ExtractDefinition = refOf(opslevel.YAML(etlModel.ExtractDefinition.ValueString()))
			input.TransformDefinition = refOf(opslevel.YAML(etlModel.TransformDefinition.ValueString()))
		}
	}
	return input
}

func (r *IntegrationKubernetesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[IntegrationKubernetesResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := newKubernetesIntegrationInput(ctx, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	kubernetesIntegration, err := r.client.CreateIntegrationKubernetes(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Kubernetes integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationKubernetesResourceModel(ctx, *kubernetesIntegration, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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

	verifiedStateModel := NewIntegrationKubernetesResourceModel(ctx, *kubernetesIntegration, stateModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a Kubernetes integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationKubernetesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[IntegrationKubernetesResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := newKubernetesIntegrationInput(ctx, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	kubernetesIntegration, err := r.client.UpdateIntegrationKubernetes(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Kubernetes integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationKubernetesResourceModel(ctx, *kubernetesIntegration, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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
