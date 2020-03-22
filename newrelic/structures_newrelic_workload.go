package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

func expandWorkloadCreateInput(d *schema.ResourceData) workloads.CreateInput {
	createInput := workloads.CreateInput{
		Name: d.Get("name").(string),
	}

	if e, ok := d.GetOk("entity_guids"); ok {
		createInput.EntityGUIDs = expandWorkloadEntityGUIDs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("entity_search_query"); ok {
		createInput.EntitySearchQueries = expandWorkloadEntitySearchQueryInputs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("scope_account_ids"); ok {
		createInput.ScopeAccountsInput = expandWorkloadScopeAccountsInput(e.(*schema.Set).List())
	}

	return createInput
}

func expandWorkloadUpdateInput(d *schema.ResourceData) workloads.UpdateInput {
	name := d.Get("name").(string)
	updateInput := workloads.UpdateInput{
		Name: &name,
	}

	if e, ok := d.GetOk("entity_guids"); ok {
		updateInput.EntityGUIDs = expandWorkloadEntityGUIDs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("entity_search_query"); ok {
		updateInput.EntitySearchQueries = expandWorkloadEntitySearchQueryInputs(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("scope_account_ids"); ok {
		updateInput.ScopeAccountsInput = expandWorkloadScopeAccountsInput(e.(*schema.Set).List())
	}

	return updateInput
}

func expandWorkloadEntityGUIDs(cfg []interface{}) []string {
	if len(cfg) == 0 {
		return []string{}
	}

	perms := make([]string, len(cfg))

	for i, rawCfg := range cfg {
		perms[i] = rawCfg.(string)
	}

	return perms
}

func expandWorkloadEntitySearchQueryInputs(cfg []interface{}) []workloads.EntitySearchQueryInput {
	if len(cfg) == 0 {
		return []workloads.EntitySearchQueryInput{}
	}

	perms := make([]workloads.EntitySearchQueryInput, len(cfg))

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		entitySearchQuery := expandWorkloadEntitySearchQueryInput(cfg)

		perms[i] = entitySearchQuery
	}

	return perms
}

func expandWorkloadEntitySearchQueryInput(cfg map[string]interface{}) workloads.EntitySearchQueryInput {
	queryInput := workloads.EntitySearchQueryInput{}

	if name, ok := cfg["name"]; ok {
		nameStr := name.(string)
		queryInput.Name = &nameStr
	}

	if query, ok := cfg["query"]; ok {
		queryInput.Query = query.(string)
	}

	return queryInput
}

func expandWorkloadScopeAccountsInput(cfg []interface{}) *workloads.ScopeAccountsInput {
	scopeAccounts := workloads.ScopeAccountsInput{}

	for _, a := range cfg {
		scopeAccounts.AccountIDs = append(scopeAccounts.AccountIDs, a.(int))
	}

	return &scopeAccounts
}

func flattenWorkload(workload *workloads.Workload, d *schema.ResourceData) error {
	d.Set("account_id", workload.Account.ID)
	d.Set("guid", workload.GUID)
	d.Set("workload_id", workload.ID)
	d.Set("name", workload.Name)
	d.Set("permalink", workload.Permalink)
	d.Set("composite_entity_search_query", workload.EntitySearchQuery)

	d.Set("entity_guids", flattenWorkloadEntityGUIDs(workload.Entities))
	d.Set("entity_search_query", flattenWorkloadEntitySearchQueries(workload.EntitySearchQueries))
	d.Set("scope_account_ids", workload.ScopeAccounts.AccountIDs)

	return nil
}

func flattenWorkloadEntityGUIDs(in []workloads.EntityRef) interface{} {
	out := make([]interface{}, len(in))

	for i, e := range in {
		out[i] = e.GUID
	}

	return out
}

func flattenWorkloadEntitySearchQueries(in []workloads.EntitySearchQuery) interface{} {
	out := make([]interface{}, len(in))

	for i, e := range in {
		m := make(map[string]interface{})
		m["name"] = e.Name
		m["query"] = e.Query

		out[i] = m
	}

	return out
}
