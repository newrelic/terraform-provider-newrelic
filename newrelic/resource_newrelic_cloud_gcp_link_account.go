package newrelic

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicCloudGcpLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudGcpLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudGcpLinkAccountRead,
		UpdateContext: resourceNewRelicCloudGcpLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudGcpLinkAccountDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "accountID of newrelic account",
				Required:    true,
			},
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
		},
	}
}

func resourceNewRelicCloudGcpLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkAccountInput := expandGcpCloudLinkAccountInput(d)

	var diags diag.Diagnostics

	//cloudLinkAccountWithContext func which links Gcp account with Newrelic
	//which returns payload and error
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

	//Storing the linked account id using setId func after creating the resource.
	d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))

	if len(diags) > 0 {
		return diags
	}

	return nil
}

//expand function to extract inputs from the schema.
//Here it takes ResourceData as input and returns cloudLinkCloudAccountsInput.
func expandGcpCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {

	gcpAccount := cloud.CloudGcpLinkAccountInput{}

	if name, ok := d.GetOk("name"); ok {
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

//readGcpLinkedAccount function to store name and ExternalId.
//Using set func to store the output values.
func readGcpLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {

	_ = d.Set("name", result.Name)
	_ = d.Set("project_id", result.ExternalId)
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

	//CloudRenameAccount to rename the name of linkedAccount
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

	//CloudUnlinkAccountWithContext func to unlink the GCP account with Newrelic
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
	//Setting up the linked account id to null after destroying the resource.
	d.SetId("")

	return nil

}
