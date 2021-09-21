//go:build integration
// +build integration

package newrelic

import (
	"fmt"
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
