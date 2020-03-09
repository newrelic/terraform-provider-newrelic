package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicApplicationData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicApplicationDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicApplicationExists("data.newrelic_application.app"),
				),
			},
		},
	})
}

func testAccCheckNewRelicApplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get an application from New Relic")
		}

		if a["name"] != testAccExpectedApplicationName {
			return fmt.Errorf("expected the application name to be: %s, but got: %s", testAccExpectedApplicationName, a["name"])
		}

		return nil
	}
}

// The test application for this data source is created in provider_test.go
func testAccNewRelicApplicationDataConfig() string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%s"
}
`, testAccExpectedApplicationName)
}
