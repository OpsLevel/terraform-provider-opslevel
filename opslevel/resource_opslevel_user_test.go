package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hasura/go-graphql-client"
	"github.com/opslevel/opslevel-go/v2023"
	"testing"
)

func testAccCheckUser(c *opslevel.Client, id string) (interface{}, error) {
	return c.GetUser(graphql.ID(id))
}

func TestAccOpsLevelUser(t *testing.T) {
	newTestAcc(t, testAccCheckDestroy("opslevel_user", "User with id '.*' does not exist on this account", testAccCheckUser),
		resource.TestStep{
			Config: testAccLoadResource("opslevel_user"),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("opslevel_user.john", "name", "John Doe"),
				resource.TestCheckResourceAttr("opslevel_user.john", "email", "john.doe@example.com"),
				resource.TestCheckResourceAttr("opslevel_user.john", "role", "user"),
			),
		})
}
