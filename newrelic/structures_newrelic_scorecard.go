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

func expandNRQLEngineCreate(raw []interface{}) *scorecards.EntityManagementNRQLRuleEngineCreateInput {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	engine := &scorecards.EntityManagementNRQLRuleEngineCreateInput{
		Query: m["query"].(string),
	}
	for _, a := range m["accounts"].([]interface{}) {
		engine.Accounts = append(engine.Accounts, a.(int))
	}
	if v, ok := m["join_accounts"].([]interface{}); ok {
		for _, a := range v {
			engine.JoinAccounts = append(engine.JoinAccounts, a.(int))
		}
	}
	return engine
}

func expandNRQLEngineUpdate(raw []interface{}) *scorecards.EntityManagementNRQLRuleEngineUpdateInput {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	engine := &scorecards.EntityManagementNRQLRuleEngineUpdateInput{
		Query: m["query"].(string),
	}
	for _, a := range m["accounts"].([]interface{}) {
		engine.Accounts = append(engine.Accounts, a.(int))
	}
	if v, ok := m["join_accounts"].([]interface{}); ok {
		for _, a := range v {
			engine.JoinAccounts = append(engine.JoinAccounts, a.(int))
		}
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
