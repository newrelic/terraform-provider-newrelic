// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func TestAccNewRelicOneDashboard_Basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicOneDashboardConfig(rName, strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicOneDashboardExists("newrelic_one_dashboard.bar"),
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

/*
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
				),
			},
			{
				Config: testAccCheckNewRelicDashboardConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
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
				),
			},
			{
				Config: testAccCheckNewRelicDashboardWidgetConfigAdded(rName, widgetNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					// logState(t),
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
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
*/

func testAccCheckNewRelicOneDashboardConfig(dashboardName string, accountID string) string {
	return `
resource "newrelic_one_dashboard" "bar" {
  name = "` + dashboardName + `"
  permissions = "private"

  page {
    name = "` + dashboardName + `"

    widget_area {
      title = "Area 51"
      row = 1
      column = 1
      height = 3
      width = 12

      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT 51 TIMESERIES"
      }
    }

    widget_bar {
      title = "foo"
      row = 4
      column = 1
      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT count(*) FACET name"
      }
    }

    widget_billboard {
      title = "top 40"
      row = 4
      column = 5
      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT count(*)"
      }

      warning = 1
      critical = 2
    }

    widget_line {
      title = "over the"
      row = 4
      column = 9
      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT 1 TIMESERIES"
      }
      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT 2 TIMESERIES"
      }
    }

    widget_markdown {
      title = "My cool widget"
      row = 7
      column = 1
      text = "# Header text"
    }

    widget_pie {
      title = "3.14"
      row = 7
      column = 5
      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT count(*) FACET name"
      }
    }

    widget_table {
      title = "Round"
      row = 7
      column = 9
      query {
        account_id = ` + accountID + `
        nrql = "FROM Transaction SELECT *"
      }
    }
  }
}
`
}

/*
func testAccCheckNewRelicDashboardWidgetConfigAdded(dashboardName string, widgetName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title                = "%s"
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
*/
func testAccCheckNewRelicOneDashboardExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

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

/*
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
*/
