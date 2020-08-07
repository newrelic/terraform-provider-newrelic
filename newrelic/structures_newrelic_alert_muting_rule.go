package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandAlertMutingRule(d *schema.ResourceData) (*alerts.MutingRule, error) {
	mutingRule := alerts.Condition{
		Type:      alerts.MutingRule(d.Get("type").(string)),
		AccountId: d.Get("accountId").(int),
		// Condition: d.Get("condition").()
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// condition.Entities = expandAlertConditionEntities(d.Get("entities").(*schema.Set).List())
	// condition.Terms = expandAlertConditionTerms(d.Get("term").(*schema.Set).List())

	// if violationCloseTimer, ok := d.GetOk("violation_close_timer"); ok {
	// 	if condition.Type == "apm_app_metric" && condition.Scope == "application" {
	// 		return nil, fmt.Errorf("violation_close_timer only supported for apm_app_metric when condition_scope = 'instance'")
	// 	}

	// 	condition.ViolationCloseTimer = violationCloseTimer.(int)
	// }

	// if attr, ok := d.GetOk("runbook_url"); ok {
	// 	condition.RunbookURL = attr.(string)
	// }

	// if attr, ok := d.GetOk("user_defined_metric"); ok {
	// 	condition.UserDefined.Metric = attr.(string)
	// }

	// if attr, ok := d.GetOk("user_defined_value_function"); ok {
	// 	condition.UserDefined.ValueFunction = alerts.ValueFunctionType(attr.(string))
	// }

	return &condition, nil
}

func flattenAlertMutingRule(mutingRule *alerts.Alerts, d *schema.ResourceData) error {
	d.Set("account_id", mutingRule.Account.ID)
	d.Set("rule", mutingRule.Rule)
	d.Set("condition", mutingRule.Condition) // Need to fix
	d.Set("createdAt", mutingRule.CreatedAt)
	d.Set("createdBy", mutingRule.CreatedBy)
	d.Set("description", mutingRule.Description)

	d.Set("id", mutingRule.ID)
	d.Set("name", mutingRule.Name)
	d.Set("updatedAt", mutingRule.UpdatedAt)
	d.Set("updatedBy", mutingRule.UpdatedBy)

	return nil
}
