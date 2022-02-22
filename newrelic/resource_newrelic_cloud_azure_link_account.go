package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicCloudAzureLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAzureLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudAzureLinkAccountRead,
		UpdateContext: resourceNewRelicCloudAzureLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudAzureLinkAccountDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to link the Azure account.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Description: "application id for Azure account",
				Required:    true,
				ForceNew:    true,
			},
			"client_secret_id": {
				Type:        schema.TypeString,
				Description: "Value of the client secret from Azure",
				Required:    true,
				ForceNew:    true,
			},

			"name": {
				Type:        schema.TypeString,
				Description: "name of the linked account",
				Required:    true,
				ForceNew:    false,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "subscription id for the Azure account",
				Required:    true,
				ForceNew:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "tenant id for the Azure account",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceNewRelicCloudAzureLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkAccountInput := expandAzureCloudLinkAccountInput(d)

	var diags diag.Diagnostics

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)
		if err != nil {
			return resource.NonRetryableError(err)

		}

		if len(cloudLinkAccountPayload.Errors) > 0 {
			for _, err := range cloudLinkAccountPayload.Errors {
				if strings.Contains(err.Message, "We encountered an error") {
					return resource.RetryableError(fmt.Errorf("%s : %s", err.Type, err.Message))

				}

				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  err.Type + " " + err.Message,
				})
			}

		}

		d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))

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

func expandAzureCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {

	azureAccount := cloud.CloudAzureLinkAccountInput{}

	if applicationID, ok := d.GetOk("application_id"); ok {
		azureAccount.ApplicationID = applicationID.(string)
	}

	if clientSecretID, ok := d.GetOk("client_secret_id"); ok {
		azureAccount.ClientSecret = cloud.SecureValue(clientSecretID.(string))
	}

	if name, ok := d.GetOk("name"); ok {
		azureAccount.Name = name.(string)
	}

	if subscriptionID, ok := d.GetOk("subscription_id"); ok {
		azureAccount.SubscriptionId = subscriptionID.(string)
	}

	if tenantID, ok := d.GetOk("tenant_id"); ok {
		azureAccount.TenantId = tenantID.(string)
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Azure: []cloud.CloudAzureLinkAccountInput{azureAccount},
	}

	return input
}
func resourceNewRelicCloudAzureLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	readAzureLinkedAccount(d, linkedAccount)

	return nil
}

////////

func readAzureLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("application_id", result.ID)
	_ = d.Set("client_secret_id", result.ID)
	_ = d.Set("name", result.Name)
	_ = d.Set("subscription_id", result.ID)
	_ = d.Set("tenant_id", result.ID)
}

func resourceNewRelicCloudAzureLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

		return diags
	}

	return nil
}

func resourceNewRelicCloudAzureLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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
