package newrelic

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func resourceNewRelicFleetConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetConfigurationVersionCreate,
		ReadContext:   resourceNewRelicFleetConfigurationVersionRead,
		DeleteContext: resourceNewRelicFleetConfigurationVersionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID (GUID) of the configuration to add this version to.",
			},
			"configuration_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"configuration_content"},
				Description:   "Path to a file containing the configuration content. Mutually exclusive with configuration_content.",
			},
			"configuration_content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"configuration_file_path"},
				Description:   "Inline configuration content. Mutually exclusive with configuration_file_path.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. If not provided, it will be auto-fetched from the account.",
			},
			"version_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version number of this configuration version.",
			},
		},
	}
}

func resourceNewRelicFleetConfigurationVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Build custom headers with the configuration GUID
	// This tells the API to add a version to the existing configuration
	// Different from create which sends name/agentType/managedEntityType
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": fmt.Sprintf(
				`{"agentConfiguration": "%s"}`,
				d.Get("configuration_id").(string),
			),
		},
	}

	// Add version using the same create API but with different headers
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfiguration(
		configBody,
		customHeaders,
		organizationID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID to the version entity GUID
	d.SetId(result.ConfigurationVersion.ConfigurationVersionEntityGUID)

	// Flatten the result
	if err := flattenFleetConfigurationVersion(result, d, organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetConfigurationVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Get configuration ID from state
	configurationID := d.Get("configuration_id").(string)
	if configurationID == "" {
		return diag.Errorf("configuration_id is required for read operation")
	}

	// Fetch the specific version by version number
	// Use ConfigVersionEntity mode to get version-specific data
	mode := fleetcontrol.GetConfigurationModeTypes.ConfigVersionEntity

	// Get version number from state if available, otherwise use 0 to get latest
	versionNumber := d.Get("version_number").(int)

	_, err := providerConfig.NewClient.FleetControl.FleetControlGetConfiguration(
		d.Id(),
		organizationID,
		mode,
		versionNumber,
	)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Set organization_id if it was fetched
	if err := d.Set("organization_id", organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetConfigurationVersionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	log.Printf("[INFO] Deleting New Relic Fleet Configuration Version %s", d.Id())

	err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfigurationVersion(
		d.Id(),
		organizationID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
