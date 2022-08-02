package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsStepMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsStepMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsStepMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsStepMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsStepMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			syntheticsMonitorCommonSchema(),
			syntheticsStepMonitorSchema(),
		),
	}
}

func syntheticsStepMonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_screenshot_on_failure_and_script": {
			Type:        schema.TypeBool,
			Description: "Capture a screenshot during job execution.",
			Optional:    true,
		},
		"locations_private": {
			Type:        schema.TypeSet,
			Description: "",
			Optional:    true,
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
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Optional:    true,
			Description: "The public location(s) that the monitor will run jobs from.",
		},
		"steps": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "The steps that make up the script the monitor will run",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ordinal": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "The position of the step within the script ranging from 0-100",
						//SchemaValidateDiagFunc: // TODO: add validation to ensure value is between 0 and 100 (inclusive)
					},
					"type": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The type of step to be added to the script.",
						//SchemaValidateDiagFunc: // TODO: add valid step types via enum values in client
					},
					"values": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Optional:    true,
						Description: "",
					},
				},
			},
		},
	}
}

func expandSyntheticsMonitorSteps(steps []interface{}) []synthetics.SyntheticsStepInput {
	stepsOut := []synthetics.SyntheticsStepInput{}

	for _, s := range steps {
		st := s.(map[string]interface{})

		stepsOut = append(stepsOut, synthetics.SyntheticsStepInput{
			Ordinal: st["ordinal"].(int),
			Type:    synthetics.SyntheticsStepType(st["type"].(string)),
			Values:  expandStringSlice(st["values"].([]interface{})),
		})
	}

	return stepsOut
}

func expandPrivateLocations(locations []interface{}) []synthetics.SyntheticsPrivateLocationInput {
	pl := []synthetics.SyntheticsPrivateLocationInput{}

	for _, v := range locations {
		loc := v.(map[string]string)
		pl = append(pl, synthetics.SyntheticsPrivateLocationInput{
			GUID:        loc["guid"],
			VsePassword: synthetics.SecureValue(loc["vse_password"]),
		})
	}

	return pl
}

func buildSyntheticsStepMonitorCreateInput(d *schema.ResourceData) *synthetics.SyntheticsCreateStepMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateStepMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Steps:  expandSyntheticsMonitorSteps(d.Get("steps").([]interface{})),
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandPrivateLocations(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		v := attr.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &v
	}

	return &input
}

func buildSyntheticsStepMonitorUpdateInput(d *schema.ResourceData) *synthetics.SyntheticsUpdateStepMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateStepMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Steps:  expandSyntheticsMonitorSteps(d.Get("steps").([]interface{})),
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandPrivateLocations(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		v := attr.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &v
	}

	return &input
}

func resourceNewRelicSyntheticsStepMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	monitorInput := buildSyntheticsStepMonitorCreateInput(d)
	resp, err := client.Synthetics.SyntheticsCreateStepMonitorWithContext(ctx, accountID, *monitorInput)
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
	_ = d.Set("locations_public", resp.Monitor.Locations.Public)

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"guid":   string(resp.Monitor.GUID),
		"name":   resp.Monitor.Name,
		"period": string(resp.Monitor.Period),
		"status": string(resp.Monitor.Status),
	})

	return diag.FromErr(err)
}

func flattenSyntheticsMonitorSteps(stepsIn []synthetics.SyntheticsStep) []map[string]interface{} {
	steps := []map[string]interface{}{}

	// Note: This might need further flattening for TF
	for _, s := range stepsIn {
		step := map[string]interface{}{
			"ordinal": s.Ordinal,
			"type":    string(s.Type),
			"values":  s.Values,
		}

		steps = append(steps, step)
	}

	return steps
}

func resourceNewRelicSyntheticsStepMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		entity := (*resp).(*entities.SyntheticMonitorEntity)
		stepsResp, errr := client.Synthetics.GetSteps(accountID, synthetics.EntityGUID(entity.GetGUID()))
		if err != nil {
			return diag.FromErr(errr)
		}

		steps := flattenSyntheticsMonitorSteps(*stepsResp)

		d.SetId(string(e.GUID))
		_ = d.Set("account_id", accountID)
		_ = d.Set("locations_public", getPublicLocationsFromEntityTags(entity.GetTags()))
		_ = d.Set("steps", steps)

		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"guid":   string(e.GUID),
			"name":   entity.Name,
			"period": string(syntheticsMonitorPeriodValueMap[int(entity.GetPeriod())]),
			"status": string(entity.MonitorSummary.Status),
		})
	}

	return diag.FromErr(err)
}

func resourceNewRelicSyntheticsStepMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	monitorInput := buildSyntheticsStepMonitorUpdateInput(d)
	resp, err := client.Synthetics.SyntheticsUpdateStepMonitorWithContext(ctx, synthetics.EntityGUID(d.Id()), *monitorInput)
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildUpdateSyntheticsMonitorResponseErrors(resp.Errors)
	if len(errors) > 0 {
		return errors
	}

	_ = d.Set("locations_public", resp.Monitor.Locations.Public)
	_ = d.Set("steps", flattenSyntheticsMonitorSteps(resp.Monitor.Steps))

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"guid":   string(resp.Monitor.GUID),
		"name":   resp.Monitor.Name,
		"period": string(resp.Monitor.Period),
		"status": string(resp.Monitor.Status),
	})

	return diag.FromErr(err)
}

func resourceNewRelicSyntheticsStepMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	_, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.FromErr(err)
}
