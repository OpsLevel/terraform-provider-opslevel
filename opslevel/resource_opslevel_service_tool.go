package opslevel

import (
	"context"
	"fmt"

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

var _ resource.ResourceWithConfigure = &ServiceToolResource{}

var _ resource.ResourceWithImportState = &ServiceToolResource{}

func NewServiceToolResource() resource.Resource {
	return &ServiceToolResource{}
}

// ServiceToolResource defines the resource implementation.
type ServiceToolResource struct {
	CommonResourceClient
}

// ServiceToolResourceModel describes the ServiceTool managed resource.
type ServiceToolResourceModel struct {
	Category     types.String `tfsdk:"category"`
	Environment  types.String `tfsdk:"environment"`
	Id           types.String `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	Name         types.String `tfsdk:"name"`
	Service      types.String `tfsdk:"service"`
	ServiceAlias types.String `tfsdk:"service_alias"`
	Url          types.String `tfsdk:"url"`
}

func NewServiceToolResourceModel(ctx context.Context, serviceTool opslevel.Tool, planModel ServiceToolResourceModel) ServiceToolResourceModel {
	stateModel := ServiceToolResourceModel{
		Category:    RequiredStringValue(string(serviceTool.Category)),
		Environment: OptionalStringValue(serviceTool.Environment),
		Id:          ComputedStringValue(string(serviceTool.Id)),
		Name:        RequiredStringValue(serviceTool.DisplayName),
		Url:         RequiredStringValue(serviceTool.Url),
	}
	if planModel.Service.ValueString() == string(serviceTool.Service.Id) {
		stateModel.Service = OptionalStringValue(string(serviceTool.Service.Id))
	}
	stateModel.ServiceAlias = planModel.ServiceAlias
	return stateModel
}

func (r *ServiceToolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_tool"
}

func (r *ServiceToolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ServiceTool Resource",

		Attributes: map[string]schema.Attribute{
			"category": schema.StringAttribute{
				Description: "The category that the tool belongs to.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllToolCategory...)},
			},
			"environment": schema.StringAttribute{
				Description: "The environment that the tool belongs to.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the serviceTool.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the tool.",
				Required:    true,
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
			"url": schema.StringAttribute{
				Description: "The URL of the tool.",
				Required:    true,
			},
		},
	}
}

func (r *ServiceToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel ServiceToolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var service *opslevel.Service
	serviceId := planModel.Service.ValueString()
	if opslevel.IsID(serviceId) {
		service, err = r.client.GetService(opslevel.ID(serviceId))
	} else {
		service, err = r.client.GetServiceWithAlias(planModel.ServiceAlias.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service during create, got error: %s", err))
		return
	}

	serviceTool, err := r.client.CreateTool(opslevel.ToolCreateInput{
		Category:    opslevel.ToolCategory(planModel.Category.ValueString()),
		DisplayName: planModel.Name.ValueString(),
		Environment: planModel.Environment.ValueStringPointer(),
		ServiceId:   &service.Id,
		Url:         planModel.Url.ValueString(),
	})
	if err != nil || serviceTool == nil || string(serviceTool.Id) == "" {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service tool, got error: %s", err))
		return
	}
	stateModel := NewServiceToolResourceModel(ctx, *serviceTool, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a service tool resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var currentStateModel ServiceToolResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &currentStateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var service *opslevel.Service
	if serviceId := currentStateModel.Service.ValueString(); opslevel.IsID(serviceId) {
		service, err = r.client.GetService(opslevel.ID(serviceId))
	} else {
		service, err = r.client.GetServiceWithAlias(currentStateModel.ServiceAlias.ValueString())
	}
	if err != nil || service == nil || string(service.Id) == "" {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}

	var serviceTool *opslevel.Tool
	id := currentStateModel.Id.ValueString()
	for _, tool := range service.Tools.Nodes {
		if string(tool.Id) == id {
			serviceTool = &tool
			break
		}
	}
	if serviceTool == nil || string(serviceTool.Id) == "" {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to find tool with id '%s' on service with id '%s'", id, service.Id))
		return
	}

	verifiedStateModel := NewServiceToolResourceModel(ctx, *serviceTool, currentStateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *ServiceToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel ServiceToolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceTool, err := r.client.UpdateTool(opslevel.ToolUpdateInput{
		Category:    opslevel.RefOf(opslevel.ToolCategory(planModel.Category.ValueString())),
		DisplayName: planModel.Name.ValueStringPointer(),
		Environment: opslevel.RefOf(planModel.Environment.ValueString()),
		Id:          opslevel.ID(planModel.Id.ValueString()),
		Url:         planModel.Url.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service tool, got error: %s", err))
		return
	}

	stateModel := NewServiceToolResourceModel(ctx, *serviceTool, planModel)
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a service tool resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel ServiceToolResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteTool(opslevel.ID(planModel.Id.ValueString())); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service tool, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a serviceTool resource")
}

func (r *ServiceToolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
