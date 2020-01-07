package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicSyntheticsMonitorScript_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor_script.foo_script"
	rName := acctest.RandString(5)
	scriptText := acctest.RandString(5)
	scriptTextUpdated := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorScriptDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMonitorScriptConfig(rName, scriptText),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorScriptExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMonitorScriptConfig(rName, scriptTextUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorScriptExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsMonitorScriptExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor script ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Synthetics

		script, err := client.GetMonitorScript(rs.Primary.ID)
		if err != nil {
			return err
		}

		if script.Text != rs.Primary.Attributes["text"] {
			return fmt.Errorf("synthetics monitor script text does not match: %v \n\n %v", script.Text, rs.Primary.Attributes["text"])
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorScriptDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Synthetics
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_monitor_script" {
			continue
		}

		monitorScript, err := client.GetMonitorScript(r.Primary.ID)
		if err == nil && monitorScript != nil {
			return fmt.Errorf("synthetics monitor script text still exists")
		}
	}
	return nil
}

func testAccNewRelicSyntheticsMonitorScriptConfig(name string, scriptText string) string {
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
  monitor_id = newrelic_synthetics_monitor.foo.id
  text = "%[2]s"
}
`, name, scriptText)
}
