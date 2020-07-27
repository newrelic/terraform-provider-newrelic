// +build integration unit
//
// Test helpers
//

package newrelic

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckNewRelicAlertPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_policy" {
			continue
		}

		idParts, err := parseHashedIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		policyID := strconv.Itoa(idParts[0])

		_, err = client.Alerts.QueryPolicy(testAccountID, policyID)

		if err == nil {
			return fmt.Errorf("policy still exists: %s", err)
		}

	}
	return nil
}

func testAccCheckNewRelicAlertPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no policy ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		idParts, err := parseHashedIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		policyID := strconv.Itoa(idParts[0])

		found, err := client.Alerts.QueryPolicy(testAccountID, policyID)
		if err != nil {
			return err
		}

		if found.ID != policyID {
			return fmt.Errorf("policy not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccNewRelicAlertPolicyConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%s"
}
`, name)
}

func testAccNewRelicAlertPolicyConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name                = "tf-test-updated-%s"
  incident_preference = "PER_CONDITION"
}
`, rName)
}

func testAccNewRelicAlertPolicyConfigNameEmpty() string {
	return `
provider "newrelic" {
	api_key = "foo"
}

resource "newrelic_alert_policy" "foo" {
  name = ""
}
`
}

func testAccNewRelicAccAlertPolicyConfigWithChannels(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "channel_a" {
	name = "tf-test-%[1]s-channel-a"
	type = "email"

	config {
		recipients = "no-reply+a@newrelic.com"
		include_json_attachment = "1"
	}
}

resource "newrelic_alert_channel" "channel_b" {
	name = "tf-test-%[1]s-channel-b"
	type = "email"

	config {
		recipients = "no-reply+b@newrelic.com"
		include_json_attachment = "1"
	}
}

resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
	channel_ids =  [
		newrelic_alert_channel.channel_a.id,
		newrelic_alert_channel.channel_b.id
	]
}
`, name)
}
