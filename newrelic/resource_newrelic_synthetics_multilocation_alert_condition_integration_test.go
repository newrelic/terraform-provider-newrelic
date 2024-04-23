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

func TestAccNewRelicSyntheticsMultiLocationAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_multilocation_alert_condition.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMultiLocationAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMultiLocationConditionConfigBasic(rName, "1", "2", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMultiLocationAlertConditionExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMultiLocationConditionConfigBasic(rName, "11", "12", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMultiLocationAlertConditionExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entities"},
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsMultiLocationAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no alert condition ID is set")
		}

		var err error

		ids, err := parseHashedIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		conditionID := ids[1]
		policyID := ids[0]

		found, err := client.Alerts.GetMultiLocationSyntheticsCondition(policyID, conditionID)
		if err != nil {
			return err
		}

		if found.ID != conditionID {
			return fmt.Errorf("synthetics multi-location alert condition not found: %v - %v", conditionID, found)
		}

		return nil
	}
}

func testAccNewRelicSyntheticsMultiLocationConditionConfigBasic(
	name string,
	criticalThreshold string,
	warningThreshold string,
	conditionalAttrs string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "policy" {
  name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_monitor" "monitor" {
  locations_public          = ["US_WEST_1"]
  name                      = "tf-test-%[1]s"
  period                    = "EVERY_10_MINUTES"
  status                    = "DISABLED"
  type                      = "SIMPLE"
  uri = "https://www.one.newrelic.com"
}

resource "newrelic_synthetics_multilocation_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.policy.id

  name                         = "tf-test-%[1]s"
  runbook_url                  = "https://foo.example.com"
  enabled                      = true
  violation_time_limit_seconds = "3600"

	entities = [
		newrelic_synthetics_monitor.monitor.id
	]

	critical {
    threshold = %[2]s
	}

	warning {
    threshold = %[3]s
	}

	%[4]s
}
`, name, criticalThreshold, warningThreshold, conditionalAttrs)
}

func testAccCheckNewRelicMultiLocationAlertConditionDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_nrql_alert_condition" {
			continue
		}

		var err error

		ids, err := parseHashedIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		conditionID := ids[1]
		policyID := ids[0]

		if _, err = client.Alerts.GetMultiLocationSyntheticsCondition(policyID, conditionID); err == nil {
			return fmt.Errorf("Synthetics multi-location condition still exists") //nolint:golint
		}
	}

	return nil
}

func TestAccNewRelicSyntheticsMultiLocationAlertCondition_InvalidEntities(t *testing.T) {
	resourceName := "newrelic_synthetics_multilocation_alert_condition.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicMultiLocationAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicSyntheticsMultiLocationConditionConfigInvalidEntities(rName),
				ExpectError: regexp.MustCompile(`invalid monitor ID`),
			},
		},
	})
}

func testAccNewRelicSyntheticsMultiLocationConditionConfigInvalidEntities(
	name string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "policy" {
  name = "tf-test-%[1]s"
}

resource "newrelic_synthetics_monitor" "monitor" {
  locations_public          = ["US_WEST_1"]
  name                      = "tf-test-%[1]s"
  period                    = "EVERY_10_MINUTES"
  status                    = "DISABLED"
  type                      = "SIMPLE"
  uri = "https://www.one.newrelic.com"
}

resource "newrelic_synthetics_multilocation_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.policy.id

  name                         = "tf-test-%[1]s"
  runbook_url                  = "https://foo.example.com"
  enabled                      = true
  violation_time_limit_seconds = "3600"

	entities = "1234" // Providing an invalid value here

	critical {
    threshold = "1.0"
	}

	warning {
    threshold = "0.5"
	}
}
`, name)
}
