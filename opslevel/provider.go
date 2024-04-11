package opslevel

import (
	"context"
	"fmt"
	// "log"
	"os"
	"strconv"
	"time"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/terraform"
	// "github.com/opslevel/opslevel-go/v2024"
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
var _ provider.Provider = &OpslevelProvider{}

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
				Required:    true,
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
				Description: "Value (in seconds) to use for the timeout of API calls made",
				Sensitive:   false,
			},
		},
	}
}

func configApiToken(data *OpslevelProviderModel, resp *provider.ConfigureResponse) {
	if data.ApiUrl.IsNull() || data.ApiToken.Equal(types.StringValue("")) {
		if apiToken, ok := os.LookupEnv("OPSLEVEL_API_TOKEN"); ok {
			data.ApiToken = types.StringValue(apiToken)
		} else {
			resp.Diagnostics.AddError(
				"Missing OPSLEVEL_API_TOKEN",
				"An OPSLEVEL_API_TOKEN is needed to authenticate with the opslevel client. "+
					"This can be set as an environment variable or in the provider configuration block as 'api_token'.",
			)
		}
	}
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

	if apiTimeout, ok := os.LookupEnv("OPSLEVEL_API_TIMEOUT"); ok {
		if timeout, err := strconv.Atoi(apiTimeout); err == nil {
			data.ApiTimeout = types.Int64Value(int64(timeout))
			return
		}
	}

	resp.Diagnostics.AddWarning(
		"Expected OPSLEVEL_API_TIMEOUT to be an int",
		fmt.Sprintf(
			"OPSLEVEL_API_TIMEOUT was set but not as an int. The default timeout value of %d seconds will be used.",
			defaultApiTimeout,
		),
	)
	data.ApiTimeout = types.Int64Value(defaultApiTimeout)
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
		NewCheckManualResource,
		NewCheckToolUsageResource,
		NewDomainResource,
		NewInfrastructureResource,
		NewRubricCategoryResource,
		NewRubricLevelResource,
		NewScorecardResource,
		NewSecretResource,
		NewServiceResource,
		NewUserResource,
		NewWebhookActionResource,
	}
}

func (p *OpslevelProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCategoryDataSource,
		NewDomainDataSource,
		NewDomainDataSourcesAll,
		NewFilterDataSource,
		NewIntegrationDataSource,
		NewLevelDataSource,
		NewLifecycleDataSource,
		NewPropertyDefinitionDataSource,
		NewRepositoryDataSource,
		NewScorecardDataSource,
		NewServiceDataSource,
		NewSystemDataSource,
		NewTeamDataSource,
		NewTierDataSource,
		NewUserDataSource,
		NewWebhookActionDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpslevelProvider{version: version}
	}
}

// func Provider() terraform.ResourceProvider {
// 	return &schema.Provider{
// 		Schema: map[string]*schema.Schema{
// 			"api_url": {
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_URL", "https://api.opslevel.com/"),
// 				Description: "The url of the OpsLevel API to. It can also be sourced from the OPSLEVEL_API_URL environment variable.",
// 				Sensitive:   false,
// 			},
// 			"api_token": {
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_TOKEN", ""),
// 				Description: "The API authorization token. It can also be sourced from the OPSLEVEL_API_TOKEN environment variable.",
// 				Sensitive:   true,
// 			},
// 			"api_timeout": {
// 				Type:        schema.TypeInt,
// 				Optional:    true,
// 				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_TIMEOUT", "30"),
// 				Description: "Value (in seconds) to use for the timeout of API calls made",
// 				Sensitive:   false,
// 			},
// 		},

// 		DataSourcesMap: map[string]*schema.Resource{
// 			"opslevel_domain":               datasourceDomain(),
// 			"opslevel_domains":              datasourceDomains(),
// 			"opslevel_filter":               datasourceFilter(),
// 			"opslevel_filters":              datasourceFilters(),
// 			"opslevel_integration":          datasourceIntegration(),
// 			"opslevel_integrations":         datasourceIntegrations(),
// 			"opslevel_lifecycle":            datasourceLifecycle(),
// 			"opslevel_lifecycles":           datasourceLifecycles(),
// 			"opslevel_property_definition":  datasourcePropertyDefinition(),
// 			"opslevel_property_definitions": datasourcePropertyDefinitions(),
// 			"opslevel_repository":           datasourceRepository(),
// 			"opslevel_repositories":         datasourceRepositories(),
// 			"opslevel_rubric_category":      datasourceRubricCategory(),
// 			"opslevel_rubric_categories":    datasourceRubricCategories(),
// 			"opslevel_rubric_level":         datasourceRubricLevel(),
// 			"opslevel_rubric_levels":        datasourceRubricLevels(),
// 			"opslevel_scorecard":            datasourceScorecard(),
// 			"opslevel_scorecards":           datasourceScorecards(),
// 			"opslevel_service":              datasourceService(),
// 			"opslevel_services":             datasourceServices(),
// 			"opslevel_system":               datasourceSystem(),
// 			"opslevel_systems":              datasourceSystems(),
// 			"opslevel_team":                 datasourceTeam(),
// 			"opslevel_teams":                datasourceTeams(),
// 			"opslevel_tier":                 datasourceTier(),
// 			"opslevel_tiers":                datasourceTiers(),
// 			"opslevel_users":                datasourceUsers(),
// 		},
// 		ResourcesMap: map[string]*schema.Resource{
// 			"opslevel_check_alert_source_usage":    resourceCheckAlertSourceUsage(),
// 			"opslevel_check_custom_event":          resourceCheckCustomEvent(),
// 			"opslevel_check_git_branch_protection": resourceCheckGitBranchProtection(),
// 			"opslevel_check_has_documentation":     resourceCheckHasDocumentation(),
// 			"opslevel_check_has_recent_deploy":     resourceCheckHasRecentDeploy(),
// 			"opslevel_check_manual":                resourceCheckManual(),
// 			"opslevel_check_repository_file":       resourceCheckRepositoryFile(),
// 			"opslevel_check_repository_grep":       resourceCheckRepositoryGrep(),
// 			"opslevel_check_repository_integrated": resourceCheckRepositoryIntegrated(),
// 			"opslevel_check_repository_search":     resourceCheckRepositorySearch(),
// 			"opslevel_check_service_dependency":    resourceCheckServiceDependency(),
// 			"opslevel_check_service_configuration": resourceCheckServiceConfiguration(),
// 			"opslevel_check_service_ownership":     resourceCheckServiceOwnership(),
// 			"opslevel_check_service_property":      resourceCheckServiceProperty(),
// 			"opslevel_check_tag_defined":           resourceCheckTagDefined(),
// 			"opslevel_check_tool_usage":            resourceCheckToolUsage(),
// 			"opslevel_domain":                      resourceDomain(),
// 			"opslevel_filter":                      resourceFilter(),
// 			"opslevel_infrastructure":              resourceInfrastructure(),
// 			"opslevel_integration_aws":             resourceIntegrationAWS(),
// 			"opslevel_property_assignment":         resourcePropertyAssignment(),
// 			"opslevel_property_definition":         resourcePropertyDefinition(),
// 			"opslevel_repository":                  resourceRepository(),
// 			"opslevel_rubric_level":                resourceRubricLevel(),
// 			"opslevel_rubric_category":             resourceRubricCategory(),
// 			"opslevel_scorecard":                   resourceScorecard(),
// 			"opslevel_secret":                      resourceSecret(),
// 			"opslevel_service":                     resourceService(),
// 			"opslevel_service_dependency":          resourceServiceDependency(),
// 			"opslevel_service_repository":          resourceServiceRepository(),
// 			"opslevel_service_tag":                 resourceServiceTag(),
// 			"opslevel_service_tool":                resourceServiceTool(),
// 			"opslevel_system":                      resourceSystem(),
// 			"opslevel_tag":                         resourceTag(),
// 			"opslevel_team_contact":                resourceTeamContact(),
// 			"opslevel_team_tag":                    resourceTeamTag(),
// 			"opslevel_team":                        resourceTeam(),
// 			"opslevel_trigger_definition":          resourceTriggerDefinition(),
// 			"opslevel_user":                        resourceUser(),
// 			"opslevel_webhook_action":              resourceWebhookAction(),
// 		},

// 		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
// 			url := d.Get("api_url").(string)
// 			token := d.Get("api_token").(string)
// 			timeout := d.Get("api_timeout").(int)
// 			if timeout <= 0 {
// 				timeout = 10
// 			}
// 			log.Println("[INFO] Initializing OpsLevel client")

// 			opts := make([]opslevel.Option, 0)

// 			opts = append(opts, opslevel.SetAPIToken(token))
// 			opts = append(opts, opslevel.SetURL(url))
// 			opts = append(opts, opslevel.SetUserAgentExtra(fmt.Sprintf("terraform-provider-%s", version)))
// 			opts = append(opts, opslevel.SetTimeout(time.Second*time.Duration(timeout)))

// 			client := opslevel.NewGQLClient(opts...)

// 			return client, client.Validate()
// 		},
// 	}
// }

// func GetString(d *schema.ResourceData, key string) *string {
// 	value := d.Get(key).(string)
// 	return &value
// }
