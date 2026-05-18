package newrelic

// This resource implements an opt-in ring membership model: it only tracks
// entities that are explicitly listed in entity_ids. Entities that join a
// fleet ring through other means (e.g. Agent Control instrumentation) are
// invisible to Terraform unless the user deliberately adds their GUIDs to
// entity_ids, at which point they are "adopted" and come under full lifecycle
// management.

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func resourceNewRelicFleetMembers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetMembersCreate,
		ReadContext:   resourceNewRelicFleetMembersRead,
		UpdateContext: resourceNewRelicFleetMembersUpdate,
		DeleteContext: resourceNewRelicFleetMembersDelete,
		Importer: &schema.ResourceImporter{
			// ImportStatePassthroughContext tells Terraform to use the value the
			// user passes to "terraform import" as the resource ID directly,
			// without any extra transformation. For this resource that value is
			// the fleet GUID.
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew means that if the fleet_id ever changes, Terraform
				// will destroy the existing resource and create a new one rather
				// than trying to update in place. A fleet_members resource is
				// always scoped to exactly one fleet.
				ForceNew:    true,
				Description: "The GUID of the fleet to manage entity assignments for.",
			},
			"ring": {
				Type:     schema.TypeList,
				Required: true,
				Description: "One or more ring blocks. Each block declares which entities Terraform should " +
					"maintain in that ring. At least one ring block must be specified. Only rings " +
					"explicitly declared here are managed — any other rings on the fleet are left untouched.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the ring as configured on the fleet (e.g. \"default\", \"canary\").",
						},
						"entity_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "Ordered list of entity GUIDs to assign to this ring. " +
								"Only the entities listed here are tracked by Terraform; any other entities " +
								"already in the ring through other means are not affected. Removing a GUID " +
								"from this list will remove that entity from the fleet ring on the next apply.",
						},
					},
				},
			},
		},
	}
}

// resourceNewRelicFleetMembersCreate assigns entities to the declared rings
// for the first time.
//
// Before issuing the add mutation it calls resolveEntitiesToAdd for each ring,
// which checks whether any of the requested entities are already somewhere in
// the fleet (e.g. added by Agent Control or a previous partial run). Those
// entities are skipped for the mutation — the API would reject a duplicate add
// — but if they happen to already be in the correct ring they are confirmed by
// the subsequent Read and adopted into Terraform state.
//
// nil is passed for entitiesBeingRemoved because a Create never removes
// anything; that parameter only matters during Update when entities move
// between rings.
func resourceNewRelicFleetMembersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Get("fleet_id").(string)

	rings := expandFleetMemberRings(d.Get("ring").([]interface{}))

	var diags diag.Diagnostics
	// toAdd is built up across all rings and sent in a single API call so that
	// the platform can process the assignments atomically.
	var toAdd []fleetcontrol.FleetControlFleetMemberRingInput

	for _, ring := range rings {
		// nil means "nothing is being removed in this apply" — safe for Create
		// because we are not changing any existing assignments, only adding new ones.
		actualToAdd, addDiags, err := resolveEntitiesToAdd(ctx, client, fleetID, ring.name, ring.entityIDs, nil)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		diags = append(diags, addDiags...)
		if len(actualToAdd) > 0 {
			toAdd = append(toAdd, fleetcontrol.FleetControlFleetMemberRingInput{
				Ring:      ring.name,
				EntityIds: actualToAdd,
			})
		}
	}

	if len(toAdd) > 0 {
		_, err := client.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, toAdd)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("error adding fleet members: %w", err))...)
		}
	}

	// The fleet GUID doubles as the Terraform resource ID. It is stable for the
	// lifetime of the resource because fleet_id is ForceNew.
	d.SetId(fleetID)

	// Immediately read back from the API to populate state with what was
	// actually applied (e.g. confirmed entity IDs in their server-returned order).
	return append(diags, resourceNewRelicFleetMembersRead(ctx, d, meta)...)
}

// resourceNewRelicFleetMembersRead reconciles the Terraform state against the
// API, applying the opt-in management model: only entities that are already in
// state (declared by the user) are checked. Entities present in a ring but not
// in entity_ids are silently ignored — Terraform never claimed ownership of
// them and will not surface them as drift.
//
// If a declared entity is no longer in the ring according to the API, it was
// removed out-of-band. A warning is emitted so the user is aware, and the
// entity is dropped from state. The next apply will re-add it to restore the
// declared configuration.
//
// Import path (no rings in prior state): all current fleet members are queried
// and placed under a synthetic "default" ring. This gives terraform import a
// valid starting state. Users with entities spread across multiple rings should
// update the ring blocks in their configuration after importing.
func resourceNewRelicFleetMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Id()

	rings := expandFleetMemberRings(d.Get("ring").([]interface{}))

	// When rings is empty the resource was just imported: there is no prior
	// state to diff against, so we populate a baseline by fetching all members.
	if len(rings) == 0 {
		memberIDs, err := listFleetMemberIDs(ctx, client, fleetID, "")
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading fleet members on import: %w", err))
		}
		// Convert []string to []interface{} because the Terraform SDK requires
		// interface slices when calling d.Set on a TypeList of TypeString.
		var entityIDs []interface{}
		for _, id := range memberIDs {
			entityIDs = append(entityIDs, id)
		}
		if entityIDs == nil {
			entityIDs = []interface{}{}
		}
		if err := d.Set("fleet_id", fleetID); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ring", []interface{}{
			map[string]interface{}{
				"name":       "default",
				"entity_ids": entityIDs,
			},
		}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	var diags diag.Diagnostics
	var updatedRings []interface{}

	for _, ring := range rings {
		// Fetch only the members of this specific ring from the API.
		ringMemberIDs, err := listFleetMemberIDs(ctx, client, fleetID, ring.name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading members of ring %q: %w", ring.name, err))
		}
		// Build a set for O(1) membership checks instead of scanning the slice
		// for every entity in state.
		ringMemberSet := makeStringSet(ringMemberIDs)

		// Walk the declared entity list and split it into two buckets:
		//   confirmedIDs      — still present in the API; keep in state
		//   removedExternally — absent from the API; were removed out-of-band
		var confirmedIDs []interface{}
		var removedExternally []string

		for _, id := range ring.entityIDs {
			if ringMemberSet[id] {
				confirmedIDs = append(confirmedIDs, id)
			} else {
				removedExternally = append(removedExternally, id)
			}
		}

		if len(removedExternally) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Fleet members removed outside Terraform in ring %q", ring.name),
				Detail: fmt.Sprintf(
					"The following entities were removed from fleet %q ring %q outside of Terraform. "+
						"Run terraform apply to re-add them, or remove them from entity_ids if the removal was intentional:\n  - %s",
					fleetID, ring.name, strings.Join(removedExternally, "\n  - "),
				),
			})
		}

		if confirmedIDs == nil {
			confirmedIDs = []interface{}{}
		}

		updatedRings = append(updatedRings, map[string]interface{}{
			"name":       ring.name,
			"entity_ids": confirmedIDs,
		})
	}

	if updatedRings == nil {
		updatedRings = []interface{}{}
	}

	if err := d.Set("ring", updatedRings); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// resourceNewRelicFleetMembersUpdate reconciles ring membership when the
// configuration changes. It handles three scenarios in one pass:
//
//  1. Existing ring modified — diff old vs new entity_ids to find additions
//     and removals within that ring.
//  2. New ring block added — all entities in the new block are additions.
//  3. Ring block removed — all entities previously declared for that ring
//     are removals.
//
// Removes are always executed before adds. The Fleet Control API rejects
// adding an entity that is already assigned anywhere in the fleet, so when
// an entity moves from ring A to ring B the remove-from-A must complete
// before the add-to-B is attempted.
func resourceNewRelicFleetMembersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("ring") {
		return nil
	}

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Id()

	// GetChange returns the ring list as it was before this apply (old) and
	// how it should look after (new). Both are decoded into name→[]entityID maps
	// for efficient lookup.
	oldRaw, newRaw := d.GetChange("ring")
	oldByName := indexRingsByName(expandFleetMemberRings(oldRaw.([]interface{})))
	newByName := indexRingsByName(expandFleetMemberRings(newRaw.([]interface{})))

	// Build the full removal set up front so that resolveEntitiesToAdd can
	// distinguish entities that are being moved (remove from one ring, add to
	// another in the same apply) from entities that are genuinely new to the
	// fleet. Without this, a moved entity would be flagged as "already assigned"
	// and the add would be skipped.
	entitiesBeingRemoved := collectEntitiesBeingRemoved(oldByName, newByName)

	var diags diag.Diagnostics
	var toAdd []fleetcontrol.FleetControlFleetMemberRingInput
	var toRemove []fleetcontrol.FleetControlFleetMemberRingInput

	// Process each ring that exists in the new desired state.
	for ringName, newIDs := range newByName {
		newIDSet := makeStringSet(newIDs)
		var additions []string

		if oldIDs, existed := oldByName[ringName]; existed {
			// Ring existed before — diff the two entity lists.
			oldIDSet := makeStringSet(oldIDs)
			for _, id := range newIDs {
				if !oldIDSet[id] {
					additions = append(additions, id) // present in new, absent from old
				}
			}
			var removals []string
			for _, id := range oldIDs {
				if !newIDSet[id] {
					removals = append(removals, id) // present in old, absent from new
				}
			}
			if len(removals) > 0 {
				toRemove = append(toRemove, fleetcontrol.FleetControlFleetMemberRingInput{Ring: ringName, EntityIds: removals})
			}
		} else {
			// Ring is brand new — every entity in it is an addition.
			additions = newIDs
		}

		actualToAdd, addDiags, err := resolveEntitiesToAdd(ctx, client, fleetID, ringName, additions, entitiesBeingRemoved)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		diags = append(diags, addDiags...)
		if len(actualToAdd) > 0 {
			toAdd = append(toAdd, fleetcontrol.FleetControlFleetMemberRingInput{Ring: ringName, EntityIds: actualToAdd})
		}
	}

	// Process rings that have been removed from the config entirely.
	// All entities that Terraform was tracking in those rings must be removed.
	for ringName, oldIDs := range oldByName {
		if _, stillExists := newByName[ringName]; !stillExists && len(oldIDs) > 0 {
			toRemove = append(toRemove, fleetcontrol.FleetControlFleetMemberRingInput{Ring: ringName, EntityIds: oldIDs})
		}
	}

	// Execute removes BEFORE adds so that entities moving between rings are
	// cleanly unassigned before the add mutation runs. The API enforces that an
	// entity can only belong to one ring at a time.
	if len(toRemove) > 0 {
		_, err := client.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, toRemove)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("error removing fleet members: %w", err))...)
		}
	}
	if len(toAdd) > 0 {
		_, err := client.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, toAdd)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("error adding fleet members: %w", err))...)
		}
	}

	return append(diags, resourceNewRelicFleetMembersRead(ctx, d, meta)...)
}

// resourceNewRelicFleetMembersDelete removes only the entities that Terraform
// is explicitly tracking (those listed in entity_ids). Entities that joined
// the fleet ring through other means — Agent Control instrumentation, manual
// API calls, etc. — are left untouched because Terraform never claimed
// ownership of them.
func resourceNewRelicFleetMembersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Id()

	rings := expandFleetMemberRings(d.Get("ring").([]interface{}))

	var toRemove []fleetcontrol.FleetControlFleetMemberRingInput
	for _, ring := range rings {
		if len(ring.entityIDs) > 0 {
			toRemove = append(toRemove, fleetcontrol.FleetControlFleetMemberRingInput{
				Ring:      ring.name,
				EntityIds: ring.entityIDs,
			})
		}
	}

	// Nothing to remove — all rings were already empty in state.
	if len(toRemove) == 0 {
		return nil
	}

	_, err := client.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, toRemove)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error removing fleet members: %w", err))
	}

	return nil
}

// fleetMemberRing is the in-memory representation of one ring block from the
// Terraform configuration. The SDK stores ring blocks as a raw list of maps;
// expandFleetMemberRings decodes that into this typed struct to make the
// business logic easier to follow.
type fleetMemberRing struct {
	name      string   // matches the "name" attribute in the ring block, e.g. "default"
	entityIDs []string // ordered list of GUIDs from the "entity_ids" attribute
}

// indexRingsByName converts a slice of rings into a map keyed by ring name.
// This allows O(1) lookups when diffing old vs new ring state, instead of
// scanning the slice every time.
func indexRingsByName(rings []fleetMemberRing) map[string][]string {
	m := make(map[string][]string, len(rings))
	for _, r := range rings {
		m[r.name] = r.entityIDs
	}
	return m
}

// collectEntitiesBeingRemoved returns the set of entity IDs that will be
// removed from at least one ring in this apply.
//
// This is needed to support cross-ring moves in a single apply. When an entity
// moves from ring A to ring B, it appears in both the "remove from A" list and
// the "add to B" list. If resolveEntitiesToAdd checked fleet membership without
// knowing about this pending removal, it would see the entity as "already
// assigned" and skip the add — leaving the entity in limbo after A's removal
// ran. By telling the resolver which entities are on their way out, we let them
// pass through to the add list.
func collectEntitiesBeingRemoved(oldByName, newByName map[string][]string) map[string]bool {
	entitiesBeingRemoved := make(map[string]bool)
	for ringName, oldIDs := range oldByName {
		if newIDs, stillExists := newByName[ringName]; stillExists {
			// Ring still exists — mark only the entities being dropped from it.
			newIDSet := makeStringSet(newIDs)
			for _, id := range oldIDs {
				if !newIDSet[id] {
					entitiesBeingRemoved[id] = true
				}
			}
		} else {
			// Entire ring is being removed — all its entities are leaving.
			for _, id := range oldIDs {
				entitiesBeingRemoved[id] = true
			}
		}
	}
	return entitiesBeingRemoved
}

// resolveEntitiesToAdd separates a list of candidate entity IDs into two
// groups and handles the "already assigned" case:
//
//   - Entities NOT yet in the fleet → returned as the actual add list.
//   - Entities ALREADY in the fleet (but not being removed in this apply) →
//     skipped (the API would reject a duplicate add) and surfaced as a
//     warning. If the entity happens to already be in the correct ring it
//     is confirmed by the Read that follows and adopted into Terraform state.
//
// The entitiesBeingRemoved parameter lets cross-ring moves through: an entity
// that is simultaneously being removed from one ring and added to another is
// included in the add list even though it is currently "assigned". Pass nil
// when no removals are happening (e.g. on Create).
func resolveEntitiesToAdd(
	ctx context.Context, client *nr.NewRelic,
	fleetID, ringName string,
	candidates []string,
	entitiesBeingRemoved map[string]bool,
) ([]string, diag.Diagnostics, error) {
	if len(candidates) == 0 {
		return nil, nil, nil
	}

	// Fetch the full current membership of the fleet (all rings) so we can
	// detect whether any candidate is already assigned somewhere.
	currentMembers, err := listFleetMemberIDs(ctx, client, fleetID, "")
	if err != nil {
		return nil, nil, fmt.Errorf("error checking existing fleet members: %w", err)
	}
	currentMemberSet := makeStringSet(currentMembers)

	var actualToAdd []string
	var alreadyAssigned []string
	for _, id := range candidates {
		// An entity passes through to the add list if it is either:
		//   (a) not yet in the fleet at all, or
		//   (b) in the fleet but already scheduled for removal in this apply.
		if currentMemberSet[id] && !entitiesBeingRemoved[id] {
			alreadyAssigned = append(alreadyAssigned, id)
		} else {
			actualToAdd = append(actualToAdd, id)
		}
	}

	var diags diag.Diagnostics
	if len(alreadyAssigned) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Entities already assigned in fleet — skipped add for ring %q", ringName),
			Detail: fmt.Sprintf(
				"The following entities are already assigned somewhere in fleet %q. "+
					"If they are already in ring %q, they are now Terraform-managed — removing them from entity_ids will remove them from the fleet. "+
					"If they are in a different ring, remove them from that ring first and re-apply:\n  - %s",
				fleetID, ringName, strings.Join(alreadyAssigned, "\n  - "),
			),
		})
	}
	return actualToAdd, diags, nil
}

// expandFleetMemberRings decodes the Terraform SDK's raw representation of the
// ring list into a slice of typed fleetMemberRing structs.
//
// The SDK stores TypeList[schema.Resource] attributes as []interface{} where
// each element is a map[string]interface{}. This conversion is the standard
// pattern in this provider for turning that untyped SDK data into something
// the rest of the code can work with safely.
func expandFleetMemberRings(raw []interface{}) []fleetMemberRing {
	rings := make([]fleetMemberRing, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})
		rings = append(rings, fleetMemberRing{
			name:      m["name"].(string),
			entityIDs: expandStringList(m["entity_ids"].([]interface{})),
		})
	}
	return rings
}

// listFleetMembers pages through the fleet members API and collects all
// results into a single slice.
//
// The API uses cursor-based pagination: each response includes a NextCursor
// string that points to the next page. An empty NextCursor signals the last
// page. cursor is kept as a *string (pointer to string) so the first call can
// pass nil — meaning "start from the beginning" — which is distinct from an
// empty string that would be treated as an invalid cursor value by the API.
//
// Pass an empty ring to retrieve members across all rings in the fleet.
func listFleetMembers(ctx context.Context, client *nr.NewRelic, fleetID, ring string) ([]fleetcontrol.FleetControlFleetMemberEntityResult, error) {
	var all []fleetcontrol.FleetControlFleetMemberEntityResult
	var cursor *string // nil on the first request; points to the next-page token thereafter

	for {
		filter := &fleetcontrol.FleetControlFleetMembersFilterInput{FleetId: fleetID}
		if ring != "" {
			filter.Ring = ring
		}

		result, err := client.FleetControl.GetFleetMembersWithContext(ctx, cursor, filter)
		if err != nil {
			return nil, err
		}

		all = append(all, result.Items...)

		if result.NextCursor == "" {
			break // no more pages
		}
		// Store the cursor in a local variable first so we can take its address.
		// Ranging over a loop variable and taking &loopVar would give us a
		// pointer that changes on each iteration, which is a common Go pitfall.
		next := result.NextCursor
		cursor = &next
	}

	return all, nil
}

// listFleetMemberIDs is a convenience wrapper around listFleetMembers that
// returns only the entity GUIDs, discarding the name and type fields that the
// resource logic does not need. Pass an empty ring to query across all rings.
func listFleetMemberIDs(ctx context.Context, client *nr.NewRelic, fleetID, ring string) ([]string, error) {
	items, err := listFleetMembers(ctx, client, fleetID, ring)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids, nil
}

// makeStringSet converts a string slice into a map[string]bool so that
// membership checks run in O(1) time instead of O(n). This matters when
// diffing entity lists that can contain hundreds of GUIDs — using a nested
// loop would make the diff O(n²).
func makeStringSet(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, s := range items {
		m[s] = true
	}
	return m
}
