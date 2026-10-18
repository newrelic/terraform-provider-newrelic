package newrelic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ── Validation ────────────────────────────────────────────────────────────────

// checkEnvelopeErrors validates the minimum declarative UI envelope:
//
//	{ "type": "declarative", "version": 1, "content": [ <single container> ] }
//
// It returns a slice of plain error strings (never the raw JSON) so callers can
// pool them with other validation results before surfacing them to the user.
func checkEnvelopeErrors(raw string) []string {
	if raw == "" {
		return []string{"content must not be empty"}
	}

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return []string{jsonParseHint(raw, err)}
	}

	var errs []string

	t, hasType := doc["type"]
	if !hasType {
		errs = append(errs, `missing required field "type" — expected "declarative"`)
	} else if t != "declarative" {
		errs = append(errs, fmt.Sprintf(`"type" must be "declarative", got %q`, t))
	}

	if ver, hasVersion := doc["version"]; !hasVersion {
		errs = append(errs, `missing required field "version" — expected the integer 1`)
	} else {
		// json.Unmarshal decodes numbers as float64; accept integer 1 only.
		if f, ok := ver.(float64); !ok || f != 1 {
			errs = append(errs, fmt.Sprintf(`"version" must be the integer 1, got %s`, jsonValueDesc(ver)))
		}
	}

	if _, hasContent := doc["content"]; !hasContent {
		errs = append(errs, `missing required field "content" — expected an array containing exactly one container block`)
	} else if containers, isArray := doc["content"].([]interface{}); !isArray {
		errs = append(errs, `"content" must be a JSON array`)
	} else if len(containers) > 1 {
		errs = append(errs, fmt.Sprintf(
			`"content" array must contain exactly one container block, got %d — the New Relic UI only renders the first container; place all widgets inside a single container`,
			len(containers),
		))
	} else if len(containers) == 1 {
		// Validate the single container has at least one widget.
		if container, ok := containers[0].(map[string]interface{}); ok {
			widgets, _ := container["content"].([]interface{})
			if len(widgets) == 0 {
				errs = append(errs, `the container in "content[0]" has no widgets — add at least one widget block inside the container's "content" array`)
			}
		}
	}

	return errs
}

// jsonParseHint converts a json.Unmarshal error into a user-friendly message.
// It surfaces line/column info from json.SyntaxError (far more useful than a
// raw byte offset) and adds actionable hints for common mistakes such as using
// comments or trailing commas in raw JSON strings.
func jsonParseHint(raw string, err error) string {
	var synErr *json.SyntaxError
	if !errors.As(err, &synErr) {
		return fmt.Sprintf("not valid JSON: %s", err)
	}

	offset := synErr.Offset
	line, col := byteOffsetToLineCol(raw, offset)
	base := fmt.Sprintf("JSON syntax error at line %d, column %d", line, col)

	if offset > 0 && int(offset) <= len(raw) {
		ch := rune(raw[offset-1])
		switch {
		case ch == '#':
			base += " — found '#': JSON does not support comments; use jsonencode({...}) in Terraform instead of a raw JSON string"
		case ch == '/' && int(offset) < len(raw) && raw[offset] == '/':
			// SyntaxError.Offset points to the first '/' of '//'; the second '/' is raw[offset].
			base += " — found '//': JSON does not support comments; use jsonencode({...}) instead"
		default:
			base += fmt.Sprintf(" — unexpected character %q", ch)
		}
	}

	return base
}

// byteOffsetToLineCol converts a 1-based byte offset (as returned by
// json.SyntaxError.Offset) to a 1-based line and column number. Both line and
// column are 1-based so they match what editors display.
func byteOffsetToLineCol(s string, offset int64) (line, col int) {
	line, col = 1, 1
	for i := int64(0); i < offset-1 && i < int64(len(s)); i++ {
		if s[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

// jsonValueDesc formats a json.Unmarshal-decoded value in a user-friendly way,
// avoiding Go-internal type names (float64, bool, etc.) in error messages.
func jsonValueDesc(v interface{}) string {
	switch val := v.(type) {
	case float64:
		// JSON numbers are decoded as float64; show the numeric value.
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d (number)", int64(val))
		}
		return fmt.Sprintf("%g (number)", val)
	case string:
		return fmt.Sprintf("%q (string)", val)
	case bool:
		return fmt.Sprintf("%v (boolean)", val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%T", v)
	}
}

// customizeNotebookDiff runs all content validation at plan time.
// Using CustomizeDiff (rather than ValidateFunc) means Terraform never echoes
// the full JSON body in the error output — only the concise error messages are shown.
func customizeNotebookDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	raw := d.Get("content").(string)

	// Skip validation only when the content is unchanged AND non-empty — meaning
	// it was already validated on a previous plan. An empty string always requires
	// validation because the Terraform SDK treats nil and "" as equivalent for string
	// fields, so d.HasChange returns false even for a brand-new resource with content = "".
	if raw != "" && !d.HasChange("content") {
		return nil
	}

	// Pool envelope errors first; if any exist, skip the widget-level check
	// because the document structure may be too broken to walk safely.
	var allErrs []string
	allErrs = append(allErrs, checkEnvelopeErrors(raw)...)
	if len(allErrs) == 0 {
		allErrs = append(allErrs, checkWidgetLevelProps(raw)...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	msg := "invalid notebook content:\n"
	for i, e := range allErrs {
		msg += fmt.Sprintf("  (%d) %s\n", i+1, e)
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
