package newrelic

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandPluginsCondition(d *schema.ResourceData) *alerts.PluginsCondition {
	condition := alerts.PluginsCondition{
		Name:              d.Get("name").(string),
		Enabled:           d.Get("enabled").(bool),
		Metric:            d.Get("metric").(string),
		MetricDescription: d.Get("metric_description").(string),
		ValueFunction:     d.Get("value_function").(string),
		PolicyID:          d.Get("policy_id").(int),
	}

	condition.Entities = expandPluginsConditionEntities(d.Get("entities").(*schema.Set).List())
	condition.Terms = expandPluginsConditionTerms(d.Get("term").(*schema.Set).List())

	condition.Plugin = alerts.AlertPlugin{
		ID:   d.Get("plugin_id").(string),
		GUID: d.Get("plugin_guid").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	return &condition
}

func expandPluginsConditionEntities(entities []interface{}) []string {
	perms := make([]string, len(entities))

	for i, entity := range entities {
		perms[i] = strconv.Itoa(entity.(int))
	}

	return perms
}

func expandPluginsConditionTerms(terms []interface{}) []alerts.ConditionTerm {
	perms := make([]alerts.ConditionTerm, len(terms))

	for i, t := range terms {
		term := t.(map[string]interface{})

		perms[i] = alerts.ConditionTerm{
			Duration:     term["duration"].(int),
			Operator:     term["operator"].(string),
			Priority:     term["priority"].(string),
			Threshold:    term["threshold"].(float64),
			TimeFunction: term["time_function"].(alerts.TimeFunctionType),
		}
	}

	return perms
}

func flattenPluginsCondition(condition *alerts.PluginsCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("enabled", condition.Enabled)
	d.Set("metric", condition.Metric)
	d.Set("metric_description", condition.MetricDescription)
	d.Set("value_function", condition.ValueFunction)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("plugin_id", condition.Plugin.ID)
	d.Set("plugin_guid", condition.Plugin.GUID)

	entities, err := flattenPluginsConditionEntities(condition.Entities)
	if err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition entities: %#v", err)
	}

	d.Set("entities", entities)

	if err := d.Set("term", flattenPluginsConditionTerms(condition.Terms)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting alert condition terms: %#v", err)
	}

	return nil
}

func flattenPluginsConditionEntities(in []string) ([]int, error) {
	entities := make([]int, len(in))
	for i, entity := range in {
		v, err := strconv.ParseInt(entity, 10, 32)
		if err != nil {
			return nil, err
		}
		entities[i] = int(v)
	}

	return entities, nil
}

func flattenPluginsConditionTerms(in []alerts.ConditionTerm) []map[string]interface{} {
	var terms []map[string]interface{}

	for _, src := range in {
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
