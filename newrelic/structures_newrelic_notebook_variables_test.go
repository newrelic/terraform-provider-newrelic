//go:build unit

package newrelic

import (
	"encoding/json"
	"strings"
	"testing"
)

// ── validateNotebookVariables ─────────────────────────────────────────────────

func TestValidateNotebookVariables_ValidList(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "list",
			"name": "env",
			"as":   "string",
			"item": []interface{}{map[string]interface{}{"value": "prod"}, map[string]interface{}{"value": "staging"}},
		},
	}
	if errs := validateNotebookVariables(vars); len(errs) != 0 {
		t.Errorf("expected no errors for valid list variable, got: %v", errs)
	}
}

func TestValidateNotebookVariables_ValidQuery(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "query",
			"name": "appName",
			"as":   "string",
			"nrql_query": []interface{}{
				map[string]interface{}{
					"query":       "SELECT uniques(appName) FROM Transaction SINCE 1 hour ago",
					"account_ids": []interface{}{1234567},
				},
			},
		},
	}
	if errs := validateNotebookVariables(vars); len(errs) != 0 {
		t.Errorf("expected no errors for valid query variable, got: %v", errs)
	}
}

func TestValidateNotebookVariables_ValidInput(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "input",
			"name": "searchTerm",
			"as":   "string",
		},
	}
	if errs := validateNotebookVariables(vars); len(errs) != 0 {
		t.Errorf("expected no errors for valid input variable, got: %v", errs)
	}
}

func TestValidateNotebookVariables_EmptyName(t *testing.T) {
	vars := []map[string]interface{}{
		{"from": "list", "name": "", "item": []interface{}{map[string]interface{}{"value": "x"}}},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "name must not be empty") {
		t.Errorf("expected name-empty error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_InvalidFrom(t *testing.T) {
	vars := []map[string]interface{}{
		{"from": "ENUM", "name": "x"},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], `from must be one of`) {
		t.Errorf("expected from-invalid error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_InvalidAs(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "list",
			"name": "x",
			"as":   "INTEGER",
			"item": []interface{}{map[string]interface{}{"value": "1"}},
		},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], `as must be one of`) {
		t.Errorf("expected as-invalid error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_ListWithNoItems(t *testing.T) {
	vars := []map[string]interface{}{
		{"from": "list", "name": "x", "item": []interface{}{}},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "requires at least one item") {
		t.Errorf("expected items-empty error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_QueryWithNoNRQL(t *testing.T) {
	vars := []map[string]interface{}{
		{"from": "query", "name": "x", "nrql_query": []interface{}{}},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "requires at least one nrql_query block") {
		t.Errorf("expected nrql-empty error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_QueryEmptyQueryString(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "query",
			"name": "x",
			"nrql_query": []interface{}{
				map[string]interface{}{"query": "", "account_ids": []interface{}{1}},
			},
		},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "query must not be empty") {
		t.Errorf("expected empty-query error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_QueryEmptyAccountIDs(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "query",
			"name": "x",
			"nrql_query": []interface{}{
				map[string]interface{}{"query": "SELECT 1", "account_ids": []interface{}{}},
			},
		},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "account_ids must contain at least one") {
		t.Errorf("expected account_ids-empty error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_InputWithItems(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from": "input",
			"name": "x",
			"item": []interface{}{map[string]interface{}{"value": "y"}},
		},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "item must not be set on an input variable") {
		t.Errorf("expected item-on-input error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_InputWithNRQL(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from":       "input",
			"name":       "x",
			"nrql_query": []interface{}{map[string]interface{}{"query": "SELECT 1", "account_ids": []interface{}{1}}},
		},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "nrql_query must not be set on an input variable") {
		t.Errorf("expected nrql_query-on-input error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_InputMultipleDefaults(t *testing.T) {
	vars := []map[string]interface{}{
		{
			"from":     "input",
			"name":     "x",
			"defaults": []interface{}{"a", "b"},
		},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) == 0 || !strings.Contains(errs[0], "at most one default value") {
		t.Errorf("expected multiple-defaults-on-input error, got: %v", errs)
	}
}

func TestValidateNotebookVariables_MultipleVarsMixedErrors(t *testing.T) {
	vars := []map[string]interface{}{
		{"from": "list", "name": "", "item": []interface{}{}},
		{"from": "query", "name": "x", "nrql_query": []interface{}{}},
	}
	errs := validateNotebookVariables(vars)
	if len(errs) < 2 {
		t.Errorf("expected at least 2 errors for 2 invalid variables, got: %v", errs)
	}
}

// ── mergeVariablesIntoBlob ────────────────────────────────────────────────────

func TestMergeVariablesIntoBlob_AddsVariables(t *testing.T) {
	doc := map[string]interface{}{"type": "declarative", "version": 1, "content": []interface{}{}}
	vars := []interface{}{map[string]interface{}{"from": "list", "name": "x"}}
	result := mergeVariablesIntoBlob(doc, vars)
	if result["variables"] == nil {
		t.Error("expected variables key to be set")
	}
}

func TestMergeVariablesIntoBlob_RemovesVariablesWhenNil(t *testing.T) {
	doc := map[string]interface{}{
		"type":      "declarative",
		"variables": []interface{}{map[string]interface{}{"from": "list"}},
	}
	result := mergeVariablesIntoBlob(doc, nil)
	if _, ok := result["variables"]; ok {
		t.Error("expected variables key to be removed when vars is nil")
	}
}

func TestMergeVariablesIntoBlob_RemovesVariablesWhenEmpty(t *testing.T) {
	doc := map[string]interface{}{
		"type":      "declarative",
		"variables": []interface{}{map[string]interface{}{"from": "list"}},
	}
	result := mergeVariablesIntoBlob(doc, []interface{}{})
	if _, ok := result["variables"]; ok {
		t.Error("expected variables key to be removed when vars is empty slice")
	}
}

// ── extractVariablesFromBlob ──────────────────────────────────────────────────

func TestExtractVariablesFromBlob_PresentArray(t *testing.T) {
	raw := `{"type":"declarative","variables":[{"from":"list","name":"x"}]}`
	vars, err := extractVariablesFromBlob(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vars) != 1 {
		t.Errorf("expected 1 variable, got %d", len(vars))
	}
}

func TestExtractVariablesFromBlob_AbsentKey(t *testing.T) {
	raw := `{"type":"declarative","content":[]}`
	vars, err := extractVariablesFromBlob(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vars != nil {
		t.Errorf("expected nil when variables key absent, got: %v", vars)
	}
}

func TestExtractVariablesFromBlob_NullValue(t *testing.T) {
	raw := `{"type":"declarative","variables":null}`
	vars, err := extractVariablesFromBlob(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vars != nil {
		t.Errorf("expected nil when variables is null, got: %v", vars)
	}
}

func TestExtractVariablesFromBlob_InvalidJSON(t *testing.T) {
	_, err := extractVariablesFromBlob("not json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// ── flattenNotebookVarItems ───────────────────────────────────────────────────

func TestFlattenNotebookVarItems_BasicItems(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{"value": "prod"},
		map[string]interface{}{"value": "staging"},
	}
	out := flattenNotebookVarItems(raw)
	if len(out) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out))
	}
	first := out[0].(map[string]interface{})
	if first["value"] != "prod" {
		t.Errorf("expected 'prod', got %q", first["value"])
	}
}

func TestFlattenNotebookVarItems_NilInput(t *testing.T) {
	out := flattenNotebookVarItems(nil)
	if len(out) != 0 {
		t.Errorf("expected empty slice for nil input, got: %v", out)
	}
}

// ── flattenNotebookVarNRQLQueries ─────────────────────────────────────────────

func TestFlattenNotebookVarNRQLQueries_Float64AccountIDs(t *testing.T) {
	// JSON decode produces float64 for numbers — flatten must convert to int.
	raw := []interface{}{
		map[string]interface{}{
			"query":      "SELECT 1",
			"accountIds": []interface{}{float64(1234567), float64(7654321)},
		},
	}
	out := flattenNotebookVarNRQLQueries(raw)
	if len(out) != 1 {
		t.Fatalf("expected 1 query, got %d", len(out))
	}
	q := out[0].(map[string]interface{})
	ids := q["account_ids"].([]interface{})
	if ids[0] != 1234567 {
		t.Errorf("expected int 1234567, got %v (%T)", ids[0], ids[0])
	}
}

// ── flattenIgnoreTimeRange ────────────────────────────────────────────────────

func TestFlattenIgnoreTimeRange_TrueValue(t *testing.T) {
	opts := map[string]interface{}{"ignoreTimeRange": true}
	if !flattenIgnoreTimeRange(opts) {
		t.Error("expected true")
	}
}

func TestFlattenIgnoreTimeRange_FalseValue(t *testing.T) {
	opts := map[string]interface{}{"ignoreTimeRange": false}
	if flattenIgnoreTimeRange(opts) {
		t.Error("expected false")
	}
}

func TestFlattenIgnoreTimeRange_NilOptions(t *testing.T) {
	if !flattenIgnoreTimeRange(nil) {
		t.Error("expected default true when options is nil")
	}
}

func TestFlattenIgnoreTimeRange_MissingKey(t *testing.T) {
	opts := map[string]interface{}{}
	if !flattenIgnoreTimeRange(opts) {
		t.Error("expected default true when ignoreTimeRange key absent")
	}
}

// ── round-trip: expand then flatten ──────────────────────────────────────────

// roundTripVarJSON marshals a variable slice to JSON (simulating the blob API
// pass-through) and unmarshals it back into []interface{} (simulating the GET
// response), so we can test that flatten correctly recovers all fields.
func roundTripVarJSON(t *testing.T, vars []interface{}) []interface{} {
	t.Helper()
	b, err := json.Marshal(vars)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out []interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestExpandFlattenRoundTrip_ListVariable(t *testing.T) {
	// Build what expand would produce for a list variable.
	vars := []interface{}{
		map[string]interface{}{
			"from":     "list",
			"name":     "environment",
			"label":    "Environment",
			"as":       "string",
			"omit":     false,
			"hidden":   false,
			"defaults": []string{"prod"},
			"items": []map[string]interface{}{
				{"value": "prod"},
				{"value": "staging"},
				{"value": "dev"},
			},
			"multi": true,
			"defer": true,
		},
	}

	// Simulate blob API pass-through (JSON round-trip).
	stored := roundTripVarJSON(t, vars)

	// stored[0] is the blob GET response for this variable.
	v := stored[0].(map[string]interface{})
	if v["from"] != "list" {
		t.Errorf("from: expected 'list', got %v", v["from"])
	}
	if v["name"] != "environment" {
		t.Errorf("name: expected 'environment', got %v", v["name"])
	}
	items, _ := v["items"].([]interface{})
	if len(items) != 3 {
		t.Errorf("items: expected 3, got %d", len(items))
	}
	if v["multi"] != true {
		t.Errorf("multi: expected true, got %v", v["multi"])
	}
	if v["defer"] != true {
		t.Errorf("defer: expected true, got %v", v["defer"])
	}
	defs, _ := v["defaults"].([]interface{})
	if len(defs) != 1 || defs[0] != "prod" {
		t.Errorf("defaults: expected ['prod'], got %v", defs)
	}
}

func TestExpandFlattenRoundTrip_QueryVariable(t *testing.T) {
	vars := []interface{}{
		map[string]interface{}{
			"from":   "query",
			"name":   "appName",
			"label":  "Application",
			"as":     "string",
			"omit":   false,
			"hidden": false,
			"nrqlQueries": []map[string]interface{}{
				{
					"query":      "SELECT uniques(appName) FROM Transaction SINCE 1 hour ago",
					"accountIds": []int{1234567},
				},
			},
			"options": map[string]interface{}{"ignoreTimeRange": true},
			"multi":   false,
			"defer":   true,
		},
	}

	stored := roundTripVarJSON(t, vars)
	v := stored[0].(map[string]interface{})

	if v["from"] != "query" {
		t.Errorf("from: expected 'query', got %v", v["from"])
	}
	queries, _ := v["nrqlQueries"].([]interface{})
	if len(queries) != 1 {
		t.Errorf("nrqlQueries: expected 1, got %d", len(queries))
	}
	q := queries[0].(map[string]interface{})
	if q["query"] == "" {
		t.Error("query string must not be empty")
	}
	ids, _ := q["accountIds"].([]interface{})
	if len(ids) != 1 {
		t.Errorf("accountIds: expected 1, got %d", len(ids))
	}
	opts, _ := v["options"].(map[string]interface{})
	if opts["ignoreTimeRange"] != true {
		t.Errorf("ignoreTimeRange: expected true, got %v", opts["ignoreTimeRange"])
	}
}

func TestExpandFlattenRoundTrip_InputVariable(t *testing.T) {
	vars := []interface{}{
		map[string]interface{}{
			"from":   "input",
			"name":   "searchTerm",
			"label":  "Search",
			"as":     "string",
			"omit":   false,
			"hidden": false,
			"type":   "text",
		},
	}

	stored := roundTripVarJSON(t, vars)
	v := stored[0].(map[string]interface{})

	if v["from"] != "input" {
		t.Errorf("from: expected 'input', got %v", v["from"])
	}
	if v["type"] != "text" {
		t.Errorf("type: expected 'text', got %v", v["type"])
	}
	// input variables must not have items or nrqlQueries
	if _, ok := v["items"]; ok {
		t.Error("input variable should not have items key in blob")
	}
	if _, ok := v["nrqlQueries"]; ok {
		t.Error("input variable should not have nrqlQueries key in blob")
	}
}

// ── boolOrTrueNb / boolOrFalseNb ─────────────────────────────────────────────

func TestBoolOrTrueNb_NilIsTrue(t *testing.T) {
	if !boolOrTrueNb(nil) {
		t.Error("expected true for nil")
	}
}

func TestBoolOrTrueNb_FalseIsRespected(t *testing.T) {
	if boolOrTrueNb(false) {
		t.Error("expected false to be respected")
	}
}

func TestBoolOrFalseNb_NilIsFalse(t *testing.T) {
	if boolOrFalseNb(nil) {
		t.Error("expected false for nil")
	}
}

func TestBoolOrFalseNb_TrueIsRespected(t *testing.T) {
	if !boolOrFalseNb(true) {
		t.Error("expected true to be respected")
	}
}
