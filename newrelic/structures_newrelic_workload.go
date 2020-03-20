package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

func expandWorkload(d *schema.ResourceData) workloads.Workload {
	workload := workloads.Workload{
		Name: d.Get("name").(string),
		ID:   d.Get("workload_id").(int),
		GUID: d.Get("guid").(string),
	}

	if e, ok := d.GetOk("entity"); ok {
		workload.Entities = expandWorkloadEntities(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("entity_search_query"); ok {
		workload.EntitySearchQueries = expandWorkloadEntitySearchQueries(e.(*schema.Set).List())
	}

	if e, ok := d.GetOk("scope_accounts"); ok {
		workload.ScopeAccounts = expandWorkloadScopeAccounts(e.([]interface{})[0].(map[string]interface{}))
	}

	return workload
}

func expandWorkloadEntities(cfg []interface{}) []workloads.EntityRef {
	if len(cfg) == 0 {
		return []workloads.EntityRef{}
	}

	perms := make([]workloads.EntityRef, len(cfg))

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		guid := cfg["guid"].(string)
		entityRef := workloads.EntityRef{
			GUID: &guid,
		}

		perms[i] = entityRef
	}

	return perms
}

func expandWorkloadEntitySearchQueries(cfg []interface{}) []workloads.EntitySearchQuery {
	if len(cfg) == 0 {
		return []workloads.EntitySearchQuery{}
	}

	perms := make([]workloads.EntitySearchQuery, len(cfg))

	for i, rawCfg := range cfg {
		cfg := rawCfg.(map[string]interface{})
		entitySearchQuery := expandWorkloadEntitySearchQuery(cfg)

		perms[i] = entitySearchQuery
	}

	return perms
}

func expandWorkloadEntitySearchQuery(cfg map[string]interface{}) workloads.EntitySearchQuery {
	entitySearchQuery := workloads.EntitySearchQuery{}

	if name, ok := cfg["name"]; ok {
		entitySearchQuery.Name = name.(string)
	}

	if query, ok := cfg["query"]; ok {
		entitySearchQuery.Query = query.(string)
	}

	return entitySearchQuery
}

func expandWorkloadScopeAccounts(cfg map[string]interface{}) workloads.ScopeAccounts {
	scopeAccounts := workloads.ScopeAccounts{}

	if accountIDs, ok := cfg["account_ids"]; ok {
		for _, a := range accountIDs.(*schema.Set).List() {
			scopeAccounts.AccountIDs = append(scopeAccounts.AccountIDs, a.(int))
		}
	}

	return scopeAccounts
}

func flattenWorkload(workload *workloads.Workload, d *schema.ResourceData) error {
	d.Set("guid", workload.GUID)
	d.Set("workload_id", workload.ID)
	d.Set("name", workload.Name)
	d.Set("permalink", workload.Permalink)

	d.Set("account", flattenWorkloadAccountReference(workload.Account))
	d.Set("created_by", flattenWorkloadUserReference(workload.CreatedBy))
	d.Set("entity", flattenWorkloadEntities(workload.Entities))
	d.Set("entity_search_query", flattenWorkloadEntitySearchQueries(workload.EntitySearchQueries))
	d.Set("scope_accounts", flattenWorkloadScopeAccounts(workload.ScopeAccounts))

	return nil
}

func flattenWorkloadAccountReference(in nerdgraph.AccountReference) interface{} {
	m := make(map[string]interface{})

	m["id"] = in.ID
	m["name"] = in.Name

	return []interface{}{m}
}

func flattenWorkloadUserReference(in workloads.UserReference) interface{} {
	m := make(map[string]interface{})

	m["email"] = in.Email
	m["gravatar"] = in.Gravatar
	m["id"] = in.ID
	m["name"] = in.Name

	return []interface{}{m}
}

func flattenWorkloadEntities(in []workloads.EntityRef) interface{} {
	out := make([]interface{}, len(in))

	for i, e := range in {
		m := make(map[string]interface{})
		m["guid"] = e.GUID

		out[i] = m
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

func flattenWorkloadScopeAccounts(in workloads.ScopeAccounts) interface{} {
	m := make(map[string]interface{})

	m["account_ids"] = in.AccountIDs

	return []interface{}{m}
}
