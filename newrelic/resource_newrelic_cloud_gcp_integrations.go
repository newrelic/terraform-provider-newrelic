package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewrelicCloudGcpIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewrelicCloudGcpIntegrationsCreate,
		ReadContext:   resourceNewrelicCloudGcpIntegrationsRead,
		UpdateContext: resourceNewrelicCloudGcpIntegrationsUpdate,
		DeleteContext: resourceNewrelicCloudGcpIntegrationsDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "ID of the newrelic account",
				Computed:    true,
				Optional:    true,
			},
			"linked_account_id": {
				Type:        schema.TypeInt,
				Description: "Id of the linked gcp account in New Relic",
				Required:    true,
			},
			"app_engine": {
				Type:        schema.TypeList,
				Description: "GCP app engine service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsAppEngineSchemaElem(),
			},
		},
	}
}

func cloudGcpIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Description: "the data polling interval in seconds",
			Optional:    true,
		},
	}
}

func cloudGcpIntegrationsAppEngineSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

func resourceNewrelicCloudGcpIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	cloudGcpIntegrationInputs, _ := expandCloudGcpIntegrationsInputs(d)
	gcpIntegrationspayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudGcpIntegrationInputs)
	if err != nil {
		diag.FromErr(err)
	}
	var diags diag.Diagnostics
	if len(gcpIntegrationspayload.Errors) > 0 {
		for _, err := range gcpIntegrationspayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	if len(gcpIntegrationspayload.Integrations) > 0 {
		d.SetId(strconv.Itoa(d.Get("linked_account_id").(int)))
	}
	return nil
}

func expandCloudGcpIntegrationsInputs(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	gcpCloudIntegrations := cloud.CloudGcpIntegrationsInput{}
	gcpDisableIntegrations := cloud.CloudGcpDisableIntegrationsInput{}
	var linkedAccountID int
	if lid, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = lid.(int)
	}
	if a, ok := d.GetOk("app_engine"); ok {
		gcpCloudIntegrations.GcpAppengine = expandCloudGcpAppEngineIntegrationsInputs(a.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("app_engine"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	configureInput := cloud.CloudIntegrationsInput{
		Gcp: gcpCloudIntegrations,
	}
	disableInput := cloud.CloudDisableIntegrationsInput{
		Gcp: gcpDisableIntegrations,
	}
	return configureInput, disableInput
}

func expandCloudGcpAppEngineIntegrationsInputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpAppengineIntegrationInput {
	expanded := make([]cloud.CloudGcpAppengineIntegrationInput, len(b))
	for i, appEngine := range b {
		var appEngineInput cloud.CloudGcpAppengineIntegrationInput
		in := appEngine.(map[string]interface{})
		appEngineInput.LinkedAccountId = linkedAccountID
		if a, ok := in["metrics_polling_interval"]; ok {
			appEngineInput.MetricsPollingInterval = a.(int)
		}
		expanded[i] = appEngineInput
	}
	return expanded
}

func resourceNewrelicCloudGcpIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		diag.FromErr(err)
	}
	flattenCloudGcpLinkedAccount(d, linkedAccount)
	return nil
}

func flattenCloudGcpLinkedAccount(d *schema.ResourceData, linkedAccount *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", linkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccount.ID)
	for _, i := range linkedAccount.Integrations {
		switch t := i.(type) {
		case *cloud.CloudGcpAppengineIntegration:
			_ = d.Set("app_engine", flattenCloudGcpAppEngineIntegration(t))
		}
	}
}

func flattenCloudGcpAppEngineIntegration(in *cloud.CloudGcpAppengineIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

func resourceNewrelicCloudGcpIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	configureInput, disableInput := expandCloudGcpIntegrationsInputs(d)
	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudDisableIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudDisableIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	cloudGcpIntegrationsPayload, err := client.Cloud.CloudConfigureIntegration(accountID, configureInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(cloudGcpIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudGcpIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return nil
}

func resourceNewrelicCloudGcpIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	deleteInput := expandCloudGcpDisableInputs(d)
	gcpDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(gcpDisablePayload.Errors) > 0 {
		for _, err := range gcpDisablePayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	d.SetId("")

	return nil
}

func expandCloudGcpDisableInputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudGcpDisableInput := cloud.CloudGcpDisableIntegrationsInput{}
	var linkedAccountID int
	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("app_engine"); ok {
		cloudGcpDisableInput.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	deleteInput := cloud.CloudDisableIntegrationsInput{
		Gcp: cloudGcpDisableInput,
	}
	return deleteInput
}
