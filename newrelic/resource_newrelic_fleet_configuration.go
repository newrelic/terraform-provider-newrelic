package newrelic

import (
	"context"
	"fmt"
	"log"
	"os"

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
			StateContext: schema.ImportStatePassthroughContext,
		},
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
					"INFRASTRUCTURE",
					"KUBERNETES",
				}, false),
				Description: "The type of agent this configuration is for. Allowed values: INFRASTRUCTURE, KUBERNETES.",
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
			"configuration_file_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"configuration_file_path", "configuration_content"},
				Description:  "Path to a file containing the configuration content. Exactly one of configuration_file_path or configuration_content must be provided.",
			},
			"configuration_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"configuration_file_path", "configuration_content"},
				Description:  "Inline configuration content. Exactly one of configuration_file_path or configuration_content must be provided.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. If not provided, it will be auto-fetched from the account.",
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the configuration.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current version number of the configuration.",
			},
		},
	}
}

func resourceNewRelicFleetConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Get or fetch organization ID
	organizationID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate that exactly one of configuration_file_path or configuration_content is provided
	configFilePath, hasFilePath := d.GetOk("configuration_file_path")
	configContent, hasContent := d.GetOk("configuration_content")

	if !hasFilePath && !hasContent {
		return diag.Errorf("one of configuration_file_path or configuration_content must be provided")
	}

	if hasFilePath && hasContent {
		return diag.Errorf("configuration_file_path and configuration_content are mutually exclusive, use only one")
	}

	// Read configuration content
	var configBody []byte
	if hasFilePath {
		fileContent, readErr := os.ReadFile(configFilePath.(string))
		if readErr != nil {
			return diag.Errorf("failed to read configuration file: %v", readErr)
		}
		configBody = fileContent
	} else {
		configBody = []byte(configContent.(string))
	}

	// Build custom headers required by the API
	// These headers specify the entity name, agent type, and managed entity type
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": fmt.Sprintf(
				`{"name": "%s", "agentType": "%s", "managedEntityType": "%s"}`,
				d.Get("name").(string),
				d.Get("agent_type").(string),
				d.Get("managed_entity_type").(string),
			),
		},
	}

	// Create configuration
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
		configBody,
		customHeaders,
		organizationID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID and attributes
	d.SetId(result.ConfigurationEntityGUID)

	// Flatten the result
	if err := flattenFleetConfiguration(result, d, organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Get organization ID from state, or fetch if not present (e.g., during import)
	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Fetch configuration entity
	// Use ConfigEntity mode to get the configuration entity
	mode := fleetcontrol.GetConfigurationModeTypes.ConfigEntity
	version := 0 // Not used for ConfigEntity mode

	configContent, err := providerConfig.NewClient.FleetControl.FleetControlGetConfiguration(
		d.Id(),
		organizationID,
		mode,
		version,
	)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if configContent == nil {
		d.SetId("")
		return nil
	}

	// The API only returns the raw configuration content, not the metadata
	// We need to preserve the values already in state
	// Set organization_id if it was fetched
	if err := d.Set("organization_id", organizationID); err != nil {
		return diag.FromErr(err)
	}

	// Preserve other fields from state as API doesn't return them
	// name, agent_type, managed_entity_type, configuration_content, configuration_file_path
	// are already in state and don't need to be updated unless we got new data from API
	// (which we don't for this API)

	return nil
}

func resourceNewRelicFleetConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Only name can be updated directly
	// Configuration content changes create a new version, not handled here

	// Check if name changed
	if !d.HasChange("name") {
		return nil
	}

	// Note: The API doesn't have an update configuration name endpoint
	// This would need to be implemented via delete and recreate
	// For now, we'll return an error suggesting the user use a new resource

	return diag.Errorf("configuration name cannot be updated. To change the configuration, create a new resource or add a new version using newrelic_fleet_configuration_version")
}

func resourceNewRelicFleetConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Get organization ID from state
	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Deleting New Relic Fleet Configuration %s", d.Id())

	_, err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfiguration(
		d.Id(),
		organizationID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
