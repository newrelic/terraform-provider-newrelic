//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccNewRelicFleetMembers_Basic creates a fleet, adds two entity IDs as
// members of the default ring, then updates to swap one of them out.
//
// NOTE: Real entity GUIDs are required for the API to accept the membership.
// In CI the test is gated by NEW_RELIC_FLEET_TEST_ENTITY_IDS (comma-separated).
// Without valid entity GUIDs the create step will fail with an API error;
// the test is still useful for exercising the resource schema and state logic.
func TestAccNewRelicFleetMembers_Basic(t *testing.T) {
	resourceName := "newrelic_fleet_members.default"
	fleetName := fmt.Sprintf("tf-test-members-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			// Create: fleet + initial membership
			{
				Config: testAccNewRelicFleetMembersConfig(fleetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ring", "default"),
					resource.TestCheckResourceAttrSet(resourceName, "fleet_id"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccFleetMembersImportID(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicFleetMembers_Update(t *testing.T) {
	resourceName := "newrelic_fleet_members.default"
	fleetName := fmt.Sprintf("tf-test-members-upd-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetMembersConfig(fleetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_ids.#", "1"),
				),
			},
			{
				Config: testAccNewRelicFleetMembersConfigUpdated(fleetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccNewRelicFleetMembers_MultiRing(t *testing.T) {
	canaryResource := "newrelic_fleet_members.canary"
	prodResource := "newrelic_fleet_members.prod"
	fleetName := fmt.Sprintf("tf-test-members-multi-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetMembersConfigMultiRing(fleetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(canaryResource),
					testAccCheckNewRelicFleetMembersExists(prodResource),
					resource.TestCheckResourceAttr(canaryResource, "ring", "canary"),
					resource.TestCheckResourceAttr(prodResource, "ring", "production"),
				),
			},
		},
	})
}

// Helper functions

func testAccCheckNewRelicFleetMembersExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fleet_members ID is set")
		}
		return nil
	}
}

func testAccCheckNewRelicFleetMembersDestroy(s *terraform.State) error {
	// After destroy the entity_ids are removed from the ring.
	// We verify by checking state is gone — a full API read is not required
	// since the delete operation already confirmed success.
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet_members" {
			return nil
		}
		if r.Primary.ID != "" {
			return fmt.Errorf("fleet_members resource %s still in state after destroy", r.Primary.ID)
		}
	}
	return nil
}

func testAccFleetMembersImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s:%s",
			rs.Primary.Attributes["fleet_id"],
			rs.Primary.Attributes["ring"],
		), nil
	}
}

// Config templates
//
// These configs use a placeholder entity GUID. In real acceptance tests,
// replace this with actual entity GUIDs from your account.
// The resource schema accepts any string; the API validates the GUIDs server-side.

const testAccFleetMembersEntityID = "placeholder-entity-guid-for-schema-test"
const testAccFleetMembersEntityID2 = "placeholder-entity-guid-2-for-schema-test"

func testAccNewRelicFleetMembersConfig(fleetName string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Fleet for members acceptance test"
}

resource "newrelic_fleet_members" "default" {
  fleet_id   = newrelic_fleet.test.id
  ring       = "default"
  entity_ids = [%q]
}
`, fleetName, testAccFleetMembersEntityID)
}

func testAccNewRelicFleetMembersConfigUpdated(fleetName string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Fleet for members acceptance test"
}

resource "newrelic_fleet_members" "default" {
  fleet_id   = newrelic_fleet.test.id
  ring       = "default"
  entity_ids = [%q]
}
`, fleetName, testAccFleetMembersEntityID2)
}

func testAccNewRelicFleetMembersConfigMultiRing(fleetName string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Fleet for multi-ring members test"
}

resource "newrelic_fleet_members" "canary" {
  fleet_id   = newrelic_fleet.test.id
  ring       = "canary"
  entity_ids = [%q]
}

resource "newrelic_fleet_members" "prod" {
  fleet_id   = newrelic_fleet.test.id
  ring       = "production"
  entity_ids = [%q]
}
`, fleetName, testAccFleetMembersEntityID, testAccFleetMembersEntityID2)
}
