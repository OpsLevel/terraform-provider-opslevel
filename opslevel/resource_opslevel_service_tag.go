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
	"github.com/opslevel/opslevel-go/v2025"
)

var _ resource.ResourceWithConfigure = &ServiceTagResource{}

var _ resource.ResourceWithImportState = &ServiceTagResource{}

type ServiceTagResource struct {
	CommonResourceClient
}

func NewServiceTagResource() resource.Resource {
	return &ServiceTagResource{}
}

type ServiceTagResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Key          types.String `tfsdk:"key"`
	Service      types.String `tfsdk:"service"`
	ServiceAlias types.String `tfsdk:"service_alias"`
	Value        types.String `tfsdk:"value"`
}

func NewServiceTagResourceModel(serviceTag opslevel.Tag) ServiceTagResourceModel {
	serviceResourceModel := ServiceTagResourceModel{
		Key:   RequiredStringValue(serviceTag.Key),
		Value: RequiredStringValue(serviceTag.Value),
		Id:    ComputedStringValue(string(serviceTag.Id)),
	}

	return serviceResourceModel
}

func (serviceTagResource *ServiceTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_tag"
}

func (serviceTagResource *ServiceTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service Tag Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The tag's key.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(TagKeyRegex, TagKeyErrorMsg),
				},
			},
			"value": schema.StringAttribute{
				Description: "The tag's value.",
				Required:    true,
			},
			"service": schema.StringAttribute{
				Description: "The id of the service that this will be added to.",
				Optional:    true,
				Validators: []validator.String{
					IdStringValidator(),
					stringvalidator.ExactlyOneOf(path.MatchRoot("service"),
						path.MatchRoot("service_alias")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_alias": schema.StringAttribute{
				Description: "The alias of the service that this will be added to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (serviceTagResource *ServiceTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := read[ServiceTagResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	tagCreateInput := opslevel.TagCreateInput{
		Type:  &opslevel.TaggableResourceService,
		Key:   data.Key.ValueString(),
		Value: data.Value.ValueString(),
	}

	// use either the service ID or alias based on what is used in the config
	var serviceIdentifier string
	if data.Service.ValueString() != "" {
		serviceIdentifier = data.Service.ValueString()
		tagCreateInput.Id = opslevel.NewID(serviceIdentifier)
	} else {
		serviceIdentifier = data.ServiceAlias.ValueString()
		tagCreateInput.Alias = &serviceIdentifier
	}

	service, err := serviceTagResource.client.CreateTag(tagCreateInput)
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to create service (%s) tag (with key '%s'), got error: %s", serviceIdentifier, data.Key.ValueString(), err))
		return
	}

	createdServiceTagResourceModel := NewServiceTagResourceModel(*service)
	// use either the service ID or alias based on what is used in the config
	if opslevel.IsID(serviceIdentifier) {
		createdServiceTagResourceModel.Service = OptionalStringValue(serviceIdentifier)
	} else {
		createdServiceTagResourceModel.ServiceAlias = OptionalStringValue(serviceIdentifier)
	}
	tflog.Trace(ctx, "created a service tag resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdServiceTagResourceModel)...)
}

func (serviceTagResource *ServiceTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := read[ServiceTagResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	// use either the service ID or alias based on what is used in the config
	var serviceIdentifier string
	var service *opslevel.Service
	var err error
	if data.Service.ValueString() != "" {
		serviceIdentifier = data.Service.ValueString()
		service, err = serviceTagResource.client.GetService(serviceIdentifier)
	} else {
		serviceIdentifier = data.ServiceAlias.ValueString()
		service, err = serviceTagResource.client.GetServiceWithAlias(serviceIdentifier)
	}
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read service (%s), got error: %s", serviceIdentifier, err))
		return
	}
	_, err = service.GetTags(serviceTagResource.client, nil)
	if err != nil || service.Tags == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read tags on service (%s), got error: %s", serviceIdentifier, err))
	}
	var serviceTag *opslevel.Tag
	for _, readTag := range service.Tags.Nodes {
		if readTag.Key == data.Key.ValueString() {
			serviceTag = &readTag
			break
		}
	}
	if serviceTag == nil || serviceTag.Id == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	readServiceResourceModel := NewServiceTagResourceModel(*serviceTag)
	// use either the service ID or alias based on what is used in the config
	if opslevel.IsID(serviceIdentifier) {
		readServiceResourceModel.Service = OptionalStringValue(serviceIdentifier)
	} else {
		readServiceResourceModel.ServiceAlias = OptionalStringValue(serviceIdentifier)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &readServiceResourceModel)...)
}

func (serviceTagResource *ServiceTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := read[ServiceTagResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// use either the service ID or alias based on what is used in the config
	var serviceIdentifier string
	if data.Service.ValueString() != "" {
		serviceIdentifier = data.Service.ValueString()
	} else {
		serviceIdentifier = data.ServiceAlias.ValueString()
	}

	tagUpdateInput := opslevel.TagUpdateInput{
		Id:    opslevel.ID(data.Id.ValueString()),
		Key:   data.Key.ValueStringPointer(),
		Value: data.Value.ValueStringPointer(),
	}

	serviceTag, err := serviceTagResource.client.UpdateTag(tagUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to update service tag (with id '%s'), got error: %s", data.Id.ValueString(), err))
		return
	}

	updatedServiceTagResourceModel := NewServiceTagResourceModel(*serviceTag)
	// use either the service ID or alias based on what is used in the config
	if opslevel.IsID(serviceIdentifier) {
		updatedServiceTagResourceModel.Service = OptionalStringValue(serviceIdentifier)
	} else {
		updatedServiceTagResourceModel.ServiceAlias = OptionalStringValue(serviceIdentifier)
	}
	tflog.Trace(ctx, "updated a service tag")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedServiceTagResourceModel)...)
}

func (serviceTagResource *ServiceTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[ServiceTagResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := serviceTagResource.client.DeleteTag(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to delete service tag (with id '%s'), got error: %s", data.Id.ValueString(), err))
		return
	}
	tflog.Trace(ctx, "deleted a service tag resource")
}

func (serviceTagResource *ServiceTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !isTagValid(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid format for given Import Id",
			fmt.Sprintf("Id expected to be formatted as '<service-id>:<tag-id>'. Given '%s'", req.ID),
		)
		return
	}

	ids := strings.Split(req.ID, ":")
	serviceId := ids[0]
	tagId := ids[1]

	service, err := serviceTagResource.client.GetTaggableResource(opslevel.TaggableResourceService, serviceId)
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to read service (%s), got error: %s", serviceId, err))
		return
	}
	tags, diags := getTagsFromResource(serviceTagResource.client, service)
	resp.Diagnostics.Append(diags...)
	if tags == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to get tags from service with id '%s'", serviceId))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	serviceTag := extractTagFromTags(opslevel.ID(tagId), tags.Nodes)
	if serviceTag == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to find tag with id '%s' in service with id '%s'", tagId, serviceId))
		return
	}

	idPath := path.Root("id")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idPath, string(serviceTag.Id))...)

	keyPath := path.Root("key")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, keyPath, serviceTag.Key)...)

	servicePath := path.Root("service")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, servicePath, string(service.ResourceId()))...)

	valuePath := path.Root("value")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, valuePath, serviceTag.Value)...)
}
