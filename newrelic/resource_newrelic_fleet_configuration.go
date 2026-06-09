package newrelic

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func resourceNewRelicFleetConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetConfigurationCreate,
		ReadContext:   resourceNewRelicFleetConfigurationRead,
		UpdateContext: resourceNewRelicFleetConfigurationUpdate,
		DeleteContext: resourceNewRelicFleetConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNewRelicFleetConfigurationImportState,
		},
		CustomizeDiff: resourceNewRelicFleetConfigurationCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the configuration. Changing this forces resource recreation because the API does not support renaming.",
			},
			"agent_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NRInfra",
					"NRDOT",
					"FluentBit",
					"NRPrometheusAgent",
				}, false),
				Description: "The type of agent this configuration is for. Allowed values: NRInfra, NRDOT, FluentBit, NRPrometheusAgent.",
			},
			"managed_entity_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HOST",
					"KUBERNETESCLUSTER",
				}, false),
				Description: "The type of entities this configuration manages. Allowed values: HOST, KUBERNETESCLUSTER.",
			},
			"operating_system": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LINUX",
					"WINDOWS",
				}, false),
				Description: "The operating system this configuration targets. Required for HOST configurations. Allowed values: LINUX, WINDOWS. Must not be set for KUBERNETESCLUSTER configurations.",
			},
			"configuration_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The configuration content (YAML or JSON). Use file() to load from a file. Each change to this field creates a new immutable version on the API.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. Auto-fetched from the account if not provided.",
			},
			// Computed
			"configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration entity GUID.",
			},
			"latest_version_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The highest version number across all versions.",
			},
			"latest_version_entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entity GUID of the highest-numbered version.",
			},
			"total_versions": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of versions in this configuration.",
			},
			"version_entity_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Entity GUIDs of all versions, sorted oldest first. Use with the newrelic_fleet_configuration data source to access a specific historical version.",
			},
		},
	}
}

func resourceNewRelicFleetConfigurationCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	managedEntityType := d.Get("managed_entity_type").(string)
	_, hasOS := d.GetOk("operating_system")

	if managedEntityType == "HOST" && !hasOS {
		return fmt.Errorf("operating_system is required when managed_entity_type is HOST")
	}
	if managedEntityType == "KUBERNETESCLUSTER" && hasOS {
		return fmt.Errorf("operating_system must not be set when managed_entity_type is KUBERNETESCLUSTER")
	}

	if d.HasChange("configuration_content") {
		for _, field := range []string{"total_versions", "latest_version_number", "latest_version_entity_id", "version_entity_ids"} {
			if err := d.SetNewComputed(field); err != nil {
				return fmt.Errorf("failed to mark %s as computed: %w", field, err)
			}
		}
	}

	return nil
}

// resourceNewRelicFleetConfigurationImportState handles import via a composite ID:
//
//	<configuration_guid>:<managed_entity_type>
//
// managed_entity_type must be included because the GetEntity GraphQL fragment
// for AgentConfigurationEntity does not return managedEntityType.
func resourceNewRelicFleetConfigurationImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf(
			"invalid import ID %q: expected \"<configuration_guid>:<managed_entity_type>\" "+
				"(e.g. \"NjQy...abc:HOST\" or \"NjQy...abc:KUBERNETESCLUSTER\")",
			d.Id(),
		)
	}
	guid, managedEntityType := parts[0], parts[1]

	providerConfig := meta.(*ProviderConfig)

	entityInterface, err := providerConfig.NewClient.FleetControl.GetEntityWithContext(ctx, guid)
	missing := entityInterface == nil || *entityInterface == nil
	// Treat as not-found only if the error is absent or genuinely a not-found
	// signal; transient errors (401, network) also return nil entity but must
	// surface to the user, not be misread as "missing".
	if missing && (err == nil || isFleetNotFoundError(err)) {
		return nil, fmt.Errorf("fleet configuration entity %s not found", guid)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fleet configuration entity %s: %w", guid, err)
	}

	entity, ok := (*entityInterface).(*fleetcontrol.EntityManagementAgentConfigurationEntity)
	if !ok {
		return nil, fmt.Errorf("entity %s is not a fleet configuration", guid)
	}

	d.SetId(guid)
	_ = d.Set("name", entity.Name)
	_ = d.Set("agent_type", entity.AgentType)
	_ = d.Set("managed_entity_type", managedEntityType)
	if entity.OperatingSystem.Type != "" {
		_ = d.Set("operating_system", string(entity.OperatingSystem.Type))
	}
	if entity.Scope.ID != "" {
		_ = d.Set("organization_id", entity.Scope.ID)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceNewRelicFleetConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	organizationID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	content := d.Get("configuration_content").(string)
	entityMeta := fleetConfigBuildEntityMeta(
		d.Get("name").(string),
		d.Get("agent_type").(string),
		d.Get("managed_entity_type").(string),
		d.Get("operating_system").(string),
	)

	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
		[]byte(content),
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": entityMeta,
			},
		},
		organizationID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating fleet configuration: %w", err))
	}

	d.SetId(result.ConfigurationEntityGUID)
	_ = d.Set("configuration_id", result.ConfigurationEntityGUID)
	_ = d.Set("organization_id", organizationID)

	log.Printf("[DEBUG] Created fleet configuration: %s", result.ConfigurationEntityGUID)

	return resourceNewRelicFleetConfigurationRead(ctx, d, meta)
}

func resourceNewRelicFleetConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfig := meta.(*ProviderConfig)

	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Capture prior state to detect out-of-band drift after we sync from API.
	priorVersionEntityIDs := expandStringList(d.Get("version_entity_ids").([]interface{}))
	priorLatestVersionEntityID := d.Get("latest_version_entity_id").(string)

	_ = d.Set("configuration_id", d.Id())
	_ = d.Set("organization_id", organizationID)

	versionsResp, err := providerConfig.NewClient.FleetControl.FleetControlGetConfigurationVersions(
		d.Id(), organizationID,
	)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if versionsResp == nil || len(versionsResp.Versions) == 0 {
		d.SetId("")
		return nil
	}

	// Sort versions ascending so version_entity_ids is stable oldest-first.
	apiVersions := make([]fleetcontrol.ConfigurationVersion, len(versionsResp.Versions))
	copy(apiVersions, versionsResp.Versions)
	sort.Slice(apiVersions, func(i, j int) bool {
		ni, _ := strconv.Atoi(apiVersions[i].Version)
		nj, _ := strconv.Atoi(apiVersions[j].Version)
		return ni < nj
	})

	var latestVersionNum int
	var latestVersionEntityID string
	versionEntityIDs := make([]string, 0, len(apiVersions))
	for _, v := range apiVersions {
		versionEntityIDs = append(versionEntityIDs, v.EntityGUID)
		num, parseErr := strconv.Atoi(v.Version)
		if parseErr != nil {
			return diag.FromErr(fmt.Errorf("failed to parse version number %q: %w", v.Version, parseErr))
		}
		if num > latestVersionNum {
			latestVersionNum = num
			latestVersionEntityID = v.EntityGUID
		}
	}

	// Detect and warn about out-of-band version deletions. Only fire when prior
	// state has values (skip the first read after Create/Import).
	if len(priorVersionEntityIDs) > 0 {
		currentSet := make(map[string]bool, len(versionEntityIDs))
		for _, id := range versionEntityIDs {
			currentSet[id] = true
		}
		var deleted []string
		for _, id := range priorVersionEntityIDs {
			if !currentSet[id] {
				deleted = append(deleted, id)
			}
		}
		if len(deleted) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Fleet configuration drift: %d version(s) deleted out-of-band", len(deleted)),
				Detail: fmt.Sprintf(
					"The following version entity GUIDs were tracked in Terraform state but no longer exist in the New Relic API: %v. "+
						"State has been synced to reflect the API. If any other resource or output references these GUIDs, those values will change accordingly.",
					deleted,
				),
			})

			// Stronger warning when the latest version itself was the one that disappeared,
			// because configuration_content gets overwritten and the next apply will create a new version.
			if priorLatestVersionEntityID != "" && !currentSet[priorLatestVersionEntityID] {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Fleet configuration latest version was deleted out-of-band",
					Detail: fmt.Sprintf(
						"The previously-tracked latest version (%q) no longer exists. "+
							"configuration_content has been refreshed from the new latest version (%q, version %d). "+
							"If your declared configuration_content differs from the new latest, the next apply will create a new version restoring your declared content.",
						priorLatestVersionEntityID, latestVersionEntityID, latestVersionNum,
					),
				})
			}
		}
	}

	// Fetch content of the latest version to keep configuration_content in sync.
	// If the version GUID returned by GetConfigurationVersions is already gone
	// at this endpoint (eventual-consistency window after a deletion), surface
	// a Warning and skip the content sync rather than hard-failing the plan —
	// the next refresh will reconcile naturally.
	content, fetchErr := providerConfig.NewClient.FleetControl.FleetControlGetConfiguration(
		latestVersionEntityID, organizationID, fleetcontrol.GetConfigurationModeTypes.ConfigVersionEntity, 0,
	)
	if fetchErr != nil {
		if isFleetNotFoundError(fetchErr) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Fleet configuration latest version content fetch returned not-found",
				Detail: fmt.Sprintf(
					"The version list reported %q as the latest version but the content endpoint returned not-found. "+
						"This is typically a brief eventual-consistency window after an out-of-band version deletion. "+
						"State sync skipped for configuration_content; the next refresh will reconcile.",
					latestVersionEntityID,
				),
			})
		} else {
			return diag.FromErr(fmt.Errorf("error fetching content for latest version %s: %w", latestVersionEntityID, fetchErr))
		}
	} else if content != nil {
		if err := d.Set("configuration_content", string(*content)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("total_versions", len(apiVersions)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("latest_version_number", latestVersionNum); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("latest_version_entity_id", latestVersionEntityID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version_entity_ids", versionEntityIDs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceNewRelicFleetConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// name is ForceNew, so a name change triggers a destroy+create flow rather
	// than reaching this function. configuration_content is the only mutable field.
	if !d.HasChange("configuration_content") {
		return nil
	}

	providerConfig := meta.(*ProviderConfig)
	organizationID := d.Get("organization_id").(string)
	content := d.Get("configuration_content").(string)

	_, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
		[]byte(content),
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": fmt.Sprintf(`{"agentConfiguration": "%s"}`, d.Id()),
			},
		},
		organizationID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating new version for fleet configuration %s: %w", d.Id(), err))
	}

	log.Printf("[DEBUG] Created new version for fleet configuration: %s", d.Id())

	return resourceNewRelicFleetConfigurationRead(ctx, d, meta)
}

func resourceNewRelicFleetConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Deleting New Relic Fleet Configuration %s", d.Id())

	// Fetch current versions from the API — the API rejects config deletion if any version still exists.
	versionsResp, err := providerConfig.NewClient.FleetControl.FleetControlGetConfigurationVersions(
		d.Id(), organizationID,
	)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error fetching versions for fleet configuration %s: %w", d.Id(), err))
	}

	if versionsResp != nil {
		for _, v := range versionsResp.Versions {
			log.Printf("[INFO] Deleting version %s", v.EntityGUID)
			if deleteErr := providerConfig.NewClient.FleetControl.FleetControlDeleteConfigurationVersion(
				v.EntityGUID, organizationID,
			); deleteErr != nil {
				return diag.FromErr(fmt.Errorf("failed to delete version %s: %w", v.EntityGUID, deleteErr))
			}
		}
	}

	_, err = providerConfig.NewClient.FleetControl.FleetControlDeleteConfiguration(d.Id(), organizationID)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			log.Printf("[INFO] Fleet configuration %s already gone (auto-removed by API)", d.Id())
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

// isFleetNotFoundError detects "not found" responses across the various error
// types the fleet control SDK returns. The blob-service REST endpoint wraps
// 404s as fmt.Errorf("resource not found") instead of *nrErrors.NotFound; the
// NerdGraph entity query returns *http.GraphQLErrorResponse with errorClass
// NOT_FOUND. We need a single check that catches both.
func isFleetNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(*nrErrors.NotFound); ok {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") || strings.Contains(msg, "resource not found")
}

// fleetConfigBuildEntityMeta builds the JSON string for the Newrelic-Entity header.
func fleetConfigBuildEntityMeta(name, agentType, managedEntityType, operatingSystem string) string {
	if operatingSystem != "" {
		return fmt.Sprintf(
			`{"name": "%s", "agentType": "%s", "managedEntityType": "%s", "operatingSystem": {"type": "%s"}}`,
			name, agentType, managedEntityType, operatingSystem,
		)
	}
	return fmt.Sprintf(
		`{"name": "%s", "agentType": "%s", "managedEntityType": "%s"}`,
		name, agentType, managedEntityType,
	)
}
