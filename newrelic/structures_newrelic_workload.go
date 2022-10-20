package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
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

// Static
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

// Update Static
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

// Automatic
func expandWorkloadStatusConfigAutomaticInput(rcfg []interface{}) *workloads.WorkloadAutomaticStatusInput {
	prem := workloads.WorkloadAutomaticStatusInput{}
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})

		if x, ok := cfg["enabled"]; ok {
			prem.Enabled = x.(bool)
		}

		if x := cfg["remaining_entities_rule"]; x.(*schema.Set).Len() != 0 {
			prem.RemainingEntitiesRule = expandRemainingEntityRule(x.(*schema.Set).List())
		}

		if x, ok := cfg["rule"]; ok {
			prem.Rules = expandAutoConfigRule(x.(*schema.Set).List())
		}
	}

	return &prem
}

func expandRemainingEntityRule(rcfg []interface{}) *workloads.WorkloadRemainingEntitiesRuleInput {
	var prem workloads.WorkloadRemainingEntitiesRuleInput
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})

		if x := cfg["remaining_entities_rule_rollup"]; x.(*schema.Set).Len() != 0 {
			prem.Rollup = expandRemainingEntityRuleRollup(x.(*schema.Set).List())
		}
	}
	return &prem
}

func expandRemainingEntityRuleRollup(rcfg []interface{}) *workloads.WorkloadRemainingEntitiesRuleRollupInput {
	prem := &workloads.WorkloadRemainingEntitiesRuleRollupInput{}
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
		}

		if x, ok := cfg["threshold_value"]; ok {
			prem.ThresholdValue = x.(int)
		}
	}
	return prem
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
		}
	}

	return &inp
}

// Update Automatic
func expandWorkloadStatusConfigUpdateAutomaticInput(rcfg []interface{}) *workloads.WorkloadUpdateAutomaticStatusInput {
	prem := workloads.WorkloadUpdateAutomaticStatusInput{}
	for _, v := range rcfg {
		cfg := v.(map[string]interface{})

		if x, ok := cfg["enabled"]; ok {
			prem.Enabled = x.(bool)
		}

		if x := cfg["remaining_entities_rule"]; x.(*schema.Set).Len() != 0 {
			prem.RemainingEntitiesRule = expandRemainingEntityRule(x.(*schema.Set).List())
		}

		if x, ok := cfg["rule"]; ok {
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

func listValidWorkloadStatuses() []string {
	return []string{
		string(workloads.WorkloadStatusValueInputTypes.DEGRADED),
		string(workloads.WorkloadStatusValueInputTypes.OPERATIONAL),
		string(workloads.WorkloadStatusValueInputTypes.DISRUPTED),
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
func flattenWorkload(workload *workloads.WorkloadCollection, d *schema.ResourceData) error {
	_ = d.Set("account_id", workload.Account.ID)
	_ = d.Set("guid", workload.GUID)
	_ = d.Set("workload_id", workload.ID)
	_ = d.Set("name", workload.Name)
	_ = d.Set("permalink", workload.Permalink)
	_ = d.Set("composite_entity_search_query", workload.EntitySearchQuery)
	_ = d.Set("entity_guids", flattenWorkloadEntityGUIDs(workload.Entities))
	_ = d.Set("entity_search_query", flattenWorkloadEntitySearchQueries(workload.EntitySearchQueries))
	_ = d.Set("scope_account_ids", workload.ScopeAccounts.AccountIDs)

	if workload.Description != "" {
		_ = d.Set("description", workload.Description)
	}

	if workload.StatusConfig.Static != nil {
		statusStatic := flattenStatusConfigStatic(workload.StatusConfig.Static)
		if err := d.Set("status_config_static", statusStatic); err != nil {
			return err
		}
	}

	return nil
}

func flattenStatusConfigStatic(in []workloads.WorkloadStaticStatus) []interface{} {
	out := make([]interface{}, len(in))

	for _, p := range in {
		m := make(map[string]interface{})
		m["enabled"] = p.Enabled

		if p.Description != "" {
			m["description"] = p.Description
		}

		if p.Summary != "" {
			m["summary"] = p.Summary
		}

		m["status"] = p.Status
		out[0] = m
	}
	return out
}

func flattenWorkloadEntityGUIDs(in []workloads.WorkloadEntityRef) interface{} {
	out := make([]interface{}, len(in))
	for i, e := range in {
		out[i] = e.GUID
	}
	return out
}

func flattenWorkloadEntitySearchQueries(in []workloads.WorkloadEntitySearchQuery) interface{} {
	out := make([]interface{}, len(in))
	for i, e := range in {
		m := make(map[string]interface{})
		m["query"] = e.Query
		out[i] = m
	}
	return out
}
