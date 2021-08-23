package opslevel

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opslevel/opslevel-go"
)

const defaultUrl = "https://api.opslevel.com/graphql"

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_APITOKEN", ""),
				Description: "The OpsLevel API token to use for authentication.",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opslevel_filter":            datasourceFilter(),
			"opslevel_filters":           datasourceFilters(),
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
			"opslevel_check_custom_event":          resourceCheckCustomEvent(),
			"opslevel_check_manual":                resourceCheckManual(),
			"opslevel_check_repository_file":       resourceCheckRepositoryFile(),
			"opslevel_check_repository_integrated": resourceCheckRepositoryIntegrated(),
			"opslevel_check_repository_search":     resourceCheckRepositorySearch(),
			"opslevel_check_service_configuration": resourceCheckServiceConfiguration(),
			"opslevel_check_service_owner":         resourceCheckServiceOwnership(),
			"opslevel_check_service_property":      resourceCheckServiceProperty(),
			"opslevel_check_tag_defined":           resourceCheckTagDefined(),
			"opslevel_check_tool_usage":            resourceCheckToolUsage(),
			"opslevel_filter":                      resourceFilter(),
			"opslevel_rubric_level":                resourceRubricLevel(),
			"opslevel_rubric_category":             resourceRubricCategory(),
			"opslevel_service":                     resourceService(),
			"opslevel_service_repository":          resourceServiceRepository(),
			"opslevel_service_tag":                 resourceServiceTag(),
			"opslevel_service_tool":                resourceServiceTool(),
			"opslevel_team":                        resourceTeam(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			token := d.Get("token").(string)
			log.Println("[INFO] Initializing OpsLevel client")
			client := opslevel.NewClient(token)
			return client, nil
		},
	}
}
