package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

// dataSourceNewRelicTeamsHierarchyLevels lists all Teams hierarchy level
// entities in the org so that users can reference their IDs without having
// to look them up manually. Hierarchy levels are created automatically by
// the entityManagement API when parentId relationships are established between
// Team resources.
func dataSourceNewRelicTeamsHierarchyLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicTeamsHierarchyLevelsRead,
		Schema: map[string]*schema.Schema{
			"levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All Teams hierarchy level entities in the organisation, sorted by depth (index 0 = root/shallowest, last index = deepest).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GUID of this hierarchy level entity.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the level (e.g. 'Division', 'Squad').",
						},
					},
				},
			},
		},
	}
}

func dataSourceNewRelicTeamsHierarchyLevelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	// Fetch org settings first — its HierarchyLevelOrder provides the GUIDs in
	// canonical depth order (index 0 = shallowest/root). We then fetch each
	// level entity individually to retrieve its display name, iterating in that
	// same order so no sorting or map-lookup is needed.
	orgSettings, err := client.Scorecards.GetTeamsOrganizationSettingsWithContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var levels []map[string]interface{}
	if orgSettings != nil {
		for _, id := range orgSettings.HierarchyLevelOrder {
			entityIface, entityErr := client.Scorecards.GetEntityWithContext(ctx, id)
			if entityErr != nil || entityIface == nil || *entityIface == nil {
				continue
			}
			if level, ok := (*entityIface).(*servicearchintelligence.EntityManagementTeamsHierarchyLevelEntity); ok {
				levels = append(levels, map[string]interface{}{
					"id":   level.ID,
					"name": level.Name,
				})
			}
		}
	}

	_ = d.Set("levels", levels)
	d.SetId("teams-hierarchy-levels")
	return nil
}
