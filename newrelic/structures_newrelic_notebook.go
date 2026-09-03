package newrelic

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// normalizeNotebookContent canonicalises a JSON string so that identical
// documents with different formatting or key ordering produce the same bytes.
// Go's json.Unmarshal+json.MarshalIndent pipeline sorts map keys alphabetically
// on re-serialisation, which is all we need for deterministic state storage and
// precise line-level plan diffs.
func normalizeNotebookContent(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}

	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return "", fmt.Errorf("re-serialise JSON: %w", err)
	}

	// Encoder appends a trailing newline; trim so state comparisons are stable.
	return string(bytes.TrimRight(buf.Bytes(), "\n")), nil
}

// suppressEquivalentNotebookContent is a DiffSuppressFunc that treats two JSON
// strings as equal when they are semantically identical regardless of
// whitespace, indentation, or key ordering. Applied to both the "content" and
// "content_json" fields.
//
// This prevents cosmetic-only edits — e.g. reformatting a jsonencode({}) block
// or re-indenting a JSON file — from showing up as changes in terraform plan.
func suppressEquivalentNotebookContent(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	if oldVal == newVal {
		return true
	}
	normOld, err := normalizeNotebookContent(oldVal)
	if err != nil {
		return false
	}
	normNew, err := normalizeNotebookContent(newVal)
	if err != nil {
		return false
	}
	return normOld == normNew
}

// flattenNotebookContent writes normalised content from the API response into
// state under the given field name ("content" or "content_json"). Calling this
// on every Read ensures state always holds canonical JSON so Terraform's native
// diff highlights only lines that actually changed.
func flattenNotebookContent(raw json.RawMessage, d *schema.ResourceData, field string) error {
	if len(raw) == 0 {
		return nil
	}
	normalized, err := normalizeNotebookContent(string(raw))
	if err != nil {
		return fmt.Errorf("normalize notebook content from API: %w", err)
	}
	return d.Set(field, normalized)
}
