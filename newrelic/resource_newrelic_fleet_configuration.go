package newrelic

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicFleetConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetConfigurationCreate,
		ReadContext:   resourceNewRelicFleetConfigurationRead,
		DeleteContext: resourceNewRelicFleetConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				Description: "The type of agent this configuration targets. Allowed values: NRInfra, NRDOT, FluentBit, NRPrometheusAgent.",
			},
			"managed_entity_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HOST",
					"KUBERNETESCLUSTER",
				}, false),
				Description: "The type of entities this configuration applies to. Allowed values: HOST, KUBERNETESCLUSTER.",
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
			"entity_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity GUID of the configuration.",
			},
			"blob_version_entity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the initial version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number.",
						},
						"guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version entity GUID.",
						},
						"blob_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The blob ID.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicFleetConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Build custom header for the request
	name := d.Get("name").(string)
	agentType := d.Get("agent_type").(string)
	managedEntityType := d.Get("managed_entity_type").(string)

	// Build custom headers required by the API
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": fmt.Sprintf(
				`{"name": "%s", "agentType": "%s", "managedEntityType": "%s"}`,
				name,
				agentType,
				managedEntityType,
			),
		},
	}

	// Call the configuration creation API
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateConfigurationWithContext(
		ctx,
		content,
		customHeaders,
		organizationID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create fleet configuration: %w", err))
	}

	// Set the entity GUID as the ID
	d.SetId(result.ConfigurationEntityGUID)

	// Set computed fields
	if err := d.Set("entity_guid", result.ConfigurationEntityGUID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", organizationID); err != nil {
		return diag.FromErr(err)
	}

	// Set blob version entity
	blobVersion := []map[string]interface{}{
		{
			"version": result.ConfigurationVersion.ConfigurationVersionNumber,
			"guid":    result.ConfigurationVersion.ConfigurationVersionEntityGUID,
			"blob_id": result.BlobId,
		},
	}
	if err := d.Set("blob_version_entity", blobVersion); err != nil {
		return diag.FromErr(err)
	}

	return nil
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

	// TODO: Verify configuration exists
	// For now, trust the state
	_ = organizationID
	return nil
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

	_, err := providerConfig.NewClient.FleetControl.FleetControlDeleteConfigurationWithContext(ctx, d.Id(), organizationID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
