package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	"log"
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
		CustomizeDiff: validateSyntheticMonitorRuntimeAttributes,
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

// Scripted browser monitors have advanced options, browsers and devices fields, but scripted API monitors do not.
func syntheticsScriptBrowserMonitorAdvancedOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_screenshot_on_failure_and_script": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Capture a screenshot during job execution.",
		},
		"device_orientation": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The device orientation the user would like to represent. Valid values are LANDSCAPE, PORTRAIT, or NONE.",
		},
		"device_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The device type that a user can select. Valid values are MOBILE, TABLET, or NONE.",
		},
		"browsers": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			MinItems: 1,
			Optional: true,
			Description: "The browsers that can be used to execute script execution. Valid values are array of CHROME," +
				" EDGE, FIREFOX, and NONE.",
		},
		"devices": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			MinItems: 1,
			Optional: true,
			Description: "The devices that can be used to execute script execution. Valid values are array of DESKTOP," +
				" MOBILE_LANDSCAPE, MOBILE_PORTRAIT, TABLET_LANDSCAPE, TABLET_PORTRAIT and NONE.",
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
		},
		"runtime_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The runtime type that the monitor will run.",
		},
		"runtime_type_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The specific semver version of the runtime type.",
		},
		SyntheticsUseLegacyRuntimeAttrLabel: SyntheticsUseLegacyRuntimeSchema,
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
		_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

		attrs := map[string]string{"guid": string(resp.Monitor.GUID)}
		if err = setSyntheticsMonitorAttributes(d, attrs); err != nil {
			return diag.FromErr(err)
		}
	case string(SyntheticsMonitorTypes.SCRIPT_BROWSER):
		monitorInput := buildSyntheticsScriptBrowserMonitorInput(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptBrowserMonitorWithContext(ctx, accountID, monitorInput)
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
		_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

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

	// This should probably be in go-client, so we can use *errors.NotFound
	if *resp == nil {
		d.SetId("")
		return nil
	}

	response, error := client.Synthetics.GetScript(accountID, synthetics.EntityGUID(d.Id()))
	if error != nil {
		return diag.FromErr(error)
	}

	if response == nil {
		d.SetId("")
		return nil
	}

	error = setSyntheticsMonitorAttributes(d, map[string]string{
		"script": response.Text,
	})

	if error != nil {
		return diag.FromErr(error)
	}

	_ = d.Set("account_id", accountID)

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"name":   e.Name,
			"type":   string(e.MonitorType),
			"guid":   string(e.GUID),
			"period": string(syntheticsMonitorPeriodValueMap[int(e.GetPeriod())]),
			"status": string(e.MonitorSummary.Status),
		})

		_ = d.Set("period_in_minutes", int(e.GetPeriod()))
		for _, t := range e.Tags {
			if k, ok := syntheticsMonitorTagKeyToSchemaAttrMap[t.Key]; ok {
				if t.Key == "devices" || t.Key == "browsers" {
					_ = d.Set(k, t.Values)
				} else if len(t.Values) == 1 {
					_ = d.Set(k, t.Values[0])
				}
			}
		}
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
		log.Printf("No monitor type specified")
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
			"name":   resp.Monitor.Name,
			"guid":   string(resp.Monitor.GUID),
			"period": string(resp.Monitor.Period),
			"status": string(resp.Monitor.Status),
		})

		_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

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
			"name":   resp.Monitor.Name,
			"guid":   string(resp.Monitor.GUID),
			"period": string(resp.Monitor.Period),
			"status": string(resp.Monitor.Status),
		})

		_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

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
