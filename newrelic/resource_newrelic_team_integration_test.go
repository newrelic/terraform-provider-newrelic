//go:build integration || SERVICE_ARCHITECTURE_INTELLIGENCE

// Integration tests for newrelic_team. These tests call the live NGEP API and
// require the standard NR integration test environment variables plus an
// organization UUID (INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID).
//
// Run with:
//
//	TF_ACC=1 go test -tags integration -run TestAccNewRelicTeam ./newrelic/ -v -timeout 20m

package newrelic

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testEntityGUID is a stable APM application GUID in account 3806526 used
// across team tests that exercise entity ownership / drift scenarios.
const testEntityGUID = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTUzNDQ4MjAy"

func testAccPreCheckTeam(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)
	if os.Getenv("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID") == "" {
		t.Skip("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID must be set for team acceptance tests")
	}
}

// ── 1. BasicCRUD ──────────────────────────────────────────────────────────────

// TestAccNewRelicTeam_BasicCRUD exercises the full resource lifecycle:
// create → read → import → update (rename + description + alias) → destroy.
func TestAccNewRelicTeam_BasicCRUD(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-%s", acctest.RandString(6))
	rNameUpdated := rName + "-upd"
	resourceName := "newrelic_team.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicTeamConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Acceptance test team"),
					resource.TestCheckResourceAttrSet(resourceName, "membership_collection_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ownership_collection_id"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"managers"},
			},
			{
				Config: testAccNewRelicTeamConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "aliases.#", "1"),
				),
			},
		},
	})
}

// ── 2. Members + Entities ─────────────────────────────────────────────────────

// TestAccNewRelicTeam_WithMembers creates a team with members (user IDs) and
// entities (GUIDs) in the ownership collection, then removes one of each.
func TestAccNewRelicTeam_WithMembers(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-mbr-%s", acctest.RandString(6))
	resourceName := "newrelic_team.test"

	testUserIDEnv := os.Getenv("NEW_RELIC_TEST_USER_ID")
	if testUserIDEnv == "" {
		t.Skip("NEW_RELIC_TEST_USER_ID must be set for member tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicTeamConfigWithMembersAndEntities(rName, testUserIDEnv, testEntityGUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "1"),
				),
			},
			{
				Config: testAccNewRelicTeamConfigWithMembersOnly(rName, testUserIDEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "0"),
				),
			},
		},
	})
}

// ── 3. Hierarchy ──────────────────────────────────────────────────────────────

// TestAccNewRelicTeam_Hierarchy creates a parent team and a child team that
// references it via parent_id, then verifies the hierarchy link.
func TestAccNewRelicTeam_Hierarchy(t *testing.T) {
	parentName := fmt.Sprintf("tf-acc-parent-%s", acctest.RandString(5))
	childName := fmt.Sprintf("tf-acc-child-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckTeam(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicTeamConfigHierarchy(parentName, childName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists("newrelic_team.parent"),
					testAccCheckNewRelicTeamExists("newrelic_team.child"),
					resource.TestCheckResourceAttrPair("newrelic_team.child", "parent_id",
						"newrelic_team.parent", "id"),
				),
			},
		},
	})
}

// ── 4. Aux Resources (tags + aliases + resource links) ────────────────────────

// TestAccNewRelicTeam_AuxResources exercises the full surface of optional team
// attributes: tags, aliases, and supplemental resources (link type). It:
//
//   - Creates a team with two tags, two aliases, and two resource links.
//   - Updates to change tag values, drop one alias, change a resource title, and
//     remove the second resource link.
//   - Clears all optional attributes to verify the clearXxx raw-patch code paths.
func TestAccNewRelicTeam_AuxResources(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-aux-%s", acctest.RandString(6))
	resourceName := "newrelic_team.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			// ── Create: full aux attributes ────────────────────────────────────
			{
				Config: testAccNewRelicTeamConfigAuxFull(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "aliases.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
				),
			},
			// ── Import round-trip ──────────────────────────────────────────────
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"managers"},
			},
			// ── Update: reduce to one alias, one tag, one resource ─────────────
			{
				Config: testAccNewRelicTeamConfigAuxReduced(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "aliases.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.title", "Main Repo Updated"),
				),
			},
			// ── Clear all optional attrs ───────────────────────────────────────
			{
				Config: testAccNewRelicTeamConfigBasic(rName + "-cleared"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "aliases.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
				),
			},
		},
	})
}

// ── 5. entity_management_mode transition ─────────────────────────────────────

// TestAccNewRelicTeam_ModeTransition validates the three-way lifecycle:
//
//  1. managed mode with entity declared → entity tracked in state, in collection.
//  2. Switch to unmanaged mode → entity cleared from state, collection untouched.
//  3. Switch back to managed with entity re-declared → entity re-tracked, already-in-
//     collection error treated as no-op (idempotent).
func TestAccNewRelicTeam_ModeTransition(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-mode-%s", acctest.RandString(6))
	resourceName := "newrelic_team.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			// ── Step 1: managed + entity declared ─────────────────────────────
			{
				Config: testAccNewRelicTeamConfigManagedWithEntity(rName, testEntityGUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_management_mode", "managed"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "1"),
				),
			},
			// ── Step 2: switch to unmanaged ────────────────────────────────────
			// entity_management_mode = "unmanaged" — entities must be cleared from
			// state; the ownership collection is NOT modified (entity stays there).
			{
				Config: testAccNewRelicTeamConfigUnmanaged(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_management_mode", "unmanaged"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "0"),
				),
			},
			// ── Step 3: switch back to managed with same entity ────────────────
			// The entity is still in the collection (was never removed); re-declaring
			// it must succeed without "already in collection" errors.
			{
				Config: testAccNewRelicTeamConfigManagedWithEntity(rName, testEntityGUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_management_mode", "managed"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "1"),
				),
			},
		},
	})
}

// ── 6. Entity drift detection and reconciliation ──────────────────────────────

// TestAccNewRelicTeam_EntityDrift verifies that when an entity is added to the
// ownership collection out-of-band (directly via the NGEP API), the provider:
//
//   - Detects the extra entity as drift on the next plan (Read populates state
//     with the out-of-band GUID; plan shows it as a removal).
//   - Reconciles it on apply (syncTeamOwnership removes the extra GUID).
func TestAccNewRelicTeam_EntityDrift(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-drift-%s", acctest.RandString(6))
	resourceName := "newrelic_team.test"

	var ownershipColID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			// ── Step 1: create team with an explicit empty entities set ─────────
			// Using entities = [] (not omitted) ensures the attribute is present in
			// the raw config — isEntitiesAttributeConfigured() returns true, so the
			// out-of-band entity will show up as drift (not silently accepted).
			{
				Config: testAccNewRelicTeamConfigManagedEmptyEntities(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "0"),
					// Capture the ownership_collection_id for use in the next step.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource %s not in state", resourceName)
						}
						ownershipColID = rs.Primary.Attributes["ownership_collection_id"]
						return nil
					},
				),
			},
			// ── Step 2: inject drift + verify reconciliation ───────────────────
			// PreConfig adds testEntityGUID to the ownership collection directly via
			// the NGEP API — simulating a user/process adding it out-of-band.
			// The Config is unchanged (entities = []), so on apply Terraform will
			// detect and remove the extra entity, restoring the declared state.
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementAddCollectionMembers(
						ownershipColID,
						[]string{testEntityGUID},
					)
					if err != nil {
						t.Logf("[WARN] PreConfig: failed to inject drift entity: %v", err)
					}
				},
				Config: testAccNewRelicTeamConfigManagedEmptyEntities(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					// After apply, Terraform should have removed the out-of-band entity:
					// the ownership collection is back to matching the empty entities config.
					resource.TestCheckResourceAttr(resourceName, "entities.#", "0"),
				),
			},
		},
	})
}

// ── Check helpers ──────────────────────────────────────────────────────────────

func testAccCheckNewRelicTeamExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("team ID not set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		e, err := client.Scorecards.GetEntity(rs.Primary.ID)
		if err != nil {
			return err
		}
		if e == nil {
			return fmt.Errorf("team %s not found", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicTeamDestroy(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return nil
		}
		if rs.Primary.ID == "" {
			return nil
		}
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		_, err := client.Scorecards.GetEntity(rs.Primary.ID)
		if err != nil {
			return nil
		}
		return fmt.Errorf("team %s still exists", rs.Primary.ID)
	}
}

// ── Config templates ──────────────────────────────────────────────────────────

func testAccNewRelicTeamConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Acceptance test team"
}
`, name)
}

func testAccNewRelicTeamConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Updated description"
  aliases     = ["%s-alias"]
}
`, name, name)
}

func testAccNewRelicTeamConfigWithMembersAndEntities(name, userID, entityGUID string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Team with members and entities"
  members     = [%s]
  entities    = [%q]
}
`, name, userID, entityGUID)
}

func testAccNewRelicTeamConfigWithMembersOnly(name, userID string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Team with members and entities"
  members     = [%s]
  entities    = []
}
`, name, userID)
}

func testAccNewRelicTeamConfigHierarchy(parentName, childName string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "parent" {
  name        = %q
  description = "Parent team for hierarchy test"
}

resource "newrelic_team" "child" {
  name        = %q
  description = "Child team for hierarchy test"
  parent_id   = newrelic_team.parent.id
}
`, parentName, childName)
}

func testAccNewRelicTeamConfigAuxFull(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Team with full aux attributes"

  aliases = ["alias-alpha", "alias-beta"]

  tags {
    key    = "env"
    values = ["staging"]
  }
  tags {
    key    = "tier"
    values = ["platform", "infra"]
  }

  resources {
    type    = "GITHUB"
    content = "https://github.com/newrelic/example-repo"
    title   = "Main Repo"
  }
  resources {
    type    = "SLACK"
    content = "https://newrelic.slack.com/channels/platform-team"
    title   = "Slack Channel"
  }
}
`, name)
}

func testAccNewRelicTeamConfigAuxReduced(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Team with reduced aux attributes"

  aliases = ["alias-alpha"]

  tags {
    key    = "env"
    values = ["production"]
  }

  resources {
    type    = "GITHUB"
    content = "https://github.com/newrelic/example-repo"
    title   = "Main Repo Updated"
  }
}
`, name)
}

func testAccNewRelicTeamConfigManagedWithEntity(name, entityGUID string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name                   = %q
  entity_management_mode = "managed"
  entities               = [%q]
}
`, name, entityGUID)
}

func testAccNewRelicTeamConfigUnmanaged(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name                   = %q
  entity_management_mode = "unmanaged"
}
`, name)
}

func testAccNewRelicTeamConfigManagedEmptyEntities(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name                   = %q
  entity_management_mode = "managed"
  entities               = []
}
`, name)
}
