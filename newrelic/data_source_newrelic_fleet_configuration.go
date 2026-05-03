package newrelic

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func dataSourceNewRelicFleetConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicFleetConfigurationRead,
		Schema: map[string]*schema.Schema{
			// ── inputs (mutually exclusive) ────────────────────────────────
			"configuration_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"version_entity_id", "name"},
				Description:   "The GUID of the fleet configuration entity. Returns the content of the latest version. Populated automatically when looking up by version_entity_id.",
			},
			"version_entity_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"configuration_id", "name"},
				Description:   "The GUID of a specific configuration version entity. Returns the content of that exact version.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"configuration_id", "version_entity_id"},
				Description:   "The name of the fleet configuration. The first matching configuration is returned. Returns the content of its latest version.",
			},
			// ── optional context ──────────────────────────────────────────
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The organization ID. Resolved automatically from the provider when omitted.",
			},
			// ── outputs ───────────────────────────────────────────────────
			"configuration_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The raw configuration content (YAML/JSON) of the resolved version.",
			},
			"latest_version_entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity GUID of the latest version. Populated when looking up by configuration_id or name.",
			},
			"version_entity_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Entity GUIDs of all versions ordered by version number (oldest first). Populated when looking up by configuration_id or name.",
			},
		},
	}
}

func dataSourceNewRelicFleetConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*ProviderConfig)
	client := cfg.NewClient

	orgID, err := getOrganizationID(ctx, cfg, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var configID string
	var content *fleetcontrol.GetConfigurationResponse
	var diags diag.Diagnostics

	switch {
	case d.Get("configuration_id").(string) != "":
		configID, content, diags = dsFleetReadByConfigID(d, client, orgID)
	case d.Get("version_entity_id").(string) != "":
		configID, content, diags = dsFleetReadByVersionGUID(d, client, orgID)
	case d.Get("name").(string) != "":
		configID, orgID, content, diags = dsFleetReadByName(d, client, orgID)
	default:
		return diag.Errorf("one of configuration_id, version_entity_id, or name must be specified")
	}

	if diags.HasError() {
		return diags
	}

	if content == nil {
		return diag.Errorf("received nil content for fleet configuration %q", configID)
	}

	d.SetId(configID)
	if setErr := d.Set("organization_id", orgID); setErr != nil {
		return diag.FromErr(setErr)
	}
	if setErr := d.Set("configuration_content", string(*content)); setErr != nil {
		return diag.FromErr(setErr)
	}

	return nil
}

func dsFleetReadByConfigID(d *schema.ResourceData, client *nr.NewRelic, orgID string) (string, *fleetcontrol.GetConfigurationResponse, diag.Diagnostics) {
	configID := d.Get("configuration_id").(string)
	content, err := client.FleetControl.FleetControlGetConfiguration(
		configID, orgID, fleetcontrol.GetConfigurationModeTypes.ConfigEntity, 0,
	)
	if err != nil {
		return "", nil, diag.FromErr(fmt.Errorf("failed to fetch fleet configuration %q: %w", configID, err))
	}
	if diags := setVersionFields(d, client, configID, orgID); diags != nil {
		return "", nil, diags
	}
	return configID, content, nil
}

func dsFleetReadByVersionGUID(d *schema.ResourceData, client *nr.NewRelic, orgID string) (string, *fleetcontrol.GetConfigurationResponse, diag.Diagnostics) {
	versionGUID := d.Get("version_entity_id").(string)
	content, err := client.FleetControl.FleetControlGetConfiguration(
		versionGUID, orgID, fleetcontrol.GetConfigurationModeTypes.ConfigVersionEntity, 0,
	)
	if err != nil {
		return "", nil, diag.FromErr(fmt.Errorf("failed to fetch fleet configuration version %q: %w", versionGUID, err))
	}

	// Resolve the parent configuration GUID from the version entity
	entityIface, entityErr := client.FleetControl.GetEntity(versionGUID)
	if entityErr == nil && entityIface != nil {
		if versionEntity, ok := (*entityIface).(*fleetcontrol.EntityManagementAgentConfigurationVersionEntity); ok && versionEntity.AgentConfiguration != "" {
			if setErr := d.Set("configuration_id", versionEntity.AgentConfiguration); setErr != nil {
				return "", nil, diag.FromErr(setErr)
			}
		}
	}

	return versionGUID, content, nil
}

func dsFleetReadByName(d *schema.ResourceData, client *nr.NewRelic, orgID string) (string, string, *fleetcontrol.GetConfigurationResponse, diag.Diagnostics) {
	name := d.Get("name").(string)
	// The entity management entitySearch API only supports type-based filtering;
	// name filtering is not a valid predicate. Fetch all AGENT_CONFIGURATION
	// entities and filter by name client-side.
	result, searchErr := client.FleetControl.GetEntitySearch(
		"", "type = 'AGENT_CONFIGURATION'",
	)
	if searchErr != nil {
		return "", orgID, nil, diag.FromErr(fmt.Errorf("failed to search fleet configurations by name %q: %w", name, searchErr))
	}
	if result == nil || len(result.Entities) == 0 {
		return "", orgID, nil, diag.Errorf("no fleet configuration found with name %q", name)
	}

	var configID string
	for _, entity := range result.Entities {
		if cfgEntity, ok := entity.(*fleetcontrol.EntityManagementAgentConfigurationEntity); ok {
			if cfgEntity.Name != name {
				continue
			}
			configID = cfgEntity.ID
			if cfgEntity.Scope.Type == "ORGANIZATION" {
				orgID = cfgEntity.Scope.ID
			}
			break
		}
	}
	if configID == "" {
		return "", orgID, nil, diag.Errorf("no fleet configuration found with name %q", name)
	}

	if setErr := d.Set("configuration_id", configID); setErr != nil {
		return "", orgID, nil, diag.FromErr(setErr)
	}

	content, err := client.FleetControl.FleetControlGetConfiguration(
		configID, orgID, fleetcontrol.GetConfigurationModeTypes.ConfigEntity, 0,
	)
	if err != nil {
		return "", orgID, nil, diag.FromErr(fmt.Errorf("failed to fetch fleet configuration %q (found via name %q): %w", configID, name, err))
	}
	if diags := setVersionFields(d, client, configID, orgID); diags != nil {
		return "", orgID, nil, diags
	}

	return configID, orgID, content, nil
}

// setVersionFields fetches all versions for configID and populates
// latest_version_entity_id and version_entity_ids.
func setVersionFields(d *schema.ResourceData, client *nr.NewRelic, configID, orgID string) diag.Diagnostics {
	versionsResp, err := client.FleetControl.FleetControlGetConfigurationVersions(configID, orgID)
	if err != nil {
		// Non-fatal: log and skip rather than failing the data source read
		return nil
	}
	if versionsResp == nil || len(versionsResp.Versions) == 0 {
		return nil
	}

	// Sort by version number (Version field is a string representing an int)
	versions := versionsResp.Versions
	sort.Slice(versions, func(i, j int) bool {
		ni, _ := strconv.Atoi(versions[i].Version)
		nj, _ := strconv.Atoi(versions[j].Version)
		return ni < nj
	})

	// Collect GUIDs in sorted order
	guids := make([]string, len(versions))
	for i, v := range versions {
		guids[i] = v.EntityGUID
	}

	// Latest = highest version number = last after sort
	latestGUID := versions[len(versions)-1].EntityGUID

	if setErr := d.Set("latest_version_entity_id", latestGUID); setErr != nil {
		return diag.FromErr(setErr)
	}
	if setErr := d.Set("version_entity_ids", guids); setErr != nil {
		return diag.FromErr(setErr)
	}
	return nil
}
