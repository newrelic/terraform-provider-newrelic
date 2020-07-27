// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNewRelicAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameUpdated := fmt.Sprintf("tf-test-updated-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", rName),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "type", "apm_app_metric"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "metric", "apdex"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "condition_scope", "application"),
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
			// Test: Check no diff on re-apply
			{
				Config:             testAccNewRelicAlertConditionConfig(rName),
				ExpectNonEmptyPlan: false,
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertConditionConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "name", rNameUpdated),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "type", "apm_app_metric"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "runbook_url", "https://bar.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "entities.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "metric", "error_percentage"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "condition_scope", "application"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.operator", "above"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.priority", "critical"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.threshold", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_condition.foo", "term.3409672004.time_function", "any"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAlertCondition_ZeroThreshold(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionConfigThreshold(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_FloatThreshold(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionConfigThreshold(rName, 0.5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_AlertPolicyNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(rName),
				Config:    testAccNewRelicAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicAlertConditionExists("newrelic_alert_condition.foo"),
			},
		},
	})
}
