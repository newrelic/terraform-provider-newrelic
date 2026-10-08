package newrelic

import (
	"context"
	"fmt"
	"strings"

	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
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
// sends the full desired managers list via EntityManagementUpdateTeam.
// An empty managerUserIDs slice sends Managers: &[]string{} which explicitly
// clears all managers — the pointer field ensures the JSON encoder includes
// the empty array rather than omitting it.
func syncTeamManagers(ctx context.Context, client *nr.NewRelic, teamID string, managerUserIDs []int) error {
	if len(managerUserIDs) == 0 {
		empty := []string{}
		_, err := client.Scorecards.EntityManagementUpdateTeam(teamID,
			servicearchintelligence.EntityManagementTeamEntityUpdateInput{Managers: &empty})
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
		servicearchintelligence.EntityManagementTeamEntityUpdateInput{Managers: &managerGUIDs})
	return err
}

// syncTeamOwnership reconciles the team's ownership collection using only the
// delta between old and new entity GUIDs.
//
// "Already belongs to collection" from NGEP is treated as success — it means
// the entity is already in the right place, so the GUID is saved to state as
// if the add call succeeded. This covers several valid scenarios:
//   - Switching from unmanaged → managed: the entity stayed in the collection
//     during the unmanaged phase and is now being re-declared.
//   - Adding a tag-discovered entity to the entities block: the entity is already
//     in the collection (placed there by discovery) but the user now wants
//     Terraform to track it. The provider records it in state and from that
//     point treats it as a statically-managed entity.
func syncTeamOwnership(ctx context.Context, client *servicearchintelligence.Scorecards, ownershipColID string, oldGUIDs, newGUIDs []string) error {
	toAdd, toRemove := stringSetDelta(oldGUIDs, newGUIDs)

	if len(toAdd) > 0 {
		if _, err := client.EntityManagementAddCollectionMembers(ownershipColID, toAdd); err != nil {
			if !strings.Contains(err.Error(), "already belongs to collection") &&
				!strings.Contains(err.Error(), "already exists in one collection") {
				return fmt.Errorf("adding entities to team ownership collection %s: %w", ownershipColID, err)
			}
			// "Already belongs" = entity is in the collection, which is the desired
			// end state. The GUID will be written to state by the caller.
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
func readTeamMembershipMap(ctx context.Context, client *servicearchintelligence.Scorecards, membershipColID string) (map[string]int, error) {
	guidToUserID := make(map[string]int)
	err := pageCollectionItems(ctx, client, membershipColID, func(item servicearchintelligence.EntityManagementEntityInterface) {
		if u, ok := item.(*servicearchintelligence.EntityManagementUserEntity); ok {
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
// readTeamOwnedEntityGUIDs returns the GUID of every entity in the team's
// ownership collection, regardless of entity type.
//
// It calls item.GetID() on each element, which is defined on the
// EntityManagementEntityInterface. Every registered NGEP entity type
// implements GetID(), and the GenericEntity fallback in
// UnmarshalEntityManagementEntityInterface ensures that any future entity
// type not yet in the switch also returns a non-nil GetID(). No type
// switch is needed — this function is unconditionally future-proof.
func readTeamOwnedEntityGUIDs(ctx context.Context, client *servicearchintelligence.Scorecards, ownershipColID string) ([]string, error) {
	var guids []string
	err := pageCollectionItems(ctx, client, ownershipColID, func(item servicearchintelligence.EntityManagementEntityInterface) {
		if id := item.GetID(); id != "" {
			guids = append(guids, id)
		}
	})
	if err != nil {
		return nil, err
	}
	return guids, nil
}

// readStaticOwnershipGUIDs separates the team's ownership collection entities
// into two groups: manually-added (staticGUIDs) and tag-discovered (discoveryGUIDs).
//
// Design: BOUNDED search, not org-wide.
//
// Both manually-added and tag-matched entities coexist in the same ownership
// collection — there is no source/origin field per item. The only way to
// distinguish them is by checking whether each collection entity has a tag
// matching the team's name or an alias under the configured discovery tag keys.
//
// We do this with a single, collection-bounded query:
//
//	id IN (<allGUIDs from collection>) AND (tags.<key> IN (<name>, <alias1>, ...))
//
// This is O(collection_size), NOT O(all_org_entities). A team's collection
// typically holds < 50 entities — this is always a single API page. Contrast
// with an org-wide tag search which can return tens of thousands of results
// and make hundreds of API calls on every terraform plan.
//
// Fleet entities (EntityManagementFleetEntity) are silently dropped by the
// Go-client's unmarshal switch and never appear in allGUIDs. They therefore
// never show as drift (correct — they cannot be managed via the entities block)
// and are not listed in the warning GUIDs. This is an acceptable limitation:
// Fleet entities are NGEP-internal, they cannot be imported into Terraform, and
// customers cannot add them to the entities block regardless.
//
// declaredGUIDs must contain the entity GUIDs currently in the Terraform
// config's entities block. These are ALWAYS treated as static regardless of
// their tags — an explicit Terraform declaration takes precedence over
// tag-based discovery classification. Without this, adding a discovery tag to
// an entity that is also in the entities block would cause a perpetual diff.
//
// If discovery is disabled or orgSettings is nil, all collection entities are
// treated as static and discoveryGUIDs is empty.
func readStaticOwnershipGUIDs(
	ctx context.Context,
	scClient *servicearchintelligence.Scorecards,
	entitiesClient *entities.Entities,
	ownershipColID string,
	teamName string,
	teamAliases []string,
	declaredGUIDs []string,
	orgSettings *servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntity,
) (staticGUIDs []string, discoveryGUIDs []string, err error) {
	allGUIDs, err := readTeamOwnedEntityGUIDs(ctx, scClient, ownershipColID)
	if err != nil {
		return nil, nil, err
	}
	if len(allGUIDs) == 0 {
		return nil, nil, nil
	}

	// Entities explicitly declared in the Terraform config are ALWAYS static.
	// A Terraform declaration trumps tag-based classification — adding a
	// discovery tag to a declared entity must not create a perpetual diff.
	declaredSet := make(map[string]bool, len(declaredGUIDs))
	for _, g := range declaredGUIDs {
		declaredSet[g] = true
	}

	// If discovery is not configured or disabled, all collection entities are
	// static — nothing to classify further.
	if orgSettings == nil || !orgSettings.Discovery.Enabled || len(orgSettings.Discovery.TagKeys) == 0 {
		return allGUIDs, nil, nil
	}

	// ── Bounded query: id IN (collection_guids) AND tags.<key> IN (name/aliases) ──
	// Hits only the entities already in the collection — O(collection_size), not
	// O(all_org_entities). Collection is typically < 50 → single API page.
	quotedGUIDs := make([]string, len(allGUIDs))
	for i, g := range allGUIDs {
		quotedGUIDs[i] = "'" + g + "'"
	}
	idFilter := "id IN (" + strings.Join(quotedGUIDs, ", ") + ")"

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
	boundedQuery := idFilter + " AND (" + strings.Join(tagFragments, " OR ") + ")"

	result, err := entitiesClient.GetEntitySearchByQueryWithContext(
		ctx,
		entities.EntitySearchOptions{},
		boundedQuery,
		[]entities.EntitySearchSortCriteria{},
	)
	if err != nil {
		return allGUIDs, nil, nil
	}

	discoverySet := make(map[string]bool)
	if result != nil {
		for _, e := range result.Results.Entities {
			guid := string(e.GetGUID())
			// Skip discovery classification for declared entities — they are
			// always managed by Terraform regardless of their tags.
			if !declaredSet[guid] {
				discoverySet[guid] = true
				discoveryGUIDs = append(discoveryGUIDs, guid)
			}
		}
	}

	// Static = collection entities that are either declared in config OR not tag-matched.
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
//
// oldManagers is the prior state. The managers mutation is skipped when both
// oldManagers and newManagers are empty — this avoids a redundant API call
// on every Create without managers and on every Update where only members or
// entities changed while managers stayed empty.
//
// Call order: (1) syncTeamMembership → (2) syncTeamManagers → (3) syncTeamOwnership.
// Steps 1→2 are dependency-ordered: NGEP rejects a manager assignment unless
// the user is already a collection member, so membership MUST be committed
// before managers are set. Step 3 (ownership) is independent of the other two.
func applyTeamCollections(
	ctx context.Context,
	client *nr.NewRelic,
	teamID, membershipColID, ownershipColID string,
	oldMembers, newMembers []int,
	oldManagers, newManagers []int,
	oldEntities, newEntities []string,
) error {
	if err := syncTeamMembership(ctx, client, membershipColID, oldMembers, newMembers); err != nil {
		return err
	}
	// Skip if both old and new manager lists are empty — nothing to set or clear.
	if len(oldManagers) > 0 || len(newManagers) > 0 {
		if err := syncTeamManagers(ctx, client, teamID, newManagers); err != nil {
			return err
		}
	}
	if err := syncTeamOwnership(ctx, &client.Scorecards, ownershipColID, oldEntities, newEntities); err != nil {
		return err
	}
	return nil
}
