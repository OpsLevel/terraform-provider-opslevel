package opslevel

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

//go:embed provider_test_block.tf
var providerBlock string

func TestAccDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerBlock + testAccDomainDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opslevel_domain.test", "identifier", "my_domain"),
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
