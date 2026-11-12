//go:build unit

package newrelic

import (
	"encoding/json"
	"strings"
	"testing"
)

// notebook wraps a content string in a valid declarative envelope so tests
// don't need to repeat the outer boilerplate.
func notebookWithWidgets(widgetsJSON string) string {
	return `{"type":"declarative","version":1,"content":[{"type":"container","props":{"layout":"stack"},"content":` + widgetsJSON + `}]}`
}

// ── type field ────────────────────────────────────────────────────────────────

func TestCheckWidgetLevelProps_TypeNotWidget(t *testing.T) {
	raw := notebookWithWidgets(`[{"type":"container","content":{"type":"visualization","id":"viz.markdown"}}]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], `type must be "widget"`) {
		t.Errorf("unexpected error message: %s", errs[0])
	}
	if !strings.Contains(errs[0], `"container"`) {
		t.Errorf("error should include the bad type value: %s", errs[0])
	}
}

func TestCheckWidgetLevelProps_TypeEmpty(t *testing.T) {
	raw := notebookWithWidgets(`[{"content":{"type":"visualization","id":"viz.markdown"}}]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], `type must be "widget"`) {
		t.Errorf("unexpected error message: %s", errs[0])
	}
}

func TestCheckWidgetLevelProps_TypeErrorSkipsDeepChecks(t *testing.T) {
	// A non-widget item with no content block and no props should produce exactly
	// one error (the type error), not also a content-missing or props error.
	raw := notebookWithWidgets(`[{"type":"unknown"}]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected exactly 1 error (type only), got %d: %v", len(errs), errs)
	}
}

// ── content block ─────────────────────────────────────────────────────────────

func TestCheckWidgetLevelProps_MissingContentBlock(t *testing.T) {
	raw := notebookWithWidgets(`[{"type":"widget"}]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], `missing required "content" block`) {
		t.Errorf("unexpected error message: %s", errs[0])
	}
}

func TestCheckWidgetLevelProps_ContentBlockPresentButNoID(t *testing.T) {
	raw := notebookWithWidgets(`[{"type":"widget","content":{"type":"visualization"}}]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], `missing required "content" block`) {
		t.Errorf("unexpected error message: %s", errs[0])
	}
}

func TestCheckWidgetLevelProps_ContentErrorSkipsPropsCheck(t *testing.T) {
	// A widget with a valid type but no content block should produce exactly
	// one error, not also trigger the props check.
	raw := notebookWithWidgets(`[{"type":"widget","props":{"title":"x"}}]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected exactly 1 error (content only), got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], `missing required "content" block`) {
		t.Errorf("unexpected error: %s", errs[0])
	}
}

// ── path format ───────────────────────────────────────────────────────────────

func TestCheckWidgetLevelProps_PathReflectsPosition(t *testing.T) {
	// Two widgets: first is valid (viz.markdown), second has wrong type.
	// The error path should reference content[0].content[1].
	raw := notebookWithWidgets(`[
		{"type":"widget","content":{"type":"visualization","id":"viz.markdown"}},
		{"type":"bad"}
	]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], "content[0].content[1]") {
		t.Errorf("error path should be content[0].content[1], got: %s", errs[0])
	}
}

// ── valid documents ───────────────────────────────────────────────────────────

func TestCheckWidgetLevelProps_ValidMarkdown(t *testing.T) {
	raw := notebookWithWidgets(`[{"type":"widget","content":{"type":"visualization","id":"viz.markdown","props":{"text":"hi"}}}]`)
	if errs := checkWidgetLevelProps(raw); len(errs) != 0 {
		t.Errorf("expected no errors for valid markdown widget, got: %v", errs)
	}
}

func TestCheckWidgetLevelProps_ValidNonMarkdown(t *testing.T) {
	raw := notebookWithWidgets(`[{"type":"widget","props":{"title":""},"content":{"type":"visualization","id":"viz.billboard","props":{}}}]`)
	if errs := checkWidgetLevelProps(raw); len(errs) != 0 {
		t.Errorf("expected no errors for valid billboard widget, got: %v", errs)
	}
}

func TestCheckWidgetLevelProps_MultipleContainersAllValid(t *testing.T) {
	raw := `{"type":"declarative","version":1,"content":[
		{"type":"container","content":[{"type":"widget","content":{"type":"visualization","id":"viz.markdown","props":{"text":"a"}}}]},
		{"type":"container","content":[{"type":"widget","props":{"title":""},"content":{"type":"visualization","id":"viz.line","props":{}}}]}
	]}`
	if errs := checkWidgetLevelProps(raw); len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestCheckWidgetLevelProps_MultipleErrors(t *testing.T) {
	// Two widgets: first has wrong type, second is missing content block.
	raw := notebookWithWidgets(`[
		{"type":"unknown"},
		{"type":"widget"}
	]`)
	errs := checkWidgetLevelProps(raw)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestCheckWidgetLevelProps_InvalidJSONReturnNil(t *testing.T) {
	if errs := checkWidgetLevelProps("not json"); errs != nil {
		t.Errorf("invalid JSON should return nil (envelope validator catches it), got: %v", errs)
	}
}

// ── checkEnvelopeErrors ───────────────────────────────────────────────────────

func TestCheckEnvelopeErrors_ValidEnvelope(t *testing.T) {
	errs := checkEnvelopeErrors(`{"type":"declarative","version":1,"content":[]}`)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid envelope, got: %v", errs)
	}
}

func TestCheckEnvelopeErrors_EmptyString(t *testing.T) {
	errs := checkEnvelopeErrors("")
	if len(errs) == 0 {
		t.Error("expected error for empty string")
	}
}

func TestCheckEnvelopeErrors_InvalidJSON(t *testing.T) {
	errs := checkEnvelopeErrors("{not json}")
	if len(errs) == 0 || !strings.Contains(errs[0], "JSON syntax error") {
		t.Errorf("expected JSON parse error, got: %v", errs)
	}
}

func TestCheckEnvelopeErrors_MissingType(t *testing.T) {
	errs := checkEnvelopeErrors(`{"version":1,"content":[]}`)
	if len(errs) == 0 || !strings.Contains(errs[0], `missing required field "type"`) {
		t.Errorf("expected missing-type error, got: %v", errs)
	}
}

func TestCheckEnvelopeErrors_WrongType(t *testing.T) {
	errs := checkEnvelopeErrors(`{"type":"legacy","version":1,"content":[]}`)
	if len(errs) == 0 || !strings.Contains(errs[0], `must be "declarative"`) {
		t.Errorf("expected wrong-type error, got: %v", errs)
	}
}

func TestCheckEnvelopeErrors_MissingVersion(t *testing.T) {
	errs := checkEnvelopeErrors(`{"type":"declarative","content":[]}`)
	if len(errs) == 0 || !strings.Contains(errs[0], `missing required field "version"`) {
		t.Errorf("expected missing-version error, got: %v", errs)
	}
}

func TestCheckEnvelopeErrors_VersionAsString(t *testing.T) {
	errs := checkEnvelopeErrors(`{"type":"declarative","version":"1","content":[]}`)
	if len(errs) == 0 || !strings.Contains(errs[0], `"version" must be the integer 1`) {
		t.Errorf("expected version-type error, got: %v", errs)
	}
}

func TestCheckEnvelopeErrors_VersionWrongInteger(t *testing.T) {
	errs := checkEnvelopeErrors(`{"type":"declarative","version":2,"content":[]}`)
	if len(errs) == 0 {
		t.Error("expected error for version 2")
	}
}

func TestCheckEnvelopeErrors_ContentNotArray(t *testing.T) {
	errs := checkEnvelopeErrors(`{"type":"declarative","version":1,"content":{}}`)
	if len(errs) == 0 || !strings.Contains(errs[0], `"content" must be a JSON array`) {
		t.Errorf("expected content-not-array error, got: %v", errs)
	}
}

// ── normalizeNotebookContent ──────────────────────────────────────────────────

func TestNormalizeNotebookContent_EmptyReturnsError(t *testing.T) {
	_, _, err := normalizeNotebookContent("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestNormalizeNotebookContent_KeysSorted(t *testing.T) {
	a := `{"version":1,"type":"declarative","content":[]}`
	b := `{"content":[],"type":"declarative","version":1}`
	normA, _, _ := normalizeNotebookContent(a)
	normB, _, _ := normalizeNotebookContent(b)
	if normA != normB {
		t.Errorf("key order should not affect normalized form:\n  a=%s\n  b=%s", normA, normB)
	}
}

func TestNormalizeNotebookContent_NoTrailingNewline(t *testing.T) {
	norm, _, _ := normalizeNotebookContent(`{"type":"declarative","version":1,"content":[]}`)
	if strings.HasSuffix(norm, "\n") {
		t.Error("normalized output must not have trailing newline")
	}
}

func TestNormalizeNotebookContent_HTMLNotEscaped(t *testing.T) {
	raw := `{"type":"declarative","version":1,"content":[{"text":"<b>bold</b>"}]}`
	norm, _, err := normalizeNotebookContent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(norm, "<b>bold</b>") {
		t.Error("HTML characters should not be escaped in normalized output")
	}
}

func TestNormalizeNotebookContent_ReturnsParsedBody(t *testing.T) {
	_, body, err := normalizeNotebookContent(`{"type":"declarative","version":1,"content":[]}`)
	if err != nil || body == nil {
		t.Fatalf("expected non-nil body, err: %v", err)
	}
	m, ok := body.(map[string]interface{})
	if !ok || m["type"] != "declarative" {
		t.Errorf("unexpected body: %v", body)
	}
}

// ── suppressEquivalentNotebookContent ─────────────────────────────────────────

func TestSuppressEquivalentContent_IdenticalStrings(t *testing.T) {
	v := `{"type":"declarative","version":1,"content":[]}`
	if !suppressEquivalentNotebookContent("", v, v, nil) {
		t.Error("identical strings should suppress diff")
	}
}

func TestSuppressEquivalentContent_ReorderedKeys(t *testing.T) {
	a := `{"type":"declarative","version":1,"content":[]}`
	b := `{"content":[],"version":1,"type":"declarative"}`
	if !suppressEquivalentNotebookContent("", a, b, nil) {
		t.Error("reordered keys with same values should suppress diff")
	}
}

func TestSuppressEquivalentContent_DifferentValues(t *testing.T) {
	a := `{"type":"declarative","version":1,"content":[]}`
	b := `{"type":"declarative","version":1,"content":[{}]}`
	if suppressEquivalentNotebookContent("", a, b, nil) {
		t.Error("documents with different values should not suppress diff")
	}
}

func TestSuppressEquivalentContent_InvalidJSON(t *testing.T) {
	if suppressEquivalentNotebookContent("", "not json", "{}", nil) {
		t.Error("invalid JSON should not suppress diff")
	}
}

// ── checkNotebookVariablesJSON ────────────────────────────────────────────────

func TestCheckNotebookVariablesJSON_Empty(t *testing.T) {
	if errs := checkNotebookVariablesJSON(""); len(errs) != 0 {
		t.Errorf("empty string should produce no errors, got %v", errs)
	}
}

func TestCheckNotebookVariablesJSON_ValidList(t *testing.T) {
	raw := `[{"name":"env","from":"list","items":[{"value":"prod"}]}]`
	if errs := checkNotebookVariablesJSON(raw); len(errs) != 0 {
		t.Errorf("valid list variable should produce no errors, got %v", errs)
	}
}

func TestCheckNotebookVariablesJSON_ValidQuery(t *testing.T) {
	raw := `[{"name":"app","from":"query","nrqlQueries":[{"query":"SELECT uniques(appName) FROM Transaction","accountIds":[1234567]}]}]`
	if errs := checkNotebookVariablesJSON(raw); len(errs) != 0 {
		t.Errorf("valid query variable should produce no errors, got %v", errs)
	}
}

func TestCheckNotebookVariablesJSON_ValidInput(t *testing.T) {
	raw := `[{"name":"searchTerm","from":"input"}]`
	if errs := checkNotebookVariablesJSON(raw); len(errs) != 0 {
		t.Errorf("valid input variable should produce no errors, got %v", errs)
	}
}

func TestCheckNotebookVariablesJSON_ValidEmptyArray(t *testing.T) {
	if errs := checkNotebookVariablesJSON("[]"); len(errs) != 0 {
		t.Errorf("empty array should produce no errors, got %v", errs)
	}
}

func TestCheckNotebookVariablesJSON_AllThreeTypes(t *testing.T) {
	raw := `[
		{"name":"env","from":"list","items":[{"value":"prod"}]},
		{"name":"app","from":"query","nrqlQueries":[{"query":"SELECT 1","accountIds":[1]}]},
		{"name":"filter","from":"input"}
	]`
	if errs := checkNotebookVariablesJSON(raw); len(errs) != 0 {
		t.Errorf("all three variable types should produce no errors, got %v", errs)
	}
}

func TestCheckNotebookVariablesJSON_InvalidJSON(t *testing.T) {
	errs := checkNotebookVariablesJSON("not json")
	if len(errs) == 0 {
		t.Error("invalid JSON should produce an error")
	}
	if !strings.Contains(errs[0], "not a valid JSON array") {
		t.Errorf("expected 'not a valid JSON array' in error, got %q", errs[0])
	}
}

func TestCheckNotebookVariablesJSON_NotArray(t *testing.T) {
	errs := checkNotebookVariablesJSON(`{"name":"env","from":"list"}`)
	if len(errs) == 0 {
		t.Error("JSON object (not array) should produce an error")
	}
}

func TestCheckNotebookVariablesJSON_MissingName(t *testing.T) {
	errs := checkNotebookVariablesJSON(`[{"from":"list"}]`)
	if len(errs) == 0 {
		t.Error("missing name should produce an error")
	}
	if !strings.Contains(errs[0], `"name" is required`) {
		t.Errorf("expected name-required error, got %q", errs[0])
	}
}

func TestCheckNotebookVariablesJSON_MissingFrom(t *testing.T) {
	errs := checkNotebookVariablesJSON(`[{"name":"env"}]`)
	if len(errs) == 0 {
		t.Error("missing from should produce an error")
	}
	if !strings.Contains(errs[0], `"from" is required`) {
		t.Errorf("expected from-required error, got %q", errs[0])
	}
}

func TestCheckNotebookVariablesJSON_InvalidFrom(t *testing.T) {
	errs := checkNotebookVariablesJSON(`[{"name":"env","from":"unknown"}]`)
	if len(errs) == 0 {
		t.Error("invalid from value should produce an error")
	}
	if !strings.Contains(errs[0], `"from" must be one of`) {
		t.Errorf("expected from-must-be-one-of error, got %q", errs[0])
	}
}

func TestCheckNotebookVariablesJSON_EntryNotObject(t *testing.T) {
	errs := checkNotebookVariablesJSON(`["not-an-object"]`)
	if len(errs) == 0 {
		t.Error("non-object array entry should produce an error")
	}
	if !strings.Contains(errs[0], "must be a JSON object") {
		t.Errorf("expected object-required error, got %q", errs[0])
	}
}

func TestCheckNotebookVariablesJSON_MultipleErrors(t *testing.T) {
	raw := `[{"from":"list"},{"name":"app","from":"bad"}]`
	errs := checkNotebookVariablesJSON(raw)
	if len(errs) < 2 {
		t.Errorf("expected at least 2 errors for two invalid entries, got %d: %v", len(errs), errs)
	}
}

// ── mergeVariablesIntoBlob / extractVariablesFromBlob ─────────────────────────

func TestMergeVariablesIntoBlob_SetsKey(t *testing.T) {
	doc := map[string]interface{}{"type": "declarative", "version": 1.0, "content": []interface{}{}}
	vars := []interface{}{map[string]interface{}{"name": "env", "from": "list"}}
	result := mergeVariablesIntoBlob(doc, vars)
	v, ok := result["variables"]
	if !ok {
		t.Fatal("expected 'variables' key in merged doc")
	}
	arr, ok := v.([]interface{})
	if !ok || len(arr) != 1 {
		t.Errorf("expected 1-entry variables array, got %v", v)
	}
}

func TestMergeVariablesIntoBlob_NilVarsDeletesKey(t *testing.T) {
	doc := map[string]interface{}{
		"type":      "declarative",
		"variables": []interface{}{map[string]interface{}{"name": "old", "from": "list"}},
	}
	result := mergeVariablesIntoBlob(doc, nil)
	if _, ok := result["variables"]; ok {
		t.Error("nil vars should remove 'variables' key from blob")
	}
}

func TestMergeVariablesIntoBlob_EmptySliceDeletesKey(t *testing.T) {
	doc := map[string]interface{}{"type": "declarative", "variables": []interface{}{}}
	result := mergeVariablesIntoBlob(doc, []interface{}{})
	if _, ok := result["variables"]; ok {
		t.Error("empty vars should remove 'variables' key from blob")
	}
}

func TestExtractVariablesFromBlob_Present(t *testing.T) {
	raw := `{"type":"declarative","version":1,"content":[],"variables":[{"name":"env","from":"list"}]}`
	vars, err := extractVariablesFromBlob(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vars) != 1 {
		t.Errorf("expected 1 variable, got %d", len(vars))
	}
	m := vars[0].(map[string]interface{})
	if m["name"] != "env" {
		t.Errorf("expected name=env, got %v", m["name"])
	}
}

func TestExtractVariablesFromBlob_Absent(t *testing.T) {
	raw := `{"type":"declarative","version":1,"content":[]}`
	vars, err := extractVariablesFromBlob(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vars) != 0 {
		t.Errorf("expected 0 variables when key absent, got %d", len(vars))
	}
}

func TestExtractVariablesFromBlob_InvalidJSON(t *testing.T) {
	_, err := extractVariablesFromBlob("not json")
	if err == nil {
		t.Error("invalid JSON should return error")
	}
}

// ── suppressEquivalentNotebookVariables ───────────────────────────────────────

func TestSuppressEquivalentVariables_BothEmpty(t *testing.T) {
	if !suppressEquivalentNotebookVariables("", "", "", nil) {
		t.Error("two empty strings should suppress diff")
	}
}

func TestSuppressEquivalentVariables_NullAndEmpty(t *testing.T) {
	if !suppressEquivalentNotebookVariables("", "null", "[]", nil) {
		t.Error("null vs [] should suppress diff")
	}
}

func TestSuppressEquivalentVariables_EmptyArrayBothWays(t *testing.T) {
	if !suppressEquivalentNotebookVariables("", "[]", "", nil) {
		t.Error("[] vs empty string should suppress diff")
	}
}

func TestSuppressEquivalentVariables_FormattingOnly(t *testing.T) {
	compact := `[{"from":"list","name":"env"}]`
	pretty, _ := json.MarshalIndent([]interface{}{map[string]interface{}{"from": "list", "name": "env"}}, "", "  ")
	if !suppressEquivalentNotebookVariables("", compact, string(pretty), nil) {
		t.Error("formatting-only difference should suppress diff")
	}
}

func TestSuppressEquivalentVariables_DifferentContent(t *testing.T) {
	a := `[{"name":"env","from":"list"}]`
	b := `[{"name":"env","from":"list"},{"name":"region","from":"input"}]`
	if suppressEquivalentNotebookVariables("", a, b, nil) {
		t.Error("different variable counts should not suppress diff")
	}
}

func TestSuppressEquivalentVariables_InvalidJSON(t *testing.T) {
	if suppressEquivalentNotebookVariables("", "not json", "[]", nil) {
		t.Error("invalid JSON should not suppress diff")
	}
}

// ── flattenNotebookContent strips variables ───────────────────────────────────

func TestFlattenNotebookContent_StripsVariables(t *testing.T) {
	// flattenNotebookContent should not store variables inside the content field.
	raw := []byte(`{"type":"declarative","version":1,"content":[],"variables":[{"name":"env","from":"list"}]}`)
	d := resourceNewRelicNotebook().TestResourceData()
	_ = d.Set("content", `{"type":"declarative","version":1,"content":[]}`)

	if err := flattenNotebookContent(raw, d, "content"); err != nil {
		t.Fatalf("flattenNotebookContent error: %v", err)
	}

	stored := d.Get("content").(string)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(stored), &m); err != nil {
		t.Fatalf("stored content is not valid JSON: %v", err)
	}
	if _, ok := m["variables"]; ok {
		t.Error("flattenNotebookContent must strip 'variables' key from the content field")
	}
}
