package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Alias                      types.String    `tfsdk:"alias"`
	Aliases                    types.List      `tfsdk:"aliases"`
	ApiDocumentPath            types.String    `tfsdk:"api_document_path"`
	Description                types.String    `tfsdk:"description"`
	Framework                  types.String    `tfsdk:"framework"`
	Id                         types.String    `tfsdk:"id"`
	Language                   types.String    `tfsdk:"language"`
	LifecycleAlias             types.String    `tfsdk:"lifecycle_alias"`
	Name                       types.String    `tfsdk:"name"`
	Owner                      types.String    `tfsdk:"owner"`
	OwnerId                    types.String    `tfsdk:"owner_id"`
	PreferredApiDocumentSource types.String    `tfsdk:"preferred_api_document_source"`
	Product                    types.String    `tfsdk:"product"`
	Properties                 []propertyModel `tfsdk:"properties"`
	Repositories               types.List      `tfsdk:"repositories"`
	Tags                       types.List      `tfsdk:"tags"`
	TierAlias                  types.String    `tfsdk:"tier_alias"`
}

type propertyModel struct {
	Definition propertyDefinitionModel `tfsdk:"definition"`
	Value      types.String            `tfsdk:"value"`
}

type propertyDefinitionModel struct {
	Aliases types.List   `tfsdk:"aliases"`
	Id      types.String `tfsdk:"id"`
}

func NewPropertyModel(ctx context.Context, opslevelProperty opslevel.Property) (propertyModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, opslevelProperty.Definition.Aliases)
	if diags != nil && diags.HasError() {
		return propertyModel{}, diags
	}
	propModel := propertyModel{
		Definition: propertyDefinitionModel{
			Id:      ComputedStringValue(string(opslevelProperty.Definition.Id)),
			Aliases: aliases,
		},
	}
	if opslevelProperty.Value != nil {
		propModel.Value = ComputedStringValue(string(*opslevelProperty.Value))
	}
	return propModel, diags
}

func NewPropertiesAllModel(ctx context.Context, opslevelProperties []opslevel.Property) ([]propertyModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	propertiesModel := []propertyModel{}
	for _, property := range opslevelProperties {
		propertyModel, propertyDiag := NewPropertyModel(ctx, property)
		diags.Append(propertyDiag...)
		propertiesModel = append(propertiesModel, propertyModel)
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

func NewServiceDataSourceModel(ctx context.Context, service opslevel.Service, alias string) (ServiceDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

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

	serviceAliases, svcDiags := OptionalStringListValue(ctx, service.Aliases)
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
	var planModel, stateModel ServiceDataSourceModel
	var service opslevel.Service
	var err error

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
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

	stateModel, diags := NewServiceDataSourceModel(ctx, service, planModel.Alias.ValueString())
	resp.Diagnostics.Append(diags...)

	// NOTE: service's hydrate does not populate properties
	properties, err := service.GetProperties(d.client, nil)
	if err != nil {
		diags.AddAttributeError(
			path.Root("properties"),
			"OpsLevel Client Error",
			fmt.Sprintf("unable to read Properties for service, got error: %s", err),
		)
		return
	}
	if properties != nil {
		stateModel.Properties, diags = NewPropertiesAllModel(ctx, properties.Nodes)
	}

	if service.Repositories == nil {
		stateModel.Repositories = types.ListNull(types.StringType)
	} else {
		stateModel.Repositories, diags = types.ListValueFrom(
			ctx,
			types.StringType,
			flattenServiceRepositoriesArray(service.Repositories),
		)
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
	service, err := client.GetService(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		return opslevel.Service{}, err
	}
	return *service, nil
}
