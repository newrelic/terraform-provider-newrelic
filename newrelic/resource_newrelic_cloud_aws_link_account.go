package newrelic

import (
	"context"
	"strconv"
	"strings"

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
			"metric_collection_mode": {
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

	linkAccountInput := expandCloudLinkAccountInput(d)

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

	readAwsLinkedAccount(d, linkedAccount)

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
			Name:            d.Get("name").(string),
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

func expandCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsAccount := cloud.CloudAwsLinkAccountInput{}

	if arn, ok := d.GetOk("arn"); ok {
		awsAccount.Arn = arn.(string)
	}

	if m, ok := d.GetOk("metric_collection_mode"); ok {
		awsAccount.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	}

	if name, ok := d.GetOk("name"); ok {
		awsAccount.Name = name.(string)
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Aws: []cloud.CloudAwsLinkAccountInput{awsAccount},
	}
	return input
}

func readAwsLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("arn", result.AuthLabel)
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
}
