package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── Progress levels ───────────────────────────────────────────────────────────

// expandProgressLevels converts the progress_levels Terraform list to CreateInput slice.
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

// flattenProgressLevels converts API ProgressLevelDefinition values back to Terraform maps.
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

// expandRuleIDsFromSet extracts string rule GUIDs from a TypeSet of plain strings.
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
	// join_accounts is TypeSet — use .(*schema.Set).List() not .([]interface{})
	if s, ok := m["join_accounts"].(*schema.Set); ok && s.Len() > 0 {
		engine.JoinAccounts = expandIntListFromInterface(s.List())
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
	if s, ok := m["join_accounts"].(*schema.Set); ok && s.Len() > 0 {
		engine.JoinAccounts = expandIntListFromInterface(s.List())
	}
	return engine
}

// flattenNRQLEngine converts an API NRQLRuleEngine value back to the nrql_engine Terraform block.
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

// ── Rules collection reader ───────────────────────────────────────────────────

// readScorecardRuleGUIDs pages through the scorecard's rules collection and
// returns the GUID of every ScorecardRule entity in it.
func readScorecardRuleGUIDs(ctx context.Context, client *scorecards.Scorecards, rulesColID string) ([]string, error) {
	if rulesColID == "" {
		return nil, nil
	}
	var guids []string
	err := pageCollectionItems(ctx, client, rulesColID, func(item scorecards.EntityManagementEntityInterface) {
		if r, ok := item.(*scorecards.EntityManagementScorecardRuleEntity); ok {
			guids = append(guids, r.ID)
		}
	})
	return guids, err
}
