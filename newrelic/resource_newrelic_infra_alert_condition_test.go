package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicInfraAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_infra_alert_condition.foo"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
					resource.TestCheckNoResourceAttr("newrelic_infra_alert_condition.foo", "warning"),
				),
			},
			// Test: No diff on reapply
			{
				Config:             testAccCheckNewRelicInfraAlertConditionConfig(rName),
				ExpectNonEmptyPlan: false,
			},
			// Test: Update
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigUpdated(rName),
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
	rName := acctest.RandString(5)
	whereClause := "(`hostname` LIKE '%cassandra%')"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithWhere(rName, whereClause),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
				),
			},
			// Test: No diff on reapply
			{
				Config:             testAccCheckNewRelicInfraAlertConditionConfigWithWhere(rName, whereClause),
				ExpectNonEmptyPlan: false,
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
	rName := acctest.RandString(5)
	resourceName := "newrelic_infra_alert_condition.foo"
	integrationProvider := "Elb"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithIntegrationProvider(rName, integrationProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
				),
			},
			// Test: No diff on re-apply
			{
				Config:             testAccCheckNewRelicInfraAlertConditionConfigWithIntegrationProvider(rName, integrationProvider),
				ExpectNonEmptyPlan: false,
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
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithThreshold(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists(resourceName),
				),
			},
			// Test: No diff on re-apply
			{
				Config:             testAccCheckNewRelicInfraAlertConditionConfigWithThreshold(rName),
				ExpectNonEmptyPlan: false,
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

func TestAccNewRelicInfraAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccCheckNewRelicInfraAlertConditionConfig(rName),
				Check:     testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
			},
		},
	})
}

func testAccCheckNewRelicInfraAlertConditionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).InfraClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_infra_alert_condition" {
			continue
		}

		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		_, err = client.GetAlertInfraCondition(policyID, id)
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

		client := testAccProvider.Meta().(*ProviderConfig).InfraClient

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		id := ids[1]

		found, err := client.GetAlertInfraCondition(policyID, id)
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("alert condition not found: %v - %v", id, found)
		}

		return nil
	}
}

func testAccCheckNewRelicInfraAlertConditionConfig(rName string) string {
	return fmt.Sprintf(`

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-%[1]s"
  runbook_url     = "https://foo.example.com"
  type            = "infra_metric"
  event           = "StorageSample"
  select          = "diskFreePercent"
  comparison      = "below"

  critical {
	  duration = 10
	  value = 10
	  time_function = "any"
  }
}
`, rName)
}

func testAccCheckNewRelicInfraAlertConditionConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-updated-%[1]s"
  runbook_url     = "https://bar.example.com"
  type            = "infra_metric"
  event           = "StorageSample"
  select          = "diskFreePercent"
  comparison      = "below"

  critical {
	  duration = 10
	  value = 10
	  time_function = "any"
  }
}
`, rName)
}

func testAccCheckNewRelicInfraAlertConditionConfigWithThreshold(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-%[1]s"
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
`, rName)
}

func testAccCheckNewRelicInfraAlertConditionConfigWithThresholdUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name            = "tf-test-%[1]s"
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
`, rName)
}

func testAccCheckNewRelicInfraAlertConditionConfigWithWhere(rName, where string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name          = "tf-test-%[1]s"
  type          = "infra_process_running"
  process_where = "commandName = 'java'"
  comparison    = "equal"
  where         = "%[2]s"

  critical {
	duration = 10
	value = 0
  }
}
`, rName, where)
}

func testAccCheckNewRelicInfraAlertConditionConfigWithIntegrationProvider(rName, integrationProvider string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name                 = "tf-test-%[1]s"
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
`, rName, integrationProvider)
}
