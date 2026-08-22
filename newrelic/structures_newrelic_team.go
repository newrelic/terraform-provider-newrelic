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
//  1. Every user_id in the managers block must also appear in the members
//     block. NGEP enforces this at mutation time, but surfacing it at plan
//     time gives a much clearer error than a cryptic API 400.
//
//  2. aliases must not contain the empty string.
func resourceNewRelicTeamCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	var errs []string

	// ── managers ⊆ members ─────────────────────────────────────────────────
	managersSet, managersOk := d.GetOk("managers")
	membersSet, membersOk := d.GetOk("members")

	if managersOk && membersOk {
		memberIDs := make(map[int]bool)
		for _, raw := range membersSet.(*schema.Set).List() {
			memberIDs[raw.(map[string]interface{})["user_id"].(int)] = true
		}
		for _, raw := range managersSet.(*schema.Set).List() {
			uid := raw.(map[string]interface{})["user_id"].(int)
			if !memberIDs[uid] {
				errs = append(errs, fmt.Sprintf(
					"manager user_id %d is not listed in the members block — "+
						"every manager must first be a member of the team", uid))
			}
		}
	}

	if managersOk && !membersOk {
		for _, raw := range managersSet.(*schema.Set).List() {
			uid := raw.(map[string]interface{})["user_id"].(int)
			errs = append(errs, fmt.Sprintf(
				"manager user_id %d requires a matching members block entry — "+
					"every manager must first be a member of the team", uid))
		}
	}

	// ── aliases must not be empty strings ──────────────────────────────────
	if aliases, ok := d.GetOk("aliases"); ok {
		for _, a := range aliases.([]interface{}) {
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

// expandTeamTags converts a list of "key:value1,value2" strings into
// EntityManagementTagInput values.
func expandTeamTags(raw []interface{}) []scorecards.EntityManagementTagInput {
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

// expandTeamResources converts the Terraform resources list to CreateInput.
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

// expandTeamResourcesUpdate converts the Terraform resources list to UpdateInput.
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

// expandMemberUserIDsFromSet converts a TypeSet of {user_id} maps to []int.
func expandMemberUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, raw := range s.List() {
		out = append(out, raw.(map[string]interface{})["user_id"].(int))
	}
	return out
}

// expandManagerUserIDsFromSet converts a TypeSet of {user_id} maps to []int.
func expandManagerUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, raw := range s.List() {
		out = append(out, raw.(map[string]interface{})["user_id"].(int))
	}
	return out
}

// expandEntityGUIDsFromSet converts a TypeSet of {guid} maps to []string.
func expandEntityGUIDsFromSet(s *schema.Set) []string {
	out := make([]string, 0, s.Len())
	for _, raw := range s.List() {
		out = append(out, raw.(map[string]interface{})["guid"].(string))
	}
	return out
}

// ── Flatten helpers (API response → Terraform state) ─────────────────────────

// flattenTeamTags converts EntityManagementTag values back to
// "key:value1,value2" strings. System-managed tags (keys starting with "nr.")
// such as "nr.hierarchy.level" are filtered out — they are auto-injected by
// NGEP when a team is assigned a parent and must not appear in state.
func flattenTeamTags(tags []scorecards.EntityManagementTag) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if strings.HasPrefix(t.Key, "nr.") {
			continue
		}
		out = append(out, t.Key+":"+strings.Join(t.Values, ","))
	}
	return out
}

// flattenTeamResources converts API resource values back to Terraform maps.
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

// flattenMemberUserIDs converts []int → list-of-maps for the members TypeSet.
func flattenMemberUserIDs(userIDs []int) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(userIDs))
	for _, id := range userIDs {
		out = append(out, map[string]interface{}{"user_id": id})
	}
	return out
}

// flattenEntityGUIDs converts []string → list-of-maps for the entities TypeSet.
func flattenEntityGUIDs(guids []string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(guids))
	for _, g := range guids {
		out = append(out, map[string]interface{}{"guid": g})
	}
	return out
}

// flattenManagerUserIDs decodes the NGEP-encoded manager GUIDs back to
// integer userIds using the GUID→userId map built from the team's membership
// collection. Returns nil when the map is empty (e.g. on import before the
// collection has been read).
func flattenManagerUserIDs(managerGUIDs []string, guidToUserID map[string]int) []map[string]interface{} {
	if len(managerGUIDs) == 0 || len(guidToUserID) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(managerGUIDs))
	for _, guid := range managerGUIDs {
		if uid, ok := guidToUserID[guid]; ok {
			out = append(out, map[string]interface{}{"user_id": uid})
		}
	}
	return out
}
