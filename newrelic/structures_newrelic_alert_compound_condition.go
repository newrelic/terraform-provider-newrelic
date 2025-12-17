package newrelic

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

// expandAlertCompoundConditionCreateInput builds the create input from Terraform schema
func expandAlertCompoundConditionCreateInput(d *schema.ResourceData) (*alerts.CompoundConditionCreateInput, error) {
	input := alerts.CompoundConditionCreateInput{
		Name:              d.Get("name").(string),
		Enabled:           d.Get("enabled").(bool),
		TriggerExpression: d.Get("trigger_expression").(string),
	}

	// Expand component conditions
	componentConditions, err := expandComponentConditions(d.Get("component_conditions").(*schema.Set))
	if err != nil {
		return nil, err
	}
	input.ComponentConditions = componentConditions

	// Optional fields
	if v, ok := d.GetOk("facet_matching_behavior"); ok {
		val := v.(string)
		input.FacetMatchingBehavior = &val
	}

	if v, ok := d.GetOk("runbook_url"); ok {
		val := v.(string)
		input.RunbookURL = &val
	}

	if v, ok := d.GetOk("threshold_duration"); ok {
		val := v.(int)
		input.ThresholdDuration = &val
	}

	return &input, nil
}

// expandAlertCompoundConditionUpdateInput builds the update input from Terraform schema
func expandAlertCompoundConditionUpdateInput(d *schema.ResourceData) (*alerts.CompoundConditionUpdateInput, error) {
	policyID := strconv.Itoa(d.Get("policy_id").(int))

	input := alerts.CompoundConditionUpdateInput{
		Name:              d.Get("name").(string),
		Enabled:           d.Get("enabled").(bool),
		TriggerExpression: d.Get("trigger_expression").(string),
		PolicyID:          &policyID,
	}

	// Expand component conditions
	componentConditions, err := expandComponentConditions(d.Get("component_conditions").(*schema.Set))
	if err != nil {
		return nil, err
	}
	input.ComponentConditions = componentConditions

	// Optional fields
	if v, ok := d.GetOk("facet_matching_behavior"); ok {
		val := v.(string)
		input.FacetMatchingBehavior = &val
	}

	if v, ok := d.GetOk("runbook_url"); ok {
		val := v.(string)
		input.RunbookURL = &val
	}

	if v, ok := d.GetOk("threshold_duration"); ok {
		val := v.(int)
		input.ThresholdDuration = &val
	}

	return &input, nil
}

// expandComponentConditions converts the Terraform set to ComponentConditionInput slice
func expandComponentConditions(componentSet *schema.Set) ([]alerts.ComponentConditionInput, error) {
	components := make([]alerts.ComponentConditionInput, 0, componentSet.Len())
	aliases := make(map[string]bool)

	for _, c := range componentSet.List() {
		component := c.(map[string]interface{})

		id := component["id"].(string)
		alias := component["alias"].(string)

		// Validate unique aliases
		if aliases[alias] {
			return nil, fmt.Errorf("duplicate alias '%s' found in component_conditions", alias)
		}
		aliases[alias] = true

		components = append(components, alerts.ComponentConditionInput{
			ID:    id,
			Alias: alias,
		})
	}

	return components, nil
}

// flattenAlertCompoundCondition converts API response to Terraform state
func flattenAlertCompoundCondition(accountID int, condition *alerts.CompoundCondition, d *schema.ResourceData) error {
	policyID, err := strconv.Atoi(condition.PolicyID)
	if err != nil {
		return fmt.Errorf("error converting policy ID to int: %v", err)
	}

	_ = d.Set("account_id", accountID)
	_ = d.Set("policy_id", policyID)
	_ = d.Set("name", condition.Name)
	_ = d.Set("enabled", condition.Enabled)
	_ = d.Set("trigger_expression", condition.TriggerExpression)
	_ = d.Set("facet_matching_behavior", condition.FacetMatchingBehavior)
	_ = d.Set("runbook_url", condition.RunbookURL)
	_ = d.Set("threshold_duration", condition.ThresholdDuration)

	// Flatten component conditions - ONLY id and alias (per user requirement)
	componentConditions := flattenComponentConditions(condition.ComponentConditions)
	if err := d.Set("component_conditions", componentConditions); err != nil {
		return fmt.Errorf("error setting component_conditions: %v", err)
	}

	// Note: entity_guid might be available in the API response if needed
	// For now, we're leaving it unset as it's not at the top level of AlertCompoundCondition

	return nil
}

// flattenComponentConditions converts API component conditions to Terraform state
// Only includes ID and Alias - ignores nested NRQL condition details per user requirement
func flattenComponentConditions(components []alerts.ComponentCondition) *schema.Set {
	result := make([]interface{}, 0, len(components))

	for _, component := range components {
		c := map[string]interface{}{
			"id":    component.ID,
			"alias": component.Alias,
			// Explicitly NOT including component.Condition (nested NRQL data)
		}
		result = append(result, c)
	}

	// Use the same hash function as the schema
	return schema.NewSet(func(v interface{}) int {
		m := v.(map[string]interface{})
		return schema.HashString(m["alias"].(string))
	}, result)
}
