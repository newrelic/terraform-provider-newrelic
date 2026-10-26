package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func resourceNewRelicSyntheticsBrokenLinksMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsBrokenLinksMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsBrokenLinksMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsBrokenLinksMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsBrokenLinksMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			syntheticsBrokenLinksMonitorSchema(),
			syntheticsMonitorCommonSchema(),
			syntheticsMonitorLocationsAsStringsSchema(),
			syntheticsMonitorRuntimeOptions(),
		),
		CustomizeDiff: validateSyntheticMonitorAttributes,
	}
}

func syntheticsBrokenLinksMonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uri": {
			Type:        schema.TypeString,
			Description: "The URI the monitor runs against.",
			Required:    true,
		},
	}
}

func resourceNewRelicSyntheticsBrokenLinksMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	monitorInput, monitorInputErr := buildSyntheticsBrokenLinksMonitorCreateInput(d)
	if monitorInputErr != nil {
		return diag.FromErr(monitorInputErr)
	}
	resp, err := client.Synthetics.SyntheticsCreateBrokenLinksMonitorWithContext(ctx, accountID, *monitorInput)
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

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"guid":       string(resp.Monitor.GUID),
		"name":       resp.Monitor.Name,
		"period":     string(resp.Monitor.Period),
		"status":     string(resp.Monitor.Status),
		"uri":        resp.Monitor.Uri,
		"monitor_id": resp.Monitor.ID,
	})

	respRuntimeType := resp.Monitor.Runtime.RuntimeType
	respRuntimeTypeVersion := resp.Monitor.Runtime.RuntimeTypeVersion

	if respRuntimeType != "" {
		_ = d.Set("runtime_type", respRuntimeType)
	}

	if respRuntimeTypeVersion != "" {
		_ = d.Set("runtime_type_version", respRuntimeTypeVersion)
	}

	return diag.FromErr(err)
}

func resourceNewRelicSyntheticsBrokenLinksMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		entity := (*resp).(*entities.SyntheticMonitorEntity)

		d.SetId(string(e.GUID))
		_ = d.Set("account_id", accountID)
		_ = d.Set("locations_public", getPublicLocationsFromEntityTags(entity.GetTags()))
		_ = d.Set("period_in_minutes", int(entity.GetPeriod()))

		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"guid":       string(e.GUID),
			"name":       entity.Name,
			"period":     string(syntheticsMonitorPeriodValueMap[int(entity.GetPeriod())]),
			"status":     string(entity.MonitorSummary.Status),
			"uri":        entity.MonitoredURL,
			"monitor_id": entity.MonitorId,
		})

		runtimeType, runtimeTypeVersion := getRuntimeValuesFromEntityTags(entity.GetTags())
		if runtimeType != "" && runtimeTypeVersion != "" {
			_ = d.Set("runtime_type", runtimeType)
			_ = d.Set("runtime_type_version", runtimeTypeVersion)
		}
	}

	return diag.FromErr(err)
}

func resourceNewRelicSyntheticsBrokenLinksMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	guid := synthetics.EntityGUID(d.Id())

	monitorInput, monitorInputErr := buildSyntheticsBrokenLinksMonitorUpdateInput(d)
	if monitorInputErr != nil {
		return diag.FromErr(monitorInputErr)
	}
	resp, err := client.Synthetics.SyntheticsUpdateBrokenLinksMonitorWithContext(ctx, guid, *monitorInput)
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildUpdateSyntheticsMonitorResponseErrors(resp.Errors)
	if len(errors) > 0 {
		return errors
	}

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"guid":   string(resp.Monitor.GUID),
		"name":   resp.Monitor.Name,
		"period": string(resp.Monitor.Period),
		"status": string(resp.Monitor.Status),
		"uri":    resp.Monitor.Uri,
	})

	_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

	respRuntimeType := resp.Monitor.Runtime.RuntimeType
	respRuntimeTypeVersion := resp.Monitor.Runtime.RuntimeTypeVersion

	if respRuntimeType != "" {
		_ = d.Set("runtime_type", respRuntimeType)
	}

	if respRuntimeTypeVersion != "" {
		_ = d.Set("runtime_type_version", respRuntimeTypeVersion)
	}

	return diag.FromErr(err)
}

// NOTE: We can make rename this to reusable function for all new monitor types,
//
//	but the legacy function already has a good generic name (`resourceNewRelicSyntheticsMonitorDelete()`)
func resourceNewRelicSyntheticsBrokenLinksMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	_, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}
