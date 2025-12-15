//go:build integration || ALERTS
// +build integration ALERTS

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAlertCompoundCondition_Basic(t *testing.T) {
	resourceName := "newrelic_alert_compound_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create with AND expression
			{
				Config: testAccNewRelicAlertCompoundConditionConfigBasic(rName, "A AND B"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertCompoundConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "trigger_expression", "A AND B"),
					resource.TestCheckResourceAttr(resourceName, "component_conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "facet_matching_behavior", "FACETS_IGNORED"),
					resource.TestCheckResourceAttr(resourceName, "runbook_url", "https://example.com/runbook"),
					resource.TestCheckResourceAttr(resourceName, "threshold_duration", "120"),
				),
			},
			// Test: Update to OR expression
			{
				Config: testAccNewRelicAlertCompoundConditionConfigBasic(rName, "A OR B"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertCompoundConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "trigger_expression", "A OR B"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_ThreeComponents(t *testing.T) {
	resourceName := "newrelic_alert_compound_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertCompoundConditionConfigThreeComponents(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertCompoundConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "trigger_expression", "(A AND B) OR C"),
					resource.TestCheckResourceAttr(resourceName, "component_conditions.#", "3"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_FacetMatching(t *testing.T) {
	resourceName := "newrelic_alert_compound_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertCompoundConditionConfigWithFacetMatch(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertCompoundConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "facet_matching_behavior", "FACETS_MATCH"),
				),
			},
		},
	})
}

func testAccNewRelicAlertCompoundConditionConfigBasic(name, expression string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 5.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "condition_b" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-b-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT average(duration) FROM Transaction WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 1.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_alert_compound_condition" "foo" {
	account_id         = %[3]d
	policy_id          = newrelic_alert_policy.foo.id
	name               = "tf-test-%[1]s"
	enabled            = true
	trigger_expression = "%[2]s"
	runbook_url        = "https://example.com/runbook"
	threshold_duration = 120

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_a.id)[1]
		alias = "A"
	}

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_b.id)[1]
		alias = "B"
	}

	facet_matching_behavior = "FACETS_IGNORED"
}
`, name, expression, testAccountID)
}

func testAccNewRelicAlertCompoundConditionConfigThreeComponents(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 5.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "condition_b" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-b-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT average(duration) FROM Transaction WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 1.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "condition_c" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-c-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT percentage(count(*), WHERE error IS true) FROM Transaction WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 10.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_alert_compound_condition" "foo" {
	account_id         = %[2]d
	policy_id          = newrelic_alert_policy.foo.id
	name               = "tf-test-%[1]s"
	enabled            = true
	trigger_expression = "(A AND B) OR C"

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_a.id)[1]
		alias = "A"
	}

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_b.id)[1]
		alias = "B"
	}

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_c.id)[1]
		alias = "C"
	}
}
`, name, testAccountID)
}

func testAccNewRelicAlertCompoundConditionConfigWithFacetMatch(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction FACET appName WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 5.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "condition_b" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-b-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT average(duration) FROM Transaction FACET appName WHERE appName = 'Dummy App'"
	}

	critical {
		operator              = "above"
		threshold             = 1.0
		threshold_duration    = 300
		threshold_occurrences = "all"
	}

	violation_time_limit_seconds = 3600
}

resource "newrelic_alert_compound_condition" "foo" {
	account_id              = %[2]d
	policy_id               = newrelic_alert_policy.foo.id
	name                    = "tf-test-%[1]s"
	enabled                 = true
	trigger_expression      = "A AND B"
	facet_matching_behavior = "FACETS_MATCH"

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_a.id)[1]
		alias = "A"
	}

	component_conditions {
		id    = split(":", newrelic_nrql_alert_condition.condition_b.id)[1]
		alias = "B"
	}
}
`, name, testAccountID)
}
