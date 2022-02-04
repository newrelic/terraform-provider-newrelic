package newrelic

import (
	"context"
	"strconv"
	"strings"

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

	linkAccountInput := expandAwsGovCloudLinkAccountInput(d)

	cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	if len(cloudLinkAccountPayload.Errors) > 0 {
		for _, err := range cloudLinkAccountPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))
	return nil
}

func expandAwsGovCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsGovCloud := cloud.CloudAwsGovcloudLinkAccountInput{}
	if accessKeyID, ok := d.GetOk("access_key_id"); ok {
		awsGovCloud.AccessKeyId = accessKeyID.(string)
	}
	if awsAccountID, ok := d.GetOk("aws_account_id"); ok {
		awsGovCloud.AwsAccountId = awsAccountID.(string)
	}
	if m, ok := d.GetOk("metric_collection_mode"); ok {
		awsGovCloud.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	}
	if name, ok := d.GetOk("name"); ok {
		awsGovCloud.Name = name.(string)
	}
	if secretKeyID, ok := d.GetOk("secret_access_key"); ok {
		awsGovCloud.SecretAccessKey = secretKeyID.(cloud.SecureValue)
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
	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccountPayload, err := client.Cloud.GetLinkedAccount(accountID, linkedAccountID)

	if err != nil {
		return diag.FromErr(err)
	}
	readAwsGovCloudLinkAccount(d, linkedAccountPayload)
	return nil
}

func readAwsGovCloudLinkAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
	//
	_ = d.Set("aws_account_id", result.NrAccountId)
	_ = d.Set("access_key_id", result.ID)
	_ = d.Set("secret_access_key", result.ID)
}

func resourceNewRelicAwsGovCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	cloudRenameAccountPayload, err := client.Cloud.CloudRenameAccount(accountID, input)
	if err != nil {
		diag.FromErr(err)
	}
	var diags diag.Diagnostics

	if len(cloudRenameAccountPayload.Errors) > 0 {
		for _, err := range cloudRenameAccountPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})

		}
	}

	return nil
}

func resourceNewRelicAwsGovCloudLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	cloudUnlinkAccountPayload, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, unlinkAccountInput)
	if err != nil {
		diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if len(cloudUnlinkAccountPayload.Errors) > 0 {
		for _, err := range cloudUnlinkAccountPayload.Errors {
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
