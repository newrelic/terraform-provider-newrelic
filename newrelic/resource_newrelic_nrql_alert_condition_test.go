package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "20", "120", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// // Test: Update
			// {
			// 	Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "5", "180", ""),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
			// 	),
			// },
			// // Test: Import
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// Ignore items with deprecated fields because
			// 	// we don't set deprecated fields on import
			// 	ImportStateVerifyIgnore: []string{"term", "nrql", "violation_time_limit"},
			// 	ImportStateIdFunc:       testAccImportStateIDFunc(resourceName, "static"),
			// },
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)
	conditionType := "outlier"
	conditionalAttr := `expected_groups = 2
	open_violation_on_group_overlap = true`
	facetClause := `FACET host`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNrqlAlertConditionOutlierNerdGraphConfig(
					fmt.Sprintf("tf-test-%s", rName),
					conditionType,
					3,
					3600,
					conditionalAttr,
					facetClause,
				),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config: testAccNewRelicNrqlAlertConditionOutlierNerdGraphConfig(
					fmt.Sprintf("tf-test-%s", rName),
					conditionType,
					3,
					3600,
					conditionalAttr,
					facetClause,
				),
				Check: testAccCheckNewRelicNrqlAlertConditionExists("newrelic_nrql_alert_condition.foo"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphBaseline(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	conditionType := "baseline"
	conditionalAttr := `baseline_direction = "lower_only"` // value transformed to UPPERCASE in expand/flatten

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
	conditionalAttr := `value_function = "Single_valuE"` // value transformed to UPPERCASE in expand/flatten

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
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphOutlier(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	conditionType := "outlier"
	conditionalAttr := `expected_groups = 2
	open_violation_on_group_overlap = true`
	facetClause := `FACET host`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (NerdGraph)
			{
				Config: testAccNewRelicNrqlAlertConditionOutlierNerdGraphConfig(
					rName,
					conditionType,
					3,
					3600,
					conditionalAttr,
					facetClause,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (NerdGraph)
			{
				Config: testAccNewRelicNrqlAlertConditionOutlierNerdGraphConfig(
					rName,
					conditionType,
					3,
					1800,
					conditionalAttr,
					facetClause,
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
				},
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "outlier"),
			},
		},
	})
}

func testAccCheckNewRelicNrqlAlertConditionDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

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

		conditionID := ids[1]
		accountID = providerConfig.AccountID

		if r.Primary.Attributes["account_id"] != "" {
			accountID, err = strconv.Atoi(r.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		if _, err = client.Alerts.GetNrqlConditionQuery(accountID, strconv.Itoa(conditionID)); err == nil {
			return fmt.Errorf("NRQL Alert condition still exists") //nolint:golint
		}
	}

	return nil
}

func testAccCheckNewRelicNrqlAlertConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

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

		conditionID := ids[1]
		accountID = providerConfig.AccountID

		if rs.Primary.Attributes["account_id"] != "" {
			accountID, err = strconv.Atoi(rs.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		found, err := client.Alerts.GetNrqlConditionQuery(accountID, strconv.Itoa(conditionID))
		if err != nil {
			return err
		}

		if found.ID != strconv.Itoa(conditionID) {
			return fmt.Errorf("alert condition not found: %v - %v", conditionID, found)
		}

		return nil
	}
}

func testAccNewRelicNrqlAlertConditionConfigBasic(
	name string,
	evaluationOffset string,
	duration string,
	conditionalAttrs string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id

  name                 = "tf-test-%[1]s"
  runbook_url          = "https://foo.example.com"
  enabled              = false
  violation_time_limit = "EIGHT_HOURS"

	nrql {
    query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
    evaluation_offset = "%[2]s"
	}

	critical {
    operator              = "above"
    threshold             = 0.75
    threshold_duration    = %[3]s
    threshold_occurrences = "ALL"
	}

	warning {
		operator              = "above"
		threshold             = 0.5
		threshold_duration    = 120
		threshold_occurrences = "AT_LEAST_ONCE"
	}

	value_function  = "single_value"

	%[4]s
}
`, name, evaluationOffset, duration, conditionalAttrs)
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
		threshold     = 1.55
		time_function = "all"
	}

	term {
		duration      = 2
		operator      = "above"
		priority      = "warning"
		threshold     = 1.12
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
	violation_time_limit = "one_hour"

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
		evaluation_offset = %[3]d
	}

	critical {
    operator              = "above"
    threshold             = 1.25666
		threshold_duration    = %[4]d
		threshold_occurrences = "ALL"
	}

	warning {
    operator              = "above"
    threshold             = 1.1666
		threshold_duration    = %[4]d
		threshold_occurrences = "AT_LEAST_ONCE"
	}

	# Will be baseline_direction or value_function depending on condition type
	%[5]s
}
`, name, conditionType, nrqlEvalOffset, termDuration, conditionalAttr)
}

func testAccNewRelicNrqlAlertConditionOutlierNerdGraphConfig(
	name string,
	conditionType string,
	nrqlEvalOffset int,
	termDuration int,
	conditionalAttr string,
	facetClause string,
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
	violation_time_limit = "one_hour"

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample %[6]s"
		evaluation_offset = %[3]d
	}

	critical {
    operator              = "above"
    threshold             = 1.25
		threshold_duration    = %[4]d
		threshold_occurrences = "ALL"
	}

	# Will be one of baseline_direction, value_function, expected_groups, or open_violation_on_group_overlap depending on condition type
	%[5]s
}
`, name, conditionType, nrqlEvalOffset, termDuration, conditionalAttr, facetClause)
}
