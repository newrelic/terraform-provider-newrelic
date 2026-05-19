package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicFleetMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicFleetMembersRead,
		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the fleet to list members for.",
			},
			"ring": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter members by ring name. If omitted, all members across all rings are returned.",
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of member entities in the fleet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The entity GUID of the fleet member.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the fleet member entity.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The entity type of the fleet member.",
						},
					},
				},
			},
		},
	}
}

// dataSourceNewRelicFleetMembersRead fetches the current member list from the
// API and stores it in the computed "members" attribute.
//
// It operates in two modes depending on whether "ring" is set:
//
//   - Unfiltered (ring = ""): returns every entity across all rings in the
//     fleet. Useful for a complete inventory of who belongs to a fleet.
//   - Ring-filtered (ring = "default", "canary", etc.): returns only the
//     entities assigned to that specific ring. Useful when you need to know
//     exactly which entities are in a given rollout tier.
//
// The resource ID is set to fleetID for the unfiltered case, and to
// "fleetID:ring" for the filtered case. This ensures that two data source
// blocks targeting the same fleet but different rings can coexist in the
// same configuration without colliding on the same Terraform resource ID.
func dataSourceNewRelicFleetMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)

	fetchedMembers, err := listFleetMembers(ctx, client, fleetID, ring)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading fleet members: %w", err))
	}

	// Convert the API result slice into the []interface{} form that the
	// Terraform SDK expects when setting a TypeList[schema.Resource] attribute.
	members := make([]interface{}, len(fetchedMembers))
	for i, item := range fetchedMembers {
		members[i] = map[string]interface{}{
			"id":   item.ID,
			"name": item.Name,
			"type": item.Type,
		}
	}

	if err := d.Set("members", members); err != nil {
		return diag.FromErr(err)
	}

	// Construct a stable, unique ID. Using just fleetID would cause two data
	// source instances on the same fleet (one filtered, one not) to share the
	// same ID, which confuses Terraform's state tracking.
	id := fleetID
	if ring != "" {
		id = fmt.Sprintf("%s:%s", fleetID, ring)
	}
	d.SetId(id)

	return nil
}
