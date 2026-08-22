package newrelic

// helpers_newrelic_ngep.go contains utilities shared across all resources that
// interact with the New Relic Entity Platform (NGEP) via the entityManagement
// GraphQL surface — currently newrelic_team and newrelic_scorecard, but
// applicable to any future entityManagement-backed resource.

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

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
