//go:build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributes_EmptyValuesError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
				),
				ExpectError: regexp.MustCompile(
					constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(RUNTIME_TYPE_ATTRIBUTE_LABEL).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributes_LegacyValuesError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	runtimeTypeInConfig := "CHROME_BROWSER"
	runtimeTypeVersionInConfig := "72"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					runtimeTypeInConfig,
					runtimeTypeVersionInConfig,
				),
				ExpectError: regexp.MustCompile(
					constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
						RUNTIME_TYPE_ATTRIBUTE_LABEL,
						RUNTIME_TYPE_VERSION_ATTRIBUTE_LABEL,
						runtimeTypeInConfig,
						runtimeTypeVersionInConfig,
					).Error(),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_step_monitor" "legacy_synthetics_step_monitor_runtime_attributes_empty_strings" {
			name                                    = "%[1]s"
			period                                  = "EVERY_DAY"
			status                                  = "ENABLED"
			locations_public                        = ["US_WEST_2"]
			enable_screenshot_on_failure_and_script = true
			runtime_type							= "%[2]s"
			runtime_type_version					= "%[3]s"
			steps {
				ordinal = 0
				type    = "NAVIGATE"
				values  = ["https://google.com"]
			}
		}
`,
		name,
		runtimeType,
		runtimeTypeVersion,
	)
}
