package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure ServiceDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &ServiceDataSource{}

func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

// ServiceDataSource manages a Service data source.
type ServiceDataSource struct {
	CommonDataSourceClient
}

// ServiceDataSourceModel describes the data source data model.
type ServiceDataSourceModel struct {
	Aliases                    types.List   `tfsdk:"aliases"`
	ApiDocumentPath            types.String `tfsdk:"api_document_path"`
	Description                types.String `tfsdk:"description"`
	Framework                  types.String `tfsdk:"framework"`
	Identifier                 types.String `tfsdk:"identifier"`
	Language                   types.String `tfsdk:"language"`
	LifecycleAlias             types.String `tfsdk:"lifecycle_alias"`
	Name                       types.String `tfsdk:"name"`
	Owner                      types.String `tfsdk:"owner"`
	OwnerId                    types.String `tfsdk:"owner_id"`
	PreferredApiDocumentSource types.String `tfsdk:"preferred_api_document_source"`
	Product                    types.String `tfsdk:"product"`
	Properties                 types.List   `tfsdk:"properties"`
	Repositories               types.List   `tfsdk:"repositories"`
	Tags                       types.List   `tfsdk:"tags"`
	TierAlias                  types.String `tfsdk:"tier_alias"`
}

func NewServiceDataSourceModel(ctx context.Context, service opslevel.Service, identifier string) (ServiceDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	serviceDataSourceModel := ServiceDataSourceModel{
		ApiDocumentPath: types.StringValue(service.ApiDocumentPath),
		Description:     types.StringValue(service.Description),
		Framework:       types.StringValue(service.Framework),
		Identifier:      types.StringValue(identifier),
		Language:        types.StringValue(service.Language),
		LifecycleAlias:  types.StringValue(service.Lifecycle.Alias),
		Name:            types.StringValue(service.Name),
		Owner:           types.StringValue(service.Owner.Alias),
		OwnerId:         types.StringValue(string(service.Owner.Id)),
		Product:         types.StringValue(service.Product),
		TierAlias:       types.StringValue(service.Tier.Alias),
	}

	serviceAliases, svcDiags := types.ListValueFrom(ctx, types.StringType, service.Aliases)
	diags = append(diags, svcDiags...)
	serviceDataSourceModel.Aliases = serviceAliases

	if service.PreferredApiDocumentSource != nil {
		serviceDataSourceModel.PreferredApiDocumentSource = types.StringValue(string(*service.PreferredApiDocumentSource))
	}

	if service.Tags == nil {
		serviceDataSourceModel.Tags = types.ListNull(types.StringType)
	} else {
		serviceTags, tagsDiags := types.ListValueFrom(ctx, types.StringType, flattenTagArray(service.Tags.Nodes))
		serviceDataSourceModel.Tags = serviceTags
		diags = append(diags, tagsDiags...)
	}

	return serviceDataSourceModel, diags
}

func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service data source",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases of the service.",
				Computed:    true,
			},
			"api_document_path": schema.StringAttribute{
				Description: "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A brief description of the service.",
				Computed:    true,
			},
			"framework": schema.StringAttribute{
				Description: "The primary software development framework that the service uses.",
				Computed:    true,
			},
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the service to find.",
				Required:    true,
			},
			"language": schema.StringAttribute{
				Description: "The primary programming language that the service is written in.",
				Computed:    true,
			},
			"lifecycle_alias": schema.StringAttribute{
				Description: "The lifecycle stage of the service.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the service.",
				Computed:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The team that owns the service.",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The team ID that owns the service.",
				Computed:    true,
			},
			"preferred_api_document_source": schema.StringAttribute{
				Description: "The API document source (PUSH or PULL) used to determine the displayed document. If null, we use the order push and then pull.",
				Computed:    true,
			},
			"product": schema.StringAttribute{
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				Computed:    true,
			},
			"properties": schema.ListAttribute{
				Description: "Custom properties assigned to this service.",
				ElementType: opslevelPropertyObjectType,
				Computed:    true,
			},
			"repositories": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of repositories connected to the service.",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of tags applied to the service.",
				Computed:    true,
			},
			"tier_alias": schema.StringAttribute{
				Description: "The software tier that the service belongs to.",
				Computed:    true,
			},
		},
	}
}

func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, err := getService(*d.client, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service datasource, got error: %s", err))
		return
	}
	serviceDataModel, diags := NewServiceDataSourceModel(ctx, service, data.Identifier.ValueString())
	resp.Diagnostics.Append(diags...)

	// NOTE: service's hydrate does not populate properties
	serviceDataModel.Properties, diags = getServiceProperties(ctx, d.client, service)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if service.Repositories == nil {
		serviceDataModel.Repositories = types.ListNull(types.StringType)
	} else {
		serviceDataModel.Repositories, diags = types.ListValueFrom(
			ctx,
			types.StringType,
			flattenServiceRepositoriesArray(service.Repositories),
		)
	}

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Service data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &serviceDataModel)...)
}

func getServiceProperties(ctx context.Context, client *opslevel.Client, service opslevel.Service) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	properties, err := service.GetProperties(client, nil)
	if err != nil {
		diags.AddAttributeError(
			path.Root("properties"),
			"OpsLevel Client Error",
			fmt.Sprintf("unable to read Properties for service, got error: %s", err),
		)
		return types.ListNull(opslevelPropertyObjectType), diags
	}
	if properties == nil {
		return types.ListNull(opslevelPropertyObjectType), diags
	}

	serviceProperties, diags := opslevelPropertiesToListValue(ctx, properties.Nodes)
	if diags != nil && diags.HasError() {
		return types.ListNull(opslevelPropertyObjectType), diags
	}
	return serviceProperties, diags
}

func getService(client opslevel.Client, data ServiceDataSourceModel) (opslevel.Service, error) {
	var err error
	var service *opslevel.Service

	identifier := data.Identifier.ValueString()
	if opslevel.IsID(identifier) {
		service, err = client.GetService(opslevel.ID(identifier))
	} else {
		service, err = client.GetServiceWithAlias(identifier)
	}
	if err != nil {
		return opslevel.Service{}, err
	}

	if service == nil || service.Id == "" {
		return opslevel.Service{}, fmt.Errorf("unable to find repository with identifier=`%s`", identifier)
	}

	return *service, nil
}

func opslevelPropertiesToListValue(ctx context.Context, opslevelProperties []opslevel.Property) (basetypes.ListValue, diag.Diagnostics) {
	properties := make([]attr.Value, len(opslevelProperties))
	for i, property := range opslevelProperties {
		propertyObject, diags := opslevelPropertyToObject(property)
		if diags != nil && diags.HasError() {
			return basetypes.NewListNull(domainObjectType), diags
		}
		properties[i] = propertyObject
	}

	result, diags := types.ListValueFrom(ctx, opslevelPropertyObjectType, properties)
	if diags != nil && diags.HasError() {
		return basetypes.NewListNull(domainObjectType), diags
	}

	return result, nil
}

func opslevelPropertyToObject(opslevelProperty opslevel.Property) (basetypes.ObjectValue, diag.Diagnostics) {
	propertyAttrs := make(map[string]attr.Value)

	propertyAttrs["definition"] = types.StringValue(string(opslevelProperty.Definition.Id))
	propertyAttrs["owner"] = types.StringValue(string(opslevelProperty.Owner.Id()))
	if opslevelProperty.Value == nil {
		propertyAttrs["value"] = basetypes.NewStringNull()
	} else {
		propertyAttrs["value"] = types.StringValue(string(*opslevelProperty.Value))
	}

	parsedProperty, diags := types.ObjectValue(opslevelPropertyObjectType.AttrTypes, propertyAttrs)
	if diags != nil && diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	return parsedProperty, nil
}
