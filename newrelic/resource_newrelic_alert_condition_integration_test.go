//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAlertCondition_Basic(t *testing.T) {
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

	// resourceName := "newrelic_alert_condition.foo"
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
				),
			},
			// Test: Import
			//{
			//	ResourceName:      resourceName,
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

func TestAccNewRelicAlertCondition_ZeroThreshold(t *testing.T) {
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

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
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

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
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

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

func TestAccNewRelicAlertCondition_ApplicationScopeWithCloseTimer(t *testing.T) {
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAlertConditionApplicationScopeWithCloseTimerConfig(rName),
				ExpectError: regexp.MustCompile("violation_close_timer only supported for apm_app_metric when condition_scope = 'instance'"),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_InstanceScopeWithCloseTimer(t *testing.T) {
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionInstanceScopeWithCloseTimerConfig(rName),
			},
		},
	})
}

func TestAccNewRelicAlertCondition_APMJVMMetricApplicationScope(t *testing.T) {
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionAPMJVMMetricApplicationScopeConfig(rName),
			},
		},
	})
}
func TestAccNewRelicAlertCondition_APMJVMMetricInstanceScope(t *testing.T) {
	t.Skip("Skipping. API has been deprecated and has reached EOL.")

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertConditionAPMJVMMetricInstanceScopeConfig(rName),
			},
		},
	})
}
