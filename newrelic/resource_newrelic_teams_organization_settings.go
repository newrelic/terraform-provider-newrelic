package newrelic

// resource_newrelic_teams_organization_settings manages the singleton
// TEAMS_ORGANIZATION_SETTINGS entity for the authenticated organisation.
//
// There is exactly one of these per organisation — it is auto-created when
// Teams is first used. Terraform does not create or delete the underlying
// entity; the resource block always references the pre-existing singleton.
//
// On first `terraform apply` (without a prior import), Terraform will
// automatically locate the singleton, apply the declared configuration
// values, and emit a warning stating that the existing org settings have
// been overridden. This is the same result as a `terraform import` followed
// by `terraform apply`, so import is optional but still supported.
//
// Key capabilities managed here:
//   - discovery.enabled    — toggles tag-based automatic entity ownership
//   - discovery.tag_keys   — the tag keys used for auto-assignment (e.g. ["team"])
//   - hierarchy_levels     — ordered list of hierarchy level entities (id + name)
//   - sync_groups_enabled  — toggles automatic team creation from IdP groups
//   - sync_group_rules     — rules controlling which IdP groups create teams

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

func resourceNewRelicTeamsOrganizationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicTeamsOrgSettingsCreate,
		ReadContext:   resourceNewRelicTeamsOrgSettingsRead,
		UpdateContext: resourceNewRelicTeamsOrgSettingsUpdate,
		DeleteContext: resourceNewRelicTeamsOrgSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"discovery_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Whether tag-based entity discovery is enabled. When true, entities with a " +
					"matching tag are automatically assigned to teams.",
			},
			"discovery_tag_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Description: "Tag keys used for entity discovery (e.g. [\"team\"]). Entities with these tags " +
					"are auto-assigned to the team whose name or alias matches the tag value.",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"hierarchy_levels": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Description: "Ordered list of organisation hierarchy levels. The order here defines the visual " +
					"order in the Teams UI. Use the newrelic_teams_hierarchy_levels data source to obtain level IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "GUID of the hierarchy level entity. Obtain from the newrelic_teams_hierarchy_levels data source.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Display name of this hierarchy level (e.g. 'Division', 'Squad'). Renaming here updates the entity directly.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"sync_groups_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether automatic team creation from IdP groups is enabled.",
			},
			"sync_group_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "Rules that control which IdP groups automatically create teams. Each rule " +
					"specifies one or more match conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Conditions that a group name must satisfy for this rule to apply.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type: STARTS_WITH, ENDS_WITH, or CONTAINS.",
										ValidateFunc: validation.StringInSlice([]string{
											"STARTS_WITH", "ENDS_WITH", "CONTAINS",
										}, false),
									},
									"value": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The string to match against the group name.",
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// resourceNewRelicTeamsOrgSettingsCreate locates the pre-existing singleton,
// applies all declared configuration values, and emits a warning so the
// operator knows that existing org settings have been overridden.
//
// There is no true "create" — the underlying entity always exists. This
// function applies the same update mutation as Update, but unconditionally
// sends every field (HasChange is meaningless on a brand-new resource).
// After this call the resource behaves identically to an imported resource.
func resourceNewRelicTeamsOrgSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	// Locate the singleton entity.
	existing, err := client.Scorecards.GetTeamsOrganizationSettingsWithContext(ctx)
	if err != nil {
		return diag.Errorf("could not locate Teams organisation settings entity: %v", err)
	}
	if existing == nil {
		return diag.Errorf(
			"Teams organisation settings entity not found — ensure the Teams feature is " +
				"enabled for this New Relic organisation before managing it with Terraform.",
		)
	}

	d.SetId(existing.ID)
	log.Printf("[INFO] Found Teams organisation settings singleton %s — applying declared configuration", existing.ID)

	// Build the full update input from config. Do NOT use HasChange here —
	// on a new resource every field is "new" but HasChange may return false
	// because there is no prior state to compare against.
	tagKeys := make([]string, 0)
	for _, v := range d.Get("discovery_tag_keys").([]interface{}) {
		tagKeys = append(tagKeys, v.(string))
	}
	upd := scorecards.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{
		Discovery: &scorecards.EntityManagementDiscoverySettingsUpdateInput{
			Enabled: d.Get("discovery_enabled").(bool),
			TagKeys: tagKeys,
		},
		SyncGroups: expandSyncGroupsUpdate(
			d.Get("sync_groups_enabled").(bool),
			d.Get("sync_group_rules").([]interface{}),
		),
	}

	rawLevels := d.Get("hierarchy_levels").([]interface{})
	if len(rawLevels) > 0 {
		levelIDs := make([]string, 0, len(rawLevels))
		for _, r := range rawLevels {
			levelIDs = append(levelIDs, r.(map[string]interface{})["id"].(string))
		}
		upd.HierarchyLevelOrder = levelIDs
	}

	if _, err := client.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(d.Id(), upd); err != nil {
		return diag.Errorf("applying org settings configuration: %v", err)
	}

	// Rename hierarchy levels only when the declared name differs from the
	// current API name. Fetching each level before renaming avoids unnecessary
	// API mutations when the names are already correct (e.g. on re-apply after
	// an interrupted run or a no-op import).
	for _, r := range rawLevels {
		m := r.(map[string]interface{})
		id, wantName := m["id"].(string), m["name"].(string)

		currentName := ""
		levelIface, levelErr := client.Scorecards.GetEntityWithContext(ctx, id)
		if levelErr == nil && levelIface != nil && *levelIface != nil {
			if level, ok := (*levelIface).(*scorecards.EntityManagementTeamsHierarchyLevelEntity); ok {
				currentName = level.Name
			}
		}

		if currentName == wantName {
			continue
		}
		if _, err := client.Scorecards.EntityManagementUpdateTeamsHierarchyLevel(
			id,
			scorecards.EntityManagementTeamsHierarchyLevelEntityUpdateInput{Name: wantName},
		); err != nil {
			return diag.Errorf("renaming hierarchy level %s: %v", id, err)
		}
	}

	// Read back the live state so Terraform tracks what was applied.
	diags := resourceNewRelicTeamsOrgSettingsRead(ctx, d, meta)

	// Prepend the singleton-override warning so it appears prominently.
	return append(diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "newrelic_teams_organisation_settings: existing singleton overridden",
		Detail: "This resource manages a pre-existing singleton entity that exists in every " +
			"New Relic organisation — it was not created by Terraform. Your declared configuration " +
			"values have been applied and have overridden the previous org-level Teams settings " +
			"(discovery, sync groups, hierarchy level order).\n\n" +
			"This is the expected behaviour. Future runs of terraform plan will show any " +
			"drift between the live settings and your configuration. If you prefer an explicit " +
			"workflow, use terraform import before the first apply.",
	}}, diags...)
}

func resourceNewRelicTeamsOrgSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	entityIface, err := client.Scorecards.GetEntityWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if entityIface == nil || *entityIface == nil {
		d.SetId("")
		return nil
	}

	settings, ok := (*entityIface).(*scorecards.EntityManagementTeamsOrganizationSettingsEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamsOrganizationSettingsEntity", d.Id())
	}

	_ = d.Set("discovery_enabled", settings.Discovery.Enabled)
	_ = d.Set("discovery_tag_keys", settings.Discovery.TagKeys)
	_ = d.Set("sync_groups_enabled", settings.SyncGroups.Enabled)
	_ = d.Set("sync_group_rules", flattenSyncGroupRules(settings.SyncGroups.Rules))

	// Hierarchy levels: read each level entity to get its current name.
	levels := make([]map[string]interface{}, 0, len(settings.HierarchyLevelOrder))
	for _, levelID := range settings.HierarchyLevelOrder {
		levelIface, err := client.Scorecards.GetEntityWithContext(ctx, levelID)
		if err != nil || levelIface == nil || *levelIface == nil {
			continue
		}
		if level, ok := (*levelIface).(*scorecards.EntityManagementTeamsHierarchyLevelEntity); ok {
			levels = append(levels, map[string]interface{}{
				"id":   level.ID,
				"name": level.Name,
			})
		}
	}
	_ = d.Set("hierarchy_levels", levels)

	return nil
}

func resourceNewRelicTeamsOrgSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Updating NGEP teams organisation settings %s", d.Id())

	upd := scorecards.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{}
	updHasFields := false

	if d.HasChange("discovery_enabled") || d.HasChange("discovery_tag_keys") {
		tagKeys := make([]string, 0)
		for _, v := range d.Get("discovery_tag_keys").([]interface{}) {
			tagKeys = append(tagKeys, v.(string))
		}
		upd.Discovery = &scorecards.EntityManagementDiscoverySettingsUpdateInput{
			Enabled: d.Get("discovery_enabled").(bool),
			TagKeys: tagKeys,
		}
		updHasFields = true
	}

	if d.HasChange("hierarchy_levels") {
		rawLevels := d.Get("hierarchy_levels").([]interface{})
		levelIDs := make([]string, 0, len(rawLevels))
		for _, r := range rawLevels {
			m := r.(map[string]interface{})
			levelIDs = append(levelIDs, m["id"].(string))
		}
		upd.HierarchyLevelOrder = levelIDs
		updHasFields = true
	}

	if d.HasChange("sync_groups_enabled") || d.HasChange("sync_group_rules") {
		upd.SyncGroups = expandSyncGroupsUpdate(
			d.Get("sync_groups_enabled").(bool),
			d.Get("sync_group_rules").([]interface{}),
		)
		updHasFields = true
	}

	if updHasFields {
		if _, err := client.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(d.Id(), upd); err != nil {
			return diag.FromErr(err)
		}
	}

	// Rename hierarchy levels where name changed.
	if d.HasChange("hierarchy_levels") {
		old, newVal := d.GetChange("hierarchy_levels")
		oldMap := make(map[string]string)
		for _, r := range old.([]interface{}) {
			m := r.(map[string]interface{})
			oldMap[m["id"].(string)] = m["name"].(string)
		}
		for _, r := range newVal.([]interface{}) {
			m := r.(map[string]interface{})
			id, name := m["id"].(string), m["name"].(string)
			if oldName, exists := oldMap[id]; !exists || oldName != name {
				renameInput := scorecards.EntityManagementTeamsHierarchyLevelEntityUpdateInput{
					Name: name,
				}
				if _, err := client.Scorecards.EntityManagementUpdateTeamsHierarchyLevel(id, renameInput); err != nil {
					return diag.Errorf("renaming hierarchy level %s: %v", id, err)
				}
			}
		}
	}

	return resourceNewRelicTeamsOrgSettingsRead(ctx, d, meta)
}

// Delete only removes from state — the singleton entity cannot be destroyed.
func resourceNewRelicTeamsOrgSettingsDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[INFO] Removing newrelic_teams_organization_settings from state (singleton entity persists)")
	d.SetId("")
	return nil
}
