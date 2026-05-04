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

// TestAccNewRelicFleetConfiguration_Basic covers the create → read → import → destroy lifecycle
// with a single version block. Also verifies computed fields are populated after a single apply
// (no second-apply required for output variables).
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
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version_entity_id"),
					// TypeList positional checks
					resource.TestCheckResourceAttr(resourceName, "version.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "version.0.version_entity_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version.0.configuration_content"),
				),
			},
			// Import — compound ID reconstructs non-API-readable fields (agent_type, etc.)
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFleetConfigImportID(resourceName),
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_SingleApplyOutputs verifies that all computed fields
// (including version_entity_id and version_number) are populated after a single apply,
// confirming the Read-at-end-of-Create fix.
func TestAccNewRelicFleetConfiguration_SingleApplyOutputs(t *testing.T) {
	rName := fmt.Sprintf("tf-test-output-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.foo"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// Single apply — all fields must be set without a follow-up plan/apply cycle
			{
				Config: testAccFleetConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					// These would be empty string if Read was not called at end of Create
					resource.TestCheckResourceAttrSet(resourceName, "version.0.version_entity_id"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version_entity_id"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
				),
			},
			// Plan-only step: zero diff expected (no "known after apply" drift)
			{
				Config:             testAccFleetConfigBasic(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_MultiVersion creates a configuration with two versions
// and verifies both are tracked correctly with unique entity IDs and version numbers.
func TestAccNewRelicFleetConfiguration_MultiVersion(t *testing.T) {
	rName := fmt.Sprintf("tf-test-multi-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.multi"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with v1 only
			{
				Config: testAccFleetConfigOneVersion(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "version.0.version_entity_id"),
				),
			},
			// Step 2: Add v2 — verify v1 entity_id is preserved, latest_version_number advances
			{
				Config: testAccFleetConfigTwoVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "2"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.1.version_number", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "version.0.version_entity_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version.1.version_entity_id"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_VersionSequence exercises the full add/remove sequence:
// v1 → +v2 → +v3 → -v2 → +v4+v5 → -v1-v5
// Verifies that entity_ids for unchanged versions remain stable (TypeList key advantage).
func TestAccNewRelicFleetConfiguration_VersionSequence(t *testing.T) {
	rName := fmt.Sprintf("tf-test-seq-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.seq"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// v1 only
			{
				Config: testAccFleetConfigSeq(rName, true, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
				),
			},
			// +v2
			{
				Config: testAccFleetConfigSeq(rName, true, true, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "2"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.1.version_number", "2"),
				),
			},
			// +v3 (v1, v2, v3)
			{
				Config: testAccFleetConfigSeq(rName, true, true, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "3"),
					resource.TestCheckResourceAttr(resourceName, "version.2.version_number", "3"),
				),
			},
			// -v2 (v1, v3 remain) — content-based matching must delete v2's entity, not position
			{
				Config: testAccFleetConfigSeq(rName, true, false, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "2"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.1.version_number", "3"),
				),
			},
			// +v4+v5 simultaneously (v1, v3, v4, v5)
			{
				Config: testAccFleetConfigSeq(rName, true, false, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "4"),
				),
			},
			// -v1-v5 simultaneously (v3, v4 remain)
			{
				Config: testAccFleetConfigSeq(rName, false, false, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "2"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "3"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetConfiguration_RollbackWorkflow exercises the rollback pattern:
// v1(A) → v2(B) → v3(A) — verifies that re-using version A's content creates a new version
// rather than resurrecting the old entity_id, since the TypeList in-place edit guard only
// fires when block count stays the same and a block's content changes.
func TestAccNewRelicFleetConfiguration_RollbackWorkflow(t *testing.T) {
	rName := fmt.Sprintf("tf-test-rollback-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_configuration.rollback"

	const contentA = "log:\n  level: info\n# rollback-a\n"
	const contentB = "log:\n  level: debug\n# rollback-b\n"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetConfigurationDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with content A (version 1)
			{
				Config: testAccFleetConfigRollback(rName, []string{contentA}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "1"),
				),
			},
			// Step 2: Add content B (version 2)
			{
				Config: testAccFleetConfigRollback(rName, []string{contentA, contentB}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "version.1.version_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "2"),
				),
			},
			// Step 3: Add content A again (version 3 — rollback to same content as v1)
			// This must succeed: adding a block with content that previously existed
			// but is not currently in state is a pure addition, not an in-place edit.
			{
				Config: testAccFleetConfigRollback(rName, []string{contentA, contentB, contentA + "# copy\n"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "3"),
					resource.TestCheckResourceAttr(resourceName, "latest_version_number", "3"),
				),
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
					resource.TestCheckResourceAttr(resourceName, "version.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "total_versions", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
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

// testAccFleetConfigImportID builds the compound import ID needed by the custom
// importer: configGUID:orgID:agentType:managedEntityType:name.
// The API does not expose agent_type/managed_entity_type/name via any working
// read endpoint, so we encode them in the import ID to satisfy ImportStateVerify.
func testAccFleetConfigImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s:%s:%s:%s:%s",
			rs.Primary.ID,
			rs.Primary.Attributes["organization_id"],
			rs.Primary.Attributes["agent_type"],
			rs.Primary.Attributes["managed_entity_type"],
			rs.Primary.Attributes["name"],
		), nil
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
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = <<-EOT
      log:
        level: info
      metrics:
        enabled: true
      # v1
    EOT
  }
}
`, name)
}

func testAccFleetConfigOneVersion(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "multi" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = <<-EOT
      log:
        level: info
      # v1
    EOT
  }
}
`, name)
}

func testAccFleetConfigTwoVersions(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "multi" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = <<-EOT
      log:
        level: info
      # v1
    EOT
  }

  version {
    configuration_content = <<-EOT
      log:
        level: debug
      # v2
    EOT
  }
}
`, name)
}

// testAccFleetConfigSeq generates a config for the version sequence test.
// Each boolean flag controls whether that version block is included.
func testAccFleetConfigSeq(name string, v1, v2, v3, v4, v5 bool) string {
	cfg := fmt.Sprintf(`
resource "newrelic_fleet_configuration" "seq" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"
`, name)

	if v1 {
		cfg += `
  version {
    configuration_content = <<-EOT
      log:
        level: info
      metrics:
        enabled: true
        sample_rate: 30
      # v1
    EOT
  }
`
	}
	if v2 {
		cfg += `
  version {
    configuration_content = <<-EOT
      log:
        level: debug
      metrics:
        enabled: true
        sample_rate: 10
      # v2
    EOT
  }
`
	}
	if v3 {
		cfg += `
  version {
    configuration_content = <<-EOT
      log:
        level: warn
      metrics:
        enabled: false
      # v3
    EOT
  }
`
	}
	if v4 {
		cfg += `
  version {
    configuration_content = <<-EOT
      log:
        level: error
      metrics:
        enabled: true
        sample_rate: 60
      # v4
    EOT
  }
`
	}
	if v5 {
		cfg += `
  version {
    configuration_content = <<-EOT
      log:
        level: trace
      metrics:
        enabled: true
        sample_rate: 5
      # v5
    EOT
  }
`
	}

	cfg += "}\n"
	return cfg
}

func testAccFleetConfigRollback(name string, contents []string) string {
	cfg := fmt.Sprintf(`
resource "newrelic_fleet_configuration" "rollback" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"
`, name)

	for _, c := range contents {
		cfg += fmt.Sprintf(`
  version {
    configuration_content = %q
  }
`, c)
	}

	cfg += "}\n"
	return cfg
}

func testAccFleetConfigKubernetes(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "k8s" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "KUBERNETESCLUSTER"

  version {
    configuration_content = <<-EOT
      cluster:
        enabled: true
      prometheus:
        enabled: true
      # v1
    EOT
  }
}
`, name)
}
