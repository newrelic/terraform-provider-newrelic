//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicServiceLevelAlertHelper_FastBurn(t *testing.T) {
	resourceName := "data.newrelic_service_level_alert_helper.fast"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicServiceLevelAlertHelperFastBurnConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelAlertHelper_FastBurn(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_SlowBurn(t *testing.T) {
	resourceName := "data.newrelic_service_level_alert_helper.slow"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicServiceLevelAlertHelperSlowBurnConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelAlertHelper_SlowBurn(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_Custom(t *testing.T) {
	resourceName := "data.newrelic_service_level_alert_helper.custom"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicServiceLevelAlertHelperCustomConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelAlertHelper_Custom(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_FastBurnError(t *testing.T) {
	expectedErrorMessage := regexp.MustCompile(`For 'fast_burn' alert type do not fill 'custom_evaluation_period' or 'custom_tolerated_budget_consumption'.`)

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicServiceLevelAlertHelperFastBurnEvaluationErrorConfig(),
				ExpectError: expectedErrorMessage,
			},
			{
				Config:      testAccNewRelicServiceLevelAlertHelperFastBurnBudgetErrorConfig(),
				ExpectError: expectedErrorMessage,
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_SlowBurnError(t *testing.T) {
	expectedErrorMessage := regexp.MustCompile(`For 'slow_burn' alert type do not fill 'custom_evaluation_period' or 'custom_tolerated_budget_consumption'.`)

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicServiceLevelAlertHelperSlowBurnEvaluationErrorConfig(),
				ExpectError: expectedErrorMessage,
			},
			{
				Config:      testAccNewRelicServiceLevelAlertHelperSlowBurnBudgetErrorConfig(),
				ExpectError: expectedErrorMessage,
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_CustomError(t *testing.T) {
	expectedErrorMessage := regexp.MustCompile(`For 'custom' alert type the fields 'custom_evaluation_period' and 'custom_tolerated_budget_consumption' are mandatory.`)

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicServiceLevelAlertHelperCustomEvaluationErrorConfig(),
				ExpectError: expectedErrorMessage,
			},
			{
				Config:      testAccNewRelicServiceLevelAlertHelperCustomBudgetErrorConfig(),
				ExpectError: expectedErrorMessage,
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_CustomBadEvents(t *testing.T) {
	resourceName := "data.newrelic_service_level_alert_helper.custom"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicServiceLevelAlertHelperCustomBadEventsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicServiceLevelAlertHelper_CustomBadEvents(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicServiceLevelAlertHelperSlowBurnConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "slow" {
    alert_type = "slow_burn"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
}
`)
}

func testAccNewRelicServiceLevelAlertHelperFastBurnConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "fast" {
    alert_type = "fast_burn"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
}
`)
}

func testAccNewRelicServiceLevelAlertHelperFastBurnEvaluationErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "fastBad" {
    alert_type = "fast_burn"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
    custom_evaluation_period = 12
}
`)
}

func testAccNewRelicServiceLevelAlertHelperSlowBurnEvaluationErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "slowBad" {
    alert_type = "slow_burn"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
    custom_evaluation_period = 12
}
`)
}

func testAccNewRelicServiceLevelAlertHelperFastBurnBudgetErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "fastBad" {
    alert_type = "fast_burn"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
    custom_tolerated_budget_consumption = 34
}
`)
}

func testAccNewRelicServiceLevelAlertHelperSlowBurnBudgetErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "slowBad" {
    alert_type = "slow_burn"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
    custom_tolerated_budget_consumption = 34
}
`)
}

func testAccNewRelicServiceLevelAlertHelperCustomEvaluationErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "customBad" {
    alert_type = "custom"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
    custom_evaluation_period = 12
}
`)
}

func testAccNewRelicServiceLevelAlertHelperCustomBudgetErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "customBad" {
    alert_type = "custom"
    sli_guid = "sliGuid"
    slo_target = 99.9
    slo_period = 28
    custom_tolerated_budget_consumption = 34
}
`)
}

func testAccCheckNewRelicServiceLevelAlertHelper_FastBurn(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		a := rs.Primary.Attributes

		testCases := map[string]string{
			"slo_period":                          "28",
			"slo_target":                          "99.9",
			"alert_type":                          "fast_burn",
			"custom_evaluation_period":            "",
			"custom_tolerated_budget_consumption": "",
			"evaluation_period":                   "60",
			"tolerated_budget_consumption":        "2",
			"threshold":                           "1.3439999999999237",
			"sli_guid":                            "sliGuid",
			"nrql":                                "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = 'sliGuid'",
		}

		for attrName, expectedVal := range testCases {
			if err := runTest(a, attrName, expectedVal); err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccCheckNewRelicServiceLevelAlertHelper_SlowBurn(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		a := rs.Primary.Attributes

		testCases := map[string]string{
			"slo_period":                          "28",
			"slo_target":                          "99.9",
			"alert_type":                          "slow_burn",
			"custom_evaluation_period":            "",
			"custom_tolerated_budget_consumption": "",
			"evaluation_period":                   "360",
			"tolerated_budget_consumption":        "5",
			"threshold":                           "0.5599999999999682",
			"sli_guid":                            "sliGuid",
			"nrql":                                "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = 'sliGuid'",
		}

		for attrName, expectedVal := range testCases {
			if err := runTest(a, attrName, expectedVal); err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccNewRelicServiceLevelAlertHelperCustomConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "custom" {
    alert_type = "custom"
    sli_guid = "sliGuidCustom"
    slo_target = 98
    slo_period = 7
    custom_tolerated_budget_consumption = 5
    custom_evaluation_period = 120
}
`)
}

func testAccCheckNewRelicServiceLevelAlertHelper_Custom(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		a := rs.Primary.Attributes

		testCases := map[string]string{
			"slo_period":                          "7",
			"slo_target":                          "98",
			"alert_type":                          "custom",
			"custom_evaluation_period":            "120",
			"custom_tolerated_budget_consumption": "5",
			"evaluation_period":                   "120",
			"tolerated_budget_consumption":        "5",
			"threshold":                           "8.4",
			"sli_guid":                            "sliGuidCustom",
			"nrql":                                "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = 'sliGuidCustom'",
		}

		for attrName, expectedVal := range testCases {
			if err := runTest(a, attrName, expectedVal); err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccNewRelicServiceLevelAlertHelperCustomBadEventsConfig() string {
	return fmt.Sprintf(`
data "newrelic_service_level_alert_helper" "custom" {
    alert_type = "custom"
    sli_guid = "sliGuidCustom"
    slo_target = 98
    slo_period = 7
    custom_tolerated_budget_consumption = 5
    custom_evaluation_period = 120
    is_bad_events = true
}
`)
}

func testAccCheckNewRelicServiceLevelAlertHelper_CustomBadEvents(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		a := rs.Primary.Attributes

		testCases := map[string]string{
			"slo_period":                          "7",
			"slo_target":                          "98",
			"alert_type":                          "custom",
			"custom_evaluation_period":            "120",
			"custom_tolerated_budget_consumption": "5",
			"evaluation_period":                   "120",
			"tolerated_budget_consumption":        "5",
			"threshold":                           "8.4",
			"sli_guid":                            "sliGuidCustom",
			"nrql":                                "FROM Metric SELECT 100 - clamp_max((sum(newrelic.sli.valid) - sum(newrelic.sli.bad)) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance' WHERE sli.guid = 'sliGuidCustom'",
		}

		for attrName, expectedVal := range testCases {
			if err := runTest(a, attrName, expectedVal); err != nil {
				return err
			}
		}

		return nil
	}
}

func runTest(attributes map[string]string, name string, expected string) error {
	actual := attributes[name]
	if actual != expected {
		return fmt.Errorf("Expected %s was %s, actual was %s", name, expected, actual)
	}
	return nil
}
