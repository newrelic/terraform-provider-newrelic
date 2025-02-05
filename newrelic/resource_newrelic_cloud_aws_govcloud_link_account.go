package newrelic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewRelicAwsGovCloudLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAwsGovCloudLinkAccountCreate,
		ReadContext:   resourceNewRelicAwsGovCloudLinkAccountRead,
		UpdateContext: resourceNewRelicAwsGovCloudLinkAccountUpdate,
		DeleteContext: resourceNewRelicAwsGovCloudLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the account in New Relic.",
				// since the mutation to update cloud linked accounts does not support "changing" the account ID of a linked account,
				// we shall force re-creation of the resource if the metric_collection_mode is changed after the first apply.
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the AWS GovCloud 'Linked Account' to identify in New Relic.",
				Required:    true,
			},
			"metric_collection_mode": {
				Type:         schema.TypeString,
				Description:  "The mode by which metric data is to be collected from the linked AWS GovCloud account. Use 'PUSH' for Metric Streams and 'PULL' for API Polling based metric collection respectively.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PULL", "PUSH"}, false),
				Default:      "PULL",
				// since the mutation to update cloud linked accounts does not support updating metric collection mode,
				// we shall force re-creation of the resource if the metric_collection_mode is changed after the first apply.
				ForceNew: true,
			},
			"arn": {
				Type:        schema.TypeString,
				Description: "The ARN of the identifying AWS GovCloud account.",
				Required:    true,
			},

			// NOTE: The following arguments are no longer supported, as the establishment of a connection
			// with New Relic from AWS GovCloud is no longer supported with these credentials (an ARN is needed
			// to facilitate a working connection.

			//"aws_account_id": {
			//	Type:        schema.TypeString,
			//	Description: "The ID of the AWS GovCloud account.",
			//	Required:    true,
			//},
			//"access_key_id": {
			//	Type:        schema.TypeString,
			//	Description: "The Access Key used to programmatically access the AWS GovCloud account.",
			//	Required:    true,
			//	Sensitive:   true,
			//},
			//"secret_access_key": {
			//	Type:        schema.TypeString,
			//	Description: "The Secret Access Key used to programmatically access the AWS GovCloud account.",
			//	Required:    true,
			//	Sensitive:   true,
			//},
		},
	}
}

func resourceNewRelicAwsGovCloudLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	createAwsGovCloudLinkAccountInput := expandAwsGovCloudLinkAccountInputForCreate(d)

	cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, createAwsGovCloudLinkAccountInput)
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

	if len(cloudLinkAccountPayload.LinkedAccounts) > 0 {
		d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))
	}
	return nil
}

func resourceNewRelicAwsGovCloudLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := getLinkedAccountIDFromState(d)
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	result, err := client.Cloud.GetLinkedAccount(accountID, linkedAccountID)
	if err != nil {
		return diag.FromErr(err)
	}

	expandAwsGovCloudLinkAccountInputForRead(d, result)

	return nil
}

func resourceNewRelicAwsGovCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := getLinkedAccountIDFromState(d)
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	updateAwsGovCloudLinkAccountInput := expandAwsGovCloudLinkAccountInputForUpdate(d, linkedAccountID)

	cloudUpdateAwsGovCloudAccountPayload, err := client.Cloud.CloudUpdateAccountWithContext(ctx, accountID, updateAwsGovCloudLinkAccountInput)
	if err != nil {

		return diag.FromErr(err)
	}

	if len(cloudUpdateAwsGovCloudAccountPayload.LinkedAccounts) == 0 {
		return diag.FromErr(fmt.Errorf("no linked account with 'linked_account_id': %d found", linkedAccountID))
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
