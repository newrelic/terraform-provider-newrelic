//go:build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsStepMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_step_monitor.foo"
	rName := generateNameForIntegrationTestResource()
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
				Config: testAccNewRelicSyntheticsStepMonitorConfig(
					rName,
					"",
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserNewRuntimeTypeVersion,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicSyntheticsStepMonitorConfig(
					fmt.Sprintf("%s-updated", rName),
					updateStep,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserNewRuntimeTypeVersion,
				),
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
					"location_private",
					"tag",
					"enable_screenshot_on_failure_and_script",
					SyntheticsUseLegacyRuntimeAttrLabel,
					"browsers",
					"devices",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsStepMonitorConfig(
	name string,
	step string,
	runtimeType string,
	runtimeTypeVersion string,
) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_step_monitor" "foo" {
	name                                    = "%[1]s"
	period                                  = "EVERY_DAY"
	status                                  = "ENABLED"
	locations_public                        = ["US_WEST_2"]
	enable_screenshot_on_failure_and_script = true
	%[3]s
	%[4]s
	steps {
		ordinal = 0
		type    = "NAVIGATE"
		values  = ["https://google.com"]
	}
	browsers = ["CHROME", "FIREFOX"]
	devices = ["DESKTOP","MOBILE_PORTRAIT", "TABLET_LANDSCAPE", "MOBILE_LANDSCAPE", "TABLET_PORTRAIT"]
	
	%[2]s
}
`,
		name,
		step,
		testConfigurationStringBuilder("runtime_type", runtimeType),
		testConfigurationStringBuilder("runtime_type_version", runtimeTypeVersion),
	)
}
