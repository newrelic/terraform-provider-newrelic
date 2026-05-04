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

	var diags diag.Diagnostics

	// Check for pre-existing members in this ring before adding. These entities
	// were added outside Terraform (e.g. via the CLI instrumentation flow) and
	// are not in the caller's config. After Create, Read will sync them into
	// state, and a subsequent plan will show them as pending removal unless the
	// caller adds them to entity_ids.
	existing, err := getAllFleetMembersInRing(ctx, client, fleetID, ring)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking existing fleet members: %w", err))
	}

	configSet := stringSliceToSet(entityIDs)
	var unmanaged []string
	for _, id := range existing {
		if !configSet[id] {
			unmanaged = append(unmanaged, id)
		}
	}

	if len(unmanaged) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Pre-existing fleet members not in configuration",
			Detail: fmt.Sprintf(
				"The following entities already exist in fleet %q ring %q but are not in your entity_ids configuration. "+
					"They will be included in Terraform state after this apply. Add them to entity_ids to manage them "+
					"with Terraform, or they will be removed from the fleet on the next apply:\n  - %s",
				fleetID, ring, strings.Join(unmanaged, "\n  - "),
			),
		})
	}

	// Add entities from config. The API is idempotent for entities already present.
	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{Ring: ring, EntityIds: entityIDs},
	}

	_, err = client.FleetControl.FleetControlAddFleetMembersWithContext(ctx, fleetID, members)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("error adding fleet members: %w", err))...)
	}

	d.SetId(fmt.Sprintf("%s:%s", fleetID, ring))

	return append(diags, resourceNewRelicFleetMembersRead(ctx, d, meta)...)
}

func resourceNewRelicFleetMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)

	apiIDs, err := getAllFleetMembersInRing(ctx, client, fleetID, ring)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading fleet members: %w", err))
	}

	var diags diag.Diagnostics

	// Drift detection — only meaningful after initial creation (when state has
	// a prior entity_ids set to compare against).
	if d.Id() != "" {
		priorSet := d.Get("entity_ids").(*schema.Set)
		apiSet := stringSliceToSet(apiIDs)

		// Members that were in state but have since been removed outside Terraform.
		// These will be re-added on the next apply.
		var removedExternally []string
		for _, id := range expandStringSet(priorSet) {
			if !apiSet[id] {
				removedExternally = append(removedExternally, id)
			}
		}
		if len(removedExternally) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Fleet members removed outside Terraform",
				Detail: fmt.Sprintf(
					"The following entities were removed from fleet %q ring %q outside of Terraform. "+
						"Run terraform apply to re-add them, or remove them from entity_ids if the removal was intentional:\n  - %s",
					fleetID, ring, strings.Join(removedExternally, "\n  - "),
				),
			})
		}

		// Members present in the API but not in prior state — added outside
		// Terraform (e.g. via CLI instrumentation). They will be included in
		// state now and removed on the next apply unless added to entity_ids.
		priorStateSet := stringSliceToSet(expandStringSet(priorSet))
		var addedExternally []string
		for _, id := range apiIDs {
			if !priorStateSet[id] {
				addedExternally = append(addedExternally, id)
			}
		}
		if len(addedExternally) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Fleet members added outside Terraform",
				Detail: fmt.Sprintf(
					"The following entities were added to fleet %q ring %q outside of Terraform. "+
						"Add them to entity_ids to manage them with Terraform, or they will be removed from the fleet on the next apply:\n  - %s",
					fleetID, ring, strings.Join(addedExternally, "\n  - "),
				),
			})
		}
	}

	members := make([]interface{}, len(apiIDs))
	for i, id := range apiIDs {
		members[i] = id
	}

	if err := d.Set("entity_ids", schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), members)); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
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

// stringSliceToSet converts a string slice to a map for O(1) lookup.
func stringSliceToSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
