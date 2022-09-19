package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
	"log"
)

func expandWorkloadCreateInput(d *schema.ResourceData) workloads.WorkloadCreateInput {
	createInput := workloads.WorkloadCreateInput{
		Name:         d.Get("name").(string),
		StatusConfig: &workloads.WorkloadStatusConfigInput{},
	}

	if e, ok := d.GetOk("entity_guids"); ok {
		createInput.EntityGUIDs = expandWorkloadEntityGUIDs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("entity_search_query"); ok {
		createInput.EntitySearchQueries = expandWorkloadEntitySearchQueryInputs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("scope_account_ids"); ok {
		createInput.ScopeAccounts = expandWorkloadScopeAccountsInput(e.(*schema.Set).List())
	}
	if e, ok := d.GetOk("description"); ok {
		createInput.Description = e.(string)
	}
	if e, ok := d.GetOk("status_config_static"); ok {
		createInput.StatusConfig.Static = expandWorkloadStatusConfigStaticInput(e.(*schema.Set).List())
	}
	if e, ok := d.GetOk("status_config_automatic"); ok {
		createInput.StatusConfig.Automatic = expandWorkloadStatusConfigAutomaticInput(e.(*schema.Set).List())
	}
	log.Printf("[INFO] pks automatic value %v", createInput.StatusConfig.Automatic)
	log.Printf("[INFO] pks StatusConfig value %v", createInput.StatusConfig)

	return createInput
}

func expandWorkloadUpdateInput(d *schema.ResourceData) workloads.WorkloadUpdateInput {
	name := d.Get("name").(string)
	updateInput := workloads.WorkloadUpdateInput{
		Name:         name,
		StatusConfig: &workloads.WorkloadUpdateStatusConfigInput{},
	}

	if e, ok := d.GetOk("entity_guids"); ok {
		updateInput.EntityGUIDs = expandWorkloadEntityGUIDs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("entity_search_query"); ok {
		updateInput.EntitySearchQueries = expandWorkloadUpdateCollectionEntitySearchQueryInputs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("scope_account_ids"); ok {
		updateInput.ScopeAccounts = expandWorkloadScopeAccountsInput(e.(*schema.Set).List())
	}
	if e, ok := d.GetOk("description"); ok {
		updateInput.Description = e.(string)
	}
	if e, ok := d.GetOk("status_config_static"); ok {
		updateInput.StatusConfig.Static = expandWorkloadUpdateStatusConfigStaticInput(e.(*schema.Set).List())
	}
	if e, ok := d.GetOk("status_config_automatic"); ok {
		updateInput.StatusConfig.Automatic = expandWorkloadStatusConfigUpdateAutomaticInput(e.(*schema.Set).List())
	}
	log.Printf("[INFO] pks automatic value %v", updateInput.StatusConfig.Automatic)

	return updateInput
}

func expandWorkloadEntityGUIDs(cfg []interface{}) []common.EntityGUID {
	if len(cfg) == 0 {
		return []common.EntityGUID{}
	}

	perms := make([]common.EntityGUID, len(cfg))

	for i, rawCfg := range cfg {
		perms[i] = common.EntityGUID(rawCfg.(string))
	}

	return perms
}

func expandWorkloadEntitySearchQueryInputs(cfg []interface{}) []workloads.WorkloadEntitySearchQueryInput {
	if len(cfg) == 0 {
		return []workloads.WorkloadEntitySearchQueryInput{}
	}

	perms := make([]workloads.WorkloadEntitySearchQueryInput, len(cfg))

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		entitySearchQuery := expandWorkloadEntitySearchQueryInput(cfg)

		perms[i] = entitySearchQuery
	}

	return perms
}

func expandWorkloadEntitySearchQueryInput(cfg map[string]interface{}) workloads.WorkloadEntitySearchQueryInput {
	queryInput := workloads.WorkloadEntitySearchQueryInput{}

	if query, ok := cfg["query"]; ok {
		queryInput.Query = query.(string)
	}

	return queryInput
}

func expandWorkloadScopeAccountsInput(cfg []interface{}) *workloads.WorkloadScopeAccountsInput {
	scopeAccounts := workloads.WorkloadScopeAccountsInput{}

	for _, a := range cfg {
		scopeAccounts.AccountIDs = append(scopeAccounts.AccountIDs, a.(int))
	}

	return &scopeAccounts
}

func expandWorkloadUpdateCollectionEntitySearchQueryInputs(cfg []interface{}) []workloads.WorkloadUpdateCollectionEntitySearchQueryInput {
	if len(cfg) == 0 {
		return []workloads.WorkloadUpdateCollectionEntitySearchQueryInput{}
	}

	perms := make([]workloads.WorkloadUpdateCollectionEntitySearchQueryInput, len(cfg))

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		entitySearchQuery := expandWorkloadUpdateCollectionEntitySearchQueryInput(cfg)

		perms[i] = entitySearchQuery
	}

	return perms
}

func expandWorkloadUpdateCollectionEntitySearchQueryInput(cfg map[string]interface{}) workloads.WorkloadUpdateCollectionEntitySearchQueryInput {
	queryInput := workloads.WorkloadUpdateCollectionEntitySearchQueryInput{}

	if query, ok := cfg["query"]; ok {
		queryInput.Query = query.(string)
	}

	return queryInput
}

//Static
func expandWorkloadStatusConfigStaticInput(e []interface{}) []workloads.WorkloadStaticStatusInput {
	staticOut := make([]workloads.WorkloadStaticStatusInput, len(e))

	for i, s := range e {
		setStatic := s.(map[string]interface{})
		x := workloads.WorkloadStaticStatusInput{
			Enabled:     setStatic["enabled"].(bool),
			Status:      workloads.WorkloadStatusValueInput(setStatic["status"].(string)),
			Summary:     setStatic["summary"].(string),
			Description: setStatic["description"].(string),
		}
		staticOut[i] = x
	}

	return staticOut
}

//Update Static
func expandWorkloadUpdateStatusConfigStaticInput(cfg []interface{}) []workloads.WorkloadUpdateStaticStatusInput {
	staticOut := make([]workloads.WorkloadUpdateStaticStatusInput, len(cfg))

	for i, s := range cfg {
		setStatic := s.(map[string]interface{})
		x := workloads.WorkloadUpdateStaticStatusInput{
			Enabled:     setStatic["enabled"].(bool),
			Status:      workloads.WorkloadStatusValueInput(setStatic["status"].(string)),
			Summary:     setStatic["summary"].(string),
			Description: setStatic["description"].(string),
		}
		staticOut[i] = x
	}

	return staticOut
}

//TRY 1
//Automatic
func expandWorkloadStatusConfigAutomaticInput(rcfg []interface{}) *workloads.WorkloadAutomaticStatusInput {
	prem := workloads.WorkloadAutomaticStatusInput{
		RemainingEntitiesRule: &workloads.WorkloadRemainingEntitiesRuleInput{},
	}
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})
		if x, ok := cfg["enabled"]; ok {
			prem.Enabled = x.(bool)
		}
		if x, ok := cfg["remaining_entities_rule_rollup"]; ok {
			prem.RemainingEntitiesRule.Rollup = expandRemainingEntityRuleRollup(x.(*schema.Set).List())
		}
		if x, ok := cfg["rules"]; ok {
			prem.Rules = expandAutoConfigRule(x.(*schema.Set).List())
		}
		log.Printf("[INFO] pks RemainingEntitiesRuleRollup value %v", prem.RemainingEntitiesRule.Rollup)
		log.Printf("[INFO] pks RemainingEntitiesRule value %v", prem.RemainingEntitiesRule)
	}

	return &prem
}

func expandRemainingEntityRuleRollup(rcfg []interface{}) *workloads.WorkloadRemainingEntitiesRuleRollupInput {
	var prem workloads.WorkloadRemainingEntitiesRuleRollupInput
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})
		if x, ok := cfg["group_by"]; ok {
			prem.GroupBy = workloads.WorkloadGroupRemainingEntitiesRuleBy(x.(string))
		}
		if x, ok := cfg["strategy"]; ok {
			prem.Strategy = workloads.WorkloadRollupStrategy(x.(string))
		}
		if x, ok := cfg["threshold_type"]; ok {
			prem.ThresholdType = (workloads.WorkloadRuleThresholdType)(x.(string))
			//s := x.(string)
			//prem.ThresholdType = (*workloads.WorkloadRuleThresholdType)(&s)

		}
		if x, ok := cfg["threshold_value"]; ok {
			prem.ThresholdValue = x.(int)
		}
	}
	return &prem
}

func expandAutoConfigRule(cfg []interface{}) []workloads.WorkloadRegularRuleInput {
	if len(cfg) == 0 {
		return []workloads.WorkloadRegularRuleInput{}
	}
	perms := make([]workloads.WorkloadRegularRuleInput, len(cfg))
	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		rules := expandRules(cfg)
		perms[i] = rules
	}

	return perms
}

func expandRules(cfg map[string]interface{}) workloads.WorkloadRegularRuleInput {
	inp := workloads.WorkloadRegularRuleInput{}

	if x, ok := cfg["entity_guids"]; ok {
		inp.EntityGUIDs = expandWorkloadEntityGUIDs(x.(*schema.Set).List())
	}
	if x, ok := cfg["nrql_query"]; ok {
		inp.EntitySearchQueries = expandWorkloadEntitySearchQueryInputs(x.(*schema.Set).List())
	}
	if x, ok := cfg["rollup"]; ok {
		inp.Rollup = expandRuleRollUp(x.(*schema.Set).List())
	}
	return inp
}

func expandRuleRollUp(rcfg []interface{}) *workloads.WorkloadRollupInput {
	var inp workloads.WorkloadRollupInput
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})
		if x, ok := cfg["strategy"]; ok {
			inp.Strategy = workloads.WorkloadRollupStrategy(x.(string))
		}
		if x, ok := cfg["threshold_value"]; ok {
			inp.ThresholdValue = x.(int)
		}
		if x, ok := cfg["threshold_type"]; ok {
			inp.ThresholdType = workloads.WorkloadRuleThresholdType(x.(string))
			//s := x.(string)
			//inp.ThresholdType = (*workloads.WorkloadRuleThresholdType)(&s)
		}
	}

	return &inp
}

//Update Automatic
func expandWorkloadStatusConfigUpdateAutomaticInput(rcfg []interface{}) *workloads.WorkloadUpdateAutomaticStatusInput {
	prem := workloads.WorkloadUpdateAutomaticStatusInput{
		RemainingEntitiesRule: &workloads.WorkloadRemainingEntitiesRuleInput{},
	}
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})
		if x, ok := cfg["enabled"]; ok {
			prem.Enabled = x.(bool)
		}
		if x, ok := cfg["remaining_entities_rule_rollup"]; ok {
			prem.RemainingEntitiesRule.Rollup = expandRemainingEntityRuleRollup(x.(*schema.Set).List())
		}
		if x, ok := cfg["rules"]; ok {
			prem.Rules = expandUpdateAutoConfigRule(x.(*schema.Set).List())
		}
	}

	return &prem
}

func expandUpdateAutoConfigRule(cfg []interface{}) []workloads.WorkloadUpdateRegularRuleInput {
	if len(cfg) == 0 {
		return []workloads.WorkloadUpdateRegularRuleInput{}
	}
	perms := make([]workloads.WorkloadUpdateRegularRuleInput, len(cfg))
	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		rules := expandUpdateRules(cfg)
		perms[i] = rules
	}

	return perms
}

func expandUpdateRules(cfg map[string]interface{}) workloads.WorkloadUpdateRegularRuleInput {
	inp := workloads.WorkloadUpdateRegularRuleInput{}

	if x, ok := cfg["entity_guids"]; ok {
		inp.EntityGUIDs = expandWorkloadEntityGUIDs(x.(*schema.Set).List())
	}
	if x, ok := cfg["nrql_query"]; ok {
		inp.EntitySearchQueries = expandWorkloadUpdateCollectionEntitySearchQueryInputs(x.(*schema.Set).List())
	}
	if x, ok := cfg["rollup"]; ok {
		inp.Rollup = expandRuleRollUp(x.(*schema.Set).List())
	}
	return inp
}

// Handles setting simple string attributes in the schema. If the attribute/key is
// invalid or the value is not a correct type, an error will be returned.
func setWorkloadAttributes(d *schema.ResourceData, attributes map[string]string) error {
	for key := range attributes {
		err := d.Set(key, attributes[key])
		if err != nil {
			return err
		}
	}

	return nil
}

func listValidWorkloadStatuses() []string {
	return []string{
		string(workloads.WorkloadStatusValueTypes.DEGRADED),
		string(workloads.WorkloadStatusValueTypes.OPERATIONAL),
		string(workloads.WorkloadStatusValueTypes.UNKNOWN),
		string(workloads.WorkloadStatusValueTypes.DISRUPTED),
	}
}

func listValidWorkloadGroupBy() []string {
	return []string{
		string(workloads.WorkloadGroupRemainingEntitiesRuleByTypes.ENTITY_TYPE),
		string(workloads.WorkloadGroupRemainingEntitiesRuleByTypes.NONE),
	}
}

func listValidWorkloadStrategy() []string {
	return []string{
		string(workloads.WorkloadRollupStrategyTypes.BEST_STATUS_WINS),
		string(workloads.WorkloadRollupStrategyTypes.WORST_STATUS_WINS),
	}
}

func listValidWorkloadRuleThresholdType() []string {
	return []string{
		string(workloads.WorkloadRuleThresholdTypeTypes.FIXED),
		string(workloads.WorkloadRuleThresholdTypeTypes.PERCENTAGE),
	}
}

////Try 2
//func expandWorkloadStatusConfigAutomaticInput(rcfg []interface{}) *workloads.WorkloadAutomaticStatusInput {
//	for _, v := range rcfg {
//		cfg := v.(map[string]interface{})
//		prem := workloads.WorkloadAutomaticStatusInput{
//			Enabled: cfg["enabled"].(bool),
//			RemainingEntitiesRule: workloads.WorkloadRemainingEntitiesRuleInput{
//				Rollup: expandRemainingEntityRuleRollup(cfg["remaining_entities_rule_rollup"].(*schema.Set).List()),
//			},
//			Rules: expandAutoConfigRule(cfg["rules"].(*schema.Set).List()),
//		}
//		log.Printf("[INFO] pks RemainingEntitiesRuleRollup value %v", prem.RemainingEntitiesRule.Rollup)
//		log.Printf("[INFO] pks RemainingEntitiesRule value %v", prem.RemainingEntitiesRule)
//
//		return &prem
//	}
//
//	return nil
//}
//
//func expandRemainingEntityRuleRollup(rollResource []interface{}) *workloads.WorkloadRemainingEntitiesRuleRollupInput {
//	var x workloads.WorkloadRemainingEntitiesRuleRollupInput
//	for _, r := range rollResource {
//		setRoll := r.(map[string]interface{})
//		v := setRoll["threshold_type"].(string)
//		x = workloads.WorkloadRemainingEntitiesRuleRollupInput{
//			GroupBy:        workloads.WorkloadGroupRemainingEntitiesRuleBy(setRoll["group_by"].(string)),
//			Strategy:       workloads.WorkloadRollupStrategy(setRoll["strategy"].(string)),
//			ThresholdValue: setRoll["threshold_value"].(int),
//			ThresholdType:  (workloads.WorkloadRuleThresholdType)(v),
//		}
//	}
//	return &x
//}
//
//func expandAutoConfigRule(list []interface{}) []workloads.WorkloadRegularRuleInput {
//	ruleOut := make([]workloads.WorkloadRegularRuleInput, len(list))
//	for i, r := range list {
//		setRule := r.(map[string]interface{})
//		x := workloads.WorkloadRegularRuleInput{
//			EntityGUIDs:         expandWorkloadEntityGUIDs(setRule["entity_guids"].(*schema.Set).List()),
//			EntitySearchQueries: expandWorkloadEntitySearchQueryInputs(setRule["nrql_query"].(*schema.Set).List()),
//			Rollup:              expandRuleRollUp(setRule["rollup"].(*schema.Set).List()),
//		}
//		ruleOut[i] = x
//	}
//	return ruleOut
//}
//func expandRuleRollUp(rollResource []interface{}) *workloads.WorkloadRollupInput {
//	var x workloads.WorkloadRollupInput
//	for _, r := range rollResource {
//		setRoll := r.(map[string]interface{})
//		v := setRoll["threshold_type"].(string)
//		x = workloads.WorkloadRollupInput{
//			Strategy:       workloads.WorkloadRollupStrategy(setRoll["strategy"].(string)),
//			ThresholdValue: setRoll["threshold_value"].(int),
//			ThresholdType:  (workloads.WorkloadRuleThresholdType)(v),
//		}
//	}
//	return &x
//}
//
//func expandWorkloadStatusConfigUpdateAutomaticInput(rcfg []interface{}) *workloads.WorkloadUpdateAutomaticStatusInput {
//	for _, v := range rcfg {
//		cfg := v.(map[string]interface{})
//		prem := workloads.WorkloadUpdateAutomaticStatusInput{
//			Enabled: cfg["enabled"].(bool),
//			RemainingEntitiesRule: workloads.WorkloadRemainingEntitiesRuleInput{
//				Rollup: expandRemainingEntityRuleRollup(cfg["remaining_entities_rule_rollup"].(*schema.Set).List()),
//			},
//			Rules: expandUpdateAutoConfigRule(cfg["rules"].(*schema.Set).List()),
//		}
//		return &prem
//	}
//
//	return nil
//}
//func expandUpdateAutoConfigRule(list []interface{}) []workloads.WorkloadUpdateRegularRuleInput {
//	ruleOut := make([]workloads.WorkloadUpdateRegularRuleInput, len(list))
//	for i, r := range list {
//		setRule := r.(map[string]interface{})
//		x := workloads.WorkloadUpdateRegularRuleInput{
//			EntityGUIDs:         expandWorkloadEntityGUIDs(setRule["entity_guids"].(*schema.Set).List()),
//			EntitySearchQueries: expandWorkloadUpdateCollectionEntitySearchQueryInputs(setRule["nrql_query"].(*schema.Set).List()),
//			Rollup:              expandRuleRollUp(setRule["rollup"].(*schema.Set).List()),
//		}
//		ruleOut[i] = x
//	}
//	return ruleOut
//}
