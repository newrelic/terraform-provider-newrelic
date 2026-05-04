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
			StateContext: resourceNewRelicFleetMembersImportState,
		},
		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GUID of the fleet to manage members for.",
			},
			"ring": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ring name to manage members in (e.g. \"default\", \"canary\", \"production\").",
			},
			"entity_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of entity GUIDs to add as members of the fleet ring.",
			},
		},
	}
}

func resourceNewRelicFleetMembersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)
	entityIDs := expandStringSet(d.Get("entity_ids").(*schema.Set))

	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{Ring: ring, EntityIds: entityIDs},
	}

	_, err := client.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, members)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error adding fleet members: %w", err))
	}

	d.SetId(fmt.Sprintf("%s:%s", fleetID, ring))

	return resourceNewRelicFleetMembersRead(ctx, d, meta)
}

func resourceNewRelicFleetMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)

	entityIDs, err := getAllFleetMembersInRing(ctx, client, fleetID, ring)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading fleet members: %w", err))
	}

	members := make([]interface{}, len(entityIDs))
	for i, id := range entityIDs {
		members[i] = id
	}

	if err := d.Set("entity_ids", schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), members)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetMembersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	if !d.HasChange("entity_ids") {
		return nil
	}

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)

	oldRaw, newRaw := d.GetChange("entity_ids")
	oldSet := oldRaw.(*schema.Set)
	newSet := newRaw.(*schema.Set)

	toAdd := expandStringSet(newSet.Difference(oldSet))
	toRemove := expandStringSet(oldSet.Difference(newSet))

	if len(toAdd) > 0 {
		members := []fleetcontrol.FleetControlFleetMemberRingInput{
			{Ring: ring, EntityIds: toAdd},
		}
		_, err := client.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, members)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding fleet members: %w", err))
		}
	}

	if len(toRemove) > 0 {
		members := []fleetcontrol.FleetControlFleetMemberRingInput{
			{Ring: ring, EntityIds: toRemove},
		}
		_, err := client.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, members)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error removing fleet members: %w", err))
		}
	}

	return resourceNewRelicFleetMembersRead(ctx, d, meta)
}

func resourceNewRelicFleetMembersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)
	entityIDs := expandStringSet(d.Get("entity_ids").(*schema.Set))

	if len(entityIDs) == 0 {
		return nil
	}

	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{Ring: ring, EntityIds: entityIDs},
	}

	_, err := client.FleetControl.FleetControlRemoveFleetMembersWithContext(ctx, fleetID, members)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error removing fleet members: %w", err))
	}

	return nil
}

func resourceNewRelicFleetMembersImportState(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import ID %q: expected fleet_id:ring", d.Id())
	}

	if err := d.Set("fleet_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("ring", parts[1]); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// getAllFleetMembersInRing fetches all member entity IDs for the given fleet and ring,
// following the cursor-based pagination from GetFleetMembers.
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
