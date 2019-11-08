package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicAlertCondition_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "type", "apm_app_metric"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
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
			{
				Config: testAccCheckNewRelicAlertConditionConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1944209821.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1944209821.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1944209821.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1944209821.threshold", "0.65"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.1944209821.time_function", "all"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_ZeroThreshold(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertConditionConfigZeroThreshold(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "type", "apm_app_metric"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.971858588.duration", "5"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.971858588.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.971858588.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.971858588.threshold", "0"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.971858588.time_function", "all"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertCondition(t *testing.T) {
	resourceName := "newrelic_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertConditionConfig(rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_nameGreaterThan64Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicAlertConditionConfig("really-long-name-longer-than-sixty-four-characters-so-it-causes-an-error"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertConditionConfig(rName),
			},
			{
				PreConfig: deletePolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccCheckNewRelicAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
			},
		},
	})
}

func testAccCheckNewRelicAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
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

		_, err = client.GetAlertCondition(policyID, id)
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

		client := testAccProvider.Meta().(*ProviderConfig).Client

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		found, err := client.GetAlertCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func TestErrorThrownUponConditionNameGreaterThan64Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponConditionNameGreaterThan64Char(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccCheckNewRelicAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-%[1]s"
  enabled         = false
  type            = "apm_app_metric"
  entities        = ["${data.newrelic_application.app.id}"]
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
`, rName, testAccExpectedApplicationName)
}

func testAccCheckNewRelicAlertConditionConfigUpdated(rName string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-updated-%[1]s"
}

resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-updated-%[1]s"
  enabled         = true  
  type            = "apm_app_metric"
  entities        = ["${data.newrelic_application.app.id}"]
  metric          = "apdex"
  runbook_url     = "https://bar.example.com"
  condition_scope = "application"

  term {
    duration      = 10
    operator      = "below"
    priority      = "critical"
    threshold     = "0.65"
    time_function = "all"
  }
}
`, rName, testAccExpectedApplicationName)
}

func testAccCheckNewRelicAlertConditionConfigZeroThreshold(rName string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-%[1]s"
  enabled         = false  
  type            = "apm_app_metric"
  entities        = ["${data.newrelic_application.app.id}"]
  metric          = "apdex"
  runbook_url     = "https://foo.example.com"
  condition_scope = "application"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0"
    time_function = "all"
  }
}
`, rName, testAccExpectedApplicationName)
}

func testErrorThrownUponConditionNameGreaterThan64Char(resourceName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  api_key = "foo"
}
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"
  name            = "really-long-name-that-is-more-than-sixtyfour-characters-long-tf-test-%[1]s"
  type            = "apm_app_metric"
  entities        = ["12345"]
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
`, resourceName, testAccExpectedApplicationName)
}

func TestErrorThrownUponConditionNameLessThan1Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponConditionNameLessThan1Char(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testErrorThrownUponConditionNameLessThan1Char() string {
	return `
provider "newrelic" {
  api_key = "foo"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"
  name            = ""
  type            = "apm_app_metric"
  entities        = ["12345"]
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
`
}

// TODO: const testAccCheckNewRelicAlertConditionConfigMulti = `
