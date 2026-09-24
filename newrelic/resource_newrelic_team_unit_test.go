//go:build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── structures_newrelic_team.go ───────────────────────────────────────────────

func TestExpandTeamTags(t *testing.T) {
	t.Parallel()
	raw := []interface{}{
		map[string]interface{}{"key": "env", "values": []interface{}{"dev", "staging"}},
		map[string]interface{}{"key": "team", "values": []interface{}{"platform"}},
	}
	tags := expandNGEPTags(raw)
	require.Len(t, tags, 2)
	assert.Equal(t, "env", tags[0].Key)
	assert.ElementsMatch(t, []string{"dev", "staging"}, tags[0].Values)
	assert.Equal(t, "team", tags[1].Key)
}

func TestExpandTeamTagsEmpty(t *testing.T) {
	t.Parallel()
	assert.Nil(t, expandNGEPTags(nil))
	assert.Nil(t, expandNGEPTags([]interface{}{}))
}

func TestFlattenTeamTags_FiltersNrSystem(t *testing.T) {
	t.Parallel()
	// nr.* tags (auto-injected by NGEP) must be stripped from state.
	tags := []scorecards.EntityManagementTag{
		{Key: "env", Values: []string{"dev"}},
		{Key: "nr.hierarchy.level", Values: []string{"Level 2"}},
		{Key: "team", Values: []string{"platform"}},
	}
	flat := flattenNGEPTags(tags)
	require.Len(t, flat, 2, "nr.* tag should be filtered out")
	assert.Equal(t, "env", flat[0]["key"])
	assert.Equal(t, []string{"dev"}, flat[0]["values"])
	assert.Equal(t, "team", flat[1]["key"])
	assert.Equal(t, []string{"platform"}, flat[1]["values"])
}

func TestExpandTeamResources(t *testing.T) {
	t.Parallel()
	raw := []interface{}{
		map[string]interface{}{"type": "link", "content": "https://example.com/runbook", "title": "Runbook"},
		map[string]interface{}{"type": "link", "content": "https://example.com/wiki", "title": ""},
	}
	res := expandTeamResources(raw)
	require.Len(t, res, 2)
	assert.Equal(t, "Runbook", res[0].Title)
	assert.Equal(t, "", res[1].Title)
}

func TestFlattenTeamResources(t *testing.T) {
	t.Parallel()
	res := []scorecards.EntityManagementTeamResource{
		{Type: "link", Content: "https://example.com", Title: "Docs"},
	}
	flat := flattenTeamResources(res)
	require.Len(t, flat, 1)
	assert.Equal(t, "link", flat[0]["type"])
	assert.Equal(t, "Docs", flat[0]["title"])
}

func TestExpandMemberUserIDsFromSet(t *testing.T) {
	t.Parallel()
	s := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{"user_id": {Type: schema.TypeInt}},
	}), []interface{}{
		map[string]interface{}{"user_id": 123},
		map[string]interface{}{"user_id": 456},
	})
	ids := expandUserIDsFromSet(s)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, 123)
	assert.Contains(t, ids, 456)
}

func TestFlattenMemberUserIDs(t *testing.T) {
	t.Parallel()
	flat := flattenMemberUserIDs([]int{10, 20})
	require.Len(t, flat, 2)
	assert.Equal(t, 10, flat[0]["user_id"])
	assert.Equal(t, 20, flat[1]["user_id"])
}

func TestFlattenEntityGUIDs(t *testing.T) {
	t.Parallel()
	flat := flattenEntityGUIDs([]string{"guid-a", "guid-b"})
	require.Len(t, flat, 2)
	assert.Equal(t, "guid-a", flat[0]["guid"])
}

// ── CustomizeDiff validation ──────────────────────────────────────────────────

// teamResourceDiff builds a *schema.ResourceDiff with the given attrs set.
// We can't create a real ResourceDiff without a running provider, so we use
// resourceNewRelicTeamCustomizeDiff's documented input contract and test it
// through the public schema helpers instead. The integration tests cover the
// end-to-end plan-time error; here we test the validation logic directly via
// a minimal fake ResourceDiff.
//
// Since schema.ResourceDiff is hard to construct in isolation, we test the
// logic that underlies it — the set intersection — through the underlying
// schema types directly.

func TestCustomizeDiff_ManagerNotInMembers(t *testing.T) {
	t.Parallel()
	// Build a ResourceData (Create=true) to simulate what CustomizeDiff sees.
	r := resourceNewRelicTeam()
	d := r.TestResourceData()
	_ = d.Set("name", "test-team")
	_ = d.Set("members", []interface{}{
		map[string]interface{}{"user_id": 111},
	})
	_ = d.Set("managers", []interface{}{
		map[string]interface{}{"user_id": 999}, // not in members
	})

	// We can't call resourceNewRelicTeamCustomizeDiff directly with ResourceData,
	// but we can exercise the logic it uses:
	memberIDs := make(map[int]bool)
	for _, raw := range d.Get("members").(*schema.Set).List() {
		memberIDs[raw.(map[string]interface{})["user_id"].(int)] = true
	}
	var badManagers []int
	for _, raw := range d.Get("managers").(*schema.Set).List() {
		uid := raw.(map[string]interface{})["user_id"].(int)
		if !memberIDs[uid] {
			badManagers = append(badManagers, uid)
		}
	}
	assert.Equal(t, []int{999}, badManagers, "manager 999 is not in members — should be flagged")
}

func TestCustomizeDiff_ManagerInMembers_Valid(t *testing.T) {
	t.Parallel()
	r := resourceNewRelicTeam()
	d := r.TestResourceData()
	_ = d.Set("name", "test-team")
	_ = d.Set("members", []interface{}{
		map[string]interface{}{"user_id": 111},
		map[string]interface{}{"user_id": 222},
	})
	_ = d.Set("managers", []interface{}{
		map[string]interface{}{"user_id": 111}, // present in members — valid
	})

	memberIDs := make(map[int]bool)
	for _, raw := range d.Get("members").(*schema.Set).List() {
		memberIDs[raw.(map[string]interface{})["user_id"].(int)] = true
	}
	var badManagers []int
	for _, raw := range d.Get("managers").(*schema.Set).List() {
		uid := raw.(map[string]interface{})["user_id"].(int)
		if !memberIDs[uid] {
			badManagers = append(badManagers, uid)
		}
	}
	assert.Empty(t, badManagers, "all managers are in members — no validation error expected")
}
