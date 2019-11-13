package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

// Assemble the *newrelic.Dashboard variable.
//
// Used by the newrelic_dashboard Create and Update functions.
func expandDashboard(d *schema.ResourceData) (*newrelic.Dashboard, error) {
	metadata := newrelic.DashboardMetadata{
		Version: 1,
	}

	// TODO: Some of these should be terraform defaults and validated
	dashboard := newrelic.Dashboard{
		Title:      d.Get("title").(string),
		Metadata:   metadata,
		Icon:       d.Get("icon").(string),
		Visibility: d.Get("visibility").(string),
		Editable:   d.Get("editable").(string),
	}

	if f, ok := d.GetOk("filter"); ok {
		filter := f.([]interface{})[0].(map[string]interface{})
		dashboardFilter := newrelic.DashboardFilter{}

		if v, ok := filter["attributes"]; ok {
			attributes := v.(*schema.Set).List()
			vs := make([]string, 0, len(attributes))
			for _, a := range attributes {
				vs = append(vs, a.(string))
			}

			dashboardFilter.Attributes = vs
		}

		if v, ok := filter["event_types"]; ok {
			eventTypes := v.(*schema.Set).List()
			vs := make([]string, 0, len(eventTypes))
			for _, e := range eventTypes {
				vs = append(vs, e.(string))
			}
			dashboardFilter.EventTypes = vs
		}
		dashboard.Filter = dashboardFilter
	}

	if widgets, ok := d.GetOk("widget"); ok && widgets.(*schema.Set).Len() > 0 {
		expandedWidgets, err := expandWidgets(widgets.(*schema.Set).List())

		if err != nil {
			return nil, err
		}

		dashboard.Widgets = expandedWidgets
	}

	return &dashboard, nil
}

func expandWidgets(widgets []interface{}) ([]newrelic.DashboardWidget, error) {
	if len(widgets) < 1 {
		return []newrelic.DashboardWidget{}, nil
	}

	perms := make([]newrelic.DashboardWidget, len(widgets))

	for i, rawCfg := range widgets {
		cfg := rawCfg.(map[string]interface{})
		expandedWidget, err := expandWidget(cfg)

		if err != nil {
			return nil, err
		}

		perms[i] = *expandedWidget
	}

	return perms, nil
}

func expandWidget(cfg map[string]interface{}) (*newrelic.DashboardWidget, error) {
	widget := &newrelic.DashboardWidget{
		Visualization: cfg["visualization"].(string),
		ID:            cfg["widget_id"].(int),
	}

	err := validateWidgetData(cfg)
	if err != nil {
		return nil, err
	}

	expandedLayout, err := expandWidgetLayout(cfg)
	if err != nil {
		return nil, err
	}

	widget.Data = expandWidgetData(cfg)
	widget.Presentation = expandWidgetPresentation(cfg)
	widget.Layout = *expandedLayout

	return widget, nil
}

func validateWidgetData(cfg map[string]interface{}) error {
	visualization := cfg["visualization"].(string)

	switch visualization {
	case "gauge":
		if nrql, ok := cfg["nrql"]; !ok || nrql.(string) == "" {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}
		if red, ok := cfg["threshold_red"]; !ok || red.(float64) == 0 {
			return fmt.Errorf("threshold_red is required for %s visualization", visualization)
		}
	case "billboard", "billboard_comparison":
		if nrql, ok := cfg["nrql"]; !ok || nrql.(string) == "" {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}
	case "facet_bar_chart", "faceted_line_chart", "facet_pie_chart", "facet_table", "faceted_area_chart", "heatmap":
		if nrql, ok := cfg["nrql"]; !ok || nrql.(string) == "" {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}
	case "attribute_sheet", "single_event", "histogram", "funnel", "raw_json", "event_feed", "event_table", "uniques_list", "line_chart", "comparison_line_chart":
		if nrql, ok := cfg["nrql"]; !ok || nrql.(string) == "" {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}
	case "markdown":
		if source, ok := cfg["source"]; !ok || source.(string) == "" {
			return fmt.Errorf("source is required for %s visualization", visualization)
		}
	case "metric_line_chart":
		if _, ok := cfg["metric"]; !ok {
			return fmt.Errorf("metric is required for %s visualization", visualization)
		}
		if _, ok := cfg["entity_ids"]; !ok {
			return fmt.Errorf("entity_ids is required for %s visualization", visualization)
		}
		if _, ok := cfg["duration"]; !ok {
			return fmt.Errorf("duration is required for %s visualization", visualization)
		}
	case "application_breakdown":
		if _, ok := cfg["entity_ids"]; !ok {
			return fmt.Errorf("entity_ids is required for %s visualization", visualization)
		}
	}

	return nil
}

func expandWidgetData(cfg map[string]interface{}) []newrelic.DashboardWidgetData {
	widgetData := newrelic.DashboardWidgetData{}

	if nrql, ok := cfg["nrql"]; ok {
		widgetData.NRQL = nrql.(string)
	}

	if source, ok := cfg["source"]; ok {
		widgetData.Source = source.(string)
	}

	if duration, ok := cfg["duration"]; ok {
		widgetData.Duration = duration.(int)
	}

	if endTime, ok := cfg["end_time"]; ok {
		widgetData.EndTime = endTime.(int)
	}

	if rawMetricName, ok := cfg["raw_metric_name"]; ok {
		widgetData.RawMetricName = rawMetricName.(string)
	}

	if facet, ok := cfg["facet"]; ok {
		widgetData.Facet = facet.(string)
	}

	if orderBy, ok := cfg["order_by"]; ok {
		widgetData.OrderBy = orderBy.(string)
	}

	if limit, ok := cfg["limit"]; ok {
		widgetData.Limit = limit.(int)
	}

	if metrics, ok := cfg["metric"]; ok && metrics.(*schema.Set).Len() > 0 {
		widgetData.Metrics = expandWidgetDataMetrics(metrics.(*schema.Set).List())
	}

	if entityIds, ok := cfg["entity_ids"]; ok && entityIds.(*schema.Set).Len() > 0 {
		widgetData.EntityIds = expandIntSet(entityIds.(*schema.Set))
	}

	if compareWith, ok := cfg["compare_with"]; ok {
		widgetData.CompareWith = expandWidgetDataCompareWith(compareWith.(*schema.Set).List())
	}

	// widget data is a slice for legacy reasons
	return []newrelic.DashboardWidgetData{widgetData}
}

func expandWidgetDataMetrics(metrics []interface{}) []newrelic.DashboardWidgetDataMetric {
	if len(metrics) < 1 {
		return []newrelic.DashboardWidgetDataMetric{}
	}

	perms := make([]newrelic.DashboardWidgetDataMetric, len(metrics))

	for i, rawCfg := range metrics {
		cfg := rawCfg.(map[string]interface{})

		metric := newrelic.DashboardWidgetDataMetric{
			Name: cfg["name"].(string),
		}
		if values, ok := cfg["values"]; ok {
			metric.Values = expandStringSet(values.(*schema.Set))
		}
		if units, ok := cfg["units"]; ok {
			metric.Units = units.(string)
		}
		if scope, ok := cfg["limit"]; ok {
			metric.Scope = scope.(string)
		}

		perms[i] = metric
	}

	return perms
}

func expandWidgetDataCompareWith(windows []interface{}) []newrelic.DashboardWidgetDataCompareWith {
	if len(windows) < 1 {
		return []newrelic.DashboardWidgetDataCompareWith{}
	}

	perms := make([]newrelic.DashboardWidgetDataCompareWith, len(windows))

	for i, rawCfg := range windows {
		cfg := rawCfg.(map[string]interface{})

		perms[i] = newrelic.DashboardWidgetDataCompareWith{
			OffsetDuration: cfg["offset_duration"].(string),
			Presentation:   expandWidgetDataCompareWithPresentation(cfg["presentation"].([]interface{})[0].(map[string]interface{})),
		}
	}

	return perms
}

func expandWidgetDataCompareWithPresentation(cfg map[string]interface{}) newrelic.DashboardWidgetDataCompareWithPresentation {
	widgetDataCompareWithPresentation := newrelic.DashboardWidgetDataCompareWithPresentation{
		Name:  cfg["name"].(string),
		Color: cfg["color"].(string),
	}

	return widgetDataCompareWithPresentation
}

func expandWidgetPresentation(cfg map[string]interface{}) newrelic.DashboardWidgetPresentation {
	widgetPresentation := newrelic.DashboardWidgetPresentation{
		Title: cfg["title"].(string),
	}

	if n, ok := cfg["notes"]; ok {
		widgetPresentation.Notes = n.(string)
	}

	if d, ok := cfg["drilldown_dashboard_id"]; ok {
		widgetPresentation.DrilldownDashboardID = d.(int)
	}

	widgetThreshold := &newrelic.DashboardWidgetThreshold{}

	if red, ok := cfg["threshold_red"]; ok {
		widgetThreshold.Red = red.(float64)
	}

	if yellow, ok := cfg["threshold_yellow"]; ok {
		widgetThreshold.Yellow = yellow.(float64)
	}

	widgetPresentation.Threshold = widgetThreshold

	return widgetPresentation
}

func expandWidgetLayout(cfg map[string]interface{}) (*newrelic.DashboardWidgetLayout, error) {
	widgetLayout := &newrelic.DashboardWidgetLayout{
		Row:    cfg["row"].(int),
		Column: cfg["column"].(int),
		Width:  cfg["width"].(int),
		Height: cfg["height"].(int),
	}

	return widgetLayout, nil
}

// Unpack the *newrelic.Dashboard variable and set resource data.
//
// Used by the newrelic_dashboard Read function (resourceNewRelicDashboardRead)
func flattenDashboard(dashboard *newrelic.Dashboard, d *schema.ResourceData) error {
	d.Set("title", dashboard.Title)
	d.Set("icon", dashboard.Icon)
	d.Set("visibility", dashboard.Visibility)
	d.Set("editable", dashboard.Editable)
	d.Set("dashboard_url", dashboard.UIURL)

	if filterErr := d.Set("filter", flattenFilter(&dashboard.Filter)); filterErr != nil {
		return filterErr
	}

	if dashboard.Widgets != nil && len(dashboard.Widgets) > 0 {
		if widgetErr := d.Set("widget", flattenWidgets(&dashboard.Widgets)); widgetErr != nil {
			return widgetErr
		}
	}

	return nil
}

func flattenWidgets(in *[]newrelic.DashboardWidget) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(*in), len(*in))
	for i, w := range *in {
		m := make(map[string]interface{})
		m["widget_id"] = w.ID
		m["visualization"] = w.Visualization
		m["title"] = w.Presentation.Title
		m["notes"] = w.Presentation.Notes
		m["row"] = w.Layout.Row
		m["column"] = w.Layout.Column
		m["width"] = w.Layout.Width
		m["height"] = w.Layout.Height

		if w.Presentation.DrilldownDashboardID > 0 {
			m["drilldown_dashboard_id"] = w.Presentation.DrilldownDashboardID
		}

		if w.Presentation.Threshold != nil {
			threshold := w.Presentation.Threshold

			if threshold.Red > 0 {
				m["threshold_red"] = threshold.Red
			}

			if threshold.Yellow > 0 {
				m["threshold_yellow"] = threshold.Yellow
			}
		}

		if w.Data != nil && len(w.Data) > 0 {
			data := w.Data[0]

			if data.NRQL != "" {
				m["nrql"] = data.NRQL
			}

			if data.Source != "" {
				m["source"] = data.Source
			}

			if data.Duration > 0 {
				m["duration"] = data.Duration
			}

			if data.EndTime > 0 {
				m["end_time"] = data.EndTime
			}

			if data.RawMetricName != "" {
				m["raw_metric_name"] = data.RawMetricName
			}

			if data.Facet != "" {
				m["facet"] = data.Facet
			}

			if data.OrderBy != "" {
				m["order_by"] = data.OrderBy
			}

			if data.Limit > 0 {
				m["limit"] = data.Limit
			}

			if data.EntityIds != nil && len(data.EntityIds) > 0 {
				m["entity_ids"] = data.EntityIds
			}

			if data.CompareWith != nil && len(data.CompareWith) > 0 {
				m["compare_with"] = flattenWidgetDataCompareWith(data.CompareWith)
			}

			if data.Metrics != nil && len(data.Metrics) > 0 {
				m["metric"] = flattenWidgetDataMetrics(data.Metrics)
			}
		}

		out[i] = m
	}

	return out
}

func flattenWidgetDataCompareWith(in []newrelic.DashboardWidgetDataCompareWith) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})

		m["offset_duration"] = v.OffsetDuration
		m["presentation"] = flattenWidgetDataCompareWithPresentation(&v.Presentation)

		out[i] = m
	}

	return out
}

func flattenWidgetDataCompareWithPresentation(in *newrelic.DashboardWidgetDataCompareWithPresentation) interface{} {
	m := make(map[string]interface{})

	m["name"] = in.Name
	m["color"] = in.Color

	return []interface{}{m}
}

func flattenWidgetDataMetrics(in []newrelic.DashboardWidgetDataMetric) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})

		m["name"] = v.Name
		m["units"] = v.Units
		m["scope"] = v.Scope

		if v.Values != nil && len(v.Values) > 0 {
			m["values"] = v.Values
		}

		out[i] = m
	}

	return out
}

func flattenFilter(f *newrelic.DashboardFilter) []interface{} {
	if f == nil {
		return nil
	}

	if len(f.Attributes) == 0 && len(f.EventTypes) == 0 {
		return nil
	}

	filterResult := make(map[string]interface{})

	attributesList := make([]interface{}, 0, len(f.Attributes))
	for _, v := range f.Attributes {
		attributesList = append(attributesList, v)
	}

	eventTypesList := make([]interface{}, 0, len(f.EventTypes))
	for _, v := range f.EventTypes {
		eventTypesList = append(eventTypesList, v)
	}

	filterResult["attributes"] = schema.NewSet(schema.HashString, attributesList)
	filterResult["event_types"] = schema.NewSet(schema.HashString, eventTypesList)
	return []interface{}{filterResult}
}
