package newrelic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandNrqlAlertConditionInput(d *schema.ResourceData) (*alerts.NrqlConditionInput, error) {
	conditionType := d.Get("type").(string)

	input := alerts.NrqlConditionInput{
		NrqlConditionBase: alerts.NrqlConditionBase{
			Description:        d.Get("description").(string),
			Enabled:            d.Get("enabled").(bool),
			Name:               d.Get("name").(string),
			ViolationTimeLimit: alerts.NrqlConditionViolationTimeLimit(strings.ToUpper(d.Get("violation_time_limit").(string))),
		},
	}

	if conditionType == "baseline" {
		if attr, ok := d.GetOk("baseline_direction"); ok {
			d := alerts.NrqlBaselineDirection(attr.(string))
			input.BaselineDirection = &d
		} else {
			return nil, fmt.Errorf("attribute `%s` is required for nrql alert conditions of type `%+v`", attr, conditionType)
		}
	}

	if conditionType == "static" {
		if attr, ok := d.GetOk("value_function"); ok {
			d := alerts.NrqlConditionValueFunction(attr.(string))
			input.ValueFunction = &d
		} else {
			return nil, fmt.Errorf("attribute `%s` is required for nrql alert conditions of type `%+v`", attr, conditionType)
		}
	}

	if runbookURL, ok := d.GetOk("runbook_url"); ok {
		input.RunbookURL = runbookURL.(string)
	}

	if violationTimeLimit, ok := d.GetOkExists("violation_time_limit"); ok {
		input.ViolationTimeLimit = alerts.NrqlConditionViolationTimeLimit(violationTimeLimit.(string))
	}

	input.Nrql = expandNrql(d, input)
	input.Terms = expandNrqlTerms(d.Get("term").(*schema.Set).List())

	return &input, nil
}

func expandNrql(d *schema.ResourceData, condition alerts.NrqlConditionInput) alerts.NrqlConditionQuery {
	var nrql alerts.NrqlConditionQuery

	if nrqlQuery, ok := d.GetOk("nrql.0.query"); ok {
		nrql.Query = nrqlQuery.(string)
	}

	if evaluationOffset, ok := d.GetOk("nrql.0.evaluation_offset"); ok {
		nrql.EvaluationOffset = evaluationOffset.(int)
	}

	return nrql
}

func expandNrqlTerms(terms []interface{}) []alerts.NrqlConditionTerms {
	expanded := make([]alerts.NrqlConditionTerms, len(terms))

	for i, t := range terms {
		term := t.(map[string]interface{})

		expanded[i] = alerts.NrqlConditionTerms{
			Operator:             alerts.NrqlConditionOperator(strings.ToUpper(term["operator"].(string))),
			Priority:             alerts.NrqlConditionPriority(strings.ToUpper(term["priority"].(string))),
			Threshold:            term["threshold"].(float64),
			ThresholdDuration:    float64(term["duration"].(int) * 60),
			ThresholdOccurrences: alerts.ThresholdOccurence(strings.ToUpper(term["threshold_occurrences"].(string))),
		}
	}

	return expanded
}

func expandNrqlAlertConditionStruct(d *schema.ResourceData) *alerts.NrqlCondition {
	condition := alerts.NrqlCondition{
		Name:                d.Get("name").(string),
		Type:                d.Get("type").(string),
		Enabled:             d.Get("enabled").(bool),
		ValueFunction:       alerts.ValueFunctionType(d.Get("value_function").(string)),
		ViolationCloseTimer: d.Get("violation_time_limit_seconds").(int),
	}

	condition.Terms = expandNrqlConditionTerms(d.Get("term").(*schema.Set).List())

	if nrqlQuery, ok := d.GetOk("nrql.0.query"); ok {
		condition.Nrql.Query = nrqlQuery.(string)
	}

	if sinceValue, ok := d.GetOk("nrql.0.since_value"); ok {
		condition.Nrql.SinceValue = sinceValue.(string)
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	if attr, ok := d.GetOkExists("ignore_overlap"); ok {
		condition.IgnoreOverlap = attr.(bool)
	}

	if attr, ok := d.GetOkExists("violation_time_limit_seconds"); ok {
		condition.ViolationCloseTimer = attr.(int)
	}

	if attr, ok := d.GetOk("expected_groups"); ok {
		condition.ExpectedGroups = attr.(int)
	}

	return &condition
}

func expandNrqlConditionTerms(terms []interface{}) []alerts.ConditionTerm {
	expanded := make([]alerts.ConditionTerm, len(terms))

	for i, t := range terms {
		term := t.(map[string]interface{})

		expanded[i] = alerts.ConditionTerm{
			Duration:     term["duration"].(int),
			Operator:     alerts.OperatorType(term["operator"].(string)),
			Priority:     alerts.PriorityType(term["priority"].(string)),
			Threshold:    term["threshold"].(float64),
			TimeFunction: alerts.TimeFunctionType(term["time_function"].(string)),
		}
	}

	return expanded
}

func flattenNrqlQuery(nrql alerts.NrqlQuery) []interface{} {
	m := map[string]interface{}{
		"query":       nrql.Query,
		"since_value": nrql.SinceValue,
	}

	return []interface{}{m}
}

func flattenNrqlConditionTerms(terms []alerts.ConditionTerm) []map[string]interface{} {
	var t []map[string]interface{}

	for _, src := range terms {
		dst := map[string]interface{}{
			"duration":      src.Duration,
			"operator":      src.Operator,
			"priority":      src.Priority,
			"threshold":     src.Threshold,
			"time_function": src.TimeFunction,
		}
		t = append(t, dst)
	}

	return t
}

func flattenNrqlConditionStruct(condition *alerts.NrqlCondition, d *schema.ResourceData) error {
	d.Set("name", condition.Name)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("enabled", condition.Enabled)
	d.Set("type", condition.Type)
	d.Set("expected_groups", condition.ExpectedGroups)
	d.Set("ignore_overlap", condition.IgnoreOverlap)
	d.Set("violation_time_limit_seconds", condition.ViolationCloseTimer)

	if condition.ValueFunction == "" {
		d.Set("value_function", "single_value")
	} else {
		d.Set("value_function", condition.ValueFunction)
	}

	if err := d.Set("nrql", flattenNrqlQuery(condition.Nrql)); err != nil {
		return err
	}

	terms := flattenNrqlConditionTerms(condition.Terms)

	if err := d.Set("term", terms); err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition terms: %#v", err)
	}

	return nil
}
