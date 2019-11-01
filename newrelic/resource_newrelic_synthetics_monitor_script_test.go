package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicSyntheticsMonitorScript_Basic(t *testing.T) {
	rname := acctest.RandString(5)
	scriptText := acctest.RandString(5)
	scriptTextUpdated := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicSyntheticsMonitorScriptConfig(rname, scriptText),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorScriptExists("newrelic_synthetics_monitor_script.foo_script"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor_script.foo_script", "text", scriptText),
				),
			},
			{
				Config: testAccCheckNewRelicSyntheticsMonitorScriptConfig(rname, scriptTextUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorScriptExists("newrelic_synthetics_monitor_script.foo_script"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_monitor_script.foo_script", "text", scriptTextUpdated),
				),
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsMonitorScriptExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No synthetics monitor script ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Synthetics

		foundText, err := client.GetMonitorScript(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundText != rs.Primary.Attributes["text"] {
			return fmt.Errorf("Synthetics monitor script text does not match: %v \n\n %v", foundText, rs.Primary.Attributes["text"])
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorScriptDestroy(s *terraform.State) error {
	defer testAccCheckDestroy()

	client := testAccProvider.Meta().(*ProviderConfig).Synthetics
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_monitor_script" {
			continue
		}

		foundText, err := client.GetMonitorScript(r.Primary.ID)
		if err == nil || strings.TrimSpace(foundText) != "" {
			return fmt.Errorf("Synthetics monitor script text still exists")
		}
	}
	return nil
}

func testAccCheckNewRelicSyntheticsMonitorScriptConfig(rName string, scriptText string) string {
	return fmt.Sprintf(`

resource "newrelic_synthetics_monitor" "foo" {
  name = "%[1]s"
  type = "SCRIPT_BROWSER"
  frequency = 1
  status = "DISABLED"
  locations = ["AWS_US_EAST_1"]
  uri = "https://google.com"
}

resource "newrelic_synthetics_monitor_script" "foo_script" {
  monitor_id = "${newrelic_synthetics_monitor.foo.id}"
  text = "%[2]s"
}
`, rName, scriptText)
}
