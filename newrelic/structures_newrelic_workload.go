package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
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
	if e, ok := d.GetOk("description"); ok {
		createInput.Description = e.(string)
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
		updateInput.EntitySearchQueries = expandWorkloadUpdateCollectionEntitySearchQueryInputs(e.(*schema.Set).List())
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
		string(workloads.WorkloadResultingGroupTypeTypes.REGULAR_GROUP),
		string(workloads.WorkloadResultingGroupTypeTypes.REMAINING_ENTITIES),
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
