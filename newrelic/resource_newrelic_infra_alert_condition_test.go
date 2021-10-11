//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicInfraAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameUpdated := fmt.Sprintf("tf-test-updated-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
					resource.TestCheckNoResourceAttr("newrelic_infra_alert_condition.foo", "warning"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicInfraAlertConditionConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
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

func TestAccNewRelicInfraAlertCondition_Where(t *testing.T) {
	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	whereClause := "(`hostname` LIKE '%cassandra%')"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionConfigWithWhere(rName, whereClause),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
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

func TestAccNewRelicInfraAlertCondition_IntegrationProvider(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resourceName := "newrelic_infra_alert_condition.foo"
	integrationProvider := "Elb"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionConfigWithIntegrationProvider(rName, integrationProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
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

func TestAccNewRelicInfraAlertCondition_Thresholds(t *testing.T) {
	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionConfigWithThreshold(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithThresholdUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
					resource.TestCheckNoResourceAttr(
						"newrelic_infra_alert_condition.foo", "warning"),
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

func TestAccNewRelicInfraAlertCondition_ThresholdFloatValue(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionWithThresholdFloatValue(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
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

func TestAccNewRelicInfraAlertCondition_ViolationCloseTimer(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionWithViolationCloseTimerConfig(rName, 24),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
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

func TestAccNewRelicInfraAlertCondition_ViolationCloseTimerZeroValue(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicInfraAlertConditionWithViolationCloseTimerConfig(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
				),
				ExpectNonEmptyPlan: true,
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

func TestAccNewRelicInfraAlertCondition_MissingPolicy(t *testing.T) {
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicInfraAlertConditionConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccNewRelicInfraAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
			},
		},
	})
}

func TestAccNewRelicInfraAlertCondition_InvalidAttrsForType(t *testing.T) {
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicInfraAlertConditionInvalidAttrsForTypeConfig(rName),
				ExpectError: regexp.MustCompile("not supported"),
			},
		},
	})
}

func TestAccNewRelicInfraAlertCondition_ComputedEvent(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resourceName := "newrelic_infra_alert_condition.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicInfraAlertConditionComputedEvent(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckNewRelicInfraAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_infra_alert_condition" {
			continue
		}

		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		id := ids[1]

		_, err = client.Alerts.GetInfrastructureCondition(id)
		if err == nil {
			return fmt.Errorf("infra Alert condition still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicInfraAlertConditionExists(n string) resource.TestCheckFunc {
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

		id := ids[1]

		found, err := client.Alerts.GetInfrastructureCondition(id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccNewRelicInfraAlertConditionConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name            = "%[1]s"
  runbook_url     = "https://foo.example.com"
  type            = "infra_metric"
  event           = "StorageSample"
  select          = "diskFreePercent"
	comparison      = "below"
	description     = "test description"

  critical {
	  duration = 10
	  value = 10
	  time_function = "any"
  }
}
`, name)
}

func testAccNewRelicInfraAlertConditionConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name            = "%[1]s"
  runbook_url     = "https://bar.example.com"
  type            = "INFRA_METRIC"
  event           = "StorageSample"
  select          = "diskFreePercent"
	comparison      = "BELOW"
	description     = "test description"

  critical {
	  duration = 10
	  value = 10
	  time_function = "ANY"
  }
}
`, name)
}

func testAccNewRelicInfraAlertConditionConfigWithThreshold(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name            = "%[1]s"
  type            = "infra_metric"
  event           = "StorageSample"
  select          = "diskFreePercent"
  comparison      = "below"

  critical {
	duration = 10
	value = 10
	time_function = "any"
  }

  warning {
	duration = 10
	value = 20
	time_function = "any"
  }
}
`, name)
}

func testAccCheckNewRelicInfraAlertConditionConfigWithThresholdUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name            = "%[1]s"
  type            = "infra_metric"
  event           = "StorageSample"
  select          = "diskFreePercent"
  comparison      = "below"

  critical {
    duration = 20
	value = 15
	time_function = "all"
  }
}
`, name)
}

func testAccNewRelicInfraAlertConditionConfigWithWhere(name, where string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name          = "%[1]s"
  type          = "infra_process_running"
  process_where = "commandName = 'java'"
  comparison    = "equal"
  where         = "%[2]s"

  critical {
	duration = 10
	value = 0
  }
}
`, name, where)
}

func testAccNewRelicInfraAlertConditionConfigWithIntegrationProvider(name, integrationProvider string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name                 = "%[1]s"
  type                 = "infra_metric"
  event                = "LoadBalancerSample"
  integration_provider = "%[2]s"
  select               = "provider.healthyHostCount.Minimum"
  comparison           = "below"

  critical {
    duration      = 10
    value         = 1
    time_function = "all"
  }
}
`, name, integrationProvider)
}

func testAccNewRelicInfraAlertConditionInvalidAttrsForTypeConfig(name string) string {
	return fmt.Sprintf(`

resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name            = "%[1]s"
  runbook_url     = "https://foo.example.com"
  type            = "infra_process_running"
  event           = "StorageSample"
  select          = "diskFreePercent"
  comparison      = "below"

  critical {
	  duration = 10
	  value = 10
	  time_function = "any"
  }
}
`, name)
}

func testAccNewRelicInfraAlertConditionComputedEvent(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id
	name                 = "%[1]s"
	type                 = "infra_metric"
	select               = "nr.ingestTimeMs"
	event                = "SystemSample"
	comparison = "above"

	critical {
		duration      = "1440"
		time_function = "all"
		value         = "25"
	}
}
`, name)
}

func testAccNewRelicInfraAlertConditionWithThresholdFloatValue(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
	policy_id            = newrelic_alert_policy.foo.id
	name                 = "%[1]s"
	type                 = "infra_metric"
	select               = "nr.ingestTimeMs"
	comparison           = "above"
	event                = "SystemSample"

	critical {
		duration      = "1440"
		time_function = "all"
		value         = "1.5"
	}
}
`, name)
}

func testAccNewRelicInfraAlertConditionWithViolationCloseTimerConfig(name string, violationCloseTimer int) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
	policy_id            = newrelic_alert_policy.foo.id
	name                 = "%[1]s"
	type                 = "infra_metric"
	select               = "nr.ingestTimeMs"
	comparison           = "above"
	event                = "SystemSample"

	critical {
		duration      = "1440"
		time_function = "all"
		value         = "1.5"
	}
	violation_close_timer = %[2]d
}
`, name, violationCloseTimer)
}
