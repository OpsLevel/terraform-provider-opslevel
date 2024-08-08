package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure ServiceDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &ServiceDataSource{}

func NewServiceDependenciesDataSource() datasource.DataSource {
	return &ServiceDependenciesDataSource{}
}

// ServiceDataSource manages a Service data source.
type ServiceDependenciesDataSource struct {
	CommonDataSourceClient
}

// ServiceDependenciesModel describes the data source data model.
type ServiceDependenciesModel struct {
	Dependents   []dependentsModel   `tfsdk:"dependents"`
	Dependencies []dependenciesModel `tfsdk:"dependencies"`
	Service      types.String        `tfsdk:"service"`
}

type dependentsModel struct {
	Id     types.String `tfsdk:"id"`
	Locked types.Bool   `tfsdk:"locked"`
	Notes  types.String `tfsdk:"notes"`
}

type dependenciesModel struct {
	Id     types.String `tfsdk:"id"`
	Locked types.Bool   `tfsdk:"locked"`
	Notes  types.String `tfsdk:"notes"`
}

func NewServiceDependenciesModel(serviceIdentifier string, dependents []dependentsModel, dependencies []dependenciesModel) ServiceDependenciesModel {
	return ServiceDependenciesModel{
		Dependents:   dependents,
		Dependencies: dependencies,
		Service:      types.StringValue(serviceIdentifier),
	}
}

func (d *ServiceDependenciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_dependencies"
}

var depsAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The ID of the serviceDependency.",
		Computed:    true,
	},
	"locked": schema.BoolAttribute{
		Description: "Is the dependency locked by a service config?",
		Computed:    true,
	},
	"notes": schema.StringAttribute{
		Description: "Notes for service dependency.",
		Optional:    true,
	},
}

func (d *ServiceDependenciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Service Dependencies data source",

		Attributes: map[string]schema.Attribute{
			"dependents": schema.ListNestedAttribute{
				Description: "List of Service Dependents of a service",
				NestedObject: schema.NestedAttributeObject{
					Attributes: depsAttrs,
				},
				Computed: true,
			},
			"dependencies": schema.ListNestedAttribute{
				Description: "List of Service Dependencies of a service",
				NestedObject: schema.NestedAttributeObject{
					Attributes: depsAttrs,
				},
				Computed: true,
			},
			"service": schema.StringAttribute{
				Description: "The ID or alias of the service with the dependency.",
				Required:    true,
			},
		},
	}
}

func (d *ServiceDependenciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		dependencies      []dependenciesModel
		dependents        []dependentsModel
		err               error
		service           *opslevel.Service
		serviceIdentifier string
		stateModel        ServiceDependenciesModel
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service"), &serviceIdentifier)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve Service
	if opslevel.IsID(serviceIdentifier) {
		service, err = d.client.GetService(opslevel.ID(serviceIdentifier))
	} else {
		service, err = d.client.GetServiceWithAlias(serviceIdentifier)
	}
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}

	// List Service Dependents
	svcDependents, err := service.GetDependents(d.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get dependents for service, got error: %s", err))
		return
	}
	if svcDependents == nil || len(svcDependents.Edges) == 0 {
		dependents = []dependentsModel{}
	} else {
		for _, svcDependency := range svcDependents.Edges {
			dependents = append(dependents, dependentsModel{
				Locked: types.BoolValue(svcDependency.Locked),
				Id:     ComputedStringValue(string(svcDependency.Id)),
				Notes:  ComputedStringValue(svcDependency.Notes),
			})
		}
	}

	// List Service Dependencies
	svcDependencies, err := service.GetDependencies(d.client, nil)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get dependencies for service, got error: %s", err))
		return
	}
	if svcDependencies == nil || len(svcDependencies.Edges) == 0 {
		dependencies = []dependenciesModel{}
	} else {
		for _, svcDependency := range svcDependencies.Edges {
			dependencies = append(dependencies, dependenciesModel{
				Locked: types.BoolValue(svcDependency.Locked),
				Id:     ComputedStringValue(string(svcDependency.Id)),
				Notes:  ComputedStringValue(svcDependency.Notes),
			})
		}
	}

	stateModel = NewServiceDependenciesModel(serviceIdentifier, dependents, dependencies)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Service data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
