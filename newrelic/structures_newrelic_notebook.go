package newrelic

import (
	"bytes"
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

	// Widget-level props rules — enforced to match the convention the New Relic
	// UI uses when saving notebooks. Violating these causes drift the next time
	// the notebook is edited via the UI.
	//
	//   viz.markdown: the "props" key must be absent from the widget object.
	//                 Markdown does not support a title; the UI omits props entirely.
	//
	//   all other viz: "props.title" is required (may be an empty string "").
	//                  The UI always writes props.title; omitting it causes the UI
	//                  to inject it on the next save, surfacing false drift.
	if containers, ok := doc["content"].([]interface{}); ok {
		for ci, c := range containers {
			container, cOK := c.(map[string]interface{})
			if !cOK {
				continue
			}
			widgets, wOK := container["content"].([]interface{})
			if !wOK {
				continue
			}
			for wi, w := range widgets {
				widget, wdOK := w.(map[string]interface{})
				if !wdOK {
					continue
				}
				vizContent, _ := widget["content"].(map[string]interface{})
				vizID, _ := vizContent["id"].(string)
				if vizID == "" {
					continue
				}
				_, hasWidgetProps := widget["props"]

				if vizID == "viz.markdown" {
					if hasWidgetProps {
						errors = append(errors, fmt.Errorf(
							`%q: container[%d].widget[%d] (viz.markdown): `+
								`widget-level "props" must not be present — `+
								`markdown does not support a title; remove the "props" key entirely`,
							k, ci, wi))
					}
				} else {
					widgetProps, _ := widget["props"].(map[string]interface{})
					if widgetProps == nil {
						errors = append(errors, fmt.Errorf(
							`%q: container[%d].widget[%d] (%s): `+
								`widget-level "props" is required and must include "title" `+
								`(use "" for no title)`,
							k, ci, wi, vizID))
					} else if _, hasTitle := widgetProps["title"]; !hasTitle {
						errors = append(errors, fmt.Errorf(
							`%q: container[%d].widget[%d] (%s): `+
								`widget-level "props.title" is required (use "" for no title)`,
							k, ci, wi, vizID))
					}
				}
			}
		}
	}

	return
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
