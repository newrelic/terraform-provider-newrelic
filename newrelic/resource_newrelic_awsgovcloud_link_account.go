package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicAwsGovCloudLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAwsGovCloudLinkAccountCreate,
		ReadContext:   resourceNewRelicAwsGovCloudLinkAccountRead,
		UpdateContext: resourceNewRelicAwsGovCloudLinkAccountUpdate,
		DeleteContext: resourceNewRelicAwsGovCloudLinkAccountDelete,
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Description: "access-key-id of aws account",
				Required:    true,
				ForceNew:    true,
			},
			"aws_account_id": {
				Type:        schema.TypeString,
				Description: "aws account id",
				Required:    true,
				ForceNew:    true,
			},
			"metric_collection_mode": {
				Type:        schema.TypeString,
				Description: "push or pull",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the account",
				Required:    true,
				ForceNew:    false,
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Description: "secret access key of the aws account",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceNewRelicAwsGovCloudLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkAwsGovCloudAccountInput := expandAwsGovCloudLinkAccountInput(d)
	payload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAwsGovCloudAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}
	payloadReturned := &payload.LinkedAccounts[0]
	id := payloadReturned.ID
	d.SetId(string(rune(id)))
	return nil
}

func expandAwsGovCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsGovCloud := cloud.CloudAwsGovcloudLinkAccountInput{
		AccessKeyId:          d.Get("access_key_id").(string),
		AwsAccountId:         d.Get("aws_account_id").(string),
		MetricCollectionMode: d.Get("metric_collection_mode").(cloud.CloudMetricCollectionMode),
		Name:                 d.Get("name").(string),
		SecretAccessKey:      d.Get("secret_access_key").(cloud.SecureValue),
	}
	input := cloud.CloudLinkCloudAccountsInput{
		AwsGovcloud: []cloud.CloudAwsGovcloudLinkAccountInput{awsGovCloud},
	}
	return input
}

func resourceNewRelicAwsGovCloudLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	id, _ := strconv.Atoi(d.Id())
	var input []cloud.CloudRenameAccountsInput
	renameInput := cloud.CloudRenameAccountsInput{
		Name:            d.Get("name").(string),
		LinkedAccountId: id,
	}
	input = append(input, renameInput)
	_, err := client.Cloud.CloudRenameAccount(accountID, input)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}
func resourceNewRelicAwsGovCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNewRelicAwsGovCloudLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	var input []cloud.CloudUnlinkAccountsInput
	id := d.Id()
	input[0].LinkedAccountId, _ = strconv.Atoi(id)
	_, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, input)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}
