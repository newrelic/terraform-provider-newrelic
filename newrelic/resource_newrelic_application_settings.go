package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/newrelic/newrelic-client-go/v2/pkg/agentapplications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func resourceNewRelicApplicationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicApplicationSettingsCreate,
		ReadContext:   resourceNewRelicApplicationSettingsRead,
		UpdateContext: resourceNewRelicApplicationSettingsUpdate,
		DeleteContext: resourceNewRelicApplicationSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			applicationSettingCommonSchema(),
			apmApplicationSettingsSchema(),
		),
		CustomizeDiff: validateApplicationSettingsInput,
	}
}

func resourceNewRelicApplicationSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return resourceNewRelicApplicationSettingsUpdate(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic application %+v", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("no New Relic application found with given guid %s", d.Id()))
	}

	var dig diag.Diagnostics
	switch (*resp).(type) {
	case *entities.ApmApplicationEntity:
		entity := (*resp).(*entities.ApmApplicationEntity)
		d.SetId(string(entity.GUID))
		_ = d.Set("guid", string(entity.GUID))
		dig = diag.FromErr(setAPMApplicationValues(d, entity.ApmSettings))
	default:
		dig = diag.FromErr(fmt.Errorf("problem in retrieving application with GUID %s", d.Id()))
	}
	return dig
}

func resourceNewRelicApplicationSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	_, guidExists := d.GetOk("guid")
	_, nameExists := d.GetOk("name")
	if nameExists && !guidExists {
		entityRes, err := getEntityDetailsFromName(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		if entityRes == nil {
			return diag.FromErr(fmt.Errorf("no entities found with the provided name, please ensure the name of a valid APM entity is provided"))
		}
		_ = d.Set("guid", string((*entityRes).GetGUID()))
	}

	updateApplicationParams := expandApplication(d)

	guid := d.Get("guid").(string)
	log.Printf("[INFO] Updating New Relic application %+v with params: %+v", guid, updateApplicationParams)

	agentApplicationSettingResult, err := client.AgentApplications.AgentApplicationSettingsUpdate(common.EntityGUID(guid), *updateApplicationParams)

	if err != nil {
		return diag.FromErr(err)
	}
	if agentApplicationSettingResult == nil {
		return diag.FromErr(fmt.Errorf("something went wrong while Updating New Relic application"))
	}

	d.SetId(string(agentApplicationSettingResult.GUID))
	err = d.Set("is_imported", true)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicApplicationSettingsRead(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ProviderConfig).NewClient
	guid := d.Get("guid").(string)
	falseValue := false

	agentApplicationSettingResult, err := client.AgentApplications.AgentApplicationSettingsUpdate(
		common.EntityGUID(guid),
		agentapplications.AgentApplicationSettingsUpdateInput{
			ApmConfig: &agentapplications.AgentApplicationSettingsApmConfigInput{
				UseServerSideConfig: &falseValue,
				ApdexTarget:         0.5,
			},
			ErrorCollector: &agentapplications.AgentApplicationSettingsErrorCollectorInput{
				Enabled: &falseValue,
			},
			TransactionTracer: &agentapplications.AgentApplicationSettingsTransactionTracerInput{
				Enabled:        &falseValue,
				ExplainEnabled: &falseValue,
				RecordSql:      agentapplications.AgentApplicationSettingsRecordSqlEnumTypes.OFF,
			},
			TracerType: &agentapplications.AgentApplicationSettingsTracerTypeInput{
				// choosing "NONE" instead of "OPT_OUT", as OPT_OUT is not available as a type in the Go Client
				// since OPT_OUT is not a recognized enum value by NerdGraph
				Value: agentapplications.AgentApplicationSettingsTracerTypes.NONE,
			},
			ThreadProfiler: &agentapplications.AgentApplicationSettingsThreadProfilerInput{
				Enabled: &falseValue,
			},
			SlowSql: &agentapplications.AgentApplicationSettingsSlowSqlInput{
				Enabled: &falseValue,
			},
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if agentApplicationSettingResult == nil {
		return diag.FromErr(fmt.Errorf("something went wrong while clearing the settings of the application"))
	}

	// this function should not be removed, this exists in every delete function to "clear" Terraform state
	d.SetId("")

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "\nSince the `newrelic_application_settings` resource does not support deletion via NerdGraph, the resource has been reset to its initial state and cleared of most settings.\nYou would still find some settings of this APM entity in the New Relic UI (such as the Apdex Threshold) which cannot be cleared.",
	})
	return diags
}
