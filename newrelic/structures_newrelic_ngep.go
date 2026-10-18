package newrelic

// structures_newrelic_ngep.go contains expand/flatten helpers shared across
// all NGEP-backed resources. Add new shared transforms here rather than in a
// resource-specific structures file to avoid duplication.

import (
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

// ── Tags ──────────────────────────────────────────────────────────────────────

// expandSAITags converts a Terraform list of tag blocks (each with a "key"
// string and "values" []string) into EntityManagementTagInput values accepted
// by any entityManagement mutation.
func expandSAITags(raw []interface{}) []servicearchintelligence.EntityManagementTagInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]servicearchintelligence.EntityManagementTagInput, 0, len(raw))
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
			out = append(out, servicearchintelligence.EntityManagementTagInput{Key: key, Values: vals})
		}
	}
	return out
}

// tagsInputToFlattenedSet converts TagInput values to the flattened Terraform
// representation ready for d.Set("tags", ...). Used in Create to set state
// without a Read round-trip.
func tagsInputToFlattenedSet(tags []servicearchintelligence.EntityManagementTagInput) interface{} {
	if len(tags) == 0 {
		return nil
	}
	ts := make([]servicearchintelligence.EntityManagementTag, len(tags))
	for i, t := range tags {
		ts[i] = servicearchintelligence.EntityManagementTag(t)
	}
	return flattenSAITags(ts)
}

// flattenSAITags converts EntityManagementTag values back to a list of
// maps with "key" and "values" keys. Tags whose keys begin with "nr." are
// stripped — NGEP auto-injects system tags (e.g. "nr.hierarchy.level") that
// must not appear in Terraform state and trigger spurious plan diffs.
func flattenSAITags(tags []servicearchintelligence.EntityManagementTag) []map[string]interface{} {
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
