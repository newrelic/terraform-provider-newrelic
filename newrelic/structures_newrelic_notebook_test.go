//go:build unit

package newrelic

import (
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

// ── validateNotebookContent ───────────────────────────────────────────────────

func TestValidateNotebookContent_ValidEnvelope(t *testing.T) {
	_, errs := validateNotebookContent(`{"type":"declarative","version":1,"content":[]}`, "content")
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid envelope, got: %v", errs)
	}
}

func TestValidateNotebookContent_EmptyString(t *testing.T) {
	_, errs := validateNotebookContent("", "content")
	if len(errs) == 0 {
		t.Error("expected error for empty string")
	}
}

func TestValidateNotebookContent_InvalidJSON(t *testing.T) {
	_, errs := validateNotebookContent("{not json}", "content")
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), "not valid JSON") {
		t.Errorf("expected JSON parse error, got: %v", errs)
	}
}

func TestValidateNotebookContent_MissingType(t *testing.T) {
	_, errs := validateNotebookContent(`{"version":1,"content":[]}`, "content")
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), `missing required field "type"`) {
		t.Errorf("expected missing-type error, got: %v", errs)
	}
}

func TestValidateNotebookContent_WrongType(t *testing.T) {
	_, errs := validateNotebookContent(`{"type":"legacy","version":1,"content":[]}`, "content")
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), `must be "declarative"`) {
		t.Errorf("expected wrong-type error, got: %v", errs)
	}
}

func TestValidateNotebookContent_MissingVersion(t *testing.T) {
	_, errs := validateNotebookContent(`{"type":"declarative","content":[]}`, "content")
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), `missing required field "version"`) {
		t.Errorf("expected missing-version error, got: %v", errs)
	}
}

func TestValidateNotebookContent_VersionAsString(t *testing.T) {
	_, errs := validateNotebookContent(`{"type":"declarative","version":"1","content":[]}`, "content")
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), `"version" must be the integer 1`) {
		t.Errorf("expected version-type error, got: %v", errs)
	}
}

func TestValidateNotebookContent_VersionWrongInteger(t *testing.T) {
	_, errs := validateNotebookContent(`{"type":"declarative","version":2,"content":[]}`, "content")
	if len(errs) == 0 {
		t.Error("expected error for version 2")
	}
}

func TestValidateNotebookContent_ContentNotArray(t *testing.T) {
	_, errs := validateNotebookContent(`{"type":"declarative","version":1,"content":{}}`, "content")
	if len(errs) == 0 || !strings.Contains(errs[0].Error(), `"content" must be an array`) {
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
