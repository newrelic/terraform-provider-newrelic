package newrelic

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
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

	dash.Variables = expandDashboardVariablesInput(d.Get("variable").([]interface{}))

	// Optional, with default
	perm := d.Get("permissions").(string)
	dash.Permissions = entities.DashboardPermissions(strings.ToUpper(perm))

	// Optional
	if e, ok := d.GetOk("description"); ok {
		dash.Description = e.(string)
	}

	return &dash, nil
}

func expandDashboardVariablesInput(variables []interface{}) []dashboards.DashboardVariableInput {
	if len(variables) < 1 {
		return []dashboards.DashboardVariableInput{}
	}

	expanded := make([]dashboards.DashboardVariableInput, len(variables))

	for i, val := range variables {
		var variable dashboards.DashboardVariableInput
		v := val.(map[string]interface{})

		if d, ok := v["default_values"]; ok && len(d.([]interface{})) > 0 {
			variable.DefaultValues = expandVariableDefaultValues(d.([]interface{}))
		}

		if m, ok := v["is_multi_selection"]; ok {
			variable.IsMultiSelection = m.(bool)
		}

		if i, ok := v["item"]; ok && len(i.([]interface{})) > 0 {
			variable.Items = expandVariableItems(i.([]interface{}))
		}

		if n, ok := v["name"]; ok {
			variable.Name = n.(string)
		}

		if q, ok := v["nrql_query"]; ok && len(q.([]interface{})) > 0 {
			variable.NRQLQuery = expandVariableNRQLQuery(q.([]interface{}))
		}

		if r, ok := v["replacement_strategy"]; ok {
			variable.ReplacementStrategy = dashboards.DashboardVariableReplacementStrategy(strings.ToUpper(r.(string)))
		}

		if t, ok := v["title"]; ok {
			variable.Title = t.(string)
		}

		if ty, ok := v["type"]; ok {
			variable.Type = dashboards.DashboardVariableType(strings.ToUpper(ty.(string)))
		}

		expanded[i] = variable
	}
	return expanded
}

func expandVariableDefaultValues(in []interface{}) *[]dashboards.DashboardVariableDefaultItemInput {
	out := make([]dashboards.DashboardVariableDefaultItemInput, len(in))

	for i, v := range in {
		cfg := v.(string)
		expanded := dashboards.DashboardVariableDefaultItemInput{Value: dashboards.DashboardVariableDefaultValueInput{String: cfg}}
		out[i] = expanded
	}

	return &out
}

func expandVariableItems(in []interface{}) []dashboards.DashboardVariableEnumItemInput {
	out := make([]dashboards.DashboardVariableEnumItemInput, len(in))

	for i, v := range in {
		cfg := v.(map[string]interface{})
		expanded := dashboards.DashboardVariableEnumItemInput{
			Title: cfg["title"].(string),
			Value: cfg["value"].(string),
		}
		out[i] = expanded
	}

	return out
}

func expandVariableNRQLQuery(in []interface{}) *dashboards.DashboardVariableNRQLQueryInput {
	var out dashboards.DashboardVariableNRQLQueryInput

	for _, v := range in {
		cfg := v.(map[string]interface{})
		out = dashboards.DashboardVariableNRQLQueryInput{
			AccountIDs: expandVariableAccountIDs(cfg["account_ids"].([]interface{})),
			Query:      nrdb.NRQL(cfg["query"].(string))}
	}

	return &out
}

func expandVariableAccountIDs(in []interface{}) []int {
	out := make([]int, len(in))

	for i := range out {
		out[i] = in[i].(int)
	}

	return out
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

		page.Widgets = []dashboards.DashboardWidgetInput{}
		// For each of the widget type, we need to expand them as well
		if widgets, ok := p["widget_area"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.area")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_bar"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.bar")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_billboard"]; ok {
			for widgetIndex, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.billboard")
				if err != nil {
					return nil, err
				}

				// Set thresholds
				rawConfiguration.Thresholds = expandDashboardBillboardWidgetConfigurationInput(d, v.(map[string]interface{}), meta, pageIndex, widgetIndex)

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_bullet"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.bullet")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_funnel"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.funnel")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_heatmap"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.heatmap")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_histogram"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.histogram")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_line"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.line")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_markdown"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.markdown")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_pie"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.pie")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_table"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.table")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_log_table"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "logger.log-table-widget")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_json"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.json")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_stacked_bar"]; ok {
			for _, v := range widgets.(*schema.Set).List() {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.stacked-bar")
				if err != nil {
					return nil, err
				}

				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}

		expanded[pageIndex] = page
	}

	return expanded, nil
}

func expandDashboardBillboardWidgetConfigurationInput(d *schema.ResourceData, i map[string]interface{}, meta interface{}, pageIndex int, widgetIndex int) []dashboards.DashboardBillboardWidgetThresholdInput {
	// optional, order is important (API returns them sorted alpha)
	var thresholds = []dashboards.DashboardBillboardWidgetThresholdInput{}
	if data, ok := d.GetOk(fmt.Sprintf("page.%d.widget_billboard.%d.critical", pageIndex, widgetIndex)); ok {
		value := data.(string)
		if value != "" {
			floatValue, _ := strconv.ParseFloat(value, 64)
			thresholds = append(thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
				AlertSeverity: entities.DashboardAlertSeverityTypes.CRITICAL,
				Value:         &floatValue,
			})
		} else {
			thresholds = append(thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
				AlertSeverity: entities.DashboardAlertSeverityTypes.CRITICAL,
				Value:         nil,
			})
		}
	}

	if data, ok := d.GetOk(fmt.Sprintf("page.%d.widget_billboard.%d.warning", pageIndex, widgetIndex)); ok {
		value := data.(string)
		if value != "" {
			floatValue, _ := strconv.ParseFloat(value, 64)
			thresholds = append(thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
				AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
				Value:         &floatValue,
			})
		} else {
			thresholds = append(thresholds, dashboards.DashboardBillboardWidgetThresholdInput{
				AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
				Value:         nil,
			})
		}
	}

	return thresholds
}

// expandDashboardWidgetInput expands the common items in WidgetInput, but not the configuration
// which is specific to the widgets
func expandDashboardWidgetInput(w map[string]interface{}, meta interface{}, visualisation string) (*dashboards.DashboardWidgetInput, *dashboards.RawConfiguration, error) {
	var widget dashboards.DashboardWidgetInput
	var err error
	var cfg dashboards.RawConfiguration

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

	if q, ok := w["nrql_query"]; ok {
		cfg.NRQLQueries, err = expandDashboardWidgetNRQLQueryInput(q.([]interface{}), meta)
		if err != nil {
			return nil, nil, err
		}
	}

	if q, ok := w["legend_enabled"]; ok {
		var l dashboards.DashboardWidgetLegend
		l.Enabled = q.(bool)
		cfg.Legend = &l
	}
	if q, ok := w["facet_show_other_series"]; ok {
		var l dashboards.DashboardWidgetFacet
		l.ShowOtherSeries = q.(bool)
		cfg.Facet = &l
	}

	if visualisation != "viz.line" {
		if q, ok := w["y_axis_left_min"]; ok {
			var l dashboards.DashboardWidgetYAxisLeft
			min := q.(float64)
			l.Min = &min
			if q, ok := w["y_axis_left_max"]; ok {
				l.Max = q.(float64)
			}
			cfg.YAxisLeft = &l
		}
	} else {
		lineWidgetYAxisLeft := expandDashboardWidgetYAxisLeft(w)
		cfg.YAxisLeft = &lineWidgetYAxisLeft
	}

	cfg = expandDashboardWidgetNullValuesInput(w, cfg)
	cfg = expandDashboardWidgetColorsInput(w, cfg)
	cfg = expandDashboardWidgetUnitsInput(w, cfg)

	if l, ok := w["limit"]; ok {
		cfg.Limit = l.(float64)
	}

	if l, ok := w["ignore_time_range"]; ok {
		var platformOptions = dashboards.RawConfigurationPlatformOptions{}
		platformOptions.IgnoreTimeRange = l.(bool)
		cfg.PlatformOptions = &platformOptions
	}

	if t, ok := w["text"]; ok {
		if t.(string) != "" {
			cfg.Text = t.(string)
		}
	}

	widget.Visualization.ID = visualisation

	return &widget, &cfg, nil
}

func expandDashboardWidgetYAxisLeft(w map[string]interface{}) dashboards.DashboardWidgetYAxisLeft {
	var l dashboards.DashboardWidgetYAxisLeft
	if q, ok := w["y_axis_left_zero"]; ok {
		yAxisZero := q.(bool)
		l.Zero = &yAxisZero
		if !yAxisZero {
			if yMin, okMin := w["y_axis_left_min"]; okMin {
				if yMin.(float64) != 0 {
					min := yMin.(float64)
					l.Min = &min
				}
			}
			if yMax, okMax := w["y_axis_left_max"]; okMax {
				if yMax.(float64) != 0 {
					l.Max = yMax.(float64)
				}
			}
		} else {
			if yMin, okMin := w["y_axis_left_min"]; okMin {
				min := yMin.(float64)
				l.Min = &min
				if yMax, okMax := w["y_axis_left_max"]; okMax {
					l.Max = yMax.(float64)
				}
			}
		}
	}

	return l
}
func expandDashboardWidgetUnitsInput(w map[string]interface{}, cfg dashboards.RawConfiguration) dashboards.RawConfiguration {
	if q, ok := w["units"]; ok {
		units := q.([]interface{})
		var n dashboards.DashboardWidgetUnits
		if len(units) > 0 {
			for _, y := range units {
				if y != nil {
					z := y.(map[string]interface{})
					if v, ok := z["unit"]; ok {
						n.Unit = v.(string)
					}
					if s, ok := z["series_overrides"]; ok {
						var seriesOverrides = s.([]interface{})
						n.SeriesOverrides = make([]dashboards.DashboardWidgetUnitOverrides, len(seriesOverrides))
						for i, v := range seriesOverrides {
							var t dashboards.DashboardWidgetUnitOverrides
							k := v.(map[string]interface{})
							if n, ok := k["unit"]; ok {
								t.Unit = n.(string)
							}
							if n, ok := k["series_name"]; ok {
								t.SeriesName = n.(string)
							}
							n.SeriesOverrides[i] = t
						}
					}
				}
			}
			cfg.Units = &n
		}
	}
	return cfg
}

func expandDashboardWidgetColorsInput(w map[string]interface{}, cfg dashboards.RawConfiguration) dashboards.RawConfiguration {
	if q, ok := w["colors"]; ok {
		colors := q.([]interface{})
		var n dashboards.DashboardWidgetColors
		if len(colors) > 0 {
			for _, y := range colors {
				if y != nil {
					z := y.(map[string]interface{})
					if v, ok := z["color"]; ok {
						n.Color = v.(string)
					}
					if s, ok := z["series_overrides"]; ok {
						var seriesOverrides = s.([]interface{})
						n.SeriesOverrides = make([]dashboards.DashboardWidgetColorOverrides, len(seriesOverrides))
						for i, v := range seriesOverrides {
							var t dashboards.DashboardWidgetColorOverrides
							k := v.(map[string]interface{})
							if n, ok := k["color"]; ok {
								t.Color = n.(string)
							}
							if n, ok := k["series_name"]; ok {
								t.SeriesName = n.(string)
							}
							n.SeriesOverrides[i] = t
						}

					}
				}
			}
			cfg.Colors = &n
		}
	}
	return cfg
}

func expandDashboardWidgetNullValuesInput(w map[string]interface{}, cfg dashboards.RawConfiguration) dashboards.RawConfiguration {
	if q, ok := w["null_values"]; ok {
		nullValues := q.([]interface{})
		if len(nullValues) > 0 {
			var n dashboards.DashboardWidgetNullValues
			for _, y := range nullValues {
				if y != nil {
					z := y.(map[string]interface{})
					if v, ok := z["null_value"]; ok {
						n.NullValue = v.(string)
					}
					if s, ok := z["series_overrides"]; ok {
						var seriesOverrides = s.([]interface{})
						n.SeriesOverrides = make([]dashboards.DashboardWidgetNullValueOverrides, len(seriesOverrides))
						for i, v := range seriesOverrides {
							var t dashboards.DashboardWidgetNullValueOverrides
							k := v.(map[string]interface{})
							if n, ok := k["null_value"]; ok {
								t.NullValue = n.(string)
							}
							if n, ok := k["series_name"]; ok {
								t.SeriesName = n.(string)
							}
							n.SeriesOverrides[i] = t
						}

					}
				}
				cfg.NullValues = &n
			}

		}
	}
	return cfg
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

	if dashboard.Variables != nil && len(dashboard.Variables) > 0 {
		variables := flattenDashboardVariable(&dashboard.Variables)
		if err := d.Set("variable", variables); err != nil {
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

	if dashboard.Variables != nil && len(dashboard.Variables) > 0 {
		variables := flattenDashboardVariable(&dashboard.Variables)
		if err := d.Set("variable", variables); err != nil {
			return err
		}
	}

	return nil
}

func flattenDashboardVariable(in *[]entities.DashboardVariable) []interface{} {
	out := make([]interface{}, len(*in))

	for i, v := range *in {
		m := make(map[string]interface{})

		if v.DefaultValues != nil {
			m["default_values"] = flattenVariableDefaultValues(v.DefaultValues)
		}
		m["is_multi_selection"] = v.IsMultiSelection
		m["item"] = flattenVariableItems(v.Items)
		m["name"] = v.Name
		if v.NRQLQuery != nil {
			m["nrql_query"] = flattenVariableNRQLQuery(v.NRQLQuery)
		}
		m["replacement_strategy"] = strings.ToLower(string(v.ReplacementStrategy))
		m["title"] = v.Title
		m["type"] = strings.ToLower(string(v.Type))

		out[i] = m
	}
	return out
}

func flattenVariableDefaultValues(in *[]entities.DashboardVariableDefaultItem) []string {
	out := make([]string, len(*in))

	for i, v := range *in {
		out[i] = v.Value.String
	}
	return out
}

func flattenVariableItems(in []entities.DashboardVariableEnumItem) []interface{} {
	out := make([]interface{}, len(in))

	for i, v := range in {
		item := make(map[string]interface{})
		item["title"] = v.Title
		item["value"] = v.Value

		out[i] = item
	}
	return out
}

func flattenVariableNRQLQuery(in *entities.DashboardVariableNRQLQuery) []interface{} {
	out := make([]interface{}, 1)

	n := make(map[string]interface{})

	n["account_ids"] = in.AccountIDs
	n["query"] = in.Query

	out[0] = n

	return out
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
			widgetType, w := flattenDashboardWidget(&widget, string(p.GUID))

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
func flattenDashboardWidget(in *entities.DashboardWidget, pageGUID string) (string, map[string]interface{}) {
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

	var filterCurrentDashboard = false
	if out["linked_entity_guids"] != nil && len(out["linked_entity_guids"].([]string)) == 1 && stringInSlice(out["linked_entity_guids"].([]string), pageGUID) {
		filterCurrentDashboard = true
	}

	// Read out the rawConfiguration field for use in all widgets
	rawCfg := dashboards.RawConfiguration{}
	if len(in.RawConfiguration) > 0 {
		if err := json.Unmarshal(in.RawConfiguration, &rawCfg); err != nil {
			log.Printf("Error parsing: %s", err)
		}
	}

	// Set global raw configuration fields
	if rawCfg.PlatformOptions != nil {
		out["ignore_time_range"] = rawCfg.PlatformOptions.IgnoreTimeRange
	}

	if rawCfg.Legend != nil {
		out["legend_enabled"] = rawCfg.Legend.Enabled
	}
	if rawCfg.Facet != nil {
		out["facet_show_other_series"] = rawCfg.Facet.ShowOtherSeries
	}
	if rawCfg.YAxisLeft != nil {
		out["y_axis_left_min"] = rawCfg.YAxisLeft.Min
		out["y_axis_left_max"] = rawCfg.YAxisLeft.Max
	}
	if rawCfg.NullValues != nil {
		out["null_values"] = flattenDashboardWidgetNullValues(rawCfg.NullValues)
	}
	if rawCfg.Units != nil {
		out["units"] = flattenDashboardWidgetUnits(rawCfg.Units)
	}
	if rawCfg.Colors != nil {
		out["colors"] = flattenDashboardWidgetColors(rawCfg.Colors)
	}

	// Set widget type and arguments
	switch in.Visualization.ID {
	case "viz.area":
		widgetType = "widget_area"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.bar":
		widgetType = "widget_bar"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
		out["filter_current_dashboard"] = filterCurrentDashboard
	case "viz.billboard":
		widgetType = "widget_billboard"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
		if len(rawCfg.Thresholds) > 0 {
			for _, v := range rawCfg.Thresholds {
				// Double check if we have a value, the API sometimes returns a null
				if v.Value == nil {
					continue
				}

				switch v.AlertSeverity {
				case entities.DashboardAlertSeverityTypes.CRITICAL:
					out["critical"] = strconv.FormatFloat(*v.Value, 'f', -1, 64)
				case entities.DashboardAlertSeverityTypes.WARNING:
					out["warning"] = strconv.FormatFloat(*v.Value, 'f', -1, 64)
				}
			}
		}
	case "viz.bullet":
		widgetType = "widget_bullet"
		out["limit"] = rawCfg.Limit
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.funnel":
		widgetType = "widget_funnel"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.heatmap":
		widgetType = "widget_heatmap"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.histogram":
		widgetType = "widget_histogram"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.json":
		widgetType = "widget_json"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.line":
		widgetType = "widget_line"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
		if rawCfg.YAxisLeft != nil {
			out["y_axis_left_zero"] = rawCfg.YAxisLeft.Zero
		}
	case "viz.markdown":
		widgetType = "widget_markdown"
		out["text"] = rawCfg.Text
	case "viz.stacked-bar":
		widgetType = "widget_stacked_bar"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	case "viz.pie":
		widgetType = "widget_pie"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
		out["filter_current_dashboard"] = filterCurrentDashboard
	case "viz.table":
		widgetType = "widget_table"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
		out["filter_current_dashboard"] = filterCurrentDashboard
	case "logger.log-table-widget":
		widgetType = "widget_log_table"
		out["nrql_query"] = flattenDashboardWidgetNRQLQuery(&rawCfg.NRQLQueries)
	}

	return widgetType, out
}

func flattenDashboardWidgetNRQLQuery(in *[]dashboards.DashboardWidgetNRQLQueryInput) []interface{} {
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
				for _, widget := range widgets.(*schema.Set).List() {
					w := widget.(map[string]interface{})
					if v, ok := w["filter_current_dashboard"]; ok && v.(bool) {

						if l, ok := w["linked_entity_guids"]; ok && len(l.([]interface{})) > 1 {
							return nil, fmt.Errorf("filter_current_dashboard can't be set if linked_entity_guids is configured")
						}

						if l, ok := w["linked_entity_guids"]; ok && len(l.([]interface{})) == 1 {
							for _, le := range l.([]interface{}) {
								if le.(string) != p["guid"] {
									return nil, fmt.Errorf("filter_current_dashboard can't be set if linked_entity_guids is configured")
								}
							}
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
	selfLinkingWidgets := []string{"widget_bar", "widget_pie", "widget_table"}

	pages := d.Get("page").([]interface{})
	for i, v := range pages {
		p := v.(map[string]interface{})
		for _, widgetType := range selfLinkingWidgets {
			if widgets, ok := p[widgetType]; ok {
				for _, k := range widgets.(*schema.Set).List() {
					w := k.(map[string]interface{})
					if l, ok := w["linked_entity_guids"]; ok && len(l.([]interface{})) == 1 {
						for _, le := range l.([]interface{}) {
							if f, ok := w["filter_current_dashboard"]; ok && f == false && le.(string) == p["guid"] {
								w["linked_entity_guids"] = nil
							}
						}
					}
					for _, f := range filterWidgets {
						e := f.(map[string]interface{})
						if e["page"] == i {
							if w["title"] == e["title"] && w["column"] == e["column"] && w["row"] == e["row"] {
								guid := p["guid"].(string)
								if guid == "" {
									w["linked_entity_guids"] = nil
								} else {
									w["linked_entity_guids"] = []string{guid}
								}
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
func flattenDashboardWidgetNullValues(in *dashboards.DashboardWidgetNullValues) interface{} {
	out := make([]interface{}, 1)
	k := make(map[string]interface{})
	k["null_value"] = in.NullValue
	seriesOverrides := make([]interface{}, len(in.SeriesOverrides))
	for i, v := range in.SeriesOverrides {
		m := make(map[string]interface{})

		m["null_value"] = v.NullValue
		m["series_name"] = v.SeriesName

		seriesOverrides[i] = m
	}
	k["series_overrides"] = seriesOverrides
	out[0] = k
	return out
}

func flattenDashboardWidgetColors(in *dashboards.DashboardWidgetColors) interface{} {
	out := make([]interface{}, 1)
	k := make(map[string]interface{})
	k["color"] = in.Color
	seriesOverrides := make([]interface{}, len(in.SeriesOverrides))
	for i, v := range in.SeriesOverrides {
		m := make(map[string]interface{})

		m["color"] = v.Color
		m["series_name"] = v.SeriesName

		seriesOverrides[i] = m
	}
	k["series_overrides"] = seriesOverrides
	out[0] = k
	return out
}

func flattenDashboardWidgetUnits(in *dashboards.DashboardWidgetUnits) interface{} {
	out := make([]interface{}, 1)
	k := make(map[string]interface{})
	k["unit"] = in.Unit
	seriesOverrides := make([]interface{}, len(in.SeriesOverrides))
	for i, v := range in.SeriesOverrides {
		m := make(map[string]interface{})

		m["unit"] = v.Unit
		m["series_name"] = v.SeriesName

		seriesOverrides[i] = m
	}
	k["series_overrides"] = seriesOverrides
	out[0] = k
	return out
}
