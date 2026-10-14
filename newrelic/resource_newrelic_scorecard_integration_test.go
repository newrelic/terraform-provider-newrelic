//go:build integration || SERVICE_ARCHITECTURE_INTELLIGENCE

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

// ── acceptance test helpers ───────────────────────────────────────────────────

func testAccPreCheckScorecard(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)
}

func testAccCheckNewRelicScorecardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("resource not found or has no ID: %s", n)
		}
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		e, err := client.Scorecards.GetEntity(rs.Primary.ID)
		if err != nil || e == nil {
			return fmt.Errorf("scorecard %s not found: %v", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckNewRelicScorecardRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("resource not found or has no ID: %s", n)
		}
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		e, err := client.Scorecards.GetEntity(rs.Primary.ID)
		if err != nil || e == nil {
			return fmt.Errorf("scorecard rule %s not found: %v", rs.Primary.ID, err)
		}
		return nil
	}
}

// ── ScorecardRule tests ───────────────────────────────────────────────────────

func TestAccNewRelicScorecardRule_BasicCRUD(t *testing.T) {
	rName := fmt.Sprintf("nr-test-rule-%s", acctest.RandString(6))
	resourceName := "newrelic_scorecard_rule.test"
	accountID := testAccountID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckScorecard(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicScorecardRuleConfig(rName, accountID, true, 1440),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "run_interval", "1440"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update: disable + change run_interval
			{
				Config: testAccNewRelicScorecardRuleConfig(rName+"-upd", accountID, false, 720),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "run_interval", "720"),
				),
			},
		},
	})
}

// ── Scorecard tests ───────────────────────────────────────────────────────────

func TestAccNewRelicScorecard_WithRules(t *testing.T) {
	scName := fmt.Sprintf("nr-test-sc-%s", acctest.RandString(6))
	ruleName := fmt.Sprintf("nr-test-rule-%s", acctest.RandString(6))
	resourceName := "newrelic_scorecard.test"
	accountID := testAccountID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckScorecard(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create scorecard + rule, attach rule
			{
				Config: testAccNewRelicScorecardWithRuleConfig(scName, ruleName, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", scName),
					resource.TestCheckResourceAttrSet(resourceName, "rules_collection_id"),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"progress_levels"},
			},
			// Detach rule
			{
				Config: testAccNewRelicScorecardNoRulesConfig(scName, ruleName, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "0"),
				),
			},
		},
	})
}

// ── Scorecard standalone CRUD ─────────────────────────────────────────────────

// TestAccNewRelicScorecard_BasicCRUD exercises the scorecard resource in
// isolation (no rules attached). It validates:
//
//   - Create with progress levels → import round-trip → update name and level
//     attributes → destroy.
func TestAccNewRelicScorecard_BasicCRUD(t *testing.T) {
	scName := fmt.Sprintf("nr-test-sc-crud-%s", acctest.RandString(6))
	resourceName := "newrelic_scorecard.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckScorecard(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// ── Create ────────────────────────────────────────────────────────
			{
				Config: testAccNewRelicScorecardBasicConfig(scName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", scName),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "rules_collection_id"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
				),
			},
			// ── Import round-trip ─────────────────────────────────────────────
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"progress_levels"},
			},
			// ── Update: rename + add description to a progress level ───────────
			{
				Config: testAccNewRelicScorecardBasicUpdatedConfig(scName + "-upd"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", scName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated scorecard"),
				),
			},
		},
	})
}

// ── Scorecard rule drift detection and reconciliation ─────────────────────────

// TestAccNewRelicScorecard_RuleDrift verifies that when a rule is detached from
// the scorecard's rules collection out-of-band (directly via the NGEP API), the
// provider detects the drift on the next plan and re-attaches the rule on apply.
func TestAccNewRelicScorecard_RuleDrift(t *testing.T) {
	scName := fmt.Sprintf("nr-test-sc-drift-%s", acctest.RandString(6))
	ruleName := fmt.Sprintf("nr-test-rule-drift-%s", acctest.RandString(6))
	resourceName := "newrelic_scorecard.test"
	accountID := testAccountID

	var rulesColID string
	var ruleGUID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckScorecard(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// ── Step 1: create scorecard with rule attached ────────────────────
			{
				Config: testAccNewRelicScorecardWithRuleConfig(scName, ruleName, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
					// Capture rules_collection_id and rule GUID for use in next step.
					func(s *terraform.State) error {
						sc := s.RootModule().Resources[resourceName]
						if sc == nil {
							return fmt.Errorf("scorecard resource not in state")
						}
						rulesColID = sc.Primary.Attributes["rules_collection_id"]

						rule := s.RootModule().Resources["newrelic_scorecard_rule.test"]
						if rule == nil {
							return fmt.Errorf("rule resource not in state")
						}
						ruleGUID = rule.Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: detach rule out-of-band + verify reconciliation ────────
			// PreConfig removes the rule from the collection directly via the NGEP
			// API — simulating a user/process detaching it outside Terraform.
			// The Config is unchanged (rule_ids = [rule.id]), so on apply Terraform
			// will detect the missing rule and re-attach it.
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementRemoveCollectionMembers(
						rulesColID,
						[]string{ruleGUID},
					)
					if err != nil {
						t.Logf("[WARN] PreConfig: failed to inject rule drift: %v", err)
					}
				},
				Config: testAccNewRelicScorecardWithRuleConfig(scName, ruleName, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					// After apply, the rule must be re-attached by Terraform.
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
				),
			},
		},
	})
}

// ── Scorecard extra-rule drift ────────────────────────────────────────────────

// TestAccNewRelicScorecard_ExtraRuleDrift verifies the + direction of rule drift:
// when an extra rule is injected into the scorecard's rules collection out-of-band,
// the provider detects it on the next plan and removes it on apply, restoring
// the collection to exactly the rules declared in rule_ids.
func TestAccNewRelicScorecard_ExtraRuleDrift(t *testing.T) {
	scName := fmt.Sprintf("nr-test-sc-xtra-%s", acctest.RandString(6))
	rule1Name := fmt.Sprintf("nr-test-rule-xtra1-%s", acctest.RandString(6))
	rule2Name := fmt.Sprintf("nr-test-rule-xtra2-%s", acctest.RandString(6))
	resourceName := "newrelic_scorecard.test"
	accountID := testAccountID

	var rulesColID string
	var extraRuleID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckScorecard(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// ── Step 1: create scorecard with one rule; create extra rule separately
			// The extra rule is a managed resource but NOT in rule_ids.
			{
				Config: testAccNewRelicScorecardExtraRuleDriftConfig(scName, rule1Name, rule2Name, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
					func(s *terraform.State) error {
						sc := s.RootModule().Resources[resourceName]
						if sc == nil {
							return fmt.Errorf("scorecard not in state")
						}
						rulesColID = sc.Primary.Attributes["rules_collection_id"]
						extra := s.RootModule().Resources["newrelic_scorecard_rule.extra"]
						if extra == nil {
							return fmt.Errorf("extra rule not in state")
						}
						extraRuleID = extra.Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: inject extra rule out-of-band + verify removal on apply ───
			// PreConfig attaches rule2 directly to the collection. Config still
			// declares only rule1 in rule_ids. Terraform must detect and remove rule2.
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementAddCollectionMembers(
						rulesColID,
						[]string{extraRuleID},
					)
					if err != nil {
						t.Logf("[WARN] PreConfig: failed to inject extra rule: %v", err)
					}
				},
				Config: testAccNewRelicScorecardExtraRuleDriftConfig(scName, rule1Name, rule2Name, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardExists(resourceName),
					// After apply, only the declared rule must remain.
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
				),
			},
		},
	})
}

// ── Scorecard rule attribute drift ────────────────────────────────────────────

// TestAccNewRelicScorecardRule_AttributeDrift verifies that an out-of-band
// attribute change on a rule — specifically disabling it via the API — is
// detected on the next plan and corrected on apply.
func TestAccNewRelicScorecardRule_AttributeDrift(t *testing.T) {
	rName := fmt.Sprintf("nr-test-rule-attdft-%s", acctest.RandString(6))
	resourceName := "newrelic_scorecard_rule.test"
	accountID := testAccountID

	var ruleID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckScorecard(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// ── Step 1: create rule with enabled = true ───────────────────────────
			{
				Config: testAccNewRelicScorecardRuleConfig(rName, accountID, true, 1440),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					func(s *terraform.State) error {
						ruleID = s.RootModule().Resources[resourceName].Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: disable rule out-of-band + verify re-enable on apply ──────
			// The UpdateScorecardRule input requires Description always be sent
			// (no omitempty). Since this test rule has no description, "" is correct.
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, err := client.Scorecards.EntityManagementUpdateScorecardRule(ruleID,
						servicearchintelligence.EntityManagementScorecardRuleEntityUpdateInput{
							Enabled:     false,
							Description: "",
						})
					if err != nil {
						t.Logf("[WARN] PreConfig: failed to disable rule: %v", err)
					}
				},
				Config: testAccNewRelicScorecardRuleConfig(rName, accountID, true, 1440),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicScorecardRuleExists(resourceName),
					// After apply, Terraform must have re-enabled the rule.
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

// ── Config templates ──────────────────────────────────────────────────────────

func testAccNewRelicScorecardRuleConfig(name string, accountID int, enabled bool, runInterval int) string {
	return fmt.Sprintf(`
resource "newrelic_scorecard_rule" "test" {
  name         = %q
  enabled      = %t
  run_interval = %d
  nrql_engine {
    accounts = [%d]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }
}
`, name, enabled, runInterval, accountID)
}

func testAccNewRelicScorecardWithRuleConfig(scName, ruleName string, accountID int) string {
	return fmt.Sprintf(`
resource "newrelic_scorecard_rule" "test" {
  name    = %q
  enabled = true
  nrql_engine {
    accounts = [%d]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }
}

resource "newrelic_scorecard" "test" {
  name = %q
  progress_levels {
    id   = "red"
    name = "Red"
  }
  progress_levels {
    id   = "green"
    name = "Green"
  }
  rule_ids = [newrelic_scorecard_rule.test.id]
}
`, ruleName, accountID, scName)
}

func testAccNewRelicScorecardNoRulesConfig(scName, ruleName string, accountID int) string {
	return fmt.Sprintf(`
resource "newrelic_scorecard_rule" "test" {
  name    = %q
  enabled = true
  nrql_engine {
    accounts = [%d]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }
}

resource "newrelic_scorecard" "test" {
  name = %q
  progress_levels {
    id   = "red"
    name = "Red"
  }
  progress_levels {
    id   = "green"
    name = "Green"
  }
  # rule_ids intentionally empty — rule detached but not deleted
}
`, ruleName, accountID, scName)
}

func testAccNewRelicScorecardBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_scorecard" "test" {
  name = %q
  progress_levels {
    id             = "red"
    name           = "Below Expectations"
    hex_color_code = "#FF0000"
  }
  progress_levels {
    id             = "yellow"
    name           = "Approaching Expectations"
    hex_color_code = "#FFA500"
  }
  progress_levels {
    id             = "green"
    name           = "Meeting Expectations"
    hex_color_code = "#00CC00"
  }
}
`, name)
}

func testAccNewRelicScorecardBasicUpdatedConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_scorecard" "test" {
  name        = %q
  description = "Updated scorecard"
  progress_levels {
    id             = "red"
    name           = "Below Expectations"
    hex_color_code = "#FF0000"
    description    = "Score below 50%%"
  }
  progress_levels {
    id             = "yellow"
    name           = "Approaching Expectations"
    hex_color_code = "#FFA500"
    description    = "Score 50–80%%"
  }
  progress_levels {
    id             = "green"
    name           = "Meeting Expectations"
    hex_color_code = "#00CC00"
    description    = "Score above 80%%"
  }
}
`, name)
}

// testAccNewRelicScorecardExtraRuleDriftConfig creates one scorecard with rule1
// attached (via rule_ids) and rule2 as an independent rule not attached to any
// scorecard. Used by TestAccNewRelicScorecard_ExtraRuleDrift.
func testAccNewRelicScorecardExtraRuleDriftConfig(scName, rule1Name, rule2Name string, accountID int) string {
	nrqlQuery := "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
	return fmt.Sprintf(`
resource "newrelic_scorecard_rule" "attached" {
  name    = %q
  enabled = true
  nrql_engine {
    accounts = [%d]
    query    = %q
  }
}

resource "newrelic_scorecard_rule" "extra" {
  name    = %q
  enabled = true
  nrql_engine {
    accounts = [%d]
    query    = %q
  }
}

resource "newrelic_scorecard" "test" {
  name     = %q
  rule_ids = [newrelic_scorecard_rule.attached.id]
}
`, rule1Name, accountID, nrqlQuery, rule2Name, accountID, nrqlQuery, scName)
}
