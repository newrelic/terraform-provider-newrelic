package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicDashboard_Basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	rNameUpdated := fmt.Sprintf("%s-updated", rName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			// Check exists
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", rName),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "editable", "editable_by_all"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "icon", "bar-chart"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "visibility", "all"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.754238675.title", "Average Transaction Duration"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.754238675.height", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.754238675.width", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.754238675.row", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.754238675.column", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.754238675.visualization", "faceted_line_chart"),
				),
			},
			// Update dashboard title
			{
				Config: testAccCheckNewRelicDashboardConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", rNameUpdated),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "filter.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "filter.0.event_types.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "filter.0.event_types.4104882694", "Transaction"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "filter.0.attributes.#", "2"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "filter.0.attributes.2634578693", "appName"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "filter.0.attributes.3755723101", "envName"),
				),
			},
			// Add widget
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigAdded(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", rNameUpdated),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "2"),
				),
			},
			// Update widget
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", rNameUpdated),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "2"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.3380029573.nrql", "SELECT PERCENTILE(duration, 50) from Transaction FACET appName TIMESERIES auto"),
				),
			},
		},
	})
}

func TestAccNewRelicDashboard_import(t *testing.T) {
	resourceName := "newrelic_dashboard.foo"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNewRelicDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_dashboard" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = client.GetDashboard(int(id))

		if err == nil {
			return fmt.Errorf("Dashboard still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicDashboardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No dashboard ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Client

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.GetDashboard(int(id))
		if err != nil {
			return err
		}

		if strconv.Itoa(found.ID) != rs.Primary.ID {
			return fmt.Errorf("Dashboard not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicDashboardWidgetConfigAdded(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title                = "%s"
  filter {
    event_types = [
        "Transaction"
    ]
    attributes = [
        "appName",
        "envName"
    ]
  }
  widget {
    title         = "Average Transaction Duration"
    visualization = "faceted_line_chart"
    column        = 1
    row           = 1
    nrql          = "SELECT PERCENTILE(duration, 95) from Transaction FACET appName TIMESERIES auto"
  }
  widget {
    title         = "Page Views"
	visualization = "faceted_line_chart"
	column        = 1
	row           = 2
    nrql          = "SELECT AVERAGE(duration) from PageView FACET appName TIMESERIES auto"
  }
}
`, rName)
}

func testAccCheckNewRelicDashboardWidgetConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title                = "%s"
  filter {
    event_types = [
        "Transaction"
    ]
    attributes = [
        "appName",
        "envName"
    ]
  }
  widget {
    title         = "Average Transaction Duration"
    visualization = "faceted_line_chart"
    column        = 1
    row           = 1
    nrql          = "SELECT PERCENTILE(duration, 50) from Transaction FACET appName TIMESERIES auto"
  }
  widget {
    title         = "Page Views"
    visualization = "faceted_line_chart"
    column        = 1
    row           = 2
    nrql          = "SELECT AVERAGE(duration) from PageView FACET appName TIMESERIES auto"
  }
}
`, rName)
}

func testAccCheckNewRelicDashboardConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title                = "%s"
  filter {
    event_types = [
        "Transaction"
    ]
    attributes = [
        "appName",
        "envName"
    ]
  }
  widget {
    title         = "Average Transaction Duration"
    visualization = "faceted_line_chart"
    column        = 1
    row           = 1
    nrql          = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
  }
}
`, rName)
}

func testAccCheckNewRelicDashboardConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title = "%s"

  widget {
    title         = "Average Transaction Duration"
    visualization = "faceted_line_chart"
    column        = 1
    row           = 1
    nrql          = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
  }
}
`, rName)
}
