package opslevel

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

const defaultApiTimeout = int64(30)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.ProviderWithValidateConfig = &OpslevelProvider{}

type OpslevelProvider struct {
	version string
}

type OpslevelProviderModel struct {
	ApiToken   types.String `tfsdk:"api_token"`
	ApiUrl     types.String `tfsdk:"api_url"`
	ApiTimeout types.Int64  `tfsdk:"api_timeout"`
}

func (p *OpslevelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opslevel"
	resp.Version = p.version
}

func (p *OpslevelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Optional:    true,
				Description: "The API authorization token. It can also be sourced from the OPSLEVEL_API_TOKEN environment variable.",
				Sensitive:   true,
			},
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "The url of the OpsLevel API to. It can also be sourced from the OPSLEVEL_API_URL environment variable.",
				Sensitive:   false,
			},
			"api_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Value (in seconds) to use for the timeout of API calls made.  It can also be sourced from the OPSLEVEL_API_TIMEOUT environment variable.",
				Sensitive:   false,
			},
		},
	}
}

func (p *OpslevelProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var providerModel OpslevelProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &providerModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if providerModel.ApiToken.IsNull() && os.Getenv("OPSLEVEL_API_TOKEN") == "" {
		resp.Diagnostics.AddError(
			"Provider Config Error",
			"An OPSLEVEL_API_TOKEN is needed to authenticate with the opslevel client. "+
				"This can be set as an 'OPSLEVEL_API_TOKEN' environment variable or in the provider configuration block as 'api_token'.",
		)
	}
}

func configApiToken(data *OpslevelProviderModel, resp *provider.ConfigureResponse) {
	if data.ApiToken.ValueString() != "" {
		return
	}

	if apiToken, ok := os.LookupEnv("OPSLEVEL_API_TOKEN"); ok {
		data.ApiToken = types.StringValue(apiToken)
		return
	}

	resp.Diagnostics.AddError(
		"Missing OPSLEVEL_API_TOKEN",
		"An OPSLEVEL_API_TOKEN is needed to authenticate with the opslevel client. "+
			"This can be set as an environment variable or in the provider configuration block as 'api_token'.",
	)
}

func configApiUrl(data *OpslevelProviderModel) {
	if data.ApiUrl.IsNull() || data.ApiUrl.Equal(types.StringValue("")) {
		if apiUrl, ok := os.LookupEnv("OPSLEVEL_API_URL"); ok {
			data.ApiUrl = types.StringValue(apiUrl)
		} else {
			data.ApiUrl = types.StringValue("https://api.opslevel.com/")
		}
	}
}

func configApiTimeOut(data *OpslevelProviderModel, resp *provider.ConfigureResponse) {
	if data.ApiTimeout.ValueInt64() > 0 {
		return
	}

	apiTimeout, ok := os.LookupEnv("OPSLEVEL_API_TIMEOUT")
	if !ok {
		data.ApiTimeout = types.Int64Value(defaultApiTimeout)
		return
	}

	if timeout, err := strconv.Atoi(apiTimeout); err == nil {
		data.ApiTimeout = types.Int64Value(int64(timeout))
		return
	}

	// Display warning when OPSLEVEL_API_TIMEOUT is set to to an invalid value
	resp.Diagnostics.AddWarning(
		"Expected OPSLEVEL_API_TIMEOUT to be an int",
		fmt.Sprintf(
			"OPSLEVEL_API_TIMEOUT was set to '%s'. The default timeout value of %d seconds will be used.",
			apiTimeout,
			defaultApiTimeout,
		),
	)
}

func (p *OpslevelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpslevelProviderModel
	tflog.Info(ctx, "Initializing opslevel client")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	tflog.Debug(ctx, "Setting opslevel client API token...")
	configApiToken(&data, resp)
	tflog.Debug(ctx, "opslevel client API token is set")

	tflog.Debug(ctx, "Setting opslevel client API endpoint URL...")
	configApiUrl(&data)
	tflog.Debug(ctx, "opslevel client API endpoint URL is set")

	tflog.Debug(ctx, "Setting opslevel client API timeout...")
	configApiTimeOut(&data, resp)
	tflog.Debug(ctx, "opslevel client API timeout is set")

	opts := []opslevel.Option{
		opslevel.SetAPIToken(data.ApiToken.ValueString()),
		opslevel.SetURL(data.ApiUrl.ValueString()),
		opslevel.SetTimeout(time.Second * time.Duration(data.ApiTimeout.ValueInt64())),
		opslevel.SetUserAgentExtra(fmt.Sprintf("terraform-provider-%s", p.version)),
	}
	client := opslevel.NewGQLClient(opts...)

	tflog.Debug(ctx, "Validating OpsLevel client...")
	if err := client.Validate(); err != nil {
		tflog.Error(ctx, fmt.Sprintf("OpsLevel client validation error: %s", err))
	}
	tflog.Debug(ctx, "OpsLevel client is valid")
	tflog.Info(ctx, "OpsLevel client is initialized")

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpslevelProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAliasResource,
		NewCheckAlertSourceUsageResource,
		NewCheckCustomEventResource,
		NewCheckGitBranchProtectionResource,
		NewCheckHasDocumentationResource,
		NewCheckHasRecentDeployResource,
		NewCheckManualResource,
		NewCheckPackageVersionResource,
		NewCheckRepositoryFileResource,
		NewCheckRepositoryGrepResource,
		NewCheckRepositoryIntegratedResource,
		NewCheckRepositorySearchResource,
		NewCheckServiceConfigurationResource,
		NewCheckServiceDependencyResource,
		NewCheckServiceOwnershipResource,
		NewCheckServicePropertyResource,
		NewCheckTagDefinedResource,
		NewCheckToolUsageResource,
		NewDomainResource,
		NewFilterResource,
		NewInfrastructureResource,
		NewIntegrationAwsResource,
		NewIntegrationAzureResourcesResource,
		NewIntegrationEndpointResource,
		NewIntegrationGoogleCloudResource,
		NewPropertyAssignmentResource,
		NewPropertyDefinitionResource,
		NewRepositoryResource,
		NewRubricCategoryResource,
		NewRubricLevelResource,
		NewScorecardResource,
		NewSecretResource,
		NewServiceDependencyResource,
		NewServiceRepositoryResource,
		NewServiceResource,
		NewServiceTagResource,
		NewServiceToolResource,
		NewSystemResource,
		NewTagResource,
		NewTeamContactResource,
		NewTeamResource,
		NewTeamTagResource,
		NewTriggerDefinitionResource,
		NewUserResource,
		NewWebhookActionResource,
	}
}

func (p *OpslevelProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCategoryDataSource,
		NewCategoryDataSourcesAll,
		NewDomainDataSource,
		NewDomainDataSourcesAll,
		NewFilterDataSource,
		NewFilterDataSourcesAll,
		NewIntegrationDataSource,
		NewIntegrationDataSourcesAll,
		NewLevelDataSource,
		NewLevelDataSourcesAll,
		NewLifecycleDataSource,
		NewLifecycleDataSourcesAll,
		NewPropertyDefinitionDataSource,
		NewPropertyDefinitionDataSourcesAll,
		NewRepositoriesDataSourceAll,
		NewRepositoryDataSource,
		NewScorecardDataSource,
		NewScorecardDataSourcesAll,
		NewServiceDataSource,
		NewServiceDependenciesDataSource,
		NewServiceDataSourcesAll,
		NewSystemDataSource,
		NewSystemDataSourcesAll,
		NewTeamDataSource,
		NewTeamDataSourcesAll,
		NewTierDataSource,
		NewTierDataSourcesAll,
		NewUserDataSource,
		NewUserDataSourcesAll,
		NewWebhookActionDataSource,
		NewWebhookActionDataSourcesAll,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpslevelProvider{version: version}
	}
}
