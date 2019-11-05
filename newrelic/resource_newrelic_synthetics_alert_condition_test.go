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
						"newrelic_synthetics_alert_condition.foo", "runbook_url", "www.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckNewRelicSyntheticsAlertConditionUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsAlertConditionExists("newrelic_synthetics_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "runbook_url", "www.example2.com"),
					resource.TestCheckResourceAttr(
						"newrelic_synthetics_alert_condition.foo", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckNewRelicSyntheticsAlertConditionConfig(rName),
			},
			resource.TestStep{
				PreConfig: deletePolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccCheckNewRelicSyntheticsAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicSyntheticsAlertConditionExists("newrelic_synthetics_alert_condition.foo"),
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

resource "newrelic_synthetics_monitor" "bar" {
	name = "tf-test-synthetic-%[1]s"
	type = "SIMPLE"
	frequency = 15
	status = "DISABLED"
	locations = ["AWS_US_EAST_1"]
	uri = "https://google.com"
}

resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_alert_condition" "foo" {
	policy_id = "${newrelic_alert_policy.foo.id}"
	name            = "tf-test-%[1]s"
	monitor_id      = "${newrelic_synthetics_monitor.bar.id}"
	runbook_url     = "www.example.com"
	enabled			= "true"
}
`, rName)
}

func testAccCheckNewRelicSyntheticsAlertConditionUpdated(rName string) string {
	return fmt.Sprintf(`

resource "newrelic_synthetics_monitor" "bar" {
	name = "tf-test-synthetic-%[1]s"
	type = "SIMPLE"
	frequency = 15
	status = "DISABLED"
	locations = ["AWS_US_EAST_1"]
	uri = "https://google.com"
}

resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_alert_condition" "foo" {
	policy_id = "${newrelic_alert_policy.foo.id}"
	name            = "tf-test-%[1]s"
	monitor_id      = "${newrelic_synthetics_monitor.bar.id}"
	runbook_url     = "www.example2.com"
	enabled			= "false"
}
`, rName)
}
