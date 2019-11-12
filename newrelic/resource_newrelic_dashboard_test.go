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
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),

				// Check exists
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "title", rName),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "editable", "editable_by_all"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "icon", "bar-chart"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "visibility", "all"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.#", "4"),

					// filters
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "filter.#", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "filter.0.event_types.#", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "filter.0.event_types.4104882694", "Transaction"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "filter.0.attributes.#", "2"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "filter.0.attributes.2634578693", "appName"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "filter.0.attributes.3755723101", "envName"),

					// billboard widget
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.title", "Transaction Count"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.visualization", "billboard"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.threshold_red", "100"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.threshold_yellow", "50"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.height", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.width", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.row", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2779494853.column", "1"),

					// faceted_line_chart widget
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.title", "Average Transaction Duration"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.visualization", "facet_bar_chart"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.nrql", "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.drilldown_dashboard_id", "1234"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.height", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.width", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.row", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1942929817.column", "2"),

					// markdown widget
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.title", "Dashboard Note"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.visualization", "markdown"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.source", "#h1 Heading"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.height", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.width", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.row", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.4111231263.column", "3"),

					// metric_line_chart widget
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.title", "Apdex"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.visualization", "metric_line_chart"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.height", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.width", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.row", "2"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.column", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.#", "2"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.1986073529.offset_duration", "P7D"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.1986073529.presentation.#", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.1986073529.presentation.0.color", "#b1b6ba"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.1986073529.presentation.0.name", "Last week"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.4162716032.offset_duration", "P1D"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.4162716032.presentation.#", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.4162716032.presentation.0.color", "#77add4"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.compare_with.4162716032.presentation.0.name", "Yesterday"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.metric.1284261170.name", "Apdex"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.metric.1284261170.values.#", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.metric.1284261170.values.2136473340", "score"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.raw_metric_name", "Apdex"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.entity_ids.#", "1"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.1516746896.entity_ids.743249691", "1234"),
				),
			},

			// Import
			{
				ResourceName:      "newrelic_dashboard.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicDashboard_NoDiffOnReapply(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),
			},
			{
				Config:             testAccCheckNewRelicDashboardConfig(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccNewRelicDashboard_UpdateDashboard(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	rNameUpdated := fmt.Sprintf("%s-updated", rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "title", rName),
				),
			},
			{
				Config: testAccCheckNewRelicDashboardConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "title", rNameUpdated),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.#", "4"),
				),
			},
		},
	})
}

func TestAccNewRelicDashboard_AddWidget(t *testing.T) {
	rDashboardName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	widgetName := "Page Views"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardConfig(rDashboardName),
			},
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigAdded(rDashboardName, widgetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "title", rDashboardName),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.#", "5"),
				),
			},
		},
	})
}

func TestAccNewRelicDashboard_UpdateWidget(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	widgetName := "Page Views"
	widgetNameUpdated := "Page Views Updated"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigAdded(rName, widgetName),
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "title", rName),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.#", "5"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.3349214025.title", widgetName),
				),
			},
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigAdded(rName, widgetNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "title", rName),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.#", "5"),
					resource.TestCheckResourceAttr("newrelic_dashboard.foo", "widget.2318477184.title", widgetNameUpdated),
				),
			},
		},
	})
}

func testAccCheckNewRelicDashboardConfig(dashboardName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title = "%s"

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
    title = "Transaction Count"
    visualization = "billboard"
    nrql = "SELECT count(*) from Transaction since 5 minutes ago facet appName"
    threshold_red = 100
    threshold_yellow = 50
    row    = 1
    column = 1
  }
  widget {
    title         		   = "Average Transaction Duration"
    visualization 		   = "facet_bar_chart"
	nrql          		   = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
	drilldown_dashboard_id = 1234
    row           		   = 1
	column        		   = 2
  }
  widget {
    title         = "Dashboard Note"
    visualization = "markdown"
    source        = "#h1 Heading"
    row           = 1
    column        = 3
  }
  widget {
    title         = "Apdex"
    visualization = "metric_line_chart"
    duration = 1800000
    entity_ids = [ 1234 ]
    compare_with {
        offset_duration = "P1D"
        presentation {
            color = "#77add4"
            name = "Yesterday"
        }
    }

    compare_with {
        offset_duration = "P7D"
        presentation {
            color = "#b1b6ba"
            name = "Last week"
        }
    }
    metric {
        name = "Apdex"
        values = [ "score" ]
    }
    row           = 2
    column        = 1
  }
}
`, dashboardName)
}

func testAccCheckNewRelicDashboardWidgetConfigAdded(dashboardName string, widgetName string) string {
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
    title = "Transaction Count"
    visualization = "billboard"
    nrql = "SELECT count(*) from Transaction since 5 minutes ago facet appName"
    threshold_red = 100
    threshold_yellow = 50
    row    = 1
    column = 1
  }
  widget {
    title         		   = "Average Transaction Duration"
    visualization 		   = "facet_bar_chart"
	nrql          		   = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
	drilldown_dashboard_id = 1234
    row           		   = 1
	column        		   = 2
  }
  widget {
    title         = "Dashboard Note"
    visualization = "markdown"
    source        = "#h1 Heading"
    row           = 1
    column        = 3
  }
  widget {
    title         = "Apdex"
    visualization = "metric_line_chart"
    duration = 1800000
    entity_ids = [ 1234 ]
    compare_with {
        offset_duration = "P1D"
        presentation {
            color = "#77add4"
            name = "Yesterday"
        }
    }

    compare_with {
        offset_duration = "P7D"
        presentation {
            color = "#b1b6ba"
            name = "Last week"
        }
    }
    metric {
        name = "Apdex"
        values = [ "score" ]
    }
    row           = 2
    column        = 1
  }
  widget {
    title         = "%s"
	visualization = "faceted_line_chart"
	column        = 2
	row           = 2
    nrql          = "SELECT AVERAGE(duration) from PageView FACET appName TIMESERIES auto"
  }
}
`, dashboardName, widgetName)
}

func testAccCheckNewRelicDashboardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
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
			return fmt.Errorf("dashboard not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
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
			return fmt.Errorf("dashboard still exists")
		}

	}
	return nil
}

// A custom check function to log the state during a test run.
// This is useful to find the individual widget hash values when writing assertions against them.
func logState(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		t.Logf("State: %s\n", s)

		return nil
	}
}
