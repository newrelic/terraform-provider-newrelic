package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// resolveUserIDsToNGEPGUIDs takes a list of integer NR user IDs and returns
// the corresponding EntityManagementUserEntity GUIDs via the standard
// actor.entitySearch API.
//
// The lookup uses the query:
//
//	domain = 'NGEP' AND type = 'USER' AND tags.userId in (id1, id2, ...)
//
// as described by the Teams platform team (see Confluence Part-II comments).
// The tagFilter option is included so the response contains userId tags for
// reverse-mapping.
func resolveUserIDsToNGEPGUIDs(ctx context.Context, client *entities.Entities, userIDs []int) (map[int]string, error) {
	if len(userIDs) == 0 {
		return map[int]string{}, nil
	}

	// Build the IN clause: tags.userId in (1, 2, 3)
	idStrs := make([]string, len(userIDs))
	for i, id := range userIDs {
		idStrs[i] = fmt.Sprintf("%d", id)
	}
	query := fmt.Sprintf("domain = 'NGEP' AND type = 'USER' AND tags.userId in (%s)", strings.Join(idStrs, ", "))

	opts := entities.EntitySearchOptions{
		TagFilter: []string{"userId"},
	}
	results, err := client.GetEntitySearchByQueryWithContext(ctx, opts, query, []entities.EntitySearchSortCriteria{})
	if err != nil {
		return nil, fmt.Errorf("looking up NGEP user entities: %w", err)
	}

	guidByUserID := make(map[int]string, len(userIDs))
	if results == nil {
		return guidByUserID, nil
	}
	for _, e := range results.Results.Entities {
		// Each entity outline has a GUID; we extract the userId from its tags.
		guidByUserID[entityUserIDTag(e)] = string(e.GetGUID())
	}

	// Warn about any IDs we couldn't resolve (user doesn't exist as NGEP entity).
	var missing []string
	for _, id := range userIDs {
		if _, ok := guidByUserID[id]; !ok {
			missing = append(missing, fmt.Sprintf("%d", id))
		}
	}
	if len(missing) > 0 {
		return guidByUserID, fmt.Errorf(
			"the following user IDs could not be resolved to NGEP UserEntity GUIDs — "+
				"ensure the users exist in the organization: %s", strings.Join(missing, ", "))
	}
	return guidByUserID, nil
}

// entityUserIDTag extracts the integer userId from an entity outline's tags.
// The standard entities.EntityOutlineInterface does not expose tags as a typed
// field, so we check via a local interface assertion.
func entityUserIDTag(e entities.EntityOutlineInterface) int {
	type taggedEntity interface {
		GetTags() []entities.EntityTag
	}
	t, ok := e.(taggedEntity)
	if !ok {
		return 0
	}
	for _, tag := range t.GetTags() {
		if tag.Key == "userId" && len(tag.Values) > 0 {
			var id int
			fmt.Sscanf(tag.Values[0], "%d", &id)
			return id
		}
	}
	return 0
}

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

// flattenTeamTags converts EntityManagementTag values back to "key:value1,value2" strings.
// System-managed tags injected by NGEP (keys starting with "nr.") are filtered
// out — for example "nr.hierarchy.level" is added automatically when a team is
// given a parent_id and must not appear in state.
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

// expandTeamResources converts the Terraform resources list to API inputs.
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
		if title, ok := m["title"].(string); ok {
			res.Title = title
		}
		out = append(out, res)
	}
	return out
}

// expandTeamResourcesUpdate converts the Terraform resources list to Update API inputs.
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
		if title, ok := m["title"].(string); ok {
			res.Title = title
		}
		out = append(out, res)
	}
	return out
}

// flattenTeamResources converts API resources back to Terraform map list.
func flattenTeamResources(res []scorecards.EntityManagementTeamResource) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(res))
	for _, r := range res {
		m := map[string]interface{}{
			"type":    r.Type,
			"content": r.Content,
			"title":   r.Title,
		}
		out = append(out, m)
	}
	return out
}
