//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func TestAccNewRelicFleetConfiguration_Basic(t *testing.T) {
	resourceName := "newrelic_fleet_configuration.foo"
	rName := fmt.Sprintf("tf-test-config-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicFleetConfigurationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "INFRASTRUCTURE"),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttrSet(resourceName, "configuration_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
				),
			},
			// Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configuration_content", "configuration_file_path"},
			},
		},
	})
}

func TestAccNewRelicFleetConfiguration_Kubernetes(t *testing.T) {
	resourceName := "newrelic_fleet_configuration.k8s"
	rName := fmt.Sprintf("tf-test-k8s-config-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigurationConfigKubernetes(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "KUBERNETES"),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "KUBERNETESCLUSTER"),
				),
			},
			// Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configuration_content", "configuration_file_path"},
			},
		},
	})
}

// Error case tests

func TestAccNewRelicFleetConfiguration_MissingConfiguration(t *testing.T) {
	rName := fmt.Sprintf("tf-test-error-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicFleetConfigurationConfigMissingConfiguration(rName),
				ExpectError: regexp.MustCompile("one of configuration_file_path or configuration_content must be provided"),
			},
		},
	})
}

func TestAccNewRelicFleetConfiguration_BothConfigurations(t *testing.T) {
	rName := fmt.Sprintf("tf-test-error-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicFleetConfigurationConfigBothConfigurations(rName),
				ExpectError: regexp.MustCompile("conflicts with configuration_content|conflicts with configuration_file_path"),
			},
		},
	})
}

// Helper functions

func testAccCheckNewRelicFleetConfigurationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no configuration ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// Get organization ID from resource state
		organizationID := rs.Primary.Attributes["organization_id"]
		if organizationID == "" {
			return fmt.Errorf("no organization_id is set")
		}

		// Try to fetch the configuration
		mode := fleetcontrol.GetConfigurationModeTypes.ConfigEntity
		_, err := client.FleetControl.FleetControlGetConfiguration(
			rs.Primary.ID,
			organizationID,
			mode,
			0,
		)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckNewRelicFleetConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet_configuration" {
			continue
		}

		// Get organization ID from resource state
		organizationID := r.Primary.Attributes["organization_id"]
		if organizationID == "" {
			continue
		}

		// Try to get the configuration - should return error if not found
		mode := fleetcontrol.GetConfigurationModeTypes.ConfigEntity
		_, err := client.FleetControl.FleetControlGetConfiguration(
			r.Primary.ID,
			organizationID,
			mode,
			0,
		)
		if err == nil {
			return fmt.Errorf("fleet configuration still exists: %s", r.Primary.ID)
		}
	}

	return nil
}

// Config functions

func testAccNewRelicFleetConfigurationConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "foo" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
    metrics:
      enabled: true
  EOT
}
`, name)
}

func testAccNewRelicFleetConfigurationConfigKubernetes(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "k8s" {
  name                = "%s"
  agent_type          = "KUBERNETES"
  managed_entity_type = "KUBERNETESCLUSTER"

  configuration_content = <<-EOT
    cluster:
      enabled: true
    prometheus:
      enabled: true
  EOT
}
`, name)
}

func testAccNewRelicFleetConfigurationConfigMissingConfiguration(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "error" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"
}
`, name)
}

func testAccNewRelicFleetConfigurationConfigBothConfigurations(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "error" {
  name                = "%s"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_file_path = "/tmp/config.yml"
  configuration_content   = "log: info"
}
`, name)
}
