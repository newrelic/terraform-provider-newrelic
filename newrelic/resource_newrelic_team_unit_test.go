//go:build unit

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// TestExpandTeamTags covers the "key:value1,value2" → TagInput conversion.
func TestExpandTeamTags(t *testing.T) {
	t.Parallel()
	raw := []interface{}{"env:dev,staging", "team:platform"}
	tags := expandTeamTags(raw)
	require.Len(t, tags, 2)
	assert.Equal(t, "env", tags[0].Key)
	assert.Equal(t, []string{"dev", "staging"}, tags[0].Values)
	assert.Equal(t, "team", tags[1].Key)
	assert.Equal(t, []string{"platform"}, tags[1].Values)
}

func TestExpandTeamTagsEmpty(t *testing.T) {
	t.Parallel()
	assert.Nil(t, expandTeamTags(nil))
	assert.Nil(t, expandTeamTags([]interface{}{}))
}

// TestFlattenTeamTags verifies the EntityManagementTag → "key:values" conversion.
func TestFlattenTeamTags(t *testing.T) {
	t.Parallel()
	tags := []scorecards.EntityManagementTag{
		{Key: "env", Values: []string{"dev"}},
		{Key: "team", Values: []string{"platform"}},
	}
	flat := flattenTeamTags(tags)
	require.Len(t, flat, 2)
	assert.Equal(t, "env:dev", flat[0])
	assert.Equal(t, "team:platform", flat[1])
}

// TestExpandTeamResources converts a Terraform list into TeamResourceCreateInput.
func TestExpandTeamResources(t *testing.T) {
	t.Parallel()
	raw := []interface{}{
		map[string]interface{}{"type": "link", "content": "https://example.com/runbook", "title": "Runbook"},
		map[string]interface{}{"type": "link", "content": "https://example.com/wiki", "title": ""},
	}
	res := expandTeamResources(raw)
	require.Len(t, res, 2)
	assert.Equal(t, "link", res[0].Type)
	assert.Equal(t, "https://example.com/runbook", res[0].Content)
	assert.Equal(t, "Runbook", res[0].Title)
}

// TestFlattenTeamResources converts EntityManagementTeamResource back to maps.
func TestFlattenTeamResources(t *testing.T) {
	t.Parallel()
	res := []scorecards.EntityManagementTeamResource{
		{Type: "link", Content: "https://example.com", Title: "Docs"},
	}
	flat := flattenTeamResources(res)
	require.Len(t, flat, 1)
	assert.Equal(t, "link", flat[0]["type"])
	assert.Equal(t, "https://example.com", flat[0]["content"])
	assert.Equal(t, "Docs", flat[0]["title"])
}

// TestExpandMemberUserIDsFromSet verifies the TypeSet → []int expansion.
func TestExpandMemberUserIDsFromSet(t *testing.T) {
	t.Parallel()
	s := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{"user_id": {Type: schema.TypeInt}},
	}), []interface{}{
		map[string]interface{}{"user_id": 123},
		map[string]interface{}{"user_id": 456},
	})
	ids := expandMemberUserIDsFromSet(s)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, 123)
	assert.Contains(t, ids, 456)
}

// TestFlattenMemberUserIDs converts int slice → list-of-map.
func TestFlattenMemberUserIDs(t *testing.T) {
	t.Parallel()
	flat := flattenMemberUserIDs([]int{10, 20})
	require.Len(t, flat, 2)
	assert.Equal(t, 10, flat[0]["user_id"])
	assert.Equal(t, 20, flat[1]["user_id"])
}

// TestFlattenEntityGUIDs converts GUID slice → list-of-map.
func TestFlattenEntityGUIDs(t *testing.T) {
	t.Parallel()
	flat := flattenEntityGUIDs([]string{"guid-a", "guid-b"})
	require.Len(t, flat, 2)
	assert.Equal(t, "guid-a", flat[0]["guid"])
}

// TestIsNGEPGhostNotFound confirms the leading-colon heuristic.
func TestIsNGEPGhostNotFound(t *testing.T) {
	t.Parallel()
	assert.True(t, isNGEPGhostNotFound(fmt.Errorf(": Entity not found.")))
	assert.False(t, isNGEPGhostNotFound(fmt.Errorf("abc123: Entity not found.")))
	assert.False(t, isNGEPGhostNotFound(nil))
}
