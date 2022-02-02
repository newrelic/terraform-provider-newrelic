package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
	"strconv"
)

func resourceNewRelicCloudGcpLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudGcpLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudGcpLinkAccountRead,
		UpdateContext: resourceNewRelicCloudGcpLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudGcpLinkAccountDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the linked account",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "project id of the Gcp account",
				Required:    true,
			},
		},
	}
}

func resourceNewRelicCloudGcpLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkAccountInput := expandGcpCloudLinkAccountInput(d)
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

func expandGcpCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	gcpAccount := cloud.CloudGcpLinkAccountInput{}
	if name, ok := d.GetOk("arn"); ok {
		gcpAccount.Name = name.(string)
	}
	if projectID, ok := d.GetOk("project_id"); ok {
		gcpAccount.ProjectId = projectID.(string)
	}
	input := cloud.CloudLinkCloudAccountsInput{
		Gcp: []cloud.CloudGcpLinkAccountInput{gcpAccount},
	}
	return input
}

func resourceNewRelicCloudGcpLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	readGcpLinkedAccount(d, linkedAccount)
	return nil
}

func readGcpLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("name", result.Name)
	_ = d.Set("project_id", result.ID)
}

func resourceNewRelicCloudGcpLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceNewRelicCloudGcpLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID, convErr := strconv.Atoi(d.Id())
	if convErr != nil {
		diag.FromErr(convErr)
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
