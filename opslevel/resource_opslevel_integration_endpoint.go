package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

var _ resource.ResourceWithConfigure = &IntegrationEndpointResource{}

var _ resource.ResourceWithImportState = &IntegrationEndpointResource{}

func NewIntegrationEndpointResource() resource.Resource {
	return &IntegrationEndpointResource{}
}

// IntegrationEndpointResource defines the resource implementation.
type IntegrationEndpointResource struct {
	CommonResourceClient
}

// IntegrationEndpointResourceModel describes the CEC Integration managed resource.
type IntegrationEndpointResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func NewIntegrationEndpointResourceModel(cecIntegration opslevel.Integration, givenModel IntegrationEndpointResourceModel) IntegrationEndpointResourceModel {
	return IntegrationEndpointResourceModel{
		Id:   ComputedStringValue(string(cecIntegration.Id)),
		Name: RequiredStringValue(cecIntegration.Name),
		Type: RequiredStringValue(givenModel.Type.ValueString()),
	}
}

func (r *IntegrationEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_endpoint"
}

func (r *IntegrationEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Integration Endpoint resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Integration Endpoint.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The type of the Integration Endpoint. One of `%s`",
					strings.Join(opslevel.AllEventIntegrationEnum, "`, `"),
				),
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllEventIntegrationEnum...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *IntegrationEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel IntegrationEndpointResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.EventIntegrationInput{
		Name: planModel.Name.ValueStringPointer(),
		Type: opslevel.EventIntegrationEnum(planModel.Type.ValueString()),
	}

	cecIntegration, err := r.client.CreateEventIntegration(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create Integration Endpoint, got error: %s", err))
		return
	}

	stateModel := NewIntegrationEndpointResourceModel(*cecIntegration, planModel)

	tflog.Trace(ctx, "created a Integration Endpoint resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel IntegrationEndpointResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cecIntegration, err := r.client.GetIntegration(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read Integration Endpoint, got error: %s", err))
		return
	}

	verifiedStateModel := NewIntegrationEndpointResourceModel(*cecIntegration, stateModel)

	// Save updated data into Terraform state
	tflog.Trace(ctx, "read a Integration Endpoint resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *IntegrationEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel IntegrationEndpointResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.EventIntegrationUpdateInput{
		Id:   opslevel.ID(planModel.Id.ValueString()),
		Name: planModel.Name.ValueString(),
	}

	cecIntegration, err := r.client.UpdateEventIntegration(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update Integration Endpoint, got error: %s", err))
		return
	}

	stateModel := NewIntegrationEndpointResourceModel(*cecIntegration, planModel)

	tflog.Trace(ctx, "updated a Integration Endpoint resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *IntegrationEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationEndpointResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteIntegration(data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Custom Event  integration, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a Custom Event  integration resource")
}

func (r *IntegrationEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
