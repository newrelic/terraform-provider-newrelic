// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	expectedMonitorName = fmt.Sprintf("tf-test-synthetic-%s", acctest.RandString(5))
)

func TestAccNewRelicSyntheticsMonitorDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicSyntheticsDataSourceConfig(expectedMonitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsDataSource("data.newrelic_synthetics_monitor.bar"),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to read synthetics monitor data from New Relic")
		}

		if a["name"] != expectedMonitorName {
			return fmt.Errorf("expected the synthetics monitor name to be: %s, but got: %s", expectedMonitorName, a["name"])
		}
		return nil
	}
}

func testAccCheckNewRelicSyntheticsDataSourceConfig(name string) string {
	return fmt.Sprintf(`

resource "newrelic_synthetics_monitor" "foo" {
	name = "%[1]s"
	type = "SIMPLE"
	frequency = 15
	status = "DISABLED"
	locations = ["AWS_US_EAST_1"]
	uri = "https://google.com"
}

data "newrelic_synthetics_monitor" "bar" {
	name = newrelic_synthetics_monitor.foo.name
}
`, name)
}
