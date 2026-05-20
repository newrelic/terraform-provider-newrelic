//go:build integration

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccNewRelicAccountCardinalityLimit_Default covers the full DEFAULT lifecycle:
// create → read (live) → update (in-place upsert) → destroy (reset to 100,000).
func TestAccNewRelicAccountCardinalityLimit_Default(t *testing.T) {
	resourceName := "newrelic_account_cardinality_limit.default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityLimitDefaultDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicCardinalityLimitDefaultConfig(150000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "150000"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", ""),
				),
			},
			// Update (in-place upsert — no destroy/recreate)
			{
				Config: testAccNewRelicCardinalityLimitDefaultConfig(200000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "200000"),
				),
			},
		},
	})
}

// TestAccNewRelicAccountCardinalityLimit_DefaultImport verifies import reconstructs state correctly.
func TestAccNewRelicAccountCardinalityLimit_DefaultImport(t *testing.T) {
	resourceName := "newrelic_account_cardinality_limit.default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityLimitDefaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityLimitDefaultConfig(150000),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicAccountCardinalityLimit_PerMetric covers the PER_METRIC lifecycle:
// create → state check → destroy (reset to account-wide default).
func TestAccNewRelicAccountCardinalityLimit_PerMetric(t *testing.T) {
	resourceName := "newrelic_account_cardinality_limit.per_metric"
	metricName := "tf.test.cardinality.metric"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityLimitPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityLimitPerMetricConfig(metricName, 200000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "PER_METRIC"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", metricName),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "200000"),
				),
			},
		},
	})
}

// TestAccNewRelicAccountCardinalityLimit_PerMetricImport verifies import reconstructs mode and
// metric_name from the resource ID. cardinality_limit is excluded from verification because the
// API does not return per-metric override values.
func TestAccNewRelicAccountCardinalityLimit_PerMetricImport(t *testing.T) {
	resourceName := "newrelic_account_cardinality_limit.per_metric"
	metricName := "tf.test.cardinality.import.metric"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityLimitPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityLimitPerMetricConfig(metricName, 200000),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cardinality_limit"},
			},
		},
	})
}

// TestAccNewRelicAccountCardinalityLimit_PerMetricRequiresMetricName verifies that setting
// mode = PER_METRIC without metric_name is rejected at plan time.
func TestAccNewRelicAccountCardinalityLimit_PerMetricRequiresMetricName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityLimitPerMetricMissingMetricNameConfig(),
				ExpectError: regexp.MustCompile(`metric_name is required when mode is "PER_METRIC"`),
			},
		},
	})
}

// TestAccNewRelicAccountCardinalityLimit_DefaultMustNotHaveMetricName verifies that setting
// mode = DEFAULT with metric_name is rejected at plan time.
func TestAccNewRelicAccountCardinalityLimit_DefaultMustNotHaveMetricName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityLimitDefaultWithMetricNameConfig(),
				ExpectError: regexp.MustCompile(`metric_name must not be set when mode is "DEFAULT"`),
			},
		},
	})
}

// testAccCheckNewRelicCardinalityLimitDefaultDestroy verifies that the account-wide default has
// been reset to the New Relic platform default of 100,000 after destroy.
func testAccCheckNewRelicCardinalityLimitDefaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_account_cardinality_limit" {
			continue
		}

		accountID, metricName, err := parseCardinalityLimitID(r.Primary.ID)
		if err != nil {
			return err
		}

		if metricName != "" {
			continue // PER_METRIC — handled separately
		}

		limits, err := client.DataManagement.GetLimitsWithContext(context.Background(), accountID)
		if err != nil {
			return fmt.Errorf("error fetching limits after destroy: %w", err)
		}

		if limits != nil {
			for _, l := range *limits {
				if l.Name == cardinalityLimitName {
					if l.Value != cardinalityLimitPlatformDefault {
						return fmt.Errorf(
							"expected cardinality limit to be reset to %d after destroy, got %d",
							cardinalityLimitPlatformDefault, l.Value,
						)
					}
					return nil
				}
			}
		}
	}

	return nil
}

// testAccCheckNewRelicCardinalityLimitPerMetricDestroy is a no-op: per-metric override values
// cannot be read back via the API so we cannot verify the reset value. The TF resource is
// removed from state by the framework automatically.
func testAccCheckNewRelicCardinalityLimitPerMetricDestroy(_ *terraform.State) error {
	return nil
}

func testAccNewRelicCardinalityLimitDefaultConfig(limit int) string {
	return fmt.Sprintf(`
resource "newrelic_account_cardinality_limit" "default" {
  mode              = "DEFAULT"
  cardinality_limit = %d
}
`, limit)
}

func testAccNewRelicCardinalityLimitPerMetricConfig(metricName string, limit int) string {
	return fmt.Sprintf(`
resource "newrelic_account_cardinality_limit" "per_metric" {
  mode              = "PER_METRIC"
  metric_name       = %q
  cardinality_limit = %d
}
`, metricName, limit)
}

func testAccNewRelicCardinalityLimitPerMetricMissingMetricNameConfig() string {
	return `
resource "newrelic_account_cardinality_limit" "per_metric" {
  mode              = "PER_METRIC"
  cardinality_limit = 200000
}
`
}

func testAccNewRelicCardinalityLimitDefaultWithMetricNameConfig() string {
	return `
resource "newrelic_account_cardinality_limit" "default" {
  mode              = "DEFAULT"
  metric_name       = "some.metric"
  cardinality_limit = 150000
}
`
}
