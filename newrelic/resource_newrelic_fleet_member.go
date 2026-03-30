package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func resourceNewRelicFleetMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetMemberCreate,
		ReadContext:   resourceNewRelicFleetMemberRead,
		UpdateContext: resourceNewRelicFleetMemberUpdate,
		DeleteContext: resourceNewRelicFleetMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNewRelicFleetMemberImport,
		},
		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the fleet.",
			},
			"ring": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ring name within the fleet.",
			},
			"entity_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of entity IDs to add to the fleet ring.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNewRelicFleetMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)
	entityIDs := expandStringSet(d.Get("entity_ids").(*schema.Set))

	// Build input for adding members
	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{
			Ring:      ring,
			EntityIds: entityIDs,
		},
	}

	// Add members to fleet
	result, err := providerConfig.NewClient.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, members)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create composite ID: fleet_id:ring
	d.SetId(fmt.Sprintf("%s:%s", fleetID, ring))

	// Set entity IDs from result
	if len(result.Members) > 0 {
		if err := d.Set("entity_ids", result.Members[0].EntityIds); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceNewRelicFleetMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Parse composite ID
	fleetID, ring, err := parseFleetMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Set fleet_id and ring from ID
	if err := d.Set("fleet_id", fleetID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ring", ring); err != nil {
		return diag.FromErr(err)
	}

	// Note: For now, we'll trust the state since we don't have a direct API to verify individual ring members
	// The entity_ids are already in state from create/update
	// A full Read would require querying all fleet members and filtering by ring
	// which is not efficient. For now, we'll keep the state as is.

	return nil
}

func resourceNewRelicFleetMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !d.HasChange("entity_ids") {
		return nil
	}

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)

	old, new := d.GetChange("entity_ids")
	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

	// Entities to add (in new but not in old)
	toAdd := newSet.Difference(oldSet)
	if toAdd.Len() > 0 {
		entityIDs := expandStringSet(toAdd)
		members := []fleetcontrol.FleetControlFleetMemberRingInput{
			{
				Ring:      ring,
				EntityIds: entityIDs,
			},
		}

		_, err := providerConfig.NewClient.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, members)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Entities to remove (in old but not in new)
	toRemove := oldSet.Difference(newSet)
	if toRemove.Len() > 0 {
		entityIDs := expandStringSet(toRemove)
		members := []fleetcontrol.FleetControlFleetMemberRingInput{
			{
				Ring:      ring,
				EntityIds: entityIDs,
			},
		}

		_, err := providerConfig.NewClient.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, members)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNewRelicFleetMemberRead(ctx, d, meta)
}

func resourceNewRelicFleetMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)
	entityIDs := expandStringSet(d.Get("entity_ids").(*schema.Set))

	// Remove members from fleet
	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{
			Ring:      ring,
			EntityIds: entityIDs,
		},
	}

	log.Printf("[INFO] Removing members from New Relic Fleet %s ring %s", fleetID, ring)

	_, err := providerConfig.NewClient.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, members)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetMemberImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Expected format: fleet_id:ring
	fleetID, ring, err := parseFleetMemberID(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("fleet_id", fleetID); err != nil {
		return nil, err
	}
	if err := d.Set("ring", ring); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// parseFleetMemberID parses the composite ID format "fleet_id:ring"
func parseFleetMemberID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid fleet member ID format, expected 'fleet_id:ring' but got '%s'", id)
	}
	return parts[0], parts[1], nil
}
