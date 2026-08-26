package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
)

func resourceNewRelicPathpointFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicPathpointFlowCreate,
		ReadContext:   resourceNewRelicPathpointFlowRead,
		UpdateContext: resourceNewRelicPathpointFlowUpdate,
		DeleteContext: resourceNewRelicPathpointFlowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNewRelicPathpointFlowImport,
		},
		Schema: pathpointFlowSchema(),
	}
}

func resourceNewRelicPathpointFlowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	flowInput := expandPathpointFlowInput(d, accountID)
	scopeInput := pathpoint.PathPointScopeInput{
		ID:   accountID,
		Type: pathpoint.PathPointScopeTypeTypes.ACCOUNT,
	}

	result, err := client.PathPoint.PathPointCreate(flowInput, scopeInput)
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("error creating Pathpoint flow: empty response"))
	}

	d.SetId(string(result.GUID))
	log.Printf("[INFO] Created Pathpoint flow %s (GUID: %s)", result.Name, result.GUID)

	return flattenPathpointFlowResult(d, result)
}

func resourceNewRelicPathpointFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := pathpoint.EntityGUID(d.Id())
	accountID := selectAccountID(providerConfig, d)

	result, err := client.PathPoint.GetFlow(accountID, guid)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("error reading Pathpoint flow (GUID: %s): %w", guid, err))
	}
	if result == nil || string(result.GUID) == "" {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", accountID)
	return flattenPathpointFlowResult(d, result)
}

func resourceNewRelicPathpointFlowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := pathpoint.EntityGUID(d.Id())

	versionStr := d.Get("version").(string)
	versionInt, _ := strconv.ParseInt(versionStr, 10, 64)
	version := nrtime.EpochMilliseconds(time.UnixMilli(versionInt))

	accountID := selectAccountID(providerConfig, d)
	updateInput := expandPathpointFlowUpdateInput(d, version, accountID)

	result, err := client.PathPoint.PathPointUpdate(guid, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("error updating Pathpoint flow: empty response"))
	}

	log.Printf("[INFO] Updated Pathpoint flow %s (GUID: %s)", result.Name, result.GUID)

	return flattenPathpointFlowResult(d, result)
}

func resourceNewRelicPathpointFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := pathpoint.EntityGUID(d.Id())

	_, err := client.PathPoint.PathPointDelete(guid)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted Pathpoint flow (GUID: %s)", guid)
	return nil
}

func flattenPathpointFlowResult(d *schema.ResourceData, result *pathpoint.PathPointFlowResult) diag.Diagnostics {
	_ = d.Set("guid", string(result.GUID))
	_ = d.Set("name", result.Name)
	_ = d.Set("version", strconv.FormatInt(time.Time(result.Version).UnixMilli(), 10))
	_ = d.Set("description", result.Description)
	_ = d.Set("category", result.Category)

	_ = d.Set("health_rollup", string(result.HealthRollup))
	_ = d.Set("refresh_interval", string(result.RefreshInterval))
	if result.Message != "" {
		log.Printf("[WARN] Pathpoint flow response message: %s", result.Message)
	}

	if err := d.Set("kpis", flattenPathpointKpis(result.Kpis)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting kpis: %w", err))
	}

	if err := d.Set("stages", flattenPathpointStages(result.Stages.Items)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting stages: %w", err))
	}

	return nil
}
