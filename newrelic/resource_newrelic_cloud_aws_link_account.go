package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewRelicCloudAwsAccountLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAwsAccountLinkCreate,
		ReadContext:   resourceNewRelicCloudAwsAccountLinkRead,
		UpdateContext: resourceNewRelicCloudAwsAccountLinkUpdate,
		DeleteContext: resourceNewRelicCloudAwsAccountLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to link the AWS account.",
				ForceNew:    true,
			},
			"arn": {
				Type:        schema.TypeString,
				Description: "The AWS role ARN.",
				Required:    true,
				ForceNew:    true,
			},
			"metric_collection_mode": {
				Type:         schema.TypeString,
				Description:  "How metrics will be collected. Defaults to `PULL` if empty.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PULL", "PUSH"}, false),
				ForceNew:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the linked account.",
				Required:    true,
			},
		},
	}
}

func resourceNewRelicCloudAwsAccountLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkAccountInput := expandAwsCloudLinkAccountInput(d)

	var diags diag.Diagnostics

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if len(cloudLinkAccountPayload.Errors) > 0 {
			for _, err := range cloudLinkAccountPayload.Errors {
				if strings.Contains(err.Message, "The ARN you entered does not permit the correct access to your AWS account") {
					return resource.RetryableError(fmt.Errorf("%s : %s", err.Type, err.Message))
				}
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  err.Type + " " + err.Message,
				})
			}
		}

		if len(cloudLinkAccountPayload.LinkedAccounts) > 0 {
			d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	if len(diags) > 0 {
		return diags
	}

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

func resourceNewRelicCloudAwsAccountLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	readAwsLinkedAccount(d, linkedAccount)

	return nil
}

func readAwsLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("arn", result.AuthLabel)
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
}

func resourceNewRelicCloudAwsAccountLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	cloudRenameAccountPayload, err := client.Cloud.CloudRenameAccountWithContext(ctx, accountID, input)
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
		return diags
	}

	return nil
}

func resourceNewRelicCloudAwsAccountLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
