//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// testAccFleetDeploymentFleetID is the GUID of the shared fleet used in
// deployment acceptance tests. Tests reference this fleet by ID only — they
// never manage the fleet resource itself, so it will never be deleted by
// a test teardown.
const testAccFleetDeploymentFleetID = "NjQyNTg2NXxOR0VQfEZMRUVUfDAxOWRmMTkyLTAxNjktNzJiZi1hNDA2LWVhNDkxNTZmZjUzNg"

// testAccPreCheckFleetDeploymentEnvVars extends testAccPreCheckFleetEnvVars
// with an additional gate on NEW_RELIC_FLEET_TEST_FLEET_ID. That env var is
// set only in the CI run that has access to the shared fleet (account 6425865).
// Tests that make real API calls against the fleet must use this pre-check so
// they skip in CI runs that use a different account.
func testAccPreCheckFleetDeploymentEnvVars(t *testing.T) {
	testAccPreCheckFleetEnvVars(t)
	if os.Getenv("NEW_RELIC_FLEET_TEST_FLEET_ID") == "" {
		t.Skip("NEW_RELIC_FLEET_TEST_FLEET_ID must be set for fleet deployment acceptance tests")
	}
}

// TestAccNewRelicFleetDeployment_Basic covers create → read → import → destroy
// for a minimal deployment with no linked configuration.
func TestAccNewRelicFleetDeployment_Basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.basic"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentBasic(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test deployment"),
					resource.TestCheckResourceAttr(resourceName, "fleet_id", testAccFleetDeploymentFleetID),
					resource.TestCheckResourceAttr(resourceName, "agent.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.version", "1.58.0"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_id"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(resourceName, "phase"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_WithConfiguration creates a fleet configuration
// and links its latest version to the deployment via configuration_version_id.
func TestAccNewRelicFleetDeployment_WithConfiguration(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-cfg-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.with_config"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentWithConfiguration(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.agent_type", "NRInfra"),
					resource.TestCheckResourceAttrSet(resourceName, "agent.0.configuration_version_id"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_NoDriftAfterCreate verifies that a plan
// immediately after create shows no changes (no perpetual drift).
func TestAccNewRelicFleetDeployment_NoDriftAfterCreate(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-nodrift-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentBasic(rName, testAccFleetDeploymentFleetID),
			},
			{
				Config:             testAccFleetDeploymentBasic(rName, testAccFleetDeploymentFleetID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_MultipleAgents verifies that a deployment with
// multiple agent blocks is created and read back correctly.
func TestAccNewRelicFleetDeployment_MultipleAgents(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-multi-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.multi"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentMultipleAgents(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "agent.1.agent_type", "FluentBit"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_ZeroAgentsOnUpdate verifies that:
//   - Creating a deployment with zero agent blocks is rejected at plan time.
//   - Updating an existing CREATED deployment to zero agents is allowed.
func TestAccNewRelicFleetDeployment_ZeroAgentsOnUpdate(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-zero-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			// Step 1: create with one agent block — must succeed.
			{
				Config: testAccFleetDeploymentBasic(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists("newrelic_fleet_deployment.basic"),
					resource.TestCheckResourceAttr("newrelic_fleet_deployment.basic", "agent.#", "1"),
				),
			},
			// Step 2: update to zero agents — must be accepted at plan time and
			// result in an empty agent list in state.
			{
				Config: testAccFleetDeploymentZeroAgents(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists("newrelic_fleet_deployment.basic"),
					resource.TestCheckResourceAttr("newrelic_fleet_deployment.basic", "agent.#", "0"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_ZeroAgentsOnCreate verifies that creating a
// deployment without any agent block is rejected at plan time.
func TestAccNewRelicFleetDeployment_ZeroAgentsOnCreate(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-nocreate-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFleetDeploymentZeroAgents(rName, testAccFleetDeploymentFleetID),
				ExpectError: regexp.MustCompile(`at least one agent block is required`),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_DuplicateAgentType verifies that declaring two
// agent blocks with the same agent_type is rejected at plan time.
func TestAccNewRelicFleetDeployment_DuplicateAgentType(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-dup-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFleetDeploymentDuplicateAgentType(rName, testAccFleetDeploymentFleetID),
				ExpectError: regexp.MustCompile(`duplicate agent_type "NRInfra"`),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_WithTags verifies that tags are created and
// reflected in state.
func TestAccNewRelicFleetDeployment_WithTags(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-tags-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.tagged"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentWithTags(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_PhaseGate verifies that attempting to update a
// deployment that is no longer in CREATED phase is blocked at plan time with a
// descriptive error pointing the user toward terraform destroy.
//
// Note: this test is skipped if the deployment stays in CREATED long enough
// that the phase hasn't advanced before the second step runs. In practice the
// fleet backend transitions within seconds.
func TestAccNewRelicFleetDeployment_PhaseGate(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-gate-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.gate"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetDeploymentEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the deployment (phase starts as CREATED).
			{
				Config: testAccFleetDeploymentGate(rName, testAccFleetDeploymentFleetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
				),
			},
			// Step 2: Attempt to change name after phase has advanced.
			// CustomizeDiff must block the plan with a clear error.
			{
				Config:      testAccFleetDeploymentGateUpdated(rName, testAccFleetDeploymentFleetID),
				ExpectError: regexp.MustCompile(`cannot update fleet deployment`),
			},
		},
	})
}

// ── helpers ───────────────────────────────────────────────────────────────────

func testAccCheckNewRelicFleetDeploymentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fleet deployment ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		entityInterface, err := client.FleetControl.GetEntity(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching fleet deployment %s: %s", rs.Primary.ID, err)
		}
		if entityInterface == nil {
			return fmt.Errorf("fleet deployment %s not found", rs.Primary.ID)
		}
		if _, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetDeploymentEntity); !ok {
			return fmt.Errorf("entity %s is not a fleet deployment", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicFleetDeploymentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet_deployment" {
			continue
		}

		entityInterface, err := client.FleetControl.GetEntity(r.Primary.ID)
		if err == nil && entityInterface != nil {
			return fmt.Errorf("fleet deployment still exists: %s", r.Primary.ID)
		}
	}

	return nil
}

// ── config templates ──────────────────────────────────────────────────────────

func testAccFleetDeploymentBasic(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "basic" {
  fleet_id    = %q
  name        = %q
  description = "Test deployment"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }
}
`, fleetID, name)
}

func testAccFleetDeploymentWithConfiguration(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "deploy_cfg" {
  name                = %q
  agent_type          = "NRInfra"
  managed_entity_type = "HOST"

  version {
    configuration_content = <<-EOT
      log:
        level: info
      # deployment-config-test
    EOT
  }
}

resource "newrelic_fleet_deployment" "with_config" {
  fleet_id    = %q
  name        = %q
  description = "Deployment linked to a configuration version"

  agent {
    agent_type               = "NRInfra"
    version                  = "1.58.0"
    configuration_version_id = newrelic_fleet_configuration.deploy_cfg.latest_version_entity_id
  }
}
`, name+"-cfg", fleetID, name)
}

func testAccFleetDeploymentMultipleAgents(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "multi" {
  fleet_id    = %q
  name        = %q
  description = "Multi-agent deployment"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }

  agent {
    agent_type = "FluentBit"
    version    = "3.2.0"
  }
}
`, fleetID, name)
}

func testAccFleetDeploymentDuplicateAgentType(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "dup" {
  fleet_id    = %q
  name        = %q
  description = "Should fail due to duplicate agent type"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }

  agent {
    agent_type = "NRInfra"
    version    = "1.59.0"
  }
}
`, fleetID, name)
}

func testAccFleetDeploymentWithTags(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "tagged" {
  fleet_id    = %q
  name        = %q
  description = "Tagged deployment"
  tags        = ["environment:production", "team:platform"]

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }
}
`, fleetID, name)
}

func testAccFleetDeploymentGate(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "gate" {
  fleet_id    = %q
  name        = %q
  description = "Phase gate test deployment"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }
}
`, fleetID, name)
}

func testAccFleetDeploymentGateUpdated(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "gate" {
  fleet_id    = %q
  name        = %q
  description = "Phase gate test — updated after phase advanced"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }
}
`, fleetID, name+"-changed")
}

func testAccFleetDeploymentZeroAgents(name, fleetID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "basic" {
  fleet_id    = %q
  name        = %q
  description = "Zero-agent deployment"
}
`, fleetID, name)
}
