//go:build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsStepMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_step_monitor.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	updateStep := `steps {
		ordinal = 1
		type    = "ASSERT_TITLE"
		values  = ["==", "Google"]
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsStepMonitorConfig(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicSyntheticsStepMonitorConfig(fmt.Sprintf("%s-updated", rName), updateStep),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"locations_private",
					"tag",
					"enable_screenshot_on_failure_and_script",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsStepMonitorConfig(name string, step string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_step_monitor" "foo" {
	name                                    = "%[1]s"
	period                                  = "EVERY_DAY"
	status                                  = "ENABLED"
	locations_public                        = ["US_WEST_2"]
	enable_screenshot_on_failure_and_script = true

	steps {
		ordinal = 0
		type    = "NAVIGATE"
		values  = ["https://google.com"]
	}

	%[2]s
}
`, name, step)
}
