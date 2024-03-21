package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	Alias                      types.String `tfsdk:"alias"`
	Aliases                    types.List   `tfsdk:"aliases"`
	ApiDocumentPath            types.String `tfsdk:"api_document_path"`
	Description                types.String `tfsdk:"description"`
	Framework                  types.String `tfsdk:"framework"`
	Id                         types.String `tfsdk:"id"`
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

func NewServiceDataSourceModel(ctx context.Context, service opslevel.Service) (ServiceDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	serviceDataSourceModel := ServiceDataSourceModel{
		ApiDocumentPath: types.StringValue(service.ApiDocumentPath),
		Description:     types.StringValue(service.Description),
		Id:              types.StringValue(string(service.Id)),
		Framework:       types.StringValue(service.Framework),
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

	if service.Properties == nil {
		serviceDataSourceModel.Properties = types.ListNull(types.StringType)
	} else {
		serviceProperties, propsDiags := types.ListValueFrom(ctx, types.StringType, service.Properties.Nodes)
		serviceDataSourceModel.Properties = serviceProperties
		diags = append(diags, propsDiags...)
	}

	if service.Repositories == nil {
		serviceDataSourceModel.Repositories = types.ListNull(types.StringType)
	} else {
		repositories, tagsDiags := types.ListValueFrom(ctx, types.StringType, flattenServiceRepositoriesArray(service.Repositories))
		serviceDataSourceModel.Repositories = repositories
		diags = append(diags, tagsDiags...)
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
				MarkdownDescription: "The id of the service to find.",
				Optional:            true,
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
				ElementType: types.StringType,
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
	var data ServiceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, err := d.getService(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service datasource, got error: %s", err))
		return
	}
	serviceDataModel, diags := NewServiceDataSourceModel(ctx, service)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Service data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &serviceDataModel)...)
}

func (d *ServiceDataSource) getService(data ServiceDataSourceModel) (opslevel.Service, error) {
	var err error
	var service *opslevel.Service
	alias := data.Alias.ValueString()
	id := data.Id.ValueString()

	if id != "" {
		service, err = d.client.GetService(opslevel.ID(id))
	} else if alias != "" {
		service, err = d.client.GetServiceWithAlias(alias)
	} else {
		return opslevel.Service{}, fmt.Errorf("unable to read config for service datasource - must have alias or id")
	}
	if err != nil {
		return opslevel.Service{}, err
	}

	if service == nil || service.Id == "" {
		return opslevel.Service{}, fmt.Errorf("unable to find repository with alias=`%s` or id=`%s`", alias, id)
	}

	return *service, nil
}
