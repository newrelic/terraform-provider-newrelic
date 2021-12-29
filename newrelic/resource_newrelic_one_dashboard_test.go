//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
)

// TestAccNewRelicOneDashboard_CreateOnePage Ensure that we can create a NR1 Dashboard
func TestAccNewRelicOneDashboard_CreateOnePage(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFull(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicOneDashboard_CreateTwoPages Ensure we can create a Two page NR1 Dashboard
func TestAccNewRelicOneDashboard_CreateTwoPages(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_TwoPageBasic(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicOneDashboard_CrossAccountQueries Ensures we can have different account IDs for NRQL queries
func TestAccNewRelicOneDashboard_CrossAccountQueries(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_TwoPageBasic(rName, "1"), // Hard-coded accountID for NRQL queries
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicOneDashboard_PageRename Ensure we can change the name of a NR1 Dashboard
func TestAccNewRelicOneDashboard_PageRename(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	rNameUpdated := fmt.Sprintf("%s-updated", rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFull(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFull(rNameUpdated, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 5), // Sleep waiting for entity re-indexing
				),
			},
		},
	})
}

// TestAccNewRelicOneDashboard_UpdateInvalidNRQL Ensure we catch and display richer error messages on update
func TestAccNewRelicOneDashboard_UpdateInvalidNRQL(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_PageValidNRQL(rName),
				Check:  resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Test: Update
			{
				Config: 		 testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL(rName),
				ExpectError: regexp.MustCompile("Invalid widget input"),
			},
		},
	})
}

// TestAccNewRelicOneDashboard_InvalidNRQL checks for proper response if a widget is not configured correctly
func TestAccNewRelicOneDashboard_InvalidNRQL(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config:      testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL(rName),
				ExpectError: regexp.MustCompile("Invalid widget input"),
			},
		},
	})
}

// TestAccNewRelicOneDashboard_FilterCurrentDashboard Checks if linked_entity_guid is set after updating
func TestAccNewRelicOneDashboard_FilterCurrentDashboard(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicOneDashboardConfig_FilterCurrentDashboard(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_FilterCurrentDashboard("newrelic_one_dashboard.bar", 5),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccNewRelicOneDashboard_RawWidgets checks if raw widgets are accepted and returned correcly
func TestAccNewRelicOneDashboard_RawWidgets(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageRaw(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckNewRelicOneDashboard_FilterCurrentDashboard fetches the dashboard resource after creation, with an optional sleep time
// used when we know the async nature of the API will mess with consistent testing. The filter_current_dashboard requires a second call to update
// the linked_entity_guid to add the page GUID. This also checks to make sure the page GUID matches what has been added.
func testAccCheckNewRelicOneDashboard_FilterCurrentDashboard(name string, sleepSeconds int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

		time.Sleep(time.Duration(sleepSeconds) * time.Second)

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		found, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if string(found.GUID) != rs.Primary.ID {
			return fmt.Errorf("dashboard not found: %v - %v", rs.Primary.ID, found)
		}

		if found.Pages[0].Widgets[0].LinkedEntities == nil {
			return fmt.Errorf("No linked entities found")
		}

		if len(found.Pages[0].Widgets[0].LinkedEntities) > 1 {
			return fmt.Errorf("Greater than 1 linked entity found: %d", len(found.Pages[0].Widgets[0].LinkedEntities))
		}

		if found.Pages[0].Widgets[0].LinkedEntities[0].GetGUID() != found.Pages[0].GUID {
			return fmt.Errorf("Page GUID did not match LinkedEntity: %s", found.Pages[0].Widgets[0].LinkedEntities[0].GetGUID())
		}

		return nil
	}
}

// testAccCheckNewRelicOneDashboardConfig contains raw widget configurations
func testAccCheckNewRelicOneDashboardConfig_OnePageRaw(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"
` + testAccCheckNewRelicOneDashboardConfig_RawWidgets(dashboardName, accountID) + `
}`
}

// testAccCheckNewRelicOneDashboardConfig_TwoPageBasic generates a TF config snippet for a simple
// two page dashboard.
func testAccCheckNewRelicOneDashboardConfig_TwoPageBasic(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageFull(dashboardName, accountID) + `
` + testAccCheckNewRelicOneDashboardConfig_PageSimple("Page 2") + `
}
`
}

// testAccCheckNewRelicOneDashboardConfig contains all the config options for a single page dashboard
func testAccCheckNewRelicOneDashboardConfig_OnePageFull(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageFull(dashboardName, accountID) + `
}`
}

// testAccCheckNewRelicOneDashboardConfig_PageSimple generates a basic dashboard page
func testAccCheckNewRelicOneDashboardConfig_PageSimple(pageName string) string {
	return `
  page {
    name = "` + pageName + `"

    widget_bar {
      title = "foo"
      row = 4
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name"
      }
    }
  }
`
}

func testAccCheckNewRelicOneDashboardConfig_FilterCurrentDashboard(dashboardName string, accountID string) string {
	return `
	resource "newrelic_one_dashboard" "bar" {

		name = "` + dashboardName + `"
	  
		page {
		  name = "` + dashboardName + `"
	  
		  widget_bar {
			title = "Average transaction duration, by application"
			row = 1
			column = 1
	  
			nrql_query {
			  account_id = ` + accountID + `
			  query      = "FROM Transaction SELECT average(duration) FACET appName"
			}
	  
			# Linking to self
			filter_current_dashboard = true
		  }
		}
	  }
	  
`
}

// testAccCheckNewRelicOneDashboardConfig_PageFull generates a TF config snippet that is
// an entire dashboard page, with all widget types
func testAccCheckNewRelicOneDashboardConfig_PageFull(pageName string, accountID string) string {
	return `
  page {
    name = "` + pageName + `"

    widget_area {
      title = "area widget"
      row = 1
      column = 1
      height = 3
      width = 12

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT 51 TIMESERIES"
      }
    }

    widget_bar {
      title = "bar widget"
      row = 4
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name"
      }
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
	}

    widget_billboard {
      title = "billboard widget"
      row = 4
      column = 5
      nrql_query {
        query      = "FROM Transaction SELECT count(*)"
      }

      warning = 0
      critical = 2
    }

    widget_bullet {
      title = "bullet widget"
      row = 4
      column = 9
      limit = 1.5
      nrql_query {
        query  = "FROM Transaction SELECT count(*)"
      }
    }

    widget_funnel {
      title = "funnel widget"
      row = 7
      column = 1
      nrql_query {
        query = "FROM Transaction SELECT funnel(response.status, WHERE name = 'WebTransaction/Expressjs/GET//', WHERE name = 'WebTransaction/Expressjs/GET//api/inventory')"
      }
    }

    widget_heatmap {
      title = "heatmap widget"
      row = 7
      column = 5
      nrql_query {
        query = "FROM Transaction SELECT histogram(duration, buckets: 100, width: 0.1) FACET appName"
      }
    }

    widget_histogram {
      title = "histogram widget"
      row = 7
      column = 9
      nrql_query {
        query = "FROM Transaction SELECT histogram(duration * 100, buckets: 500, width: 1)"
      }
    }

    widget_line {
      title = "line widget"
      row = 10
      column = 1
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT 1 TIMESERIES"
      }
      nrql_query {
        query      = "FROM Transaction SELECT 2 TIMESERIES"
      }
    }

    widget_markdown {
      title = "markdown widget"
      row = 10
      column = 5
      text = "# Header text"
    }

    widget_pie {
      title = "3.14 widget"
      row = 10
      column = 9
      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name"
      }
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_table {
      title = "table widget"
      row = 13
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT average(duration) FACET appName"
      }
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_json {
      title = "JSON widget"
      row = 13
      column = 2
      nrql_query {
        query      = "FROM Transaction SELECT average(duration) FACET appName"
      }
    }

	widget_stacked_bar {
		title = "stacked bar widget"
		row = 14
		column = 1
		nrql_query {
		  query      = "FROM Transaction SELECT average(duration) FACET appName TIMESERIES"
		}
	}
  }
`
}

func testAccCheckNewRelicOneDashboardConfig_PageValidNRQL(dashboardName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

  page {
    name = "` + dashboardName + `"

    widget_line {
      title = "foo"
      row = 1
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT 2 TIMESERIES"
      }
    }
  }
}`
}

// testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL generates a basic dashboard page
func testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL(dashboardName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

  page {
    name = "` + dashboardName + `"

    widget_line {
      title = "foo"
      row = 1
      column = 1
      nrql_query {
        query      = "THIS IS INVALID NRQL"
      }
    }
  }
}`
}

// testAccCheckNewRelicOneDashboardConfig_RawWidgets generates a TF config snippet that is
// an entire dashboard page, with a combination of raw widget types
// we should be able to accept any raw widget, as we don't do any checking on the input types
func testAccCheckNewRelicOneDashboardConfig_RawWidgets(pageName string, accountID string) string {
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
  }
`
}

// testAccCheckNewRelicOneDashboardExists fetches the dashboard back, with an optional sleep time
// used when we know the async nature of the API will mess with consistent testing.
func testAccCheckNewRelicOneDashboardExists(name string, sleepSeconds int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

		time.Sleep(time.Duration(sleepSeconds) * time.Second)

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		found, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if string(found.GUID) != rs.Primary.ID {
			return fmt.Errorf("dashboard not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

// testAccCheckNewRelicOneDashboardDestroy expects the dashboard read to fail,
// and errors if we DO get the dashboard back.
func testAccCheckNewRelicOneDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_one_dashboard" {
			continue
		}

		_, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(r.Primary.ID))
		if err == nil {
			return fmt.Errorf("one_dashboard still exists")
		}

	}
	return nil
}
