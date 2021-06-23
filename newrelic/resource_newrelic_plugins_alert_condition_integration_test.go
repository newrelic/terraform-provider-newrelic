// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicPluginsAlertCondition_Basic(t *testing.T) {
	t.Skip()

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

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
	t.Skip()

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

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

func TestAccNewRelicPluginsAlertCondition_MissingPolicy(t *testing.T) {
	t.Skip()

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

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
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccCheckNewRelicPluginsAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicPluginsAlertConditionExists("newrelic_plugins_alert_condition.foo"),
			},
		},
	})
}
