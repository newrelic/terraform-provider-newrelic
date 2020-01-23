package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameUpdated := fmt.Sprintf("tf-test-updated-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", rName),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "type", "apm_app_metric"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "metric", "apdex"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "condition_scope", "application"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1025554152.duration", "5"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1025554152.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1025554152.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1025554152.threshold", "0.75"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1025554152.time_function", "all"),
				),
			},
			// Test: Check no diff on re-apply
			{
				Config:             testAccNewRelicAlertConditionConfig(rName),
				ExpectNonEmptyPlan: false,
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertConditionConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", rNameUpdated),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "type", "apm_app_metric"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "metric", "error_percentage"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "condition_scope", "application"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.operator", "above"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.threshold", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.time_function", "any"),
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

func TestAccNewRelicAlertCondition_ZeroThreshold(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionConfigThreshold(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_AlertPolicyNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(rName),
				Config:    testAccNewRelicAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_ShortTermDuration(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfigDuration(rName, 4),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_LongTermDuration(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfigDuration(rName, 121),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_LongName(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfig("really-long-name-longer-than-sixty-four-characters-so-it-causes-an-error"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_EmptyName(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`name must not be empty`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionConfig(""),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

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

func testAccNewRelicAlertConditionConfigThreshold(name string, threshold int) string {
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
		threshold     = "%d"
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
