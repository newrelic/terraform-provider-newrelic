//go:build integration || INGEST

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCardinalityManagement_Default(t *testing.T) {
	resourceName := "newrelic_cardinality_management.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementDefaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementDefaultConfig(200000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "200000"),
				),
			},
			// cardinality_limit is not ForceNew, so this is an in-place update.
			{
				Config: testAccNewRelicCardinalityManagementDefaultConfig(250000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "250000"),
				),
			},
			// Import by account_id:DEFAULT.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_PerMetric_Single(t *testing.T) {
	resourceName := "newrelic_cardinality_management.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig([]testMetricEntry{
					{name: "test.cardinality.single.tf", limit: 150000},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "PER_METRIC"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.0.name", "test.cardinality.single.tf"),
					resource.TestCheckResourceAttr(resourceName, "metric.0.cardinality_limit", "150000"),
				),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_PerMetric_Multiple(t *testing.T) {
	resourceName := "newrelic_cardinality_management.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig([]testMetricEntry{
					{name: "test.cardinality.multi.one.tf", limit: 150000},
					{name: "test.cardinality.multi.two.tf", limit: 200000},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "PER_METRIC"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "metric.0.name", "test.cardinality.multi.one.tf"),
					resource.TestCheckResourceAttr(resourceName, "metric.0.cardinality_limit", "150000"),
					resource.TestCheckResourceAttr(resourceName, "metric.1.name", "test.cardinality.multi.two.tf"),
					resource.TestCheckResourceAttr(resourceName, "metric.1.cardinality_limit", "200000"),
				),
			},
			// Add a third metric without recreating the resource.
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig([]testMetricEntry{
					{name: "test.cardinality.multi.one.tf", limit: 150000},
					{name: "test.cardinality.multi.two.tf", limit: 200000},
					{name: "test.cardinality.multi.three.tf", limit: 175000},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metric.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "metric.2.name", "test.cardinality.multi.three.tf"),
					resource.TestCheckResourceAttr(resourceName, "metric.2.cardinality_limit", "175000"),
				),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_InvalidDefault_WithMetricBlocks(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementInvalidDefaultWithMetricBlocksConfig(),
				ExpectError: regexp.MustCompile(`metric blocks must not be set when mode is "DEFAULT"`),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_InvalidDefault_MissingLimit(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementInvalidDefaultMissingLimitConfig(),
				ExpectError: regexp.MustCompile(`cardinality_limit is required when mode is "DEFAULT"`),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_InvalidPerMetric_WithTopLevelLimit(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementInvalidPerMetricWithTopLevelLimitConfig(),
				ExpectError: regexp.MustCompile(`cardinality_limit must not be set at the top level when mode is "PER_METRIC"`),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_InvalidPerMetric_NoMetricBlocks(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementInvalidPerMetricNoMetricBlocksConfig(),
				ExpectError: regexp.MustCompile(`at least one metric block is required when mode is "PER_METRIC"`),
			},
		},
	})
}

// testAccCheckNewRelicCardinalityManagementDefaultDestroy verifies the account-wide limit
// was reset to the platform default (100,000) after destroying a DEFAULT-mode resource.
func testAccCheckNewRelicCardinalityManagementDefaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cardinality_management" {
			continue
		}

		accountID, _, err := parseCardinalityLimitID(r.Primary.ID)
		if err != nil {
			return err
		}

		limits, err := client.DataManagement.GetLimitsWithContext(context.Background(), accountID)
		if err != nil {
			return err
		}

		if limits != nil {
			for _, l := range *limits {
				if l.Name == cardinalityLimitName {
					if l.Value != cardinalityLimitPlatformDefault {
						return fmt.Errorf(
							"expected account-wide cardinality limit to be reset to %d after destroy, got %d",
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

// testAccCheckNewRelicCardinalityManagementPerMetricDestroy cannot read back per-metric
// override values via the API, so no meaningful post-destroy assertion is possible.
func testAccCheckNewRelicCardinalityManagementPerMetricDestroy(_ *terraform.State) error {
	return nil
}

// testMetricEntry holds a single metric name + limit for use in test configs.
type testMetricEntry struct {
	name  string
	limit int
}

func testAccNewRelicCardinalityManagementDefaultConfig(limit int) string {
	return fmt.Sprintf(`
resource "newrelic_cardinality_management" "test" {
  mode              = "DEFAULT"
  cardinality_limit = %d
}
`, limit)
}

func testAccNewRelicCardinalityManagementPerMetricConfig(metrics []testMetricEntry) string {
	var blocks string
	for _, m := range metrics {
		blocks += fmt.Sprintf(`
  metric {
    name              = %q
    cardinality_limit = %d
  }`, m.name, m.limit)
	}
	return fmt.Sprintf(`
resource "newrelic_cardinality_management" "test" {
  mode = "PER_METRIC"
%s
}
`, blocks)
}

func testAccNewRelicCardinalityManagementInvalidDefaultWithMetricBlocksConfig() string {
	return `
resource "newrelic_cardinality_management" "test" {
  mode              = "DEFAULT"
  cardinality_limit = 100000
  metric {
    name              = "should.not.be.here"
    cardinality_limit = 50000
  }
}
`
}

func testAccNewRelicCardinalityManagementInvalidDefaultMissingLimitConfig() string {
	return `
resource "newrelic_cardinality_management" "test" {
  mode = "DEFAULT"
}
`
}

func testAccNewRelicCardinalityManagementInvalidPerMetricWithTopLevelLimitConfig() string {
	return `
resource "newrelic_cardinality_management" "test" {
  mode              = "PER_METRIC"
  cardinality_limit = 100000
  metric {
    name              = "some.metric"
    cardinality_limit = 150000
  }
}
`
}

func testAccNewRelicCardinalityManagementInvalidPerMetricNoMetricBlocksConfig() string {
	return `
resource "newrelic_cardinality_management" "test" {
  mode = "PER_METRIC"
}
`
}
