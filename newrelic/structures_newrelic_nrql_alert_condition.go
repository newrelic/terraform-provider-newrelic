package newrelic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

var (
	// old:new
	timeFunctionMap = map[string]alerts.ThresholdOccurrence{
		"all": alerts.ThresholdOccurrences.All,
		"any": alerts.ThresholdOccurrences.AtLeastOnce,
	}

	// new:old
	timeFunctionMapNewOld = map[alerts.ThresholdOccurrence]string{
		alerts.ThresholdOccurrences.All:         "all",
		alerts.ThresholdOccurrences.AtLeastOnce: "any",
	}

	// old:new
	fillOptionMap = map[string]*alerts.AlertsFillOption{
		"none":       &alerts.AlertsFillOptionTypes.NONE,
		"last_value": &alerts.AlertsFillOptionTypes.LAST_VALUE,
		"static":     &alerts.AlertsFillOptionTypes.STATIC,
	}

	// new:old
	fillOptionMapNewOld = map[alerts.AlertsFillOption]string{
		alerts.AlertsFillOptionTypes.NONE:       "none",
		alerts.AlertsFillOptionTypes.LAST_VALUE: "last_value",
		alerts.AlertsFillOptionTypes.STATIC:     "static",
	}

	// old:new
	aggregationMethodMap = map[string]*alerts.NrqlConditionAggregationMethod{
		"cadence":     &alerts.NrqlConditionAggregationMethodTypes.Cadence,
		"event_flow":  &alerts.NrqlConditionAggregationMethodTypes.EventFlow,
		"event_timer": &alerts.NrqlConditionAggregationMethodTypes.EventTimer,
	}

	// new:old
	aggregationMethodMapNewOld = map[alerts.NrqlConditionAggregationMethod]string{
		alerts.NrqlConditionAggregationMethodTypes.Cadence:    "cadence",
		alerts.NrqlConditionAggregationMethodTypes.EventFlow:  "event_flow",
		alerts.NrqlConditionAggregationMethodTypes.EventTimer: "event_timer",
	}
)

// NerdGraph
func expandNrqlAlertConditionCreateInput(d *schema.ResourceData) (*alerts.NrqlConditionCreateInput, error) {
	input := alerts.NrqlConditionCreateInput{
		NrqlConditionCreateBase: alerts.NrqlConditionCreateBase{
			Description: d.Get("description").(string),
			Enabled:     d.Get("enabled").(bool),
			Name:        d.Get("name").(string),
		},
	}

	conditionType := strings.ToLower(d.Get("type").(string))

	if conditionType == "baseline" {
		if attr, ok := d.GetOk("baseline_direction"); ok {
			direction := alerts.NrqlBaselineDirection(strings.ToUpper(attr.(string)))
			input.BaselineDirection = &direction
		} else {
			return nil, fmt.Errorf("attribute `%s` is required for nrql alert conditions of type `%+v`", "baseline_direction", conditionType)
		}
	}

	if conditionType == "baseline" {
		if attr, ok := d.GetOk("signal_seasonality"); ok {
			seasonality := alerts.NrqlSignalSeasonality(strings.ToUpper(attr.(string)))
			input.SignalSeasonality = &seasonality
		} else {
			// Null is equivalent to the default value of `NEW_RELIC_CALCULATION` in the API
			seasonality := alerts.NrqlSignalSeasonalities.NewRelicCalculation
			input.SignalSeasonality = &seasonality
		}
	}

	if runbookURL, ok := d.GetOk("runbook_url"); ok {
		input.RunbookURL = runbookURL.(string)
	}

	if titleTemplate, ok := d.GetOk("title_template"); ok {
		template := titleTemplate.(string)
		input.TitleTemplate = &template
	}

	if violationTimeLimitSec, ok := d.GetOk("violation_time_limit_seconds"); ok {
		input.ViolationTimeLimitSeconds = violationTimeLimitSec.(int)
	} else if violationTimeLimit, ok := d.GetOk("violation_time_limit"); ok {
		input.ViolationTimeLimit = alerts.NrqlConditionViolationTimeLimit(strings.ToUpper(violationTimeLimit.(string)))
	}

	nrql, err := expandCreateNrql(d, input)
	if err != nil {
		return nil, err
	}

	input.Nrql = *nrql

	terms, err := expandNrqlTerms(d, conditionType)
	if err != nil {
		return nil, err
	}

	input.Terms = terms

	if input.Expiration, err = expandExpiration(d); err != nil {
		return nil, err
	}

	if input.Signal, err = expandCreateSignal(d); err != nil {
		return nil, err
	}

	return &input, nil
}

func expandNrqlAlertConditionUpdateInput(d *schema.ResourceData) (*alerts.NrqlConditionUpdateInput, error) {
	input := alerts.NrqlConditionUpdateInput{
		NrqlConditionUpdateBase: alerts.NrqlConditionUpdateBase{
			Description: d.Get("description").(string),
			Enabled:     d.Get("enabled").(bool),
			Name:        d.Get("name").(string),
		},
	}

	conditionType := strings.ToLower(d.Get("type").(string))

	if conditionType == "baseline" {
		if attr, ok := d.GetOk("baseline_direction"); ok {
			direction := alerts.NrqlBaselineDirection(strings.ToUpper(attr.(string)))
			input.BaselineDirection = &direction
		} else {
			return nil, fmt.Errorf("attribute `%s` is required for nrql alert conditions of type `%+v`", "baseline_direction", conditionType)
		}
	}

	if conditionType == "baseline" {
		if attr, ok := d.GetOk("signal_seasonality"); ok {
			seasonality := alerts.NrqlSignalSeasonality(strings.ToUpper(attr.(string)))
			input.SignalSeasonality = &seasonality
		} else {
			// Null is equivalent to the default value of `NEW_RELIC_CALCULATION` in the API
			// so this effectively allows the user to null out the signal seasonality on update
			seasonality := alerts.NrqlSignalSeasonalities.NewRelicCalculation
			input.SignalSeasonality = &seasonality
		}
	}

	if runbookURL, ok := d.GetOk("runbook_url"); ok {
		input.RunbookURL = runbookURL.(string)
	}

	if titleTemplate, ok := d.GetOk("title_template"); ok {
		template := titleTemplate.(string)
		input.TitleTemplate = &template
	}

	if violationTimeLimitSec, ok := d.GetOk("violation_time_limit_seconds"); ok {
		input.ViolationTimeLimitSeconds = violationTimeLimitSec.(int)
	} else if violationTimeLimit, ok := d.GetOk("violation_time_limit"); ok {
		input.ViolationTimeLimit = alerts.NrqlConditionViolationTimeLimit(strings.ToUpper(violationTimeLimit.(string)))
	}

	nrql, err := expandUpdateNrql(d, input)
	if err != nil {
		return nil, err
	}

	input.Nrql = *nrql

	terms, err := expandNrqlTerms(d, conditionType)
	if err != nil {
		return nil, err
	}

	input.Terms = terms

	if input.Expiration, err = expandExpiration(d); err != nil {
		return nil, err
	}

	if input.Signal, err = expandUpdateSignal(d); err != nil {
		return nil, err
	}

	return &input, nil
}

// NerdGraph
func expandCreateNrql(d *schema.ResourceData, condition alerts.NrqlConditionCreateInput) (*alerts.NrqlConditionCreateQuery, error) {
	var nrql alerts.NrqlConditionCreateQuery

	if nrqlQuery, ok := d.GetOk("nrql.0.query"); ok {
		nrql.Query = nrqlQuery.(string)
	}

	if dataAccountID, ok := d.GetOk("nrql.0.data_account_id"); ok {
		v := dataAccountID.(int)
		nrql.DataAccountId = &v
	}

	if sinceValue, ok := d.GetOk("nrql.0.since_value"); ok {
		sv, err := strconv.Atoi(sinceValue.(string))
		if err != nil {
			return nil, err
		}

		nrql.EvaluationOffset = &sv
	} else if evalOffset, ok := d.GetOk("nrql.0.evaluation_offset"); ok {
		eo := evalOffset.(int)
		nrql.EvaluationOffset = &eo
	}

	return &nrql, nil
}

// NerdGraph
func expandUpdateNrql(d *schema.ResourceData, condition alerts.NrqlConditionUpdateInput) (*alerts.NrqlConditionUpdateQuery, error) {
	var nrql alerts.NrqlConditionUpdateQuery

	if nrqlQuery, ok := d.GetOk("nrql.0.query"); ok {
		nrql.Query = nrqlQuery.(string)
	}

	if dataAccountID, ok := d.GetOk("nrql.0.data_account_id"); ok {
		v := dataAccountID.(int)
		nrql.DataAccountId = &v
	}

	if sinceValue, ok := d.GetOk("nrql.0.since_value"); ok {
		sv, err := strconv.Atoi(sinceValue.(string))
		if err != nil {
			return nil, err
		}

		nrql.EvaluationOffset = &sv
	} else if evalOffset, ok := d.GetOk("nrql.0.evaluation_offset"); ok {
		eo := evalOffset.(int)
		nrql.EvaluationOffset = &eo
	}

	return &nrql, nil
}

// NerdGraph
func expandNrqlConditionTerm(term map[string]interface{}, conditionType, priority string) (*alerts.NrqlConditionTerm, error) {
	var durationIn int
	if attr, ok := term["duration"]; ok {
		durationIn = attr.(int)
	}

	thresholdDurationIn := term["threshold_duration"].(int)

	if durationIn == 0 && thresholdDurationIn == 0 {
		return nil, fmt.Errorf("one of `duration` or `threshold_duration` must be configured for block `term`")
	}

	if durationIn > 0 && thresholdDurationIn > 0 {
		return nil, fmt.Errorf("one of `duration` or `threshold_duration` must be configured for block `term`, but not both")
	}

	operator := alerts.AlertsNRQLConditionTermsOperator(strings.ToUpper(term["operator"].(string)))

	switch conditionType {
	case "baseline":
		if operator != alerts.AlertsNRQLConditionTermsOperatorTypes.ABOVE {
			return nil, fmt.Errorf("only ABOVE operator is allowed for `baseline` condition types")
		}
	}

	var duration int
	if durationIn > 0 {
		duration = durationIn * 60 // convert min to sec
	} else {
		duration = thresholdDurationIn
	}

	// required
	threshold := term["threshold"].(float64)

	thresholdOccurrences, err := expandNrqlThresholdOccurrences(term)
	if err != nil {
		return nil, err
	}

	// If we have not been passed a priority, then we should inspect the term we've received.
	if priority == "" {
		if termPriority, ok := term["priority"].(string); ok {
			if termPriority != "" {
				priority = termPriority
			}
		}
	}

	expandedTerm := alerts.NrqlConditionTerm{
		Operator:             operator,
		Priority:             alerts.NrqlConditionPriority(strings.ToUpper(priority)),
		Threshold:            &threshold,
		ThresholdDuration:    duration,
		ThresholdOccurrences: *thresholdOccurrences,
	}

	if term["disable_health_status_reporting"] != nil {
		disableHealthStatusReporting := term["disable_health_status_reporting"].(bool)
		expandedTerm.DisableHealthStatusReporting = &disableHealthStatusReporting
	}

	if conditionType == "baseline" {
		return &expandedTerm, nil
	}

	// Set prediction fields if they're present
	prediction, err := expandNrqlThresholdPrediction(term)
	if err != nil {
		return nil, err
	}

	expandedTerm.Prediction = prediction

	return &expandedTerm, nil
}

// Terraform config => NerdGraph payload
func expandNrqlThresholdPrediction(term map[string]interface{}) (*alerts.NrqlConditionThresholdPrediction, error) {
	if term["prediction"] == nil {
		return nil, nil
	}

	predictionSet := term["prediction"].(*schema.Set)

	if predictionSet.Len() == 0 {
		return nil, nil
	}

	if predictionSet.Len() > 1 {
		return nil, fmt.Errorf("only one `prediction` is allowed per `term` block")
	}

	predictionMap := predictionSet.List()[0].(map[string]interface{})

	var prediction alerts.NrqlConditionThresholdPrediction

	if predictBy, ok := predictionMap["predict_by"].(int); ok {
		prediction.PredictBy = predictBy
	}

	prediction.PreferPredictionViolation = predictionMap["prefer_prediction_violation"].(bool)

	return &prediction, nil
}

func expandNrqlThresholdOccurrences(term map[string]interface{}) (*alerts.ThresholdOccurrence, error) {
	var timeFunctionIn string
	if attr, ok := term["time_function"]; ok {
		timeFunctionIn = attr.(string)
	}

	thresholdOccurrencesIn := term["threshold_occurrences"].(string)

	if timeFunctionIn == "" && thresholdOccurrencesIn == "" {
		return nil, fmt.Errorf("one of `time_function` or `threshold_occurrences` must be configured for block `term`")
	}

	if timeFunctionIn != "" && thresholdOccurrencesIn != "" {
		return nil, fmt.Errorf("one of `time_function` or `threshold_occurrences` must be configured for block `term`, but not both")
	}

	var thresholdOccurrences alerts.ThresholdOccurrence
	if timeFunctionIn != "" {
		thresholdOccurrences = timeFunctionMap[timeFunctionIn]
	} else {
		thresholdOccurrences = alerts.ThresholdOccurrence(strings.ToUpper(thresholdOccurrencesIn))
	}

	return &thresholdOccurrences, nil
}

// Terraform config => NerdGraph payload
func expandNrqlTerms(d *schema.ResourceData, conditionType string) ([]alerts.NrqlConditionTerm, error) {
	var expandedTerms []alerts.NrqlConditionTerm
	var err error
	var errs []string

	terms := d.Get("term").(*schema.Set).List()

	for _, t := range terms {
		term := t.(map[string]interface{})
		var nrqlConditionTerm *alerts.NrqlConditionTerm

		nrqlConditionTerm, err = expandNrqlConditionTerm(term, conditionType, "")
		if err != nil {
			errs = append(errs, fmt.Sprintf("unable to expand NRQL condition term: %s", err))
		}

		if nrqlConditionTerm != nil {
			expandedTerms = append(expandedTerms, *nrqlConditionTerm)
		}
	}

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ", "))
		return expandedTerms, err
	}

	if len(expandedTerms) == 0 {
		if critical, ok := d.GetOk("critical"); ok {
			x := critical.([]interface{})
			// A critical attribute is a list, but is limited to a single item in the schema.
			if len(x) > 0 {
				single := x[0].(map[string]interface{})

				criticalTerm, err := expandNrqlConditionTerm(single, conditionType, "critical")
				if err != nil {
					return nil, err
				}
				if criticalTerm != nil {
					expandedTerms = append(expandedTerms, *criticalTerm)
				}
			}
		}

		if warning, ok := d.GetOk("warning"); ok {
			x := warning.([]interface{})
			// A warning attribute is a list, but is limited to a single item in the schema.
			if len(x) > 0 {
				single := x[0].(map[string]interface{})

				warningTerm, err := expandNrqlConditionTerm(single, conditionType, "warning")
				if err != nil {
					return nil, err
				}

				if warningTerm != nil {
					expandedTerms = append(expandedTerms, *warningTerm)
				}
			}
		}
	}

	return expandedTerms, nil
}

// NerdGraph
func expandExpiration(d *schema.ResourceData) (*alerts.AlertsNrqlConditionExpiration, error) {
	var expiration alerts.AlertsNrqlConditionExpiration

	expiration.OpenViolationOnExpiration = d.Get("open_violation_on_expiration").(bool)
	expiration.CloseViolationsOnExpiration = d.Get("close_violations_on_expiration").(bool)
	expiration.IgnoreOnExpectedTermination = d.Get("ignore_on_expected_termination").(bool)

	// 0 is not a valid expiration duration so don't set it if it's nonexistent
	if expirationDuration, ok := d.GetOk("expiration_duration"); ok {
		v := expirationDuration.(int)
		expiration.ExpirationDuration = &v
	}

	return &expiration, nil
}

// NerdGraph
func expandCreateSignal(d *schema.ResourceData) (*alerts.AlertsNrqlConditionCreateSignal, error) {
	signal := alerts.AlertsNrqlConditionCreateSignal{
		FillOption: fillOptionMap[strings.ToLower(d.Get("fill_option").(string))],
	}

	// Due to the way that nulls are handled as zeros in Terraform 0.11, add another check that a 0 fill_value
	// can only be applied when the fill_option is static
	if fillValue, ok := d.GetOkExists("fill_value"); ok {
		v := fillValue.(float64)
		if v != 0 || (signal.FillOption != nil && *signal.FillOption == alerts.AlertsFillOptionTypes.STATIC) {
			signal.FillValue = &v
		}
	}

	if aggregationWindow, ok := d.GetOk("aggregation_window"); ok {
		v := aggregationWindow.(int)
		signal.AggregationWindow = &v
	}

	if slideBy, ok := d.GetOk("slide_by"); ok {
		v := slideBy.(int)
		signal.SlideBy = &v
	}

	if _, ok := d.GetOk("aggregation_method"); ok {
		v := aggregationMethodMap[strings.ToLower(d.Get("aggregation_method").(string))]

		signal.AggregationMethod = v
	}

	if aggregationDelay, ok := d.GetOk("aggregation_delay"); ok {
		value := aggregationDelay.(string)
		if value != "" {
			v, _ := strconv.Atoi(value)
			signal.AggregationDelay = &v
		}
	}

	if aggregationTimer, ok := d.GetOk("aggregation_timer"); ok {
		value := aggregationTimer.(string)
		if value != "" {
			v, _ := strconv.Atoi(value)
			signal.AggregationTimer = &v
		}
	}

	if evaluationDelay, ok := d.GetOk("evaluation_delay"); ok {
		value := evaluationDelay.(int)
		signal.EvaluationDelay = &value
	}

	return &signal, nil
}

// NerdGraph
func expandUpdateSignal(d *schema.ResourceData) (*alerts.AlertsNrqlConditionUpdateSignal, error) {
	signal := alerts.AlertsNrqlConditionUpdateSignal{
		FillOption: fillOptionMap[strings.ToLower(d.Get("fill_option").(string))],
	}

	// Due to the way that nulls are handled as zeros in Terraform 0.11, add another check that a 0 fill_value
	// can only be applied when the fill_option is static
	if fillValue, ok := d.GetOkExists("fill_value"); ok {
		v := fillValue.(float64)
		if v != 0 || (signal.FillOption != nil && *signal.FillOption == alerts.AlertsFillOptionTypes.STATIC) {
			signal.FillValue = &v
		}
	}

	if aggregationWindow, ok := d.GetOk("aggregation_window"); ok {
		v := aggregationWindow.(int)
		signal.AggregationWindow = &v
	}

	if slideBy, ok := d.GetOk("slide_by"); ok {
		v := slideBy.(int)
		signal.SlideBy = &v
	}

	if _, ok := d.GetOk("aggregation_method"); ok {
		v := aggregationMethodMap[strings.ToLower(d.Get("aggregation_method").(string))]

		signal.AggregationMethod = v
	}

	if aggregationDelay, ok := d.GetOk("aggregation_delay"); ok {
		value := aggregationDelay.(string)
		if value != "" {
			v, _ := strconv.Atoi(value)
			signal.AggregationDelay = &v
		}
	}

	if aggregationTimer, ok := d.GetOk("aggregation_timer"); ok {
		value := aggregationTimer.(string)
		if value != "" {
			v, _ := strconv.Atoi(value)
			signal.AggregationTimer = &v
		}
	}

	if evaluationDelay, ok := d.GetOk("evaluation_delay"); ok {
		value := evaluationDelay.(int)
		signal.EvaluationDelay = &value
	}

	return &signal, nil
}

// NerdGraph response => Terraform state
func flattenNrqlAlertCondition(accountID int, condition *alerts.NrqlAlertCondition, d *schema.ResourceData) error {
	policyID, err := strconv.Atoi(condition.PolicyID)
	if err != nil {
		return err
	}

	conditionType := strings.ToLower(string(condition.Type))

	_ = d.Set("account_id", accountID)
	_ = d.Set("type", conditionType)
	_ = d.Set("description", condition.Description)
	_ = d.Set("policy_id", policyID)
	_ = d.Set("name", condition.Name)
	_ = d.Set("runbook_url", condition.RunbookURL)
	_ = d.Set("title_template", condition.TitleTemplate)
	_ = d.Set("enabled", condition.Enabled)
	_ = d.Set("entity_guid", condition.EntityGUID)

	if conditionType == "baseline" {
		_ = d.Set("baseline_direction", string(*condition.BaselineDirection))

		if condition.SignalSeasonality != nil {
			_ = d.Set("signal_seasonality", string(*condition.SignalSeasonality))
		}
	}

	configuredNrql := d.Get("nrql.0").(map[string]interface{})
	if err := d.Set("nrql", flattenNrql(condition.Nrql, configuredNrql)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `nrql`: %v", err)
	}

	// setting terms explicitly, critical/warning are not set
	configuredTerms := d.Get("term").(*schema.Set).List()

	conditionTerms := flattenNrqlTerms(condition.Terms, configuredTerms)

	if len(configuredTerms) > 0 {
		if err := d.Set("term", conditionTerms); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `term`: %v", err)
		}
	} else {
		// Handle the named condition priorities.

		for _, term := range conditionTerms {
			switch term["priority"].(string) {
			case "critical":
				t := term
				delete(t, "priority")
				var terms []map[string]interface{}
				terms = append(terms, t)
				if err := d.Set("critical", terms); err != nil {
					return fmt.Errorf("[DEBUG] Error setting nrql alert condition `critical`: %v", err)
				}
			case "warning":
				t := term
				delete(t, "priority")
				var terms []map[string]interface{}
				terms = append(terms, t)
				if err := d.Set("warning", terms); err != nil {
					return fmt.Errorf("[DEBUG] Error setting nrql alert condition `warning`: %v", err)
				}
			}
		}
	}

	if condition.ViolationTimeLimitSeconds != 0 {
		_ = d.Set("violation_time_limit_seconds", condition.ViolationTimeLimitSeconds)
	}
	if condition.ViolationTimeLimit != "" {
		_ = d.Set("violation_time_limit", condition.ViolationTimeLimit)
	}

	if err := flattenExpiration(d, condition.Expiration); err != nil {
		return err
	}

	if err := flattenSignal(d, condition.Signal); err != nil {
		return err
	}

	return nil
}

// NerdGraph
func flattenExpiration(d *schema.ResourceData, expiration *alerts.AlertsNrqlConditionExpiration) error {
	if expiration == nil {
		return nil
	}

	if err := d.Set("open_violation_on_expiration", expiration.OpenViolationOnExpiration); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `open_violation_on_expiration`: %v", err)
	}

	if err := d.Set("close_violations_on_expiration", expiration.CloseViolationsOnExpiration); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `close_violations_on_expiration`: %v", err)
	}

	if err := d.Set("expiration_duration", expiration.ExpirationDuration); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `expiration_duration`: %v", err)
	}

	if err := d.Set("ignore_on_expected_termination", expiration.IgnoreOnExpectedTermination); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `ignore_on_expected_termination`: %v", err)
	}

	return nil
}

// NerdGraph
func flattenSignal(d *schema.ResourceData, signal *alerts.AlertsNrqlConditionSignal) error {
	if signal == nil {
		return nil
	}

	if err := d.Set("aggregation_window", signal.AggregationWindow); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `aggregation_window`: %v", err)
	}
	if signal.SlideBy != nil {
		if err := d.Set("slide_by", signal.SlideBy); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `slide_by`: %v", err)
		}
	}

	if err := d.Set("fill_value", signal.FillValue); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `fill_value`: %v", err)
	}

	if signal.FillOption != nil {
		if err := d.Set("fill_option", fillOptionMapNewOld[*signal.FillOption]); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `fill_option`: %v", err)
		}
	}

	if signal.AggregationMethod != nil {
		if err := d.Set("aggregation_method", aggregationMethodMapNewOld[*signal.AggregationMethod]); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `aggregation_method`: %v", err)
		}
	}

	if signal.AggregationDelay != nil {
		delay := strconv.Itoa(*signal.AggregationDelay)
		if err := d.Set("aggregation_delay", delay); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `aggregation_delay`: %v", err)
		}
	}

	if signal.AggregationTimer != nil {
		timer := strconv.Itoa(*signal.AggregationTimer)
		if err := d.Set("aggregation_timer", timer); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `aggregation_timer`: %v", err)
		}
	}

	if signal.EvaluationDelay != nil {
		if err := d.Set("evaluation_delay", signal.EvaluationDelay); err != nil {
			return fmt.Errorf("[DEBUG] Error setting nrql alert condition `evaluation_delay`: %v", err)
		}

	}

	return nil
}

// NerdGraph
func flattenNrql(nrql alerts.NrqlConditionQuery, configNrql map[string]interface{}) []interface{} {
	out := map[string]interface{}{
		"query": nrql.Query,
	}

	if nrql.DataAccountId != nil {
		out["data_account_id"] = nrql.DataAccountId
	}

	svRaw := configNrql["since_value"]

	// Handle deprecated
	if svRaw != nil && svRaw.(string) != "" && nrql.EvaluationOffset != nil {
		evalOffset := nrql.EvaluationOffset
		out["since_value"] = strconv.Itoa(*evalOffset)
	} else {
		out["evaluation_offset"] = nrql.EvaluationOffset
	}

	return []interface{}{out}
}

// NerdGraph response => Terraform state
func flattenNrqlTerms(terms []alerts.NrqlConditionTerm, configTerms []interface{}) []map[string]interface{} {
	// Represents the built terms to be saved in state
	var out []map[string]interface{}

	// Import scenario
	if len(terms) > 0 && len(configTerms) == 0 {
		return handleImportFlattenNrqlTerms(terms)
	}

	// Represents the terms set in the user's .tf config file
	configuredTerms := getConfiguredTerms(configTerms)

	for i, term := range terms {
		dst := map[string]interface{}{
			"operator":  strings.ToLower(string(term.Operator)),
			"priority":  strings.ToLower(string(term.Priority)),
			"threshold": term.Threshold,
		}

		if term.Prediction != nil {
			dst["prediction"] = make([]interface{}, 1)
			predictionBlockContents := make(map[string]interface{})
			predictionBlockContents["predict_by"] = term.Prediction.PredictBy
			predictionBlockContents["prefer_prediction_violation"] = term.Prediction.PreferPredictionViolation
			dst["prediction"] = []interface{}{predictionBlockContents}
		}

		if i < len(configuredTerms) {
			setDuration := configuredTerms[i]["duration"]
			if setDuration != nil && setDuration.(int) > 0 {
				dst["duration"] = term.ThresholdDuration / 60 // convert to minutes for old way
			} else {
				dst["threshold_duration"] = term.ThresholdDuration
			}
		} else {
			dst["threshold_duration"] = term.ThresholdDuration
		}

		if i < len(configuredTerms) {
			setTimeFunction := configuredTerms[i]["time_function"]
			if setTimeFunction != nil && setTimeFunction.(string) != "" {
				dst["time_function"] = timeFunctionMapNewOld[term.ThresholdOccurrences]
			} else {
				dst["threshold_occurrences"] = strings.ToLower(string(term.ThresholdOccurrences))
			}
		} else {
			dst["threshold_occurrences"] = strings.ToLower(string(term.ThresholdOccurrences))
		}

		if term.DisableHealthStatusReporting != nil {
			dst["disable_health_status_reporting"] = term.DisableHealthStatusReporting
		}

		out = append(out, dst)
	}

	return out
}

// Note: We DO NOT set deprecated attributes on import for NRQL alert conditions.
func handleImportFlattenNrqlTerms(terms []alerts.NrqlConditionTerm) []map[string]interface{} {
	var out []map[string]interface{}

	for _, term := range terms {
		dst := map[string]interface{}{
			"operator":              strings.ToLower(string(term.Operator)),
			"priority":              strings.ToLower(string(term.Priority)),
			"threshold":             term.Threshold,
			"threshold_duration":    term.ThresholdDuration,
			"threshold_occurrences": strings.ToLower(string(term.ThresholdOccurrences)),
		}

		if term.Prediction != nil {
			dst["prediction"] = make([]interface{}, 1)
			predictionBlockContents := make(map[string]interface{})
			predictionBlockContents["predict_by"] = term.Prediction.PredictBy
			predictionBlockContents["prefer_prediction_violation"] = term.Prediction.PreferPredictionViolation
			dst["prediction"] = []interface{}{predictionBlockContents}
		}

		if term.DisableHealthStatusReporting != nil {
			dst["disable_health_status_reporting"] = term.DisableHealthStatusReporting
		}

		out = append(out, dst)
	}

	return out
}

// Returns the term attributes that were configured by the user in their .tf config file
func getConfiguredTerms(configTerms []interface{}) []map[string]interface{} {
	var setTerms []map[string]interface{}

	for _, tm := range configTerms {
		t := tm.(map[string]interface{})
		trm := map[string]interface{}{
			"operator":      t["operator"],
			"priority":      t["priority"],
			"threshold":     t["threshold"],
			"duration":      t["duration"],
			"time_function": t["time_function"],

			// NerdGraph fields
			"threshold_duration":              t["threshold_duration"],
			"threshold_occurrences":           t["threshold_occurrences"],
			"disable_health_status_reporting": t["disable_health_status_reporting"],
		}

		predictionSet := t["prediction"].(*schema.Set)
		if predictionSet.Len() > 0 {
			predictionMap := predictionSet.List()[0].(map[string]interface{})
			predictionMapOut := map[string]interface{}{}
			if predictBy, ok := predictionMap["predict_by"].(int); ok {
				predictionMapOut["predict_by"] = predictBy
			}
			if preferPredictionViolation, ok := predictionMap["prefer_prediction_violation"].(bool); ok {
				predictionMapOut["prefer_prediction_violation"] = preferPredictionViolation
			}
			trm["prediction"] = predictionMapOut
		}

		setTerms = append(setTerms, trm)
	}

	return setTerms
}

func validateNrqlConditionAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []error

	_, conditionType := d.GetChange("type")
	if conditionType != nil {
		isNotBaselineCondition := !strings.Contains(conditionType.(string), "baseline")
		if isNotBaselineCondition {
			err := validateSignalSeasonality(d)
			if err != nil {
				errorsList = append(errorsList, err)
			}
		}
	}

	if len(errorsList) == 0 {
		return nil
	}

	errorsString := "the following validation errors have been identified with the configuration of the nrql alert condition: \n"

	for index, val := range errorsList {
		errorsString += fmt.Sprintf("(%d): %s\n", index+1, val)
	}

	return errors.New(errorsString)
}

func validateSignalSeasonality(d *schema.ResourceDiff) error {
	rawConfiguration := d.GetRawConfig()

	signalSeasonalityIsNotNil := !rawConfiguration.GetAttr("signal_seasonality").IsNull()

	if signalSeasonalityIsNotNil {
		return fmt.Errorf(`'signal_seasonality' is only valid on baseline conditions. Please remove this field or change the condition type`)
	}
	return nil
}
