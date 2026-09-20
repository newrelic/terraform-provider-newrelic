package newrelic

// structures_newrelic_ngep.go contains expand/flatten helpers shared across
// all NGEP-backed resources. Add new shared transforms here rather than in a
// resource-specific structures file to avoid duplication.

import (
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── Tags ──────────────────────────────────────────────────────────────────────

// expandNGEPTags converts a Terraform list of "key:value1,value2" strings into
// EntityManagementTagInput values accepted by any entityManagement mutation.
func expandNGEPTags(raw []interface{}) []scorecards.EntityManagementTagInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]scorecards.EntityManagementTagInput, 0, len(raw))
	for _, r := range raw {
		s, ok := r.(string)
		if !ok {
			continue
		}
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		vals := strings.Split(strings.TrimSpace(parts[1]), ",")
		for i := range vals {
			vals[i] = strings.TrimSpace(vals[i])
		}
		out = append(out, scorecards.EntityManagementTagInput{Key: key, Values: vals})
	}
	return out
}

// flattenNGEPTags converts EntityManagementTag values back to "key:value1,value2"
// strings. Tags whose keys begin with "nr." are stripped — NGEP auto-injects
// system tags (e.g. "nr.hierarchy.level") that must not appear in Terraform
// state and trigger spurious plan diffs.
//
// The output is sorted alphabetically by the full "key:values" string to
// guarantee a stable ordering across API calls. The NGEP API does not preserve
// tag order, so without sorting every plan could show spurious reordering diffs.
func flattenNGEPTags(tags []scorecards.EntityManagementTag) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if strings.HasPrefix(t.Key, "nr.") {
			continue
		}
		out = append(out, t.Key+":"+strings.Join(t.Values, ","))
	}
	return out
}
