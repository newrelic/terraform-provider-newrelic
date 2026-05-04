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

// TestAccNewRelicFleetDeployment_Basic covers create → read → import → destroy
// for a minimal deployment (no configuration versions).
func TestAccNewRelicFleetDeployment_Basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.basic"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test deployment"),
					resource.TestCheckResourceAttr(resourceName, "agent.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.version", "1.58.0"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_id"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fleet_id"),
				),
			},
			// Import by deployment GUID
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_Update verifies that name, description, and
// agent version can be updated in place.
func TestAccNewRelicFleetDeployment_Update(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.basic"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "agent.0.version", "1.58.0"),
				),
			},
			{
				Config: testAccFleetDeploymentUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-updated", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "agent.0.version", "1.59.0"),
				),
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
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentMultipleAgents(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent.#", "2"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetDeployment_WithTags verifies that tags are created and
// surfaced in state.
func TestAccNewRelicFleetDeployment_WithTags(t *testing.T) {
	rName := fmt.Sprintf("tf-test-deploy-tags-%s", acctest.RandString(5))
	resourceName := "newrelic_fleet_deployment.tagged"

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetDeploymentWithTags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetDeploymentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
		},
	})
}

// Helper functions

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

// Config functions

// testAccFleetDeploymentFleetID returns the fleet GUID to use in deployment tests.
// Replace or populate via environment variable in your test environment.
const testAccFleetDeploymentFleetID = "REPLACE_WITH_FLEET_GUID"

func testAccFleetDeploymentBasic(name string) string {
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
`, testAccFleetDeploymentFleetID, name)
}

func testAccFleetDeploymentUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "basic" {
  fleet_id    = %q
  name        = "%s-updated"
  description = "Updated description"

  agent {
    agent_type = "NRInfra"
    version    = "1.59.0"
  }
}
`, testAccFleetDeploymentFleetID, name)
}

func testAccFleetDeploymentMultipleAgents(name string) string {
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
`, testAccFleetDeploymentFleetID, name)
}

func testAccFleetDeploymentWithTags(name string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_deployment" "tagged" {
  fleet_id    = %q
  name        = %q
  description = "Tagged deployment"

  agent {
    agent_type = "NRInfra"
    version    = "1.58.0"
  }

  tags = ["environment:production", "team:platform"]
}
`, testAccFleetDeploymentFleetID, name)
}
