package opslevel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &ServiceDependencyResource{}

var _ resource.ResourceWithImportState = &ServiceDependencyResource{}

func NewServiceDependencyResource() resource.Resource {
	return &ServiceDependencyResource{}
}

// ServiceDependencyResource defines the resource implementation.
type ServiceDependencyResource struct {
	CommonResourceClient
}

// ServiceDependencyResourceModel describes the ServiceDependency managed resource.
type ServiceDependencyResourceModel struct {
	DependsUpon types.String `tfsdk:"depends_upon"`
	Id          types.String `tfsdk:"id"`
	Note        types.String `tfsdk:"note"`
	Service     types.String `tfsdk:"service"`
}

func NewServiceDependencyResourceModel(ctx context.Context, serviceDependency opslevel.ServiceDependency) ServiceDependencyResourceModel {
	return ServiceDependencyResourceModel{
		Id:   ComputedStringValue(string(serviceDependency.Id)),
		Note: OptionalStringValue(serviceDependency.Notes),
	}
}

func (r *ServiceDependencyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_dependency"
}

func (r *ServiceDependencyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ServiceDependency Resource",

		Attributes: map[string]schema.Attribute{
			"depends_upon": schema.StringAttribute{
				Description: "The ID or alias of the service that is depended upon.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the serviceDependency.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"note": schema.StringAttribute{
				Description: "Notes for service dependency.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service": schema.StringAttribute{
				Description: "The ID or alias of the service with the dependency.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ServiceDependencyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel ServiceDependencyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceDependency, err := r.client.CreateServiceDependency(opslevel.ServiceDependencyCreateInput{
		DependencyKey: opslevel.ServiceDependencyKey{
			DestinationIdentifier: opslevel.NewIdentifier(planModel.DependsUpon.ValueString()),
			SourceIdentifier:      opslevel.NewIdentifier(planModel.Service.ValueString()),
		},
		Notes: planModel.Note.ValueStringPointer(),
	})
	if err != nil || serviceDependency == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create serviceDependency, got error: %s", err))
		return
	}
	stateModel := NewServiceDependencyResourceModel(ctx, *serviceDependency)

	expectedDependsOn := planModel.DependsUpon.ValueString()
	if expectedDependsOn != string(serviceDependency.DependsOn.Id) &&
		!slices.Contains(serviceDependency.DependsOn.Aliases, expectedDependsOn) {
		resp.Diagnostics.AddError("Plan error", fmt.Sprintf("Created service dependency returned with unexpected 'depends_upon'. '%s'", expectedDependsOn))
		return
	}
	stateModel.DependsUpon = planModel.DependsUpon

	expectedService := planModel.Service.ValueString()
	if string(serviceDependency.Service.Id) != expectedService &&
		!slices.Contains(serviceDependency.Service.Aliases, expectedService) {
		resp.Diagnostics.AddError("Plan error", fmt.Sprintf("Created service dependency returned with unexpected service. '%s'", expectedService))
		return
	}
	stateModel.Service = planModel.Service

	tflog.Trace(ctx, "created a service dependency resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceDependencyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel ServiceDependencyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var service *opslevel.Service
	serviceIdentifier := planModel.Service.ValueString()
	if opslevel.IsID(serviceIdentifier) {
		service, err = r.client.GetService(opslevel.ID(serviceIdentifier))
	} else {
		service, err = r.client.GetServiceWithAlias(serviceIdentifier)
	}
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}

	dependencies, err := service.GetDependencies(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service dependencies, got error: %s", err))
		return
	}
	extractedServiceDependency := extractServiceDependency(planModel.Id.ValueString(), *dependencies)
	if extractedServiceDependency == nil || extractedServiceDependency.Id == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	serviceDependency := opslevel.ServiceDependency{
		Id:    opslevel.ID(planModel.Id.ValueString()),
		Notes: extractedServiceDependency.Notes,
	}
	readServiceDependencyResourceModel := NewServiceDependencyResourceModel(ctx, serviceDependency)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readServiceDependencyResourceModel)...)
}

func extractServiceDependency(id string, serviceDependencies opslevel.ServiceDependenciesConnection) *opslevel.ServiceDependenciesEdge {
	for _, serviceDependency := range serviceDependencies.Edges {
		if id == string(serviceDependency.Id) {
			return &serviceDependency
		}
	}
	return nil
}

func (r *ServiceDependencyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("terraform plugin error",
		"property assignments should never be updated, only replaced.\nplease file a bug report including your .tf file at: github.com/OpsLevel/terraform-provider-opslevel")
}

func (r *ServiceDependencyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel ServiceDependencyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteServiceDependency(opslevel.ID(planModel.Id.ValueString())); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service dependency, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a serviceDependency resource")
}

func (r *ServiceDependencyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
