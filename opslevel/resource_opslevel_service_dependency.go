package opslevel

import (
	"context"
	"fmt"
	"slices"
	"strings"

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

func NewServiceDependencyResourceModel(serviceDependency opslevel.ServiceDependency, givenModel ServiceDependencyResourceModel) (ServiceDependencyResourceModel, diag.Diagnostics) {
	var diag diag.Diagnostics
	serviceDependencyResourceModel := ServiceDependencyResourceModel{
		Id:   ComputedStringValue(string(serviceDependency.Id)),
		Note: OptionalStringValue(serviceDependency.Notes),
	}

	dependsUpon := identifierFromServiceId(givenModel.DependsUpon.ValueString(), serviceDependency.DependsOn)
	if dependsUpon == "" {
		diag.AddError("opslevel client error", fmt.Sprintf("expected depends_upon '%s' got '%s'", givenModel.DependsUpon.ValueString(), dependsUpon))
	}
	serviceDependencyResourceModel.DependsUpon = RequiredStringValue(dependsUpon)

	serviceIdentifier := identifierFromServiceId(givenModel.Service.ValueString(), serviceDependency.Service)
	if serviceIdentifier == "" {
		diag.AddError("opslevel client error", fmt.Sprintf("expected service '%s' got '%s'", givenModel.Service.ValueString(), serviceIdentifier))
	}
	serviceDependencyResourceModel.Service = RequiredStringValue(serviceIdentifier)

	return serviceDependencyResourceModel, diag
}

func identifierFromServiceId(identifier string, serviceId opslevel.ServiceId) string {
	if opslevel.IsID(identifier) {
		return string(serviceId.Id)
	} else if slices.Contains(serviceId.Aliases, identifier) {
		return identifier
	}
	return ""
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
	planModel := read[ServiceDependencyResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceDependencyCreateInput := opslevel.ServiceDependencyCreateInput{
		DependencyKey: opslevel.ServiceDependencyKey{
			DestinationIdentifier: opslevel.NewIdentifier(planModel.DependsUpon.ValueString()),
			SourceIdentifier:      opslevel.NewIdentifier(planModel.Service.ValueString()),
		},
	}
	if !planModel.Note.IsNull() && !planModel.Note.IsUnknown() {
		serviceDependencyCreateInput.Notes = planModel.Note.ValueStringPointer()
	}
	serviceDependency, err := r.client.CreateServiceDependency(serviceDependencyCreateInput)
	if err != nil || serviceDependency == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create serviceDependency, got error: %s", err))
		return
	}
	stateModel, diag := NewServiceDependencyResourceModel(*serviceDependency, planModel)
	if diag.HasError() {
		return
	}

	tflog.Trace(ctx, "created a service dependency resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceDependencyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	planModel := read[ServiceDependencyResourceModel](ctx, &resp.Diagnostics, req.State)
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
	if extractedServiceDependency == nil {
		resp.Diagnostics.AddError("opslevel client error", "Unable to extract service dependency")
		return
	}

	var dependsOn opslevel.ServiceId
	if opslevel.IsID(planModel.DependsUpon.ValueString()) {
		dependsOn.Id = opslevel.ID(planModel.DependsUpon.ValueString())
	} else {
		dependsOn.Aliases = []string{planModel.DependsUpon.ValueString()}
	}
	serviceDependency := opslevel.ServiceDependency{
		DependsOn: dependsOn,
		Id:        opslevel.ID(planModel.Id.ValueString()),
		Notes:     extractedServiceDependency.Notes,
		Service:   *extractedServiceDependency.Node,
	}

	readServiceDependencyResourceModel, diag := NewServiceDependencyResourceModel(serviceDependency, planModel)
	if diag.HasError() {
		return
	}

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
		"service dependencies should never be updated, only replaced.\nplease file a bug report including your .tf file at: github.com/OpsLevel/terraform-provider-opslevel")
}

func (r *ServiceDependencyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[ServiceDependencyResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteServiceDependency(opslevel.ID(stateModel.Id.ValueString())); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service dependency, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a serviceDependency resource")
}

func (r *ServiceDependencyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !isTagValid(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<service-id>:<dependency-id>'. Given '%s'", req.ID),
		)
		return
	}

	ids := strings.Split(req.ID, ":")
	serviceId := ids[0]
	dependencyId := ids[1]

	service, err := r.client.GetService(opslevel.ID(serviceId))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	dependencies, err := service.GetDependencies(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service dependencies, got error: %s", err))
		return
	}

	var foundSvcDependencyId string
	for _, serviceDependency := range dependencies.Edges {
		if serviceDependency.Node != nil && serviceDependency.Node.Id == opslevel.ID(dependencyId) {
			foundSvcDependencyId = string(serviceDependency.Id)
			break
		}
	}
	if foundSvcDependencyId == "" {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf(
			"Unable to get service dependency of service '%s' dependent upon '%s'",
			serviceId,
			dependencyId,
		))
	}

	idPath := path.Root("id")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idPath, foundSvcDependencyId)...)

	servicePath := path.Root("service")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, servicePath, serviceId)...)

	dependencyPath := path.Root("depends_upon")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, dependencyPath, dependencyId)...)
}
