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

	linkAccountInput := expandCloundLinkAccountInput(d)

	payload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountId, linkAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(payload.LinkedAccounts[0].ID))

	return nil
}

func resourceNewRelicCloudLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountId := selectAccountID(providerConfig, d)

	linkedAccountId, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountId, linkedAccountId)
	
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountId := selectAccountID(providerConfig, d)

	linkedAccountId, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	renameAccountInput := []cloud.CloudRenameAccountsInput{
		{
			LinkedAccountId: linkedAccountId,
			Name: d.Get("name").(string),
		},
	}

	_, err := client.Cloud.CloudRenameAccountWithContext(ctx, accountId, renameAccountInput)
	
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicCloudLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountId := selectAccountID(providerConfig, d)

	linkedAccountId, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	unlinkAccountInput := []cloud.CloudUnlinkAccountsInput{
		{LinkedAccountId: linkedAccountId},
	}

	_, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountId, unlinkAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func expandCloundLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsAccount := cloud.CloudAwsLinkAccountInput{
		Arn:                  d.Get("arn").(string),
		MetricCollectionMode: d.Get("metricCollectionMode").(cloud.CloudMetricCollectionMode),
		Name:                 d.Get("name").(string),
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Aws: []cloud.CloudAwsLinkAccountInput{awsAccount},
	}
	return input
}

func flattenAwsLinkedAccount(d *schema.ResourceData, result cloud.CloudLinkedAccount) error {
	_ = d.Set("arn", result.AuthLabel)
	_ = d.Set("metricCollectionMode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
}
