package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"slices"
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
						regexp.MustCompile(`^[^\/].*[^\/]$`),
						"path must not start or end with '/'",
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
	planModel := read[ServiceRepositoryResourceModel](ctx, &resp.Diagnostics, req.Plan)
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

	tflog.Trace(ctx, "created a service repository resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	currentStateModel := read[ServiceRepositoryResourceModel](ctx, &resp.Diagnostics, req.State)
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
		resp.Diagnostics.AddError("opslevel client error", "Unable to find service repository")
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
	planModel := read[ServiceRepositoryResourceModel](ctx, &resp.Diagnostics, req.Plan)
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

	tflog.Trace(ctx, "updated a service repository resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[ServiceRepositoryResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if stateModel.Id.IsNull() {
		tflog.Trace(ctx, "ServiceRepository resource to delete already does not exist")
		return
	}

	if err := r.client.DeleteServiceRepository(opslevel.ID(stateModel.Id.ValueString())); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service repository, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a ServiceRepository resource")
}

func (r *ServiceRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.SplitN(req.ID, ":", 2)
	if len(ids) != 2 {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<service-id-or-alias>:<repository-id-or-alias>'. Given '%s'", req.ID),
		)
		return
	}

	serviceIdentifier := ids[0]
	repoIdentifier := ids[1]

	var service *opslevel.Service
	var err error
	if opslevel.IsID(serviceIdentifier) {
		service, err = r.client.GetService(opslevel.ID(serviceIdentifier))
	} else {
		service, err = r.client.GetServiceWithAlias(serviceIdentifier)
	}
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	repositories, err := service.GetRepositories(r.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service dependencies, got error: %s", err))
		return
	}

	var foundRepoIdentifier string
	for _, serviceRepoEdge := range repositories.Edges {
		if serviceRepoEdge.Node.Id != opslevel.ID(repoIdentifier) && serviceRepoEdge.Node.DefaultAlias != repoIdentifier {
			continue
		}
		for _, serviceRepo := range serviceRepoEdge.ServiceRepositories {
			if serviceRepo.Repository.Id == opslevel.ID(repoIdentifier) || serviceRepo.Repository.DefaultAlias == repoIdentifier {
				foundRepoIdentifier = string(serviceRepo.Id)
				break
			}
		}
		if foundRepoIdentifier != "" {
			break
		}
	}
	if foundRepoIdentifier == "" {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf(
			"Unable to find service '%s' with repository '%s'",
			serviceIdentifier,
			repoIdentifier,
		))
	}

	idPath := path.Root("id")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idPath, foundRepoIdentifier)...)

	if opslevel.IsID(serviceIdentifier) {
		servicePath := path.Root("service")
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, servicePath, serviceIdentifier)...)
	} else {
		serviceAliasPath := path.Root("service_alias")
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, serviceAliasPath, serviceIdentifier)...)
	}

	if opslevel.IsID(repoIdentifier) {
		repoPath := path.Root("repository")
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, repoPath, repoIdentifier)...)
	} else {
		repoAliasPath := path.Root("repository_alias")
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, repoAliasPath, repoIdentifier)...)
	}
}
