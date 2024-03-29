package opslevel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opslevel/opslevel-go/v2024"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_URL", "https://api.opslevel.com/"),
				Description: "The url of the OpsLevel API to. It can also be sourced from the OPSLEVEL_API_URL environment variable.",
				Sensitive:   false,
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_TOKEN", ""),
				Description: "The API authorization token. It can also be sourced from the OPSLEVEL_API_TOKEN environment variable.",
				Sensitive:   true,
			},
			"api_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_TIMEOUT", "30"),
				Description: "Value (in seconds) to use for the timeout of API calls made",
				Sensitive:   false,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opslevel_domain":               datasourceDomain(),
			"opslevel_domains":              datasourceDomains(),
			"opslevel_filter":               datasourceFilter(),
			"opslevel_filters":              datasourceFilters(),
			"opslevel_integration":          datasourceIntegration(),
			"opslevel_integrations":         datasourceIntegrations(),
			"opslevel_lifecycle":            datasourceLifecycle(),
			"opslevel_lifecycles":           datasourceLifecycles(),
			"opslevel_property_definition":  datasourcePropertyDefinition(),
			"opslevel_property_definitions": datasourcePropertyDefinitions(),
			"opslevel_repository":           datasourceRepository(),
			"opslevel_repositories":         datasourceRepositories(),
			"opslevel_rubric_category":      datasourceRubricCategory(),
			"opslevel_rubric_categories":    datasourceRubricCategories(),
			"opslevel_rubric_level":         datasourceRubricLevel(),
			"opslevel_rubric_levels":        datasourceRubricLevels(),
			"opslevel_scorecard":            datasourceScorecard(),
			"opslevel_scorecards":           datasourceScorecards(),
			"opslevel_service":              datasourceService(),
			"opslevel_services":             datasourceServices(),
			"opslevel_system":               datasourceSystem(),
			"opslevel_systems":              datasourceSystems(),
			"opslevel_team":                 datasourceTeam(),
			"opslevel_teams":                datasourceTeams(),
			"opslevel_tier":                 datasourceTier(),
			"opslevel_tiers":                datasourceTiers(),
			"opslevel_users":                datasourceUsers(),
			"opslevel_webhook_action":       datasourceWebhookAction(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"opslevel_check_alert_source_usage":    resourceCheckAlertSourceUsage(),
			"opslevel_check_custom_event":          resourceCheckCustomEvent(),
			"opslevel_check_git_branch_protection": resourceCheckGitBranchProtection(),
			"opslevel_check_has_documentation":     resourceCheckHasDocumentation(),
			"opslevel_check_has_recent_deploy":     resourceCheckHasRecentDeploy(),
			"opslevel_check_manual":                resourceCheckManual(),
			"opslevel_check_repository_file":       resourceCheckRepositoryFile(),
			"opslevel_check_repository_grep":       resourceCheckRepositoryGrep(),
			"opslevel_check_repository_integrated": resourceCheckRepositoryIntegrated(),
			"opslevel_check_repository_search":     resourceCheckRepositorySearch(),
			"opslevel_check_service_dependency":    resourceCheckServiceDependency(),
			"opslevel_check_service_configuration": resourceCheckServiceConfiguration(),
			"opslevel_check_service_ownership":     resourceCheckServiceOwnership(),
			"opslevel_check_service_property":      resourceCheckServiceProperty(),
			"opslevel_check_tag_defined":           resourceCheckTagDefined(),
			"opslevel_check_tool_usage":            resourceCheckToolUsage(),
			"opslevel_domain":                      resourceDomain(),
			"opslevel_filter":                      resourceFilter(),
			"opslevel_infrastructure":              resourceInfrastructure(),
			"opslevel_integration_aws":             resourceIntegrationAWS(),
			"opslevel_property_assignment":         resourcePropertyAssignment(),
			"opslevel_property_definition":         resourcePropertyDefinition(),
			"opslevel_repository":                  resourceRepository(),
			"opslevel_rubric_level":                resourceRubricLevel(),
			"opslevel_rubric_category":             resourceRubricCategory(),
			"opslevel_scorecard":                   resourceScorecard(),
			"opslevel_secret":                      resourceSecret(),
			"opslevel_service":                     resourceService(),
			"opslevel_service_dependency":          resourceServiceDependency(),
			"opslevel_service_repository":          resourceServiceRepository(),
			"opslevel_service_tag":                 resourceServiceTag(),
			"opslevel_service_tool":                resourceServiceTool(),
			"opslevel_system":                      resourceSystem(),
			"opslevel_tag":                         resourceTag(),
			"opslevel_team_contact":                resourceTeamContact(),
			"opslevel_team_tag":                    resourceTeamTag(),
			"opslevel_team":                        resourceTeam(),
			"opslevel_trigger_definition":          resourceTriggerDefinition(),
			"opslevel_user":                        resourceUser(),
			"opslevel_webhook_action":              resourceWebhookAction(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			url := d.Get("api_url").(string)
			token := d.Get("api_token").(string)
			timeout := d.Get("api_timeout").(int)
			if timeout <= 0 {
				timeout = 10
			}
			log.Println("[INFO] Initializing OpsLevel client")

			opts := make([]opslevel.Option, 0)

			opts = append(opts, opslevel.SetAPIToken(token))
			opts = append(opts, opslevel.SetURL(url))
			opts = append(opts, opslevel.SetUserAgentExtra(fmt.Sprintf("terraform-provider-%s", version)))
			opts = append(opts, opslevel.SetTimeout(time.Second*time.Duration(timeout)))

			client := opslevel.NewGQLClient(opts...)

			return client, client.Validate()
		},
	}
}

func GetString(d *schema.ResourceData, key string) *string {
	value := d.Get(key).(string)
	return &value
}
