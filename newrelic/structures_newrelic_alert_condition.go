package newrelic

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandAlertCondition(d *schema.ResourceData) *alerts.Condition {
	condition := alerts.Condition{
		Type:                d.Get("type").(string),
		Name:                d.Get("name").(string),
		Enabled:             d.Get("enabled").(bool),
		Metric:              d.Get("metric").(string),
		PolicyID:            d.Get("policy_id").(int),
		Scope:               d.Get("condition_scope").(string),
		ViolationCloseTimer: d.Get("violation_close_timer").(int),
		GCMetric:            d.Get("gc_metric").(string),
	}

	condition.Entities = expandAlertConditionEntities(d.Get("entities").(*schema.Set).List())
	condition.Terms = expandAlertConditionTerms(d.Get("term").(*schema.Set).List())

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	if attrM, ok := d.GetOk("user_defined_metric"); ok {
		if attrVF, ok := d.GetOk("user_defined_value_function"); ok {
			condition.UserDefined = alerts.ConditionUserDefined{
				Metric:        attrM.(string),
				ValueFunction: attrVF.(string),
			}
		}
	}

	return &condition
}

func expandAlertConditionEntities(entities []interface{}) []string {
	perms := make([]string, len(entities))

	for i, entity := range entities {
		perms[i] = strconv.Itoa(entity.(int))
	}

	return perms
}

func expandAlertConditionTerms(terms []interface{}) []alerts.ConditionTerm {
	perms := make([]alerts.ConditionTerm, len(terms))

	for i, term := range terms {
		term := term.(map[string]interface{})

		perms[i] = alerts.ConditionTerm{
			Duration:     term["duration"].(int),
			Operator:     term["operator"].(string),
			Priority:     term["priority"].(string),
			Threshold:    term["threshold"].(float64),
			TimeFunction: term["time_function"].(string),
		}
	}

	return perms
}

func flattenAlertCondition(condition *alerts.Condition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("enabled", condition.Enabled)
	d.Set("type", condition.Type)
	d.Set("metric", condition.Metric)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("condition_scope", condition.Scope)
	d.Set("violation_close_timer", condition.ViolationCloseTimer)
	d.Set("gc_metric", condition.GCMetric)
	d.Set("user_defined_metric", condition.UserDefined.Metric)
	d.Set("user_defined_value_function", condition.UserDefined.ValueFunction)

	entities, err := flattenAlertConditionEntities(&condition.Entities)

	if err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition entities: %#v", err)
	}

	d.Set("entities", entities)

	if err := d.Set("term", flattenAlertConditionTerms(&condition.Terms)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition terms: %#v", err)
	}

	return nil
}

func flattenAlertConditionEntities(in *[]string) ([]int, error) {
	entities := make([]int, len(*in))
	for i, entity := range *in {
		v, err := strconv.ParseInt(entity, 10, 32)
		if err != nil {
			return nil, err
		}
		entities[i] = int(v)
	}

	return entities, nil
}

func flattenAlertConditionTerms(in *[]alerts.ConditionTerm) []map[string]interface{} {
	var terms []map[string]interface{}

	for _, src := range *in {
		dst := map[string]interface{}{
			"duration":      src.Duration,
			"operator":      src.Operator,
			"priority":      src.Priority,
			"threshold":     src.Threshold,
			"time_function": src.TimeFunction,
		}
		terms = append(terms, dst)
	}

	return terms
}
