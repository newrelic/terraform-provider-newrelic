package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
)

// migrateStateV0toV1 currently facilitates migrating the `widgets`
// attribute from TypeList to TypeSet. Since the underlying
// data structure is []map[string]interface{} for both, we don't
// need to do anything other than return the state and Terraform
// will handle the rest.
func migrateStateV0toV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}

// Assemble the *dashboards.Dashboard struct.
// Used by the newrelic_dashboard Create and Update functions.
func expandDashboard(d *schema.ResourceData) (*dashboards.Dashboard, error) {
	metadata := dashboards.DashboardMetadata{
		Version: 1,
	}

	dashboard := dashboards.Dashboard{
		Title:           d.Get("title").(string),
		Metadata:        metadata,
		Icon:            dashboards.DashboardIconType(d.Get("icon").(string)),
		Visibility:      dashboards.VisibilityType(d.Get("visibility").(string)),
		Editable:        dashboards.EditableType(d.Get("editable").(string)),
		GridColumnCount: dashboards.GridColumnCountType(d.Get("grid_column_count").(int)),
	}

	if f, ok := d.GetOk("filter"); ok {
		dashboard.Filter = expandFilter(f.([]interface{})[0].(map[string]interface{}))
	}

	log.Printf("[INFO] widget schema: %+v\n", d.Get("widget"))
	if widgets, ok := d.GetOk("widget"); ok {
		expandedWidgets, err := expandWidgets(widgets.(*schema.Set).List())

		if err != nil {
			return nil, err
		}

		dashboard.Widgets = expandedWidgets
	}

	return &dashboard, nil
}

func expandFilter(filter map[string]interface{}) dashboards.DashboardFilter {
	perms := dashboards.DashboardFilter{}

	if v, ok := filter["attributes"]; ok {
		perms.Attributes = expandStringSet(v.(*schema.Set))
	}

	if v, ok := filter["event_types"]; ok {
		perms.EventTypes = expandStringSet(v.(*schema.Set))
	}

	return perms
}

func expandWidgets(widgets []interface{}) ([]dashboards.DashboardWidget, error) {
	if len(widgets) < 1 {
		return []dashboards.DashboardWidget{}, nil
	}

	perms := make([]dashboards.DashboardWidget, len(widgets))

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

func expandWidget(cfg map[string]interface{}) (*dashboards.DashboardWidget, error) {
	widget := &dashboards.DashboardWidget{
		Visualization: dashboards.VisualizationType(cfg["visualization"].(string)),
		ID:            cfg["widget_id"].(int),
	}

	if accountID, ok := cfg["account_id"]; ok {
		widget.AccountID = accountID.(int)
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

// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
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

func expandWidgetData(cfg map[string]interface{}) []dashboards.DashboardWidgetData {
	widgetData := dashboards.DashboardWidgetData{}

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

	if metrics, ok := cfg["metric"]; ok {
		widgetData.Metrics = expandWidgetDataMetrics(metrics.(*schema.Set).List())
	}

	if entityIds, ok := cfg["entity_ids"]; ok {
		widgetData.EntityIds = expandIntSet(entityIds.(*schema.Set))
	}

	if compareWith, ok := cfg["compare_with"]; ok {
		widgetData.CompareWith = expandWidgetDataCompareWith(compareWith.(*schema.Set).List())
	}

	// widget data is a slice for legacy reasons
	return []dashboards.DashboardWidgetData{widgetData}
}

func expandWidgetDataMetrics(metrics []interface{}) []dashboards.DashboardWidgetDataMetric {
	if len(metrics) < 1 {
		return []dashboards.DashboardWidgetDataMetric{}
	}

	perms := make([]dashboards.DashboardWidgetDataMetric, len(metrics))

	for i, rawCfg := range metrics {
		cfg := rawCfg.(map[string]interface{})

		metric := dashboards.DashboardWidgetDataMetric{
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

func expandWidgetDataCompareWith(windows []interface{}) []dashboards.DashboardWidgetDataCompareWith {
	if len(windows) < 1 {
		return []dashboards.DashboardWidgetDataCompareWith{}
	}

	perms := make([]dashboards.DashboardWidgetDataCompareWith, len(windows))

	for i, rawCfg := range windows {
		cfg := rawCfg.(map[string]interface{})

		perms[i] = dashboards.DashboardWidgetDataCompareWith{
			OffsetDuration: cfg["offset_duration"].(string),
			Presentation:   expandWidgetDataCompareWithPresentation(cfg["presentation"].([]interface{})[0].(map[string]interface{})),
		}
	}

	return perms
}

func expandWidgetDataCompareWithPresentation(cfg map[string]interface{}) dashboards.DashboardWidgetDataCompareWithPresentation {
	widgetDataCompareWithPresentation := dashboards.DashboardWidgetDataCompareWithPresentation{
		Name:  cfg["name"].(string),
		Color: cfg["color"].(string),
	}

	return widgetDataCompareWithPresentation
}

func expandWidgetPresentation(cfg map[string]interface{}) dashboards.DashboardWidgetPresentation {
	widgetPresentation := dashboards.DashboardWidgetPresentation{
		Title: cfg["title"].(string),
	}

	if n, ok := cfg["notes"]; ok {
		widgetPresentation.Notes = n.(string)
	}

	if d, ok := cfg["drilldown_dashboard_id"]; ok {
		widgetPresentation.DrilldownDashboardID = d.(int)
	}

	widgetThreshold := &dashboards.DashboardWidgetThreshold{}

	if red, ok := cfg["threshold_red"]; ok {
		widgetThreshold.Red = red.(float64)
	}

	if yellow, ok := cfg["threshold_yellow"]; ok {
		widgetThreshold.Yellow = yellow.(float64)
	}

	widgetPresentation.Threshold = widgetThreshold

	return widgetPresentation
}

func expandWidgetLayout(cfg map[string]interface{}) (*dashboards.DashboardWidgetLayout, error) {
	widgetLayout := &dashboards.DashboardWidgetLayout{
		Row:    cfg["row"].(int),
		Column: cfg["column"].(int),
		Width:  cfg["width"].(int),
		Height: cfg["height"].(int),
	}

	return widgetLayout, nil
}

// Unpack the *dashboards.Dashboard variable and set resource data.
//
// Used by the newrelic_dashboard Read function (resourceNewRelicDashboardRead)
func flattenDashboard(dashboard *dashboards.Dashboard, d *schema.ResourceData) error {
	d.Set("title", dashboard.Title)
	d.Set("icon", dashboard.Icon)
	d.Set("visibility", dashboard.Visibility)
	d.Set("editable", dashboard.Editable)
	d.Set("dashboard_url", dashboard.UIURL)

	if gridColumnCount, ok := d.GetOk("grid_column_count"); ok {
		d.Set("grid_column_count", gridColumnCount.(int))
	} else {
		d.Set("grid_column_count", 3)
	}

	if filterErr := d.Set("filter", flattenFilter(&dashboard.Filter)); filterErr != nil {
		return filterErr
	}

	if dashboard.Widgets != nil && len(dashboard.Widgets) > 0 {
		if widgetErr := d.Set("widget", flattenWidgets(&dashboard.Widgets, d)); widgetErr != nil {
			return widgetErr
		}
	}

	return nil
}

// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func flattenWidgets(in *[]dashboards.DashboardWidget, d *schema.ResourceData) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(*in))

	// wgts := d.Get("widget")
	// configuredWidgets := wgts.(*schema.Set).List()

	for i, w := range *in {
		m := make(map[string]interface{})

		m["widget_id"] = w.ID

		// The REST API returns cross-account widgets
		// as "inaccessible" so we can't really do anything
		// with that. An inaccessible widget does not
		// contain enough data in the response to bother setting
		// it in the state. As a consequence, if customers want to
		// use cross-account widgets, they will need to either
		// ignore the lifecycle changes for `widget` block or
		// use `-auto-approve` when running `terraform apply`.
		if w.Visualization != "inaccessible" {
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
		}

		out[i] = m
	}

	return out
}

func flattenWidgetDataCompareWith(in []dashboards.DashboardWidgetDataCompareWith) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in))
	for i, v := range in {
		m := make(map[string]interface{})

		m["offset_duration"] = v.OffsetDuration
		m["presentation"] = flattenWidgetDataCompareWithPresentation(&v.Presentation)

		out[i] = m
	}

	return out
}

func flattenWidgetDataCompareWithPresentation(in *dashboards.DashboardWidgetDataCompareWithPresentation) interface{} {
	m := make(map[string]interface{})

	m["name"] = in.Name
	m["color"] = in.Color

	return []interface{}{m}
}

func flattenWidgetDataMetrics(in []dashboards.DashboardWidgetDataMetric) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in))
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

func flattenFilter(f *dashboards.DashboardFilter) []interface{} {
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
