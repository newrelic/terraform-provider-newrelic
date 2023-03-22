//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"
	"regexp"

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
    expectedErrorMessage := regexp.MustCompile(`For fast_burn alert type do not fill 'custom_evaluation_period' or 'custom_tolerated_budget_consumption'.`)

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicServiceLevelAlertHelperFastBurnEvaluationErrorConfig(),
                ExpectError: expectedErrorMessage,
			},
			{
				Config: testAccNewRelicServiceLevelAlertHelperFastBurnBudgetErrorConfig(),
                ExpectError: expectedErrorMessage,
			},
		},
	})
}

func TestAccNewRelicServiceLevelAlertHelper_CustomError(t *testing.T) {
    expectedErrorMessage := regexp.MustCompile(`For custom alert type the fields 'custom_evaluation_period' and 'custom_tolerated_budget_consumption' are mandatory.`)

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicServiceLevelAlertHelperCustomEvaluationErrorConfig(),
                ExpectError: expectedErrorMessage,
			},
			{
				Config: testAccNewRelicServiceLevelAlertHelperCustomBudgetErrorConfig(),
                ExpectError: expectedErrorMessage,
			},
		},
	})
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
		// fmt.Println(s.RootModule())
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		a := rs.Primary.Attributes
	
        var err error
		if err = runTest(a, "slo_period", "28"); err != nil {
			return err
		}
		if err = runTest(a, "slo_target", "99.9"); err != nil {
			return err
		}
		if err = runTest(a, "alert_type", "fast_burn"); err != nil {
			return err
		}
		if err = runTest(a, "custom_evaluation_period", ""); err != nil {
			return err
		}
		if err = runTest(a, "custom_tolerated_budget_consumption", ""); err != nil {
			return err
		}
		if err = runTest(a, "evaluation_period", "60"); err != nil {
			return err
		}
		if err = runTest(a, "tolerated_budget_consumption", "2"); err != nil {
			return err
		}
		if err = runTest(a, "threshold", "1.3439999999999237"); err != nil {
			return err
		}
		if err = runTest(a, "sli_guid", "sliGuid"); err != nil {
			return err
		}
		if err = runTest(a, "nrql", "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = 'sliGuid'"); err != nil {
			return err
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
		// fmt.Println(s.RootModule())
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		a := rs.Primary.Attributes
	
        var err error
		if err = runTest(a, "slo_period", "7"); err != nil {
			return err
		}
		if err = runTest(a, "slo_target", "98"); err != nil {
			return err
		}
		if err = runTest(a, "alert_type", "custom"); err != nil {
			return err
		}
		if err = runTest(a, "custom_evaluation_period", "120"); err != nil {
			return err
		}
		if err = runTest(a, "custom_tolerated_budget_consumption", "5"); err != nil {
			return err
		}
		if err = runTest(a, "evaluation_period", "120"); err != nil {
			return err
		}
		if err = runTest(a, "tolerated_budget_consumption", "5"); err != nil {
			return err
		}
		if err = runTest(a, "threshold", "8.4"); err != nil {
			return err
		}
		if err = runTest(a, "sli_guid", "sliGuidCustom"); err != nil {
			return err
		}
		if err = runTest(a, "nrql", "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = 'sliGuidCustom'"); err != nil {
			return err
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
