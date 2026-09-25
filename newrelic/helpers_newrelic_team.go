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
			if _, err := fmt.Sscanf(tag.Values[0], "%d", &id); err != nil {
				return 0
			}
			return id
		}
	}
	return 0
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
		// Managers has omitempty in TeamEntityUpdateInput, so sending an empty
		// []string{} via the typed struct would be silently dropped by the JSON
		// encoder. Use a raw mutation to guarantee the explicit empty list is sent.
		return patchTeamField(ctx, client, teamID, "managers", []interface{}{})
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
			// Treat "entity not found in collection" as a no-op: the entity may
			// have been deleted externally (e.g. via a ForceNew recreation of a
			// scorecard), in which case NGEP already removed it from the collection.
			if !strings.Contains(err.Error(), "not found in collection") {
				return fmt.Errorf("removing entities from team ownership collection %s: %w", ownershipColID, err)
			}
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
//
// Note: tag-discovery-placed entities (e.g. FLEET entities matching the team's
// discovery tags) DO appear in the collectionElements API response but are
// unmarshalled as unknown types and silently dropped by
// UnmarshalEntityManagementEntityInterface. This means they are invisible to
// the provider and will never show as drift — which is the desired behaviour
// since discovery-managed membership should not be Terraform-authoritative.
// If a new entity type is added to UnmarshalEntityManagementEntityInterface,
// a matching case must be added here to ensure it is tracked in state.
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
		case *scorecards.EntityManagementScorecardRuleEntity:
			guids = append(guids, e.ID)
		case *scorecards.EntityManagementTeamsHierarchyLevelEntity:
			guids = append(guids, e.ID)
		case *scorecards.EntityManagementTeamsOrganizationSettingsEntity:
			guids = append(guids, e.ID)
		}
	})
	if err != nil {
		return nil, err
	}
	return guids, nil
}

// readStaticOwnershipGUIDs returns only the GUIDs of entities that were
// manually added to the ownership collection (i.e., NOT auto-assigned via
// tag-based discovery rules).
//
// NGEP's tag-based discovery places both manually-added and tag-matched entities
// in the same collectionElements response, both as EntityManagementGenericEntity.
// The only way to distinguish them is by checking whether the entity has a tag
// matching the team's name/aliases under one of the org's discovery tag keys.
//
// Entities whose tags match team name or alias → dynamic (excluded from state).
// All other ownership-collection entities       → static (tracked in state).
//
// Terraform never shows drift for dynamic entities — they are managed by NGEP's
// tag-based discovery outside Terraform's scope. If discovery is disabled or
// orgSettings is nil, all entities are treated as static.
func readStaticOwnershipGUIDs(
	ctx context.Context,
	scClient *scorecards.Scorecards,
	entitiesClient *entities.Entities,
	ownershipColID string,
	teamName string,
	teamAliases []string,
	orgSettings *scorecards.EntityManagementTeamsOrganizationSettingsEntity,
) (staticGUIDs []string, dynamicCount int, err error) {
	allGUIDs, err := readTeamOwnedEntityGUIDs(ctx, scClient, ownershipColID)
	if err != nil {
		return nil, 0, err
	}
	if len(allGUIDs) == 0 {
		return nil, 0, nil
	}

	// If discovery is not configured or disabled, all entities are static.
	if orgSettings == nil || !orgSettings.Discovery.Enabled || len(orgSettings.Discovery.TagKeys) == 0 {
		return allGUIDs, 0, nil
	}

	// Build id IN (...) filter for the batch entity search.
	quotedGUIDs := make([]string, len(allGUIDs))
	for i, g := range allGUIDs {
		quotedGUIDs[i] = "'" + g + "'"
	}
	idFilter := "id IN (" + strings.Join(quotedGUIDs, ", ") + ")"

	// Build tag value IN (...) filter using team name + aliases.
	discoveryValues := append([]string{teamName}, teamAliases...)
	quotedVals := make([]string, len(discoveryValues))
	for i, v := range discoveryValues {
		quotedVals[i] = "'" + v + "'"
	}
	tagValueList := "(" + strings.Join(quotedVals, ", ") + ")"

	tagFragments := make([]string, len(orgSettings.Discovery.TagKeys))
	for i, key := range orgSettings.Discovery.TagKeys {
		tagFragments[i] = fmt.Sprintf("`tags.%s` IN %s", key, tagValueList)
	}
	discoveryTagFilter := "(" + strings.Join(tagFragments, " OR ") + ")"

	dynamicQuery := idFilter + " AND " + discoveryTagFilter
	dynamicResult, err := entitiesClient.GetEntitySearchByQueryWithContext(
		ctx,
		entities.EntitySearchOptions{},
		dynamicQuery,
		[]entities.EntitySearchSortCriteria{},
	)
	if err != nil {
		// Non-fatal: fall back to treating all as static.
		return allGUIDs, 0, nil
	}
	if dynamicResult == nil || len(dynamicResult.Results.Entities) == 0 {
		return allGUIDs, 0, nil
	}

	dynamicCount = len(dynamicResult.Results.Entities)
	dynamicSet := make(map[string]bool, dynamicCount)
	for _, e := range dynamicResult.Results.Entities {
		dynamicSet[string(e.GetGUID())] = true
	}
	for _, g := range allGUIDs {
		if !dynamicSet[g] {
			staticGUIDs = append(staticGUIDs, g)
		}
	}
	return staticGUIDs, dynamicCount, nil
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

// patchTeamField issues a raw entityManagementUpdateTeam call that sets a
// single field to an explicit value, bypassing Go struct omitempty rules.
//
// Use this when the generated TeamEntityUpdateInput drops the field via
// omitempty but the user intends an explicit clear (e.g. description → "",
// aliases → [], tags → [], parentId → nil).
//
//	patchTeamField(ctx, client, id, "description", "")
//	patchTeamField(ctx, client, id, "aliases",     []interface{}{})
//	patchTeamField(ctx, client, id, "tags",        []interface{}{})
//	patchTeamField(ctx, client, id, "parentId",    nil)
func patchTeamField(ctx context.Context, client *nr.NewRelic, teamID, field string, value interface{}) error {
	const q = `mutation($id: ID!, $teamEntity: EntityManagementTeamEntityUpdateInput!) {
  entityManagementUpdateTeam(id: $id, teamEntity: $teamEntity) {
    entity { id }
  }
}`
	_, err := client.NerdGraph.QueryWithContext(ctx, q, map[string]interface{}{
		"id":         teamID,
		"teamEntity": map[string]interface{}{field: value},
	})
	return err
}

// Convenience wrappers around patchTeamField for the four fields that require
// explicit clearing. Each wrapper name makes the call-site self-documenting.

// clearTeamDescriptionRaw sends an explicit empty string to clear the team description field.
func clearTeamDescriptionRaw(ctx context.Context, client *nr.NewRelic, teamID string) error {
	return patchTeamField(ctx, client, teamID, "description", "")
}

// clearTeamAliasesRaw sends an explicit empty list to clear all team aliases.
func clearTeamAliasesRaw(ctx context.Context, client *nr.NewRelic, teamID string) error {
	return patchTeamField(ctx, client, teamID, "aliases", []interface{}{})
}

// clearTeamTagsRaw sends an explicit empty list to clear all user-managed team tags.
func clearTeamTagsRaw(ctx context.Context, client *nr.NewRelic, teamID string) error {
	return patchTeamField(ctx, client, teamID, "tags", []interface{}{})
}

// clearTeamParentID sends an explicit null to remove the team's parent association.
func clearTeamParentID(ctx context.Context, client *nr.NewRelic, teamID string) error {
	return patchTeamField(ctx, client, teamID, "parentId", nil)
}
