package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandMutingRuleCreateInput(d *schema.ResourceData) alerts.MutingRuleCreateInput {
	createInput := alerts.MutingRuleCreateInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("condition"); ok {
		createInput.Condition = expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))
	}

	return createInput
}

func expandMutingRuleUpdateInput(d *schema.ResourceData) alerts.MutingRuleUpdateInput {
	updateInput := alerts.MutingRuleUpdateInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("condition"); ok {
		x := expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))

		updateInput.Condition = &x
	}

	return updateInput
}

func expandMutingRuleConditionGroup(cfg map[string]interface{}) alerts.MutingRuleConditionGroup {
	conditionGroup := alerts.MutingRuleConditionGroup{}
	var expandedConditions []alerts.MutingRuleCondition

	conditions := cfg["conditions"].([]interface{})

	for _, c := range conditions {
		var y = expandMutingRuleCondition(c)
		expandedConditions = append(expandedConditions, y)
	}

	conditionGroup.Conditions = expandedConditions

	if operator, ok := cfg["operator"]; ok {
		conditionGroup.Operator = operator.(string)
	}

	return conditionGroup
}

func expandMutingRuleCondition(cfg interface{}) alerts.MutingRuleCondition {
	conditionCfg := cfg.(map[string]interface{})
	condition := alerts.MutingRuleCondition{}

	if attribute, ok := conditionCfg["attribute"]; ok {
		condition.Attribute = attribute.(string)
	}

	if operator, ok := conditionCfg["operator"]; ok {
		condition.Operator = operator.(string)
	}
	if values, ok := conditionCfg["values"]; ok {
		condition.Values = expandMutingRuleValues(values.([]interface{}))
	}

	return condition
}

func expandMutingRuleValues(values []interface{}) []string {
	perms := make([]string, len(values))

	for i, values := range values {
		perms[i] = values.(string)
	}

	return perms
}

func flattenMutingRule(mutingRule *alerts.MutingRule, d *schema.ResourceData) error {
	d.Set("enabled", mutingRule.Enabled)
	d.Set("condition", flattenMutingRuleConditionGroup(mutingRule.Condition))
	d.Set("description", mutingRule.Description)
	d.Set("name", mutingRule.Name)

	return nil
}

func flattenMutingRuleConditionGroup(in alerts.MutingRuleConditionGroup) []map[string]interface{} {
	condition := []map[string]interface{}{
		{
			"operator": in.Operator,
		},
	}

	for _, src := range in.Conditions {
		dst := map[string]interface{}{
			"conditions": flattenMutingRuleCondition(src),
		}
		condition = append(condition, dst)
	}

	return condition
}

func flattenMutingRuleCondition(in alerts.MutingRuleCondition) map[string]interface{} {
	return map[string]interface{}{
		"attribute": in.Attribute,
		"operator":  in.Operator,
		"values":    in.Values,
	}
}
