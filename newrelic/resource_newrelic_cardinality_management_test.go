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
					resource.TestCheckResourceAttr(resourceName, "metric_name", ""),
				),
			},
			// Update cardinality_limit in-place (no ForceNew on that field).
			{
				Config: testAccNewRelicCardinalityManagementDefaultConfig(250000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "250000"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_PerMetric(t *testing.T) {
	resourceName := "newrelic_cardinality_management.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCardinalityManagementPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCardinalityManagementPerMetricConfig(
					"test.cardinality.metric.tf",
					150000,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "PER_METRIC"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", "test.cardinality.metric.tf"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "150000"),
				),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_InvalidDefaultWithMetricName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementInvalidDefaultWithMetricNameConfig(),
				ExpectError: regexp.MustCompile(`metric_name must not be set when mode is "DEFAULT"`),
			},
		},
	})
}

func TestAccNewRelicCardinalityManagement_InvalidPerMetricWithoutMetricName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicCardinalityManagementInvalidPerMetricWithoutMetricNameConfig(),
				ExpectError: regexp.MustCompile(`metric_name is required when mode is "PER_METRIC"`),
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
// override values via the API, so no meaningful post-destroy assertion is possible here.
func testAccCheckNewRelicCardinalityManagementPerMetricDestroy(_ *terraform.State) error {
	return nil
}

func testAccNewRelicCardinalityManagementDefaultConfig(limit int) string {
	return fmt.Sprintf(`
resource "newrelic_cardinality_management" "test" {
  mode              = "DEFAULT"
  cardinality_limit = %d
}
`, limit)
}

func testAccNewRelicCardinalityManagementPerMetricConfig(metricName string, limit int) string {
	return fmt.Sprintf(`
resource "newrelic_cardinality_management" "test" {
  mode              = "PER_METRIC"
  metric_name       = %q
  cardinality_limit = %d
}
`, metricName, limit)
}

func testAccNewRelicCardinalityManagementInvalidDefaultWithMetricNameConfig() string {
	return `
resource "newrelic_cardinality_management" "test" {
  mode              = "DEFAULT"
  metric_name       = "should.not.be.set"
  cardinality_limit = 100000
}
`
}

func testAccNewRelicCardinalityManagementInvalidPerMetricWithoutMetricNameConfig() string {
	return `
resource "newrelic_cardinality_management" "test" {
  mode              = "PER_METRIC"
  cardinality_limit = 100000
}
`
}
