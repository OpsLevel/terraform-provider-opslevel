package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure ServiceDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &ServiceDataSource{}

func NewServiceDataSourcesAll() datasource.DataSource {
	return &ServiceDataSourcesAll{}
}

// ServiceDataSource manages a Service data source.
type ServiceDataSourcesAll struct {
	CommonDataSourceClient
}

// serviceDataSourcesAllModel describes the data source data model.
type serviceDataSourcesAllModel struct {
	Filter   *filterBlockModel                `tfsdk:"filter"`
	Services []serviceMinimalaDataSourceModel `tfsdk:"services"`
}

func NewServiceDataSourcesAllModel(services []opslevel.Service) serviceDataSourcesAllModel {
	serviceDataSourcesModel := []serviceMinimalaDataSourceModel{}
	for _, service := range services {
		serviceModel := serviceMinimalaDataSourceModel{
			Id:   ComputedStringValue(string(service.Id)),
			Name: ComputedStringValue(service.Name),
			Url:  ComputedStringValue(service.HtmlURL),
		}
		serviceDataSourcesModel = append(serviceDataSourcesModel, serviceModel)
	}
	return serviceDataSourcesAllModel{Services: serviceDataSourcesModel}
}

func (d *ServiceDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

var serviceMinimalSchemaAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The id of the service",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The display name of the service.",
		Computed:    true,
	},
	"url": schema.StringAttribute{
		Description: "A link to the HTML page for the resource",
		Computed:    true,
	},
}

type serviceMinimalaDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Url  types.String `tfsdk:"url"`
}

func (d *ServiceDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"filter", "framework", "language", "lifecycle", "owner", "product", "tag", "tier"}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Services data source",

		Attributes: map[string]schema.Attribute{
			"filter": schema.SingleNestedAttribute{
				Description: fmt.Sprintf(
					"Used to filter services by one of '%s'",
					strings.Join(validFieldNames, "`, `"),
				),
				Optional:   true,
				Attributes: FilterAttrs(validFieldNames),
			},
			"services": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: serviceMinimalSchemaAttrs,
				},
				Description: "List of Service data sources",
				Computed:    true,
			},
		},
	}
}

func (d *ServiceDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel serviceDataSourcesAllModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var services *opslevel.ServiceConnection
	var err error

	if planModel.Filter == nil {
		services, err = d.client.ListServices(nil)
	} else {
		switch planModel.Filter.Field.ValueString() {
		case "filter":
			filterId := planModel.Filter.Value.ValueString()
			if opslevel.IsID(filterId) {
				services, err = d.client.ListServicesWithFilter(filterId, nil)
			} else {
				resp.Diagnostics.AddError("Config Error",
					fmt.Sprintf("'value' field in filter block must be a valid ID. Given '%s'", filterId),
				)
			}
		case "framework":
			services, err = d.client.ListServicesWithFramework(planModel.Filter.Value.ValueString(), nil)
		case "language":
			services, err = d.client.ListServicesWithLanguage(planModel.Filter.Value.ValueString(), nil)
		case "lifecycle":
			services, err = d.client.ListServicesWithLifecycle(planModel.Filter.Value.ValueString(), nil)
		case "owner":
			services, err = d.client.ListServicesWithOwner(planModel.Filter.Value.ValueString(), nil)
		case "product":
			services, err = d.client.ListServicesWithProduct(planModel.Filter.Value.ValueString(), nil)
		case "tag":
			tagArgs, tagArgsErr := opslevel.NewTagArgs(planModel.Filter.Value.ValueString())
			if tagArgsErr != nil {
				resp.Diagnostics.AddError("Client Error",
					fmt.Sprintf("Unable to create TagArgs from '%s', got error: '%s'", planModel.Filter.Value.ValueString(), err),
				)
				return
			}
			services, err = d.client.ListServicesWithTag(tagArgs, nil)
		case "tier":
			services, err = d.client.ListServicesWithTier(planModel.Filter.Value.ValueString(), nil)
		default:
			services, err = d.client.ListServices(nil)
		}
	}

	if opslevel.HasBadHttpStatus(err) {
		resp.Diagnostics.AddError("HTTP status error", fmt.Sprintf("Unable to list services datasource, got error: %s", err))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list services datasource, got error: %s", err))
		return
	}

	if services == nil {
		stateModel = NewServiceDataSourcesAllModel([]opslevel.Service{})
	} else {
		stateModel = NewServiceDataSourcesAllModel(services.Nodes)
	}
	stateModel.Filter = planModel.Filter

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all OpsLevel Service data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
