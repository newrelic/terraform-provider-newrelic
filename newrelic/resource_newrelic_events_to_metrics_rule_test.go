//go:build integration || EVENTS

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicEventsToMetricsRule_Basic(t *testing.T) {
	rand := acctest.RandString(5)
	name := fmt.Sprintf("events_to_metrics_rule_%s", rand)
	resourceName := "newrelic_events_to_metrics_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicEventsToMetricsRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicEventsToMetricsRuleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEventsToMetricsRuleExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicEventsToMetricsRuleConfigUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEventsToMetricsRuleExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

// TestAccNewRelicEventsToMetricsRule_ProviderAccountIDDefault verifies that
// omitting account_id causes the resource to fall back to the provider's
// configured account ID, not send 0 to the API (the root cause of NR-630695).
func TestAccNewRelicEventsToMetricsRule_ProviderAccountIDDefault(t *testing.T) {
	rand := acctest.RandString(5)
	name := fmt.Sprintf("events_to_metrics_rule_default_acct_%s", rand)
	resourceName := "newrelic_events_to_metrics_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicEventsToMetricsRuleDestroy,
		Steps: []resource.TestStep{
			// Create without account_id — must succeed and default to provider account.
			{
				Config: testAccNewRelicEventsToMetricsRuleConfigNoAccountID(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEventsToMetricsRuleExists(resourceName),
					// account_id must be set to the provider account, not zero.
					resource.TestCheckResourceAttr(resourceName, "account_id", strconv.Itoa(testAccountID)),
				),
			},
		},
	})
}

func testAccCheckNewRelicEventsToMetricsRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_events_to_metrics_rule" {
			continue
		}

		accountID, ruleID, err := getEventsToMetricsRuleIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.EventsToMetrics.GetRule(accountID, ruleID)

		if err == nil {
			return fmt.Errorf("events to metrics rule still exists: %s", err)
		}
	}

	return nil
}

func testAccCheckNewRelicEventsToMetricsRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		accountID, ruleID, err := getEventsToMetricsRuleIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.EventsToMetrics.GetRule(accountID, ruleID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicEventsToMetricsRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_events_to_metrics_rule" "foo" {
  account_id = "%d"
  name = "%s"
  description = "test description"
  nrql = "SELECT uniqueCount(account_id) AS `+"`"+"Transaction.account_id"+"`"+` FROM Transaction FACET appName, name"
}

`, testAccountID, name)
}

func testAccNewRelicEventsToMetricsRuleConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_events_to_metrics_rule" "foo" {
  account_id = "%d"
  name = "%s"
  description = "test description"
  nrql = "SELECT uniqueCount(account_id) AS `+"`"+"Transaction.account_id"+"`"+` FROM Transaction FACET appName, name"
  enabled = false
}
`, testAccountID, name)
}

// testAccNewRelicEventsToMetricsRuleConfigNoAccountID intentionally omits
// account_id to exercise the provider-level default fallback.
func testAccNewRelicEventsToMetricsRuleConfigNoAccountID(name string) string {
	return fmt.Sprintf(`
resource "newrelic_events_to_metrics_rule" "foo" {
  name        = "%s"
  description = "test description — account_id omitted, should default to provider account"
  nrql        = "SELECT uniqueCount(account_id) AS `+"`"+"Transaction.account_id"+"`"+` FROM Transaction FACET appName, name"
}
`, name)
}
