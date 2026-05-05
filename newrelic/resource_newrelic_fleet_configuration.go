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
				Description: "The name of the configuration.",
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
			// TypeList — preserves insertion order so CustomizeDiff sees all blocks (including
			// duplicates) and positional removal is unambiguous (surviving blocks encode which
			// one was removed). This also makes rollback sequences like v1(A)→v2(B)→v3(A) safe,
			// because identical-content blocks are NOT collapsed before CustomizeDiff runs.
			//
			// FUTURE: TypeSet alternative preserved below if content-hash deduplication is preferred.
			"version": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Description: "Configuration versions. At least one required. Each version must have unique content. " +
					"Use file() to load content: configuration_content = file(\"${path.module}/config.yaml\").",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_content": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Configuration content for this version (YAML or JSON). " +
								"Content must be unique across version blocks. " +
								"Use file() to load from a file: file(\"${path.module}/config.yaml\").",
						},
						"version_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version number assigned by the API (1, 2, 3, ...).",
						},
						"version_entity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version entity GUID.",
						},
					},
				},
			},
			/* TypeSet alternative — uncomment to switch back if content-hash deduplication is preferred:
			"version": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Description: "Configuration versions. At least one required. Each version must have unique content. " +
					"Use file() to load content: configuration_content = file(\"${path.module}/config.yaml\").",
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["configuration_content"].(string))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version_entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			*/
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. Auto-fetched from the account if not provided.",
			},
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
				Description: "Total number of versions.",
			},
		},
	}
}

// resourceNewRelicFleetConfigurationCustomizeDiff validates version blocks and
// marks derived scalar fields as unknown when the version list changes.
//
// Three checks are performed:
//
//  1. In-place content modification guard — versions are immutable in the API;
//     there is no "update content" mutation. Editing configuration_content of an
//     existing block would silently delete that version and create a new one with a
//     higher version number. We detect this and surface a clear error instead.
//
//     Detection heuristic: for each position i present in both old and new lists,
//     if old[i] has a persisted version_entity_id (it was previously applied) AND
//     new[i] has different content AND that new content did not exist anywhere in
//     the old list (ruling out blocks that merely shifted position due to a removal),
//     we treat it as an attempted in-place edit.
//
//  2. Duplicate content guard — with TypeList the SDK does NOT deduplicate blocks
//     before this function runs, so this is the primary uniqueness check.
//
//  3. SetNewComputed — marks derived scalar fields as unknown whenever the version
//     list changes so the plan accurately shows they will be re-evaluated.
func resourceNewRelicFleetConfigurationCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	managedEntityType := d.Get("managed_entity_type").(string)
	_, hasOS := d.GetOk("operating_system")

	if managedEntityType == "KUBERNETESCLUSTER" && hasOS {
		return fmt.Errorf("operating_system must not be set when managed_entity_type is KUBERNETESCLUSTER")
	}

	versionsRaw, ok := d.GetOk("version")
	if !ok {
		return nil
	}

	// --- Check 1: in-place content modification ---
	//
	// Only meaningful when the list length is unchanged. If the count differs,
	// the change is a pure addition or removal — positional comparison would
	// produce false positives (e.g. a version deleted from the API is removed
	// from state by Read, making the counts differ on the next plan).
	if d.HasChange("version") {
		oldRaw, newRaw := d.GetChange("version")
		oldList := oldRaw.([]interface{})
		newList := newRaw.([]interface{})

		if len(oldList) > 0 && len(oldList) == len(newList) {
			// Build a set of all content values present in the old state so we can
			// distinguish a genuinely new content value from a block that shifted
			// position due to a prior removal that was already applied.
			oldContentSet := make(map[string]bool, len(oldList))
			for _, v := range oldList {
				content := v.(map[string]interface{})["configuration_content"].(string)
				oldContentSet[content] = true
			}

			var violations []string
			for i := 0; i < len(oldList); i++ {
				oldMap := oldList[i].(map[string]interface{})
				newMap := newList[i].(map[string]interface{})

				entityID, _ := oldMap["version_entity_id"].(string)
				oldContent, _ := oldMap["configuration_content"].(string)
				newContent, _ := newMap["configuration_content"].(string)

				// A persisted block whose content changed to something not present
				// anywhere in the old list is an in-place edit attempt.
				if entityID != "" && oldContent != newContent && !oldContentSet[newContent] {
					violations = append(violations, fmt.Sprintf(
						"  - index %d: to replace this content, add a new version block with "+
							"the updated content (apply), then remove the old one (apply again)",
						i,
					))
				}
			}
			if len(violations) > 0 {
				return fmt.Errorf(
					"configuration_content cannot be modified in place — versions are immutable:\n%s",
					strings.Join(violations, "\n"),
				)
			}
		}
	}

	// --- Check 2: duplicate content ---
	// seen maps configuration_content → version_number of first occurrence (0 if not yet assigned).
	seen := make(map[string]int)
	for _, v := range versionsRaw.([]interface{}) {
		vMap := v.(map[string]interface{})
		content := vMap["configuration_content"].(string)
		vNum, _ := vMap["version_number"].(int)
		if existingNum, exists := seen[content]; exists {
			if existingNum > 0 {
				return fmt.Errorf(
					"duplicate configuration_content detected — matches version %d; "+
						"add a distinguishing comment if you intend to roll back to that content",
					existingNum,
				)
			}
			return fmt.Errorf(
				"duplicate configuration_content detected across version blocks — " +
					"each version must have unique content; add a distinguishing comment " +
					"if two versions are otherwise identical",
			)
		}
		seen[content] = vNum
	}

	// --- Check 3: mark derived scalar fields as unknown ---
	// When the version list changes, mark derived scalar fields as unknown so
	// Terraform re-evaluates them (and dependent outputs) after apply.
	if d.HasChange("version") {
		for _, field := range []string{"total_versions", "latest_version_number", "latest_version_entity_id"} {
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
// for AgentConfigurationEntity does not return managedEntityType — the field has
// conflicting nullability across entity types in the union and the API rejects
// any query that includes it. All other ForceNew fields (name, agent_type,
// operating_system, organization_id) are available from GetEntity and are set
// here so Terraform does not plan spurious replacements after import.
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
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fleet configuration entity %s: %w", guid, err)
	}
	if entityInterface == nil {
		return nil, fmt.Errorf("fleet configuration entity %s not found", guid)
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

	versions := d.Get("version").([]interface{})

	// Create configuration atomically with first version (v1).
	firstVersion := versions[0].(map[string]interface{})
	configBody := []byte(firstVersion["configuration_content"].(string))

	entityMeta := fleetConfigBuildEntityMeta(
		d.Get("name").(string),
		d.Get("agent_type").(string),
		d.Get("managed_entity_type").(string),
		d.Get("operating_system").(string),
	)
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
		configBody,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": entityMeta,
			},
		},
		organizationID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ConfigurationEntityGUID)

	// content → version data map for state assembly
	contentToVersion := map[string]map[string]interface{}{
		firstVersion["configuration_content"].(string): {
			"configuration_content": firstVersion["configuration_content"],
			"version_number":        result.ConfigurationVersion.ConfigurationVersionNumber,
			"version_entity_id":     result.ConfigurationVersion.ConfigurationVersionEntityGUID,
		},
	}

	// Create additional versions (v2+).
	for i := 1; i < len(versions); i++ {
		vMap := versions[i].(map[string]interface{})
		content := vMap["configuration_content"].(string)

		vr, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
			[]byte(content),
			map[string]interface{}{
				"x-newrelic-client-go-custom-headers": map[string]string{
					"Newrelic-Entity": fmt.Sprintf(`{"agentConfiguration": "%s"}`, result.ConfigurationEntityGUID),
				},
			},
			organizationID,
		)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to create version %d: %w", i+1, err))
		}
		contentToVersion[content] = map[string]interface{}{
			"configuration_content": content,
			"version_number":        vr.ConfigurationVersion.ConfigurationVersionNumber,
			"version_entity_id":     vr.ConfigurationVersion.ConfigurationVersionEntityGUID,
		}
	}

	versionsList := fleetConfigBuildVersionsList(versions, contentToVersion)

	if err := d.Set("version", versionsList); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configuration_id", result.ConfigurationEntityGUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", organizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("total_versions", len(versionsList)); err != nil {
		return diag.FromErr(err)
	}

	fleetConfigSetLatestVersionFields(d, versionsList)
	// Call Read to ensure all computed fields (version_entity_id, version_number) are
	// populated from the API in a single apply — avoids the output-only drift that
	// otherwise requires a second apply to resolve.
	return resourceNewRelicFleetConfigurationRead(ctx, d, meta)
}

func resourceNewRelicFleetConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Always sync computed fields that can be derived without an extra API call.
	if setErr := d.Set("configuration_id", d.Id()); setErr != nil {
		return diag.FromErr(setErr)
	}
	if setErr := d.Set("organization_id", organizationID); setErr != nil {
		return diag.FromErr(setErr)
	}

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

	// Sort API versions by version number (ascending) so state order is stable.
	apiVersions := make([]fleetcontrol.ConfigurationVersion, len(versionsResp.Versions))
	copy(apiVersions, versionsResp.Versions)
	sort.Slice(apiVersions, func(i, j int) bool {
		ni, _ := strconv.Atoi(apiVersions[i].Version)
		nj, _ := strconv.Atoi(apiVersions[j].Version)
		return ni < nj
	})

	// Track the highest version number and build API entity-id set.
	apiVersionSet := make(map[string]bool, len(apiVersions))
	var latestVersionNum int
	var latestVersionEntityID string
	for _, v := range apiVersions {
		apiVersionSet[v.EntityGUID] = true
		num, parseErr := strconv.Atoi(v.Version)
		if parseErr != nil {
			return diag.FromErr(fmt.Errorf("failed to parse version number %q: %w", v.Version, parseErr))
		}
		if num > latestVersionNum {
			latestVersionNum = num
			latestVersionEntityID = v.EntityGUID
		}
	}

	// Index current state versions by entity_id so we can carry content forward.
	stateVersions := d.Get("version").([]interface{})
	stateByEntityID := make(map[string]map[string]interface{}, len(stateVersions))
	for _, sv := range stateVersions {
		svMap := sv.(map[string]interface{})
		entityID, _ := svMap["version_entity_id"].(string)
		if entityID != "" {
			stateByEntityID[entityID] = svMap
		}
	}

	// Warn about state versions that no longer exist in the API (externally deleted).
	var diags diag.Diagnostics
	for i, sv := range stateVersions {
		svMap := sv.(map[string]interface{})
		entityID, _ := svMap["version_entity_id"].(string)
		if entityID != "" && !apiVersionSet[entityID] {
			log.Printf("[WARN] Version %s (index %d) no longer exists in API, removing from state", entityID, i)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Version at index %d was deleted externally", i),
				Detail: fmt.Sprintf(
					"Version entity %s (index %d, version_number=%v) was not found in the API — "+
						"it was likely deleted outside of Terraform.\n\n"+
						"If this deletion was intentional, remove the corresponding version block "+
						"from your configuration to avoid it being recreated on the next apply.\n\n"+
						"If it was not intentional, run 'terraform apply' to recreate it.",
					entityID, i, svMap["version_number"],
				),
			})
		}
	}

	// Rebuild version list from API order. For versions already in state, carry content
	// forward. For versions not in state (import case), fetch content from the API.
	syncedVersions := make([]map[string]interface{}, 0, len(apiVersions))
	for _, apiVer := range apiVersions {
		versionNum, _ := strconv.Atoi(apiVer.Version)
		if svMap, ok := stateByEntityID[apiVer.EntityGUID]; ok {
			syncedVersions = append(syncedVersions, svMap)
		} else {
			// Version not in state — fetch content so import reconstructs full state.
			content, fetchErr := providerConfig.NewClient.FleetControl.FleetControlGetConfiguration(
				apiVer.EntityGUID, organizationID, fleetcontrol.GetConfigurationModeTypes.ConfigVersionEntity, 0,
			)
			var contentStr string
			if fetchErr != nil {
				log.Printf("[WARN] Could not fetch content for version entity %s: %v", apiVer.EntityGUID, fetchErr)
			} else if content != nil {
				contentStr = string(*content)
			}
			syncedVersions = append(syncedVersions, map[string]interface{}{
				"version_entity_id":     apiVer.EntityGUID,
				"version_number":        versionNum,
				"configuration_content": contentStr,
			})
		}
	}

	if err := d.Set("version", syncedVersions); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("total_versions", len(versionsResp.Versions)); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("latest_version_entity_id", latestVersionEntityID); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("latest_version_number", latestVersionNum); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceNewRelicFleetConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("name") {
		return diag.Errorf("configuration name cannot be updated — create a new configuration resource instead")
	}

	if !d.HasChange("version") {
		return nil
	}

	oldRaw, newRaw := d.GetChange("version")
	oldList := oldRaw.([]interface{})
	newList := newRaw.([]interface{})

	providerConfig := meta.(*ProviderConfig)
	organizationID := d.Get("organization_id").(string)

	// Index old state by content so we can look up entity_ids for deletion and
	// identify which contents already exist (no API call needed).
	// For the edge case of duplicate content in old state (should not happen after
	// CustomizeDiff, but defensive), first occurrence wins.
	oldByContent := make(map[string]map[string]interface{})
	for _, v := range oldList {
		vMap := v.(map[string]interface{})
		content := vMap["configuration_content"].(string)
		if _, exists := oldByContent[content]; !exists {
			oldByContent[content] = vMap
		}
	}

	// Determine which content values are present in the new desired list.
	newContentSet := make(map[string]bool)
	for _, v := range newList {
		newContentSet[v.(map[string]interface{})["configuration_content"].(string)] = true
	}

	// Delete versions whose content no longer appears in the new list.
	for content, vMap := range oldByContent {
		if newContentSet[content] {
			continue
		}
		entityID, _ := vMap["version_entity_id"].(string)
		if entityID == "" {
			continue
		}
		log.Printf("[INFO] Deleting removed version %s", entityID)
		if err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfigurationVersion(
			entityID, organizationID,
		); err != nil {
			return diag.FromErr(fmt.Errorf("failed to delete version %s: %w", entityID, err))
		}
	}

	// Create versions whose content does not exist in old state.
	newVersionData := make(map[string]map[string]interface{})
	for _, v := range newList {
		vMap := v.(map[string]interface{})
		content := vMap["configuration_content"].(string)
		if _, exists := oldByContent[content]; exists {
			continue
		}

		vr, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
			[]byte(content),
			map[string]interface{}{
				"x-newrelic-client-go-custom-headers": map[string]string{
					"Newrelic-Entity": fmt.Sprintf(`{"agentConfiguration": "%s"}`, d.Id()),
				},
			},
			organizationID,
		)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to create new version: %w", err))
		}

		newVersionData[content] = map[string]interface{}{
			"configuration_content": content,
			"version_number":        vr.ConfigurationVersion.ConfigurationVersionNumber,
			"version_entity_id":     vr.ConfigurationVersion.ConfigurationVersionEntityGUID,
		}
	}

	// Merge old (unchanged) and new (just created) version data, then assemble
	// final state ordered by the new desired list.
	mergedData := make(map[string]map[string]interface{}, len(oldByContent)+len(newVersionData))
	for k, v := range oldByContent {
		mergedData[k] = v
	}
	for k, v := range newVersionData {
		mergedData[k] = v
	}

	versionsList := fleetConfigBuildVersionsList(newList, mergedData)

	if err := d.Set("version", versionsList); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("total_versions", len(versionsList)); err != nil {
		return diag.FromErr(err)
	}

	fleetConfigSetLatestVersionFields(d, versionsList)
	// Call Read to ensure all computed fields are populated from the API after
	// the update — avoids the output-only drift requiring a second apply.
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

	// Delete all versions first — the API rejects configuration deletion if any version still exists.
	for _, v := range d.Get("version").([]interface{}) {
		vMap := v.(map[string]interface{})
		entityID, _ := vMap["version_entity_id"].(string)
		if entityID == "" {
			continue
		}
		log.Printf("[INFO] Deleting version %s", entityID)
		if err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfigurationVersion(
			entityID, organizationID,
		); err != nil {
			return diag.FromErr(fmt.Errorf("failed to delete version %s: %w", entityID, err))
		}
	}

	// Delete the configuration entity.
	// The API auto-deletes the config when its last version is removed, so NotFound is a success.
	_, err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfiguration(d.Id(), organizationID)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			log.Printf("[INFO] Fleet configuration %s already gone (auto-removed by API)", d.Id())
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

// fleetConfigBuildVersionsList assembles the []map to write into state.
// It iterates the desired version elements and backfills computed fields
// (version_entity_id, version_number) from contentToVersion, keyed by configuration_content.
func fleetConfigBuildVersionsList(versions []interface{}, contentToVersion map[string]map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(versions))
	for _, v := range versions {
		content := v.(map[string]interface{})["configuration_content"].(string)
		if known, ok := contentToVersion[content]; ok {
			out = append(out, known)
		} else {
			// Should not happen in normal flow; preserve whatever came in.
			out = append(out, v.(map[string]interface{}))
		}
	}
	return out
}

// fleetConfigSetLatestVersionFields scans all versions and sets latest_version_number
// and latest_version_entity_id to the entry with the highest version_number.
// Required because the API returns versions in unsorted order.
func fleetConfigSetLatestVersionFields(d *schema.ResourceData, versions []map[string]interface{}) {
	var latestNum int
	var latestEntityID string
	for _, v := range versions {
		num, _ := v["version_number"].(int)
		if num > latestNum {
			latestNum = num
			latestEntityID, _ = v["version_entity_id"].(string)
		}
	}
	_ = d.Set("latest_version_number", latestNum)
	_ = d.Set("latest_version_entity_id", latestEntityID)
}

// fleetConfigBuildEntityMeta builds the JSON string for the Newrelic-Entity header.
// When operatingSystem is non-empty, the operatingSystem object is included.
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
