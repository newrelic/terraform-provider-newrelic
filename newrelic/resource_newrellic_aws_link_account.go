package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicAwsAccountLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAwsAccountLinkCreate,
		ReadContext:   resourceNewRelicAwsAccountLinkRead,
		UpdateContext: resourceNewRelicAwsAccountLinkUpdate,
		DeleteContext: resourceNewRelicAwsAccountLinkDelete,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:        schema.TypeString,
				Description: "aws iam role arn",
				Required:    true,
				ForceNew:    true,
			},
			"metric_collection_mode": {
				Type:        schema.TypeString,
				Description: "push or pull metric collection mode",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the cloud link account",
				Required:    true,
				ForceNew:    false,
			},
		},
	}
}

func resourceNewRelicAwsAccountLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkAccount := expandAwsCloudLinkAccountInput(d)
	payload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccount)
	if err != nil {
		return diag.FromErr(err)
	}
	payloadReturned := &payload.LinkedAccounts[0]
	id := payloadReturned.ID
	d.SetId(string(rune(id)))
	return nil
}

func expandAwsCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsAccount := cloud.CloudAwsLinkAccountInput{
		Arn:                  d.Get("arn").(string),
		MetricCollectionMode: d.Get("metric_collection_mode").(cloud.CloudMetricCollectionMode),
		Name:                 d.Get("name").(string),
	}
	input := cloud.CloudLinkCloudAccountsInput{
		Aws: []cloud.CloudAwsLinkAccountInput{awsAccount},
	}
	return input
}

func resourceNewRelicAwsAccountLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicAwsAccountLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceNewRelicAwsAccountLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
