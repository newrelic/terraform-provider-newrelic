package newrelic

import (
	"context"
	"log"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func resourceNewRelicSyntheticsScriptMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsScriptMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsScriptMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsScriptMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsScriptMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			syntheticsMonitorCommonSchema(),
			syntheticsScriptMonitorCommonSchema(),
			syntheticsScriptMonitorLocationsSchema(),
			syntheticsScriptBrowserMonitorAdvancedOptionsSchema(),
		),
	}
}

func syntheticsScriptMonitorLocationsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location_private": {
			Type:         schema.TypeSet,
			Description:  "",
			Optional:     true, // Note: Optional
			AtLeastOneOf: []string{"location_private", "locations_public"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:        schema.TypeString,
						Description: "The unique identifier for the Synthetics private location in New Relic.",
						Required:    true,
					},
					"vse_password": {
						Type:        schema.TypeString,
						Description: "The location's Verified Script Execution password (Only necessary if Verified Script Execution is enabled for the location).",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},
		"locations_public": {
			Type:         schema.TypeSet,
			Elem:         &schema.Schema{Type: schema.TypeString},
			MinItems:     1,
			Optional:     true,
			Description:  "The public location(s) that the monitor will run jobs from.",
			AtLeastOneOf: []string{"location_private", "locations_public"},
		},
	}
}

// Scripted browser monitors have advanced options, but scripted API monitors do not.
func syntheticsScriptBrowserMonitorAdvancedOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_screenshot_on_failure_and_script": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Capture a screenshot during job execution.",
		},
	}
}

// Returns common schema attributes shared by both scripted browser and scripted API monitors.
func syntheticsScriptMonitorCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "The monitor type. Valid values are SCRIPT_BROWSER, and SCRIPT_API.",
			ValidateFunc: validation.StringInSlice(listValidSyntheticsScriptMonitorTypes(), false),
		},
		"script": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The script that the monitor runs.",
		},
		"script_language": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The programing language that should execute the script.",
			Default:     "JAVASCRIPT",
		},
		"runtime_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The runtime type that the monitor will run.",
			Default:     "CHROME_BROWSER",
		},
		"runtime_type_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The specific semver version of the runtime type.",
			Default:     "100",
		},
	}
}

// CREATE
func resourceNewRelicSyntheticsScriptMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	monitorType, ok := d.GetOk("type")
	if !ok {
		log.Printf("attribute `type` is required and must be one of 'SCRIPT_API' or 'SCRIPT_BROWSER'")
	}

	switch monitorType {
	case string(SyntheticsMonitorTypes.SCRIPT_API):
		monitorInput := buildSyntheticsScriptAPIMonitorInput(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptAPIMonitorWithContext(ctx, accountID, monitorInput)
		if err != nil {
			return diag.FromErr(err)
		}

		errors := buildCreateSyntheticsMonitorResponseErrors(resp.Errors)
		if len(errors) > 0 {
			return errors
		}

		// Set attributes
		d.SetId(string(resp.Monitor.GUID))
		_ = d.Set("account_id", accountID)
		attrs := map[string]string{"guid": string(resp.Monitor.GUID)}
		if err = setSyntheticsMonitorAttributes(d, attrs); err != nil {
			return diag.FromErr(err)
		}
	case string(SyntheticsMonitorTypes.SCRIPT_BROWSER):
		monitorInput := buildSyntheticsScriptBrowserMonitorInput(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptBrowserMonitorWithContext(ctx, accountID, monitorInput)
		if err != nil {
			diag.FromErr(err)
		}

		errors := buildCreateSyntheticsMonitorResponseErrors(resp.Errors)
		if len(errors) > 0 {
			return errors
		}

		// Set attributes
		d.SetId(string(resp.Monitor.GUID))
		_ = d.Set("account_id", accountID)
		attrs := map[string]string{"guid": string(resp.Monitor.GUID)}
		if err = setSyntheticsMonitorAttributes(d, attrs); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// READ
func resourceNewRelicSyntheticsScriptMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// This should probably be in go-client so we can use *errors.NotFound
	if *resp == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", accountID)

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"name": e.Name,
			"type": string(e.MonitorType),
			"guid": string(e.GUID),
		})
	}

	return diag.FromErr(err)
}

// UPDATE
func resourceNewRelicSyntheticsScriptMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	guid := synthetics.EntityGUID(d.Id())

	monitorType, ok := d.GetOk("type")
	if !ok {
		log.Printf("Not Monitor type specified")
	}

	switch monitorType {
	case string(SyntheticsMonitorTypes.SCRIPT_API):
		monitorInput := buildSyntheticsScriptAPIMonitorUpdateInput(d)
		resp, err := client.Synthetics.SyntheticsUpdateScriptAPIMonitorWithContext(ctx, guid, monitorInput)
		if err != nil {
			return diag.FromErr(err)
		}

		errors := buildUpdateSyntheticsMonitorResponseErrors(resp.Errors)
		if len(errors) > 0 {
			return errors
		}

		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"name": resp.Monitor.Name,
			"guid": string(resp.Monitor.GUID),
		})
		if err != nil {
			return diag.FromErr(err)
		}

	case string(SyntheticsMonitorTypes.SCRIPT_BROWSER):
		monitorInput := buildSyntheticsScriptBrowserUpdateInput(d)
		resp, err := client.Synthetics.SyntheticsUpdateScriptBrowserMonitorWithContext(ctx, guid, monitorInput)
		if err != nil {
			return diag.FromErr(err)
		}

		errors := buildUpdateSyntheticsMonitorResponseErrors(resp.Errors)
		if len(errors) > 0 {
			return errors
		}

		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"name": resp.Monitor.Name,
			"guid": string(resp.Monitor.GUID),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// DELETE
func resourceNewRelicSyntheticsScriptMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	_, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
