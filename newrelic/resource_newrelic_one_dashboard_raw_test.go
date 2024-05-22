//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

// TestAccNewRelicOneDashboardRaw_CreateOnePage Ensure that we can create a NR1 Dashboard
func TestAccNewRelicOneDashboardRaw_CreateOnePage(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardRawConfig_OnePageFull(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard_raw.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard_raw.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicOneDashboardRaw_EmptyPage tests successful creation of a dashboard comprising a page with no widgets
func TestAccNewRelicOneDashboardRaw_EmptyPage(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardRawConfig_EmptyPage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard_raw.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard_raw.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckNewRelicOneDashboardRawConfig_EmptyPage contains the configuration to test the creation of a dashboard
// with a page comprising no widgets
func testAccCheckNewRelicOneDashboardRawConfig_EmptyPage(dashboardName string) string {
	return `
	resource "newrelic_one_dashboard_raw" "bar" {
  		name = "` + dashboardName + `"
  		permissions = "private"
  		page {
  			name = "` + dashboardName + `_page_one"
  		}
	}`
}

// testAccCheckNewRelicOneDashboardRawConfig_OnePageFull contains all the config options for a single page dashboard
func testAccCheckNewRelicOneDashboardRawConfig_OnePageFull(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard_raw" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardRawConfig_PageFull(dashboardName, accountID) + `
}`
}

// testAccCheckNewRelicOneDashboardRawConfig_PageFull generates a TF config snippet that is
// an entire dashboard page, with all widget types
func testAccCheckNewRelicOneDashboardRawConfig_PageFull(pageName string, accountID string) string {
	return `
  page {
    name = "` + pageName + `"
    widget {
      title = "Custom widget"
      row = 1
      column = 1
      width = 1
      height = 1
      visualization_id = "viz.custom"
      configuration = <<EOT
      {
        "legend": {
          "enabled": false
        },
        "nrqlQueries": [
          {
            "accountId": ` + accountID + `,
            "query": "SELECT average(loadAverageOneMinute), average(loadAverageFiveMinute), average(loadAverageFifteenMinute) from SystemSample SINCE 60 minutes ago    TIMESERIES"
          }
        ],
        "yAxisLeft": {
          "max": 100,
          "min": 50,
          "zero": false
        }
      }
      EOT
    }
    widget {
      title = "Server CPU"
      row = 1
      column = 2
      width = 1
      height = 1
      visualization_id = "viz.testing"
      configuration = <<EOT
      {
        "nrqlQueries": [
          {
            "accountId": ` + accountID + `,
            "query": "SELECT average(cpuPercent) FROM SystemSample since 3 hours ago facet hostname limit 400"
          }
        ]
      }
      EOT
    }
    widget {
      title  = "Docker Server CPU"
      row    = 1
      column = 3
      height = 1
      width  = 1
      visualization_id = "viz.bar"
      configuration = jsonencode(
        {
          "facet": {
            "showOtherSeries": false
          },
          "nrqlQueries": [
            {
	          "accountId": ` + accountID + `,
	          "query": "SELECT average(cpuPercent) FROM SystemSample since 3 hours ago facet hostname limit 400"
            }
          ]
        }
      )
      linked_entity_guids = ["MzgwNjUyNnxWSVp8REFTSEJPQVJEfGRhOjI5Mjk0MjA"]
    }
  }
`
}

// testAccCheckNewRelicOneDashboardDestroy expects the dashboard read to fail,
// and errors if we DO get the dashboard back.
func testAccCheckNewRelicOneDashboardRawDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_one_dashboard_raw" {
			continue
		}

		_, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(r.Primary.ID))
		if err == nil {
			return fmt.Errorf("one_dashboard still exists")
		}

	}
	return nil
}
