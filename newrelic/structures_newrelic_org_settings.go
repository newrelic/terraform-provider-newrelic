package newrelic

import "github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"

// flattenSyncGroupRules converts API SyncGroupRule values to the Terraform list
// shape used by the sync_group_rules attribute.
func flattenSyncGroupRules(rules []scorecards.EntityManagementSyncGroupRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		conds := make([]map[string]interface{}, 0, len(rule.Conditions))
		for _, c := range rule.Conditions {
			conds = append(conds, map[string]interface{}{
				"type":  string(c.Type),
				"value": c.Value,
			})
		}
		out = append(out, map[string]interface{}{"conditions": conds})
	}
	return out
}

// expandSyncGroupsUpdate converts the Terraform sync_group_rules list and
// sync_groups_enabled bool into the API update input type.
func expandSyncGroupsUpdate(enabled bool, rawRules []interface{}) scorecards.EntityManagementSyncGroupsSettingsUpdateInput {
	rules := make([]scorecards.EntityManagementSyncGroupRuleUpdateInput, 0, len(rawRules))
	for _, rr := range rawRules {
		rm := rr.(map[string]interface{})
		rawConds := rm["conditions"].([]interface{})
		conds := make([]scorecards.EntityManagementSyncGroupRuleConditionUpdateInput, 0, len(rawConds))
		for _, rc := range rawConds {
			cm := rc.(map[string]interface{})
			conds = append(conds, scorecards.EntityManagementSyncGroupRuleConditionUpdateInput{
				Type:  scorecards.EntityManagementSyncGroupRuleConditionType(cm["type"].(string)),
				Value: cm["value"].(string),
			})
		}
		rules = append(rules, scorecards.EntityManagementSyncGroupRuleUpdateInput{Conditions: conds})
	}
	return scorecards.EntityManagementSyncGroupsSettingsUpdateInput{
		Enabled: enabled,
		Rules:   rules,
	}
}
