package newrelic

import (
	"context"
	"fmt"
	"strings"

	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── User ID ↔ NGEP GUID resolution ───────────────────────────────────────────

// resolveUserIDsToNGEPGUIDs converts integer NR user IDs to their
// EntityManagementUserEntity GUIDs via the standard actor.entitySearch API
// using the query:
//
//	domain = 'NGEP' AND type = 'USER' AND tags.userId in (id1, id2, …)
//
// This is the same lookup the Teams UI frontend performs when a user selects
// members to add to a team (see Confluence Part-II Q&A, Pablo's comment).
func resolveUserIDsToNGEPGUIDs(ctx context.Context, client *entities.Entities, userIDs []int) (map[int]string, error) {
	if len(userIDs) == 0 {
		return map[int]string{}, nil
	}

	idStrs := make([]string, len(userIDs))
	for i, id := range userIDs {
		idStrs[i] = fmt.Sprintf("%d", id)
	}
	query := fmt.Sprintf("domain = 'NGEP' AND type = 'USER' AND tags.userId in (%s)", strings.Join(idStrs, ", "))

	results, err := client.GetEntitySearchByQueryWithContext(ctx,
		entities.EntitySearchOptions{TagFilter: []string{"userId"}},
		query,
		[]entities.EntitySearchSortCriteria{})
	if err != nil {
		return nil, fmt.Errorf("looking up NGEP user entities: %w", err)
	}

	guidByUserID := make(map[int]string, len(userIDs))
	if results != nil {
		for _, e := range results.Results.Entities {
			guidByUserID[entityUserIDTag(e)] = string(e.GetGUID())
		}
	}

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

// entityUserIDTag extracts the integer userId tag from an entity outline.
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

// ── Collection sync operations ────────────────────────────────────────────────

// syncMembers reconciles the membership collection. Newly added user IDs are
// resolved to NGEP GUIDs before calling AddCollectionMembers; removed ones
// are resolved and removed. Only the delta is sent.
func syncMembers(ctx context.Context, client *nr.NewRelic, membershipColID string, oldUserIDs, newUserIDs []int) error {
	oldSet := make(map[int]bool, len(oldUserIDs))
	for _, id := range oldUserIDs {
		oldSet[id] = true
	}
	newSet := make(map[int]bool, len(newUserIDs))
	for _, id := range newUserIDs {
		newSet[id] = true
	}

	var toAdd, toRemove []int
	for id := range newSet {
		if !oldSet[id] {
			toAdd = append(toAdd, id)
		}
	}
	for id := range oldSet {
		if !newSet[id] {
			toRemove = append(toRemove, id)
		}
	}

	if len(toAdd) > 0 {
		guids, err := resolveUserIDsToNGEPGUIDs(ctx, &client.Entities, toAdd)
		if err != nil {
			return err
		}
		addList := make([]string, 0, len(guids))
		for _, g := range guids {
			addList = append(addList, g)
		}
		if _, err := client.Scorecards.EntityManagementAddCollectionMembers(membershipColID, addList); err != nil {
			return fmt.Errorf("adding members to collection %s: %w", membershipColID, err)
		}
	}

	if len(toRemove) > 0 {
		guids, err := resolveUserIDsToNGEPGUIDs(ctx, &client.Entities, toRemove)
		if err != nil {
			return err
		}
		removeList := make([]string, 0, len(guids))
		for _, g := range guids {
			removeList = append(removeList, g)
		}
		if _, err := client.Scorecards.EntityManagementRemoveCollectionMembers(membershipColID, removeList); err != nil {
			return fmt.Errorf("removing members from collection %s: %w", membershipColID, err)
		}
	}
	return nil
}

// syncManagers resolves the declared manager user IDs to NGEP GUIDs and
// updates the team's managers list. An empty slice explicitly clears all managers.
func syncManagers(ctx context.Context, client *nr.NewRelic, teamID string, managerUserIDs []int) error {
	if len(managerUserIDs) == 0 {
		_, err := client.Scorecards.EntityManagementUpdateTeam(teamID,
			scorecards.EntityManagementTeamEntityUpdateInput{Managers: []string{}})
		return err
	}

	guids, err := resolveUserIDsToNGEPGUIDs(ctx, &client.Entities, managerUserIDs)
	if err != nil {
		return fmt.Errorf("resolving manager user IDs to NGEP GUIDs: %w", err)
	}
	managerGUIDs := make([]string, 0, len(guids))
	for _, g := range guids {
		managerGUIDs = append(managerGUIDs, g)
	}
	_, err = client.Scorecards.EntityManagementUpdateTeam(teamID,
		scorecards.EntityManagementTeamEntityUpdateInput{Managers: managerGUIDs})
	return err
}

// syncOwnedEntities reconciles the ownership collection by diffing old vs new
// GUIDs and only applying the delta.
func syncOwnedEntities(_ context.Context, client *scorecards.Scorecards, ownershipColID string, oldGUIDs, newGUIDs []string) error {
	oldSet := make(map[string]bool, len(oldGUIDs))
	for _, g := range oldGUIDs {
		oldSet[g] = true
	}
	newSet := make(map[string]bool, len(newGUIDs))
	for _, g := range newGUIDs {
		newSet[g] = true
	}

	var toAdd, toRemove []string
	for g := range newSet {
		if !oldSet[g] {
			toAdd = append(toAdd, g)
		}
	}
	for g := range oldSet {
		if !newSet[g] {
			toRemove = append(toRemove, g)
		}
	}

	if len(toAdd) > 0 {
		if _, err := client.EntityManagementAddCollectionMembers(ownershipColID, toAdd); err != nil {
			return fmt.Errorf("adding entities to ownership collection %s: %w", ownershipColID, err)
		}
	}
	if len(toRemove) > 0 {
		if _, err := client.EntityManagementRemoveCollectionMembers(ownershipColID, toRemove); err != nil {
			return fmt.Errorf("removing entities from ownership collection %s: %w", ownershipColID, err)
		}
	}
	return nil
}

// ── Collection readers ────────────────────────────────────────────────────────

// readCollectionMembersMap pages through a collection and returns a map of
// NGEP entity GUID → integer userId for every EntityManagementUserEntity.
// This is used to decode the team.Managers GUID list back to integer userIds
// so managers can be stored in state and compared idempotently.
func readCollectionMembersMap(_ context.Context, client *scorecards.Scorecards, colID string) (map[string]int, error) {
	if colID == "" {
		return nil, nil
	}
	guidToUserID := make(map[string]int)
	cursor := ""
	for {
		result, err := client.GetCollectionElements(cursor,
			scorecards.EntityManagementCollectionElementsFilter{
				CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: colID},
			}, 100)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}
		for _, item := range result.Items {
			if u, ok := item.(*scorecards.EntityManagementUserEntity); ok {
				guidToUserID[u.ID] = u.UserID
			}
		}
		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}
	return guidToUserID, nil
}

// readCollectionUserIDs pages through a collection and returns the userId
// integer for every EntityManagementUserEntity it contains.
func readCollectionUserIDs(_ context.Context, client *scorecards.Scorecards, colID string) ([]int, error) {
	if colID == "" {
		return nil, nil
	}
	var userIDs []int
	cursor := ""
	for {
		result, err := client.GetCollectionElements(cursor,
			scorecards.EntityManagementCollectionElementsFilter{
				CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: colID},
			}, 100)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}
		for _, item := range result.Items {
			if u, ok := item.(*scorecards.EntityManagementUserEntity); ok {
				userIDs = append(userIDs, u.UserID)
			}
		}
		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}
	return userIDs, nil
}

// readCollectionEntityGUIDs pages through a collection and returns the ID of
// every entity it contains (any type).
func readCollectionEntityGUIDs(_ context.Context, client *scorecards.Scorecards, colID string) ([]string, error) {
	if colID == "" {
		return nil, nil
	}
	var guids []string
	cursor := ""
	for {
		result, err := client.GetCollectionElements(cursor,
			scorecards.EntityManagementCollectionElementsFilter{
				CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: colID},
			}, 100)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}
		for _, item := range result.Items {
			switch e := item.(type) {
			case *scorecards.EntityManagementGenericEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementUserEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementTeamEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementCollectionEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementScorecardEntity:
				guids = append(guids, e.ID)
			}
		}
		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}
	return guids, nil
}

// ── Miscellaneous helpers ─────────────────────────────────────────────────────

// clearTeamParentID sends an explicit null parentId to NGEP to detach a team
// from its parent. The generated EntityManagementTeamEntityUpdateInput uses
// json:"parentId,omitempty" which silently omits an empty string, so we
// bypass it with a raw NerdGraph call that sends parentId: null.
func clearTeamParentID(ctx context.Context, client *nr.NewRelic, teamID string) error {
	const q = `mutation($id: ID!, $teamEntity: EntityManagementTeamEntityUpdateInput!) {
  entityManagementUpdateTeam(id: $id, teamEntity: $teamEntity) {
    entity { id parentId }
  }
}`
	_, err := client.NerdGraph.QueryWithContext(ctx, q, map[string]interface{}{
		"id":         teamID,
		"teamEntity": map[string]interface{}{"parentId": nil},
	})
	return err
}

// isNGEPGhostNotFound detects the transient "ghost" NOT_FOUND that NGEP
// returns for freshly-created entities before they are fully indexed.
// A real deletion carries the entity id in the error prefix
// (e.g. "abc123: Entity not found."); the ghost version has an empty prefix
// (": Entity not found.").
func isNGEPGhostNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) > 0 && msg[0] == ':' && len(msg) > 2 && msg[1] == ' '
}
