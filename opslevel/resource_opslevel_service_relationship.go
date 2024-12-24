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

var _ resource.ResourceWithConfigure = &ServiceRelationshipResource{}

var _ resource.ResourceWithImportState = &ServiceRelationshipResource{}

func NewServiceRelationshipResource() resource.Resource {
	return &ServiceRelationshipResource{}
}

// ServiceRelationshipResource defines the resource implementation.
type ServiceRelationshipResource struct {
	CommonResourceClient
}

// ServiceRelationshipResourceModel describes the Service managed resource.
type ServiceRelationshipResourceModel struct {
	Service types.String `tfsdk:"service"`
	System  types.String `tfsdk:"system"`
}

func NewServiceRelationshipResourceModel(service *opslevel.Service, givenModel ServiceRelationshipResourceModel) ServiceRelationshipResourceModel {
	return ServiceRelationshipResourceModel{
		Service: givenModel.Service,
		System:  givenModel.System,
	}
}

func (r *ServiceRelationshipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_relationship"
}

func (r *ServiceRelationshipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Service Relationship Resource",

		Attributes: map[string]schema.Attribute{
			"service": schema.StringAttribute{
				Description: "The ID or alias of the service with the system.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system": schema.StringAttribute{
				Description: "The ID or alias of the system tied to the service.",
				Required:    true,
			},
		},
	}
}

func (r *ServiceRelationshipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var diag diag.Diagnostics

	planModel := read[ServiceRelationshipResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}
	serviceIdentifier := planModel.Service.ValueString()
	systemIdentifier := planModel.System.ValueString()
	service, err := getService(r.client, serviceIdentifier)
	if err != nil {
		diag.AddError("opslevel client error", fmt.Sprintf("Unable to read service '%s', got error: %s", serviceIdentifier, err))
		return
	}

	svcUpdate := opslevel.ServiceUpdateInputV2{
		Id:     &service.Id,
		Parent: opslevel.NewIdentifier(systemIdentifier),
	}
	if _, err := r.client.UpdateService(svcUpdate); err != nil {
		diag.AddError(
			"opslevel client error",
			fmt.Sprintf(
				"Unable to set parent system '%s' for service '%s', got error: %s",
				systemIdentifier,
				serviceIdentifier,
				err,
			))
		return
	}

	stateModel := NewServiceRelationshipResourceModel(service, planModel)

	tflog.Trace(ctx, "created a service relationship resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRelationshipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[ServiceRelationshipResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}
	serviceIdentifier := stateModel.Service.ValueString()
	systemIdentifier := stateModel.System.ValueString()

	service, err := getService(r.client, serviceIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	if service.Parent == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if (opslevel.IsID(systemIdentifier) && string(service.Parent.Id) != systemIdentifier) &&
		(!opslevel.IsID(systemIdentifier) && !slices.Contains(service.Parent.Aliases, systemIdentifier)) {
		resp.Diagnostics.AddError(
			"opslevel client error",
			fmt.Sprintf("Expected service '%s' to have parent system '%s' but it does not.", serviceIdentifier, systemIdentifier),
		)
		return
	}

	stateModel = NewServiceRelationshipResourceModel(service, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRelationshipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var diag diag.Diagnostics

	planModel := read[ServiceRelationshipResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}
	serviceIdentifier := planModel.Service.ValueString()
	systemIdentifier := planModel.System.ValueString()
	service, err := getService(r.client, serviceIdentifier)
	if err != nil {
		diag.AddError("opslevel client error", fmt.Sprintf("Unable to read service '%s', got error: %s", serviceIdentifier, err))
		return
	}

	svcUpdate := opslevel.ServiceUpdateInputV2{
		Id:     &service.Id,
		Parent: opslevel.NewIdentifier(systemIdentifier),
	}
	if _, err := r.client.UpdateService(svcUpdate); err != nil {
		diag.AddError(
			"opslevel client error",
			fmt.Sprintf(
				"Unable to set parent system '%s' for service '%s', got error: %s",
				systemIdentifier,
				serviceIdentifier,
				err,
			))
		return
	}

	stateModel := NewServiceRelationshipResourceModel(service, planModel)

	tflog.Trace(ctx, "updated a service relationship resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceRelationshipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[ServiceRelationshipResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}
	serviceIdentifier := stateModel.Service.ValueString()
	systemIdentifier := stateModel.System.ValueString()

	service, err := getService(r.client, serviceIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	if string(service.Parent.Id) == systemIdentifier || slices.Contains(service.Parent.Aliases, systemIdentifier) {
		svcUpdate := opslevel.ServiceUpdateInputV2{
			Id:     &service.Id,
			Parent: opslevel.NewIdentifier(),
		}
		if _, err = r.client.UpdateService(svcUpdate); err != nil {
			resp.Diagnostics.AddWarning("opslevel client error",
				fmt.Sprintf(
					"Issue removing parent system '%s' from service '%s', got error: %s",
					systemIdentifier,
					serviceIdentifier,
					err,
				))
		}
	}

	tflog.Trace(ctx, "deleted a service relationship resource")
}

func (r *ServiceRelationshipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !hasTagFormat(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<service-identifier>:<system-identifier>'. Given '%s'", req.ID),
		)
		return
	}

	ids := strings.Split(req.ID, ":")
	serviceIdentifier := ids[0]
	systemIdentifier := ids[1]

	if _, err := getService(r.client, serviceIdentifier); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	if _, err := r.client.GetSystem(systemIdentifier); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read system, got error: %s", err))
		return
	}

	servicePath := path.Root("service")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, servicePath, serviceIdentifier)...)

	dependencyPath := path.Root("system")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, dependencyPath, systemIdentifier)...)
}
