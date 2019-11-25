package newrelic

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicInfraAlertCondition_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "runbook_url", "https://foo.example.com"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.duration", "10"),
					resource.TestCheckNoResourceAttr(
						"newrelic_infra_alert_condition.foo", "warning"),
				),
			},
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "runbook_url", "https://bar.example.com"),
				),
			},
		},
	})
}

func TestAccNewRelicInfraAlertCondition_Where(t *testing.T) {
	rName := acctest.RandString(5)
	whereClause := "(`hostname` LIKE '%cassandra%')"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithWhere(rName, whereClause),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.value", "0"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "where", whereClause),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "process_where", "commandName = 'java'"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "comparison", "equal"),
				),
			},
		},
	})
}

func TestAccNewRelicInfraAlertCondition_IntegrationProvider(t *testing.T) {
	key := "ENABLE_NEWRELIC_INTEGRATION_PROVIDER"
	enableNewRelicIntegrationProvider := os.Getenv(key)
	if enableNewRelicIntegrationProvider == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	rName := acctest.RandString(5)
	integrationProvider := "Elb"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithIntegrationProvider(rName, integrationProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.value", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "integration_provider", integrationProvider),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "comparison", "below"),
				),
			},
		},
	})
}
func TestAccNewRelicInfraAlertCondition_Thresholds(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithThreshold(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.duration", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.value", "10"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.time_function", "any"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "warning.0.value", "20"),
				),
			},
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfigWithThresholdUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInfraAlertConditionExists("newrelic_infra_alert_condition.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.duration", "20"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.value", "15"),
					resource.TestCheckResourceAttr(
						"newrelic_infra_alert_condition.foo", "critical.0.time_function", "all"),
					resource.TestCheckNoResourceAttr(
						"newrelic_infra_alert_condition.foo", "warning"),
				),
			},
		},
	})
}

func TestAccNewRelicInfraAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicInfraAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInfraAlertConditionConfig(rName),
			},
			{
				PreConfig: deletePolicy(fmt.Sprintf("tf-test-%s", rName)),
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
