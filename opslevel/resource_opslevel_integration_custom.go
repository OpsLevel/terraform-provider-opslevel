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
	"github.com/opslevel/opslevel-go/v2026"
)

var (
	_ resource.ResourceWithConfigure   = &IntegrationCustomResource{}
	_ resource.ResourceWithImportState = &IntegrationCustomResource{}
)

func NewIntegrationCustomResource() resource.Resource {
	return &IntegrationCustomResource{}
}

// IntegrationCustomResource defines the resource implementation.
type IntegrationCustomResource struct {
	CommonResourceClient
}

// IntegrationCustomResourceModel describes the Custom Integration managed resource.
type IntegrationCustomResourceModel struct {
	EtlDefinition types.Object `tfsdk:"etl_definition"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WebhookURL    types.String `tfsdk:"webhook_url"`
}

func NewIntegrationCustomResourceModel(ctx context.Context, customIntegration opslevel.Integration, givenModel IntegrationCustomResourceModel, diags *diag.Diagnostics) IntegrationCustomResourceModel {
	// Every custom integration has a webhook URL, independent of whether any
	// extractor is push based. The nil check only guards against a malformed
	// response; RequiredStringValue below keeps "" out of state as null, which
	// UseStateForUnknown would otherwise pin against a later non-null value.
	webhookURL := ""
	if customIntegration.WebhookURL != nil {
		webhookURL = *customIntegration.WebhookURL
	}

	return IntegrationCustomResourceModel{
		EtlDefinition: reconcileIntegrationEtlDefinition(ctx, customIntegration.CustomIntegrationFragment.ExtractDefinition, customIntegration.CustomIntegrationFragment.TransformDefinition, givenModel.EtlDefinition, diags),
		Id:            ComputedStringValue(string(customIntegration.Id)),
		Name:          RequiredStringValue(customIntegration.Name),
		WebhookURL:    RequiredStringValue(webhookURL),
	}
}

func (r *IntegrationCustomResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_custom"
}

func (r *IntegrationCustomResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Custom Integration resource",

		Attributes: map[string]schema.Attribute{
			"etl_definition": integrationEtlDefinitionSchemaAttribute("If omitted, the integration is created without any mapping and will not ingest data until definitions are set."),
			"id": schema.StringAttribute{
				Description: "The ID of the Custom integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "The endpoint to send data to via webhook. Always present for a custom integration; it does not indicate whether any extractor is push based.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Both definitions are omitted when unset; empty strings are rejected by the schema validators.
func newCustomIntegrationInput(ctx context.Context, planModel IntegrationCustomResourceModel, diags *diag.Diagnostics) opslevel.CustomIntegrationInput {
	input := opslevel.CustomIntegrationInput{
		Name: nullable(planModel.Name.ValueStringPointer()),
	}
	extractDefinition, transformDefinition := integrationEtlDefinitionInput(ctx, planModel.EtlDefinition, diags)
	if extractDefinition != nil || transformDefinition != nil {
		input.ExtractDefinition = extractDefinition
		input.TransformDefinition = transformDefinition
	}
	return input
}

func (r *IntegrationCustomResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[IntegrationCustomResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := newCustomIntegrationInput(ctx, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	customIntegration, err := r.client.CreateIntegrationCustom(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Custom integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationCustomResourceModel(ctx, *customIntegration, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a Custom integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationCustomResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[IntegrationCustomResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	customIntegration, err := r.client.GetIntegration(asID(stateModel.Id))
	if err != nil {
		if (customIntegration == nil || customIntegration.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Custom integration, got error: %s", err))
		return
	}
	if !verifyIntegrationType(&resp.Diagnostics, *customIntegration, "custom") {
		return
	}

	verifiedStateModel := NewIntegrationCustomResourceModel(ctx, *customIntegration, stateModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a Custom integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationCustomResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[IntegrationCustomResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := newCustomIntegrationInput(ctx, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	customIntegration, err := r.client.UpdateIntegrationCustom(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Custom integration, got error: %s", err))
		return
	}

	stateModel := NewIntegrationCustomResourceModel(ctx, *customIntegration, planModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updated a Custom integration resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationCustomResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[IntegrationCustomResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteIntegration(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Custom integration, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a Custom integration resource")
}

func (r *IntegrationCustomResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
