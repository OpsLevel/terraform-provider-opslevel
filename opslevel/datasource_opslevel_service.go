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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
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
	Alias                      types.String           `tfsdk:"alias"`
	Aliases                    types.List             `tfsdk:"aliases"`
	ApiDocumentPath            types.String           `tfsdk:"api_document_path"`
	Description                types.String           `tfsdk:"description"`
	Framework                  types.String           `tfsdk:"framework"`
	Id                         types.String           `tfsdk:"id"`
	Language                   types.String           `tfsdk:"language"`
	LifecycleAlias             types.String           `tfsdk:"lifecycle_alias"`
	Name                       types.String           `tfsdk:"name"`
	Owner                      types.String           `tfsdk:"owner"`
	OwnerId                    types.String           `tfsdk:"owner_id"`
	PreferredApiDocumentSource types.String           `tfsdk:"preferred_api_document_source"`
	Product                    types.String           `tfsdk:"product"`
	Properties                 []propertyModel        `tfsdk:"properties"`
	Repositories               types.List             `tfsdk:"repositories"`
	System                     *systemDataSourceModel `tfsdk:"system"`
	Tags                       types.List             `tfsdk:"tags"`
	TierAlias                  types.String           `tfsdk:"tier_alias"`
}

type propertyModel struct {
	Definition propertyDefinitionModel `tfsdk:"definition"`
	Value      types.String            `tfsdk:"value"`
}

type propertyDefinitionModel struct {
	Aliases types.List   `tfsdk:"aliases"`
	Id      types.String `tfsdk:"id"`
}

func NewPropertyModel(opslevelProperty opslevel.Property) propertyModel {
	aliases := OptionalStringListValue(opslevelProperty.Definition.Aliases)
	propModel := propertyModel{
		Definition: propertyDefinitionModel{
			Id:      ComputedStringValue(string(opslevelProperty.Definition.Id)),
			Aliases: aliases,
		},
	}
	if opslevelProperty.Value != nil {
		propModel.Value = ComputedStringValue(string(*opslevelProperty.Value))
	}
	return propModel
}

func NewPropertiesAllModel(ctx context.Context, opslevelProperties []opslevel.Property) ([]propertyModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	propertiesModel := []propertyModel{}
	for _, property := range opslevelProperties {
		propertiesModel = append(propertiesModel, NewPropertyModel(property))
	}
	return propertiesModel, diags
}

var opslevelPropertyAttrs = map[string]schema.Attribute{
	"definition": schema.SingleNestedAttribute{
		Description: "",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				Description: "A list of human-friendly, unique identifiers of the property definition.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Description: "The id of the property definition.",
				Computed:    true,
			},
		},
	},
	"value": schema.StringAttribute{
		Description: "The value of the custom property.",
		Computed:    true,
	},
}

func NewServiceDataSourceModel(ctx context.Context, service opslevel.Service, alias string) ServiceDataSourceModel {
	serviceDataSourceModel := ServiceDataSourceModel{
		Alias:           OptionalStringValue(alias),
		ApiDocumentPath: ComputedStringValue(service.ApiDocumentPath),
		Description:     ComputedStringValue(service.Description),
		Framework:       ComputedStringValue(service.Framework),
		Id:              OptionalStringValue(string(service.Id)),
		Language:        ComputedStringValue(service.Language),
		LifecycleAlias:  ComputedStringValue(service.Lifecycle.Alias),
		Name:            ComputedStringValue(service.Name),
		Owner:           ComputedStringValue(service.Owner.Alias),
		OwnerId:         ComputedStringValue(string(service.Owner.Id)),
		Product:         ComputedStringValue(service.Product),
		TierAlias:       ComputedStringValue(service.Tier.Alias),
	}

	serviceAliases := OptionalStringListValue(service.Aliases)
	serviceDataSourceModel.Aliases = serviceAliases

	if service.PreferredApiDocumentSource != nil {
		serviceDataSourceModel.PreferredApiDocumentSource = types.StringValue(string(*service.PreferredApiDocumentSource))
	}

	if service.Tags == nil {
		serviceDataSourceModel.Tags = types.ListNull(types.StringType)
	} else {
		serviceDataSourceModel.Tags = OptionalStringListValue(flattenTagArray(service.Tags.Nodes))
	}

	return serviceDataSourceModel
}

func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service data source",

		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "An alias of the service to find by.",
				Computed:    true,
				Optional:    true,
			},
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
			"id": schema.StringAttribute{
				Description: "The id of the service to find",
				Computed:    true,
				Optional:    true,
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
			"properties": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: opslevelPropertyAttrs,
				},
				Description: "Custom properties assigned to this service.",
				Computed:    true,
			},
			"repositories": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of repositories connected to the service.",
				Computed:    true,
			},
			"system": schema.ObjectAttribute{
				Description: "The system that the service belongs to.",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"aliases":     types.ListType{ElemType: types.StringType},
					"description": types.StringType,
					"domain":      types.StringType,
					"id":          types.StringType,
					"name":        types.StringType,
					"owner":       types.StringType,
				},
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
	var diags diag.Diagnostics
	var service opslevel.Service
	var err error

	planModel := read[ServiceDataSourceModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	if opslevel.IsID(planModel.Id.ValueString()) {
		service, err = getServiceWithId(*d.client, planModel)
	} else if planModel.Alias.ValueString() != "" {
		service, err = getServiceWithAlias(*d.client, planModel)
	} else {
		resp.Diagnostics.AddError("Config Error", "'alias' or valid 'id' for opslevel_service datasource must be set")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service datasource, got error: %s", err))
		return
	}

	stateModel := NewServiceDataSourceModel(ctx, service, planModel.Alias.ValueString())

	// Get full system data using the fixed GetSystem() method
	if service.Parent != nil {
		system, err := service.GetSystem(d.client, nil)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("system"),
				"OpsLevel Client Error",
				fmt.Sprintf("unable to read System for service, got error: %s", err),
			)
		} else if system != nil {
			systemModel := newSystemDataSourceModel(*system)
			stateModel.System = &systemModel
		}
	}

	// NOTE: service's hydrate does not populate properties
	properties, err := service.GetProperties(d.client, nil)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("properties"),
			"OpsLevel Client Error",
			fmt.Sprintf("unable to read Properties for service, got error: %s", err),
		)
	} else if properties != nil {
		stateModel.Properties, diags = NewPropertiesAllModel(ctx, properties.Nodes)
		resp.Diagnostics.Append(diags...)
	}

	if service.Repositories == nil {
		stateModel.Repositories = types.ListNull(types.StringType)
	} else {
		stateModel.Repositories, diags = types.ListValueFrom(
			ctx,
			types.StringType,
			flattenServiceRepositoriesArray(service.Repositories),
		)
		resp.Diagnostics.Append(diags...)
	}

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Service data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func getServiceWithAlias(client opslevel.Client, data ServiceDataSourceModel) (opslevel.Service, error) {
	service, err := client.GetServiceWithAlias(data.Alias.ValueString())
	if err != nil {
		return opslevel.Service{}, err
	}
	return *service, nil
}

func getServiceWithId(client opslevel.Client, data ServiceDataSourceModel) (opslevel.Service, error) {
	service, err := client.GetService(data.Id.ValueString())
	if err != nil {
		return opslevel.Service{}, err
	}
	return *service, nil
}
