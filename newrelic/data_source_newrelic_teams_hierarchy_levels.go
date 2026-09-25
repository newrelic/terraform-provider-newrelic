package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// dataSourceNewRelicTeamsHierarchyLevels lists all Teams hierarchy level
// entities in the org so that users can reference their IDs without having
// to look them up manually. Hierarchy levels are created automatically by
// NGEP when parentId relationships are established between Team resources.
func dataSourceNewRelicTeamsHierarchyLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicTeamsHierarchyLevelsRead,
		Schema: map[string]*schema.Schema{
			"levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All Teams hierarchy level entities in the organisation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NGEP GUID of this level. Use as the id when importing newrelic_teams_hierarchy_level.",
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

	searchResult, err := client.Scorecards.GetEntitySearchWithContext(ctx, "", "type = 'TEAMS_HIERARCHY_LEVEL'")
	if err != nil {
		return diag.FromErr(err)
	}

	var levels []map[string]interface{}
	if searchResult != nil {
		for _, e := range searchResult.Entities {
			if level, ok := e.(*scorecards.EntityManagementTeamsHierarchyLevelEntity); ok {
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
