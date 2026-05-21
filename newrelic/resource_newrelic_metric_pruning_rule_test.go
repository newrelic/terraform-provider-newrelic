//go:build integration || INGEST

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicMetricPruningRule_Basic(t *testing.T) {
	resourceName := "newrelic_metric_pruning_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMetricPruningRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicMetricPruningRuleConfig(
					"SELECT containerId FROM Metric WHERE metricName = 'scooter.speed.kmph'",
					"tf-test pruning rule for scooter speed metric",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMetricPruningRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "tf-test pruning rule for scooter speed metric"),
					resource.TestCheckResourceAttr(resourceName, "nrql", "SELECT containerId FROM Metric WHERE metricName = 'scooter.speed.kmph'"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicMetricPruningRule_NoDescription(t *testing.T) {
	resourceName := "newrelic_metric_pruning_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMetricPruningRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicMetricPruningRuleNoDescriptionConfig(
					"SELECT rider_id FROM Metric WHERE metricName = 'scooter.engine.temp.celsius'",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMetricPruningRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
				),
			},
		},
	})
}

func TestAccNewRelicMetricPruningRule_AccountIDInheritance(t *testing.T) {
	resourceName := "newrelic_metric_pruning_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMetricPruningRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicMetricPruningRuleNoAccountIDConfig(
					"SELECT zone FROM Metric WHERE metricName = 'scooter.fuel.level.percent'",
					"tf-test inherited account_id",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMetricPruningRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
				),
			},
		},
	})
}

func TestAccNewRelicMetricPruningRule_InvalidNRQL(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMetricPruningRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicMetricPruningRuleConfig("FROM Metric DONT WORK rider_id", "invalid nrql"),
				ExpectError: regexp.MustCompile(`INVALID_QUERY`),
			},
		},
	})
}

func testAccCheckNewRelicMetricPruningRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		accountID, ruleID, err := parseNRQLDropRuleIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		rules, err := client.Pruningrules.GetListWithContext(context.Background(), accountID)
		if err != nil {
			return err
		}

		for _, r := range rules.Rules {
			if r.ID == ruleID {
				return nil
			}
		}

		return fmt.Errorf("metric pruning rule %s not found", ruleID)
	}
}

func testAccCheckNewRelicMetricPruningRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_metric_pruning_rule" {
			continue
		}

		accountID, ruleID, err := parseNRQLDropRuleIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		rules, err := client.Pruningrules.GetListWithContext(context.Background(), accountID)
		if err != nil {
			return err
		}

		for _, rule := range rules.Rules {
			if rule.ID == ruleID {
				return fmt.Errorf("metric pruning rule %s still exists", ruleID)
			}
		}
	}

	return nil
}

func testAccNewRelicMetricPruningRuleConfig(nrql, description string) string {
	return fmt.Sprintf(`
resource "newrelic_metric_pruning_rule" "foo" {
  account_id  = %d
  nrql        = %q
  description = %q
}
`, testAccountID, nrql, description)
}

func testAccNewRelicMetricPruningRuleNoDescriptionConfig(nrql string) string {
	return fmt.Sprintf(`
resource "newrelic_metric_pruning_rule" "foo" {
  account_id = %d
  nrql       = %q
}
`, testAccountID, nrql)
}

func testAccNewRelicMetricPruningRuleNoAccountIDConfig(nrql, description string) string {
	return fmt.Sprintf(`
resource "newrelic_metric_pruning_rule" "foo" {
  nrql        = %q
  description = %q
}
`, nrql, description)
}
