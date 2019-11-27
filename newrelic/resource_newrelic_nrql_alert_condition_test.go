package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicNrqlAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNrqlAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "term.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "term.0.duration", "5"),
					resource.TestCheckResourceAttr(resourceName, "term.0.operator", "below"),
					resource.TestCheckResourceAttr(resourceName, "term.0.priority", "critical"),
					resource.TestCheckResourceAttr(resourceName, "term.0.threshold", "0.75"),
					resource.TestCheckResourceAttr(resourceName, "term.0.time_function", "all"),
					resource.TestCheckResourceAttr(resourceName, "nrql.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.query", "SELECT uniqueCount(hostname) FROM ComputeSample"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.since_value", "20"),
				),
			},
			{
				Config: testAccNewRelicNrqlAlertConditionConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "term.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "term.0.duration", "10"),
					resource.TestCheckResourceAttr(resourceName, "term.0.operator", "below"),
					resource.TestCheckResourceAttr(resourceName, "term.0.priority", "critical"),
					resource.TestCheckResourceAttr(resourceName, "term.0.threshold", "0.65"),
					resource.TestCheckResourceAttr(resourceName, "term.0.time_function", "all"),
					resource.TestCheckResourceAttr(resourceName, "nrql.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.query", "SELECT uniqueCount(hostname) as Hosts FROM ComputeSample"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.since_value", "3"),
				),
			},
			{
				Config: testAccNewRelicNrqlAlertConditionConfigUpdatedWithStatic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "term.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "term.0.duration", "10"),
					resource.TestCheckResourceAttr(resourceName, "term.0.operator", "below"),
					resource.TestCheckResourceAttr(resourceName, "term.0.priority", "critical"),
					resource.TestCheckResourceAttr(resourceName, "term.0.threshold", "0.65"),
					resource.TestCheckResourceAttr(resourceName, "term.0.time_function", "all"),
					resource.TestCheckResourceAttr(resourceName, "nrql.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.query", "SELECT uniqueCount(hostname) as Hosts FROM ComputeSample"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.since_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "type", "static"),
				),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_Outlier(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNrqlAlertConditionConfigWithOutlier(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "term.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "term.0.duration", "10"),
					resource.TestCheckResourceAttr(resourceName, "term.0.operator", "above"),
					resource.TestCheckResourceAttr(resourceName, "term.0.priority", "critical"),
					resource.TestCheckResourceAttr(resourceName, "term.0.threshold", "0.65"),
					resource.TestCheckResourceAttr(resourceName, "term.0.time_function", "all"),
					resource.TestCheckResourceAttr(resourceName, "nrql.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.query", "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip"),
					resource.TestCheckResourceAttr(resourceName, "nrql.0.since_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "type", "outlier"),
					resource.TestCheckResourceAttr(resourceName, "expected_groups", "2"),
					resource.TestCheckResourceAttr(resourceName, "ignore_overlap", "true"),
				),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNrqlAlertConditionConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccNewRelicNrqlAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicNrqlAlertConditionExists("newrelic_nrql_alert_condition.foo"),
			},
		},
	})
}

// TODO: func_ TestAccNewRelicNrqlAlertCondition_Multi(t *testing.T) {

func testAccCheckNewRelicNrqlAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
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
			return fmt.Errorf("NRQL Alert condition still exists") //nolint:golint
		}

	}
	return nil
}

func testAccCheckNewRelicNrqlAlertConditionExists(n string) resource.TestCheckFunc {
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

		found, err := client.GetAlertNrqlCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccNewRelicNrqlAlertConditionConfig(rName string) string {
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
    since_value   = "20"
  }
  value_function  = "single_value"
}
`, rName)
}

func testAccNewRelicNrqlAlertConditionConfigUpdated(rName string) string {
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
}
`, rName)
}

func testAccNewRelicNrqlAlertConditionConfigUpdatedWithStatic(rName string) string {
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
  type = "static"
}
`, rName)
}

func testAccNewRelicNrqlAlertConditionConfigWithOutlier(rName string) string {
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
    operator      = "above"
    priority      = "critical"
    threshold     = "0.65"
    time_function = "all"
  }
  nrql {
    query         = "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip"
    since_value   = "3"
  }
  type            = "outlier"
  expected_groups = 2
  ignore_overlap  = true
}
`, rName)
}

// TODO: const testAccCheckNewRelicNrqlAlertConditionConfigMulti = `
