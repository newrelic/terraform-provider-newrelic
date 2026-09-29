package newrelic

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── CustomizeDiff functions ───────────────────────────────────────────────────

// customizeScorecardRuleDiff validates the nrql_engine accounts/join_accounts
// lists for semantic correctness at plan time:
//
//  1. No account ID ≤ 0 in either list.
//  2. No duplicates within accounts.
//  3. No duplicates within join_accounts.
//  4. No account ID appears in both accounts AND join_accounts.
func customizeScorecardRuleDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	rawEngine := d.Get("nrql_engine").([]interface{})
	if len(rawEngine) == 0 || rawEngine[0] == nil {
		return nil
	}
	m := rawEngine[0].(map[string]interface{})

	toIntSlice := func(raw []interface{}) []int {
		out := make([]int, 0, len(raw))
		for _, v := range raw {
			out = append(out, v.(int))
		}
		return out
	}

	accounts := toIntSlice(m["accounts"].([]interface{}))
	joinAccounts := toIntSlice(m["join_accounts"].([]interface{}))

	// 1. No account ID ≤ 0
	for _, id := range accounts {
		if id <= 0 {
			return fmt.Errorf("nrql_engine.accounts: account ID %d is invalid (must be > 0)", id)
		}
	}
	for _, id := range joinAccounts {
		if id <= 0 {
			return fmt.Errorf("nrql_engine.join_accounts: account ID %d is invalid (must be > 0)", id)
		}
	}

	// 2. No duplicates within accounts
	seen := make(map[int]bool, len(accounts))
	for _, id := range accounts {
		if seen[id] {
			return fmt.Errorf("nrql_engine.accounts: duplicate account ID %d", id)
		}
		seen[id] = true
	}

	// 3. No duplicates within join_accounts
	seenJoin := make(map[int]bool, len(joinAccounts))
	for _, id := range joinAccounts {
		if seenJoin[id] {
			return fmt.Errorf("nrql_engine.join_accounts: duplicate account ID %d", id)
		}
		seenJoin[id] = true
	}

	// 4. No intersection between accounts and join_accounts
	for _, id := range joinAccounts {
		if seen[id] {
			return fmt.Errorf("nrql_engine: account ID %d appears in both accounts and join_accounts — these lists must be disjoint", id)
		}
	}

	return nil
}

// ── Exclusive-membership error detection ──────────────────────────────────────

// isExclusiveMembershipError reports whether an error from
// EntityManagementAddCollectionMembers indicates that one or more rules are
// already attached to a different scorecard. The NGEP API returns two
// different messages for this condition depending on the backend version.
func isExclusiveMembershipError(err error) bool {
	return strings.Contains(err.Error(), "already belongs to collection") ||
		strings.Contains(err.Error(), "already exists in one collection")
}

// ── Progress levels ───────────────────────────────────────────────────────────

// expandProgressLevels converts the progress_levels Terraform set to CreateInput slice.
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

// expandProgressLevelsUpdate converts the progress_levels Terraform set to UpdateInput slice.
// The UpdateInput type has identical fields to CreateInput; it differs only in name.
func expandProgressLevelsUpdate(raw []interface{}) []scorecards.EntityManagementProgressLevelDefinitionUpdateInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]scorecards.EntityManagementProgressLevelDefinitionUpdateInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		pl := scorecards.EntityManagementProgressLevelDefinitionUpdateInput{
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

// progressLevelsCreateToRead converts CreateInput slices to the Definition type used by
// flattenProgressLevels. The two types have identical fields so each element is a
// direct cast. Used by Create to set state without issuing a Read round-trip.
func progressLevelsCreateToRead(in []scorecards.EntityManagementProgressLevelDefinitionCreateInput) []scorecards.EntityManagementProgressLevelDefinition {
	out := make([]scorecards.EntityManagementProgressLevelDefinition, len(in))
	for i, p := range in {
		out[i] = scorecards.EntityManagementProgressLevelDefinition(p)
	}
	return out
}

// flattenProgressLevels converts API ProgressLevelDefinition values back to Terraform maps.
// Sorted alphabetically by id for deterministic plan output.
func flattenProgressLevels(levels []scorecards.EntityManagementProgressLevelDefinition) []map[string]interface{} {
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
	params := &nrqlEngineParams{
		Query: m["query"].(string),
	}
	// accounts and join_accounts are TypeList — use .([]interface{})
	if accounts, ok := m["accounts"].([]interface{}); ok {
		params.Accounts = expandIntListFromInterface(accounts)
	}
	if joinAccounts, ok := m["join_accounts"].([]interface{}); ok && len(joinAccounts) > 0 {
		params.JoinAccounts = expandIntListFromInterface(joinAccounts)
	}
	return params
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
