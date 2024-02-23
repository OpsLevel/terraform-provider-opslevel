package opslevel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opslevel_domain.test", "name", "one"),
					resource.TestCheckResourceAttr("opslevel_domain.test", "defaulted", "example value when not configured"),
					resource.TestCheckResourceAttr("opslevel_domain.test", "id", "example-id"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "opslevel_domain.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// This is not normally necessary, but is here because this
			// 	// example code does not have an actual upstream service.
			// 	// Once the Read method is able to refresh information from
			// 	// the upstream service, this can be removed.
			// 	ImportStateVerifyIgnore: []string{"name", "defaulted"},
			// },
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opslevel_domain.test", "configurable_attribute", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(name string) string {
	// TODO: fails since "opslevel_team" not defined. Here for illustration purposes
	return fmt.Sprintf(`
data "opslevel_team" "devs" {
  alias = "developers"
}

resource "opslevel_domain" "test" {
  name = %[1]q
  description = "example description"
  owner = data.opslevel_team.devs.id
  note = "example note"
}
`, name)
}
