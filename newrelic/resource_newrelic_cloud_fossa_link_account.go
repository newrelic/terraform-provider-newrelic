package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewRelicCloudFossaLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFossaLinkAccountCreate,
		ReadContext:   resourceNewRelicFossaLinkAccountRead,
		UpdateContext: resourceNewRelicFossaLinkAccountUpdate,
		DeleteContext: resourceNewRelicFossaLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Description: "The Fossa account application api key(bearer token)",
				Required:    true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Description: "The Fossa account identifier",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The linked account name",
				Required:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Disable the linked account.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceNewRelicFossaLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkAccountInput := expandFossaCloudLinkAccountInput(d)
	var diags diag.Diagnostics

	cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)

	if err != nil {
		return diag.FromErr(err)
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

func expandFossaCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	fossaAccount := cloud.CloudFossaLinkAccountInput{}

	if apiKey, ok := d.GetOk("api_key"); ok {
		fossaAccount.APIKey = cloud.SecureValue(apiKey.(string))
	}

	if externerId, ok := d.GetOk("external_id"); ok {
		fossaAccount.ExternalId = externerId.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		fossaAccount.Name = name.(string)
	}
	input := cloud.CloudLinkCloudAccountsInput{
		Fossa: []cloud.CloudFossaLinkAccountInput{fossaAccount},
	}

	return input
}

func resourceNewRelicFossaLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	readFossaLinkedAccount(d, linkedAccount)

	return nil
}
func readFossaLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("external_id", result.ExternalId)
	_ = d.Set("name", result.Name)
	_ = d.Set("disabled", result.Disabled)
}

func resourceNewRelicFossaLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID, _ := strconv.Atoi(d.Id())

	input := cloud.CloudUpdateCloudAccountsInput{
		Fossa: []cloud.CloudFossaUpdateAccountInput{
			{
				APIKey:          cloud.SecureValue(d.Get("api_key").(string)),
				ExternalId:      d.Get("external_id").(string),
				Name:            d.Get("name").(string),
				Disabled:        d.Get("disabled").(bool),
				LinkedAccountId: linkedAccountID,
			},
		},
	}

	cloudUpdateAccountPayload, err := client.Cloud.CloudUpdateAccountWithContext(ctx, accountID, input)

	if err != nil {

		return diag.FromErr(err)
	}
	if len(cloudUpdateAccountPayload.LinkedAccounts) == 0 {
		return diag.FromErr(fmt.Errorf("no linked account with 'linked_account_id': %d found", linkedAccountID))
	}
	return nil
}

func resourceNewRelicFossaLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
