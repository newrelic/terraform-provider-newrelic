//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicNRQLDropRule_Data(t *testing.T) {
	rand := acctest.RandString(5)
	description := fmt.Sprintf("nrql_drop_rule_%s", rand)
	resourceName := "newrelic_nrql_drop_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNRQLDropRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNRQLDropRuleConfig(description, "drop_data", "SELECT * FROM MyCustomEvent WHERE appName='LoadGeneratingApp' AND environment='development'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNRQLDropRuleExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestAccNewRelicNRQLDropRule_Attributes(t *testing.T) {
	rand := acctest.RandString(5)
	description := fmt.Sprintf("nrql_drop_rule_%s", rand)
	resourceName := "newrelic_nrql_drop_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNRQLDropRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNRQLDropRuleConfig(description, "drop_attributes", "SELECT userEmail, userName FROM MyCustomEvent"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNRQLDropRuleExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestAccNewRelicNRQLDropRule_AttributesInvalidNRQL(t *testing.T) {
	rand := acctest.RandString(5)
	description := fmt.Sprintf("nrql_drop_rule_%s", rand)
	expectedErrorMsg, _ := regexp.Compile(`drop rule create result wasn't returned`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNRQLDropRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config:      testAccNewRelicNRQLDropRuleConfig(description, "drop_attributes", "FROM ContainerSample DONT WORK commandLine"),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicNRQLDropRule_AttributesFromMetricAggregates(t *testing.T) {
	rand := acctest.RandString(5)
	description := fmt.Sprintf("nrql_drop_rule_%s", rand)
	resourceName := "newrelic_nrql_drop_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNRQLDropRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNRQLDropRuleConfig(description, "drop_attributes_from_metric_aggregates", "SELECT containerId FROM Metric WHERE metricName = 'some.metric'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNRQLDropRuleExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestAccNewRelicNRQLDropRule_AccountIDInheritance(t *testing.T) {
	rand := acctest.RandString(5)
	description := fmt.Sprintf("nrql_drop_rule_%s", rand)
	resourceName := "newrelic_nrql_drop_rule.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNRQLDropRuleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNRQLDropRuleAccountIDInheritanceConfig(description, "drop_attributes", "SELECT userEmail, userName FROM MyCustomEvent"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNRQLDropRuleExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func testAccCheckNewRelicNRQLDropRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_nrql_drop_rule" {
			continue
		}

		accountID, ruleID, err := parseNRQLDropRuleIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		_, err = getNRQLDropRuleByID(context.Background(), client, accountID, ruleID)

		if err == nil {
			return fmt.Errorf("drop rule still exists: %s", err)
		}
	}

	return nil
}

func testAccCheckNewRelicNRQLDropRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		accountID, _, err := parseNRQLDropRuleIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.Nrqldroprules.GetList(accountID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicNRQLDropRuleConfig(description, action, nrql string) string {
	return fmt.Sprintf(`
resource "newrelic_nrql_drop_rule" "foo" {
  account_id = "%d"
  description = "%s"
  action = "%s"
  nrql = "%s"
}

`, testAccountID, description, action, nrql)
}

func testAccNewRelicNRQLDropRuleAccountIDInheritanceConfig(description, action, nrql string) string {
	return fmt.Sprintf(`
resource "newrelic_nrql_drop_rule" "foo" {
  description = "%s"
  action = "%s"
  nrql = "%s"
}

`, description, action, nrql)
}
