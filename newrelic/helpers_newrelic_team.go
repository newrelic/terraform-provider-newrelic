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

// readStaticOwnershipGUIDs returns the GUIDs of manually-added entities in the
// ownership collection and, separately, the GUIDs of entities that are managed
// by NGEP's tag-based discovery rules.
//
// NGEP places two kinds of entities in the ownership collection:
//  1. Manually added (via Terraform or the UI) — returned in staticGUIDs.
//  2. Auto-assigned via tag discovery — NOT in staticGUIDs, their GUIDs are
//     returned in discoveryGUIDs so callers can surface them as warnings.
//
// Tag-matched entities may be returned as EntityManagementGenericEntity (APM,
// etc.) or EntityManagementFleetEntity (Fleet entities). The Fleet type is not
// registered in the Go-client unmarshal switch and therefore does not appear in
// the collection read at all. To find these, we run a separate entity search
// for ALL entities with a tag matching the team's name/aliases — that covers
// both GenericEntity and Fleet types.
//
// Callers use staticGUIDs to set d.Set("entities", ...) so that Terraform only
// tracks manually-added entities and shows drift when one is added out-of-band.
// discoveryGUIDs are used solely for the non-authoritative warning message.
//
// If discovery is disabled or orgSettings is nil, all collection entities are
// treated as static and discoveryGUIDs is empty.
func readStaticOwnershipGUIDs(
	ctx context.Context,
	scClient *scorecards.Scorecards,
	entitiesClient *entities.Entities,
	ownershipColID string,
	teamName string,
	teamAliases []string,
	orgSettings *scorecards.EntityManagementTeamsOrganizationSettingsEntity,
) (staticGUIDs []string, discoveryGUIDs []string, err error) {
	allGUIDs, err := readTeamOwnedEntityGUIDs(ctx, scClient, ownershipColID)
	if err != nil {
		return nil, nil, err
	}

	// If discovery is not configured or disabled, all collection entities are
	// static and there are no tag-matched entities to warn about.
	if orgSettings == nil || !orgSettings.Discovery.Enabled || len(orgSettings.Discovery.TagKeys) == 0 {
		return allGUIDs, nil, nil
	}

	// Build the tag-value IN (...) clause for the team's name and all aliases.
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

	// Paginated discovery search — no id IN filter, finds all entities with
	// matching team tags across the org regardless of entity type (including
	// Fleet entities that never appear in the collection collectionElements
	// response). Capped at maxEntitySearchPages to bound API call count.
	allDiscoveryGUIDs, truncated, err := entitiesClient.GetAllEntitySearchGUIDsByQueryWithContext(
		ctx,
		discoveryTagFilter,
	)
	if err != nil {
		// Non-fatal: fall back to treating all collection entities as static.
		return allGUIDs, nil, nil
	}
	if len(allDiscoveryGUIDs) == 0 && !truncated {
		// No tag-matched entities — all collection entities are static.
		return allGUIDs, nil, nil
	}

	// Build the discovery set. If the search was truncated, callers receive the
	// partial list and should note the truncation in any warning they emit.
	discoverySet := make(map[string]bool, len(allDiscoveryGUIDs))
	for _, guid := range allDiscoveryGUIDs {
		discoverySet[guid] = true
		discoveryGUIDs = append(discoveryGUIDs, guid)
	}
	if truncated {
		// Sentinel value appended so callers can detect truncation without an
		// extra return variable in the existing (staticGUIDs, discoveryGUIDs)
		// signature. Read function checks for this and appends a note.
		discoveryGUIDs = append(discoveryGUIDs, "__truncated__")
	}

	// Static = collection entities that are NOT in the discovery set.
	for _, g := range allGUIDs {
		if !discoverySet[g] {
			staticGUIDs = append(staticGUIDs, g)
		}
	}
	return staticGUIDs, discoveryGUIDs, nil
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
