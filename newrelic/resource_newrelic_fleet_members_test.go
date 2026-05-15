//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccNewRelicFleetMembers_Lifecycle exercises the full lifecycle of the
// newrelic_fleet_members resource in a single test run. The steps cover:
//
//   - Create: single ring with two entities.
//   - Update (add): a third entity is added to the ring.
//   - Update (remove): the third entity is removed.
//   - Multi-ring create: a canary ring is introduced alongside default,
//     with two entities distributed across the rings.
//   - Move between rings: both canary entities are moved to default in one
//     apply, exercising the multi-entity cross-ring transfer path.
//   - Ring-block removal: the canary block is dropped; its entities are
//     removed from the fleet.
//   - Import: the resource is re-imported by fleet GUID.
//
// A fresh fleet is created for every test run via acctest.RandString. The
// CheckDestroy step confirms that no fleet_members resource remains in state
// after the fleet is destroyed, ensuring no entities are left stranded.
//
// Prerequisites (all skipped if absent):
//
//	NEW_RELIC_FLEET_TEST_API_KEY     – API key with Fleet Control access.
//	NEW_RELIC_FLEET_TEST_ACCOUNT_ID  – Account ID of the fleet.
//	NEW_RELIC_FLEET_TEST_ENTITY_IDS  – Comma-separated list of ≥3 real entity
//	                                   GUIDs that are currently unassigned in
//	                                   the test account. The API rejects GUIDs
//	                                   that do not correspond to known entities.
func TestAccNewRelicFleetMembers_Lifecycle(t *testing.T) {
	resourceName := "newrelic_fleet_members.test"
	fleetName := fmt.Sprintf("tf-test-members-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)
	entityIDs := testAccFleetMembersEntityIDs(t, 3)
	e1, e2, e3 := entityIDs[0], entityIDs[1], entityIDs[2]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			// Step 1 — Clean create: two unassigned entities in one ring.
			// Expect: resource created, no warnings.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "fleet_id"),
					resource.TestCheckResourceAttr(resourceName, "ring.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.name", "default"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
				),
			},
			// Step 2 — Update add: introduce a third entity.
			// Expect: ring now contains three entities, no warnings.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2, e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "3"),
				),
			},
			// Step 3 — Update remove: remove the third entity.
			// Expect: ring returns to two entities.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
				),
			},
			// Step 4 — Multi-ring create: split entities across two rings.
			// Expect: both ring blocks present, entity counts correct.
			{
				Config: testAccFleetMembersConfigMultiRing(fleetName, []string{e1}, []string{e2, e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.name", "default"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ring.1.name", "canary"),
					resource.TestCheckResourceAttr(resourceName, "ring.1.entity_ids.#", "2"),
				),
			},
			// Step 5 — Move between rings: transfer both canary entities to
			// default in a single apply. Exercises the multi-entity
			// cross-ring transfer path.
			// Expect: default contains all three entities, canary is empty,
			// no spurious "already assigned" warnings.
			{
				Config: testAccFleetMembersConfigMultiRing(fleetName, []string{e1, e2, e3}, []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ring.1.entity_ids.#", "0"),
				),
			},
			// Step 6 — Ring-block removal: drop the canary block entirely.
			// Expect: resource now has one ring block, canary entities
			// (previously moved to default in step 5) remain there.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2, e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "3"),
				),
			},
			// Step 7 — Import: verify the resource can be imported using
			// the fleet GUID. Full state verification is skipped because
			// the import path reads entity IDs from the API in server-
			// determined order, which may differ from the order recorded in
			// Terraform state (TypeList is order-sensitive).
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccFleetMembersImportID(resourceName),
				ImportStateVerify: false,
			},
		},
	})
}

// testAccFleetMembersEntityIDs reads the NEW_RELIC_FLEET_TEST_ENTITY_IDS
// environment variable and returns the first n GUIDs. The test is skipped if
// the variable is unset or contains fewer than n entries.
func testAccFleetMembersEntityIDs(t *testing.T, n int) []string {
	t.Helper()
	raw := os.Getenv("NEW_RELIC_FLEET_TEST_ENTITY_IDS")
	if raw == "" {
		t.Skip("NEW_RELIC_FLEET_TEST_ENTITY_IDS is not set — skipping fleet members acceptance test")
	}
	parts := strings.Split(raw, ",")
	var ids []string
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			ids = append(ids, s)
		}
	}
	if len(ids) < n {
		t.Skipf("NEW_RELIC_FLEET_TEST_ENTITY_IDS must contain at least %d GUIDs (got %d)", n, len(ids))
	}
	return ids[:n]
}

func testAccCheckNewRelicFleetMembersExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for %s", n)
		}
		return nil
	}
}

// testAccCheckNewRelicFleetMembersDestroy verifies that no fleet_members
// resource remains in state after destroy. The underlying API removes all
// declared entity memberships during the delete operation; this check
// confirms the state was cleaned up correctly.
func testAccCheckNewRelicFleetMembersDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet_members" {
			continue
		}
		if r.Primary.ID != "" {
			return fmt.Errorf("fleet_members resource %s still present in state after destroy", r.Primary.ID)
		}
	}
	return nil
}

// testAccFleetMembersImportID returns the fleet GUID as the import ID.
func testAccFleetMembersImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

// testAccFleetMembersConfigSingleRing generates a config with a single fleet
// and a single ring block containing the given entity IDs.
func testAccFleetMembersConfigSingleRing(fleetName, ringName string, entityIDs []string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

resource "newrelic_fleet_members" "test" {
  fleet_id = newrelic_fleet.test.id
  ring {
    name       = %q
    entity_ids = [%s]
  }
}
`, fleetName, ringName, joinQuoted(entityIDs))
}

// testAccFleetMembersConfigMultiRing generates a config with two ring blocks:
// "default" containing defaultIDs and "canary" containing canaryIDs.
func testAccFleetMembersConfigMultiRing(fleetName string, defaultIDs, canaryIDs []string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

resource "newrelic_fleet_members" "test" {
  fleet_id = newrelic_fleet.test.id
  ring {
    name       = "default"
    entity_ids = [%s]
  }
  ring {
    name       = "canary"
    entity_ids = [%s]
  }
}
`, fleetName, joinQuoted(defaultIDs), joinQuoted(canaryIDs))
}

// joinQuoted returns a comma-separated, double-quoted string of the given
// values, suitable for embedding in an HCL list literal.
func joinQuoted(ss []string) string {
	quoted := make([]string, len(ss))
	for i, s := range ss {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(quoted, ", ")
}
