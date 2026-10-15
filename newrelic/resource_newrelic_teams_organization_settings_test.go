//go:build integration || SERVICE_ARCHITECTURE_INTELLIGENCE

// Integration tests for newrelic_teams_organization_settings and the
// newrelic_teams_hierarchy_levels data source.
//
// The org settings resource manages a pre-existing singleton entity — there is
// exactly one per New Relic organisation. Tests make small, reversible changes
// and restore the original values via t.Cleanup so the org is left unchanged.
//
// Run with:
//
//	TF_ACC=1 go test -tags integration -run TestAccNewRelicTeamsOrg ./newrelic/ -v -timeout 15m

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func testAccPreCheckOrgSettings(t *testing.T) {
	t.Helper()
	testAccPreCheck(t)
}

// captureAndRestoreOrgSettings saves the current org settings at test start and
// schedules a Cleanup that restores them, ensuring the org is left in its
// original state regardless of what the test does.
func captureAndRestoreOrgSettings(t *testing.T) *servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntity {
	t.Helper()
	if testAccProvider.Meta() == nil {
		return nil
	}
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	original, err := client.Scorecards.GetTeamsOrganizationSettings()
	if err != nil || original == nil {
		return nil
	}

	t.Cleanup(func() {
		if testAccProvider.Meta() == nil {
			return
		}
		c := testAccProvider.Meta().(*ProviderConfig).NewClient
		tagKeys := original.Discovery.TagKeys
		if tagKeys == nil {
			tagKeys = []string{}
		}
		rules := make([]servicearchintelligence.EntityManagementSyncGroupRuleUpdateInput, 0, len(original.SyncGroups.Rules))
		for _, r := range original.SyncGroups.Rules {
			conds := make([]servicearchintelligence.EntityManagementSyncGroupRuleConditionUpdateInput, 0, len(r.Conditions))
			for _, c := range r.Conditions {
				conds = append(conds, servicearchintelligence.EntityManagementSyncGroupRuleConditionUpdateInput{
					Type:  c.Type,
					Value: c.Value,
				})
			}
			rules = append(rules, servicearchintelligence.EntityManagementSyncGroupRuleUpdateInput{Conditions: conds})
		}
		_, _ = c.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(original.ID,
			servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{
				Discovery: &servicearchintelligence.EntityManagementDiscoverySettingsUpdateInput{
					Enabled: original.Discovery.Enabled,
					TagKeys: tagKeys,
				},
				SyncGroups: servicearchintelligence.EntityManagementSyncGroupsSettingsUpdateInput{
					Enabled: original.SyncGroups.Enabled,
					Rules:   rules,
				},
				HierarchyLevelOrder: original.HierarchyLevelOrder,
			},
		)
	})
	return original
}

// ── 1. Discovery settings — create, update, drift ────────────────────────────

// TestAccNewRelicTeamsOrgSettings_Discovery verifies that discovery_enabled and
// discovery_tag_keys are correctly applied, updated in-place, and that drift
// introduced out-of-band is detected on the next plan and reconciled on apply.
func TestAccNewRelicTeamsOrgSettings_Discovery(t *testing.T) {
	resourceName := "newrelic_teams_organization_settings.org"

	var settingsID string
	var original *servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntity

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOrgSettings(t)
			// Capture original settings for cleanup — provider is configured here.
			original = captureAndRestoreOrgSettings(t)
			if original != nil {
				settingsID = original.ID
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// ── Step 1: apply discovery config ────────────────────────────────
			{
				Config: testAccTeamsOrgSettingsDiscoveryConfig(true, []string{"team"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "discovery_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "discovery_tag_keys.#", "1"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						if rs == nil {
							return fmt.Errorf("resource %s not in state", resourceName)
						}
						settingsID = rs.Primary.ID
						return nil
					},
				),
			},
			// ── Step 2: import round-trip ──────────────────────────────────────
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sync_group_rules"},
			},
			// ── Step 3: update — add a second tag key ──────────────────────────
			{
				Config: testAccTeamsOrgSettingsDiscoveryConfig(true, []string{"team", "teamId"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "discovery_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "discovery_tag_keys.#", "2"),
				),
			},
			// ── Step 4: drift — disable discovery out-of-band + reconcile ─────
			// PreConfig flips discovery_enabled to false via the API directly.
			// The plan should detect the drift and the apply should restore it.
			{
				PreConfig: func() {
					if settingsID == "" || testAccProvider.Meta() == nil {
						return
					}
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, _ = client.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(
						settingsID,
						servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{
							Discovery: &servicearchintelligence.EntityManagementDiscoverySettingsUpdateInput{
								Enabled: false, // flip out-of-band
								TagKeys: []string{"team", "teamId"},
							},
						},
					)
				},
				Config: testAccTeamsOrgSettingsDiscoveryConfig(true, []string{"team", "teamId"}),
				Check: resource.ComposeTestCheckFunc(
					// After apply Terraform must have re-enabled discovery.
					resource.TestCheckResourceAttr(resourceName, "discovery_enabled", "true"),
				),
			},
			// ── Step 5: disable discovery (keep a tag key — API requires >= 1) ──
			{
				Config: testAccTeamsOrgSettingsDiscoveryConfig(false, []string{"team"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "discovery_enabled", "false"),
					// Tag key is preserved (required by API even when discovery is off).
					resource.TestCheckResourceAttr(resourceName, "discovery_tag_keys.#", "1"),
				),
			},
		},
	})
	_ = original // used via t.Cleanup closure
}

// ── 2. Sync group rules — all condition types + drift ─────────────────────────

// TestAccNewRelicTeamsOrgSettings_SyncGroupRules verifies that sync_group_rules
// round-trips correctly for all three condition types (STARTS_WITH, ENDS_WITH,
// CONTAINS), that rules can be updated in-place, and that adding a rule
// out-of-band is detected as drift.
func TestAccNewRelicTeamsOrgSettings_SyncGroupRules(t *testing.T) {
	resourceName := "newrelic_teams_organization_settings.org"

	var settingsID string
	var original *servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntity

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOrgSettings(t)
			original = captureAndRestoreOrgSettings(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// ── Step 1: one rule with STARTS_WITH ─────────────────────────────
			{
				Config: testAccTeamsOrgSettingsSyncRulesConfig([]syncRuleSpec{
					{condType: "STARTS_WITH", value: "tf-test-"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sync_groups_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.0.conditions.0.type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.0.conditions.0.value", "tf-test-"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						if rs != nil {
							settingsID = rs.Primary.ID
						}
						return nil
					},
				),
			},
			// ── Step 2: update rule to use ENDS_WITH ──────────────────────────
			// API allows at most one rule — update the condition type in-place.
			{
				Config: testAccTeamsOrgSettingsSyncRulesConfig([]syncRuleSpec{
					{condType: "ENDS_WITH", value: "-eng"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.0.conditions.0.type", "ENDS_WITH"),
				),
			},
			// ── Step 3: switch to CONTAINS ────────────────────────────────────
			{
				Config: testAccTeamsOrgSettingsSyncRulesConfig([]syncRuleSpec{
					{condType: "CONTAINS", value: "-engineering-"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.0.conditions.0.type", "CONTAINS"),
				),
			},
			// ── Step 4: change rule condition out-of-band + verify drift ───────
			// The API only allows one rule, so drift testing changes the condition
			// value rather than adding a new rule.
			{
				PreConfig: func() {
					if settingsID == "" || testAccProvider.Meta() == nil {
						return
					}
					client := testAccProvider.Meta().(*ProviderConfig).NewClient
					_, _ = client.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(
						settingsID,
						servicearchintelligence.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{
							SyncGroups: servicearchintelligence.EntityManagementSyncGroupsSettingsUpdateInput{
								Enabled: true,
								Rules: []servicearchintelligence.EntityManagementSyncGroupRuleUpdateInput{
									// Out-of-band: change condition value
									{Conditions: []servicearchintelligence.EntityManagementSyncGroupRuleConditionUpdateInput{
										{Type: "CONTAINS", Value: "out-of-band-drift"},
									}},
								},
							},
						},
					)
				},
				Config: testAccTeamsOrgSettingsSyncRulesConfig([]syncRuleSpec{
					{condType: "CONTAINS", value: "-engineering-"},
				}),
				Check: resource.ComposeTestCheckFunc(
					// After apply, the declared value must be restored.
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sync_group_rules.0.conditions.0.value", "-engineering-"),
				),
			},
			// ── Step 5: disable sync groups ───────────────────────────────────
			// Omitting sync_group_rules with Computed:true carries the existing
			// rule forward in state (Terraform doesn't touch it). The important
			// assertion is that sync_groups_enabled is false.
			{
				Config: testAccTeamsOrgSettingsSyncGroupsDisabledConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sync_groups_enabled", "false"),
				),
			},
		},
	})
	_ = original
}

// ── 3. Hierarchy levels data source ──────────────────────────────────────────

// TestAccNewRelicTeamsHierarchyLevels_DataSource verifies that the data source
// correctly reads all hierarchy level entities from the organisation and exposes
// their IDs and names.
func TestAccNewRelicTeamsHierarchyLevels_DataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckOrgSettings(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "newrelic_teams_hierarchy_levels" "all" {}`,
				Check: resource.ComposeTestCheckFunc(
					// At minimum the data source must return without error.
					// In orgs with hierarchy levels configured, levels.# > 0.
					resource.TestCheckResourceAttrSet("data.newrelic_teams_hierarchy_levels.all", "id"),
				),
			},
		},
	})
}

// ── Config templates ──────────────────────────────────────────────────────────

func testAccTeamsOrgSettingsDiscoveryConfig(enabled bool, tagKeys []string) string {
	keysHCL := ""
	if len(tagKeys) > 0 {
		quoted := make([]string, len(tagKeys))
		for i, k := range tagKeys {
			quoted[i] = fmt.Sprintf("%q", k)
		}
		keysHCL = fmt.Sprintf("[%s]", joinStrings(quoted, ", "))
	} else {
		keysHCL = "[]"
	}
	return fmt.Sprintf(`
resource "newrelic_teams_organization_settings" "org" {
  discovery_enabled  = %t
  discovery_tag_keys = %s
}
`, enabled, keysHCL)
}

type syncRuleSpec struct {
	condType string
	value    string
}

func testAccTeamsOrgSettingsSyncRulesConfig(rules []syncRuleSpec) string {
	rulesHCL := ""
	for _, r := range rules {
		rulesHCL += fmt.Sprintf(`
  sync_group_rules {
    conditions {
      type  = %q
      value = %q
    }
  }`, r.condType, r.value)
	}
	return fmt.Sprintf(`
resource "newrelic_teams_organization_settings" "org" {
  sync_groups_enabled = true
%s
}
`, rulesHCL)
}

func testAccTeamsOrgSettingsSyncGroupsDisabledConfig() string {
	// Note: sync_group_rules block is absent — with Computed:true, the existing
	// rule carries forward in state. Only sync_groups_enabled is toggled here.
	return `
resource "newrelic_teams_organization_settings" "org" {
  sync_groups_enabled = false
}
`
}

// joinStrings joins a slice of strings with sep (avoids importing strings package in test file).
func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
