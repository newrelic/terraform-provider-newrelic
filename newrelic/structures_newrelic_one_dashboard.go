package newrelic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

// Assemble the *dashboards.DashboardInput struct.
// Used by the newrelic_one_dashboard Create function.
func expandDashboardInput(d *schema.ResourceData) (*dashboards.DashboardInput, error) {
	var err error

	dash := dashboards.DashboardInput{
		Name: d.Get("name").(string),
	}

	dash.Pages, err = expandDashboardPageInput(d.Get("page").([]interface{}))
	if err != nil {
		return nil, err
	}

	// Optional, with default
	perm := d.Get("permissions").(string)
	dash.Permissions = entities.DashboardPermissions(strings.ToUpper(perm))

	// Optional
	if e, ok := d.GetOk("description"); ok {
		dash.Description = e.(string)
	}

	return &dash, nil
}

// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func expandDashboardPageInput(pages []interface{}) ([]dashboards.DashboardPageInput, error) {
	if len(pages) < 1 {
		return []dashboards.DashboardPageInput{}, nil
	}

	expanded := make([]dashboards.DashboardPageInput, len(pages))

	for i, v := range pages {
		var page dashboards.DashboardPageInput
		p := v.(map[string]interface{})

		if name, ok := p["name"]; ok {
			page.Name = name.(string)
		} else {
			return nil, fmt.Errorf("name required for dashboard page")
		}

		if desc, ok := p["description"]; ok {
			page.Description = desc.(string)
		}

		// GUID exists for Update, null for new page
		if guid, ok := p["guid"]; ok {
			page.GUID = entities.EntityGUID(guid.(string))
		}

		// For each of the widget type, we need to expand them as well
		if widgets, ok := p["widget_area"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Area, err = expandDashboardAreaWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_bar"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Bar, err = expandDashboardBarWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_billboard"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Billboard, err = expandDashboardBillboardWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_line"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Line, err = expandDashboardLineWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_markdown"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Markdown, err = expandDashboardMarkdownWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_pie"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Pie, err = expandDashboardPieWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_table"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				widget.Configuration.Table, err = expandDashboardTableWidgetConfigurationInput(v.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}

		expanded[i] = page
	}

	return expanded, nil
}

func expandDashboardAreaWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardAreaWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardAreaWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["query"]; ok {
		cfg.Queries, err = expandDashboardWidgetQueryInput(q.([]interface{}))
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}
func expandDashboardBarWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardBarWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardBarWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["query"]; ok {
		cfg.Queries, err = expandDashboardWidgetQueryInput(q.([]interface{}))
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}

func expandDashboardBillboardWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardBillboardWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardBillboardWidgetConfigurationInput
	var err error

	if q, ok := i["query"]; ok {
		cfg.Queries, err = expandDashboardWidgetQueryInput(q.([]interface{}))
		if err != nil {
			return nil, err
		}
	}

	// optional, order is important (API returns them sorted alpha)
	cfg.Thresholds = []dashboards.DashboardBillboardWidgetThresholdInput{}
	if t, ok := i["critical"]; ok {
		cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
			AlertSeverity: entities.DashboardAlertSeverityTypes.CRITICAL,
			Value:         t.(float64),
		})
	}

	if t, ok := i["warning"]; ok {
		cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
			AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
			Value:         t.(float64),
		})
	}

	return &cfg, nil
}

func expandDashboardLineWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardLineWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardLineWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["query"]; ok {
		cfg.Queries, err = expandDashboardWidgetQueryInput(q.([]interface{}))
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}
func expandDashboardMarkdownWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardMarkdownWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardMarkdownWidgetConfigurationInput

	if t, ok := i["text"]; ok {
		if t.(string) != "" {
			cfg.Text = t.(string)
		}

		return &cfg, nil
	}
	return nil, nil
}
func expandDashboardPieWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardPieWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardPieWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["query"]; ok {
		cfg.Queries, err = expandDashboardWidgetQueryInput(q.([]interface{}))
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}
func expandDashboardTableWidgetConfigurationInput(i map[string]interface{}) (*dashboards.DashboardTableWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardTableWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["query"]; ok {
		cfg.Queries, err = expandDashboardWidgetQueryInput(q.([]interface{}))
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}

// expandDashboardWidgetInput expands the common items in WidgetInput, but not the configuration
// which is specific to the widgets
func expandDashboardWidgetInput(w map[string]interface{}) (dashboards.DashboardWidgetInput, error) {
	var widget dashboards.DashboardWidgetInput

	if i, ok := w["id"]; ok {
		widget.ID = i.(string)
	}
	if i, ok := w["column"]; ok {
		widget.Layout.Column = i.(int)
	}
	if i, ok := w["height"]; ok {
		widget.Layout.Height = i.(int)
	}
	if i, ok := w["row"]; ok {
		widget.Layout.Row = i.(int)
	}
	if i, ok := w["width"]; ok {
		widget.Layout.Width = i.(int)
	}
	if i, ok := w["title"]; ok {
		widget.Title = i.(string)
	}

	return widget, nil
}

func expandDashboardWidgetQueryInput(queries []interface{}) ([]dashboards.DashboardWidgetQueryInput, error) {
	if len(queries) < 1 {
		return []dashboards.DashboardWidgetQueryInput{}, nil
	}

	expanded := make([]dashboards.DashboardWidgetQueryInput, len(queries))

	for i, v := range queries {
		var query dashboards.DashboardWidgetQueryInput
		q := v.(map[string]interface{})

		if acct, ok := q["account_id"]; ok {
			query.AccountID = acct.(int)
		}

		if nrql, ok := q["nrql"]; ok {
			query.NRQL = nrdb.NRQL(nrql.(string))
		}

		expanded[i] = query
	}

	return expanded, nil
}

// Unpack the *dashboards.Dashboard variable and set resource data.
//
// Used by the newrelic_dashboard Read function (resourceNewRelicDashboardRead)
func flattenOneDashboard(dashboard *entities.DashboardEntity, d *schema.ResourceData) error {
	d.Set("account_id", dashboard.AccountID)
	d.Set("guid", dashboard.GUID)
	d.Set("name", dashboard.Name)
	d.Set("permalink", dashboard.Permalink)
	d.Set("permissions", strings.ToLower(string(dashboard.Permissions)))

	if dashboard.Description != "" {
		d.Set("description", dashboard.Description)
	}

	if dashboard.Pages != nil && len(dashboard.Pages) > 0 {
		pages := flattenDashboardPage(&dashboard.Pages)
		if err := d.Set("page", pages); err != nil {
			return err
		}
	}

	return nil
}

// return []interface{} because Page is a SetList
func flattenDashboardPage(in *[]entities.DashboardPage) []interface{} {
	out := make([]interface{}, len(*in))

	for i, p := range *in {
		m := make(map[string]interface{})

		m["guid"] = p.GUID
		m["name"] = p.Name

		if p.Description != "" {
			m["description"] = p.Description
		}

		m["widget_area"] = []interface{}{}
		m["widget_bar"] = []interface{}{}
		m["widget_billboard"] = []interface{}{}
		m["widget_line"] = []interface{}{}
		m["widget_markdown"] = []interface{}{}
		m["widget_pie"] = []interface{}{}
		m["widget_table"] = []interface{}{}

		for _, widget := range p.Widgets {
			var widgetType string
			w := flattenDashboardWidget(&widget)

			switch widget.Visualization.ID {
			case "viz.area":
				widgetType = "widget_area"
			case "viz.bar":
				widgetType = "widget_bar"
			case "viz.billboard":
				widgetType = "widget_billboard"
			case "viz.line":
				widgetType = "widget_line"
			case "viz.markdown":
				widgetType = "widget_markdown"
			case "viz.pie":
				widgetType = "widget_pie"
			case "viz.table":
				widgetType = "widget_table"
			}

			if widgetType != "" {
				m[widgetType] = append(m[widgetType].([]interface{}), w)
			}
		}

		out[i] = m
	}

	log.Printf("flattenDashboardPage: '%+v'", out)
	return out
}

func flattenDashboardWidget(in *entities.DashboardWidget) map[string]interface{} {
	out := make(map[string]interface{})

	out["id"] = in.ID
	out["column"] = in.Layout.Column
	out["height"] = in.Layout.Height
	out["row"] = in.Layout.Row
	out["width"] = in.Layout.Width
	if in.Title != "" {
		out["title"] = in.Title
	}

	switch in.Visualization.ID {
	case "viz.area":
		if len(in.Configuration.Area.Queries) > 0 {
			out["query"] = flattenDashboardWidgetQuery(&in.Configuration.Area.Queries)
		}
	case "viz.bar":
		if len(in.Configuration.Bar.Queries) > 0 {
			out["query"] = flattenDashboardWidgetQuery(&in.Configuration.Bar.Queries)
		}
	case "viz.billboard":
		if len(in.Configuration.Billboard.Queries) > 0 {
			out["query"] = flattenDashboardWidgetQuery(&in.Configuration.Billboard.Queries)
		}
		if len(in.Configuration.Billboard.Thresholds) > 0 {
			for _, v := range in.Configuration.Billboard.Thresholds {
				switch v.AlertSeverity {
				case entities.DashboardAlertSeverityTypes.CRITICAL:
					out["critical"] = v.Value
				case entities.DashboardAlertSeverityTypes.WARNING:
					out["warning"] = v.Value
				}
			}
		}
	case "viz.line":
		if len(in.Configuration.Line.Queries) > 0 {
			out["query"] = flattenDashboardWidgetQuery(&in.Configuration.Line.Queries)
		}
	case "viz.markdown":
		if in.Configuration.Markdown.Text != "" {
			out["text"] = in.Configuration.Markdown.Text
		}
	case "viz.pie":
		if len(in.Configuration.Pie.Queries) > 0 {
			out["query"] = flattenDashboardWidgetQuery(&in.Configuration.Pie.Queries)
		}
	case "viz.table":
		if len(in.Configuration.Table.Queries) > 0 {
			out["query"] = flattenDashboardWidgetQuery(&in.Configuration.Table.Queries)
		}
	}

	return out
}

func flattenDashboardWidgetQuery(in *[]entities.DashboardWidgetQuery) []interface{} {
	out := make([]interface{}, len(*in))

	for i, v := range *in {
		m := make(map[string]interface{})

		m["account_id"] = v.AccountID
		m["nrql"] = v.NRQL

		out[i] = m
	}

	return out
}
