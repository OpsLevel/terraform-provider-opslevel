package opslevel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDomainDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opslevel_domain.test", "identifier", "my_domain"),
				),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDomainsAllDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAllDomainsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opslevel_domains.all", "identifier", "my_domain"),
				),
				PlanOnly: true,
			},
		},
	})
}

const testAccDomainDataSourceConfig = `
data "opslevel_domain" "test" {
  identifier = "my_domain"
}
`

const testAccAllDomainsDataSourceConfig = `
data "opslevel_domains" "all" {
}
`
