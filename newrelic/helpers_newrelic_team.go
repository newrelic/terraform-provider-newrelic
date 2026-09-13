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

// lookupUserNGEPGUIDs converts integer NR user IDs to their
// EntityManagementUserEntity GUIDs via the standard actor.entitySearch API:
//
//	domain = 'NGEP' AND type = 'USER' AND tags.userId in (id1, id2, …)
//
// The lookup strategy matches what the Teams UI performs when adding members
// (see Confluence "NRTF x Teams Part II" Q&A, Pablo's comment on mapping
// IAM users to EntityManagementUserEntity GUIDs).
func lookupUserNGEPGUIDs(ctx context.Context, client *entities.Entities, userIDs []int) (map[int]string, error) {
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
		return nil, fmt.Errorf("looking up NGEP UserEntity GUIDs: %w", err)
	}

	guidByUserID := make(map[int]string, len(userIDs))
	if results != nil {
		for _, e := range results.Results.Entities {
			guidByUserID[extractUserIDFromEntityTags(e)] = string(e.GetGUID())
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

// extractUserIDFromEntityTags reads the integer userId from the tags of an
// entity outline returned by actor.entitySearch with TagFilter: ["userId"].
func extractUserIDFromEntityTags(e entities.EntityOutlineInterface) int {
	type tagged interface {
		GetTags() []entities.EntityTag
	}
	t, ok := e.(tagged)
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

// ── Ownership collection: static/dynamic entity split ────────────────────────

// readStaticOwnershipGUIDs returns only the "statically managed" entity GUIDs
// from the team's ownership collection — i.e., those that do NOT match the
// organisation's tag-based discovery rules.
//
// Background: entities can land in the ownership collection via two paths:
//  1. Manually (via Terraform or the UI) — these are "static".
//  2. Automatically, because the entity has a tag matching the team's name or
//     aliases (e.g. tags.team = "platform") — these are "dynamic".
//
// The Teams UI distinguishes them using the same filter logic in
// useTeamEntityScope / ManageEntitiesRootContent.tsx:
//   - Static  = id IN (collection) AND NOT (tags.<key> IN (name, alias…))
//   - Dynamic = id IN (collection) AND     (tags.<key> IN (name, alias…))
//
// Terraform only manages the static subset (non-authoritative over dynamic
// entities). Calling this function instead of readTeamOwnedEntityGUIDs ensures
// the `entities {}` block never diffs against auto-assigned entries.
//
// If discovery is disabled for the org (orgSettings is nil or
// Discovery.Enabled is false), all collection entities are treated as static
// and the full GUID list is returned.
func readStaticOwnershipGUIDs(
	ctx context.Context,
	scClient *scorecards.Scorecards,
	entitiesClient *entities.Entities,
	ownershipColID string,
	teamName string,
	teamAliases []string,
	orgSettings *scorecards.EntityManagementTeamsOrganizationSettingsEntity,
) (staticGUIDs []string, dynamicCount int, err error) {

	// Read all GUIDs in the ownership collection first.
	allGUIDs, err := readTeamOwnedEntityGUIDs(ctx, scClient, ownershipColID)
	if err != nil {
		return nil, 0, err
	}
	if len(allGUIDs) == 0 {
		return nil, 0, nil
	}

	// If discovery is not configured / disabled, all entities are static.
	if orgSettings == nil || !orgSettings.Discovery.Enabled || len(orgSettings.Discovery.TagKeys) == 0 {
		return allGUIDs, 0, nil
	}

	// Build the id IN (...) clause.
	quotedGUIDs := make([]string, len(allGUIDs))
	for i, g := range allGUIDs {
		quotedGUIDs[i] = "'" + g + "'"
	}
	idFilter := "id IN (" + strings.Join(quotedGUIDs, ", ") + ")"

	// Build the tag value IN (...) clause using the team's name + aliases.
	discoveryValues := append([]string{teamName}, teamAliases...)
	quotedVals := make([]string, len(discoveryValues))
	for i, v := range discoveryValues {
		quotedVals[i] = "'" + v + "'"
	}
	tagValueList := "(" + strings.Join(quotedVals, ", ") + ")"

	// One filter fragment per configured discovery tag key (OR between them).
	tagFragments := make([]string, len(orgSettings.Discovery.TagKeys))
	for i, key := range orgSettings.Discovery.TagKeys {
		tagFragments[i] = fmt.Sprintf("`tags.%s` IN %s", key, tagValueList)
	}
	discoveryTagFilter := "(" + strings.Join(tagFragments, " OR ") + ")"

	// Query for DYNAMIC entities (tag-matched) — we want the COUNT and their
	// GUIDs so we can exclude them from static state and emit a warning.
	dynamicQuery := idFilter + " AND " + discoveryTagFilter
	dynamicResult, err := entitiesClient.GetEntitySearchByQueryWithContext(
		ctx,
		entities.EntitySearchOptions{},
		dynamicQuery,
		[]entities.EntitySearchSortCriteria{},
	)
	if err != nil {
		// Non-fatal: fall back to treating all as static rather than erroring.
		return allGUIDs, 0, nil
	}

	if dynamicResult == nil || len(dynamicResult.Results.Entities) == 0 {
		// No dynamic entities — all are static.
		return allGUIDs, 0, nil
	}

	dynamicCount = len(dynamicResult.Results.Entities)
	dynamicSet := make(map[string]bool, dynamicCount)
	for _, e := range dynamicResult.Results.Entities {
		dynamicSet[string(e.GetGUID())] = true
	}

	// Static = all GUIDs minus the dynamic ones.
	for _, g := range allGUIDs {
		if !dynamicSet[g] {
			staticGUIDs = append(staticGUIDs, g)
		}
	}
	return staticGUIDs, dynamicCount, nil
}

// ── Team collection sync ──────────────────────────────────────────────────────

// syncTeamMembership reconciles the team's membership collection. Added user
// IDs are resolved to NGEP UserEntity GUIDs before calling
// AddCollectionMembers; removed ones are resolved and removed. Only the delta
// is sent — the full list is never replaced wholesale.
func syncTeamMembership(ctx context.Context, client *nr.NewRelic, membershipColID string, oldUserIDs, newUserIDs []int) error {
	toAdd, toRemove := intSetDelta(oldUserIDs, newUserIDs)

	if len(toAdd) > 0 {
		guids, err := lookupUserNGEPGUIDs(ctx, &client.Entities, toAdd)
		if err != nil {
			return err
		}
		addGUIDs := make([]string, 0, len(guids))
		for _, g := range guids {
			addGUIDs = append(addGUIDs, g)
		}
		if _, err := client.Scorecards.EntityManagementAddCollectionMembers(membershipColID, addGUIDs); err != nil {
			return fmt.Errorf("adding members to team membership collection %s: %w", membershipColID, err)
		}
	}

	if len(toRemove) > 0 {
		guids, err := lookupUserNGEPGUIDs(ctx, &client.Entities, toRemove)
		if err != nil {
			return err
		}
		removeGUIDs := make([]string, 0, len(guids))
		for _, g := range guids {
			removeGUIDs = append(removeGUIDs, g)
		}
		if _, err := client.Scorecards.EntityManagementRemoveCollectionMembers(membershipColID, removeGUIDs); err != nil {
			return fmt.Errorf("removing members from team membership collection %s: %w", membershipColID, err)
		}
	}
	return nil
}

// syncTeamManagers resolves the declared manager user IDs to NGEP GUIDs and
// calls entityManagementUpdateTeam with the full desired managers list.
// An empty slice explicitly clears all current managers.
func syncTeamManagers(ctx context.Context, client *nr.NewRelic, teamID string, managerUserIDs []int) error {
	if len(managerUserIDs) == 0 {
		_, err := client.Scorecards.EntityManagementUpdateTeam(teamID,
			scorecards.EntityManagementTeamEntityUpdateInput{Managers: []string{}})
		return err
	}

	guids, err := lookupUserNGEPGUIDs(ctx, &client.Entities, managerUserIDs)
	if err != nil {
		return fmt.Errorf("resolving team manager user IDs to NGEP GUIDs: %w", err)
	}
	managerGUIDs := make([]string, 0, len(guids))
	for _, g := range guids {
		managerGUIDs = append(managerGUIDs, g)
	}
	_, err = client.Scorecards.EntityManagementUpdateTeam(teamID,
		scorecards.EntityManagementTeamEntityUpdateInput{Managers: managerGUIDs})
	return err
}

// syncTeamOwnership reconciles the team's ownership collection using only the
// delta between old and new entity GUIDs.
func syncTeamOwnership(ctx context.Context, client *scorecards.Scorecards, ownershipColID string, oldGUIDs, newGUIDs []string) error {
	toAdd, toRemove := stringSetDelta(oldGUIDs, newGUIDs)

	if len(toAdd) > 0 {
		if _, err := client.EntityManagementAddCollectionMembers(ownershipColID, toAdd); err != nil {
			return fmt.Errorf("adding entities to team ownership collection %s: %w", ownershipColID, err)
		}
	}
	if len(toRemove) > 0 {
		if _, err := client.EntityManagementRemoveCollectionMembers(ownershipColID, toRemove); err != nil {
			return fmt.Errorf("removing entities from team ownership collection %s: %w", ownershipColID, err)
		}
	}
	return nil
}

// ── Team collection readers ───────────────────────────────────────────────────

// readTeamMembershipMap pages through the team's membership collection and
// returns a GUID→userID map for every EntityManagementUserEntity it contains.
// This dual-purpose map is used both to populate the members block in state
// and to decode the team's manager GUIDs back to integer userIDs for
// idempotent round-tripping (avoiding a second API call for managers).
func readTeamMembershipMap(ctx context.Context, client *scorecards.Scorecards, membershipColID string) (map[string]int, error) {
	guidToUserID := make(map[string]int)
	err := pageCollectionItems(ctx, client, membershipColID, func(item scorecards.EntityManagementEntityInterface) {
		if u, ok := item.(*scorecards.EntityManagementUserEntity); ok {
			guidToUserID[u.ID] = u.UserID
		}
	})
	if err != nil {
		return nil, err
	}
	return guidToUserID, nil
}

// readTeamOwnedEntityGUIDs pages through the team's ownership collection and
// returns the GUID of every entity it contains, regardless of entity type.
func readTeamOwnedEntityGUIDs(ctx context.Context, client *scorecards.Scorecards, ownershipColID string) ([]string, error) {
	var guids []string
	err := pageCollectionItems(ctx, client, ownershipColID, func(item scorecards.EntityManagementEntityInterface) {
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
	})
	if err != nil {
		return nil, err
	}
	return guids, nil
}

// ── Collection orchestration ──────────────────────────────────────────────────

// applyTeamCollections reconciles all three of a team's collection-backed
// attributes (members, managers, entities) in the correct dependency order.
// It is called from both Create and Update so neither duplicates this logic.
//
// Pass nil/empty slices for oldMembers and oldEntities during Create — nothing
// existed before the entity was created. During Update, pass the previous
// state values so only the actual delta is applied.
//
// Managers receive the full desired list on every call (NGEP replaces the
// whole list, so there is no meaningful old/new delta for them). They must
// be applied after membership is reconciled because NGEP validates that every
// manager is already a collection member.
func applyTeamCollections(
	ctx context.Context,
	client *nr.NewRelic,
	teamID, membershipColID, ownershipColID string,
	oldMembers, newMembers []int,
	newManagers []int,
	oldEntities, newEntities []string,
) error {
	if err := syncTeamMembership(ctx, client, membershipColID, oldMembers, newMembers); err != nil {
		return err
	}
	if err := syncTeamManagers(ctx, client, teamID, newManagers); err != nil {
		return err
	}
	if err := syncTeamOwnership(ctx, &client.Scorecards, ownershipColID, oldEntities, newEntities); err != nil {
		return err
	}
	return nil
}

// ── Miscellaneous team helpers ────────────────────────────────────────────────

// clearTeamTagsRaw sends tags: [] via a raw NerdGraph call to clear all user
// tags on a team that has no system-managed tags. Required because the
// generated struct uses omitempty on Tags, which would silently drop an empty
// slice — preventing the user from clearing tags via the normal struct path.
func clearTeamTagsRaw(ctx context.Context, client *nr.NewRelic, teamID string) error {
	const q = `mutation($id: ID!, $teamEntity: EntityManagementTeamEntityUpdateInput!) {
  entityManagementUpdateTeam(id: $id, teamEntity: $teamEntity) {
    entity { id tags { key values } }
  }
}`
	_, err := client.NerdGraph.QueryWithContext(ctx, q, map[string]interface{}{
		"id":         teamID,
		"teamEntity": map[string]interface{}{"tags": []interface{}{}},
	})
	return err
}

// clearTeamParentID sends an explicit parentId: null to NGEP to detach a team
// from its parent hierarchy. The generated EntityManagementTeamEntityUpdateInput
// uses json:"parentId,omitempty" which silently drops an empty string; this
// function bypasses that by issuing a raw NerdGraph call with null.
func clearTeamParentID(ctx context.Context, client *nr.NewRelic, teamID string) error {
	const mutation = `mutation($id: ID!, $teamEntity: EntityManagementTeamEntityUpdateInput!) {
  entityManagementUpdateTeam(id: $id, teamEntity: $teamEntity) {
    entity { id parentId }
  }
}`
	_, err := client.NerdGraph.QueryWithContext(ctx, mutation, map[string]interface{}{
		"id":         teamID,
		"teamEntity": map[string]interface{}{"parentId": nil},
	})
	return err
}
