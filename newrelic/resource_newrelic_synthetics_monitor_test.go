package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicSyntheticsMonitor_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicSyntheticsMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists("newrelic_synthetics_monitor.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "name", rName),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "type", "SIMPLE"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "frequency", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "status", "DISABLED"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "uri", "https://google.com"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "locations.#", "1"),
				),
			},
			{
				Config: testAccCheckNewRelicSyntheticsMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists("newrelic_synthetics_monitor.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "name", fmt.Sprintf("%s-updated", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "type", "SIMPLE"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "frequency", "5"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "status", "ENABLED"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "uri", "https://github.com"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor.foo", "locations.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsMonitorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Synthetics

		found, err := client.GetMonitor(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Synthetics monitor not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Synthetics
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_monitor" {
			continue
		}

		_, err := client.GetMonitor(r.Primary.ID)
		if err == nil {
			return fmt.Errorf("Synthetics monitor still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicSyntheticsMonitorConfig(rName string) string {
	return fmt.Sprintf(`

resource "newrelic_synthetics_monitor" "foo" {
  name = "%[1]s"
  type = "SIMPLE"
  frequency = 1
  status = "DISABLED"
  locations = ["AWS_US_EAST_1"]
  uri = "https://google.com"
}
`, rName)
}

func testAccCheckNewRelicSyntheticsMonitorConfigUpdated(rName string) string {
	return fmt.Sprintf(`

resource "newrelic_synthetics_monitor" "foo" {
  name = "%[1]s-updated"
  type = "SIMPLE"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1", "AWS_US_WEST_1"]
  uri = "https://github.com"
}
`, rName)
}
