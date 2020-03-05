package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandNrqlAlertConditionStruct(d *schema.ResourceData) *alerts.NrqlCondition {
	condition := alerts.NrqlCondition{
		Name:                d.Get("name").(string),
		Type:                d.Get("type").(string),
		Enabled:             d.Get("enabled").(bool),
		PolicyID:            d.Get("policy_id").(int),
		ValueFunction:       d.Get("value_function").(string),
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
	trms := make([]alerts.ConditionTerm, len(terms))

	for i, t := range terms {
		term := t.(map[string]interface{})

		trms[i] = alerts.ConditionTerm{
			Duration:     term["duration"].(int),
			Operator:     term["operator"].(string),
			Priority:     term["priority"].(string),
			Threshold:    term["threshold"].(float64),
			TimeFunction: term["time_function"].(alerts.TimeFunctionType),
		}
	}

	return trms
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
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	d.Set("policy_id", policyID)
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
