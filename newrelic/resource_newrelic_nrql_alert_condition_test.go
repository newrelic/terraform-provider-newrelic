//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicNrqlAlertCondition_Basic(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "120", "120", "sTaTiC", "0", "", "60", "30", "259200"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicNrqlAlertConditionConfigBasic(rName, "300", "180", "last_value", "null", "", "60", "30", "259200"),
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
					"term",                 // contains nested attributes that are deprecated
					"nrql",                 // contains nested attributes that are deprecated
					"violation_time_limit", // deprecated in favor of violation_time_limit_seconds
				},
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "static"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_MissingPolicy(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"cadence",
					"60",
					"null",
				),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(fmt.Sprintf("tf-test-%s", rName)),
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"cadence",
					"60",
					"null",
				),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphThresholdDurationValidationErrors(t *testing.T) {
	rNameBaseline := acctest.RandString(5)
	conditionalAttrBaseline := `baseline_direction = "lower_only"`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Baseline condition invalid `threshold_duration`
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfig(
					rNameBaseline,
					"baseline",
					120,
					86460, // outside of accepted range [120, 86400] to test error handling
					"static",
					"0",
					conditionalAttrBaseline,
				),
				ExpectError: regexp.MustCompile("Validation Error"),
			},
			// Test: Baseline condition invalid `threshold_duration`
			{
				Config: testAccNewRelicNrqlAlertConditionNerdGraphConfig(
					rNameBaseline,
					"baseline",
					120,
					60, // outside of accepted range [120, 86400] to test error handling
					"static",
					"0",
					conditionalAttrBaseline,
				),
				ExpectError: regexp.MustCompile("Validation Error"),
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
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
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
					"last_value",
					"null",
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
					"static",
					"0",
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
					300,
					3600,
					"last_value",
					"null",
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
					1200,
					1800,
					"static",
					"0",
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
	conditionalAttr := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
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
					"last_value",
					"null",
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
					"static",
					"0",
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
					300,
					3600,
					"last_value",
					"null",
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
					1200,
					1800,
					"static",
					"0",
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
					"term",                 // contains nested attributes that are deprecated
					"nrql",                 // contains nested attributes that are deprecated
					"violation_time_limit", // deprecated in favor of violation_time_limit_seconds
				},
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "static"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphStreamingMethods(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	delay := "null"
	timer := "null"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (NerdGraph) condition with streaming method cadence
			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"cadence",
					"60",
					timer,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (NerdGraph) condition with streaming method event_timer
			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"event_timer",
					delay,
					"120",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_AggregationDelayZero(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	timer := "null"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (NerdGraph) condition with streaming method cadence

			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"cadence",
					"0",
					timer,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (NerdGraph) condition with streaming method event_timer
			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"event_flow",
					"0",
					timer,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphNrqlEvaluationOffset(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)
	delay := "null"
	timer := "null"
	method := "null"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (NerdGraph) condition with nrql.evaluationOffset
			{
				Config: testAccNewRelicNrqlAlertConditionNrqlEvaluationOffsetNerdGraphConfig(
					rName,
					method,
					delay,
					timer,
					20,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Update (NerdGraph) condition to remove nrql.evaluationOffset and replace with streaming methods
			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
					rName,
					"event_timer",
					delay,
					"120",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			//Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore items with deprecated fields because
				// we don't set deprecated fields on import
				ImportStateVerifyIgnore: []string{
					"term",                 // contains nested attributes that are deprecated
					"nrql",                 // contains nested attributes that are deprecated
					"violation_time_limit", // deprecated in favor of violation_time_limit_seconds
				},
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, "static"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_NerdGraphValidationErrorBadUserInputOnCreate(t *testing.T) {
	rNameStatic := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create fails with error: `fillOption` of static must be provided with `fillValue`
			{
				Config: testAccNewRelicNrqlAlertStaticNerdGraphConfigInvalid(
					rNameStatic,
					"20",
					"120",
					"",
					"60",
				),

				ExpectError: regexp.MustCompile("Validation Error: BAD_USER_INPUT"),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_RevertToDeprecatedSinceValue(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.since_value"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create a condition with since_value
			{
				Config:             testAccNewRelicNrqlAlertConditionSinceValue(rName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			//Test: Update (NerdGraph) condition with streaming method event_timer
			{
				Config:             testAccNewRelicNrqlAlertConditionSinceValueUpdateWithStreamingMethods(rName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Revert to initial condition with since_value
			{
				Config:             testAccNewRelicNrqlAlertConditionSinceValue(rName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
		},
	})
}

// https://github.com/newrelic/terraform-provider-newrelic/issues/1900
func TestAccNewRelicNrqlAlertCondition_StreamingMethodsDefaults(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create a condition without streaming methods
			{
				Config: testAccNewRelicNrqlAlertConditionStreamingMethodsDefaults(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
			// Test: Check no diff on re-apply
			{
				Config:             testAccNewRelicNrqlAlertConditionStreamingMethodsDefaults(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_StaticConditionSlideByNoValueFunction(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create (NerdGraph) static condition with slide by and no value function
			{
				Config: testAccNewRelicNrqlAlertConditionStaticWithSlideByNerdGraphConfig(
					rName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicNrqlAlertCondition_StaticConditionEvaluationDelay(t *testing.T) {
	resourceName := "newrelic_nrql_alert_condition.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNrqlAlertConditionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicNrqlAlertConditionStaticWithEvaluationDelay(
					rName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNrqlAlertConditionExists(resourceName),
				),
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
	aggregationDelay string,
	duration string,
	fillOption string,
	fillValue string,
	conditionalAttrs string,
	aggregationWindow string,
	slideBy string,
	violation_time_limit_seconds string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name                           = "tf-test-%[1]s"
  runbook_url                    = "https://foo.example.com"
  enabled                        = false
  fill_option                    = "%[4]s"
  fill_value                     = %[5]s
  aggregation_window             = %[7]s
  slide_by                       = %[8]s
  aggregation_delay              = %[2]s
	aggregation_method             = "EVENT_FLOW"
  close_violations_on_expiration = true
  open_violation_on_expiration   = true
  expiration_duration            = 120
  violation_time_limit_seconds   = %[9]s

  nrql {
    query = "SELECT uniqueCount(hostname) FROM ComputeSample"
  }

  critical {
    operator              = "above"
    threshold             = 0.75
    threshold_duration    = %[3]s
    threshold_occurrences = "all"
  }

  warning {
    operator              = "equals"
    threshold             = 0.5
    threshold_duration    = 120
    threshold_occurrences = "AT_LEAST_ONCE"
  }

  %[6]s
}
`, name, aggregationDelay, duration, fillOption, fillValue, conditionalAttrs, aggregationWindow, slideBy, violation_time_limit_seconds)
}

// Uses deprecated attributes for test case
func testAccNewRelicNrqlAlertConditionNerdGraphConfigDeprecated(
	name string,
	conditionType string,
	nrqlEvalOffset int,
	termDuration int,
	fillOption string,
	fillValue string,
	conditionalAttr string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id   = newrelic_alert_policy.foo.id

	name                           = "tf-test-%[1]s"
	type                           = "%[2]s"
	runbook_url                    = "https://foo.example.com"
	enabled                        = false
	description                    = "test description"
  fill_option                    = "%[5]s"
	fill_value                     = %[6]s
	close_violations_on_expiration = true
	open_violation_on_expiration   = true
	expiration_duration            = 120
	aggregation_window             = 60

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

	violation_time_limit_seconds = 86400

	%[7]s
}
`, name, conditionType, nrqlEvalOffset, termDuration, fillOption, fillValue, conditionalAttr)
}

// Uses new attributes for test case
func testAccNewRelicNrqlAlertConditionNerdGraphConfig(
	name string,
	conditionType string,
	aggregationDelay int,
	termDuration int,
	fillOption string,
	fillValue string,
	conditionalAttr string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id   = newrelic_alert_policy.foo.id

	name                           = "tf-test-%[1]s"
	type                           = "%[2]s"
	runbook_url                    = "https://foo.example.com"
	enabled                        = false
	description                    = "test description"
	violation_time_limit_seconds   = 3600
	fill_option                    = "%[5]s"
	fill_value                     = %[6]s
	close_violations_on_expiration = true
	open_violation_on_expiration   = true
	expiration_duration            = 120
	aggregation_window             = 60
	aggregation_delay              = %[3]d
	aggregation_method             = "EVENT_FLOW"

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
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

	# Will be baseline_direction depending on condition type
	%[7]s
}
`, name, conditionType, aggregationDelay, termDuration, fillOption, fillValue, conditionalAttr)
}

func testAccNewRelicNrqlAlertConditionStreamingMethodsNerdGraphConfig(
	name string,
	method string,
	delay string,
	timer string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id   = newrelic_alert_policy.foo.id

	name                           = "tf-test-%[1]s"
	type                           = "static"
	runbook_url                    = "https://foo.example.com"
	enabled                        = false
	description                    = "test description"
	violation_time_limit_seconds   = 3600
	close_violations_on_expiration = true
	open_violation_on_expiration   = true
	expiration_duration            = 120
	aggregation_window             = 60

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
	}

	critical {
    operator              = "above"
    threshold             = 0
		threshold_duration    = 120
		threshold_occurrences = "ALL"
	}
	aggregation_method = "%[2]s"
	aggregation_delay = %[3]s
	aggregation_timer = %[4]s
}
`, name, method, delay, timer)
}

func testAccNewRelicNrqlAlertConditionStreamingMethodsDefaults(
	name string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
    name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
    policy_id                      = newrelic_alert_policy.foo.id
    name                           = "tf-test-%[1]s"
    description                    = "Test desc"
    runbook_url                    = "REDACTED"
    enabled                        = true
    violation_time_limit_seconds   = 86400

    critical {
	    operator              = "above"
	    threshold             = 80
	    threshold_duration    = 5 * 60
	    threshold_occurrences = "all"
    }

    nrql {
        query = <<-EOT
			SELECT count(*) FROM TestEvent
		EOT
	}
}
`, name)
}

func testAccNewRelicNrqlAlertConditionSinceValue(
	name string,
) string {
	return fmt.Sprintf(`
	resource "newrelic_alert_policy" "foo" {
      name = "tf-test-%[1]s"
	}


	resource "newrelic_nrql_alert_condition" "since_value" {
      policy_id                      = newrelic_alert_policy.foo.id
	  type                           = "static"
      name                           = "tf-test-%[1]s"
	  description                    = "Test desc"
	  runbook_url                    = "REDACTED"
	  enabled                        = true
	  violation_time_limit_seconds   = 86400
	  expiration_duration            = 120
	  close_violations_on_expiration = true

	  nrql {
		query       = <<-EOT
			SELECT count(*) FROM TestEvent
		EOT
		since_value = 1
	  }
	  critical {
		operator           = "above"
		threshold          = 50
		threshold_duration = 300
		time_function      = "any"
	  }
	  warning {
		operator           = "above"
		threshold          = 10
		threshold_duration = 300
		time_function      = "any"
	  }
	}
`, name)
}

func testAccNewRelicNrqlAlertConditionSinceValueUpdateWithStreamingMethods(
	name string,
) string {
	return fmt.Sprintf(`
	resource "newrelic_alert_policy" "foo" {
      name = "tf-test-%[1]s"
	}

	resource "newrelic_nrql_alert_condition" "since_value" {
      policy_id                      = newrelic_alert_policy.foo.id
	  type                           = "static"
      name                           = "tf-test-%[1]s"
	  description                    = "Test desc"
	  runbook_url                    = "REDACTED"
	  enabled                        = true
	  violation_time_limit_seconds   = 86400
	  expiration_duration            = 120
	  close_violations_on_expiration = true
      aggregation_method             = "event_timer"
      aggregation_timer              = 1200

	  nrql {
		query       = <<-EOT
			SELECT count(*) FROM TestEvent
		EOT
	  }
	  critical {
		operator           = "above"
		threshold          = 50
		threshold_duration = 300
		time_function      = "any"
	  }
	  warning {
		operator           = "above"
		threshold          = 10
		threshold_duration = 300
		time_function      = "any"
	  }
	}
`, name)
}

func testAccNewRelicNrqlAlertConditionNrqlEvaluationOffsetNerdGraphConfig(
	name string,
	method string,
	delay string,
	timer string,
	nrqlEvalOffset int,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id   = newrelic_alert_policy.foo.id

	name                           = "tf-test-%[1]s"
	type                           = "static"
	runbook_url                    = "https://foo.example.com"
	enabled                        = false
	description                    = "test description"
	violation_time_limit_seconds   = 3600
	close_violations_on_expiration = true
	open_violation_on_expiration   = true
	expiration_duration            = 120
	aggregation_window             = 60

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
		evaluation_offset = %[5]d
	}

	critical {
    operator              = "above"
    threshold             = 0
		threshold_duration    = 120
		threshold_occurrences = "ALL"
	}
	aggregation_delay = %[3]s
	aggregation_timer = %[4]s
}
`, name, method, delay, timer, nrqlEvalOffset)
}

// `fill_option` of `static` must provide `fill_value`
func testAccNewRelicNrqlAlertStaticNerdGraphConfigInvalid(
	name string,
	evaluationOffset string,
	duration string,
	conditionalAttrs string,
	aggregationWindow string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}
resource "newrelic_nrql_alert_condition" "foo" {
	policy_id = newrelic_alert_policy.foo.id
  name                           = "tf-test-%[1]s"
  runbook_url                    = "https://foo.example.com"
  enabled                        = false
  violation_time_limit_seconds   = 28800
  fill_option                    = "static"
  aggregation_window             = %[5]s
  close_violations_on_expiration = true
  open_violation_on_expiration   = true
  expiration_duration            = 120
	nrql {
    query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
    evaluation_offset = "%[2]s"
	}
	critical {
    operator              = "above"
    threshold             = 0.75
    threshold_duration    = %[3]s
    threshold_occurrences = "all"
	}
	warning {
		operator              = "equals"
		threshold             = 0.5
		threshold_duration    = 120
		threshold_occurrences = "AT_LEAST_ONCE"
	}

	%[4]s
}
`, name, evaluationOffset, duration, conditionalAttrs, aggregationWindow)
}

func testAccNewRelicNrqlAlertConditionStaticWithSlideByNerdGraphConfig(
	name string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id   = newrelic_alert_policy.foo.id

	name                           = "tf-test-%[1]s"
	type                           = "static"
	runbook_url                    = "https://foo.example.com"
	enabled                        = false
	description                    = "test description"
	violation_time_limit_seconds   = 3600
	close_violations_on_expiration = true
	open_violation_on_expiration   = true
	expiration_duration            = 120
 	aggregation_delay              = 120
	aggregation_method             = "event_flow"
	slide_by					   = 30

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
	}

	critical {
    operator              = "above"
    threshold             = 0
		threshold_duration    = 120
		threshold_occurrences = "ALL"
	}
}
`, name)
}

func testAccNewRelicNrqlAlertConditionStaticWithEvaluationDelay(
	name string,
) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
	name = "tf-test-%[1]s"
}

resource "newrelic_nrql_alert_condition" "foo" {
	policy_id   = newrelic_alert_policy.foo.id

	name                           = "tf-test-%[1]s"
	type                           = "static"
	runbook_url                    = "https://foo.example.com"
	enabled                        = false
	description                    = "test description"
	violation_time_limit_seconds   = 3600
	close_violations_on_expiration = true
	open_violation_on_expiration   = true
	expiration_duration            = 120
 	aggregation_delay              = 120
	aggregation_method             = "event_flow"
	slide_by					   = 30
	evaluation_delay			   = 60

	nrql {
		query             = "SELECT uniqueCount(hostname) FROM ComputeSample"
	}

	critical {
    operator              = "above"
    threshold             = 0
		threshold_duration    = 120
		threshold_occurrences = "ALL"
	}
}
`, name)
}
