//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicSyntheticsMonitorLocationDataSource(t *testing.T) {
	t.Parallel()

	// Temporary until we can provision a private location for our tests
	testMonitorLocationLabel := "oac-integration-test-location"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			// TODO: Create a test private location to fetch with the data source
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigDataSourceSyntheticsLocation("label", testMonitorLocationLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_monitor_location.bar", "label"),
				),
			},
			{
				Config: testConfigDataSourceSyntheticsLocation("name", testMonitorLocationLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_monitor_location.bar", "name"),
				),
			},
		},
	})

	// TODO: Cleanup test private location after test has executed
	// defer func() { // cleanup code }
}

func testAccNewRelicSyntheticsLocationDataSource(n string, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if _, ok := a[attr]; !ok {
			return fmt.Errorf("expected to read synthetics monitor location data from New Relic using attribute `%s`", attr)
		}
		return nil
	}
}

func testConfigDataSourceSyntheticsLocation(attr string, value string) string {
	return fmt.Sprintf(`
data "newrelic_synthetics_monitor_location" "bar" {
	%s = "%s"
}
`, attr, value)
}
