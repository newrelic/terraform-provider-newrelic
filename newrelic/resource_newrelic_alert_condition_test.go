// +build integration unit
//
// Test helpers
//

package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckNewRelicAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_condition" {
			continue
		}

		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		_, err = client.Alerts.GetCondition(policyID, id)
		if err == nil {
			return fmt.Errorf("alert condition still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no alert condition ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		found, err := client.Alerts.GetCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccNewRelicAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	api_key = "%[3]s"
}
data "newrelic_application" "app" {
	name = "%[2]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = true
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.app.id]
	metric          = "apdex"
	runbook_url     = "https://foo.example.com"
	condition_scope = "application"

	term {
		duration      = 5
		operator      = "below"
		priority      = "critical"
		threshold     = "0.75"
		time_function = "all"
	}
}
`, rName, testAccExpectedApplicationName, testAccAPIKey)
}

func testAccNewRelicAlertConditionConfigUpdated(name string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = false
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.app.id]
	metric          = "error_percentage"
	runbook_url     = "https://bar.example.com"
	condition_scope = "application"

	term {
		duration      = 10
		operator      = "above"
		priority      = "critical"
		threshold     = "1.00"
		time_function = "any"
	}
}
`, name, testAccExpectedApplicationName)
}

func testAccNewRelicAlertConditionConfigThreshold(name string, threshold float64) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[3]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = false
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.app.id]
	metric          = "apdex"
	runbook_url     = "https://foo.example.com"
	condition_scope = "application"

	term {
		duration      = 5
		operator      = "below"
		priority      = "critical"
		threshold     = "%f"
		time_function = "all"
	}
}
`, name, threshold, testAccExpectedApplicationName)
}

func testAccNewRelicAlertConditionConfigDuration(name string, duration int) string {
	return fmt.Sprintf(`
provider "newrelic" {
	api_key = "foo"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id
	name            = "test-term-duration"
	type            = "apm_app_metric"
	entities        = ["12345"]
	metric          = "apdex"
	runbook_url     = "https://foo.example.com"
	condition_scope = "application"

	term {
		duration      = %[2]d
		operator      = "below"
		priority      = "critical"
		threshold     = "0.75"
		time_function = "all"
	}
}
`, name, duration)
}

func testAccNewRelicAlertConditionApplicationScopeWithCloseTimerConfig(rName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	api_key = "%[3]s"
}
data "newrelic_application" "app" {
	name = "%[2]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = true
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.app.id]
	metric          = "apdex"
	runbook_url     = "https://foo.example.com"
	condition_scope = "application"
	violation_close_timer = 24

	term {
		duration      = 5
		operator      = "below"
		priority      = "critical"
		threshold     = "0.75"
		time_function = "all"
	}
}
`, rName, testAccExpectedApplicationName, testAccAPIKey)
}

func testAccNewRelicAlertConditionInstanceScopeWithCloseTimerConfig(rName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	api_key = "%[3]s"
}
data "newrelic_application" "app" {
	name = "%[2]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = true
	type            = "apm_app_metric"
	entities        = [317250408]
	metric          = "apdex"
	runbook_url     = "https://foo.example.com"
	condition_scope = "instance"
	violation_close_timer = 24

	term {
		duration      = 5
		operator      = "below"
		priority      = "critical"
		threshold     = "0.75"
		time_function = "all"
	}
}
`, rName, testAccExpectedApplicationName, testAccAPIKey)
}

func testAccNewRelicAlertConditionAPMJVMMetricApplicationScopeConfig(rName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	api_key = "%[3]s"
}
data "newrelic_application" "app" {
	name = "%[2]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = true
	type            = "apm_jvm_metric"
	entities        = [317250408]
	metric          = "heap_memory_usage"
	runbook_url     = "https://foo.example.com"
	condition_scope = "application"
	violation_close_timer = 24

	term {
		duration      = 5
		operator      = "below"
		priority      = "critical"
		threshold     = "0.75"
		time_function = "all"
	}
}
`, rName, testAccExpectedApplicationName, testAccAPIKey)
}

func testAccNewRelicAlertConditionAPMJVMMetricInstanceScopeConfig(rName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	api_key = "%[3]s"
}
data "newrelic_application" "app" {
	name = "%[2]s"
}
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}
resource "newrelic_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name            = "%[1]s"
	enabled         = true
	type            = "apm_jvm_metric"
	entities        = [317250408]
	metric          = "heap_memory_usage"
	runbook_url     = "https://foo.example.com"
	condition_scope = "instance"
	violation_close_timer = 24

	term {
		duration      = 5
		operator      = "below"
		priority      = "critical"
		threshold     = "0.75"
		time_function = "all"
	}
}
`, rName, testAccExpectedApplicationName, testAccAPIKey)
}
