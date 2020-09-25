// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
)

func TestAccNewRelicDashboard_CrossAccountWidget(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigCrossAccountWidget(rName),
				// Check:  resource.ComposeTestCheckFunc(
				// 	logState(t),
				// ),
			},
		},
	})
}

func testAccCheckNewRelicDashboardWidgetConfigCrossAccountWidget(dashboardName string) string {
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
		account_id = "2508259"
    title = "Transaction Count"
    visualization = "billboard"
    nrql = "SELECT count(*) from Transaction since 5 minutes ago FACET appName"
    threshold_red = 100
    threshold_yellow = 50
    row    = 1
    column = 1
  }

	grid_column_count = 12
}
`, dashboardName)
}

func TestAccNewRelicDashboard_Basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
				),
			},
			// Import
			{
				ResourceName:      "newrelic_dashboard.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// grid_column_count is not returned in the GET response
				ImportStateVerifyIgnore: []string{"grid_column_count"},
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

func TestNewRelicDashboard_WidgetValidation(t *testing.T) {
	cases := []struct {
		cfg            map[string]interface{}
		visualizations []string
		condition      string
	}{
		{
			condition: "nrql field missing",
			visualizations: []string{
				"attribute_sheet",
				"billboard",
				"billboard_comparison",
				"comparison_line_chart",
				"event_feed",
				"event_table",
				"facet_bar_chart",
				"facet_pie_chart",
				"facet_table",
				"faceted_area_chart",
				"faceted_line_chart",
				"funnel",
				"gauge",
				"heatmap",
				"histogram",
				"line_chart",
				"raw_json",
				"single_event",
				"uniques_list",
			},
			cfg: map[string]interface{}{
				"title":         "title",
				"widget_id":     1234,
				"threshold_red": 1,
				"row":           1,
				"column":        1,
				"width":         1,
				"height":        1,
			},
		},
		{
			condition: "threshold_red field missing",
			visualizations: []string{
				"gauge",
			},
			cfg: map[string]interface{}{
				"title":     "title",
				"nrql":      "nrql",
				"widget_id": 1234,
				"row":       1,
				"column":    1,
				"width":     1,
				"height":    1,
			},
		},
		{
			condition: "source field missing",
			visualizations: []string{
				"markdown",
			},
			cfg: map[string]interface{}{
				"title":     "title",
				"widget_id": 1234,
				"row":       1,
				"column":    1,
				"width":     1,
				"height":    1,
			},
		},
		{
			condition: "metric field missing",
			visualizations: []string{
				"metric_line_chart",
			},
			cfg: map[string]interface{}{
				"title":      "title",
				"widget_id":  1234,
				"entity_ids": schema.NewSet(schema.HashInt, []interface{}{1234}),
				"duration":   1800000,
				"row":        1,
				"column":     1,
				"width":      1,
				"height":     1,
			},
		},
		{
			condition: "entity_ids field missing",
			visualizations: []string{
				"application_breakdown",
				"metric_line_chart",
			},
			cfg: map[string]interface{}{
				"title":     "title",
				"widget_id": 1234,
				"metric":    schema.NewSet(schema.HashString, []interface{}{}),
				"duration":  1800000,
				"row":       1,
				"column":    1,
				"width":     1,
				"height":    1,
			},
		},
		{
			condition: "duration field missing",
			visualizations: []string{
				"metric_line_chart",
			},
			cfg: map[string]interface{}{
				"title":      "title",
				"widget_id":  1234,
				"entity_ids": schema.NewSet(schema.HashInt, []interface{}{1234}),
				"metric":     schema.NewSet(schema.HashString, []interface{}{}),
				"row":        1,
				"column":     1,
				"width":      1,
				"height":     1,
			},
		},
	}

	for _, c := range cases {
		for _, v := range c.visualizations {
			c.cfg["visualization"] = v

			_, err := expandWidget(c.cfg)

			if err == nil {
				t.Errorf("validation error expected when %s for %s visualization", c.condition, v)
			}
		}
	}
}

func TestAccNewRelicDashboard_MissingDashboard(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicDashboardConfig(rName),
			},
			{
				PreConfig: deleteDashboard(rName),
				Config:    testAccCheckNewRelicDashboardConfig(rName),
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

	grid_column_count = 12
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

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.Dashboards.GetDashboard(int(id))
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
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_dashboard" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = client.Dashboards.GetDashboard(int(id))

		if err == nil {
			return fmt.Errorf("dashboard still exists")
		}

	}
	return nil
}

func deleteDashboard(title string) func() {
	return func() {
		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		params := dashboards.ListDashboardsParams{
			Title: title,
		}

		dashboards, _ := client.Dashboards.ListDashboards(&params)

		for _, d := range dashboards {
			if d.Title == title {
				_, _ = client.Dashboards.DeleteDashboard(d.ID)
				break
			}
		}
	}
}
