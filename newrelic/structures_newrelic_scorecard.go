package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── Progress levels ───────────────────────────────────────────────────────────

func expandProgressLevels(raw []interface{}) []scorecards.EntityManagementProgressLevelDefinitionCreateInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]scorecards.EntityManagementProgressLevelDefinitionCreateInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		pl := scorecards.EntityManagementProgressLevelDefinitionCreateInput{
			ID:   m["id"].(string),
			Name: m["name"].(string),
		}
		if v, ok := m["description"].(string); ok {
			pl.Description = v
		}
		if v, ok := m["hex_color_code"].(string); ok {
			pl.HexColorCode = v
		}
		out = append(out, pl)
	}
	return out
}

func flattenProgressLevels(levels []scorecards.EntityManagementProgressLevelDefinition) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(levels))
	for _, l := range levels {
		out = append(out, map[string]interface{}{
			"id":             l.ID,
			"name":           l.Name,
			"description":    l.Description,
			"hex_color_code": l.HexColorCode,
		})
	}
	return out
}

// ── Rule IDs (for scorecard ↔ rule attachment) ────────────────────────────────

func expandRuleIDsFromSet(s *schema.Set) []string {
	out := make([]string, 0, s.Len())
	for _, raw := range s.List() {
		out = append(out, raw.(string))
	}
	return out
}

// ── NRQL engine ───────────────────────────────────────────────────────────────

// expandIntListFromInterface converts a Terraform []interface{} of ints
// to []int. Shared between NRQL engine accounts and join_accounts.
func expandIntListFromInterface(raw []interface{}) []int {
	out := make([]int, 0, len(raw))
	for _, v := range raw {
		out = append(out, v.(int))
	}
	return out
}

// expandNRQLEngineCreate maps the nrql_engine block to the Create input type.
func expandNRQLEngineCreate(raw []interface{}) *scorecards.EntityManagementNRQLRuleEngineCreateInput {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	engine := &scorecards.EntityManagementNRQLRuleEngineCreateInput{
		Query:    m["query"].(string),
		Accounts: expandIntListFromInterface(m["accounts"].([]interface{})),
	}
	if v, ok := m["join_accounts"].([]interface{}); ok && len(v) > 0 {
		engine.JoinAccounts = expandIntListFromInterface(v)
	}
	return engine
}

// expandNRQLEngineUpdate maps the nrql_engine block to the Update input type.
// The input types differ in name only; the field set is identical.
func expandNRQLEngineUpdate(raw []interface{}) *scorecards.EntityManagementNRQLRuleEngineUpdateInput {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	engine := &scorecards.EntityManagementNRQLRuleEngineUpdateInput{
		Query:    m["query"].(string),
		Accounts: expandIntListFromInterface(m["accounts"].([]interface{})),
	}
	if v, ok := m["join_accounts"].([]interface{}); ok && len(v) > 0 {
		engine.JoinAccounts = expandIntListFromInterface(v)
	}
	return engine
}

func flattenNRQLEngine(engine scorecards.EntityManagementNRQLRuleEngine) []map[string]interface{} {
	accounts := make([]int, len(engine.Accounts))
	copy(accounts, engine.Accounts)

	joinAccounts := make([]int, len(engine.JoinAccounts))
	copy(joinAccounts, engine.JoinAccounts)

	return []map[string]interface{}{{
		"query":         engine.Query,
		"accounts":      accounts,
		"join_accounts": joinAccounts,
	}}
}
