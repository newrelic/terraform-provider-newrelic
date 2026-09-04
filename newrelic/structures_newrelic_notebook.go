package newrelic

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// normalizeNotebookContent converts any valid JSON string to a canonical form
// with alphabetically sorted keys and consistent 2-space indentation. Storing
// and comparing this canonical form means that two documents with the same
// semantic content but different formatting are always treated as equal,
// preventing spurious plan diffs when a user reformats their HCL or JSON file.
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
		return "", fmt.Errorf("could not re-serialize JSON: %w", err)
	}

	// json.Encoder appends a trailing newline; remove it so comparisons are stable.
	return string(bytes.TrimRight(buf.Bytes(), "\n")), nil
}

// suppressEquivalentNotebookContent tells Terraform to ignore the difference
// between two JSON strings that are semantically identical. This allows users
// to freely reformat their content or content_json values (for example,
// reordering keys or changing indentation) without triggering a planned update.
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

// flattenNotebookContent stores the content received from the API into the
// correct state field ("content" or "content_json"). It normalizes the value
// first so that subsequent plan operations only highlight lines that genuinely
// changed, not formatting differences introduced during the round-trip.
func flattenNotebookContent(raw json.RawMessage, d *schema.ResourceData, field string) error {
	if len(raw) == 0 {
		return nil
	}
	normalized, err := normalizeNotebookContent(string(raw))
	if err != nil {
		return fmt.Errorf("could not normalize notebook content returned by the API: %w", err)
	}
	return d.Set(field, normalized)
}
