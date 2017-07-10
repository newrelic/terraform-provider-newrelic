package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	newrelic "github.com/paultyng/go-newrelic/api"
)

func TestAccNewRelicNrqlAlertCondition_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckNewRelicNrqlAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists("newrelic_nrql_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.duration", "5"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.threshold", "0.75"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.time_function", "all"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "nrql.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "nrql.0.query", "SELECT uniqueCount(hostname) FROM ComputeSample"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "nrql.0.since_value", "5"),
				),
			},
			resource.TestStep{
				Config: testAccCheckNewRelicNrqlAlertConditionConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists("newrelic_nrql_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.operator", "below"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.threshold", "0.65"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "term.0.time_function", "all"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "nrql.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "nrql.0.query", "SELECT uniqueCount(hostname) as Hosts FROM ComputeSample"),
					resource.TestCheckResourceAttr(
						"newrelic_nrql_alert_condition.foo", "nrql.0.since_value", "3"),
				),
			},
		},
	})
}

// TODO: func TestAccNewRelicNrqlAlertCondition_Multi(t *testing.T) {

func testAccCheckNewRelicNrqlAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*newrelic.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_nrql_alert_condition" {
			continue
		}

		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		_, err = client.GetAlertNrqlCondition(policyID, id)
		if err == nil {
			return fmt.Errorf("NRQL Alert condition still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicNrqlAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No alert condition ID is set")
		}

		client := testAccProvider.Meta().(*newrelic.Client)

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		found, err := client.GetAlertNrqlCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("Alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccCheckNewRelicNrqlAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-%[1]s"
  runbook_url     = "https://foo.example.com"
  enabled         = false

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
  nrql {
    query         = "SELECT uniqueCount(hostname) FROM ComputeSample"
    since_value   = "5"
  }
  value_function  = "single_value"
}
`, rName)
}

func testAccCheckNewRelicNrqlAlertConditionConfigUpdated(rName string) string {
	return fmt.Sprintf(`

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-updated-%[1]s"
  runbook_url     = "https://bar.example.com"
  enabled         = false

  term {
    duration      = 10
    operator      = "below"
    priority      = "critical"
    threshold     = "0.65"
    time_function = "all"
  }
  nrql {
    query         = "SELECT uniqueCount(hostname) as Hosts FROM ComputeSample"
    since_value   = "3"
  }
  value_function  = "single_value"
}
`, rName)
}

// TODO: const testAccCheckNewRelicNrqlAlertConditionConfigMulti = `
