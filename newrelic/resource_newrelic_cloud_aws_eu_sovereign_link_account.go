package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewRelicCloudAwsEuSovereignLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAwsEuSovereignLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudAwsEuSovereignLinkAccountRead,
		UpdateContext: resourceNewRelicCloudAwsEuSovereignLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudAwsEuSovereignLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the account in New Relic.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the AWS EU Sovereign account in New Relic.",
			},
			"metric_collection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PUSH",
				ForceNew:     true,
				Description:  "How metrics are collected. PULL or PUSH.",
				ValidateFunc: validation.StringInSlice([]string{"PULL", "PUSH"}, false),
			},
			"arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARN of the IAM role.",
			},
		},
	}
}

func resourceNewRelicCloudAwsEuSovereignLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	createInput := expandAwsEuSovereignLinkAccountInputForCreate(d)

	cloudLinkedAccount, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if len(cloudLinkedAccount.Errors) > 0 {
		for _, err := range cloudLinkedAccount.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	var linkedAccountID int
	if len(cloudLinkedAccount.LinkedAccounts) > 0 {
		linkedAccountID = cloudLinkedAccount.LinkedAccounts[0].ID
	}

	if linkedAccountID == 0 {
		return diag.FromErr(fmt.Errorf("failed to create AWS EU Sovereign linked account: no linked account ID returned"))
	}

	d.SetId(strconv.Itoa(linkedAccountID))

	return resourceNewRelicCloudAwsEuSovereignLinkAccountRead(ctx, d, meta)
}

func resourceNewRelicCloudAwsEuSovereignLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return diag.FromErr(flattenAwsEuSovereignLinkAccountForRead(linkedAccount, d, accountID))
}

func resourceNewRelicCloudAwsEuSovereignLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	updateInput := expandAwsEuSovereignLinkAccountInputForUpdate(d, linkedAccountID)

	_, err := client.Cloud.CloudUpdateAccountWithContext(ctx, accountID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicCloudAwsEuSovereignLinkAccountRead(ctx, d, meta)
}

func resourceNewRelicCloudAwsEuSovereignLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, _ := strconv.Atoi(d.Id())

	unlinkInput := []cloud.CloudUnlinkAccountsInput{
		{LinkedAccountId: linkedAccountID},
	}

	cloudUnlinkAccountPayload, err := client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, unlinkInput)
	if err != nil {
		return diag.FromErr(err)
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
