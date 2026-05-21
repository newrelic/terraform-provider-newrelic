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

func TestAccNewRelicAccountCardinalityLimit_Default(t *testing.T) {
	resourceName := "newrelic_account_cardinality_limit.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAccountCardinalityLimitDefaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountCardinalityLimitDefaultConfig(200000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "cardinality_limit", "200000"),
					resource.TestCheckResourceAttr(resourceName, "metric_name", ""),
				),
			},
			// Update cardinality_limit in-place (no ForceNew on that field).
			{
				Config: testAccNewRelicAccountCardinalityLimitDefaultConfig(250000),
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

func TestAccNewRelicAccountCardinalityLimit_PerMetric(t *testing.T) {
	resourceName := "newrelic_account_cardinality_limit.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAccountCardinalityLimitPerMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountCardinalityLimitPerMetricConfig(
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

func TestAccNewRelicAccountCardinalityLimit_InvalidDefaultWithMetricName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAccountCardinalityLimitInvalidDefaultWithMetricNameConfig(),
				ExpectError: regexp.MustCompile(`metric_name must not be set when mode is "DEFAULT"`),
			},
		},
	})
}

func TestAccNewRelicAccountCardinalityLimit_InvalidPerMetricWithoutMetricName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAccountCardinalityLimitInvalidPerMetricWithoutMetricNameConfig(),
				ExpectError: regexp.MustCompile(`metric_name is required when mode is "PER_METRIC"`),
			},
		},
	})
}

// testAccCheckNewRelicAccountCardinalityLimitDefaultDestroy verifies the account-wide limit
// was reset to the platform default (100,000) after destroying a DEFAULT-mode resource.
func testAccCheckNewRelicAccountCardinalityLimitDefaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_account_cardinality_limit" {
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

// testAccCheckNewRelicAccountCardinalityLimitPerMetricDestroy cannot read back per-metric
// override values via the API, so no meaningful post-destroy assertion is possible here.
func testAccCheckNewRelicAccountCardinalityLimitPerMetricDestroy(_ *terraform.State) error {
	return nil
}

func testAccNewRelicAccountCardinalityLimitDefaultConfig(limit int) string {
	return fmt.Sprintf(`
resource "newrelic_account_cardinality_limit" "test" {
  mode              = "DEFAULT"
  cardinality_limit = %d
}
`, limit)
}

func testAccNewRelicAccountCardinalityLimitPerMetricConfig(metricName string, limit int) string {
	return fmt.Sprintf(`
resource "newrelic_account_cardinality_limit" "test" {
  mode              = "PER_METRIC"
  metric_name       = %q
  cardinality_limit = %d
}
`, metricName, limit)
}

func testAccNewRelicAccountCardinalityLimitInvalidDefaultWithMetricNameConfig() string {
	return `
resource "newrelic_account_cardinality_limit" "test" {
  mode              = "DEFAULT"
  metric_name       = "should.not.be.set"
  cardinality_limit = 100000
}
`
}

func testAccNewRelicAccountCardinalityLimitInvalidPerMetricWithoutMetricNameConfig() string {
	return `
resource "newrelic_account_cardinality_limit" "test" {
  mode              = "PER_METRIC"
  cardinality_limit = 100000
}
`
}
