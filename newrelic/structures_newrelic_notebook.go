package newrelic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateNotebookContent is the ValidateFunc for both content and content_json.
// It enforces that the value is non-empty, valid JSON, and follows the
// minimum declarative UI envelope required by the Notebooks platform:
//
//	{ "type": "declarative", "version": 1, "content": [...] }
//
// This catches mistyped or stale schema formats (e.g. the legacy "blocks" key)
// at plan time so the user gets a clear error before any API call is made.
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

// nbDoc / nbContainer / nbWidget are lightweight typed structs for navigating
// the notebook document tree. Declaring only the fields needed for validation
// avoids the repeated interface{} type assertions that the untyped approach
// requires, and gives us a single clean parse pass through json.Unmarshal.
// Unrecognised JSON keys pass through untouched via the encoding/json package.

type nbDoc struct {
	Content []nbContainer `json:"content"`
}

type nbContainer struct {
	Content []nbWidget `json:"content"`
}

type nbWidget struct {
	// Props is a pointer so we can distinguish "absent" (nil) from "present but
	// empty" ({}) — both are structurally different for validation purposes.
	Props   *nbWidgetProps `json:"props,omitempty"`
	Content *nbVizContent  `json:"content,omitempty"`
}

// nbWidgetProps holds the widget-level props. Title is a pointer to
// json.RawMessage so we can distinguish key-absent (nil) from key-present-
// but-empty-string ("") — both are syntactically different.
type nbWidgetProps struct {
	Title *json.RawMessage `json:"title,omitempty"`
}

type nbVizContent struct {
	ID string `json:"id"`
}

// checkWidgetLevelProps walks the parsed notebook document and returns a
// human-readable error string for each widget that violates the props rules:
//
//   - viz.markdown: the widget-level "props" key must be absent entirely.
//   - all other viz: "props.title" is required (may be an empty string "").
//
// Using typed structs (nbDoc/nbContainer/nbWidget) instead of interface{} maps
// means a single json.Unmarshal call populates the full tree with no type
// assertions; validation is then a straightforward typed field check.
//
// The caller is responsible for surfacing the returned strings as Terraform
// diagnostics or errors with the appropriate format.
func checkWidgetLevelProps(raw string) []string {
	var doc nbDoc
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil // JSON errors are caught by the envelope validator
	}

	var errs []string
	for ci, container := range doc.Content {
		for wi, widget := range container.Content {
			if widget.Content == nil || widget.Content.ID == "" {
				continue
			}
			vizID := widget.Content.ID
			path := fmt.Sprintf("content[%d].widget[%d] (%s)", ci, wi, vizID)

			if vizID == "viz.markdown" {
				if widget.Props != nil {
					errs = append(errs, path+`: widget-level "props" must not be present; markdown widgets do not support a title`)
				}
			} else {
				if widget.Props == nil {
					errs = append(errs, path+`: widget-level "props" is required and must include "title" (use "" for no title)`)
				} else if widget.Props.Title == nil {
					errs = append(errs, path+`: "props.title" is required (use "" for no title)`)
				}
			}
		}
	}
	return errs
}

// customizeNotebookDiff enforces widget-level props rules at plan time.
// Placing this in CustomizeDiff (rather than ValidateFunc) means Terraform
// shows only the attribute name in the error location, not the full JSON body —
// the body can be many hundreds of characters and its presence in the error
// output makes the actual violation harder to find.
func customizeNotebookDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip the parse-and-walk if neither content field changed — there is nothing
	// new to validate and the previous plan already caught any violations.
	if !d.HasChanges("content", "content_json") {
		return nil
	}

	raw := d.Get("content").(string)
	if raw == "" {
		raw = d.Get("content_json").(string)
	}
	if raw == "" {
		return nil // ValidateFunc already reports missing content
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

// parseAndNormalizeContent unmarshals raw JSON, produces the canonical
// normalized string for state storage, and returns the parsed value for direct
// use in API calls — avoiding the double-parse that would occur if callers ran
// normalizeNotebookContent followed by a separate json.Unmarshal.
func parseAndNormalizeContent(raw string) (normalized string, body interface{}, err error) {
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
	normalized = string(bytes.TrimRight(buf.Bytes(), "\n"))
	return normalized, body, nil
}

// isNotebookNotFoundError returns true when an error from the Notebooks client
// indicates the requested notebook does not exist. Centralised here so both the
// resource and data source use the same detection logic.
func isNotebookNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "404") ||
		strings.Contains(msg, "Blob not found")
}
