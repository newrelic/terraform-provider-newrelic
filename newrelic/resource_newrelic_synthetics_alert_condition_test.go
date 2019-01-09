package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicSyntheticsAlertCondition_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicSyntheticsAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsAlertConditionExists("newrelic_synthetics_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "policy_id", "0"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "monitor_id", "derp"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "runbook_url", "www.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_alert_condition" {
			continue
		}

		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		_, err = client.GetAlertSyntheticsCondition(policyID, id)
		if err == nil {
			return fmt.Errorf("Synthetics alert condition still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicSyntheticsAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No alert condition ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Client

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		found, err := client.GetAlertSyntheticsCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("Alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`

data "newrelic_synthetics_monitor" "bar" {
  name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

	name            = "tf-test-%[1]s"
	monitor_id      = "${data.newrelic_synthetics_monitor.bar.id}"
	runbook_url     = "https://foo.example.com"
}
`, rName, testAccExpectedApplicationName)
}

func testAccCheckNewRelicSyntheticsAlertConditionConfigUpdated(rName string) string {
	return fmt.Sprintf(`

data "newrelic_synthetics_monitor" "bar" {
  name = "%[2]s"
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

	name            = "tf-test-updated-%[1]s"
	monitor_id      = "${data.newrelic_synthetics_monitor.bar.id}"
  runbook_url     = "https://bar.example.com"
}
`, rName, testAccExpectedApplicationName)
}
