package opslevel

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opslevel/opslevel-go/v2022"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_URL", "https://api.opslevel.com/"),
				Description: "The url of the OpsLevel API to. It can also be sourced from the OPSLEVEL_API_URL environment variable.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_API_TOKEN", ""),
				Description: "The API authorization token. It can also be sourced from the OPSLEVEL_API_TOKEN environment variable.",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opslevel_filter":            datasourceFilter(),
			"opslevel_filters":           datasourceFilters(),
			"opslevel_group":             datasourceGroup(),
			"opslevel_groups":            datasourceGroups(),
			"opslevel_integration":       datasourceIntegration(),
			"opslevel_integrations":      datasourceIntegrations(),
			"opslevel_lifecycle":         datasourceLifecycle(),
			"opslevel_lifecycles":        datasourceLifecycles(),
			"opslevel_repository":        datasourceRepository(),
			"opslevel_repositories":      datasourceRepositories(),
			"opslevel_rubric_category":   datasourceRubricCategory(),
			"opslevel_rubric_categories": datasourceRubricCategories(),
			"opslevel_rubric_level":      datasourceRubricLevel(),
			"opslevel_rubric_levels":     datasourceRubricLevels(),
			"opslevel_service":           datasourceService(),
			"opslevel_services":          datasourceServices(),
			"opslevel_team":              datasourceTeam(),
			"opslevel_teams":             datasourceTeams(),
			"opslevel_tier":              datasourceTier(),
			"opslevel_tiers":             datasourceTiers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"opslevel_check_alert_source_usage":    resourceCheckAlertSourceUsage(),
			"opslevel_check_custom_event":          resourceCheckCustomEvent(),
			"opslevel_check_git_branch_protection": resourceCheckGitBranchProtection(),
			"opslevel_check_has_documentation":     resourceCheckHasDocumentation(),
			"opslevel_check_has_recent_deploy":     resourceCheckHasRecentDeploy(),
			"opslevel_check_manual":                resourceCheckManual(),
			"opslevel_check_repository_file":       resourceCheckRepositoryFile(),
			"opslevel_check_repository_integrated": resourceCheckRepositoryIntegrated(),
			"opslevel_check_repository_search":     resourceCheckRepositorySearch(),
			"opslevel_check_service_dependency":    resourceCheckServiceDependency(),
			"opslevel_check_service_configuration": resourceCheckServiceConfiguration(),
			"opslevel_check_service_ownership":     resourceCheckServiceOwnership(),
			"opslevel_check_service_property":      resourceCheckServiceProperty(),
			"opslevel_check_tag_defined":           resourceCheckTagDefined(),
			"opslevel_check_tool_usage":            resourceCheckToolUsage(),
			"opslevel_filter":                      resourceFilter(),
			"opslevel_group":                       resourceGroup(),
			"opslevel_rubric_level":                resourceRubricLevel(),
			"opslevel_rubric_category":             resourceRubricCategory(),
			"opslevel_service":                     resourceService(),
			"opslevel_service_repository":          resourceServiceRepository(),
			"opslevel_service_tag":                 resourceServiceTag(),
			"opslevel_service_tool":                resourceServiceTool(),
			"opslevel_team_contact":                resourceTeamContact(),
			"opslevel_team":                        resourceTeam(),
			"opslevel_user":                        resourceUser(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			url := d.Get("api_url").(string)
			token := d.Get("api_token").(string)
			log.Println("[INFO] Initializing OpsLevel client")
			client := opslevel.NewGQLClient(opslevel.SetAPIToken(token), opslevel.SetURL(url), opslevel.SetUserAgentExtra(fmt.Sprintf("terraform-provider-%s", version)))
			return client, nil
		},
	}
}
