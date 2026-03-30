package newrelic

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "The configuration entity ID to add a version to.",
			},
			"configuration_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"configuration_content"},
				Description:   "Path to the configuration file (JSON/YAML). Mutually exclusive with configuration_content.",
			},
			"configuration_content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"configuration_file"},
				Description:   "Inline configuration content (JSON/YAML). Mutually exclusive with configuration_file.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. If not provided, it will be auto-fetched.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version number.",
			},
			"blob_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The blob ID.",
			},
		},
	}
}

func resourceNewRelicFleetConfigurationVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Get organization ID
	organizationID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate that at least one configuration source is provided
	configFile, hasFile := d.GetOk("configuration_file")
	configContent, hasContent := d.GetOk("configuration_content")

	if !hasFile && !hasContent {
		return diag.Errorf("either configuration_file or configuration_content must be provided")
	}

	// Read configuration content
	var content []byte
	if hasFile {
		content, err = os.ReadFile(configFile.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to read configuration file: %w", err))
		}
	} else {
		content = []byte(configContent.(string))
	}

	configurationID := d.Get("configuration_id").(string)

	// Build custom headers for adding a version to existing configuration
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": fmt.Sprintf(
				`{"agentConfiguration": "%s"}`,
				configurationID,
			),
		},
	}

	// Call the configuration version creation API
	// Note: Same method as create but with different headers
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfigurationWithContext(
		ctx,
		content,
		customHeaders,
		organizationID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to add configuration version: %w", err))
	}

	// Set the version GUID as the ID
	d.SetId(result.ConfigurationVersion.ConfigurationVersionEntityGUID)

	if err := d.Set("version", result.ConfigurationVersion.ConfigurationVersionNumber); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("blob_id", result.BlobId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetConfigurationVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Version exists if it has an ID - these are immutable once created
	// No separate read API for a specific version, so we trust the state
	return nil
}

func resourceNewRelicFleetConfigurationVersionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Deleting New Relic Fleet Configuration Version %s", d.Id())

	err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfigurationVersionWithContext(ctx, d.Id(), organizationID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
