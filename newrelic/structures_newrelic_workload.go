package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
)

func expandWorkloadCreateInput(d *schema.ResourceData) workloads.WorkloadCreateInput {
	createInput := workloads.WorkloadCreateInput{
		Name: d.Get("name").(string),
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

	return createInput
}

func expandWorkloadUpdateInput(d *schema.ResourceData) workloads.WorkloadUpdateInput {
	name := d.Get("name").(string)
	updateInput := workloads.WorkloadUpdateInput{
		Name: name,
	}

	if e, ok := d.GetOk("entity_guids"); ok {
		updateInput.EntityGUIDs = expandWorkloadEntityGUIDs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("entity_search_query"); ok {
		updateInput.EntitySearchQueries = expandWorkloadEntitySearchQueryUpdateInputs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("scope_account_ids"); ok {
		updateInput.ScopeAccounts = expandWorkloadScopeAccountsInput(e.(*schema.Set).List())
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

func expandWorkloadEntitySearchQueryUpdateInputs(cfg []interface{}) []workloads.WorkloadUpdateCollectionEntitySearchQueryInput {
	if len(cfg) == 0 {
		return []workloads.WorkloadUpdateCollectionEntitySearchQueryInput{}
	}

	perms := make([]workloads.WorkloadUpdateCollectionEntitySearchQueryInput, len(cfg))

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		entitySearchQuery := expandWorkloadEntitySearchQueryUpdateInput(cfg)

		perms[i] = entitySearchQuery
	}

	return perms
}

func expandWorkloadEntitySearchQueryUpdateInput(cfg map[string]interface{}) workloads.WorkloadUpdateCollectionEntitySearchQueryInput {
	queryInput := workloads.WorkloadUpdateCollectionEntitySearchQueryInput{}

	if query, ok := cfg["query"]; ok {
		queryInput.Query = query.(string)
	}

	return queryInput
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

	return nil
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
