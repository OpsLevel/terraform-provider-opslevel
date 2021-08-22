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

		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap:   map[string]*schema.Resource{},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			token := d.Get("token").(string)
			log.Println("[INFO] Initializing OpsLevel client")
			client := opslevel.NewClient(token)
			return client, nil
		},
	}
}
