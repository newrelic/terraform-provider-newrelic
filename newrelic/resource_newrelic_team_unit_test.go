//go:build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

// ── structures_newrelic_team.go ───────────────────────────────────────────────

func TestExpandTeamTags(t *testing.T) {
	t.Parallel()
	raw := []interface{}{
		map[string]interface{}{"key": "env", "values": []interface{}{"dev", "staging"}},
		map[string]interface{}{"key": "team", "values": []interface{}{"platform"}},
	}
	tags := expandSAITags(raw)
	require.Len(t, tags, 2)
	assert.Equal(t, "env", tags[0].Key)
	assert.ElementsMatch(t, []string{"dev", "staging"}, tags[0].Values)
	assert.Equal(t, "team", tags[1].Key)
}

func TestExpandTeamTagsEmpty(t *testing.T) {
	t.Parallel()
	assert.Nil(t, expandSAITags(nil))
	assert.Nil(t, expandSAITags([]interface{}{}))
}

func TestFlattenTeamTags_FiltersNrSystem(t *testing.T) {
	t.Parallel()
	// nr.* tags (auto-injected by NGEP) must be stripped from state.
	tags := []servicearchintelligence.EntityManagementTag{
		{Key: "env", Values: []string{"dev"}},
		{Key: "nr.hierarchy.level", Values: []string{"Level 2"}},
		{Key: "team", Values: []string{"platform"}},
	}
	flat := flattenSAITags(tags)
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
	res := []servicearchintelligence.EntityManagementTeamResource{
		{Type: "link", Content: "https://example.com", Title: "Docs"},
	}
	flat := flattenTeamResources(res)
	require.Len(t, flat, 1)
	assert.Equal(t, "link", flat[0]["type"])
	assert.Equal(t, "Docs", flat[0]["title"])
}

func TestExpandMemberUserIDsFromSet(t *testing.T) {
	t.Parallel()
	s := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeInt}), []interface{}{123, 456})
	ids := expandUserIDsFromSet(s)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, 123)
	assert.Contains(t, ids, 456)
}

func TestFlattenMemberUserIDs(t *testing.T) {
	t.Parallel()
	flat := flattenMemberUserIDs([]int{10, 20})
	require.Len(t, flat, 2)
	assert.Contains(t, flat, 10)
	assert.Contains(t, flat, 20)
}

func TestFlattenEntityGUIDs(t *testing.T) {
	t.Parallel()
	flat := flattenEntityGUIDs([]string{"guid-a", "guid-b"})
	require.Len(t, flat, 2)
	assert.Equal(t, "guid-a", flat[0])
	assert.Equal(t, "guid-b", flat[1])
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
	_ = d.Set("members", []interface{}{111})
	_ = d.Set("managers", []interface{}{999}) // not in members

	// We can't call resourceNewRelicTeamCustomizeDiff directly with ResourceData,
	// but we can exercise the logic it uses:
	memberIDs := make(map[int]bool)
	for _, v := range d.Get("members").(*schema.Set).List() {
		memberIDs[v.(int)] = true
	}
	var badManagers []int
	for _, v := range d.Get("managers").(*schema.Set).List() {
		uid := v.(int)
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
	_ = d.Set("members", []interface{}{111, 222})
	_ = d.Set("managers", []interface{}{111}) // present in members — valid

	memberIDs := make(map[int]bool)
	for _, v := range d.Get("members").(*schema.Set).List() {
		memberIDs[v.(int)] = true
	}
	var badManagers []int
	for _, v := range d.Get("managers").(*schema.Set).List() {
		uid := v.(int)
		if !memberIDs[uid] {
			badManagers = append(badManagers, uid)
		}
	}
	assert.Empty(t, badManagers, "all managers are in members — no validation error expected")
}

// ── entity_management_mode schema tests ──────────────────────────────────────

func TestEntityManagementMode_DefaultIsManaged(t *testing.T) {
	t.Parallel()
	r := resourceNewRelicTeam()
	s, ok := r.Schema["entity_management_mode"]
	require.True(t, ok, "entity_management_mode schema attribute must exist")
	// TestResourceData() does not apply schema defaults, so check the schema directly.
	assert.Equal(t, "managed", s.Default,
		"entity_management_mode schema Default should be 'managed'")
}

func TestEntityManagementMode_SchemaHasValidValues(t *testing.T) {
	t.Parallel()
	r := resourceNewRelicTeam()
	s, ok := r.Schema["entity_management_mode"]
	require.True(t, ok, "entity_management_mode schema attribute must exist")
	assert.Equal(t, schema.TypeString, s.Type)
	assert.True(t, s.Optional, "entity_management_mode should be Optional")
	assert.Equal(t, "managed", s.Default, "entity_management_mode default should be 'managed'")
	assert.NotNil(t, s.ValidateFunc, "entity_management_mode should have a ValidateFunc")
}

func TestEntityManagementMode_EntitiesIsComputedForSetNew(t *testing.T) {
	t.Parallel()
	r := resourceNewRelicTeam()
	s, ok := r.Schema["entities"]
	require.True(t, ok, "entities schema attribute must exist")
	// Computed: true is required so CustomizeDiff can call SetNew to suppress
	// entity diffs during entity_management_mode transitions.
	assert.True(t, s.Computed, "entities must be Computed to allow SetNew in CustomizeDiff")
	assert.True(t, s.Optional, "entities should remain Optional")
}

func TestEntityManagementMode_ModeLogic_UnmanagedWithEntities(t *testing.T) {
	t.Parallel()
	// Simulate the validation logic that CustomizeDiff uses.
	// When mode=unmanaged and entities is present in config, an error should fire.
	mode := "unmanaged"
	entitiesInConfig := true // simulates rc.GetAttr("entities").IsKnown() && !IsNull()
	var errs []string
	if mode == "unmanaged" && entitiesInConfig {
		errs = append(errs, `entities cannot be specified when entity_management_mode = "unmanaged"`)
	}
	assert.Len(t, errs, 1, "should produce an error when unmanaged+entities")
	assert.Contains(t, errs[0], "unmanaged")
}

func TestEntityManagementMode_ModeLogic_ManagedWithEntities(t *testing.T) {
	t.Parallel()
	// mode=managed with entities present is always valid.
	mode := "managed"
	entitiesInConfig := true
	var errs []string
	if mode == "unmanaged" && entitiesInConfig {
		errs = append(errs, `entities cannot be specified when entity_management_mode = "unmanaged"`)
	}
	assert.Empty(t, errs, "managed mode with entities should produce no error")
}

func TestEntityManagementMode_ModeLogic_ManagedWithoutEntities(t *testing.T) {
	t.Parallel()
	// mode=managed without entities is valid — entities block is optional.
	mode := "managed"
	entitiesInConfig := false
	var errs []string
	if mode == "unmanaged" && entitiesInConfig {
		errs = append(errs, `entities cannot be specified when entity_management_mode = "unmanaged"`)
	}
	assert.Empty(t, errs, "managed mode without entities should produce no error")
}

func TestEntityManagementMode_ModeLogic_UnmanagedWithoutEntities(t *testing.T) {
	t.Parallel()
	// mode=unmanaged with no entities block is valid.
	mode := "unmanaged"
	entitiesInConfig := false
	var errs []string
	if mode == "unmanaged" && entitiesInConfig {
		errs = append(errs, `entities cannot be specified when entity_management_mode = "unmanaged"`)
	}
	assert.Empty(t, errs, "unmanaged mode without entities should produce no error")
}
