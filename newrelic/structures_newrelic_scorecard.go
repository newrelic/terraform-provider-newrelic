package newrelic

import (
	"context"
	"sort"

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
// The output is sorted alphabetically by "id" so that the state always has a canonical
// order. This prevents spurious TypeList diffs when the config lists the same levels
// in a different order — as long as the config is also written in alphabetical-by-id order
// the plan will show "No changes" after create/update.
func flattenProgressLevels(levels []scorecards.EntityManagementProgressLevelDefinition) []map[string]interface{} {
	// Sort a copy so we don't mutate the caller's slice.
	sorted := make([]scorecards.EntityManagementProgressLevelDefinition, len(levels))
	copy(sorted, levels)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })

	out := make([]map[string]interface{}, 0, len(sorted))
	for _, l := range sorted {
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

// nrqlEngineParams holds the parsed values from an nrql_engine block, shared
// between the Create and Update expand functions to avoid duplication.
type nrqlEngineParams struct {
	Query        string
	Accounts     []int
	JoinAccounts []int
}

// extractNRQLEngineParams parses a raw nrql_engine Terraform block into the
// shared nrqlEngineParams struct.
//
// nrql_engine is Required in the schema; Terraform validates its presence
// before invoking Create/Update, so raw[0] == nil is unreachable in practice.
func extractNRQLEngineParams(raw []interface{}) *nrqlEngineParams {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	p := &nrqlEngineParams{
		Query: m["query"].(string),
	}
	// accounts and join_accounts are TypeList — use .([]interface{})
	if accts, ok := m["accounts"].([]interface{}); ok {
		p.Accounts = expandIntListFromInterface(accts)
	}
	if joinAccts, ok := m["join_accounts"].([]interface{}); ok && len(joinAccts) > 0 {
		p.JoinAccounts = expandIntListFromInterface(joinAccts)
	}
	return p
}

// expandNRQLEngineCreate maps the nrql_engine block to the Create input type.
func expandNRQLEngineCreate(raw []interface{}) *scorecards.EntityManagementNRQLRuleEngineCreateInput {
	p := extractNRQLEngineParams(raw)
	if p == nil {
		return nil
	}
	return &scorecards.EntityManagementNRQLRuleEngineCreateInput{
		Query:        p.Query,
		Accounts:     p.Accounts,
		JoinAccounts: p.JoinAccounts,
	}
}

// expandNRQLEngineUpdate maps the nrql_engine block to the Update input type.
// The input types differ in name only; the field set is identical.
func expandNRQLEngineUpdate(raw []interface{}) *scorecards.EntityManagementNRQLRuleEngineUpdateInput {
	p := extractNRQLEngineParams(raw)
	if p == nil {
		return nil
	}
	return &scorecards.EntityManagementNRQLRuleEngineUpdateInput{
		Query:        p.Query,
		Accounts:     p.Accounts,
		JoinAccounts: p.JoinAccounts,
	}
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
