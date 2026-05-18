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

func dataSourceNewRelicFleetMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	fleetID := d.Get("fleet_id").(string)
	ring := d.Get("ring").(string)

	allItems, err := queryFleetMembers(ctx, client, fleetID, ring)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading fleet members: %w", err))
	}

	members := make([]interface{}, len(allItems))
	for i, item := range allItems {
		members[i] = map[string]interface{}{
			"id":   item.ID,
			"name": item.Name,
			"type": item.Type,
		}
	}

	if err := d.Set("members", members); err != nil {
		return diag.FromErr(err)
	}

	id := fleetID
	if ring != "" {
		id = fmt.Sprintf("%s:%s", fleetID, ring)
	}
	d.SetId(id)

	return nil
}
