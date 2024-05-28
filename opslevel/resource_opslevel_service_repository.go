package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"slices"

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

var _ resource.ResourceWithConfigure = &ServiceRepositoryResource{}

var _ resource.ResourceWithImportState = &ServiceRepositoryResource{}

func NewServiceRepositoryResource() resource.Resource {
	return &ServiceRepositoryResource{}
}

// ServiceRepositoryResource defines the resource implementation.
type ServiceRepositoryResource struct {
	CommonResourceClient
}

// ServiceRepositoryResourceModel describes the Servicerepository managed resource.
type ServiceRepositoryResourceModel struct {
	BaseDirectory   types.String `tfsdk:"base_directory"`
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Name            types.String `tfsdk:"name"`
	Repository      types.String `tfsdk:"repository"`
	RepositoryAlias types.String `tfsdk:"repository_alias"`
	Service         types.String `tfsdk:"service"`
	ServiceAlias    types.String `tfsdk:"service_alias"`
}

func NewServiceRepositoryResourceModel(ctx context.Context, serviceRepository opslevel.ServiceRepository, planModel ServiceRepositoryResourceModel) ServiceRepositoryResourceModel {
	stateModel := ServiceRepositoryResourceModel{
		BaseDirectory: OptionalStringValue(serviceRepository.BaseDirectory),
		Id:            ComputedStringValue(string(serviceRepository.Id)),
		LastUpdated:   planModel.LastUpdated,
		Name:          OptionalStringValue(serviceRepository.DisplayName),
	}
	if planModel.Repository.ValueString() == string(serviceRepository.Repository.Id) {
		stateModel.Repository = OptionalStringValue(string(serviceRepository.Repository.Id))
	}
	if planModel.RepositoryAlias.ValueString() == serviceRepository.Repository.DefaultAlias {
		stateModel.RepositoryAlias = OptionalStringValue(serviceRepository.Repository.DefaultAlias)
	}
	if planModel.Service.ValueString() == string(serviceRepository.Service.Id) {
		stateModel.Service = OptionalStringValue(string(serviceRepository.Service.Id))
	}
	if slices.Contains(serviceRepository.Service.Aliases, planModel.ServiceAlias.ValueString()) {
		stateModel.ServiceAlias = planModel.ServiceAlias
	}
	return stateModel
}

func (r *ServiceRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_repository"
}

func (r *ServiceRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ServiceRepository Resource",

		Attributes: map[string]schema.Attribute{
			"base_directory": schema.StringAttribute{
				Description: "The directory in the repository containing opslevel.yml.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^\/].*`),
						"path must not start with '/'",
					),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the Service Repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name displayed in the UI for the service repository.",
				Optional:    true,
			},
			"repository": schema.StringAttribute{
				Description: "The id of the repository that this will be added to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					IdStringValidator(),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("repository"),
						path.MatchRoot("repository_alias")),
				},
			},
			"repository_alias": schema.StringAttribute{
				Description: "The alias of the repository that this will be added to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("repository"),
						path.MatchRoot("repository_alias")),
				},
			},
			"service": schema.StringAttribute{
				Description: "The id of the service that this will be added to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					IdStringValidator(),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("service"),
						path.MatchRoot("service_alias")),
				},
			},
			"service_alias": schema.StringAttribute{
				Description: "The alias of the service that this will be added to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("service"),
						path.MatchRoot("service_alias")),
				},
			},
		},
	}
}

func (r *ServiceRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel ServiceRepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositoryIdentifier := opslevel.IdentifierInput{
		Alias: planModel.RepositoryAlias.ValueStringPointer(),
	}
	if planModel.Repository.ValueString() != "" {
		repositoryIdentifier.Id = opslevel.NewID(planModel.Repository.ValueString())
	}
	serviceIdentifier := opslevel.IdentifierInput{
		Alias: planModel.ServiceAlias.ValueStringPointer(),
	}
	if planModel.Service.ValueString() != "" {
		serviceIdentifier.Id = opslevel.NewID(planModel.Service.ValueString())
	}
	serviceRepository, err := r.client.CreateServiceRepository(opslevel.ServiceRepositoryCreateInput{
		BaseDirectory: opslevel.RefOf(planModel.BaseDirectory.ValueString()),
		DisplayName:   opslevel.RefOf(planModel.Name.ValueString()),
		Repository:    repositoryIdentifier,
		Service:       serviceIdentifier,
	})
	if err != nil || serviceRepository == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create serviceRepository, got error: %s", err))
		return
	}
	stateModel := NewServiceRepositoryResourceModel(ctx, *serviceRepository, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a service repository resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var currentStateModel ServiceRepositoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &currentStateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var service *opslevel.Service
	serviceId := currentStateModel.Service.ValueString()
	if opslevel.IsID(serviceId) {
		service, err = r.client.GetService(opslevel.ID(serviceId))
	} else {
		service, err = r.client.GetServiceWithAlias(currentStateModel.ServiceAlias.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}

	var serviceRepository *opslevel.ServiceRepository
	for _, edge := range service.Repositories.Edges {
		for _, repository := range edge.ServiceRepositories {
			if string(repository.Id) == currentStateModel.Id.ValueString() {
				serviceRepository = &repository
				break
			}
		}
		if serviceRepository != nil {
			break
		}
	}
	if serviceRepository == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("expected ServiceRepository not found for service '%s'", currentStateModel.Service.ValueString()))
		return
	}

	verifiedStateModel := NewServiceRepositoryResourceModel(ctx, *serviceRepository, currentStateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func extractServiceRepository(id string, serviceDependencies opslevel.ServiceDependenciesConnection) *opslevel.ServiceDependenciesEdge {
	for _, serviceRepository := range serviceDependencies.Edges {
		if id == string(serviceRepository.Id) {
			return &serviceRepository
		}
	}
	return nil
}

func (r *ServiceRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel ServiceRepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameBeforeUpdate types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &nameBeforeUpdate)...)
	if !nameBeforeUpdate.IsNull() && planModel.Name.IsNull() {
		resp.Diagnostics.AddError("Known error", "Unable to unset 'name' field for now. We have a planned fix for this.")
		return
	}

	serviceRepository, err := r.client.UpdateServiceRepository(opslevel.ServiceRepositoryUpdateInput{
		BaseDirectory: opslevel.RefOf(planModel.BaseDirectory.ValueString()),
		DisplayName:   opslevel.RefOf(planModel.Name.ValueString()),
		Id:            opslevel.ID(planModel.Id.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service repository, got error: %s", err))
		return
	}

	stateModel := NewServiceRepositoryResourceModel(ctx, *serviceRepository, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a service repository resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel ServiceRepositoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteServiceRepository(opslevel.ID(planModel.Id.ValueString())); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service repository, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a serviceRepository resource")
}

func (r *ServiceRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
