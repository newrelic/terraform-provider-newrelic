//go:build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
