//go:build unit
// +build unit

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAlertCompoundCondition_ThresholdDurationTooShort(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected threshold_duration to be in the range \(30 - 1440\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertCompoundConditionConfigInvalidThresholdDuration(rName, 29),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_ThresholdDurationTooLong(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected threshold_duration to be in the range \(30 - 1440\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertCompoundConditionConfigInvalidThresholdDuration(rName, 1441),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_InvalidAliasStartsWithNumber(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`alias must start with a letter`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertCompoundConditionConfigInvalidAlias(rName, "1A", "B"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_InvalidAliasSpecialCharacters(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`alias must start with a letter`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertCompoundConditionConfigInvalidAlias(rName, "A-B", "C"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_InvalidFacetMatchingBehavior(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected facet_matching_behavior to be one of \[FACETS_MATCH FACETS_IGNORED\]`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertCompoundConditionConfigInvalidFacetMatching(rName, "INVALID_VALUE"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCompoundCondition_TooFewComponentConditions(t *testing.T) {
	avoidEmptyAccountID()
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`At least 2 "component_conditions" blocks are required`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertCompoundConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertCompoundConditionConfigSingleComponent(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccNewRelicAlertCompoundConditionConfigInvalidThresholdDuration(name string, duration int) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction"
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
		query = "SELECT average(duration) FROM Transaction"
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
	trigger_expression = "A AND B"
	threshold_duration = %[2]d

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_a.id
		alias = "A"
	}

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_b.id
		alias = "B"
	}
}
`, name, duration, testAccountID)
}

func testAccNewRelicAlertCompoundConditionConfigInvalidAlias(name, aliasA, aliasB string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction"
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
		query = "SELECT average(duration) FROM Transaction"
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
	account_id         = %[4]d
	policy_id          = newrelic_alert_policy.foo.id
	name               = "tf-test-%[1]s"
	enabled            = true
	trigger_expression = "A AND B"

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_a.id
		alias = "%[2]s"
	}

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_b.id
		alias = "%[3]s"
	}
}
`, name, aliasA, aliasB, testAccountID)
}

func testAccNewRelicAlertCompoundConditionConfigInvalidFacetMatching(name, facetBehavior string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction"
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
		query = "SELECT average(duration) FROM Transaction"
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
	account_id              = %[3]d
	policy_id               = newrelic_alert_policy.foo.id
	name                    = "tf-test-%[1]s"
	enabled                 = true
	trigger_expression      = "A AND B"
	facet_matching_behavior = "%[2]s"

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_a.id
		alias = "A"
	}

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_b.id
		alias = "B"
	}
}
`, name, facetBehavior, testAccountID)
}

func testAccNewRelicAlertCompoundConditionConfigSingleComponent(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "condition_a" {
	policy_id = newrelic_alert_policy.foo.id
	name      = "tf-test-condition-a-%[1]s"
	enabled   = true

	nrql {
		query = "SELECT count(*) FROM Transaction"
	}

	critical {
		operator              = "above"
		threshold             = 5.0
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
	trigger_expression = "A"

	component_conditions {
		id    = newrelic_nrql_alert_condition.condition_a.id
		alias = "A"
	}
}
`, name, testAccountID)
}
