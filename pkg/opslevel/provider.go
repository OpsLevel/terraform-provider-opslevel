package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opslevel/kubectl-opslevel/opslevel"
)

const defaultUrl = "https://api.opslevel.com/graphql"

type provider struct {
	client *opslevel.Client
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPSLEVEL_TOKEN", ""),
				Description: "The OpsLevel API token to use for authentication.",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opslevel_service": datasourceOpsLevelService(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"opslevel_service": resourceOpsLevelService(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)

	client := opslevel.NewClient(token)

	return provider{client}, nil
}
