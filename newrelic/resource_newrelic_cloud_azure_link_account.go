package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Description: "Application ID for Azure account",
				Required:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "Value of the client secret from Azure",
				Required:    true,
				Sensitive:   true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the linked account",
				Required:    true,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "Subscription ID for the Azure account",
				Required:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "Tenant ID for the Azure account",
				Required:    true,
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

	cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)

	if err != nil {
		diag.FromErr(err)
	}

	if len(cloudLinkAccountPayload.Errors) > 0 {
		for _, err := range cloudLinkAccountPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
	}

	if len(cloudLinkAccountPayload.LinkedAccounts) > 0 {
		d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))
	}
	return diags
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

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)

	if err != nil {
		return diag.FromErr(err)
	}

	readAzureLinkedAccount(d, linkedAccount)

	return nil
}

func readAzureLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("name", result.Name)
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
