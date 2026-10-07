package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

// ── CustomizeDiff ─────────────────────────────────────────────────────────────

// resourceNewRelicTeamCustomizeDiff runs plan-time validations for
// newrelic_team:
//
//  1. managers ⊆ members — every user_id in the managers block must also
//     appear in the members block. NGEP enforces this at mutation time, but
//     surfacing it at plan time gives a clearer error than an API 400.
//
//  2. aliases must not contain empty or whitespace-only strings.
func resourceNewRelicTeamCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	var errs []string

	// ── managers ⊆ members ─────────────────────────────────────────────────
	managersSet, managersOk := d.GetOk("managers")
	membersSet, membersOk := d.GetOk("members")

	memberIDs := make(map[int]bool)
	if membersOk {
		for _, v := range membersSet.(*schema.Set).List() {
			memberIDs[v.(int)] = true
		}
	}
	if managersOk {
		for _, v := range managersSet.(*schema.Set).List() {
			uid := v.(int)
			if !memberIDs[uid] {
				errs = append(errs, fmt.Sprintf(
					"manager user_id %d is not in the members list — "+
						"every manager must first be a member of the team", uid))
			}
		}
	}

	// ── aliases must not be empty ───────────────────────────────────────────
	if aliases, ok := d.GetOk("aliases"); ok {
		for _, a := range aliases.(*schema.Set).List() {
			if strings.TrimSpace(a.(string)) == "" {
				errs = append(errs, "aliases must not contain empty strings")
				break
			}
		}
	}

	// ── entity_management_mode vs entities consistency ───────────────────────
	// We use GetRawConfig() here (which IS reliable in CustomizeDiff) to detect
	// whether the entities attribute is present in the user's config.
	mode := d.Get("entity_management_mode").(string)
	if mode == "" {
		mode = "managed"
	}

	rc := d.GetRawConfig()
	if rc.IsKnown() && !rc.IsNull() {
		entitiesAttr := rc.GetAttr("entities")
		entitiesInConfig := entitiesAttr.IsKnown() && !entitiesAttr.IsNull()

		if mode == "unmanaged" && entitiesInConfig {
			errs = append(errs, "entities cannot be specified when entity_management_mode = \"unmanaged\" — "+
				"in unmanaged mode, entity ownership is controlled entirely by tag-based discovery and/or the Teams UI, not Terraform")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	// ── Suppress / align entity diff during mode transitions ─────────────────
	// When transitioning to "unmanaged", plan entities as [] so the plan agrees
	// with what Update will do (clear entities from state tracking). Without
	// this, Terraform Core rejects Update's d.Set("entities", []) because the
	// plan committed to the old value via an implicit "no change".
	//
	// When transitioning FROM "unmanaged" to "managed", no suppression is needed:
	// Read (running under the old unmanaged mode) returns [] for entities, and
	// the config carries the desired entities — the diff is correct as-is.
	oldMode, newMode := d.GetChange("entity_management_mode")
	if oldMode.(string) != "" && newMode.(string) == "unmanaged" {
		if err := d.SetNew("entities", []interface{}{}); err != nil {
			return err
		}
	}

	return nil
}

// ── Config-aware warning helper ───────────────────────────────────────────────

// declaredSetFromConfig returns the set of entity GUIDs the user has declared
// in their Terraform config. It reads d.GetRawConfig() so that the out-of-band
// warning is suppressed the moment the user adds an entity to their entities
// block — even before they run apply — preventing a confusing "one extra
// warning" on the apply that should clear the drift.
//
// Falls back to stateGUIDs when the raw config is unavailable (e.g., during
// terraform import or terraform refresh without plan context).
func declaredSetFromConfig(d *schema.ResourceData, stateGUIDs []string) map[string]bool {
	rc := d.GetRawConfig()
	if rc.IsKnown() && !rc.IsNull() {
		attr := rc.GetAttr("entities")
		if attr.IsKnown() && !attr.IsNull() {
			result := make(map[string]bool)
			for it := attr.ElementIterator(); it.Next(); {
				_, v := it.Element()
				if v.IsKnown() && !v.IsNull() {
					result[v.AsString()] = true
				}
			}
			return result
		}
	}
	result := make(map[string]bool, len(stateGUIDs))
	for _, g := range stateGUIDs {
		result[g] = true
	}
	return result
}

// ── Expand helpers (Terraform state → API input) ──────────────────────────────

// expandTeamResources converts the resources Terraform list to a CreateInput slice.
func expandTeamResources(raw []interface{}) []servicearchintelligence.EntityManagementTeamResourceCreateInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]servicearchintelligence.EntityManagementTeamResourceCreateInput, 0, len(raw))
	for _, r := range raw {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, servicearchintelligence.EntityManagementTeamResourceCreateInput{
			Type:    m["type"].(string),
			Content: m["content"].(string),
			Title:   m["title"].(string),
		})
	}
	return out
}

// expandTeamResourcesUpdate converts the resources Terraform list to an UpdateInput slice.
func expandTeamResourcesUpdate(raw []interface{}) []servicearchintelligence.EntityManagementTeamResourceUpdateInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]servicearchintelligence.EntityManagementTeamResourceUpdateInput, 0, len(raw))
	for _, r := range raw {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, servicearchintelligence.EntityManagementTeamResourceUpdateInput{
			Type:    m["type"].(string),
			Content: m["content"].(string),
			Title:   m["title"].(string),
		})
	}
	return out
}

// expandUserIDsFromSet extracts integer user IDs from a TypeSet of plain ints.
// Used for both the members and managers attributes.
func expandUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, v := range s.List() {
		out = append(out, v.(int))
	}
	return out
}

// expandEntityGUIDsFromSet extracts entity GUIDs from a TypeSet of plain strings.
func expandEntityGUIDsFromSet(s *schema.Set) []string {
	out := make([]string, 0, s.Len())
	for _, v := range s.List() {
		out = append(out, v.(string))
	}
	return out
}

// ── Flatten helpers (API response → Terraform state) ─────────────────────────

// flattenTeamResources converts API TeamResource values back to Terraform maps.
func flattenTeamResources(res []servicearchintelligence.EntityManagementTeamResource) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(res))
	for _, r := range res {
		out = append(out, map[string]interface{}{
			"type":    r.Type,
			"content": r.Content,
			"title":   r.Title,
		})
	}
	return out
}

// flattenMemberUserIDs returns the integer user ID slice for use with d.Set("members", ...).
func flattenMemberUserIDs(userIDs []int) []int {
	return userIDs
}

// flattenEntityGUIDs returns the GUID string slice for use with d.Set("entities", ...).
func flattenEntityGUIDs(guids []string) []string {
	return guids
}

// decodeManagerGUIDsToUserIDs converts the NGEP manager GUID list from the
// TeamEntity back to integer user IDs using the GUID→userID map built by
// readTeamMembershipMap. This avoids a second API call on every Read and keeps
// the manager state idempotent.
//
// Returns nil when either argument is empty; d.Set("managers", nil) sets the
// attribute to an empty set, which is correct when no managers are configured.
func decodeManagerGUIDsToUserIDs(managerGUIDs []string, memberGUIDToUserID map[string]int) []int {
	if len(managerGUIDs) == 0 || len(memberGUIDToUserID) == 0 {
		return nil
	}
	out := make([]int, 0, len(managerGUIDs))
	for _, guid := range managerGUIDs {
		if uid, ok := memberGUIDToUserID[guid]; ok {
			out = append(out, uid)
		}
	}
	return out
}
