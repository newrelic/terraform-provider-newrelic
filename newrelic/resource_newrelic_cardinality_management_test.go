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

// TestAccNewRelicCardinalityManagement_Default covers the full DEFAULT lifecycle:
// create → read (live) → update (in-place upsert) → destroy (reset to 100,000).
func TestAccNewRelicCardinalityManagement_Default(t *testing.T) {
	resourceName := "newrelic_cardinality_management.default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementDefaultDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicCardinalityManagementDefaultConfig(150000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "150000"),
				),
			},
			// Update (in-place upsert — no destroy/recreate)
			{
				Config: testAccNewRelicCardinalityManagementDefaultConfig(200000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "200000"),
				),
			},
		},
	})
}

// TestAccNewRelicCardinalityManagement_DefaultImport verifies import reconstructs state correctly.
func TestAccNewRelicCardinalityManagement_DefaultImport(t *testing.T) {
	resourceName := "newrelic_cardinality_management.default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementDefaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementDefaultConfig(150000),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicCardinalityManagement_PerMetric covers the PER_METRIC lifecycle with a
// single metric: create → state check → destroy (reset to 100,000).
func TestAccNewRelicCardinalityManagement_PerMetric(t *testing.T) {
	resourceName := "newrelic_cardinality_management.per_metric"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig(
					"tf.test.cardinality.metric.single",
					200000,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "PER_METRIC"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
				),
			},
		},
	})
}

// TestAccNewRelicCardinalityManagement_PerMetricMultiple verifies that multiple metric
// blocks can be managed within a single resource.
func TestAccNewRelicCardinalityManagement_PerMetricMultiple(t *testing.T) {
	resourceName := "newrelic_cardinality_management.per_metric"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementPerMetricDestroy,
		Steps: []resource.TestStep{
			// Create with two metrics
			{
				Config: testAccNewRelicCardinalityManagementMultiMetricConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "PER_METRIC"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "2"),
				),
			},
			// Remove one metric (should reset the removed metric to platform default)
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig(
					"tf.test.cardinality.metric.a",
					200000,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
				),
			},
		},
	})
}

// TestAccNewRelicCardinalityManagement_PerMetricImport verifies import reconstructs mode
// from the resource ID. The metric set cannot be read back via the API so it is excluded
// from verification.
func TestAccNewRelicCardinalityManagement_PerMetricImport(t *testing.T) {
	resourceName := "newrelic_cardinality_management.per_metric"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig(
					"tf.test.cardinality.import.metric",
					200000,
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metric"},
			},
		},
	})
}

// TestAccNewRelicCardinalityManagement_PerMetricRequiresMetricBlocks verifies that setting
// mode = PER_METRIC without any metric blocks is rejected at plan time.
func TestAccNewRelicCardinalityManagement_PerMetricRequiresMetricBlocks(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementPerMetricMissingBlocksConfig(),
				ExpectError: regexp.MustCompile(`at least one metric block is required when mode is "PER_METRIC"`),
			},
		},
	})
}

// TestAccNewRelicCardinalityManagement_DefaultMustNotHaveMetricBlocks verifies that setting
// mode = DEFAULT with metric blocks is rejected at plan time.
func TestAccNewRelicCardinalityManagement_DefaultMustNotHaveMetricBlocks(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementDefaultWithMetricBlockConfig(),
				ExpectError: regexp.MustCompile(`metric blocks must not be set when mode is "DEFAULT"`),
			},
		},
	})
}

// testAccCheckNewRelicCardinalityManagementDefaultDestroy verifies that the account-wide
// default has been reset to the New Relic platform default of 100,000 after destroy.
func testAccCheckNewRelicCardinalityManagementDefaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cardinality_management" {
			continue
		}

		accountID, mode, err := parseCardinalityManagementID(r.Primary.ID)
		if err != nil {
			return err
		}

		if mode != cardinalityModeDefault {
			continue
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

// testAccCheckNewRelicCardinalityManagementPerMetricDestroy is a no-op: per-metric override
// values cannot be read back via the API so we cannot verify the reset value.
func testAccCheckNewRelicCardinalityManagementPerMetricDestroy(_ *terraform.State) error {
	return nil
}

func testAccNewRelicCardinalityManagementDefaultConfig(limit int) string {
	return fmt.Sprintf(`
resource "newrelic_cardinality_management" "default" {
  mode              = "DEFAULT"
  cardinality_limit = %d
}
`, limit)
}

func testAccNewRelicCardinalityManagementPerMetricConfig(metricName string, limit int) string {
	return fmt.Sprintf(`
resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"
  metric {
    name              = %q
    cardinality_limit = %d
  }
}
`, metricName, limit)
}

func testAccNewRelicCardinalityManagementMultiMetricConfig() string {
	return `
resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"
  metric {
    name              = "tf.test.cardinality.metric.a"
    cardinality_limit = 200000
  }
  metric {
    name              = "tf.test.cardinality.metric.b"
    cardinality_limit = 300000
  }
}
`
}

func testAccNewRelicCardinalityManagementPerMetricMissingBlocksConfig() string {
	return `
resource "newrelic_cardinality_management" "per_metric" {
  mode              = "PER_METRIC"
  cardinality_limit = 200000
}
`
}

func testAccNewRelicCardinalityManagementDefaultWithMetricBlockConfig() string {
	return `
resource "newrelic_cardinality_management" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
  metric {
    name              = "some.metric"
    cardinality_limit = 200000
  }
}
`
}
