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

func TestAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributes_EmptyValuesError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
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

func TestAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributes_LegacyValuesError(t *testing.T) {
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
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
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

func TestAccNewRelicSyntheticsScriptedBrowserMonitor_CreateWithLegacyRuntimeAttributes_EmptyValuesError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					"SCRIPT_BROWSER",
				),
				ExpectError: regexp.MustCompile(
					constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(RUNTIME_TYPE_ATTRIBUTE_LABEL).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptedBrowserMonitor_CreateWithLegacyRuntimeAttributes_LegacyValuesError(t *testing.T) {
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
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					runtimeTypeInConfig,
					runtimeTypeVersionInConfig,
					"SCRIPT_BROWSER",
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

func TestAccNewRelicSyntheticsScriptedAPIMonitor_CreateWithLegacyRuntimeAttributes_EmptyValuesError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					"SCRIPT_API",
				),
				ExpectError: regexp.MustCompile(
					constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(RUNTIME_TYPE_ATTRIBUTE_LABEL).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptedAPIMonitor_CreateWithLegacyRuntimeAttributes_LegacyValuesError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	runtimeTypeInConfig := "NODE_API"
	runtimeTypeVersionInConfig := "10"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					runtimeTypeInConfig,
					runtimeTypeVersionInConfig,
					"SCRIPT_API",
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

func testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "legacy_synthetics_monitor_runtime_attributes_empty_strings" {
			name                                    = "%[1]s"
			period                                  = "EVERY_DAY"
			type									= "BROWSER"
			status                                  = "ENABLED"
			locations_public                        = ["US_WEST_2"]
			enable_screenshot_on_failure_and_script = true
			runtime_type							= "%[2]s"
			runtime_type_version					= "%[3]s"
		}
`,
		name,
		runtimeType,
		runtimeTypeVersion,
	)
}

func testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
	scriptType string,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "legacy_synthetics_script_monitor_runtime_attributes_empty_strings" {
			name                                    = "%[1]s"
			period                                  = "EVERY_DAY"
			type									= "%[4]s"
			status                                  = "ENABLED"
			locations_public                        = ["US_WEST_2"]
			enable_screenshot_on_failure_and_script = true
			script									= "console.log('');"
			runtime_type							= "%[2]s"
			runtime_type_version					= "%[3]s"
		}
`,
		name,
		runtimeType,
		runtimeTypeVersion,
		scriptType,
	)
}
