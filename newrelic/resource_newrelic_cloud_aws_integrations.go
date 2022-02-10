package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNewRelicCloudAwsIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAwsIntegrationsCreate,
		ReadContext:   resourceNewRelicCloudAwsIntegrationsRead,
		UpdateContext: resourceNewRelicCloudAwsIntegrationsUpdate,
		DeleteContext: resourceNewRelicCloudAwsIntegrationsDelete,
		Schema:        map[string]*schema.Schema{},
	}
}

func resourceNewRelicCloudAwsIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudAwsIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudAwsIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudAwsIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
