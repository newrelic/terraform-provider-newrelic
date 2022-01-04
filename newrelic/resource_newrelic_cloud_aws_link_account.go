package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicCloudAwsLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudLinkAccountRead,
		UpdateContext: resourceNewRelicCloudLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS role ARN",
			},
			"metricCollectionMode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "How metrics will be collected.",
				ValidateFunc: validation.StringInSlice([]string{"PULL", "PUSH"}, true),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The linked account name",
			},
		},
	}
}

func resourceNewRelicCloudLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountId := selectAccountID(providerConfig, d)

	linkAccountInput, err := expandCloundLinkAccountInput(d)

	if err != nil {
		return diag.FromErr(err)
	}

	payload, reqErr := client.Cloud.CloudLinkAccountWithContext(ctx, accountId, linkAccountInput)
	if reqErr != nil {
		return diag.FromErr(reqErr)
	}

	d.SetId(strconv.Itoa(payload.LinkedAccounts[0].ID))

	return nil
}

func resourceNewRelicCloudLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandCloundLinkAccountInput(d *schema.ResourceData) (cloud.CloudLinkCloudAccountsInput, error) {
	awsAccount := cloud.CloudAwsLinkAccountInput{
		Arn:                  d.Get("arn").(string),
		MetricCollectionMode: d.Get("metricCollectionMode").(cloud.CloudMetricCollectionMode),
		Name:                 d.Get("name").(string),
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Aws: []cloud.CloudAwsLinkAccountInput{awsAccount},
	}
	return input, nil
}
