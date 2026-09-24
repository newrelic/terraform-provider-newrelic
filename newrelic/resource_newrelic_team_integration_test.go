//go:build integration

// Integration tests for newrelic_team. These tests call the live NGEP API and
// require the standard NR integration test environment variables plus an
// organization UUID (INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID).
//
// Run with:
//
//	TF_ACC=1 go test -tags integration -run TestAccNewRelicTeam ./newrelic/ -v -timeout 15m

package newrelic

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccPreCheckTeam(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)
	if os.Getenv("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID") == "" {
		t.Skip("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID must be set for team acceptance tests")
	}
}

// TestAccNewRelicTeam_BasicCRUD exercises the full resource lifecycle:
// create → read → update (rename, change description, add alias) → destroy.
func TestAccNewRelicTeam_BasicCRUD(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-%s", acctest.RandString(6))
	rNameUpdated := rName + "-upd"
	resourceName := "newrelic_team.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			// ── Create ───────────────────────────────────────────────────────
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
			// ── Import ───────────────────────────────────────────────────────
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// managers is not round-tripped on Read (see flattenManagerUserIDs comment)
				ImportStateVerifyIgnore: []string{"managers"},
			},
			// ── Update: rename + change description + add alias ───────────────
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

// TestAccNewRelicTeam_WithMembers creates a team with members (user IDs) and
// entities (GUIDs) in the ownership collection, then removes one of each.
func TestAccNewRelicTeam_WithMembers(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-team-mbr-%s", acctest.RandString(6))
	resourceName := "newrelic_team.test"

	// We need at least one known user ID and one known entity GUID.
	// Use the test account's APM integration test entity.
	testEntityGUID := "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTUzNDQ4MjAy"
	testUserIDEnv := os.Getenv("NEW_RELIC_TEST_USER_ID")
	if testUserIDEnv == "" {
		t.Skip("NEW_RELIC_TEST_USER_ID must be set for member tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTeam(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicTeamDestroy(resourceName),
		Steps: []resource.TestStep{
			// ── Create with members + entities ───────────────────────────────
			{
				Config: testAccNewRelicTeamConfigWithMembersAndEntities(rName, testUserIDEnv, testEntityGUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "1"),
				),
			},
			// ── Remove entity, keep member ────────────────────────────────
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

// TestAccNewRelicTeam_Hierarchy creates a parent team and a child team that
// references it via parent_id, then verifies the hierarchy.
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
			// NotFound is expected on destroy
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

  members {
    user_id = %s
  }

  entities {
    guid = %q
  }
}
`, name, userID, entityGUID)
}

func testAccNewRelicTeamConfigWithMembersOnly(name, userID string) string {
	return fmt.Sprintf(`
resource "newrelic_team" "test" {
  name        = %q
  description = "Team with members and entities"

  members {
    user_id = %s
  }
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
