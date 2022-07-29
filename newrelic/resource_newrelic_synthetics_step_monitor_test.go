//go:build integration

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsStepMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_step_monitor.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsStepMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Update
			{
				PreConfig: func() {
					// Unfortunately we still have to wait due to async delay with entity indexing :(
					time.Sleep(10 * time.Second)
				},
				Config: testAccNewRelicSyntheticsStepMonitorConfig(fmt.Sprintf("%s-updated", rName)),
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
					// NOTE: Need to ignore steps and enable_screenshot_on_failure_and_script
					// on import until the new endpoints work.
					"steps",
					"enable_screenshot_on_failure_and_script",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsStepMonitorConfig(name string) string {
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
}
`, name)
}
