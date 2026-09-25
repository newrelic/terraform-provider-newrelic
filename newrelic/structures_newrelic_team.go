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

// expandUserIDsFromSet extracts integer user IDs from a TypeSet of plain ints.
// Used for both the members and managers attributes.
func expandUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, v := range s.List() {
		out = append(out, v.(int))
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

// flattenMemberUserIDs returns the []int directly for the members TypeSet,
// which now stores plain ints rather than maps.
func flattenMemberUserIDs(userIDs []int) []int {
	return userIDs
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
// Returns nil when either argument is empty. Callers pass the result directly
// to d.Set("managers", ...) — a nil value sets the managers attribute to an
// empty set (count=0), which is correct when no managers have been configured.
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
