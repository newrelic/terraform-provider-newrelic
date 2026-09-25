package newrelic

// resource_newrelic_teams_organization_settings manages the singleton
// TEAMS_ORGANIZATION_SETTINGS entity for the authenticated organisation.
//
// This entity is auto-created by NGEP when Teams is first used; there is
// exactly one per organisation. Terraform cannot create or delete it — use
// terraform import to bring the existing entity under management.
//
//	terraform import newrelic_teams_organization_settings.this <entity_id>
//
// The entity_id can be found by running:
//
//	newrelic_teams_organization_settings data source, or
//	entitySearch(query: "type = 'TEAMS_ORGANIZATION_SETTINGS'") via NerdGraph
//
// Key capabilities managed here:
//   - discovery.enabled  — toggles tag-based automatic entity ownership assignment
//   - discovery.tag_keys — the tag keys used for auto-assignment (e.g. ["team"])

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

func resourceNewRelicTeamsOrganizationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicTeamsOrganizationSettingsCreate,
		ReadContext:   resourceNewRelicTeamsOrganizationSettingsRead,
		UpdateContext: resourceNewRelicTeamsOrganizationSettingsUpdate,
		DeleteContext: resourceNewRelicTeamsOrganizationSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// discovery controls the tag-based automatic entity ownership assignment.
			"discovery_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Whether tag-based automatic entity ownership assignment is enabled " +
					"for this organisation. When true, entities with tags matching a team's name " +
					"or aliases are automatically added to that team's ownership collection.",
			},
			"discovery_tag_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Description: "The tag keys used to match entities to teams when discovery is enabled. " +
					"For example, ['team'] means any entity with tag 'team:<team-name>' is automatically " +
					"owned by the team named '<team-name>'.",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// Create is not supported — this entity is a singleton created by NGEP.
func resourceNewRelicTeamsOrganizationSettingsCreate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{{
		Severity: diag.Error,
		Summary:  "newrelic_teams_organization_settings does not support create",
		Detail: "The Teams organisation settings entity is created automatically by New Relic. " +
			"Use 'terraform import newrelic_teams_organization_settings.<name> <entity_id>' to " +
			"bring the existing settings under Terraform management. The entity ID can be found " +
			"via: entitySearch(query: \"type = 'TEAMS_ORGANIZATION_SETTINGS'\") in NerdGraph.",
	}}
}

func resourceNewRelicTeamsOrganizationSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	entityIface, err := client.Scorecards.GetEntityWithContext(ctx, d.Id())
	if err != nil {
		var notFound *nrErrors.NotFound
		// If entity not found (deleted outside Terraform), remove from state.
		if errors.As(err, &notFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if entityIface == nil {
		d.SetId("")
		return nil
	}

	settings, ok := (*entityIface).(*scorecards.EntityManagementTeamsOrganizationSettingsEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamsOrganizationSettingsEntity", d.Id())
	}

	_ = d.Set("discovery_enabled", settings.Discovery.Enabled)
	_ = d.Set("discovery_tag_keys", settings.Discovery.TagKeys)

	return nil
}

func resourceNewRelicTeamsOrganizationSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Updating NGEP teams organisation settings %s", d.Id())

	tagKeys := make([]string, 0)
	for _, k := range d.Get("discovery_tag_keys").([]interface{}) {
		tagKeys = append(tagKeys, k.(string))
	}

	upd := scorecards.EntityManagementTeamsOrganizationSettingsEntityUpdateInput{
		Discovery: scorecards.EntityManagementDiscoverySettingsUpdateInput{
			Enabled: d.Get("discovery_enabled").(bool),
			TagKeys: tagKeys,
		},
	}

	if _, err := client.Scorecards.EntityManagementUpdateTeamsOrganizationSettings(d.Id(), upd); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicTeamsOrganizationSettingsRead(ctx, d, meta)
}

// Delete only removes from state — the singleton entity cannot be destroyed.
func resourceNewRelicTeamsOrganizationSettingsDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[INFO] Removing teams org settings %s from state (entity is a singleton managed by NGEP)", d.Id())
	d.SetId("")
	return nil
}
