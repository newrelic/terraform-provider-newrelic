package newrelic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

// Assemble the *dashboards.DashboardInput struct.
// Used by the newrelic_one_dashboard Create function.
func expandDashboardInput(d *schema.ResourceData, meta interface{}) (*dashboards.DashboardInput, error) {
	var err error

	dash := dashboards.DashboardInput{
		Name: d.Get("name").(string),
	}

	dash.Pages, err = expandDashboardPageInput(d, d.Get("page").([]interface{}), meta)
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
func expandDashboardPageInput(d *schema.ResourceData, pages []interface{}, meta interface{}) ([]dashboards.DashboardPageInput, error) {
	if len(pages) < 1 {
		return []dashboards.DashboardPageInput{}, nil
	}

	expanded := make([]dashboards.DashboardPageInput, len(pages))

	for pageIndex, v := range pages {
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
			page.GUID = common.EntityGUID(guid.(string))
		}

		// For each of the widget type, we need to expand them as well
		if widgets, ok := p["widget_area"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Area, err = expandDashboardAreaWidgetConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_bar"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Bar, err = expandDashboardBarWidgetConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_billboard"]; ok {
			for widgetIndex, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Billboard, err = expandDashboardBillboardWidgetConfigurationInput(d, v.(map[string]interface{}), meta, pageIndex, widgetIndex)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_bullet"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.RawConfiguration, err = expandDashboardBulletWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.Visualization.ID = "viz.bullet"

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_funnel"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.RawConfiguration, err = expandDashboardFunnelWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.Visualization.ID = "viz.funnel"

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_heatmap"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.RawConfiguration, err = expandDashboardHeatmapWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.Visualization.ID = "viz.heatmap"

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_histogram"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.RawConfiguration, err = expandDashboardHistogramWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.Visualization.ID = "viz.histogram"

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_line"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Line, err = expandDashboardLineWidgetConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_markdown"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Markdown, err = expandDashboardMarkdownWidgetConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_pie"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Pie, err = expandDashboardPieWidgetConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_table"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				widget.Configuration.Table, err = expandDashboardTableWidgetConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_json"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.RawConfiguration, err = expandDashboardJSONWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.Visualization.ID = "viz.json"

				page.Widgets = append(page.Widgets, widget)
			}
		}
		if widgets, ok := p["widget_stacked_bar"]; ok {
			for _, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.RawConfiguration, err = expandDashboardStackedBarWidgetRawConfigurationInput(v.(map[string]interface{}), meta)
				if err != nil {
					return nil, err
				}
				widget.Visualization.ID = "viz.stacked-bar"

				page.Widgets = append(page.Widgets, widget)
			}
		}

		expanded[pageIndex] = page
	}

	return expanded, nil
}

func expandDashboardAreaWidgetConfigurationInput(i map[string]interface{}, meta interface{}) (*dashboards.DashboardAreaWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardAreaWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}
func expandDashboardBarWidgetConfigurationInput(i map[string]interface{}, meta interface{}) (*dashboards.DashboardBarWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardBarWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}

func expandDashboardBillboardWidgetConfigurationInput(d *schema.ResourceData, i map[string]interface{}, meta interface{}, pageIndex int, widgetIndex int) (*dashboards.DashboardBillboardWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardBillboardWidgetConfigurationInput
	var err error

	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}

	// optional, order is important (API returns them sorted alpha)
	cfg.Thresholds = []dashboards.DashboardBillboardWidgetThresholdInput{}
	if _, ok := d.GetOk(fmt.Sprintf("page.%d.widget_billboard.%d.critical", pageIndex, widgetIndex)); ok {
		if t, ok := i["critical"]; ok {
			value := t.(float64)
			cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
				AlertSeverity: entities.DashboardAlertSeverityTypes.CRITICAL,
				Value:         &value,
			})
		}
	} else {
		cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
			AlertSeverity: entities.DashboardAlertSeverityTypes.CRITICAL,
			Value:         nil,
		})
	}

	if _, ok := d.GetOk(fmt.Sprintf("page.%d.widget_billboard.%d.warning", pageIndex, widgetIndex)); ok {
		if t, ok := i["warning"]; ok {
			value := t.(float64)
			cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
				AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
				Value:         &value,
			})
		}
	} else {
		cfg.Thresholds = append(cfg.Thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
			AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
			Value:         nil,
		})
	}

	return &cfg, nil
}

func expandDashboardBulletWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		Limit       float64                                    `json:"limit,omitempty"`
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput `json:"nrqlQueries"`
	}{}

	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}
	if l, ok := i["limit"]; ok {
		cfg.Limit = l.(float64)
	}

	return json.Marshal(cfg)
}

func expandDashboardFunnelWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput `json:"nrqlQueries"`
	}{}

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(cfg)
}

func expandDashboardHeatmapWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput `json:"nrqlQueries"`
	}{}

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(cfg)
}

func expandDashboardHistogramWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput `json:"nrqlQueries"`
	}{}

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(cfg)
}

func expandDashboardJSONWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput `json:"nrqlQueries"`
	}{}

	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(cfg)
}

func expandDashboardLineWidgetConfigurationInput(i map[string]interface{}, meta interface{}) (*dashboards.DashboardLineWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardLineWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}

func expandDashboardMarkdownWidgetConfigurationInput(i map[string]interface{}, meta interface{}) (*dashboards.DashboardMarkdownWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardMarkdownWidgetConfigurationInput

	if t, ok := i["text"]; ok {
		if t.(string) != "" {
			cfg.Text = t.(string)
		}

		return &cfg, nil
	}
	return nil, nil
}

func expandDashboardStackedBarWidgetRawConfigurationInput(i map[string]interface{}, meta interface{}) ([]byte, error) {
	var err error
	cfg := struct {
		NRQLQueries []dashboards.DashboardWidgetNRQLQueryInput `json:"nrqlQueries"`
	}{}

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(cfg)
}

func expandDashboardPieWidgetConfigurationInput(i map[string]interface{}, meta interface{}) (*dashboards.DashboardPieWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardPieWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
}
func expandDashboardTableWidgetConfigurationInput(i map[string]interface{}, meta interface{}) (*dashboards.DashboardTableWidgetConfigurationInput, error) {
	var cfg dashboards.DashboardTableWidgetConfigurationInput
	var err error

	// just has queries
	if q, ok := i["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, nil
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

func expandLinkedEntityGUIDs(guids []interface{}) []common.EntityGUID {
	out := make([]common.EntityGUID, len(guids))

	for i := range out {
		out[i] = common.EntityGUID(guids[i].(string))
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
			widgetType, w := flattenDashboardWidget(&widget, p.GUID)

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
func flattenDashboardWidget(in *entities.DashboardWidget, pageGUID common.EntityGUID) (string, map[string]interface{}) {
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
		for _, entity := range in.LinkedEntities {
			if entity.GetGUID() == pageGUID {
				out["filter_current_dashboard"] = true
			}
		}

		if _, ok := out["filter_current_dashboard"]; ok && out["filter_current_dashboard"] != true {
			out["linked_entity_guids"] = flattenLinkedEntityGUIDs(in.LinkedEntities)
		}

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
	case "viz.stacked-bar":
		widgetType = "widget_stacked_bar"
		if len(in.RawConfiguration) > 0 {
			cfg := struct {
				NRQLQueries []entities.DashboardWidgetNRQLQuery `json:"nrqlQueries"`
			}{}
			if err := json.Unmarshal(in.RawConfiguration, &cfg); err == nil {
				out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&cfg.NRQLQueries)
			}
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

// Function to find all of the widgets that have filter_current_dashboard set and return the title and layout location to identify later.
func findDashboardWidgetFilterCurrentDashboard(d *schema.ResourceData) ([]interface{}, error) {
	var widgetList []interface{}

	pages := d.Get("page").([]interface{})
	selfLinkingWidgets := []string{"widget_bar", "widget_pie", "widget_table"}

	for i, v := range pages {
		p := v.(map[string]interface{})
		// For each of the widget type, we need to expand them as well
		for _, widgetType := range selfLinkingWidgets {
			if widgets, ok := p[widgetType]; ok {
				for _, widget := range widgets.([]interface{}) {
					w := widget.(map[string]interface{})
					if v, ok := w["filter_current_dashboard"]; ok && v.(bool) {

						if l, ok := w["linked_entity_guids"]; ok && len(l.([]interface{})) > 0 {
							return nil, fmt.Errorf("err: filter_current_dashboard can't be set if linked_entity_guids is configured")
						}

						unqWidget := make(map[string]interface{})
						if t, ok := w["title"]; ok {
							unqWidget["title"] = t.(string)
						}
						if r, ok := w["row"]; ok {
							unqWidget["row"] = r.(int)
						}
						if c, ok := w["column"]; ok {
							unqWidget["column"] = c.(int)
						}

						unqWidget["page"] = i

						widgetList = append(widgetList, unqWidget)
					}
				}
			}
		}

	}

	return widgetList, nil

}

// Function to set the page guid as the linked entity now that the page is created
func setDashboardWidgetFilterCurrentDashboardLinkedEntity(d *schema.ResourceData, filterWidgets []interface{}) error {

	if len(filterWidgets) < 1 {
		log.Printf("[INFO] Empty list of widgets to filter")
		return nil
	}

	selfLinkingWidgets := []string{"widget_bar", "widget_pie", "widget_table"}

	pages := d.Get("page").([]interface{})
	for i, v := range pages {
		p := v.(map[string]interface{})
		for _, widgetType := range selfLinkingWidgets {
			if widgets, ok := p[widgetType]; ok {
				for _, k := range widgets.([]interface{}) {
					w := k.(map[string]interface{})
					for _, f := range filterWidgets {
						e := f.(map[string]interface{})
						if e["page"] == i {
							if w["title"] == e["title"] && w["column"] == e["column"] && w["row"] == e["row"] {
								w["linked_entity_guids"] = []string{p["guid"].(string)}
							}
						}
					}
				}
			}
		}
	}

	if err := d.Set("page", pages); err != nil {
		return err
	}

	return nil
}
