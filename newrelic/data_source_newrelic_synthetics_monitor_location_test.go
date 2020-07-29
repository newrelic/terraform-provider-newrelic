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

func TestAccNewRelicSyntheticsMonitorLocationDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicSyntheticsLocationDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_monitor.bar"),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsLocationDataSource(n string) resource.TestCheckFunc {
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

func testAccCheckNewRelicSyntheticsLocationDataSourceConfig() string {
	return fmt.Sprintf(`

data "newrelic_synthetics_monitor" "bar" {
	label = Cape Town, ZA
}
`)
}
