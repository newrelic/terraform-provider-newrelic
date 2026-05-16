package newrelic

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
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GUID of the fleet to manage member assignments for.",
			},
			"ring": {
				Type:     schema.TypeList,
				Required: true,
				Description: "One or more ring blocks. Each block declares which entities to place in that ring. " +
					"Only rings explicitly declared here are managed — any other rings on the fleet are left untouched.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ring name (e.g. \"default\", \"canary\").",
						},
						"entity_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of entity GUIDs to assign to this ring.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicFleetMembersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Get("fleet_id").(string)

	alreadyInFleet, err := getAllFleetMembers(ctx, client, fleetID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking existing fleet members: %w", err))
	}
	assignedSet := stringSliceToSet(alreadyInFleet)

	rings := expandFleetMemberRings(d.Get("ring").([]interface{}))

	var diags diag.Diagnostics
	var toAdd []fleetcontrol.FleetControlFleetMemberRingInput

	for _, ring := range rings {
		var actualToAdd []string
		var alreadyAssigned []string

		for _, id := range ring.entityIDs {
			if assignedSet[id] {
				alreadyAssigned = append(alreadyAssigned, id)
			} else {
				actualToAdd = append(actualToAdd, id)
			}
		}

		if len(alreadyAssigned) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Entities already assigned in fleet — skipped add for ring %q", ring.name),
				Detail: fmt.Sprintf(
					"The following entities are already assigned somewhere in fleet %q. "+
						"If they are already in ring %q, they are now Terraform-managed — removing them from entity_ids will remove them from the fleet. "+
						"If they are in a different ring, remove them from that ring first and re-apply:\n  - %s",
					fleetID, ring.name, strings.Join(alreadyAssigned, "\n  - "),
				),
			})
		}

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

	d.SetId(fleetID)
	return append(diags, resourceNewRelicFleetMembersRead(ctx, d, meta)...)
}

func resourceNewRelicFleetMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Id()

	rings := expandFleetMemberRings(d.Get("ring").([]interface{}))

	// Import path: no rings declared in prior state. Query all fleet members
	// and surface them under a "default" ring so that plan -generate-config-out
	// produces valid config. For multi-ring fleets, adjust ring blocks manually.
	if len(rings) == 0 {
		apiIDs, err := getAllFleetMembers(ctx, client, fleetID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading fleet members on import: %w", err))
		}
		var entityIDs []interface{}
		for _, id := range apiIDs {
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
		apiIDs, err := getAllFleetMembersInRing(ctx, client, fleetID, ring.name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading members of ring %q: %w", ring.name, err))
		}
		apiSet := stringSliceToSet(apiIDs)

		var confirmedIDs []interface{}
		var removedExternally []string

		for _, id := range ring.entityIDs {
			if apiSet[id] {
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

func resourceNewRelicFleetMembersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("ring") {
		return nil
	}

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	fleetID := d.Id()

	oldRaw, newRaw := d.GetChange("ring")
	oldByName := ringsByName(expandFleetMemberRings(oldRaw.([]interface{})))
	newByName := ringsByName(expandFleetMemberRings(newRaw.([]interface{})))

	// Pre-compute removals so entities moving between rings are excluded from
	// the "already assigned" pre-check on the add path.
	beingRemovedInThisApply := computeBeingRemoved(oldByName, newByName)

	var diags diag.Diagnostics
	var toAdd []fleetcontrol.FleetControlFleetMemberRingInput
	var toRemove []fleetcontrol.FleetControlFleetMemberRingInput

	for ringName, newIDs := range newByName {
		newIDSet := stringSliceToSet(newIDs)
		var additions []string

		if oldIDs, existed := oldByName[ringName]; existed {
			oldIDSet := stringSliceToSet(oldIDs)
			for _, id := range newIDs {
				if !oldIDSet[id] {
					additions = append(additions, id)
				}
			}
			var removals []string
			for _, id := range oldIDs {
				if !newIDSet[id] {
					removals = append(removals, id)
				}
			}
			if len(removals) > 0 {
				toRemove = append(toRemove, fleetcontrol.FleetControlFleetMemberRingInput{Ring: ringName, EntityIds: removals})
			}
		} else {
			additions = newIDs
		}

		actualToAdd, addDiags, err := fleetMembersComputeAdds(ctx, client, fleetID, ringName, additions, beingRemovedInThisApply)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		diags = append(diags, addDiags...)
		if len(actualToAdd) > 0 {
			toAdd = append(toAdd, fleetcontrol.FleetControlFleetMemberRingInput{Ring: ringName, EntityIds: actualToAdd})
		}
	}

	// Ring blocks removed from config — remove all their previously declared entities.
	for ringName, oldIDs := range oldByName {
		if _, stillExists := newByName[ringName]; !stillExists && len(oldIDs) > 0 {
			toRemove = append(toRemove, fleetcontrol.FleetControlFleetMemberRingInput{Ring: ringName, EntityIds: oldIDs})
		}
	}

	// Execute removes BEFORE adds so that entities moving between rings are
	// cleanly unassigned before the add mutation is attempted.
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

	if len(toRemove) == 0 {
		return nil
	}

	_, err := client.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, toRemove)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error removing fleet members: %w", err))
	}

	return nil
}

type fleetMemberRing struct {
	name      string
	entityIDs []string
}

func ringsByName(rings []fleetMemberRing) map[string][]string {
	m := make(map[string][]string, len(rings))
	for _, r := range rings {
		m[r.name] = r.entityIDs
	}
	return m
}

// computeBeingRemoved returns the set of entity IDs that will be removed from
// any ring in this apply, so cross-ring moves are not blocked by the
// "already assigned" pre-check.
func computeBeingRemoved(oldByName, newByName map[string][]string) map[string]bool {
	out := make(map[string]bool)
	for ringName, oldIDs := range oldByName {
		if newIDs, stillExists := newByName[ringName]; stillExists {
			newIDSet := stringSliceToSet(newIDs)
			for _, id := range oldIDs {
				if !newIDSet[id] {
					out[id] = true
				}
			}
		} else {
			for _, id := range oldIDs {
				out[id] = true
			}
		}
	}
	return out
}

// fleetMembersComputeAdds checks which candidates are already assigned in the
// fleet, emits a warning for those that are (and are not being removed in this
// apply), and returns the subset that should be passed to the add mutation.
func fleetMembersComputeAdds(
	ctx context.Context, client *nr.NewRelic,
	fleetID, ringName string,
	candidates []string,
	beingRemovedInThisApply map[string]bool,
) ([]string, diag.Diagnostics, error) {
	if len(candidates) == 0 {
		return nil, nil, nil
	}
	alreadyInFleet, err := getAllFleetMembers(ctx, client, fleetID)
	if err != nil {
		return nil, nil, fmt.Errorf("error checking existing fleet members: %w", err)
	}
	assignedSet := stringSliceToSet(alreadyInFleet)

	var actualToAdd []string
	var alreadyAssigned []string
	for _, id := range candidates {
		if assignedSet[id] && !beingRemovedInThisApply[id] {
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

// getAllFleetMembers returns all member entity IDs across all rings of a fleet.
func getAllFleetMembers(ctx context.Context, client *nr.NewRelic, fleetID string) ([]string, error) {
	var allEntityIDs []string
	var cursor *string

	for {
		filter := &fleetcontrol.FleetControlFleetMembersFilterInput{
			FleetId: fleetID,
		}

		result, err := client.FleetControl.GetFleetMembersWithContext(ctx, cursor, filter)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			allEntityIDs = append(allEntityIDs, item.ID)
		}

		if result.NextCursor == "" {
			break
		}

		next := result.NextCursor
		cursor = &next
	}

	return allEntityIDs, nil
}

// getAllFleetMembersInRing returns all member entity IDs for a specific ring.
func getAllFleetMembersInRing(ctx context.Context, client *nr.NewRelic, fleetID, ring string) ([]string, error) {
	var allEntityIDs []string
	var cursor *string

	for {
		filter := &fleetcontrol.FleetControlFleetMembersFilterInput{
			FleetId: fleetID,
			Ring:    ring,
		}

		result, err := client.FleetControl.GetFleetMembersWithContext(ctx, cursor, filter)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			allEntityIDs = append(allEntityIDs, item.ID)
		}

		if result.NextCursor == "" {
			break
		}

		next := result.NextCursor
		cursor = &next
	}

	return allEntityIDs, nil
}

// stringSliceToSet converts a string slice to a map for O(1) lookup.
func stringSliceToSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
