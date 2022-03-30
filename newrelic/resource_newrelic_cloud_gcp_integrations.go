package newrelic

import (
	"context"
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
	cloudGcpIntegrationInputs := expandCloudGcpInputs(d)
	_, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudGcpIntegrationInputs)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}

func expandCloudGcpInputs(d *schema.ResourceData) cloud.CloudIntegrationsInput {
	gcpCloudIntegrations := cloud.CloudGcpIntegrationsInput{}
	var _ int
	if AccountID, ok := d.GetOk("linked_account_id"); ok {
		_ = AccountID.(int)
	}
	input := cloud.CloudIntegrationsInput{
		Gcp: gcpCloudIntegrations,
	}
	return input
}

func expandCloudGcpAppEngineInputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpAppengineIntegration {
	return nil
}

func resourceNewrelicCloudGcpIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewrelicCloudGcpIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewrelicCloudGcpIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
