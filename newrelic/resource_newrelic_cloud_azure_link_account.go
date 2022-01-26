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
		CreateContext: resourceNewRelicAzureLinkAccountCreate,
		ReadContext:   resourceNewRelicAzureLinkAccountRead,
		UpdateContext: resourceNewRelicAzureLinkAccountUpdate,
		DeleteContext: resourceNewRelicAzureLinkAccountDelete,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Description: "application id",
				Required:    true,
				ForceNew:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "client secret",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the account",
				Required:    true,
				ForceNew:    false,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "subscription id",
				Required:    true,
				ForceNew:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "tenant id",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceNewRelicAzureLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkAZureAccountInput := expandAzureCloudLinkAccountInput(d)
	payload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAZureAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}
	payloadReturned := &payload.LinkedAccounts[0]
	id := payloadReturned.ID
	d.SetId(string(rune(id)))
	return nil
}

func expandAzureCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	azureAccout := cloud.CloudAzureLinkAccountInput{
		ApplicationID:  d.Get("application_id").(string),
		ClientSecret:   d.Get("client_secret").(cloud.SecureValue),
		Name:           d.Get("name").(string),
		SubscriptionId: d.Get("subscription_id").(string),
		TenantId:       d.Get("tenant_id").(string),
	}
	input := cloud.CloudLinkCloudAccountsInput{
		Azure: []cloud.CloudAzureLinkAccountInput{azureAccout},
	}
	return input
}

func resourceNewRelicAzureLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicAzureLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceNewRelicAzureLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
