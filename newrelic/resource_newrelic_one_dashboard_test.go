//go:build integration || DASHBOARDS
// +build integration DASHBOARDS

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
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
				Config: testAccCheckNewRelicOneDashboardConfig_TwoPageBasic(rName, "3814156"), // Hard-coded accountID for NRQL queries
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
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 60), // Sleep waiting for entity re-indexing
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Test: Update
			{
				Config:      testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL(rName),
				ExpectError: regexp.MustCompile(`query error; Input: \[\[PageInput: \[\[WidgetInput: \[Error parsing 'THIS IS INVALID NRQL', cause: Invalid nrql query]]]]]`),
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
				ExpectError: regexp.MustCompile(`query error; Input: \[\[PageInput: \[\[WidgetInput: \[Error parsing 'THIS IS INVALID NRQL', cause: Invalid nrql query]]]]]`),
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
				Config: testAccCheckNewRelicOneDashboardConfig_FilterCurrentDashboard(rName, strconv.Itoa(testAccountID), "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_FilterCurrentDashboard("newrelic_one_dashboard.bar"),
				),
			},
		},
	})
}

// TestAccNewRelicOneDashboard_BillboardThresholds Checks if critical and warning are set correctly for billboard widget
func TestAccNewRelicOneDashboard_BillboardThresholds(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	rWidgetName := fmt.Sprintf("tf-test-widget-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicOneDashboardConfig_BillboardWithThresholds(rName, rWidgetName, 100, 200),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_BillboardCriticalWarning("newrelic_one_dashboard.bar", rWidgetName, false, 100, 200),
				),
			},
			{
				Config: testAccCheckNewRelicOneDashboardConfig_BillboardWithThresholds(rName, rWidgetName, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_BillboardCriticalWarning("newrelic_one_dashboard.bar", rWidgetName, false, 0, 0),
				),
			},
			{
				Config: testAccCheckNewRelicOneDashboardConfig_BillboardWithoutThresholds(rName, rWidgetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_BillboardCriticalWarning("newrelic_one_dashboard.bar", rWidgetName, true, 0, 0),
				),
			},
		},
	})
}

func TestAccNewRelicOneDashboard_UnlinkFilterCurrentDashboard(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicOneDashboardConfig_FilterCurrentDashboard(rName, strconv.Itoa(testAccountID), "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_FilterCurrentDashboard("newrelic_one_dashboard.bar"),
				),
			},
			{
				Config: testAccCheckNewRelicOneDashboardConfig_FilterCurrentDashboard(rName, strconv.Itoa(testAccountID), "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboard_UnlinkFilterCurrentDashboard("newrelic_one_dashboard.bar"),
				),
			},
		},
	})
}

// TestAccNewRelicOneDashboard_ChangeCheck Ensures that all changes are coming through well
func TestAccNewRelicOneDashboard_ChangeCheck(t *testing.T) {
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
			// Make lots of changes
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFullChanged(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
		},
	})
}

// TestAccNewRelicOneDashboard_EmptyPage tests successful creation of a dashboard comprising a page with no widgets
func TestAccNewRelicOneDashboard_EmptyPage(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_EmptyPage(rName),
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

// TestAccNewRelicOneDashboard_Tooltip tests creating and updating dashboards with tooltip configuration
func TestAccNewRelicOneDashboard_Tooltip(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create dashboard with tooltip
			{
				Config: testAccCheckNewRelicOneDashboardConfig_WithTooltip(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.tooltip_test", 0),
				),
			},
			// Update tooltip configuration
			{
				Config: testAccCheckNewRelicOneDashboardConfig_WithTooltipUpdated(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.tooltip_test", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard.tooltip_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicOneDashboard_InvalidTooltipMode checks for proper response if a tooltip is not configured correctly
func TestAccNewRelicOneDashboard_InvalidTooltipMode(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicOneDashboardConfig_InvalidTooltipNRQL(rName),
				ExpectError: regexp.MustCompile(`expected page.0.widget_line.0.tooltip.0.mode to be one of \[all single hidden\], got invalid`),
			},
		},
	})
}

// testAccCheckNewRelicOneDashboard_FilterCurrentDashboard fetches the dashboard resource after creation, with an optional sleep time
// used when we know the async nature of the API will mess with consistent testing. The filter_current_dashboard requires a second call to update
// the linked_entity_guid to add the page GUID. This also checks to make sure the page GUID matches what has been added.
func testAccCheckNewRelicOneDashboard_FilterCurrentDashboard(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		retryErr := resource.RetryContext(context.Background(), 5*time.Second, func() *resource.RetryError {
			found, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(rs.Primary.ID))
			if err != nil {
				return resource.RetryableError(err)
			}

			if string(found.GUID) != rs.Primary.ID {
				return resource.RetryableError(fmt.Errorf("dashboard not found: %v - %v", rs.Primary.ID, found))
			}

			if found.Pages[0].Widgets[0].LinkedEntities == nil {
				return resource.NonRetryableError(fmt.Errorf("No linked entities found"))
			}

			if len(found.Pages[0].Widgets[0].LinkedEntities) > 1 {
				return resource.NonRetryableError(fmt.Errorf("Greater than 1 linked entity found: %d", len(found.Pages[0].Widgets[0].LinkedEntities)))
			}

			if found.Pages[0].Widgets[0].LinkedEntities[0].GetGUID() != found.Pages[0].GUID {
				return resource.NonRetryableError(fmt.Errorf("Page GUID did not match LinkedEntity: %s", found.Pages[0].Widgets[0].LinkedEntities[0].GetGUID()))
			}

			return nil
		})

		if retryErr != nil {
			return retryErr
		}

		return nil
	}
}

// helper function to check if the values of critical and warning are set correctly for the Billboard widget type
func testAccCheckNewRelicOneDashboard_BillboardCriticalWarning(resourceName string, widgetTitle string, empty bool, critical float64, warning float64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		retryErr := resource.RetryContext(context.Background(), 5*time.Second, func() *resource.RetryError {
			found, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(rs.Primary.ID))
			if err != nil {
				return resource.RetryableError(err)
			}

			if string(found.GUID) != rs.Primary.ID {
				return resource.RetryableError(fmt.Errorf("dashboard not found: %v - %v", rs.Primary.ID, found))
			}

			foundWidget := false
			for _, page := range found.Pages {
				for _, widget := range page.Widgets {
					if widget.Title == widgetTitle {
						foundWidget = true
						if empty {
							if len(widget.Configuration.Billboard.Thresholds) > 0 {
								return resource.NonRetryableError(fmt.Errorf("Found thresholds on billboard, but none should be set: %s", widgetTitle))
							}
						} else {
							for _, threshold := range widget.Configuration.Billboard.Thresholds {
								if threshold.AlertSeverity == entities.DashboardAlertSeverityTypes.CRITICAL && threshold.Value != critical {
									return resource.NonRetryableError(fmt.Errorf("The value of critical is incorrect for widget: %s", widgetTitle))
								}
								if threshold.AlertSeverity == entities.DashboardAlertSeverityTypes.WARNING && threshold.Value != warning {
									return resource.NonRetryableError(fmt.Errorf("The value of warning is incorrect for widget: %s", widgetTitle))
								}
							}
						}
					}
				}
			}

			if !foundWidget {
				return resource.NonRetryableError(fmt.Errorf("Unable to find widget: %s", widgetTitle))
			}

			return nil
		})

		if retryErr != nil {
			return retryErr
		}

		return nil
	}
}

// testAccCheckNewRelicOneDashboard_UnlinkFilterCurrentDashboard fetches the dashboard resource after update
// and checks that entities were unlinked
func testAccCheckNewRelicOneDashboard_UnlinkFilterCurrentDashboard(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		found, err := client.Dashboards.GetDashboardEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if string(found.GUID) != rs.Primary.ID {
			return fmt.Errorf("dashboard not found: %v - %v", rs.Primary.ID, found)
		}

		if found.Pages[0].Widgets[0].LinkedEntities != nil {
			return fmt.Errorf("Entities still linked")
		}

		return nil
	}
}

func TestAccNewRelicOneDashboard_VariablesNRQL(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesNRQL(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Test: Update
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesNRQLUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_one_dashboard.bar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"variable.0.options",
				},
			},
		},
	})
}

func TestAccNewRelicOneDashboard_VariablesEnum(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesEnum(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			// Test: Update
			{
				Config: testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesEnumUpdated(rName),
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

// testAccCheckNewRelicOneDashboardConfig_EmptyPage contains the configuration needed to create a dashboard
// with a single page comprising no widgets
func testAccCheckNewRelicOneDashboardConfig_EmptyPage(dashboardName string) string {
	return `
		resource "newrelic_one_dashboard" "bar" {
  			name = "` + dashboardName + `"
  			permissions = "private"
  			page {
    			name = "` + dashboardName + `_page_one"
  			}
		}`
}

func testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesNRQL(dashboardName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageSimple(dashboardName) + `
` + testAccCheckNewRelicOneDashboardConfig_VariableNRQL() + `
}`
}

func testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesNRQLUpdated(dashboardName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageSimple(dashboardName) + `
` + testAccCheckNewRelicOneDashboardConfig_VariableNRQLUpdated() + `
}`
}

func testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesEnum(dashboardName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageSimple(dashboardName) + `
` + testAccCheckNewRelicOneDashboardConfig_VariableEnum() + `
}`
}

func testAccCheckNewRelicOneDashboardConfig_OnePageFullVariablesEnumUpdated(dashboardName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageSimple(dashboardName) + `
` + testAccCheckNewRelicOneDashboardConfig_VariableEnumUpdated() + `
}`
}

// testAccCheckNewRelicOneDashboardConfig_OnePageFullChanged contains all the config options for a single page dashboard with lots of
// changes compared to testAccCheckNewRelicOneDashboardConfig_OnePageFull
func testAccCheckNewRelicOneDashboardConfig_OnePageFullChanged(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

` + testAccCheckNewRelicOneDashboardConfig_PageFullChanged(dashboardName, accountID) + `
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
      refresh_rate = 30000
      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name"
      }
    }
  }
`
}

func testAccCheckNewRelicOneDashboardConfig_FilterCurrentDashboard(dashboardName string, accountID string, filterDashboard string) string {
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
			filter_current_dashboard = ` + filterDashboard + `
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

      refresh_rate = 30000

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
      data_format {
		name = "count"
		type = "decimal"
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
	  ignore_time_range = true
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
	  y_axis_left_zero = false
	  y_axis_right {
		y_axis_right_zero = false
		y_axis_right_series = ["A", "B"]
      }
	  threshold {
        name     = "Duration Threshold"
        from     = 10.8 
        to       = 30.7
        severity = "critical"
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

    widget_log_table {
      title = "Log table widget"
      row = 13
      column = 1
      nrql_query {
        query      = "FROM Log SELECT *"
      }
    }

	widget_table {
      title = "table widget"
      row = 13
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT average(duration) FACET appName"
      }
	  threshold {
		from 		= 100.1
		to 			= 200.2
		column_name = "C1"
		severity 	= "unavailable"
	  }
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
      refresh_rate = 30000
      initial_sorting {
        direction = "desc"
        name      = "appName"
      }
      data_format {
		name = "Avg duration"
		type = "decimal"
      }
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

func testAccCheckNewRelicOneDashboardConfig_VariableNRQL() string {
	return `
  variable {
    default_values = ["value"]
	is_multi_selection = true
	item {
		title = "item"
		value = "ITEM"
	}
    name = "variable"
	nrql_query {
		account_ids = [3806526]
		query = "FROM Transaction SELECT average(duration) FACET appName"
	}
	replacement_strategy = "default"
	title = "title"
	type = "nrql"
	options {
		excluded = true
		ignore_time_range = true
	}
  }
`
}

func testAccCheckNewRelicOneDashboardConfig_VariableNRQLUpdated() string {
	return `
  variable {
    default_values = ["value"]
	is_multi_selection = true
	item {
		title = "item"
		value = "ITEM"
	}
    name = "variableUpdated"
	nrql_query {
		account_ids = [3806526, 3814156]
		query = "FROM Transaction SELECT average(duration) FACET appName"
	}
	replacement_strategy = "default"
	title = "title"
	type = "nrql"
	options {
		excluded = false
		ignore_time_range = false
	}
  }
`
}

func testAccCheckNewRelicOneDashboardConfig_VariableEnum() string {
	return `
  variable {
	is_multi_selection = true
	item {
		title = "item"
		value = "ITEM"
	}
    name = "variable"
	replacement_strategy = "default"
	title = "title"
	type = "enum"
  }
`
}

func testAccCheckNewRelicOneDashboardConfig_VariableEnumUpdated() string {
	return `
  variable {
	default_values = ["default"]
	is_multi_selection = true
	item {
		title = "item"
		value = "ITEM"
	}
    name = "variableUpdated"
	replacement_strategy = "default"
	title = "title"
	type = "enum"
  }
`
}

// testAccCheckNewRelicOneDashboardConfig_PageFullChanged generates a TF config snippet that is
// an entire dashboard page, with all widget types and lots of changes compared to testAccCheckNewRelicOneDashboardConfig_PageFull
func testAccCheckNewRelicOneDashboardConfig_PageFullChanged(pageName string, accountID string) string {
	return `
  page {
    name = "` + pageName + `"

    widget_area {
      title = "area widget with new name"
      row = 1
      column = 1
      height = 4
      width = 12

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT 51 TIMESERIES LIMIT 10"
      }
    }

    widget_bar {
      title = "bar widget with new name"
      row = 2
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name LIMIT 10"
      }
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
	}

    widget_billboard {
      title = "billboard widget with new name"
      row = 4
      column = 5
      nrql_query {
        query      = "FROM Transaction SELECT count(*) LIMIT 10"
      }

      warning = -1
      critical = 0
    }

    widget_bullet {
      title = "bullet widget with new name"
      row = 4
      column = 9
      limit = 1
      nrql_query {
        query  = "FROM Transaction SELECT count(*) LIMIT 10"
      }
    }

    widget_funnel {
      title = "funnel widget with new name"
      row = 7
      column = 1
      nrql_query {
        query = "FROM Transaction SELECT funnel(response.status, WHERE name = 'WebTransaction/Expressjs/GET//', WHERE name = 'WebTransaction/Expressjs/GET//api/inventory') LIMIT 10"
      }
	  ignore_time_range = true
    }

    widget_heatmap {
      title = "heatmap widget with new name"
      row = 7
      column = 5
      nrql_query {
        query = "FROM Transaction SELECT histogram(duration, buckets: 100, width: 0.1) FACET appName"
      }
    }

    widget_histogram {
      title = "histogram widgetw with new name"
      row = 7
      column = 9
      nrql_query {
        query = "FROM Transaction SELECT histogram(duration * 100, buckets: 500, width: 1)"
      }
    }

    widget_line {
      title = "line widget with new name"
      row = 10
      column = 1
      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT 1 TIMESERIES LIMIT 10"
      }
	  nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name TIMESERIES LIMIT 10"
      }
      y_axis_left_zero = true
      y_axis_left_max = 25
	  is_label_visible = true
	  threshold {
		from     = 100.1
		to       = 200.2
		name     = "T1"
		severity = "warning"
      }
      threshold {
		from     = 201
		to       = 300
		name     = "T1 Critical"
		severity = "critical"
      }
      y_axis_right {
		y_axis_right_zero = true
		y_axis_right_series = ["A", "B"]
		y_axis_right_min = 0
		y_axis_right_max = 150
      }
    }

    widget_markdown {
      title = "markdown widget with new name"
      row = 10
      column = 5
      text = "# Header text"
    }

    widget_pie {
      title = "pizza 3.14159 widget with new name"
      row = 10
      column = 9
      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name LIMIT 10"
      }
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_table {
      title = "table widget with new name"
      row = 13
      column = 1
      nrql_query {
        query      = "FROM Transaction SELECT average(duration) FACET appName LIMIT 10"
      }
	  threshold {
		from 		= 100.2
		to 			= 200.4
		column_name = "C1"
		severity 	= "unavailable"
	  }		
      linked_entity_guids = ["MjUyMDUyOHxWSVp8REFTSEJPQVJEfDE2NDYzMDQ"]
    }

    widget_log_table {
      title = "Log table widget with a new name"
      row = 12
      column = 7
      nrql_query {
        query      = "SELECT * FROM Log"
      }
    }

    widget_json {
      title = "JSON widget parsed from yaml, generated from ini"
      row = 13
      column = 2
      nrql_query {
        query      = "FROM Transaction SELECT average(duration) FACET appName LIMIT 10"
      }
    }

	widget_stacked_bar {
		title = "stacked bar widget with new name"
		row = 14
		column = 1
		nrql_query {
		  query      = "FROM Transaction SELECT average(duration) FACET appName TIMESERIES LIMIT 10"
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

// testAccCheckNewRelicOneDashboardConfig_InvalidTooltipNRQL generates a basic dashboard with Invalid Tooltip Configuration
func testAccCheckNewRelicOneDashboardConfig_InvalidTooltipNRQL(dashboardName string) string {
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
      tooltip {
        mode = "invalid_mode"
      }
    }
  }
}`
}

// testAccCheckNewRelicOneDashboardConfig_WithTooltip generates a basic dashboard with Valid Tooltip Configuration
func testAccCheckNewRelicOneDashboardConfig_WithTooltip(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "tooltip_test" {
  name = "` + dashboardName + `"
  permissions = "private"

  page {
    name = "` + dashboardName + `_page"

    widget_line {
      title = "Line widget with tooltip"
      row = 1
      column = 1
      height = 3
      width = 6

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT average(duration) TIMESERIES"
      }

      tooltip {
        mode = "all"
      }
    }

    widget_area {
      title = "Area widget with tooltip"
      row = 1
      column = 7
      height = 3
      width = 6

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*) TIMESERIES"
      }

      tooltip {
	    mode = "all"
	  }
    }

    widget_stacked_bar {
      title = "Bar widget with tooltip"
      row = 4
      column = 1
      height = 3
      width = 12

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*) FACET name"
      }

      tooltip {
	    mode = "all"
	  }
    }
  }
}`
}

// testAccCheckNewRelicOneDashboardConfig_WithTooltipUpdated generates a basic dashboard with Updated Tooltip Configuration
func testAccCheckNewRelicOneDashboardConfig_WithTooltipUpdated(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "tooltip_test" {
  name = "` + dashboardName + `"
  permissions = "private"

  page {
    name = "` + dashboardName + `_page"

    widget_line {
      title = "Line widget with tooltip"
      row = 1
      column = 1
      height = 3
      width = 6

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT average(duration) TIMESERIES"
      }

      tooltip {
        mode = "single"
      }
    }

    widget_area {
      title = "Area widget with tooltip"
      row = 1
      column = 7
      height = 3
      width = 6

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*) TIMESERIES"
      }

      tooltip {
        mode = "single"
      }
    }

    widget_stacked_bar {
      title = "Bar widget with tooltip"
      row = 4
      column = 1
      height = 3
      width = 12

      nrql_query {
        account_id = ` + accountID + `
        query      = "FROM Transaction SELECT count(*) FACET name"
      }

      tooltip {
        mode = "single"
      }
    }
  }
}`
}

// testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL generates billboard with critical and warning set
func testAccCheckNewRelicOneDashboardConfig_BillboardWithThresholds(dashboardName string, widgetName string, critical float64, warning float64) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"

  page {
    name = "` + dashboardName + `"

    widget_billboard {
      title = "` + widgetName + `"
      row = 1
      column = 1
      nrql_query {
        query      = "SELECT count(*) FROM ProcessSample SINCE 30 MINUTES AGO TIMESERIES"
      }
      critical = ` + strconv.FormatFloat(critical, 'f', -1, 64) + `
      warning = ` + strconv.FormatFloat(warning, 'f', -1, 64) + `
    }
  }
}`
}

// testAccCheckNewRelicOneDashboardConfig_PageInvalidNRQL generates billboard without critical and warning set
func testAccCheckNewRelicOneDashboardConfig_BillboardWithoutThresholds(dashboardName string, widgetName string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"

  page {
    name = "` + dashboardName + `"

    widget_billboard {
      title = "` + widgetName + `"
      row = 1
      column = 1
      nrql_query {
        query      = "SELECT count(*) FROM ProcessSample SINCE 30 MINUTES AGO TIMESERIES"
      }
    }
  }
}`
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

func TestAccNewRelicOneDashboard_FilterCurrentDashboardPageChange(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicOneDashboard_FilterCurrentDashboardPageChange_TwoPages(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			{
				Config: testAccNewRelicOneDashboard_FilterCurrentDashboardPageChange_OnePage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
			{
				Config: testAccNewRelicOneDashboard_FilterCurrentDashboardPageChange_TwoPages(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar", 0),
				),
			},
		},
	})
}

func testAccNewRelicOneDashboard_FilterCurrentDashboardPageChange_OnePage(dashboardName string) string {
	return `
	resource "newrelic_one_dashboard" "bar" {
	  name = "` + dashboardName + `"
	  permissions = "private"

	  page {
		name = "` + dashboardName + ` page one"

		widget_bar {
			title = "Average transaction duration, by application"
			row = 1
			column = 1

			nrql_query {
			  query      = "FROM Transaction SELECT average(duration) FACET appName"
			}

			# Linking to self
			filter_current_dashboard = true
		  }
	  }
	}`
}

func testAccNewRelicOneDashboard_FilterCurrentDashboardPageChange_TwoPages(dashboardName string) string {
	return `
	resource "newrelic_one_dashboard" "bar" {
	  name = "` + dashboardName + `"
	  permissions = "private"

	  page {
		name = "` + dashboardName + ` page one"

		widget_bar {
			title = "Average transaction duration, by application"
			row = 1
			column = 1

			nrql_query {
			  query      = "FROM Transaction SELECT average(duration) FACET appName"
			}

			# Linking to self
			filter_current_dashboard = true
		  }
	  }
	  page {
		name = "` + dashboardName + ` page two"

		widget_bar {
			title = "A different widget name cause debugging was getting confusing :)"
			row = 1
			column = 1

			nrql_query {
			  query      = "FROM Transaction SELECT average(duration) FACET appName"
			}

			# Linking to self
			filter_current_dashboard = true
		  }
	  }
	}`
}

func TestAccNewRelicOneDashboard_CreateOnePage_RawConfig(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicOneDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig_RawConfig(rName, strconv.Itoa(testAccountID)),
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

func testAccCheckNewRelicOneDashboardConfig_RawConfig(dashboardName string, accountID string) string {
	return `
	resource "newrelic_one_dashboard" "bar" {
	  name = "` + dashboardName + `"
	  permissions = "private"

	  page {
		name = "` + dashboardName + ` page one"

		widget_area {
      title = "Heap Utilization"
      row = 1
      column = 1
      height = 3
      width = 4

      nrql_query {
      account_id = ` + accountID + `
        query = <<EOT
FROM JVMSampleActiveMQ SELECT latest(HeapMemoryUsage.Used) / 1000, latest(HeapMemoryUsage.Max) / 1000, latest(HeapMemoryUsage.Committed) / 1000, latest(HeapMemoryUsage.Init) / 1000 TIMESERIES AUTO SINCE 3 month ago
EOT
      }
      facet_show_other_series = true
      legend_enabled = true
      ignore_time_range = true
      y_axis_left_min = 0
      y_axis_left_max = 0

      null_values {
        null_value = "preserve"
        series_overrides {
          null_value = "default"
          series_name = "Heap Memory Usage. Used"
        }
        series_overrides {
          null_value = "zero"
          series_name = "Heap Memory Usage. Max"
        }
      }

      units {
        unit = "BYTES"
        series_overrides {
          unit = "BYTES"
          series_name = "Heap Memory Usage. Committed"
        }
      }

      colors {
        color = "#722727"
        series_overrides {
          color = "#722727"
          series_name = "Heap Memory Usage. Used"
        }
        series_overrides {
          color = "#236f70"
          series_name = "Heap Memory Usage. Max"
        }
      }
    }
	  }
	}`
}
