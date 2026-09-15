package newrelic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ── Validation ────────────────────────────────────────────────────────────────

// validateNotebookContent is the ValidateFunc for both content and content.
// It enforces that the value is non-empty, valid JSON, and matches the minimum
// declarative UI envelope: { "type": "declarative", "version": 1, "content": [...] }.
func validateNotebookContent(v interface{}, k string) (warnings []string, errors []error) {
	raw, ok := v.(string)
	if !ok || raw == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		errors = append(errors, fmt.Errorf("%q is not valid JSON: %w", k, err))
		return
	}

	t, hasType := doc["type"]
	if !hasType {
		errors = append(errors, fmt.Errorf(`%q: missing required field "type" (expected "declarative")`, k))
	} else if t != "declarative" {
		errors = append(errors, fmt.Errorf(`%q: "type" must be "declarative", got %q`, k, t))
	}

	if _, hasContent := doc["content"]; !hasContent {
		errors = append(errors, fmt.Errorf(`%q: missing required field "content" (expected an array of container blocks)`, k))
	} else if _, isArray := doc["content"].([]interface{}); !isArray {
		errors = append(errors, fmt.Errorf(`%q: "content" must be an array`, k))
	}

	if ver, hasVersion := doc["version"]; !hasVersion {
		errors = append(errors, fmt.Errorf(`%q: missing required field "version" (expected 1)`, k))
	} else {
		// json.Unmarshal decodes numbers as float64; accept integer 1 only.
		if f, ok := ver.(float64); !ok || f != 1 {
			errors = append(errors, fmt.Errorf(`%q: "version" must be the integer 1, got %#v (%T)`, k, ver, ver))
		}
	}

	return
}

// customizeNotebookDiff enforces widget-level props rules at plan time.
// Using CustomizeDiff (rather than ValidateFunc) means Terraform does not echo
// the full JSON body in the error output when a violation is found.
func customizeNotebookDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if !d.HasChange("content") {
		return nil
	}

	raw := d.Get("content").(string)
	if raw == "" {
		return nil
	}

	errs := checkWidgetLevelProps(raw)
	if len(errs) == 0 {
		return nil
	}

	msg := "the following widget props violations were found:\n"
	for i, e := range errs {
		msg += fmt.Sprintf("(%d) %s\n", i+1, e)
	}
	return fmt.Errorf("%s", strings.TrimRight(msg, "\n"))
}

// ── Widget-level props typed tree ─────────────────────────────────────────────

// nbDoc, nbContainer, nbWidget, nbWidgetProps, nbVizContent are lightweight
// typed structs for walking the notebook document tree. A single json.Unmarshal
// call populates the full tree with no per-field type assertions needed.
// Unrecognised JSON keys are preserved untouched.

type nbDoc struct {
	Content []nbContainer `json:"content"`
}

type nbContainer struct {
	Content []nbWidget `json:"content"`
}

type nbWidget struct {
	Type string `json:"type"`
	// Props pointer: nil means the key is absent; non-nil means present (even if {}).
	Props   *nbWidgetProps `json:"props,omitempty"`
	Content *nbVizContent  `json:"content,omitempty"`
}

// nbWidgetProps: Title pointer distinguishes key-absent (nil) from empty-string ("").
type nbWidgetProps struct {
	Title *json.RawMessage `json:"title,omitempty"`
}

type nbVizContent struct {
	ID string `json:"id"`
}

// checkWidgetLevelProps walks the document tree and returns an error string for
// each widget that violates the props rules the New Relic UI enforces:
//
//   - viz.markdown: widget-level "props" must be absent entirely.
//   - all other viz: "props.title" is required (may be an empty string "").
func checkWidgetLevelProps(raw string) []string {
	var doc nbDoc
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil // JSON errors are caught by the envelope validator
	}

	var errs []string
	for ci, container := range doc.Content {
		for wi, widget := range container.Content {
			path := fmt.Sprintf("content[%d].content[%d]", ci, wi)

			if widget.Type != "widget" {
				errs = append(errs, fmt.Sprintf(`%s: type must be "widget", got %q`, path, widget.Type))
				continue
			}
			if widget.Content == nil || widget.Content.ID == "" {
				errs = append(errs, path+`: missing required "content" block (must include a visualization "id")`)
				continue
			}

			vizID := widget.Content.ID
			vizPath := fmt.Sprintf("%s (%s)", path, vizID)

			if vizID == "viz.markdown" {
				if widget.Props != nil {
					errs = append(errs, vizPath+`: widget-level "props" must not be present; markdown widgets do not support a title`)
				}
			} else {
				if widget.Props == nil {
					errs = append(errs, vizPath+`: widget-level "props" is required and must include "title" (use "" for no title)`)
				} else if widget.Props.Title == nil {
					errs = append(errs, vizPath+`: "props.title" is required (use "" for no title)`)
				}
			}
		}
	}
	return errs
}

// ── JSON normalization ────────────────────────────────────────────────────────

// normalizeNotebookContent converts any valid JSON string to a canonical form
// with alphabetically sorted keys and 2-space indentation, and returns the
// parsed value for direct use in API calls. Callers that only need the string
// can discard the second return value with _.
//
// Storing the canonical form in state means two documents with the same
// semantic content but different formatting always compare as equal, preventing
// spurious plan diffs when the user reformats or reorders keys.
func normalizeNotebookContent(raw string) (normalized string, body interface{}, err error) {
	if raw == "" {
		return "", nil, fmt.Errorf("notebook content must not be empty")
	}
	if err = json.Unmarshal([]byte(raw), &body); err != nil {
		return "", nil, fmt.Errorf("invalid JSON: %w", err)
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err = enc.Encode(body); err != nil {
		return "", nil, fmt.Errorf("could not re-serialize JSON: %w", err)
	}
	// json.Encoder appends a trailing newline; strip it so state comparisons are stable.
	normalized = string(bytes.TrimRight(buf.Bytes(), "\n"))
	return normalized, body, nil
}

// suppressEquivalentNotebookContent tells Terraform to ignore differences
// between two JSON strings that are semantically identical (same keys/values,
// different formatting or key order).
func suppressEquivalentNotebookContent(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	if oldVal == newVal {
		return true
	}
	normOld, _, errOld := normalizeNotebookContent(oldVal)
	normNew, _, errNew := normalizeNotebookContent(newVal)
	if errOld != nil || errNew != nil {
		return false
	}
	return normOld == normNew
}

// flattenNotebookContent normalizes the raw JSON returned by the API and stores
// it in the correct state field ("content" or "content").
func flattenNotebookContent(raw json.RawMessage, d *schema.ResourceData, field string) error {
	if len(raw) == 0 {
		return nil
	}
	normalized, _, err := normalizeNotebookContent(string(raw))
	if err != nil {
		return fmt.Errorf("could not normalize notebook content returned by the API: %w", err)
	}
	return d.Set(field, normalized)
}

// ── Resource helpers ──────────────────────────────────────────────────────────

// notebookOrgID returns the organization ID from state when available, or
// resolves it from provider credentials as a fallback.
func notebookOrgID(ctx context.Context, guid string, d *schema.ResourceData, pc *ProviderConfig) (string, error) {
	if id, _ := d.Get("organization_id").(string); id != "" {
		return id, nil
	}
	log.Printf("[DEBUG] organization_id not in state for notebook %s, resolving from provider credentials", guid)
	return getOrganizationID(ctx, pc, "")
}

// ── Error helpers ─────────────────────────────────────────────────────────────

// isNotebookNotFoundError returns true when the error indicates the notebook
// does not exist on the platform.
func isNotebookNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "404") ||
		strings.Contains(msg, "Blob not found")
}
