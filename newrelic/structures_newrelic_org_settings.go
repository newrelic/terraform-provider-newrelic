package newrelic

import "github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"

// flattenSyncGroupRules converts API SyncGroupRule values to the Terraform list
// shape used by the sync_group_rules attribute.
func flattenSyncGroupRules(rules []servicearchintelligence.EntityManagementSyncGroupRule) []map[string]interface{} {
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
func expandSyncGroupsUpdate(enabled bool, rawRules []interface{}) servicearchintelligence.EntityManagementSyncGroupsSettingsUpdateInput {
	rules := make([]servicearchintelligence.EntityManagementSyncGroupRuleUpdateInput, 0, len(rawRules))
	for _, rr := range rawRules {
		if rr == nil {
			continue
		}
		rm, ok := rr.(map[string]interface{})
		if !ok {
			continue
		}
		rawConds, _ := rm["conditions"].([]interface{})
		conds := make([]servicearchintelligence.EntityManagementSyncGroupRuleConditionUpdateInput, 0, len(rawConds))
		for _, rc := range rawConds {
			if rc == nil {
				continue
			}
			cm, ok := rc.(map[string]interface{})
			if !ok {
				continue
			}
			conds = append(conds, servicearchintelligence.EntityManagementSyncGroupRuleConditionUpdateInput{
				Type:  servicearchintelligence.EntityManagementSyncGroupRuleConditionType(cm["type"].(string)),
				Value: cm["value"].(string),
			})
		}
		rules = append(rules, servicearchintelligence.EntityManagementSyncGroupRuleUpdateInput{Conditions: conds})
	}
	return servicearchintelligence.EntityManagementSyncGroupsSettingsUpdateInput{
		Enabled: enabled,
		Rules:   rules,
	}
}
