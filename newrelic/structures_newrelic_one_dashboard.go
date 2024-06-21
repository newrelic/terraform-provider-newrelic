package newrelic

import (
	"context"
	"encoding/json"
	"errors"
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
func expandDashboardInput(d *schema.ResourceData, meta interface{}, dashboardNameCustom string) (*dashboards.DashboardInput, error) {
	var err error

	dash := dashboards.DashboardInput{}

	if dashboardNameCustom != "" {
		dash.Name = dashboardNameCustom
	} else {
		dash.Name = d.Get("name").(string)
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

		if options, ok := v["options"]; ok && len(options.([]interface{})) > 0 {
			variable.Options = expandVariableOptions(options.([]interface{}))
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
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for widgetIndex, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for widgetIndex, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.line")
				if err != nil {
					return nil, err
				}

				// Set thresholds
				rawConfiguration.Thresholds = expandDashboardLineWidgetConfigurationThresholdInput(d, pageIndex, widgetIndex)
				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_markdown"]; ok {
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for widgetIndex, v := range widgets.([]interface{}) {
				// Get generic properties set
				widget, rawConfiguration, err := expandDashboardWidgetInput(v.(map[string]interface{}), meta, "viz.table")
				if err != nil {
					return nil, err
				}

				// Set thresholds
				rawConfiguration.Thresholds = expandDashboardTableWidgetConfigurationThresholdInput(d, pageIndex, widgetIndex)
				widget.RawConfiguration, err = json.Marshal(rawConfiguration)
				if err != nil {
					return nil, err
				}

				page.Widgets = append(page.Widgets, *widget)
			}
		}
		if widgets, ok := p["widget_log_table"]; ok {
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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
			for _, v := range widgets.([]interface{}) {
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

func expandDashboardLineWidgetConfigurationThresholdInput(d *schema.ResourceData, pageIndex int, widgetIndex int) dashboards.DashboardLineWidgetThresholdInput {
	// initialize a root object of the DashboardLineWidgetThresholdInput class, which is expected to include IsLabelVisible and Thresholds
	var lineWidgetThresholdsRoot dashboards.DashboardLineWidgetThresholdInput

	// check if 'is_label_visible' has been specified in the configuration of the 'widget_line' widget currently referenced
	// if so, assign the specified value of 'is_label_visible' to IsLabelVisible in the object created above

	lineWidget, lineWidgetOk := d.GetOk(fmt.Sprintf("page.%d.widget_line.%d", pageIndex, widgetIndex))
	if lineWidgetOk {
		lineWidgetAttributes := lineWidget.(map[string]interface{})
		isLabelVisible := lineWidgetAttributes["is_label_visible"]
		if isLabelVisible != nil {
			isLabelVisibleBoolean := isLabelVisible.(bool)
			lineWidgetThresholdsRoot.IsLabelVisible = &isLabelVisibleBoolean
		}
	}

	// initialize a list of 'DashboardLineWidgetThresholdThresholdInput', which would be populated with thresholds specified in the configuration
	// and eventually assigned to the 'Thresholds' attribute of the root object of 'DashboardLineWidgetThresholdInput' specified above

	var lineWidgetThresholdsToBeAdded []dashboards.DashboardLineWidgetThresholdThresholdInput

	// check if 'threshold' has been specified in the configuration of the 'widget_line' widget currently referenced
	// if so, continue with additional logic specified below to iterate through the configuration to find each
	// 'threshold' block specified and get values of attributes specified in each threshold to add to the object of DashboardLineWidgetThresholdThresholdInput

	lineWidgetThresholdsInInput, lineWidgetThresholdsInInputOk := d.GetOk(fmt.Sprintf("page.%d.widget_line.%d.threshold", pageIndex, widgetIndex))
	if lineWidgetThresholdsInInputOk {
		// convert the thresholds obtained into a list of interfaces, in order to fetch the number of threshold blocks specified
		// using the length of this list, "threshold" blocks in the 'widget_line' widget shall be iterated through to get values
		// specified in each threshold

		lineWidgetThresholdsInInputInterface := lineWidgetThresholdsInInput.([]interface{})
		for i := 0; i < len(lineWidgetThresholdsInInputInterface); i++ {
			lineWidgetThresholdInInputSingular, lineWidgetInputThresholdSingularOk := d.GetOk(fmt.Sprintf("page.%d.widget_line.%d.threshold.%d", pageIndex, widgetIndex, i))
			lineWidgetThresholdInInputSingularInterface := lineWidgetThresholdInInputSingular.(map[string]interface{})

			// initialize a DashboardLineWidgetThresholdThresholdInput object to which the values found in the current threshold shall
			// be assigned to respective attributes. Multiple such objects would land into the 'lineWidgetThresholdsToBeAdded' list specified above,
			// which shall eventually be assigned to the 'Thresholds' attribute of the root object of 'DashboardLineWidgetThresholdInput' specified above
			lineWidgetThresholdToBeAdded := dashboards.DashboardLineWidgetThresholdThresholdInput{}

			if lineWidgetInputThresholdSingularOk {
				// if the specified threshold exists, obtain values of attributes specified in the "threshold" block of line widgets
				// and assign them to respective attributes of the DashboardLineWidgetThresholdThresholdInput object

				if v, ok := lineWidgetThresholdInInputSingularInterface["from"]; ok {
					t := v.(int)
					lineWidgetThresholdToBeAdded.From = &t
				}
				if v, ok := lineWidgetThresholdInInputSingularInterface["to"]; ok {
					t := v.(int)
					lineWidgetThresholdToBeAdded.To = &t
				}
				if v, ok := lineWidgetThresholdInInputSingularInterface["name"]; ok {
					lineWidgetThresholdToBeAdded.Name = v.(string)
				}
				if v, ok := lineWidgetThresholdInInputSingularInterface["severity"]; ok {
					lineWidgetThresholdToBeAdded.Severity = dashboards.DashboardLineTableWidgetsAlertSeverity(v.(string))
				}

				// add the threshold to the list of thresholds to be added
				lineWidgetThresholdsToBeAdded = append(lineWidgetThresholdsToBeAdded, lineWidgetThresholdToBeAdded)
			}
		}
	}

	// assign the specified thresholds to 'Thresholds' in the root object created above
	lineWidgetThresholdsRoot.Thresholds = lineWidgetThresholdsToBeAdded

	return lineWidgetThresholdsRoot
}

func expandDashboardTableWidgetConfigurationThresholdInput(d *schema.ResourceData, pageIndex int, widgetIndex int) []dashboards.DashboardTableWidgetThresholdInput {
	// initialize an object of []DashboardTableWidgetThresholdInput, which would include a list of tableWidgetThresholdsToBeAdded as specified
	// in the Terraform configuration, with the attribute "threshold" in table widgets
	var tableWidgetThresholdsToBeAdded []dashboards.DashboardTableWidgetThresholdInput

	// check if 'threshold' has been specified in the configuration of the 'widget_table' widget currently referenced
	// if so, continue with additional logic specified below to iterate through the configuration to find each
	// 'threshold' block specified and get values of attributes specified in each threshold to add to the object of DashboardTableWidgetThresholdInput

	tableWidgetThresholdsInInput, tableWidgetThresholdsInInputOk := d.GetOk(fmt.Sprintf("page.%d.widget_table.%d.threshold", pageIndex, widgetIndex))
	if tableWidgetThresholdsInInputOk {
		// convert the thresholds obtained into a list of interfaces, in order to fetch the number of threshold blocks specified
		// using the length of this list, "threshold" blocks in the 'widget_table' widget shall be iterated through to get values
		// specified in each threshold

		tableWidgetThresholdsInInputInterface := tableWidgetThresholdsInInput.([]interface{})
		for i := 0; i < len(tableWidgetThresholdsInInputInterface); i++ {
			tableWidgetThresholdInInputSingular, tableWidgetThresholdInInputSingularOk := d.GetOk(fmt.Sprintf("page.%d.widget_table.%d.threshold.%d", pageIndex, widgetIndex, i))
			tableWidgetThresholdInInputSingularInterface := tableWidgetThresholdInInputSingular.(map[string]interface{})

			// initialize a DashboardTableWidgetThresholdInput object to which the values found in the current threshold shall
			// be assigned to respective attributes. Multiple such objects would land into the 'tableWidgetThresholdsToBeAdded' list specified above

			tableWidgetThresholdToBeAdded := dashboards.DashboardTableWidgetThresholdInput{}
			if tableWidgetThresholdInInputSingularOk {
				// if the specified threshold exists, obtain values of attributes specified in the "threshold" block of table widgets
				// and assign them to respective attributes of the DashboardTableWidgetThresholdInput object

				if v, ok := tableWidgetThresholdInInputSingularInterface["from"]; ok {
					t := v.(int)
					tableWidgetThresholdToBeAdded.From = &t
				}
				if v, ok := tableWidgetThresholdInInputSingularInterface["to"]; ok {
					t := v.(int)
					tableWidgetThresholdToBeAdded.To = &t
				}
				if v, ok := tableWidgetThresholdInInputSingularInterface["column_name"]; ok {
					tableWidgetThresholdToBeAdded.ColumnName = v.(string)
				}
				if v, ok := tableWidgetThresholdInInputSingularInterface["severity"]; ok {
					tableWidgetThresholdToBeAdded.Severity = dashboards.DashboardLineTableWidgetsAlertSeverity(v.(string))
				}
				tableWidgetThresholdsToBeAdded = append(tableWidgetThresholdsToBeAdded, tableWidgetThresholdToBeAdded)
			}
		}
	}

	return tableWidgetThresholdsToBeAdded
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

	cfg = expandDashboardWidgetYAxisAttributesVizClassified(w, cfg, visualisation)
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

func expandDashboardWidgetYAxisAttributesVizClassified(w map[string]interface{}, cfg dashboards.RawConfiguration, visualisation string) dashboards.RawConfiguration {
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

		lineWidgetYAxisRight := expandDashboardWidgetYAxisRight(w)
		if lineWidgetYAxisRight.Series != nil || lineWidgetYAxisRight.Zero != nil {
			cfg.YAxisRight = &lineWidgetYAxisRight
		}
	}
	return cfg
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

func expandDashboardWidgetYAxisRight(w map[string]interface{}) dashboards.DashboardWidgetYAxisRight {
	// create an object of 'DashboardWidgetYAxisRight', to which we would assign values of attributes associated with
	// 'y_axis_right' in the configuration of the line widget, and return from this function
	dashboardYAxisRightToBeAdded := dashboards.DashboardWidgetYAxisRight{}

	if q, ok := w["y_axis_right"]; ok && len(q.([]interface{})) > 0 {
		// if "y_axis_right" exists in the Terraform configuration and is not empty, proceed with the logic needed
		// to parse "y_axis_right_zero" and "y_axis_right_series", and assign them to respective attributes

		dashboardYAxisRightInInput := q.([]interface{})[0].(map[string]interface{})

		// if "y_axis_right_zero" exists in the map derived from "y_axis_right" (above), fetch the value assigned to
		// this and assign it to the attribute "Zero" of the 'DashboardWidgetYAxisRight' object
		if dashboardYAxisRightZeroInInput, dashboardYAxisRightZeroInInputOk := dashboardYAxisRightInInput["y_axis_right_zero"]; dashboardYAxisRightZeroInInputOk {
			dashboardYAxisRightZeroInInputBoolean := dashboardYAxisRightZeroInInput.(bool)
			dashboardYAxisRightToBeAdded.Zero = &dashboardYAxisRightZeroInInputBoolean
		}

		// if "y_axis_right_series" exists in the map derived from "y_axis_right" (above), fetch the value assigned to
		// this (which is a list of strings) and marshal it accordingly into expected structures within 'DashboardWidgetYAxisRight'
		if dashboardYAxisRightSeriesInInput, dashboardYAxisRightSeriesInInputOk := dashboardYAxisRightInInput["y_axis_right_series"]; dashboardYAxisRightSeriesInInputOk {
			var dashboardYAxisRightSeriesToBeAdded []dashboards.DashboardWidgetYAxisRightSeries
			dashboardYAxisRightSeriesInInputAsList := dashboardYAxisRightSeriesInInput.(*schema.Set).List()
			for _, item := range dashboardYAxisRightSeriesInInputAsList {
				dashboardYAxisRightSeriesToBeAdded = append(
					dashboardYAxisRightSeriesToBeAdded,
					dashboards.DashboardWidgetYAxisRightSeries{
						Name: dashboards.DashboardWidgetYAxisRightSeriesName(item.(string)),
					},
				)
			}

			// eventually, assign the marshalled series from the above logic to the attribute "Series" of the 'DashboardWidgetYAxisRight' object
			dashboardYAxisRightToBeAdded.Series = dashboardYAxisRightSeriesToBeAdded
		}

		if *dashboardYAxisRightToBeAdded.Zero {
			if yMin, okMin := dashboardYAxisRightInInput["y_axis_right_min"]; okMin {
				if yMin.(float64) != 0 {
					min := yMin.(float64)
					dashboardYAxisRightToBeAdded.Min = &min
				}
			}
			if yMax, okMax := dashboardYAxisRightInInput["y_axis_right_max"]; okMax {
				if yMax.(float64) != 0 {
					dashboardYAxisRightToBeAdded.Max = yMax.(float64)
				}
			}
		} else {
			if yMin, okMin := dashboardYAxisRightInInput["y_axis_right_min"]; okMin {
				min := yMin.(float64)
				dashboardYAxisRightToBeAdded.Min = &min
				if yMax, okMax := dashboardYAxisRightInInput["y_axis_right_max"]; okMax {
					dashboardYAxisRightToBeAdded.Max = yMax.(float64)
				}
			}
		}

	}

	// return the 'DashboardWidgetYAxisRight' object into which the contents of "y_axis_right" in the Terraform configuration have been repackaged
	return dashboardYAxisRightToBeAdded
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

func expandVariableOptions(in []interface{}) *dashboards.DashboardVariableOptionsInput {
	var out dashboards.DashboardVariableOptionsInput

	for _, v := range in {
		cfg := v.(map[string]interface{})
		out = dashboards.DashboardVariableOptionsInput{
			IgnoreTimeRange: cfg["ignore_time_range"].(bool),
		}
	}

	return &out
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
		variables := flattenDashboardVariable(&dashboard.Variables, d)
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
		variables := flattenDashboardVariable(&dashboard.Variables, d)
		if err := d.Set("variable", variables); err != nil {
			return err
		}
	}

	return nil
}

func flattenDashboardVariable(in *[]entities.DashboardVariable, d *schema.ResourceData) []interface{} {
	out := make([]interface{}, len(*in))

	for i, v := range *in {
		m := make(map[string]interface{})

		if v.DefaultValues != nil {
			m["default_values"] = flattenVariableDefaultValues(&v.DefaultValues)
		}
		m["is_multi_selection"] = v.IsMultiSelection
		m["item"] = flattenVariableItems(v.Items)
		m["name"] = v.Name
		if &v.NRQLQuery != nil {
			m["nrql_query"] = flattenVariableNRQLQuery(&v.NRQLQuery)
		}
		m["replacement_strategy"] = strings.ToLower(string(v.ReplacementStrategy))
		m["title"] = v.Title
		m["type"] = strings.ToLower(string(v.Type))
		if &v.Options != nil {
			options := flattenVariableOptions(&v.Options, d, i)
			if options != nil {
				// set options -> ignore_time_range to the state only if they already exist in the configuration
				// needed to make this backward compatible with configurations which do not yet have options -> ignore_time_range
				m["options"] = options
			}
		}

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

func flattenVariableOptions(in *entities.DashboardVariableOptions, d *schema.ResourceData, index int) []interface{} {
	// fetching the contents of the variable at the specified index (in the list of variables)
	// and subsequently, finding its options field
	variableFetched := d.Get(fmt.Sprintf("variable.%d", index)).(map[string]interface{})
	optionsOfVariableFetched, optionsOfVariableFetchedOk := variableFetched["options"]
	if !optionsOfVariableFetchedOk {
		return nil
	}
	options := optionsOfVariableFetched.([]interface{})
	if len(options) == 0 {
		// if nothing exists in the options list in the state (configuration), "do nothing", to avoid drift
		// this is required to make options -> ignore_time_range backward compatible and show no drift
		// when customer configuration does not comprise these attributes
		return nil
	}

	// else, only if the options list has something, get the value of ignore_time_range
	// (set it to the value of ignore_time_range seen in the response returned by the API)
	out := make([]interface{}, 1)
	n := make(map[string]interface{})
	n["ignore_time_range"] = in.IgnoreTimeRange
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
		if rawCfg.Thresholds != nil {
			rawCfgThresholdsFetched := rawCfg.Thresholds.([]interface{})
			if len(rawCfgThresholdsFetched) > 0 {
				for _, t := range rawCfgThresholdsFetched {
					thresholdFetched := t.(map[string]interface{})
					if thresholdFetched["value"] == nil {
						continue
					}

					switch thresholdFetched["alertSeverity"].(string) {
					case string(entities.DashboardAlertSeverityTypes.CRITICAL):
						out["critical"] = strconv.FormatFloat(thresholdFetched["value"].(float64), 'f', -1, 64)
					case string(entities.DashboardAlertSeverityTypes.WARNING):
						out["warning"] = strconv.FormatFloat(thresholdFetched["value"].(float64), 'f', -1, 64)
					}
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

		// check if 'YAxisRight' in the rawConfiguration of the fetched widget is not null. If it isn't, proceed
		// with extracting values out of the fetched 'YAxisRight' and assigning them to respective attributes in thje configuration
		if rawCfg.YAxisRight != nil {
			out["y_axis_right"] = flattenDashboardLineWidgetYAxisRight(rawCfg.YAxisRight)
		}

		if rawCfg.Thresholds != nil {
			isLabelVisible, thresholds := flattenDashboardLineWidgetThresholds(rawCfg.Thresholds)
			if isLabelVisible != nil {
				out["is_label_visible"] = isLabelVisible.(bool)
			}
			if thresholds != nil {
				out["threshold"] = thresholds
			}
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
		if rawCfg.Thresholds != nil {
			thresholds := flattenDashboardTableWidgetThresholds(rawCfg.Thresholds)
			if thresholds != nil {
				out["threshold"] = thresholds
			}
		}

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

func flattenDashboardLineWidgetYAxisRight(yAxisRight *dashboards.DashboardWidgetYAxisRight) []interface{} {
	// define a map[string]interface{} which would hold key-value pairs (keys: 'y_axis_right_zero' and 'y_axis_right_series')
	// define an []interface{} that would hold the map defined above, which can then, be assigned to "y_axis_right" as it
	// would expect an []interface{}
	var yAxisRightFetched = make(map[string]interface{})
	var yAxisRightFetchedInterface []interface{}

	// if 'Zero' exists yAxisRight 'YAxisRight', assign it to "y_axis_right_zero" of the map
	if yAxisRight.Zero != nil {
		yAxisRightFetched["y_axis_right_zero"] = yAxisRight.Zero
	}

	// assign 'Min' and 'Max' of obtained to their respective attributes in the Terraform configuration
	yAxisRightFetched["y_axis_right_min"] = yAxisRight.Min
	yAxisRightFetched["y_axis_right_max"] = yAxisRight.Max

	// if 'Series' exists yAxisRight 'YAxisRight', assign it to "y_axis_right_series" of the map
	if yAxisRight.Series != nil {
		var yAxisRightSeriesFetched []string
		for _, item := range yAxisRight.Series {
			yAxisRightSeriesFetched = append(yAxisRightSeriesFetched, string(item.Name))
		}
		yAxisRightFetched["y_axis_right_series"] = yAxisRightSeriesFetched
	}

	// add the map containing 'Zero' and/or 'Series' to the []interface{}
	// eventually, assign the []interface{} to "y_axis_right" of the referenced widget
	yAxisRightFetchedInterface = append(yAxisRightFetchedInterface, yAxisRightFetched)
	return yAxisRightFetchedInterface
}

func flattenDashboardLineWidgetThresholds(thresholds interface{}) (interface{}, []map[string]interface{}) {
	var thresholdsConsolidated []map[string]interface{}
	var isLabelVisible interface{}

	thresholdsFetched := thresholds.(map[string]interface{})
	if thresholdsFetched["isLabelVisible"] != nil {
		isLabelVisible = thresholdsFetched["isLabelVisible"]
	}

	thresholdsFetchedList := thresholdsFetched["thresholds"]
	if thresholdsFetchedList != nil {
		thresholdsFetchedListInterface := thresholdsFetchedList.([]interface{})
		if len(thresholdsFetchedListInterface) > 0 {
			for _, item := range thresholdsFetchedListInterface {
				thresholdSingle := item.(map[string]interface{})
				thresholdSingleToBeFormatted := map[string]interface{}{}

				//t := reflect.TypeOf(dashboards.DashboardLineWidgetThresholdThresholdInput{})
				//for i := 0; i < t.NumField(); i++ {
				//	field := t.Field(i)
				//	name := strings.Split(field.Tag.Get("json"), ",")[0]
				//	if val, ok := thresholdSingular[name]; ok {
				//		newt[name] = val
				//	}
				//}

				for key, terraformSchemaKey := range lineWidgetThresholdAttributesJSON {
					if value, ok := thresholdSingle[key]; ok {
						thresholdSingleToBeFormatted[terraformSchemaKey] = value
					}
				}
				thresholdsConsolidated = append(thresholdsConsolidated, thresholdSingleToBeFormatted)
			}
		}
	}
	return isLabelVisible, thresholdsConsolidated
}

func flattenDashboardTableWidgetThresholds(thresholds interface{}) []map[string]interface{} {
	var thresholdsConsolidated []map[string]interface{}
	thresholdsFetchedListInterface := thresholds.([]interface{})
	if len(thresholdsFetchedListInterface) > 0 {
		for _, item := range thresholdsFetchedListInterface {
			thresholdSingle := item.(map[string]interface{})
			thresholdSingleToBeFormatted := map[string]interface{}{}
			for key, terraformSchemaKey := range tableWidgetThresholdAttributesJSON {
				if value, ok := thresholdSingle[key]; ok {
					thresholdSingleToBeFormatted[terraformSchemaKey] = value
				}
			}
			thresholdsConsolidated = append(thresholdsConsolidated, thresholdSingleToBeFormatted)
		}
	}
	return thresholdsConsolidated
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
				for _, k := range widgets.([]interface{}) {
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

func validateDashboardArguments(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []string

	err := validateDashboardVariableOptions(d)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	// add any other validation functions here

	if len(errorsList) == 0 {
		return nil
	}

	errorsString := "the following validation errors have been identified: \n"

	for index, val := range errorsList {
		errorsString += fmt.Sprintf("(%d): %s\n", index+1, val)
	}

	return errors.New(errorsString)
}

func validateDashboardVariableOptions(d *schema.ResourceDiff) error {
	_, variablesListObtained := d.GetChange("variable")
	vars := variablesListObtained.([]interface{})

	for _, v := range vars {
		variableMap := v.(map[string]interface{})
		options, optionsOk := variableMap["options"]
		if optionsOk {
			optionsInterface := options.([]interface{})
			if len(optionsInterface) > 1 {
				return errors.New("only one set of `options` may be specified per variable")
			}
			for _, o := range optionsInterface {
				if o == nil {
					return errors.New("`options` block(s) specified cannot be empty")
				}
				optionMap := o.(map[string]interface{})
				_, ignoreTimeRangeOk := optionMap["ignore_time_range"]
				variableType, variableTypeOk := variableMap["type"]
				if ignoreTimeRangeOk && variableTypeOk && variableType != "nrql" {
					return errors.New("`ignore_time_range` in `options` can only be used with the variable type `nrql`")
				}

			}
		}
	}

	return nil
}
