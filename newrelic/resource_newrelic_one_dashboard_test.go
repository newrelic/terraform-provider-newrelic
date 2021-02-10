// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
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
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*) FACET name"
			}
			linked_entity_guids =["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_billboard {
      title = "billboard widget"
      row = 4
      column = 5
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*)"
      }

      warning = 1
			critical = 2
    }

    widget_line {
      title = "line widget"
      row = 4
      column = 9
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT 1 TIMESERIES"
      }
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT 2 TIMESERIES"
			}
			linked_entity_guids =["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_markdown {
      title = "markdown widget"
      row = 7
      column = 1
			text = "# Header text"
    }

    widget_pie {
      title = "3.14 widget"
      row = 7
      column = 5
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*) FACET name"
			}
			linked_entity_guids =["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_table {
      title = "table widget"
      row = 7
      column = 9
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT *"
			}
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

		found, err := client.Dashboards.GetDashboardEntity(entities.EntityGUID(rs.Primary.ID))
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

		_, err := client.Dashboards.GetDashboardEntity(entities.EntityGUID(r.Primary.ID))
		if err == nil {
			return fmt.Errorf("one_dashboard still exists")
		}

	}
	return nil
}
