//go:build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── expandIntListFromInterface ────────────────────────────────────────────────

func TestExpandIntListFromInterface(t *testing.T) {
	t.Parallel()
	raw := []interface{}{1, 2, 3}
	assert.Equal(t, []int{1, 2, 3}, expandIntListFromInterface(raw))
}

func TestExpandIntListFromInterface_Empty(t *testing.T) {
	t.Parallel()
	assert.Empty(t, expandIntListFromInterface(nil))
	assert.Empty(t, expandIntListFromInterface([]interface{}{}))
}

// ── expandNRQLEngineCreate ────────────────────────────────────────────────────

func TestExpandNRQLEngineCreate_Basic(t *testing.T) {
	t.Parallel()
	raw := []interface{}{map[string]interface{}{
		"query":         "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id",
		"accounts":      []interface{}{12345},
		"join_accounts": []interface{}{},
	}}
	result := expandNRQLEngineCreate(raw)
	require.NotNil(t, result)
	assert.Equal(t, "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id", result.Query)
	assert.Equal(t, []int{12345}, result.Accounts)
	assert.Empty(t, result.JoinAccounts)
}

func TestExpandNRQLEngineCreate_WithJoinAccounts(t *testing.T) {
	t.Parallel()
	raw := []interface{}{map[string]interface{}{
		"query":         "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id",
		"accounts":      []interface{}{100},
		"join_accounts": []interface{}{200, 300},
	}}
	result := expandNRQLEngineCreate(raw)
	require.NotNil(t, result)
	assert.Equal(t, []int{100}, result.Accounts)
	assert.Equal(t, []int{200, 300}, result.JoinAccounts)
}

func TestExpandNRQLEngineCreate_Nil(t *testing.T) {
	t.Parallel()
	assert.Nil(t, expandNRQLEngineCreate(nil))
	assert.Nil(t, expandNRQLEngineCreate([]interface{}{nil}))
}

// ── expandNRQLEngineUpdate ────────────────────────────────────────────────────

func TestExpandNRQLEngineUpdate_MatchesCreateShape(t *testing.T) {
	t.Parallel()
	raw := []interface{}{map[string]interface{}{
		"query":         "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id",
		"accounts":      []interface{}{42},
		"join_accounts": []interface{}{},
	}}
	result := expandNRQLEngineUpdate(raw)
	require.NotNil(t, result)
	assert.Equal(t, "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id", result.Query)
	assert.Equal(t, []int{42}, result.Accounts)
}

// ── flattenNRQLEngine ─────────────────────────────────────────────────────────

func TestFlattenNRQLEngine_RoundTrip(t *testing.T) {
	t.Parallel()
	engine := scorecards.EntityManagementNRQLRuleEngine{
		Query:        "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id",
		Accounts:     []int{10, 20},
		JoinAccounts: []int{30},
	}
	flat := flattenNRQLEngine(engine)
	require.Len(t, flat, 1)
	assert.Equal(t, engine.Query, flat[0]["query"])
	assert.Equal(t, []int{10, 20}, flat[0]["accounts"])
	assert.Equal(t, []int{30}, flat[0]["join_accounts"])
}

func TestFlattenNRQLEngine_EmptyJoinAccounts(t *testing.T) {
	t.Parallel()
	engine := scorecards.EntityManagementNRQLRuleEngine{
		Query:    "SELECT if(1=1, 1, 0) AS 'score' FROM Transaction FACET id",
		Accounts: []int{1},
	}
	flat := flattenNRQLEngine(engine)
	require.Len(t, flat, 1)
	assert.Empty(t, flat[0]["join_accounts"])
}

// ── expandProgressLevels ──────────────────────────────────────────────────────

func TestExpandProgressLevels_Basic(t *testing.T) {
	t.Parallel()
	raw := []interface{}{
		map[string]interface{}{"id": "red", "name": "Red", "description": "Needs work", "hex_color_code": "#FF0000"},
		map[string]interface{}{"id": "green", "name": "Green", "description": "Good", "hex_color_code": "#00CC00"},
	}
	levels := expandProgressLevels(raw)
	require.Len(t, levels, 2)
	assert.Equal(t, "red", levels[0].ID)
	assert.Equal(t, "Red", levels[0].Name)
	assert.Equal(t, "#FF0000", levels[0].HexColorCode)
	assert.Equal(t, "green", levels[1].ID)
}

func TestExpandProgressLevels_Empty(t *testing.T) {
	t.Parallel()
	assert.Nil(t, expandProgressLevels(nil))
	assert.Nil(t, expandProgressLevels([]interface{}{}))
}

// ── flattenProgressLevels ─────────────────────────────────────────────────────

func TestFlattenProgressLevels_RoundTrip(t *testing.T) {
	t.Parallel()
	levels := []scorecards.EntityManagementProgressLevelDefinition{
		{ID: "red", Name: "Red", Description: "Needs work", HexColorCode: "#FF0000"},
		{ID: "green", Name: "Green", Description: "Good", HexColorCode: "#00CC00"},
	}
	flat := flattenProgressLevels(levels)
	require.Len(t, flat, 2)
	assert.Equal(t, "red", flat[0]["id"])
	assert.Equal(t, "Red", flat[0]["name"])
	assert.Equal(t, "#FF0000", flat[0]["hex_color_code"])
}

// ── expandRuleIDsFromSet ──────────────────────────────────────────────────────

func TestExpandRuleIDsFromSet(t *testing.T) {
	t.Parallel()
	s := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), []interface{}{"rule-id-1", "rule-id-2"})
	ids := expandRuleIDsFromSet(s)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, "rule-id-1")
	assert.Contains(t, ids, "rule-id-2")
}

func TestExpandRuleIDsFromSet_Empty(t *testing.T) {
	t.Parallel()
	s := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), []interface{}{})
	assert.Empty(t, expandRuleIDsFromSet(s))
}
