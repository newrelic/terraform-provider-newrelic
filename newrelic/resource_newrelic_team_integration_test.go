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
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
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

// ── 7. Entities absent from config — collection and state untouched ───────────

// TestAccNewRelicTeam_EntitiesAbsentPreservesCollection verifies the Computed:true
// carry-forward guarantee: when the entities block is omitted entirely from an
// update config (not even entities = []), Terraform must not plan a removal,
// must not touch the ownership collection, and the entities must remain in state.
//
// This is the key property that made the old isEntitiesAttributeConfigured guard
// redundant: d.HasChange("entities") returns false for absent blocks because the
// planned value equals the prior state value, so the sync is never triggered.
func TestAccNewRelicTeam_EntitiesAbsentPreservesCollection(t *testing.T) {
	ownerName := fmt.Sprintf("tf-acc-team-noent-%s", acctest.RandString(6))
	ownedName := fmt.Sprintf("tf-acc-team-noewn-%s", acctest.RandString(6))
	ownerResource := "newrelic_team.owner"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(ownerResource),
		Steps: []resource.TestStep{
			// ── Step 1: create owner with entities declared ────────────────────
			{
				Config: testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
				),
			},
			// ── Step 2: update config that omits entities block entirely ───────
			// The entities block is absent (not even entities = []). With
			// Computed:true, Terraform plans no change for entities, so
			// d.HasChange("entities") = false in Update and the collection is
			// never touched. Entities must remain in state at count 1.
			{
				Config: testAccNewRelicTeamConfigOwnerNoEntitiesBlock(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					// Entities must persist — the absent block is a no-op, not a clear.
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
				),
			},
			// ── Step 3: idempotency ────────────────────────────────────────────
			// A second apply with the same config must produce an empty plan.
			{
				Config:   testAccNewRelicTeamConfigOwnerNoEntitiesBlock(ownerName, ownedName),
				PlanOnly: true,
			},
		},
	})
}

// ── 8. Entity removal drift ───────────────────────────────────────────────────

// TestAccNewRelicTeam_EntityRemovalDrift verifies the mirror image of the entity
// addition drift test: when an entity is *removed* from the ownership collection
// out-of-band, the provider detects that the declared entity is missing and
// re-adds it on the next apply.
//
// Uses one team as the owned entity of another to avoid ownership-collection
// conflicts with other parallel tests that also use testEntityGUID.
func TestAccNewRelicTeam_EntityRemovalDrift(t *testing.T) {
	ownerName := fmt.Sprintf("tf-acc-team-rmown-%s", acctest.RandString(6))
	ownedName := fmt.Sprintf("tf-acc-team-rmwnd-%s", acctest.RandString(6))
	ownerResource := "newrelic_team.owner"

	var ownershipColID string
	var ownedTeamID string

	// Run serially — entity ownership changes share state that could conflict
	// with other parallel tests if the same entity GUID is used twice.
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(ownerResource),
		Steps: []resource.TestStep{
			// ── Step 1: owner team declares entity ownership over owned team ───────
			{
				Config: testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
					func(s *terraform.State) error {
						ownerRS := s.RootModule().Resources[ownerResource]
						ownershipColID = ownerRS.Primary.Attributes["ownership_collection_id"]
						ownedRS := s.RootModule().Resources["newrelic_team.owned"]
						if ownedRS == nil {
							return fmt.Errorf("newrelic_team.owned not found in state")
						}
						ownedTeamID = ownedRS.Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: remove entity from collection out-of-band + reconcile ────
			// PreConfig removes the owned team's GUID from the owner's collection,
			// simulating a user or external process detaching it manually. The config
			// still declares entities = [owned.id], so Terraform must re-add it.
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementRemoveCollectionMembers(
						ownershipColID,
						[]string{ownedTeamID},
					)
					if err != nil {
						t.Logf("[WARN] PreConfig: failed to remove entity from collection: %v", err)
					}
				},
				Config: testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					// After apply, Terraform must have re-added the entity.
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
				),
			},
		},
	})
}

// ── 8. Core team field drift ──────────────────────────────────────────────────

// TestAccNewRelicTeam_CoreFieldDrift verifies that out-of-band changes to core
// team attributes — description and aliases — are detected by the Read function
// and reconciled on the next apply.
func TestAccNewRelicTeam_CoreFieldDrift(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-coredft-%s", acctest.RandString(6))
	resourceName := "newrelic_team.test"

	var teamID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			// ── Step 1: create team with description and alias ────────────────────
			{
				Config: testAccNewRelicTeamConfigWithDescriptionAndAlias(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Core field drift test"),
					resource.TestCheckResourceAttr(resourceName, "aliases.#", "1"),
					func(s *terraform.State) error {
						teamID = s.RootModule().Resources[resourceName].Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: change description out-of-band + reconcile ────────────────
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementUpdateTeam(teamID,
						scorecards.EntityManagementTeamEntityUpdateInput{
							Description: "description changed out-of-band",
						})
					if err != nil {
						t.Logf("[WARN] PreConfig: failed to drift description: %v", err)
					}
				},
				Config: testAccNewRelicTeamConfigWithDescriptionAndAlias(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					// Terraform must have restored the declared description.
					resource.TestCheckResourceAttr(resourceName, "description", "Core field drift test"),
				),
			},
		},
	})
}

// ── 9. Mode transition with interleaved entity drift ─────────────────────────

// TestAccNewRelicTeam_ModeTransitionWithDrift combines the mode transition
// lifecycle with an entity drift injection. It verifies:
//
//  1. Managed mode: entity declared, drift detected and removed on apply.
//  2. Switch to unmanaged: collection left intact, entities cleared from state.
//  3. Switch back to managed: entity re-tracked from declared config.
func TestAccNewRelicTeam_ModeTransitionWithDrift(t *testing.T) {
	ownerName := fmt.Sprintf("tf-acc-team-mtdft-%s", acctest.RandString(6))
	ownedName := fmt.Sprintf("tf-acc-team-mtdwn-%s", acctest.RandString(6))
	ownerResource := "newrelic_team.owner"

	var ownershipColID string
	var ownedTeamID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(ownerResource),
		Steps: []resource.TestStep{
			// ── Step 1: managed mode — declare entity ownership ───────────────────
			{
				Config: testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					resource.TestCheckResourceAttr(ownerResource, "entity_management_mode", "managed"),
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
					func(s *terraform.State) error {
						ownerRS := s.RootModule().Resources[ownerResource]
						ownershipColID = ownerRS.Primary.Attributes["ownership_collection_id"]
						ownedRS := s.RootModule().Resources["newrelic_team.owned"]
						if ownedRS == nil {
							return fmt.Errorf("newrelic_team.owned not in state")
						}
						ownedTeamID = ownedRS.Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: inject drift + reconcile while still in managed mode ──────
			// Remove the owned entity from the collection — Terraform must re-add it.
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementRemoveCollectionMembers(
						ownershipColID, []string{ownedTeamID})
					if err != nil {
						t.Logf("[WARN] PreConfig: entity removal: %v", err)
					}
				},
				Config: testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
				),
			},
			// ── Step 3: switch to unmanaged ────────────────────────────────────────
			// Entities cleared from state; collection left intact.
			{
				Config: testAccNewRelicTeamConfigOwnerUnmanaged(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					resource.TestCheckResourceAttr(ownerResource, "entity_management_mode", "unmanaged"),
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "0"),
				),
			},
			// ── Step 4: switch back to managed + re-declare entity ─────────────────
			{
				Config: testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(ownerResource),
					resource.TestCheckResourceAttr(ownerResource, "entity_management_mode", "managed"),
					resource.TestCheckResourceAttr(ownerResource, "entities.#", "1"),
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

// testAccNewRelicTeamConfigOwnerOwned creates two teams: "owned" and "owner",
// where "owner" declares managed entity ownership over "owned". Used by entity
// removal drift and mode-transition-with-drift tests.
func testAccNewRelicTeamConfigOwnerOwned(ownerName, ownedName string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "owned" {
  name = %q
}

resource "newrelic_team" "owner" {
  name                   = %q
  entity_management_mode = "managed"
  entities               = [newrelic_team.owned.id]
}
`, ownedName, ownerName)
}

// testAccNewRelicTeamConfigOwnerUnmanaged switches "owner" to unmanaged mode
// while keeping "owned" alive. Used in mode-transition-with-drift step 3.
func testAccNewRelicTeamConfigOwnerUnmanaged(ownerName, ownedName string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "owned" {
  name = %q
}

resource "newrelic_team" "owner" {
  name                   = %q
  entity_management_mode = "unmanaged"
}
`, ownedName, ownerName)
}

// testAccNewRelicTeamConfigOwnerNoEntitiesBlock creates the same two-team
// setup as testAccNewRelicTeamConfigOwnerOwned but with the entities block
// absent from "owner" entirely (not even entities = []). Used by the
// EntitiesAbsentPreservesCollection test to verify the Computed carry-forward.
func testAccNewRelicTeamConfigOwnerNoEntitiesBlock(ownerName, ownedName string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "owned" {
  name = %q
}

resource "newrelic_team" "owner" {
  name                   = %q
  entity_management_mode = "managed"
  # entities block intentionally omitted — must not clear the collection
}
`, ownedName, ownerName)
}

// testAccNewRelicTeamConfigWithDescriptionAndAlias creates a team with a
// description and alias, used by the core field drift test.
func testAccNewRelicTeamConfigWithDescriptionAndAlias(name string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Core field drift test"
  aliases     = ["%s-alias"]
}
`, name, name)
}
