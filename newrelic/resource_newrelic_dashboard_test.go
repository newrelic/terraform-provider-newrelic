package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	newrelic "github.com/paultyng/go-newrelic/api"
)

func TestAccNewRelicDashboard_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			// Check exists
			resource.TestStep{
				Config: testAccCheckNewRelicDashboardConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "editable", "editable_by_all"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "icon", "bar-chart"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "visibility", "all"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.2870742112.title", "Average Transaction Duration"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.2870742112.height", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.2870742112.width", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.2870742112.row", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.2870742112.column", "1"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.2870742112.visualization", "faceted_line_chart"),
				),
			},
			// Update dashboard title
			resource.TestStep{
				Config: testAccCheckNewRelicDashboardConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "1"),
				),
			},
			// Add widget
			resource.TestStep{
				Config: testAccCheckNewRelicDashboardWidgetConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "widget.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNewRelicDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*newrelic.Client)
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

		client := testAccProvider.Meta().(*newrelic.Client)

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

func testAccCheckNewRelicDashboardWidgetConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title                = "tf-test-updated-%s"
  widget {
    title         = "Average Transaction Duration"
	visualization = "faceted_line_chart"
	column		  = 1
	row			  = 1
    nrql          = "SELECT PERCENTILE(duration, 95) from Transaction FACET appName TIMESERIES auto"
  }
  widget {
    title         = "Page Views"
	visualization = "faceted_line_chart"
	column		  = 1
	row			  = 2
    nrql          = "SELECT AVERAGE(duration) from PageView FACET appName TIMESERIES auto"
  }
}
`, rName)
}

func testAccCheckNewRelicDashboardConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title                = "tf-test-updated-%s"
  widget {
    title         = "Average Transaction Duration"
	visualization = "faceted_line_chart"
	column		  = 1
	row			  = 1
    nrql          = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
  }
}
`, rName)
}

func testAccCheckNewRelicDashboardConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title = "tf-test-%s"
  widget {
    title         = "Average Transaction Duration"
	visualization = "faceted_line_chart"
	column		  = 1
	row			  = 1
    nrql          = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
  }
}
`, rName)
}
