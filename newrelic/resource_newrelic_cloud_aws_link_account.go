package newrelic

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicCloudAwsAccountLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAwsAccountLinkCreate,
		ReadContext:   resourceNewRelicAwsAccountLinkRead,
		UpdateContext: resourceNewRelicAwsAccountLinkUpdate,
		DeleteContext: resourceNewRelicAwsAccountLinkDelete,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:        schema.TypeString,
				Description: "The AWS role arn",
				Required:    true,
			},
			"metric_collection_mode": {
				Type:         schema.TypeString,
				Description:  "How metrics will be collected",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PULL", "PUSH"}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the linked account",
				Required:    true,
			},
		},
	}
}

func resourceNewRelicAwsAccountLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkAccountInput := expandAwsCloudLinkAccountInput(d)

	payload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(payload.LinkedAccounts[0].ID))

	return nil
}

func expandAwsCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
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

func resourceNewRelicAwsAccountLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccount(accountID, linkedAccountID)

	if err != nil {
		return diag.FromErr(err)
	}

	readAwsLinkedAccount(d, linkedAccount)

	return nil
}

func readAwsLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("arn", result.AuthLabel)
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
}

func resourceNewRelicAwsAccountLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	id, _ := strconv.Atoi(d.Id())
	input := []cloud.CloudRenameAccountsInput{
		{
			Name:            d.Get("name").(string),
			LinkedAccountId: id,
		},
	}
	_, err := client.Cloud.CloudRenameAccount(accountID, input)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}

func resourceNewRelicAwsAccountLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	unlinkAccountInput := []cloud.CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
		},
	}
	_, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, unlinkAccountInput)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
