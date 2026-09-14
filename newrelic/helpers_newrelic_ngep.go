package newrelic

// helpers_newrelic_ngep.go contains utilities shared across all resources that
// interact with the New Relic Entity Platform (NGEP) via the entityManagement
// GraphQL surface — currently newrelic_team and newrelic_scorecard, but
// applicable to any future entityManagement-backed resource.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// ── Tag helpers ───────────────────────────────────────────────────────────────

// mergeWithSystemTags merges user-supplied tags with any existing system tags
// (those whose keys begin with "nr.") from the entity's current state. This is
// necessary because the NGEP API rejects an update that would implicitly remove
// system-managed tags — for example, "nr.hierarchy.level" which is injected
// automatically when a Team is assigned a parentId.
//
// Pass the current entity tags from a prior Read; the function returns the
// merged list safe to send in any entityManagement*Update mutation.
func mergeWithSystemTags(
	userTags []scorecards.EntityManagementTagInput,
	currentEntityTags []scorecards.EntityManagementTag,
) []scorecards.EntityManagementTagInput {
	out := make([]scorecards.EntityManagementTagInput, 0, len(userTags))
	out = append(out, userTags...)

	for _, t := range currentEntityTags {
		if strings.HasPrefix(t.Key, "nr.") {
			out = append(out, scorecards.EntityManagementTagInput{
				Key:    t.Key,
				Values: t.Values,
			})
		}
	}
	return out
}

// ngepTagged is a local interface satisfied by every NGEP entity type that
// carries a Tags slice. Using it avoids a fragile type-switch that would need
// updating whenever a new resource type is added.
type ngepTagged interface {
	GetNGEPTags() []scorecards.EntityManagementTag
}

// fetchEntitySystemTags reads the current entity and returns only its
// system-managed tags (keys prefixed "nr."). These must be preserved in every
// tags-update mutation because the NGEP API rejects updates that would
// implicitly remove them.
//
// Falls back to direct type assertions since the generated types do not
// implement a common interface for Tags. If the entity type is not one of
// the known types, returns nil — meaning no system tags to preserve, which
// is safe (the caller will send only the user-provided tags).
func fetchEntitySystemTags(ctx context.Context, client *scorecards.Scorecards, entityID string) []scorecards.EntityManagementTag {
	iface, err := client.GetEntityWithContext(ctx, entityID)
	if err != nil || iface == nil {
		return nil
	}
	var tags []scorecards.EntityManagementTag
	switch e := (*iface).(type) {
	case *scorecards.EntityManagementTeamEntity:
		tags = e.Tags
	case *scorecards.EntityManagementScorecardEntity:
		tags = e.Tags
	case *scorecards.EntityManagementScorecardRuleEntity:
		tags = e.Tags
	default:
		return nil
	}
	return filterSystemTags(tags)
}

func filterSystemTags(tags []scorecards.EntityManagementTag) []scorecards.EntityManagementTag {
	var out []scorecards.EntityManagementTag
	for _, t := range tags {
		if strings.HasPrefix(t.Key, "nr.") {
			out = append(out, t)
		}
	}
	return out
}

// ── Collection pagination ─────────────────────────────────────────────────────

// pageCollectionItems pages through an NGEP collection and calls visit for
// every item returned. Pagination stops when NextCursor is empty. The caller
// supplies the collection ID and a visitor function; this function owns the
// retry/paging loop so individual readers don't duplicate it.
func pageCollectionItems(
	_ context.Context,
	client *scorecards.Scorecards,
	colID string,
	visit func(item scorecards.EntityManagementEntityInterface),
) error {
	if colID == "" {
		return nil
	}
	cursor := ""
	for {
		result, err := client.GetCollectionElements(cursor,
			scorecards.EntityManagementCollectionElementsFilter{
				CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: colID},
			}, 100)
		if err != nil {
			return err
		}
		if result == nil {
			break
		}
		for _, item := range result.Items {
			visit(item)
		}
		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}
	return nil
}

// ── Set delta helpers ─────────────────────────────────────────────────────────

// stringSetDelta returns the elements to add and remove when transitioning
// from oldItems to newItems. Both slices are treated as unordered sets —
// duplicates are collapsed.
func stringSetDelta(oldItems, newItems []string) (toAdd, toRemove []string) {
	oldSet := make(map[string]bool, len(oldItems))
	for _, s := range oldItems {
		oldSet[s] = true
	}
	newSet := make(map[string]bool, len(newItems))
	for _, s := range newItems {
		newSet[s] = true
	}
	for s := range newSet {
		if !oldSet[s] {
			toAdd = append(toAdd, s)
		}
	}
	for s := range oldSet {
		if !newSet[s] {
			toRemove = append(toRemove, s)
		}
	}
	return toAdd, toRemove
}

// intSetDelta is the integer equivalent of stringSetDelta.
func intSetDelta(oldItems, newItems []int) (toAdd, toRemove []int) {
	oldSet := make(map[int]bool, len(oldItems))
	for _, i := range oldItems {
		oldSet[i] = true
	}
	newSet := make(map[int]bool, len(newItems))
	for _, i := range newItems {
		newSet[i] = true
	}
	for i := range newSet {
		if !oldSet[i] {
			toAdd = append(toAdd, i)
		}
	}
	for i := range oldSet {
		if !newSet[i] {
			toRemove = append(toRemove, i)
		}
	}
	return toAdd, toRemove
}

// ── NGEP entity indexing gate ─────────────────────────────────────────────────

// waitForNGEPEntityIndexed polls the standard actor.entitySearch API until the
// entity is visible, then returns. The Teams UI performs the same poll after
// create; without it Terraform may complete Create before the entity is
// queryable, causing subsequent reads to appear empty.
//
// Blocks up to timeout. Returns a non-nil error if the entity never appears.
func waitForNGEPEntityIndexed(
	ctx context.Context,
	entitiesClient *entities.Entities,
	entityID string,
	timeout time.Duration,
) error {
	return resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		res, err := entitiesClient.GetEntitySearchByQueryWithContext(
			ctx,
			entities.EntitySearchOptions{},
			"id IN ('"+entityID+"')",
			[]entities.EntitySearchSortCriteria{},
		)
		if err != nil {
			return resource.NonRetryableError(
				fmt.Errorf("entitySearch while waiting for %s to be indexed: %w", entityID, err))
		}
		if res == nil || len(res.Results.Entities) == 0 {
			return resource.RetryableError(
				fmt.Errorf("entity %s not yet visible in entitySearch — retrying", entityID))
		}
		return nil
	})
}

// ── Error classification ──────────────────────────────────────────────────────

// isNGEPGhostNotFound detects the transient "ghost" NOT_FOUND that NGEP
// returns for freshly-created entities before they are fully indexed.
//
// A real deletion carries the entity id in the error message prefix
// (e.g. "abc123: Entity not found."); the ghost version has an empty prefix
// (": Entity not found."). Callers should retry on a ghost but propagate a
// real NOT_FOUND.
func isNGEPGhostNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) > 1 && msg[0] == ':' && msg[1] == ' '
}
