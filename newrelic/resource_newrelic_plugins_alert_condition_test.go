// +build integration unit
//
// Test Helpers
//

package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckNewRelicPluginsAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
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

		_, err = client.Alerts.GetPluginsCondition(policyID, id)
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

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		found, err := client.Alerts.GetPluginsCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccCheckNewRelicPluginsAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`
data "newrelic_application" "app" {
	name = "%[2]s"
}
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name               = "tf-test-%[1]s"
	enabled            = false
	entities           = [data.newrelic_plugin_component.foo.id]
	metric             = "Component/Connection/Clients[connections]"
	runbook_url        = "https://foo.example.com"
	metric_description = "my-metric-description"
	plugin_id          = "21709"
	plugin_guid        = "net.kenjij.newrelic_redis_plugin"
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
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-updated-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name               = "tf-test-updated-%[1]s"
	enabled            = true
	entities           = [data.newrelic_plugin_component.foo.id]
	runbook_url        = "https://bar.example.com"
	metric             = "Component/Connection/Clients[connections]"
	metric_description = "my-metric-description"
	plugin_id          = data.newrelic_plugin.foo.id
	plugin_guid        = data.newrelic_plugin.foo.guid
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

func testAccNewRelicPluginsAlertConditionNameGreaterThan64Char(resourceName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  api_key = "foo"
}
data "newrelic_plugin" "foo" {
  guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name               = "really-long-name-that-is-more-than-sixtyfour-characters-long-tf-test-%[1]s"
  enabled            = false
  entities           = [data.newrelic_plugin_component.foo.id]
  runbook_url        = "https://foo.example.com"
  metric             = "Component/Connection/Clients[connections]"
  metric_description = "my-metric-description"
  plugin_id          = data.newrelic_plugin.foo.id
  plugin_guid        = data.newrelic_plugin.foo.guid
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

func testAccNewRelicPluginsAlertConditionNameLessThan1Char() string {
	return `
provider "newrelic" {
	api_key = "foo"
}
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id
	
	name               = ""
	enabled            = false
	entities           = [data.newrelic_plugin_component.foo.id]
	runbook_url        = "https://foo.example.com"
	metric             = "Component/Connection/Clients[connections]"
	metric_description = "my-metric-description"
	plugin_id          = data.newrelic_plugin.foo.id
	plugin_guid        = data.newrelic_plugin.foo.guid
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

func testAccNewRelicPluginsAlertConditionTermDurationGreaterThan120() string {
	return `
provider "newrelic" {
	api_key = "foo"
}
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name               = "tf-test-%[1]s"
	enabled            = false
	entities           = [data.newrelic_plugin_component.foo.id]
	runbook_url        = "https://foo.example.com"
	metric             = "Component/Connection/Clients[connections]"
	metric_description = "my-metric-description"
	plugin_id          = data.newrelic_plugin.foo.id
	plugin_guid        = data.newrelic_plugin.foo.guid
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

func testAccNewRelicPluginsAlertConditionTermDurationLessThan5() string {
	return `
provider "newrelic" {
	api_key = "foo"
}
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}
resource "newrelic_plugins_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

	name               = "tf-test-%[1]s"
	enabled            = false
	entities           = [data.newrelic_plugin_component.foo.id]
	runbook_url        = "https://foo.example.com"
	metric             = "Component/Connection/Clients[connections]"
	metric_description = "my-metric-description"
	plugin_id          = data.newrelic_plugin.foo.id
	plugin_guid        = data.newrelic_plugin.foo.guid
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
