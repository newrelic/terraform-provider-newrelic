package newrelic

// structures_newrelic_ngep.go contains expand/flatten helpers shared across
// all NGEP-backed resources. Add new shared transforms here rather than in a
// resource-specific structures file to avoid duplication.

import (
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── Tags ──────────────────────────────────────────────────────────────────────

// expandNGEPTags converts a Terraform list of tag blocks (each with a "key"
// string and "values" []string) into EntityManagementTagInput values accepted
// by any entityManagement mutation.
func expandNGEPTags(raw []interface{}) []scorecards.EntityManagementTagInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]scorecards.EntityManagementTagInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		key := m["key"].(string)
		valuesRaw := m["values"].([]interface{})
		vals := make([]string, 0, len(valuesRaw))
		for _, v := range valuesRaw {
			if s, ok := v.(string); ok && s != "" {
				vals = append(vals, s)
			}
		}
		if len(vals) > 0 {
			out = append(out, scorecards.EntityManagementTagInput{Key: key, Values: vals})
		}
	}
	return out
}

// flattenNGEPTags converts EntityManagementTag values back to a list of
// maps with "key" and "values" keys. Tags whose keys begin with "nr." are
// stripped — NGEP auto-injects system tags (e.g. "nr.hierarchy.level") that
// must not appear in Terraform state and trigger spurious plan diffs.
func flattenNGEPTags(tags []scorecards.EntityManagementTag) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(tags))
	for _, t := range tags {
		if strings.HasPrefix(t.Key, "nr.") {
			continue
		}
		out = append(out, map[string]interface{}{
			"key":    t.Key,
			"values": t.Values,
		})
	}
	return out
}
