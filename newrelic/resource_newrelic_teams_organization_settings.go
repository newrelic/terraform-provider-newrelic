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
//   - discovery_enabled    — toggles tag-based automatic entity ownership
//   - discovery_tag_keys   — the tag keys used for auto-assignment (e.g. ["team"])
//   - hierarchy_levels     — ordered list of hierarchy level entities (id + name)
//   - sync_groups_enabled  — toggles automatic team creation from IdP groups
//   - sync_group_rules     — rules controlling which IdP groups create teams

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
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
			// TypeSet (not TypeList) so that ordering differences between the
			// provider's declared list and the API's returned order do not produce
			// phantom diffs. Discovery tag keys are semantically unordered.
			// MinItems: 1 — the NGEP API rejects an empty tagKeys array in the
			// discovery update input. Always provide at least one key.
			"discovery_tag_keys": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Description: "Tag keys used for entity discovery (e.g. [\"team\"]). Entities with these tags " +
					"are auto-assigned to the team whose name or alias matches the tag value. " +
					"Must contain at least one key when declared.",
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// hierarchy_levels is TypeList because the order of entries matters:
			// the NGEP API requires HierarchyLevelOrder to be sorted root-to-leaf
			// matching the team parent_id chain depth. TypeSet would scramble the
			// order via hash iteration and cause the API to reject the payload.
			// Phantom ordering diffs are prevented by the Read function, which
			// only includes declared levels in state.
			"hierarchy_levels": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Description: "Ordered list of organisation hierarchy levels to manage (root first, " +
					"leaf last). Use the newrelic_teams_hierarchy_levels data source to obtain level GUIDs.",
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
							Description:  "Display name of this hierarchy level (e.g. 'Division', 'Squad'). Renaming here updates the entity in the API.",
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
			// Computed:true so that when sync_group_rules is absent from config,
			// Terraform carries the prior state forward (no phantom removal of rules).
			// To explicitly clear all rules, declare sync_group_rules = [] (empty block).
			// MaxItems: 1 — the NGEP API currently enforces a single rule per organisation.
			"sync_group_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Description: "Rules that control which IdP groups automatically create teams. Each rule " +
					"specifies one or more match conditions. Omitting this block preserves existing rules; " +
					"set to [] to remove all rules. The NGEP API currently allows at most one rule.",
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

// orgSettingsAttrInConfig returns true when the given attribute is explicitly
// declared in the user's Terraform config. Used by Create to skip fields the
// user did not intend to configure, so that pre-existing org settings are not
// accidentally overridden.
//
// Falls back to true when the raw config is unavailable (e.g. terraform import)
// to preserve the old behavior of sending all fields.
func orgSettingsAttrInConfig(d *schema.ResourceData, key string) bool {
	rc := d.GetRawConfig()
	if !rc.IsKnown() || rc.IsNull() {
		return true
	}
	attr := rc.GetAttr(key)
	return attr.IsKnown() && !attr.IsNull()
}

// resourceNewRelicTeamsOrgSettingsCreate locates the pre-existing singleton,
// applies only the configuration attributes that are explicitly declared in the
// user's config, and emits a warning so the operator knows existing settings
// have been overridden.
//
// Only declared fields are sent — this prevents accidentally resetting sync
// group rules, hierarchy levels, or discovery settings that the user has not
// included in their Terraform configuration.
func resourceNewRelicTeamsOrgSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

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

	upd := servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{}
	updHasFields := false

	// Discovery — send only when at least one discovery attribute is declared.
	if orgSettingsAttrInConfig(d, "discovery_enabled") || orgSettingsAttrInConfig(d, "discovery_tag_keys") {
		tagKeys := make([]string, 0)
		for _, v := range d.Get("discovery_tag_keys").(*schema.Set).List() {
			tagKeys = append(tagKeys, v.(string))
		}
		upd.Discovery = &servicearchintelligence.EntityManagementDiscoverySettingsUpdateInput{
			Enabled: d.Get("discovery_enabled").(bool),
			TagKeys: tagKeys,
		}
		updHasFields = true
	}

	// Hierarchy levels — send only when hierarchy_levels is declared.
	// Always include ALL existing levels in the order: declared levels first
	// (in config order), then any remaining existing levels appended.
	// The NGEP API rejects a partial HierarchyLevelOrder that omits any
	// existing hierarchy level entity.
	rawLevels := d.Get("hierarchy_levels").([]interface{})
	if orgSettingsAttrInConfig(d, "hierarchy_levels") && len(rawLevels) > 0 {
		declaredIDs := make([]string, 0, len(rawLevels))
		for _, r := range rawLevels {
			declaredIDs = append(declaredIDs, r.(map[string]interface{})["id"].(string))
		}
		upd.HierarchyLevelOrder = completeHierarchyLevelOrder(declaredIDs, existing.HierarchyLevelOrder)
		updHasFields = true
	}

	// Sync groups — send only when at least one sync attribute is declared.
	if orgSettingsAttrInConfig(d, "sync_groups_enabled") || orgSettingsAttrInConfig(d, "sync_group_rules") {
		upd.SyncGroups = expandSyncGroupsUpdate(
			d.Get("sync_groups_enabled").(bool),
			d.Get("sync_group_rules").([]interface{}),
		)
		updHasFields = true
	}

	if updHasFields {
		if _, err := client.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(d.Id(), upd); err != nil {
			return diag.Errorf("applying org settings configuration: %v", err)
		}
	}

	// Rename hierarchy levels where the declared name differs from the current name.
	if orgSettingsAttrInConfig(d, "hierarchy_levels") {
		for _, r := range rawLevels {
			m := r.(map[string]interface{})
			id, wantName := m["id"].(string), m["name"].(string)

			currentName := ""
			levelIface, levelErr := client.Scorecards.GetEntityWithContext(ctx, id)
			if levelErr == nil && levelIface != nil && *levelIface != nil {
				if level, ok := (*levelIface).(*servicearchintelligence.EntityManagementTeamsHierarchyLevelEntity); ok {
					currentName = level.Name
				}
			}

			if currentName == wantName {
				continue
			}
			if _, err := client.Scorecards.EntityManagementUpdateTeamsHierarchyLevel(
				id,
				servicearchintelligence.EntityManagementTeamsHierarchyLevelEntityUpdateInput{Name: wantName},
			); err != nil {
				return diag.Errorf("renaming hierarchy level %s: %v", id, err)
			}
		}
	}

	diags := resourceNewRelicTeamsOrgSettingsRead(ctx, d, meta)

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

	settings, ok := (*entityIface).(*servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamsOrganizationSettingsEntity", d.Id())
	}

	_ = d.Set("discovery_enabled", settings.Discovery.Enabled)
	_ = d.Set("discovery_tag_keys", settings.Discovery.TagKeys)
	_ = d.Set("sync_groups_enabled", settings.SyncGroups.Enabled)
	_ = d.Set("sync_group_rules", flattenSyncGroupRules(settings.SyncGroups.Rules))

	// Hierarchy levels: read each level entity to get its current name.
	// Determine which hierarchy level IDs to reflect in state.
	//
	// We use the IDs already in state (or config on first apply) as the
	// authoritative filter. This prevents undeclared levels — created by other
	// tools or earlier Terraform runs — from appearing in state and producing a
	// phantom "remove" diff on every plan.
	//
	// completeHierarchyLevelOrder in Create/Update always sends ALL org-level
	// IDs to the API, so undeclared levels are preserved without being tracked.
	//
	// Import path: prior state is empty, so we fall back to all API levels so
	// the operator can see everything and adopt levels into config as needed.
	stateRaw := d.Get("hierarchy_levels").([]interface{})
	stateIDs := make(map[string]bool, len(stateRaw))
	for _, item := range stateRaw {
		if m, ok := item.(map[string]interface{}); ok {
			if id, ok := m["id"].(string); ok && id != "" {
				stateIDs[id] = true
			}
		}
	}

	levels := make([]map[string]interface{}, 0, len(settings.HierarchyLevelOrder))
	for _, levelID := range settings.HierarchyLevelOrder {
		// When stateIDs is non-empty (steady-state plan/update), skip any level
		// not already tracked. When empty (import), include all levels.
		if len(stateIDs) > 0 && !stateIDs[levelID] {
			continue
		}
		levelIface, err := client.Scorecards.GetEntityWithContext(ctx, levelID)
		if err != nil || levelIface == nil || *levelIface == nil {
			continue
		}
		if level, ok := (*levelIface).(*servicearchintelligence.EntityManagementTeamsHierarchyLevelEntity); ok {
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

	// Fetch the current org settings upfront so Update has the full
	// HierarchyLevelOrder from the API. This is needed by
	// completeHierarchyLevelOrder to append any undeclared levels — we cannot
	// rely on state because Read only stores declared levels.
	currentSettings, settingsErr := client.Scorecards.GetTeamsOrganizationSettingsWithContext(ctx)
	var currentHierarchyOrder []string
	if settingsErr == nil && currentSettings != nil {
		currentHierarchyOrder = currentSettings.HierarchyLevelOrder
	}

	upd := servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{}
	updHasFields := false

	if d.HasChange("discovery_enabled") || d.HasChange("discovery_tag_keys") {
		tagKeys := make([]string, 0)
		for _, v := range d.Get("discovery_tag_keys").(*schema.Set).List() {
			tagKeys = append(tagKeys, v.(string))
		}
		upd.Discovery = &servicearchintelligence.EntityManagementDiscoverySettingsUpdateInput{
			Enabled: d.Get("discovery_enabled").(bool),
			TagKeys: tagKeys,
		}
		updHasFields = true
	}

	if d.HasChange("hierarchy_levels") {
		_, newVal := d.GetChange("hierarchy_levels")
		declaredIDs := make([]string, 0)
		for _, r := range newVal.([]interface{}) {
			declaredIDs = append(declaredIDs, r.(map[string]interface{})["id"].(string))
		}
		// Always send the complete order using the live API list as the
		// authoritative source of existing levels, not state (which only
		// contains declared levels since the Read fix).
		upd.HierarchyLevelOrder = completeHierarchyLevelOrder(declaredIDs, currentHierarchyOrder)
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
				if _, err := client.Scorecards.EntityManagementUpdateTeamsHierarchyLevel(id,
					servicearchintelligence.EntityManagementTeamsHierarchyLevelEntityUpdateInput{Name: name},
				); err != nil {
					return diag.Errorf("renaming hierarchy level %s: %v", id, err)
				}
			}
		}
	}

	return resourceNewRelicTeamsOrgSettingsRead(ctx, d, meta)
}

// completeHierarchyLevelOrder builds a HierarchyLevelOrder slice that satisfies
// the NGEP API requirement: every existing hierarchy level in the organisation
// must be included in the order. Callers provide the levels the user explicitly
// declared (in the order they want) and the current full order from the API.
//
// The returned list is: [declared IDs in config order] + [any existing IDs not
// in the declared set, preserving their current relative order].
//
// This prevents the "Hierarchy levels are managed automatically, manual changes
// are not allowed" API error that occurs when HierarchyLevelOrder omits one or
// more existing hierarchy level entities.
func completeHierarchyLevelOrder(declaredIDs []string, existingOrder []string) []string {
	declared := make(map[string]bool, len(declaredIDs))
	for _, id := range declaredIDs {
		declared[id] = true
	}
	result := make([]string, 0, len(existingOrder))
	result = append(result, declaredIDs...)
	for _, id := range existingOrder {
		if !declared[id] {
			result = append(result, id)
		}
	}
	return result
}

// Delete only removes from state — the singleton entity cannot be destroyed.
func resourceNewRelicTeamsOrgSettingsDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[INFO] Removing newrelic_teams_organization_settings from state (singleton entity persists)")
	d.SetId("")
	return nil
}
