package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
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
		for _, raw := range membersSet.(*schema.Set).List() {
			memberIDs[raw.(map[string]interface{})["user_id"].(int)] = true
		}
	}
	if managersOk {
		for _, raw := range managersSet.(*schema.Set).List() {
			uid := raw.(map[string]interface{})["user_id"].(int)
			if !memberIDs[uid] {
				errs = append(errs, fmt.Sprintf(
					"manager user_id %d is not listed in the members block — "+
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

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}

// ── Expand helpers (Terraform state → API input) ──────────────────────────────

// expandTeamResources converts the resources Terraform list to CreateInput.
func expandTeamResources(raw []interface{}) []scorecards.EntityManagementTeamResourceCreateInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]scorecards.EntityManagementTeamResourceCreateInput, 0, len(raw))
	for _, r := range raw {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		res := scorecards.EntityManagementTeamResourceCreateInput{
			Type:    m["type"].(string),
			Content: m["content"].(string),
		}
		if title, ok := m["title"].(string); ok && title != "" {
			res.Title = title
		}
		out = append(out, res)
	}
	return out
}

// expandTeamResourcesUpdate converts the resources Terraform list to UpdateInput.
func expandTeamResourcesUpdate(raw []interface{}) []scorecards.EntityManagementTeamResourceUpdateInput {
	if len(raw) == 0 {
		return nil
	}
	out := make([]scorecards.EntityManagementTeamResourceUpdateInput, 0, len(raw))
	for _, r := range raw {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		res := scorecards.EntityManagementTeamResourceUpdateInput{
			Type:    m["type"].(string),
			Content: m["content"].(string),
		}
		if title, ok := m["title"].(string); ok && title != "" {
			res.Title = title
		}
		out = append(out, res)
	}
	return out
}

// expandUserIDsFromSet extracts integer user IDs from a TypeSet whose elements
// are maps with a single "user_id" int key. Used for both the members and
// managers blocks which have the same shape.
func expandUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, raw := range s.List() {
		out = append(out, raw.(map[string]interface{})["user_id"].(int))
	}
	return out
}

// expandEntityGUIDsFromSet extracts string GUIDs from a TypeSet whose elements
// are maps with a single "guid" string key.
func expandEntityGUIDsFromSet(s *schema.Set) []string {
	out := make([]string, 0, s.Len())
	for _, raw := range s.List() {
		out = append(out, raw.(map[string]interface{})["guid"].(string))
	}
	return out
}

// ── Flatten helpers (API response → Terraform state) ─────────────────────────

// flattenTeamResources converts API TeamResource values back to Terraform maps.
func flattenTeamResources(res []scorecards.EntityManagementTeamResource) []map[string]interface{} {
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

// flattenMemberUserIDs converts a []int of user IDs into the list-of-maps
// shape that the members TypeSet expects.
func flattenMemberUserIDs(userIDs []int) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(userIDs))
	for _, id := range userIDs {
		out = append(out, map[string]interface{}{"user_id": id})
	}
	return out
}

// flattenEntityGUIDs converts a []string of entity GUIDs into the list-of-maps
// shape that the entities TypeSet expects.
func flattenEntityGUIDs(guids []string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(guids))
	for _, g := range guids {
		out = append(out, map[string]interface{}{"guid": g})
	}
	return out
}

// decodeManagerGUIDsToUserIDs converts the NGEP-encoded manager GUID list from
// the TeamEntity back to integer userIDs using the GUID→userID map produced by
// readTeamMembershipMap. This avoids a second API call on every Read and enables
// idempotent manager state storage.
//
// Returns nil when either argument is empty, which causes Terraform to keep
// whatever was last written to state for the managers block (safe for import
// and for the case where managers haven't been set yet).
func decodeManagerGUIDsToUserIDs(managerGUIDs []string, memberGUIDToUserID map[string]int) []map[string]interface{} {
	if len(managerGUIDs) == 0 || len(memberGUIDToUserID) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(managerGUIDs))
	for _, guid := range managerGUIDs {
		if uid, ok := memberGUIDToUserID[guid]; ok {
			out = append(out, map[string]interface{}{"user_id": uid})
		}
	}
	return out
}
