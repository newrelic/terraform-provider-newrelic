package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicPluginsAlertCondition_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPluginsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicPluginsAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPluginsAlertConditionExists("newrelic_plugins_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1025554152.duration", "5"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1025554152.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1025554152.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1025554152.threshold", "0.75"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1025554152.time_function", "all"),
				),
			},
			{
				Config: testAccCheckNewRelicPluginsAlertConditionConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPluginsAlertConditionExists("newrelic_plugins_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1944209821.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1944209821.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1944209821.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1944209821.threshold", "0.65"),
					resource.TestCheckResourceAttr(
						"newrelic_plugins_alert_condition.foo", "term.1944209821.time_function", "all"),
				),
			},
		},
	})
}

func TestAccNewRelicPluginsAlertCondition(t *testing.T) {
	resourceName := "newrelic_plugins_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPluginsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicPluginsAlertConditionConfig(rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicPluginsAlertCondition_nameGreaterThan64Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicPluginsAlertConditionConfig("really-long-name-longer-than-sixty-four-characters-so-it-causes-an-error"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicPluginsAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicPluginsAlertConditionConfig(rName),
			},
			{
				PreConfig: deletePolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccCheckNewRelicPluginsAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicPluginsAlertConditionExists("newrelic_plugins_alert_condition.foo"),
			},
		},
	})
}

func testAccCheckNewRelicPluginsAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
	for _, r := range s.RootModule().Resources {
        if r.Type != "newrelic_plugins_alert_condition" {
           continue
        }
		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		_, err = client.GetAlertPluginsCondition(policyID, id)
		if err == nil {
			return fmt.Errorf("alert condition still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicPluginsAlertConditionExists(n string) resource.TestCheckFunc {
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

		found, err := client.GetAlertPluginsCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func TestErrorThrownUponPluginsConditionNameGreaterThan64Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponPluginsConditionNameGreaterThan64Char(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccCheckNewRelicPluginsAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name               = "tf-test-%[1]s"
  enabled            = false
  entities           = ["${data.newrelic_application.app.id}"]
  metric             = "my-metric"
  runbook_url        = "https://foo.example.com"
  metric_description = "my-metric-description"
  plugin_id          = "123"
  plugin_guid        = "com.example.plugin"
  value_function     = "average"

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

func testAccCheckNewRelicPluginsAlertConditionConfigUpdated(rName string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-updated-%[1]s"
}

resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name               = "tf-test-updated-%[1]s"
  enabled            = true
  entities           = ["${data.newrelic_application.app.id}"]
  runbook_url        = "https://bar.example.com"
  metric             = "my-metric"
  metric_description = "my-metric-description"
  plugin_id          = "123"
  plugin_guid        = "com.example.plugin"
  value_function     = "average"

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

func testErrorThrownUponPluginsConditionNameGreaterThan64Char(resourceName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  api_key = "foo"
}
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name               = "really-long-name-that-is-more-than-sixtyfour-characters-long-tf-test-%[1]s"
  enabled            = false
  entities           = ["12345"]
  runbook_url        = "https://foo.example.com"
  metric             = "my-metric"
  metric_description = "my-metric-description"
  plugin_id          = "123"
  plugin_guid        = "com.example.plugin"
  value_function     = "average"

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

func TestErrorThrownUponPluginsConditionNameLessThan1Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponPluginsConditionNameLessThan1Char(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testErrorThrownUponPluginsConditionNameLessThan1Char() string {
	return `
provider "newrelic" {
  api_key = "foo"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name               = ""
  enabled            = false
  entities           = ["12345"]
  runbook_url        = "https://foo.example.com"
  metric             = "my-metric"
  metric_description = "my-metric-description"
  plugin_id          = "123"
  plugin_guid        = "com.example.plugin"
  value_function     = "average"

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

func TestErrorThrownUponPluginsConditionTermDurationGreaterThan120(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponPluginsConditionTermDurationGreaterThan120(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testErrorThrownUponPluginsConditionTermDurationGreaterThan120() string {
	return `
provider "newrelic" {
  api_key = "foo"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name               = "tf-test-%[1]s"
  enabled            = false
  entities           = ["12345"]
  runbook_url        = "https://foo.example.com"
  metric             = "my-metric"
  metric_description = "my-metric-description"
  plugin_id          = "123"
  plugin_guid        = "com.example.plugin"
  value_function     = "average"

  term {
    duration      = 121
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
`
}

func TestAccNewRelicAlertPluginsCondition_TermDurationLessThan5(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected term.0.duration to be in the range \(5 - 120\)`)
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      TestAccNewRelicAlertPluginsCondition_TermDurationLessThan5(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertPluginsCondition_TermDurationLessThan5() string {
	return `
provider "newrelic" {
  api_key = "foo"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name               = "tf-test-%[1]s"
  enabled            = false
  entities           = ["12345"]
  runbook_url        = "https://foo.example.com"
  metric             = "my-metric"
  metric_description = "my-metric-description"
  plugin_id          = "123"
  plugin_guid        = "com.example.plugin"
  value_function     = "average"

  term {
    duration      = 4
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
`
}
