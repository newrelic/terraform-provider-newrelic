package newrelic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
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
	violationTimeLimitMap = map[int]alerts.NrqlConditionViolationTimeLimit{
		3600:  alerts.NrqlConditionViolationTimeLimits.OneHour,
		7200:  alerts.NrqlConditionViolationTimeLimits.TwoHours,
		14400: alerts.NrqlConditionViolationTimeLimits.FourHours,
		28800: alerts.NrqlConditionViolationTimeLimits.EightHours,
		43200: alerts.NrqlConditionViolationTimeLimits.TwelveHours,
		86400: alerts.NrqlConditionViolationTimeLimits.TwentyFourHours,
	}

	// new:old
	violationTimeLimitMapNewOld = map[alerts.NrqlConditionViolationTimeLimit]int{
		alerts.NrqlConditionViolationTimeLimits.OneHour:         3600,
		alerts.NrqlConditionViolationTimeLimits.TwoHours:        7200,
		alerts.NrqlConditionViolationTimeLimits.FourHours:       1440,
		alerts.NrqlConditionViolationTimeLimits.EightHours:      2880,
		alerts.NrqlConditionViolationTimeLimits.TwelveHours:     4320,
		alerts.NrqlConditionViolationTimeLimits.TwentyFourHours: 8640,
	}
)

// NerdGraph
func expandNrqlAlertConditionInput(d *schema.ResourceData) (*alerts.NrqlConditionInput, error) {
	input := alerts.NrqlConditionInput{
		NrqlConditionBase: alerts.NrqlConditionBase{
			Description:        d.Get("description").(string),
			Enabled:            d.Get("enabled").(bool),
			Name:               d.Get("name").(string),
			ViolationTimeLimit: alerts.NrqlConditionViolationTimeLimit(strings.ToUpper(d.Get("violation_time_limit").(string))),
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

	if conditionType == "static" {
		if attr, ok := d.GetOk("value_function"); ok {
			valFn := alerts.NrqlConditionValueFunction(strings.ToUpper(attr.(string)))
			input.ValueFunction = &valFn
		} else {
			return nil, fmt.Errorf("attribute `%s` is required for nrql alert conditions of type `%+v`", "value_function", conditionType)
		}
	}

	if conditionType == "outlier" {
		defaultExpectedGroups := 1
		if expectedGroups, ok := d.GetOk("expected_groups"); ok {
			expectedGroupsValue := expectedGroups.(int)
			input.ExpectedGroups = &expectedGroupsValue
		} else {
			input.ExpectedGroups = &defaultExpectedGroups
		}

		var openViolationOnOverlap bool
		if ignoreOverlap, ok := d.GetOkExists("ignore_overlap"); ok {
			// Note: ignore_overlap is the inverse of open_violation_on_group_overlap
			openViolationOnOverlap = !ignoreOverlap.(bool)
		} else if violationOnOverlap, ok := d.GetOkExists("open_violation_on_group_overlap"); ok {
			openViolationOnOverlap = violationOnOverlap.(bool)
		}

		input.OpenViolationOnGroupOverlap = &openViolationOnOverlap
	}

	if runbookURL, ok := d.GetOk("runbook_url"); ok {
		input.RunbookURL = runbookURL.(string)
	}

	if violationTimeLimitSec, ok := d.GetOk("violation_time_limit_seconds"); ok {
		input.ViolationTimeLimit = violationTimeLimitMap[violationTimeLimitSec.(int)]
	} else if violationTimeLimit, ok := d.GetOk("violation_time_limit"); ok {
		input.ViolationTimeLimit = alerts.NrqlConditionViolationTimeLimit(strings.ToUpper(violationTimeLimit.(string)))
	}

	nrql, err := expandNrql(d, input)
	if err != nil {
		return nil, err
	}

	input.Nrql = *nrql

	terms, err := expandNrqlTerms(d.Get("term").(*schema.Set).List(), conditionType)
	if err != nil {
		return nil, err
	}

	input.Terms = terms

	return &input, nil
}

// NerdGraph
func expandNrql(d *schema.ResourceData, condition alerts.NrqlConditionInput) (*alerts.NrqlConditionQuery, error) {
	var nrql alerts.NrqlConditionQuery

	if nrqlQuery, ok := d.GetOk("nrql.0.query"); ok {
		nrql.Query = nrqlQuery.(string)
	}

	if sinceValue, ok := d.GetOk("nrql.0.since_value"); ok {
		sv, err := strconv.Atoi(sinceValue.(string))
		if err != nil {
			return nil, err
		}

		nrql.EvaluationOffset = sv
	} else if evalOffset, ok := d.GetOk("nrql.0.evaluation_offset"); ok {
		nrql.EvaluationOffset = evalOffset.(int)
	} else {
		return nil, fmt.Errorf("one of `since_value` or `evaluation_offset` must be configured for block `nrql`")
	}

	return &nrql, nil
}

// NerdGraph
func expandNrqlTerms(terms []interface{}, conditionType string) ([]alerts.NrqlConditionTerm, error) {
	expanded := make([]alerts.NrqlConditionTerm, len(terms))

	for i, t := range terms {
		term := t.(map[string]interface{})
		durationIn := term["duration"].(int)
		thresholdDurationIn := term["threshold_duration"].(int)

		if durationIn == 0 && thresholdDurationIn == 0 {
			return nil, fmt.Errorf("one of `duration` or `threshold_duration` must be configured for block `term`")
		}

		if durationIn > 0 && thresholdDurationIn > 0 {
			return nil, fmt.Errorf("one of `duration` or `threshold_duration` must be configured for block `term`, but not both")
		}

		var duration int
		if durationIn > 0 {
			duration = durationIn * 60 // convert min to sec
		} else {
			duration = thresholdDurationIn
		}

		threshold := term["threshold"].(float64)

		if conditionType == "baseline" {
			if duration < 120 || duration > 3600 {
				return nil, fmt.Errorf("for baseline conditions duration must be in range %v, got %v", "[2, 60]", duration)
			}

			if threshold < 1 || threshold > 1000 {
				return nil, fmt.Errorf("for baseline conditions threshold must be in range %v, got %v", "[1, 1000]", threshold)
			}
		}

		timeFunctionIn := term["time_function"].(string)
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

		expanded[i] = alerts.NrqlConditionTerm{
			Operator:             alerts.NrqlConditionOperator(strings.ToUpper(term["operator"].(string))),
			Priority:             alerts.NrqlConditionPriority(strings.ToUpper(term["priority"].(string))),
			Threshold:            threshold,
			ThresholdDuration:    duration,
			ThresholdOccurrences: thresholdOccurrences,
		}
	}

	return expanded, nil
}

// NerdGraph
func flattenNrqlAlertCondition(accountID int, condition *alerts.NrqlAlertCondition, d *schema.ResourceData) error {
	policyID, err := strconv.Atoi(condition.PolicyID)
	if err != nil {
		return err
	}

	conditionType := strings.ToLower(string(condition.Type))

	d.Set("account_id", accountID)
	d.Set("type", conditionType)
	d.Set("description", condition.Description)
	d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("enabled", condition.Enabled)

	if conditionType == "baseline" {
		d.Set("baseline_direction", string(*condition.BaselineDirection))
	}

	if conditionType == "static" {
		d.Set("value_function", string(*condition.ValueFunction))
	}

	if conditionType == "outlier" {
		d.Set("expected_groups", *condition.ExpectedGroups)

		openViolationOnGroupOverlap := *condition.OpenViolationOnGroupOverlap
		if _, ok := d.GetOkExists("ignore_overlap"); ok {
			d.Set("ignore_overlap", !openViolationOnGroupOverlap)
		} else {
			d.Set("open_violation_on_group_overlap", openViolationOnGroupOverlap)
		}
	}

	configuredNrql := d.Get("nrql.0").(map[string]interface{})
	if err := d.Set("nrql", flattenNrql(condition.Nrql, configuredNrql)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `nrql`: %v", err)
	}

	configuredTerms := d.Get("term").(*schema.Set).List()
	if err := d.Set("term", flattenNrqlTerms(condition.Terms, configuredTerms)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting nrql alert condition `term`: %v", err)
	}

	if _, ok := d.GetOk("violation_time_limit_seconds"); ok {
		d.Set("violation_time_limit_seconds", violationTimeLimitMapNewOld[condition.ViolationTimeLimit])
	} else {
		d.Set("violation_time_limit", condition.ViolationTimeLimit)
	}

	return nil
}

// NerdGraph
func flattenNrql(nrql alerts.NrqlConditionQuery, configNrql map[string]interface{}) []interface{} {
	out := map[string]interface{}{
		"query": nrql.Query,
	}

	svRaw := configNrql["since_value"]

	// Handle deprecated
	if svRaw != nil && svRaw.(string) != "" {
		out["since_value"] = strconv.Itoa(nrql.EvaluationOffset)
	} else {
		out["evaluation_offset"] = nrql.EvaluationOffset
	}

	return []interface{}{out}
}

// NerdGraph
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

		setDuration := configuredTerms[i]["duration"]
		if setDuration != nil && setDuration.(int) > 0 {
			dst["duration"] = term.ThresholdDuration / 60 // convert to minutes for old way
		} else {
			dst["threshold_duration"] = term.ThresholdDuration
		}

		setTimeFunction := configuredTerms[i]["time_function"]
		if setTimeFunction != nil && setTimeFunction.(string) != "" {
			dst["time_function"] = timeFunctionMapNewOld[term.ThresholdOccurrences]
		} else {
			dst["threshold_occurrences"] = term.ThresholdOccurrences
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
			"threshold_occurrences": term.ThresholdOccurrences,
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
			"threshold_duration":    t["threshold_duration"],
			"threshold_occurrences": t["threshold_occurrences"],
		}

		setTerms = append(setTerms, trm)
	}

	return setTerms
}
