package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicCloudGcpLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicGcpLinkAccountCreate,
		ReadContext:   resourceNewRelicGcpLinkAccountRead,
		UpdateContext: resourceNewRelicGcpLinkAccountUpdate,
		DeleteContext: resourceNewRelicGcpLinkAccountDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the aaccount",
				ForceNew:    false,
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "project id",
				ForceNew:    true,
				Required:    true,
			},
		},
	}
}

func resourceNewRelicGcpLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkGcpAccountInput := expandGcpCloudLinkAccountInput(d)
	payload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkGcpAccountInput)
	if err != nil {
		return diag.FromErr(err)
	}
	payloadReturned := &payload.LinkedAccounts[0]
	id := payloadReturned.ID
	d.SetId(string(rune(id)))
	return nil
}

func expandGcpCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	gcpAccount := cloud.CloudGcpLinkAccountInput{
		Name:      d.Get("name").(string),
		ProjectId: d.Get("project_id").(string),
	}
	input := cloud.CloudLinkCloudAccountsInput{
		Gcp: []cloud.CloudGcpLinkAccountInput{gcpAccount},
	}
	return input
}

func resourceNewRelicGcpLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicGcpLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceNewRelicGcpLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
