//go:build integration || FLEET

// Package-level note: the tests in this file are intended for local execution
// only and will be automatically skipped in CI. CI does not set the fleet-
// specific environment variables (NEW_RELIC_FLEET_TEST_API_KEY,
// NEW_RELIC_FLEET_TEST_ACCOUNT_ID, NEW_RELIC_FLEET_TEST_ENTITY_IDS,
// NEW_RELIC_FLEET_TEST_AC_ENTITY_IDS) because the tests require real, named
// host entities that exist in a specific account — infrastructure that cannot
// be reproduced generically in a shared CI environment. Each test checks for
// these variables at the top of its setup and calls t.Skip if any are absent.

package newrelic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccNewRelicFleetMembers_Lifecycle exercises the full lifecycle of the
// newrelic_fleet_members resource in the same order as manual scenarios
// S1, S3, S4, S5, S6, S7, S8, S9, S12.
//
// S2 and S10 (Agent Control entity adoption) are covered separately in
// TestAccNewRelicFleetMembers_Adoption.
//
// Prerequisites (all skipped if absent):
//
//	NEW_RELIC_FLEET_TEST_API_KEY     – API key with Fleet Control access.
//	NEW_RELIC_FLEET_TEST_ACCOUNT_ID  – Account ID of the fleet.
//	NEW_RELIC_FLEET_TEST_ENTITY_IDS  – Comma-separated list of ≥3 real entity
//	                                   GUIDs that are currently unassigned in
//	                                   the test account.
func TestAccNewRelicFleetMembers_Lifecycle(t *testing.T) {
	resourceName := "newrelic_fleet_members.test"
	fleetName := fmt.Sprintf("tf-test-members-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)
	entityIDs := testAccFleetMembersEntityIDs(t, 3)
	e1, e2, e3 := entityIDs[0], entityIDs[1], entityIDs[2]
	testAccEnsureEntitiesUnassigned(t, entityIDs)

	// Captured from Step 1 for use in the S5 drift PreConfig.
	var capturedFleetID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			// ── S1 ── Clean create: two unassigned entities in one ring.
			// Expect: resource created, no warnings.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "fleet_id"),
					resource.TestCheckResourceAttr(resourceName, "ring.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.name", "default"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
					func(s *terraform.State) error {
						if rs, ok := s.RootModule().Resources["newrelic_fleet.test"]; ok {
							capturedFleetID = rs.Primary.ID
						}
						return nil
					},
				),
			},
			// ── DS check after S1 ── Data source (no ring filter) must reflect
			// the two members just created by the resource.
			{
				Config: testAccFleetMembersConfigSingleRingWithDS(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.newrelic_fleet_members.test", "members.#", "2"),
				),
			},
			// ── S3 ── Update add: introduce e3.
			// Expect: ring now has three entities, no warnings.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2, e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "3"),
				),
			},
			// ── S4 ── Update remove: remove e3, ring returns to two entities.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
				),
			},
			// ── S5 (plan) ── Drift detection: remove e2 out-of-band before the
			// plan runs and expect Terraform to surface a non-empty diff.
			{
				PreConfig: func() {
					testAccFleetMembersRemoveOutOfBand(t, capturedFleetID, "default", e2)
				},
				Config:             testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// ── S5 (apply) ── Apply re-adds e2, restoring declared state.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
				),
			},
			// ── S6 ── Multi-ring: e1 in default, e2+e3 in canary.
			// Expect: both ring blocks present with correct entity counts.
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
			// ── DS check after S6 ── Two data sources in one step: unfiltered view
			// sees all 3 members across both rings; ring-filtered view sees only the
			// 1 member in "default". This exercises both data source modes together.
			{
				Config: testAccFleetMembersConfigMultiRingWithDualDS(fleetName, []string{e1}, []string{e2, e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.newrelic_fleet_members.all", "members.#", "3"),
					resource.TestCheckResourceAttr("data.newrelic_fleet_members.default_ring", "members.#", "1"),
				),
			},
			// ── S7 ── Move both canary entities to default in one apply.
			// Exercises the multi-entity buildFleetRemovalSet exclusion so
			// that e2 and e3 are not blocked by the "already assigned" pre-check
			// during the add phase.
			{
				Config: testAccFleetMembersConfigMultiRing(fleetName, []string{e1, e2, e3}, []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ring.1.entity_ids.#", "0"),
				),
			},
			// ── Transition ── Drop the now-empty canary block to set up S8.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2, e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "3"),
				),
			},
			// ── S8 ── Add ring block: add a canary ring alongside default.
			// e3 is moved from default to canary, exercising the add-ring-block path.
			{
				Config: testAccFleetMembersConfigMultiRing(fleetName, []string{e1, e2}, []string{e3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ring.1.entity_ids.#", "1"),
				),
			},
			// ── S9 ── Remove ring block: drop canary entirely.
			// e3 is removed from the fleet; default is untouched.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{e1, e2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "2"),
				),
			},
			// ── S12 ── Import: verify re-import by fleet GUID.
			// ImportStateVerify is disabled because the API returns entity IDs in
			// server-determined order, which may differ from the TypeList order
			// recorded in state.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccFleetMembersImportID(resourceName),
				ImportStateVerify: false,
			},
		},
	})
}

// TestAccNewRelicFleetMembers_Adoption covers scenarios S2 and S10 — declaring
// an entity that is already present in the fleet (simulated by a direct API
// add before Terraform runs) alongside unassigned entities.
//
// The entity-already-in-fleet condition triggers an "already assigned" warning
// on create (S2) and on update (S10), after which the entity is adopted into
// Terraform state.
//
// Additional prerequisites (skipped if absent):
//
//	NEW_RELIC_FLEET_TEST_AC_ENTITY_IDS – Comma-separated list of ≥2 entity
//	                                      GUIDs to use as the "already in fleet"
//	                                      entities for S2 and S10 respectively.
//	                                      These must be valid entity GUIDs in
//	                                      the test account but need not be
//	                                      Agent Control managed.
func TestAccNewRelicFleetMembers_Adoption(t *testing.T) {
	resourceName := "newrelic_fleet_members.test"
	fleetName := fmt.Sprintf("tf-test-adoption-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)
	acEntityIDs := testAccFleetMembersACEntityIDs(t, 2)
	ac1, ac2 := acEntityIDs[0], acEntityIDs[1]
	testAccEnsureEntitiesUnassigned(t, acEntityIDs)

	// Fleet ID captured after the fleet is created (Step 1) so that Steps 2
	// and 4 can pre-add the AC entities via API before Terraform runs.
	var capturedFleetID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			// ── Setup ── Create the fleet only (no fleet_members resource yet)
			// so that we can capture the fleet GUID for out-of-band operations.
			{
				Config: testAccFleetMembersConfigFleetOnly(fleetName),
				Check: func(s *terraform.State) error {
					if rs, ok := s.RootModule().Resources["newrelic_fleet.test"]; ok {
						capturedFleetID = rs.Primary.ID
					}
					return nil
				},
			},
			// ── S2 ── Adoption on create: pre-add ac1 to the fleet out-of-band,
			// then declare it in a Terraform Create call.
			// Expect: warning fires for ac1 (already assigned), ac1 is adopted into state.
			{
				PreConfig: func() {
					testAccFleetMembersAddOutOfBand(t, capturedFleetID, "default", ac1)
				},
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{ac1}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFleetMembersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "1"),
				),
			},
			// ── Cleanup ── Remove ac1 from the declared set so the fleet is clean for S10.
			{
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "0"),
				),
			},
			// ── S10 ── Adoption on update: pre-add ac2 to the fleet out-of-band,
			// then add it to the existing declared set via an Update call.
			// Expect: warning fires for ac2, ac2 is adopted into state.
			{
				PreConfig: func() {
					testAccFleetMembersAddOutOfBand(t, capturedFleetID, "default", ac2)
				},
				Config: testAccFleetMembersConfigSingleRing(fleetName, "default", []string{ac2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ring.0.entity_ids.#", "1"),
				),
			},
		},
	})
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// testAccFleetMembersEntityIDs returns ≥n entity GUIDs from
// NEW_RELIC_FLEET_TEST_ENTITY_IDS, skipping the test if unavailable.
func testAccFleetMembersEntityIDs(t *testing.T, n int) []string {
	return testAccFleetEntityIDsFromEnv(t, "NEW_RELIC_FLEET_TEST_ENTITY_IDS", n)
}

// testAccFleetMembersACEntityIDs returns ≥n entity GUIDs from
// NEW_RELIC_FLEET_TEST_AC_ENTITY_IDS, skipping the test if unavailable.
func testAccFleetMembersACEntityIDs(t *testing.T, n int) []string {
	return testAccFleetEntityIDsFromEnv(t, "NEW_RELIC_FLEET_TEST_AC_ENTITY_IDS", n)
}

// testAccFleetEntityIDsFromEnv reads a comma-separated env var and returns
// the first n GUIDs, skipping the test if the variable is unset or too short.
func testAccFleetEntityIDsFromEnv(t *testing.T, envVar string, n int) []string {
	t.Helper()
	raw := os.Getenv(envVar)
	if raw == "" {
		t.Skipf("%s is not set — skipping fleet members test", envVar)
	}
	ids := splitTrimmed(raw)
	if len(ids) < n {
		t.Skipf("%s must contain at least %d GUIDs (got %d)", envVar, n, len(ids))
	}
	return ids[:n]
}

func splitTrimmed(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// testAccEnsureEntitiesUnassigned checks each entity's nr.fleet tag and, if
// any are still assigned to a fleet from a previous failed run, removes them
// so the current test starts from a clean state.
func testAccEnsureEntitiesUnassigned(t *testing.T, entityIDs []string) {
	t.Helper()
	apiKey := os.Getenv("NEW_RELIC_FLEET_TEST_API_KEY")
	if apiKey == "" {
		return
	}
	for _, entityID := range entityIDs {
		fleetID := testAccFleetMembersGetEntityFleet(t, apiKey, entityID)
		if fleetID == "" {
			continue
		}
		t.Logf("testAccEnsureEntitiesUnassigned: %s is in fleet %s from a previous run — removing", entityID, fleetID)
		for _, ring := range []string{"default", "canary"} {
			testAccFleetMembersRemoveOutOfBand(t, fleetID, ring, entityID)
		}
	}
}

// testAccFleetMembersGetEntityFleet queries the nr.fleet entity tag to find
// which fleet the entity currently belongs to. Returns "" if unassigned.
func testAccFleetMembersGetEntityFleet(t *testing.T, apiKey, entityID string) string {
	t.Helper()
	body := fmt.Sprintf(
		`{"query":"query($g:EntityGuid!){actor{entity(guid:$g){tags{key values}}}}","variables":{"g":%q}}`,
		entityID,
	)
	req, _ := http.NewRequest("POST", "https://api.newrelic.com/graphql", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Actor struct {
				Entity struct {
					Tags []struct {
						Key    string   `json:"key"`
						Values []string `json:"values"`
					} `json:"tags"`
				} `json:"entity"`
			} `json:"actor"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}
	for _, tag := range result.Data.Actor.Entity.Tags {
		if tag.Key == "nr.fleet" && len(tag.Values) > 0 {
			return tag.Values[0]
		}
	}
	return ""
}

// testAccFleetMembersMutateOutOfBand executes an add or remove fleet members
// mutation directly against the GraphQL API, bypassing Terraform. Used to
// simulate out-of-band changes for drift detection and adoption tests.
func testAccFleetMembersMutateOutOfBand(t *testing.T, fleetID, ring, entityID, op string) {
	t.Helper()
	if fleetID == "" {
		t.Logf("testAccFleetMembersMutateOutOfBand: fleetID not yet captured, skipping")
		return
	}
	apiKey := os.Getenv("NEW_RELIC_FLEET_TEST_API_KEY")
	if apiKey == "" {
		return
	}
	mutationName := "fleetControlRemoveFleetMembers"
	if op == "add" {
		mutationName = "fleetControlAddFleetMembers"
	}
	body := fmt.Sprintf(
		`{"query":"mutation($f:ID!,$m:[FleetControlFleetMemberRingInput!]){%s(fleetId:$f,members:$m){fleetId}}","variables":{"f":%q,"m":[{"ring":%q,"entityIds":[%q]}]}}`,
		mutationName, fleetID, ring, entityID,
	)
	req, _ := http.NewRequest("POST", "https://api.newrelic.com/graphql", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Logf("testAccFleetMembersMutateOutOfBand: request failed: %v", err)
		return
	}
	resp.Body.Close()
}

// testAccFleetMembersRemoveOutOfBand removes a single entity from a ring via
// the GraphQL API, simulating an out-of-band change for drift-detection tests.
func testAccFleetMembersRemoveOutOfBand(t *testing.T, fleetID, ring, entityID string) {
	testAccFleetMembersMutateOutOfBand(t, fleetID, ring, entityID, "remove")
}

// testAccFleetMembersAddOutOfBand adds a single entity to a ring via the
// GraphQL API, simulating a pre-assignment (e.g. Agent Control) for adoption tests.
func testAccFleetMembersAddOutOfBand(t *testing.T, fleetID, ring, entityID string) {
	testAccFleetMembersMutateOutOfBand(t, fleetID, ring, entityID, "add")
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

// testAccCheckNewRelicFleetMembersDestroy verifies via the API that fleet
// members were actually removed. The SDK calls CheckDestroy with the
// pre-destroy state (IDs are always set), so state-based checks are not
// meaningful here — we query the API instead.
func testAccCheckNewRelicFleetMembersDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_fleet_members" {
			continue
		}
		fleetID := r.Primary.ID
		if fleetID == "" {
			continue
		}
		members, err := fleetMemberIDs(context.Background(), client, fleetID, "")
		if err != nil {
			// Fleet no longer exists — members are implicitly removed.
			continue
		}
		if len(members) > 0 {
			return fmt.Errorf("fleet %s still has %d member(s) after destroy", fleetID, len(members))
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

// ── Config templates ──────────────────────────────────────────────────────────

// testAccFleetMembersFleetBlock returns the HCL for a newrelic_fleet "test"
// resource, reused across all config templates.
func testAccFleetMembersFleetBlock(fleetName string) string {
	return fmt.Sprintf(`resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}`, fleetName)
}

func testAccFleetMembersConfigFleetOnly(fleetName string) string {
	return testAccFleetMembersFleetBlock(fleetName) + "\n"
}

func testAccFleetMembersConfigSingleRing(fleetName, ringName string, entityIDs []string) string {
	return testAccFleetMembersFleetBlock(fleetName) + fmt.Sprintf(`

resource "newrelic_fleet_members" "test" {
  fleet_id = newrelic_fleet.test.id
  ring {
    name       = %q
    entity_ids = [%s]
  }
}
`, ringName, joinQuoted(entityIDs))
}

func testAccFleetMembersConfigMultiRing(fleetName string, defaultIDs, canaryIDs []string) string {
	return testAccFleetMembersFleetBlock(fleetName) + fmt.Sprintf(`

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
`, joinQuoted(defaultIDs), joinQuoted(canaryIDs))
}

func joinQuoted(ss []string) string {
	quoted := make([]string, len(ss))
	for i, s := range ss {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(quoted, ", ")
}

func testAccFleetMembersConfigSingleRingWithDS(fleetName, ringName string, entityIDs []string) string {
	return testAccFleetMembersFleetBlock(fleetName) + fmt.Sprintf(`

resource "newrelic_fleet_members" "test" {
  fleet_id = newrelic_fleet.test.id
  ring {
    name       = %q
    entity_ids = [%s]
  }
}

data "newrelic_fleet_members" "test" {
  fleet_id   = newrelic_fleet.test.id
  depends_on = [newrelic_fleet_members.test]
}
`, ringName, joinQuoted(entityIDs))
}

func testAccFleetMembersConfigMultiRingWithDualDS(fleetName string, defaultIDs, canaryIDs []string) string {
	return testAccFleetMembersFleetBlock(fleetName) + fmt.Sprintf(`

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

data "newrelic_fleet_members" "all" {
  fleet_id   = newrelic_fleet.test.id
  depends_on = [newrelic_fleet_members.test]
}

data "newrelic_fleet_members" "default_ring" {
  fleet_id   = newrelic_fleet.test.id
  ring       = "default"
  depends_on = [newrelic_fleet_members.test]
}
`, joinQuoted(defaultIDs), joinQuoted(canaryIDs))
}
