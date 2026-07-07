//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// TestAccNewRelicFleetConfiguration_Basic covers the create → read → import → destroy lifecycle.
// Also verifies zero-diff on a subsequent plan (no drift after a single apply).
func TestAccNewRelicFleetConfiguration_Basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-config-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.foo"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttrSet(resourceName, "configuration_id"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(resourceName, "configuration_content"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version_entity_id"),
					resource.TestCheckResourceAttr(resourceName, "version_entity_ids.#", "1"),
				),
			},
			// No drift expected after a single apply.
			{
				Config:             testAccFleetConfigBasic(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Import — composite ID "<guid>:<managed_entity_type>" reconstructs all attributes.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFleetConfigImportID(resourceName),
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_ContentUpdate verifies the launch-template-style update path:
// changing configuration_content creates a new immutable API version transparently.
// The resource ID stays constant; total_versions, latest_version_number, latest_version_entity_id,
// and version_entity_ids all advance accordingly.
func TestAccNewRelicFleetConfiguration_ContentUpdate(t *testing.T) {
	rName := fmt.Sprintf("tf-test-update-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.foo"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with initial content (version 1).
			{
				Config: testAccFleetConfigWithContent(rName, "log:\n  level: info\n# v1\n"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "configuration_id"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version_entity_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version_entity_id"),
				),
			},
			// Step 2: Change content — provider creates a new API version transparently.
			{
				Config: testAccFleetConfigWithContent(rName, "log:\n  level: debug\n# v2\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "total_versions", "2"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "version_entity_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version_entity_id"),
				),
			},
			// Step 3: Same content — no new version, no diff.
			{
				Config:             testAccFleetConfigWithContent(rName, "log:\n  level: debug\n# v2\n"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_WithOperatingSystem verifies that a HOST configuration
// can be created and imported with the operating_system attribute set.
func TestAccNewRelicFleetConfiguration_WithOperatingSystem(t *testing.T) {
	rName := fmt.Sprintf("tf-test-os-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.with_os"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetConfigWithOperatingSystem(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttr(resourceName, "operating_system", "LINUX"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "version_entity_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFleetConfigImportID(resourceName),
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_Kubernetes verifies a Kubernetes-targeted configuration.
func TestAccNewRelicFleetConfiguration_Kubernetes(t *testing.T) {
	rName := fmt.Sprintf("tf-test-k8s-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.k8s"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetConfigKubernetes(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "KUBERNETESCLUSTER"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version_entity_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFleetConfigImportID(resourceName),
			},
		},
	})
}

// ── helpers ──────────────────────────────────────────────────────────────────

// testAccFleetConfigImportID returns "<guid>:<managed_entity_type>" as the import ID.
// managed_entity_type is not returned by the GetEntity GraphQL query and must be
// embedded in the import ID.
func testAccFleetConfigImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		managedEntityType := rs.Primary.Attributes["managed_entity_type"]
		return rs.Primary.ID + ":" + managedEntityType, nil
	}
}

func testAccCheckNewRelicFleetConfigurationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no configuration ID is set")
		}

		organizationID := rs.Primary.Attributes["organization_id"]
		if organizationID == "" {
			return fmt.Errorf("no organization_id is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.FleetControl.FleetControlGetConfigurationVersions(
			rs.Primary.ID, organizationID,
		)
		if err != nil {
			return err
		}
		if resp == nil || len(resp.Versions) == 0 {
			return fmt.Errorf("fleet configuration %s has no versions", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicFleetConfigurationDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet_configuration" {
			continue
		}

		organizationID := r.Primary.Attributes["organization_id"]
		if organizationID == "" {
			continue
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		mode := fleetcontrol.GetConfigurationModeTypes.ConfigEntity
		_, err := client.FleetControl.FleetControlGetConfiguration(
			r.Primary.ID, organizationID, mode, 0,
		)
		if err == nil {
			return fmt.Errorf("fleet configuration still exists: %s", r.Primary.ID)
		}
	}
	return nil
}

// ── config templates ─────────────────────────────────────────────────────────

func testAccFleetConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "foo" {
  name                  = %q
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = <<-EOT
    log:
      level: info
    metrics:
      enabled: true
    # v1
  EOT
}
`, name)
}

func testAccFleetConfigWithContent(name, content string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "foo" {
  name                  = %q
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = %q
}
`, name, content)
}

func testAccFleetConfigKubernetes(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "k8s" {
  name                  = %q
  agent_type            = "NRInfra"
  managed_entity_type   = "KUBERNETESCLUSTER"
  configuration_content = <<-EOT
    cluster:
      enabled: true
    prometheus:
      enabled: true
    # v1
  EOT
}
`, name)
}

func testAccFleetConfigWithOperatingSystem(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "with_os" {
  name                  = %q
  agent_type            = "NRInfra"
  managed_entity_type   = "HOST"
  operating_system      = "LINUX"
  configuration_content = <<-EOT
    log:
      level: info
    # os-test-v1
  EOT
}
`, name)
}
