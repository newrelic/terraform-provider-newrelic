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
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "type", "SIMPLE"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "frequency", "5"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "uri", "http://www.example.com"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "locations.#", "1"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "locations.0", "AWS_US_EAST_1"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "status", "DISABLED"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "sla_threshold", "10"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "validation_string", "Passed"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "verify_ssl", "false"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "bypass_head", "false"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "redirect_is_fail", "false"),
				),
			},
			{
				Config: testAccCheckNewRelicSyntheticsMonitorUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists("newrelic_synthetics_monitor.foo"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "type", "SIMPLE"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "frequency", "10"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "uri", "http://www.example2.com"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "locations.#", "1"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "locations.0", "AWS_US_EAST_2"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "status", "MUTED"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "sla_threshold", "3"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "validation_string", "Ok"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "verify_ssl", "true"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "bypass_head", "true"),
					resource.TestCheckResourceAttr("newrelic_synthetics_monitor.foo", "redirect_is_fail", "true"),
				),
			},
		},
	})
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

func testAccCheckNewRelicSyntheticsMonitorConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name = "tf-test-%s"
	type                = "SIMPLE"
    frequency           = 5
    uri                 = "http://www.example.com"
    locations           = ["AWS_US_EAST_1"]
    status              = "DISABLED"
    sla_threshold       = 10
    validation_string   = "Passed"
    verify_ssl          = "false"
    bypass_head         = "false"
    redirect_is_fail    = "false"
}
`, rName)
}

func testAccCheckNewRelicSyntheticsMonitorUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name = "tf-test-updated-%s"
	type                = "SIMPLE"
    frequency           = 10
    uri                 = "http://www.example2.com"
    locations           = ["AWS_US_EAST_2"]
    status              = "MUTED"
    sla_threshold       = 3
    validation_string   = "Ok"
    verify_ssl          = "true"
    bypass_head         = "true"
    redirect_is_fail    = "true"
}
`, rName)
}
