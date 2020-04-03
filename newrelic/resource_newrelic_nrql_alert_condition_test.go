package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func TestAccNewRelicNrqlAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "20", "2", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "5", "10", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore items with deprecated fields because
				// we don't set deprecated fields on import
				ImportStateVerifyIgnore: []string{"term", "nrql", "violation_time_limit"},
				ImportStateIdFunc:       testAccImportStateIDFunc(resourceName, "static"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_TypeOutlier(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNrqlAlertConditionConfigTypeOutlier(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore items with deprecated fields because
				// we don't set deprecated fields on import
				ImportStateVerifyIgnore: []string{"account_id", "term", "nrql", "violation_time_limit"},
				ImportStateIdFunc:       testAccImportStateIDFunc(resourceName, "outlier"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNrqlAlertConditionConfigTypeOutlier(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config:    testAccNewRelicNrqlAlertConditionConfigTypeOutlier(rName),
				Check:     testAccCheckNewRelicNrqlAlertConditionExists("newrelic_nrql_alert_condition.foo"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphBaseline(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	conditionType := "baseline"
	conditionalAttr := `baseline_direction = "LOWER_ONLY"`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (Deprecated)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfigDeprecated(
					rName,
					conditionType,
					1,
					60,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (Deprecated)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfigDeprecated(
					rName,
					conditionType,
					20,
					30,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Create (NerdGraph)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfig(
					rName,
					conditionType,
					5,
					3600,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (NerdGraph)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfig(
					rName,
					conditionType,
					20,
					1800,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},

			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore items with deprecated fields because
				// we don't set deprecated fields on import
				ImportStateVerifyIgnore: []string{
					"term",
					"nrql",
					"violation_time_limit",
					"value_function", // does not exist for type `baseline`
				},
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "baseline"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphStatic(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	conditionType := "static"
	conditionalAttr := `value_function = "single_value"`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (Deprecated)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfigDeprecated(
					rName,
					conditionType,
					5,
					60,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (Deprecated)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfigDeprecated(
					rName,
					conditionType,
					20,
					30,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Create (NerdGraph)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfig(
					rName,
					conditionType,
					5,
					3600,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (NerdGraph)
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfig(
					rName,
					conditionType,
					20,
					1800,
					conditionalAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore items with deprecated fields because
				// we don't set deprecated fields on import
				ImportStateVerifyIgnore: []string{
					"term", // contains nested attributes that are deprecated
					"nrql", // contains nested attributes that are deprecated
					"violation_time_limit",
				},
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "static"),
			},

			// TODO: TEST ERROR SCENARIOS!!!!!!!!!
		},
	})
}

func testAccImportStateIDFunc(resourceName string, metadata string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		idWithMetadata := fmt.Sprintf("%s:%s", rs.Primary.ID, metadata)

		return idWithMetadata, nil
	}
}

func testAccCheckNewRelicNrqlAlertConditionDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient
	hasNerdGraphCreds := providerConfig.hasNerdGraphCredentials()

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_nrql_alert_condition" {
			continue
		}

		var accountID int
		var err error

		ids, err := parseHashedIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		var policyID int
		var conditionID int

		if len(ids) > 2 {
			policyID = ids[1]
			conditionID = ids[2]
		} else {
			policyID = ids[0]
			conditionID = ids[1]
		}

		if hasNerdGraphCreds {
			accountID = providerConfig.AccountID

			if r.Primary.Attributes["account_id"] != "" {
				accountID, err = strconv.Atoi(r.Primary.Attributes["account_id"])
				if err != nil {
					return err
				}
			}

			_, err = client.Alerts.GetNrqlConditionQuery(accountID, conditionID)
			if err == nil {
				return fmt.Errorf("NRQL Alert condition still exists") //nolint:golint
			}
		} else {
			_, err = client.Alerts.GetNrqlCondition(policyID, conditionID)
			if err == nil {
				return fmt.Errorf("NRQL Alert condition still exists") //nolint:golint
			}
		}
	}

	return nil
}

func testAccCheckNewRelicNrqlAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient
		hasNerdGraphCreds := providerConfig.hasNerdGraphCredentials()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no alert condition ID is set")
		}

		var accountID int
		var err error

		ids, err := parseHashedIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		var policyID int
		var conditionID int

		if len(ids) > 2 {
			policyID = ids[1]
			conditionID = ids[2]
		} else {
			policyID = ids[0]
			conditionID = ids[1]
		}

		if hasNerdGraphCreds && rs.Primary.Attributes["type"] != "outlier" {
			accountID = providerConfig.AccountID

			if rs.Primary.Attributes["account_id"] != "" {
				accountID, err = strconv.Atoi(rs.Primary.Attributes["account_id"])
				if err != nil {
					return err
				}
			}

			var found *alerts.NrqlAlertCondition
			found, err = client.Alerts.GetNrqlConditionQuery(accountID, conditionID)
			if err != nil {
				return err
			}

			if found.ID != strconv.Itoa(conditionID) {
				return fmt.Errorf("alert condition not found: %v - %v", conditionID, found)
			}

			return nil
		}

		found, err := client.Alerts.GetNrqlCondition(policyID, conditionID)
		if err != nil {
			return err
		}

		if found.ID != conditionID {
			return fmt.Errorf("alert condition not found: %v - %v", conditionID, found)
		}

		return nil
	}
}

func testAccNewRelicNrqlAlertConditionConfigBasic(
	name string,
	sinceValue string,
	duration string,
	conditionalAttrs string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

  name            = "tf-test-%[1]s"
  runbook_url     = "https://foo.example.com"
  enabled         = false
  violation_time_limit_seconds = 3600

	nrql {
    query         = "SELECT uniqueCount(hostname) FROM ComputeSample"
    since_value   = "%[2]s"
	}

	term {
    duration      = %[3]s
    operator      = "above"
    priority      = "critical"
    threshold     = 0.75
    time_function = "all"
	}

	term {
		duration      = 3
		operator      = "above"
		priority      = "warning"
		threshold     = 0.5
		time_function = "any"
	}

	value_function  = "single_value"

	%[4]s
}
`, name, sinceValue, duration, conditionalAttrs)
}

func testAccNewRelicNrqlAlertConditionConfigTypeOutlier(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id
	name            = "tf-test-outlier-%[1]s"
	type            = "outlier"
  expected_groups = 2
	ignore_overlap  = true

  runbook_url     = "https://bar.example.com"
  enabled         = false
  violation_time_limit_seconds = 7200

	nrql {
    query         = "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip"
    since_value   = "3"
  }

  term {
    duration      = 10
    operator      = "above"
    priority      = "critical"
    threshold     = "0.65"
    time_function = "all"
  }
}
`, name)
}

// Uses deprecated attributes for test case
func testAccNewRelicNrqlAlertConditionNerdGraphConfigDeprecated(
	name string,
	conditionType string,
	nrqlEvalOffset int,
	termDuration int,
	conditionalAttr string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	account_id  = 2520528
	policy_id   = newrelic_alert_policy.foo.id

	name        = "tf-test-%[1]s"
	type        = "%[2]s"
	runbook_url = "https://foo.example.com"
	enabled     = false
	description = "test description"

	nrql {
		query       = "SELECT uniqueCount(hostname) FROM ComputeSample"
		since_value = "%[3]d"
	}

	term {
		duration      = %[4]d
		operator      = "above"
		priority      = "critical"
		threshold     = 1.5
		time_function = "all"
	}

	term {
		duration      = 2
		operator      = "above"
		priority      = "warning"
		threshold     = 1.1
		time_function = "any"
	}

	violation_time_limit_seconds = 3600

	%[5]s
}
`, name, conditionType, nrqlEvalOffset, termDuration, conditionalAttr)
}

// Uses new attributes for test case
func testAccNewRelicNrqlAlertConditionNerdGraphConfig(
	name string,
	conditionType string,
	nrqlEvalOffset int,
	termDuration int,
	conditionalAttr string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	account_id  = 2520528
	policy_id   = newrelic_alert_policy.foo.id

	name                 = "tf-test-%[1]s"
	type                 = "%[2]s"
	runbook_url          = "https://foo.example.com"
	enabled              = false
	description          = "test description"
	violation_time_limit = "ONE_HOUR"

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
		evaluation_offset = %[3]d
	}

	term {
    operator              = "above"
    priority              = "critical"
    threshold             = 1.25
		threshold_duration    = %[4]d
		threshold_occurrences = "ALL"
	}

	term {
    operator              = "above"
    priority              = "warning"
    threshold             = 1.1
		threshold_duration    = %[4]d
		threshold_occurrences = "AT_LEAST_ONCE"
	}

	# Will be baseline_direction or value_function depending on condition type
	%[5]s
}
`, name, conditionType, nrqlEvalOffset, termDuration, conditionalAttr)
}
