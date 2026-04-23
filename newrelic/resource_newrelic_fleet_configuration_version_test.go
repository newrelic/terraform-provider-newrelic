//go:build integration || FLEET

package newrelic

// NOTE: Fleet Configuration API Access
//
// Tests that create actual fleet configurations and versions require access to the
// Fleet Configuration API (blob-api.service.newrelic.com). This API requires special
// enablement on the test account.
//
// If tests fail with "POST https://blob-api.service.newrelic.com/.../AgentConfigurations giving up",
// it indicates the test account does not have Configuration API access enabled.
//
// Validation tests (TestAccNewRelicFleetConfigurationVersion_MissingConfiguration,
// TestAccNewRelicFleetConfigurationVersion_BothConfigurations) should pass as they only
// validate schema-level requirements without making API calls.

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// TestAccNewRelicFleetConfiguration_WithVersions tests the full lifecycle:
// 1. Create config (v1)
// 2. Add version 2
// 3. Add version 3
// 4. Verify all versions exist
// 5. Delete versions and config
func TestAccNewRelicFleetConfiguration_WithVersions(t *testing.T) {
	configResourceName := "newrelic_fleet_configuration.test"
	version2ResourceName := "newrelic_fleet_configuration_version.v2"
	version3ResourceName := "newrelic_fleet_configuration_version.v3"
	rName := fmt.Sprintf("tf-test-versioned-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create config with v1
			{
				Config: testAccNewRelicFleetConfigurationWithVersionsStep1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(configResourceName),
					resource.TestCheckResourceAttr(configResourceName, "name", rName),
					resource.TestCheckResourceAttr(configResourceName, "version", "1"),
				),
			},
			// Step 2: Add version 2
			{
				Config: testAccNewRelicFleetConfigurationWithVersionsStep2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(configResourceName),
					testAccCheckNewRelicFleetConfigurationVersionExists(version2ResourceName),
					resource.TestCheckResourceAttr(configResourceName, "version", "1"),
					resource.TestCheckResourceAttr(version2ResourceName, "version_number", "2"),
				),
			},
			// Step 3: Add version 3
			{
				Config: testAccNewRelicFleetConfigurationWithVersionsStep3(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(configResourceName),
					testAccCheckNewRelicFleetConfigurationVersionExists(version2ResourceName),
					testAccCheckNewRelicFleetConfigurationVersionExists(version3ResourceName),
					resource.TestCheckResourceAttr(version3ResourceName, "version_number", "3"),
				),
			},
		},
	})
}

func TestAccNewRelicFleetConfigurationVersion_MissingConfiguration(t *testing.T) {
	rName := fmt.Sprintf("tf-test-error-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicFleetConfigurationVersionMissingConfig(rName),
				ExpectError: regexp.MustCompile("one of.*configuration_content.*configuration_file_path.*must be specified"),
			},
		},
	})
}

func TestAccNewRelicFleetConfigurationVersion_BothConfigurations(t *testing.T) {
	rName := fmt.Sprintf("tf-test-error-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicFleetConfigurationVersionBothConfigs(rName),
				ExpectError: regexp.MustCompile("only one of.*configuration_content.*configuration_file_path.*can be specified"),
			},
		},
	})
}

// Helper functions

func testAccCheckNewRelicFleetConfigurationVersionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no configuration version ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// Get organization ID from resource state
		organizationID := rs.Primary.Attributes["organization_id"]
		if organizationID == "" {
			return fmt.Errorf("no organization_id is set")
		}

		// Try to fetch the version
		mode := fleetcontrol.GetConfigurationModeTypes.ConfigVersionEntity
		versionNumber := 0 // Get latest if not specified
		if v := rs.Primary.Attributes["version_number"]; v != "" {
			fmt.Sscanf(v, "%d", &versionNumber)
		}

		_, err := client.FleetControl.FleetControlGetConfiguration(
			rs.Primary.ID,
			organizationID,
			mode,
			versionNumber,
		)
		if err != nil {
			return err
		}

		return nil
	}
}

// Config functions - Progressive steps showing version addition

func testAccNewRelicFleetConfigurationWithVersionsStep1(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "test" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
    version: 1
  EOT
}
`, name)
}

func testAccNewRelicFleetConfigurationWithVersionsStep2(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "test" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
    version: 1
  EOT
}

resource "newrelic_fleet_configuration_version" "v2" {
  configuration_id = newrelic_fleet_configuration.test.configuration_id

  configuration_content = <<-EOT
    log:
      level: debug
    version: 2
  EOT
}
`, name)
}

func testAccNewRelicFleetConfigurationWithVersionsStep3(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "test" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
    version: 1
  EOT
}

resource "newrelic_fleet_configuration_version" "v2" {
  configuration_id = newrelic_fleet_configuration.test.configuration_id

  configuration_content = <<-EOT
    log:
      level: debug
    version: 2
  EOT
}

resource "newrelic_fleet_configuration_version" "v3" {
  configuration_id = newrelic_fleet_configuration.test.configuration_id

  configuration_content = <<-EOT
    log:
      level: warn
    metrics:
      enabled: true
    version: 3
  EOT
}
`, name)
}

func testAccNewRelicFleetConfigurationVersionMissingConfig(name string) string {
	// Use a fake configuration_id to avoid creating actual resources
	// This test only validates schema-level requirements
	return `
resource "newrelic_fleet_configuration_version" "error" {
  configuration_id = "fake-config-id-for-validation"
  # Missing both configuration_file_path and configuration_content
}
`
}

func testAccNewRelicFleetConfigurationVersionBothConfigs(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "test" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = "log: info"
}

resource "newrelic_fleet_configuration_version" "error" {
  configuration_id = newrelic_fleet_configuration.test.configuration_id

  configuration_file_path = "/tmp/config.yml"
  configuration_content   = "log: debug"
}
`, name)
}
