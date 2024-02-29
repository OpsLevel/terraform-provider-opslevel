package opslevel

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: commonProviderBlock + testAccDomainDataSourceConfig,
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
