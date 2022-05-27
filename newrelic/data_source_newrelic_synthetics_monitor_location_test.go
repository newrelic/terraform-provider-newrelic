//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testMonitorLocationLabel = "oac-integration-test-location"
)

func TestAccNewRelicSyntheticsMonitorLocationDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicSyntheticsLocationDataSourceConfig(testMonitorLocationLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_monitor_location.bar"),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsLocationDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["label"] == "" {
			return fmt.Errorf("expected to read synthetics monitor location data from New Relic")
		}
		return nil
	}
}

func testAccCheckNewRelicSyntheticsLocationDataSourceConfig(label string) string {
	return fmt.Sprintf(`

data "newrelic_synthetics_monitor_location" "bar" {
	label = "%s"
}
`, label)
}
