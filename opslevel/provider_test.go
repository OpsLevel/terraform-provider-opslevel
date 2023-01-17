package opslevel

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/opslevel/opslevel-go/v2022"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func testAccPreCheck(t *testing.T) func() {
	return func() {
		//if v := os.Getenv("OPSLEVEL_API_URL"); v == "" {
		//	t.Fatal("OPSLEVEL_API_URL must be set for acceptance tests")
		//}
		if v := os.Getenv("OPSLEVEL_API_TOKEN"); v == "" {
			t.Fatal("OPSLEVEL_API_TOKEN must be set for acceptance tests")
		}
	}
}

func newTestAcc(t *testing.T, checkDestroy resource.TestCheckFunc, steps ...resource.TestStep) {
	resource.Test(t, resource.TestCase{
		PreCheck:     testAccPreCheck(t),
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps:        steps,
	})
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"opslevel": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

type OpsLevelCheckDestroy func(c *opslevel.Client, id string) (interface{}, error)

func testAccCheckDestroy(typeName string, errRegex string, callback OpsLevelCheckDestroy) func(*terraform.State) error {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*opslevel.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != typeName {
				continue
			}
			id := rs.Primary.ID
			_, err := callback(c, id)
			if err == nil {
				return fmt.Errorf("'%s' - '%s' still exists", typeName, id)
			}
			expectedErr := regexp.MustCompile(errRegex)
			if !expectedErr.Match([]byte(err.Error())) {
				return fmt.Errorf("expected '%s', got '%s'", errRegex, err)
			}
		}

		return nil
	}
}

func readFixture(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func testAccLoadDatasource(name string) string {
	return readFixture("../examples/data-sources/" + name + "/data-source.tf")
}

func testAccLoadResource(name string) string {
	return readFixture("../examples/resources/" + name + "/resource.tf")
}
