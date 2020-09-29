package newrelic

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
)

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
		expandedWidgets, err := expandWidgets(widgets)
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

func expandWidgets(in interface{}) ([]dashboards.DashboardWidget, error) {
	widgetsIn := in.([]interface{})
	if len(widgetsIn) < 1 {
		return []dashboards.DashboardWidget{}, nil
	}

	expanded := make([]dashboards.DashboardWidget, len(widgetsIn))

	for i, wg := range widgetsIn {
		w := wg.(map[string]interface{})

		expandedWidget, err := expandWidget(w)
		if err != nil {
			return nil, err
		}

		expanded[i] = *expandedWidget
	}

	return expanded, nil
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

	// IMPORTANT! Sorting the widgets before storing in state helps prevent drift
	// in multiple scenarios, such as when/if the API returns widgets in a different
	// order or if the user changes the order the HCL resource configuration.
	sort.SliceStable(dashboard.Widgets, func(i, j int) bool {
		return dashboard.Widgets[i].ID < dashboard.Widgets[j].ID
	})

	if dashboard.Widgets != nil && len(dashboard.Widgets) > 0 {
		if widgetErr := d.Set("widget", flattenWidgets(&dashboard.Widgets, d)); widgetErr != nil {
			return widgetErr
		}
	}

	return nil
}

func isValidViz(viz dashboards.VisualizationType) bool {
	vizString := string(viz)
	for _, vizType := range validWidgetVisualizationValues {
		if vizString == vizType {
			return true
		}
	}

	return false
}

func flattenWidgets(widgetsIn *[]dashboards.DashboardWidget, d *schema.ResourceData) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(*widgetsIn))

	widgetCfg, ok := d.GetOk("widget")

	for i, w := range *widgetsIn {
		if !ok {
			// If not widgets are configured, we need
			// to provide an empty map to populate
			// using the incoming API response widget data.
			wgt := map[string]interface{}{}
			out[i] = flattenWidget(w, wgt)
		} else {
			widgetConfig := widgetCfg.([]interface{})
			wgt := widgetConfig[i].(map[string]interface{})
			out[i] = flattenWidget(w, wgt)
		}
	}

	return out
}

// SUPPORTING CROSS-ACCOUNT WIDGETS WITH THE REST API USING APIKS KEYS
//
// If a user sets `account_id` to a subaccount that's scoped outside of
// the user's API key, the API returns the widget as "inaccessible" and omits data.
// This function attempts to avoid configuration drift when certain configuration
// scenarios are presented.
//
// If a user sets `account_id` that's "inaccessible" per the API, we avoid setting
// with the API response data and just use the same ID the user provided in the HCL.
//
// If the user sets `account_id` to an accessible account associated with the API key,
// we need to set this in the Terraform state to avoid drift.
//
// If the user does not set `account_id`, then the user is basically using the default
// behavior and we don't need to set it in the state since it's not in the HCL.
// nolint:gocyclo
func flattenWidget(w dashboards.DashboardWidget, widgetCfg map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	wgtConfigAcctID := getConfiguredWidgetAcctID(widgetCfg)

	if wgtConfigAcctID > 0 && w.AccountID != 0 {
		m["account_id"] = w.AccountID
	} else {
		m["account_id"] = wgtConfigAcctID
	}

	if w.ID != 0 {
		m["widget_id"] = w.ID
	}

	// Cross-account widgets will have a visualization
	// set to "inaccessible" in some cases, so we must
	// ensure a valid visualization is provided in the
	// API's widget response.
	if isValidViz(w.Visualization) {
		m["visualization"] = w.Visualization
	} else {
		m["visualization"] = widgetCfg["visualization"]
	}

	if w.Presentation.Title != "" {
		m["title"] = w.Presentation.Title
	} else {
		m["title"] = widgetCfg["title"]
	}

	if w.Presentation.Notes != "" {
		m["notes"] = w.Presentation.Notes
	} else {
		m["notes"] = widgetCfg["notes"]
	}

	if w.Layout.Row != 0 {
		m["row"] = w.Layout.Row
	} else {
		m["row"] = widgetCfg["row"]
	}

	if w.Layout.Column != 0 {
		m["column"] = w.Layout.Column
	} else {
		m["column"] = widgetCfg["column"]
	}

	if w.Layout.Width != 0 {
		m["width"] = w.Layout.Width
	} else {
		m["width"] = widgetCfg["width"]
	}

	if w.Layout.Height != 0 {
		m["height"] = w.Layout.Height
	} else {
		m["height"] = widgetCfg["height"]
	}

	if w.Presentation.DrilldownDashboardID > 0 {
		m["drilldown_dashboard_id"] = w.Presentation.DrilldownDashboardID
	} else {
		m["drilldown_dashboard_id"] = widgetCfg["drilldown_dashboard_id"]
	}

	if w.Presentation.Threshold != nil {
		threshold := w.Presentation.Threshold

		if threshold.Red > 0 {
			m["threshold_red"] = threshold.Red
		} else {
			m["threshold_red"] = widgetCfg["threshold_red"]
		}

		if threshold.Yellow > 0 {
			m["threshold_yellow"] = threshold.Yellow
		} else {
			m["threshold_yellow"] = widgetCfg["threshold_yellow"]
		}
	} else {
		m["threshold_red"] = widgetCfg["threshold_red"]
		m["threshold_yellow"] = widgetCfg["threshold_yellow"]
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
	} else {
		m["nrql"] = widgetCfg["nrql"]
	}

	return m
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

// A helper function to get the value of the `account_id` from the HCL attribute if it was set.
func getConfiguredWidgetAcctID(widgetConfig map[string]interface{}) int {
	val, ok := widgetConfig["account_id"]
	if !ok {
		return 0
	}

	return val.(int)
}
