package newrelic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

// knownWidgetTypes is a list of widget blocks in the form of `widget_<type>` that
// all get handled the same.
var knownWidgetTypes = []string{
	"area",
	"bar",
	"billboard",
	"bullet",
	"funnel",
	"heatmap",
	"histogram",
	"json",
	"line",
	"markdown",
	"pie",
	"table",
}

// Assemble the *dashboards.DashboardInput struct.
// Used by the newrelic_one_dashboard Create function.
func expandDashboardInput(d *schema.ResourceData, meta interface{}) (*dashboards.DashboardInput, error) {
	var err error

	dash := dashboards.DashboardInput{
		Name: d.Get("name").(string),
	}

	dash.Pages, err = expandDashboardPageInput(d.Get("page").([]interface{}), meta)
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

// expandDashboardPageInput converts the HCL for a Page into the complex GraphQL structure required
// to create a page.
func expandDashboardPageInput(pages []interface{}, meta interface{}) ([]dashboards.DashboardPageInput, error) {
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
		if widgets, ok := p["widget"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				properties := v.(map[string]interface{})
				widget, err := expandDashboardWidgetInput(properties, meta)
				if err != nil {
					return nil, err
				}

				// Get and set raw widget properties
				if q, ok := properties["configuration"]; ok {
					widget.RawConfiguration = entities.DashboardWidgetRawConfiguration(q.(string))
				}

				// Get and set widget visualization_id
				if q, ok := properties["visualization_id"]; ok {
					widget.Visualization.ID = q.(string)
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}

		// For `widget_<type>` blocks, loop through and handle them all the same
		for _, wType := range knownWidgetTypes {
			if widgets, ok := p["widget_"+wType]; ok {
				for _, v := range widgets.([]interface{}) {
					// Get generic properties set
					widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
					if err != nil {
						return nil, err
					}

					widget.Visualization.ID = "viz." + wType
					widget.RawConfiguration, err = expandDashboardWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
					if err != nil {
						return nil, err
					}

					page.Widgets = append(page.Widgets, widget)
				}
			}
		}

		expanded[i] = page
	}

	return expanded, nil
}

//
// expandDashboardWidgetRawConfigurationInput is a generic HCL Block to Raw configuration translator
//
func expandDashboardWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		Limit       float64                                             `json:"limit,omitempty"`
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput          `json:"nrqlQueries,omitempty"`
		Thresholds  []dashboards.DashboardBillboardWidgetThresholdInput `json:"thresholds,omitempty"`
		Text        string                                              `json:"text,omitempty"`
	}{}

	// All except Markdown
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}
	// bar, bullet,
	if l, ok := i["limit"]; ok {
		cfg.Limit = l.(float64)
	}

	// markdown
	if l, ok := i["text"]; ok {
		if l.(string) != "" {
			cfg.Text = l.(string)
		}
	}

	// billboard
	if t, ok := i["critical"]; ok {
		cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
			AlertSeverity: entities.DashboardAlertSeverityTypes.CRITICAL,
			Value:         t.(float64),
		})
	}

	// billboard
	if t, ok := i["warning"]; ok {
		cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
			AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
			Value:         t.(float64),
		})
	}

	return json.Marshal(cfg)
}

// expandDashboardWidgetInput expands the common items in WidgetInput, but not the configuration
// which is specific to the widgets
func expandDashboardWidgetInput(w map[string]interface{}, meta interface{}) (dashboards.DashboardWidgetInput, error) {
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

	if i, ok := w["linked_entity_guids"]; ok {
		widget.LinkedEntityGUIDs = expandLinkedEntityGUIDs(i.([]interface{}))
	}

	return widget, nil
}

func expandLinkedEntityGUIDs(guids []interface{}) []entities.EntityGUID {
	out := make([]entities.EntityGUID, len(guids))

	for i := range out {
		out[i] = entities.EntityGUID(guids[i].(string))
	}

	return out
}

func expandDashboardWidgetNRQLQueryInput(queries []interface{}, meta interface{}) ([]dashboards.DashboardWidgetNRQLQueryInput, error) {
	if len(queries) < 1 {
		return []dashboards.DashboardWidgetNRQLQueryInput{}, nil
	}

	expanded := make([]dashboards.DashboardWidgetNRQLQueryInput, len(queries))

	for i, v := range queries {
		var query dashboards.DashboardWidgetNRQLQueryInput
		q := v.(map[string]interface{})

		if acct, ok := q["account_id"]; ok {
			query.AccountID = acct.(int)
		}

		if query.AccountID < 1 {
			defs := meta.(map[string]interface{})
			if acct, ok := defs["account_id"]; ok {
				query.AccountID = acct.(int)
			}
		}

		if nrql, ok := q["query"]; ok {
			query.Query = nrdb.NRQL(nrql.(string))
		}

		expanded[i] = query
	}

	return expanded, nil
}

// Unpack the *dashboards.Dashboard variable and set resource data.
//
// Used by the newrelic_one_dashboard Read function (resourceNewRelicOneDashboardRead)
func flattenDashboardEntity(dashboard *entities.DashboardEntity, d *schema.ResourceData) error {
	_ = d.Set("account_id", dashboard.AccountID)
	_ = d.Set("guid", dashboard.GUID)
	_ = d.Set("name", dashboard.Name)
	_ = d.Set("permalink", dashboard.Permalink)
	_ = d.Set("permissions", strings.ToLower(string(dashboard.Permissions)))

	if dashboard.Description != "" {
		_ = d.Set("description", dashboard.Description)
	}

	if dashboard.Pages != nil && len(dashboard.Pages) > 0 {
		pages := flattenDashboardPage(&dashboard.Pages)
		if err := d.Set("page", pages); err != nil {
			return err
		}
	}

	return nil
}

// Unpack the *dashboards.Dashboard variable and set resource data.
//
// Used by the newrelic_one_dashboard Read function (resourceNewRelicOneDashboardRead)
func flattenDashboardUpdateResult(result *dashboards.DashboardUpdateResult, d *schema.ResourceData) error {
	if result == nil {
		return fmt.Errorf("can not flatten nil DashboardUpdateResult")
	}

	dashboard := result.EntityResult // dashboard.DashboardEntityResult

	_ = d.Set("account_id", dashboard.AccountID)
	_ = d.Set("guid", dashboard.GUID)
	_ = d.Set("name", dashboard.Name)
	//d.Set("permalink", dashboard.Permalink)
	_ = d.Set("permissions", strings.ToLower(string(dashboard.Permissions)))

	if dashboard.Description != "" {
		_ = d.Set("description", dashboard.Description)
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

		for _, widget := range p.Widgets {
			widgetType, w := flattenDashboardWidget(&widget)

			if widgetType != "" {
				if _, ok := m[widgetType]; !ok {
					m[widgetType] = []interface{}{}
				}

				m[widgetType] = append(m[widgetType].([]interface{}), w)
			}
		}

		out[i] = m
	}

	return out
}

func flattenLinkedEntityGUIDs(linkedEntities []entities.EntityOutlineInterface) []string {
	out := make([]string, len(linkedEntities))

	for i, entity := range linkedEntities {
		out[i] = string(entity.GetGUID())
	}

	return out
}

// nolint:gocyclo
func flattenDashboardWidget(in *entities.DashboardWidget) (string, map[string]interface{}) {
	var widgetType string
	out := make(map[string]interface{})

	out["id"] = in.ID
	out["column"] = in.Layout.Column
	out["height"] = in.Layout.Height
	out["row"] = in.Layout.Row
	out["width"] = in.Layout.Width
	if in.Title != "" {
		out["title"] = in.Title
	}

	// NOTE: The widget types that currently support linked entities
	// are faceted widgets - i.e. bar, line, pie
	if len(in.LinkedEntities) > 0 {
		out["linked_entity_guids"] = flattenLinkedEntityGUIDs(in.LinkedEntities)
	}

	switch in.Visualization.ID {
	case "viz.area":
		widgetType = "widget_area"
		if len(in.Configuration.Area.NRQLQueries) > 0 {
			out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&in.Configuration.Area.NRQLQueries)
		}
	case "viz.bar":
		widgetType = "widget_bar"
		if len(in.Configuration.Bar.NRQLQueries) > 0 {
			out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&in.Configuration.Bar.NRQLQueries)
		}
	case "viz.billboard":
		widgetType = "widget_billboard"
		if len(in.Configuration.Billboard.NRQLQueries) > 0 {
			out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&in.Configuration.Billboard.NRQLQueries)
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
	case "viz.bullet":
		widgetType = "widget_bullet"
		if len(in.RawConfiguration) > 0 {
			cfg := struct {
				Limit       float64                             `json:"limit"`
				NRQLQueries []entities.DashboardWidgetNRQLQuery `json:"nrqlQueries"`
			}{}
			if err := json.Unmarshal(in.RawConfiguration, &cfg); err == nil {
				out["limit"] = cfg.Limit
				out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&cfg.NRQLQueries)
			}
		}
	case "viz.funnel":
		widgetType = "widget_funnel"
		if len(in.RawConfiguration) > 0 {
			cfg := struct {
				NRQLQueries []entities.DashboardWidgetNRQLQuery `json:"nrqlQueries"`
			}{}
			if err := json.Unmarshal(in.RawConfiguration, &cfg); err == nil {
				out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&cfg.NRQLQueries)
			}
		}
	case "viz.heatmap":
		widgetType = "widget_heatmap"
		if len(in.RawConfiguration) > 0 {
			cfg := struct {
				NRQLQueries []entities.DashboardWidgetNRQLQuery `json:"nrqlQueries"`
			}{}
			if err := json.Unmarshal(in.RawConfiguration, &cfg); err == nil {
				out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&cfg.NRQLQueries)
			}
		}
	case "viz.histogram":
		widgetType = "widget_histogram"
		if len(in.RawConfiguration) > 0 {
			cfg := struct {
				NRQLQueries []entities.DashboardWidgetNRQLQuery `json:"nrqlQueries"`
			}{}
			if err := json.Unmarshal(in.RawConfiguration, &cfg); err == nil {
				out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&cfg.NRQLQueries)
			}
		}
	case "viz.json":
		widgetType = "widget_json"
		if len(in.RawConfiguration) > 0 {
			cfg := struct {
				NRQLQueries []entities.DashboardWidgetNRQLQuery `json:"nrqlQueries"`
			}{}
			if err := json.Unmarshal(in.RawConfiguration, &cfg); err == nil {
				out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&cfg.NRQLQueries)
			}
		}
	case "viz.line":
		widgetType = "widget_line"
		if len(in.Configuration.Line.NRQLQueries) > 0 {
			out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&in.Configuration.Line.NRQLQueries)
		}
	case "viz.markdown":
		widgetType = "widget_markdown"
		if in.Configuration.Markdown.Text != "" {
			out["text"] = in.Configuration.Markdown.Text
		}
	case "viz.pie":
		widgetType = "widget_pie"
		if len(in.Configuration.Pie.NRQLQueries) > 0 {
			out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&in.Configuration.Pie.NRQLQueries)
		}
	case "viz.table":
		widgetType = "widget_table"
		if len(in.Configuration.Table.NRQLQueries) > 0 {
			out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&in.Configuration.Table.NRQLQueries)
		}
	// This probably means we have a dynamic widget
	default:
		widgetType = "widget"
		out["visualization_id"] = in.Visualization.ID
		if len(in.RawConfiguration) > 0 {
			out["configuration"] = string(in.RawConfiguration)
		}
	}

	return widgetType, out
}

func flattenDashboardWidgetNRQLQuery(in *[]entities.DashboardWidgetNRQLQuery) []interface{} {
	out := make([]interface{}, len(*in))

	for i, v := range *in {
		m := make(map[string]interface{})

		m["account_id"] = v.AccountID
		m["query"] = v.Query

		out[i] = m
	}

	return out
}
