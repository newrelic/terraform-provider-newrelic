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
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Description: "access-key-id of awsGovcloud account",
				Required:    true,
				Sensitive:   true,
			},
			"aws_account_id": {
				Type:        schema.TypeString,
				Description: "awsGovcloud account id",
				Required:    true,
			},
			"metric_collection_mode": {
				Type:        schema.TypeString,
				Description: "push or pull",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the account",
				Required:    true,
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Description: "secret access key of the awsGovcloud account",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceNewRelicAwsGovCloudLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	provider := d.Get("cloud_provider").(string)
	name := d.Get("name").(string)

	linkAccountInput := expandAwsGovCloudLinkAccountInput(d)

	cloudLinkAccountPayload, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkAccountInput)

	var diags diag.Diagnostics

	if err != nil {
		return diag.FromErr(err)
	}

	if cloudLinkAccountPayload == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] cloudLinkAccountPayload was nil"))
	}

	if len(cloudLinkAccountPayload.Errors) > 0 {
		for _, err := range cloudLinkAccountPayload.Errors {
			if string(err.Type) == "ERR_INVALID_DATA" && err.LinkedAccountId == 0 {
				accounts, getLinkedAccountsErr := client.Cloud.GetLinkedAccounts(provider)
				if getLinkedAccountsErr != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  err.Type + " " + err.Message + " " + getLinkedAccountsErr.Error(),
					})
				}

				var account *cloud.CloudLinkedAccount

				for _, a := range *accounts {
					if a.NrAccountId == accountID && strings.EqualFold(a.Name, name) {
						account = &a
						break
					}
				}

				if account == nil {
					return diag.FromErr(fmt.Errorf("the name '%s' does not match any account for provider '%s", name, provider))
				}

				d.SetId(strconv.Itoa(account.ID))
			} else if err.LinkedAccountId != 0 {
				d.SetId(strconv.Itoa(err.LinkedAccountId))
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  err.Type + " " + err.Message,
				})
			}
		}
	}

	if len(cloudLinkAccountPayload.LinkedAccounts) > 0 {
		d.SetId(strconv.Itoa(cloudLinkAccountPayload.LinkedAccounts[0].ID))
	}

	return diags
}

// Extracting the AWSGovCloud account  credentials from Schema using expandAzureCloudLinkAccountInput
func expandAwsGovCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsGovCloud := cloud.CloudAwsGovCloudLinkAccountInput{}
	if accessKeyID, ok := d.GetOk("access_key_id"); ok {
		awsGovCloud.AccessKeyId = accessKeyID.(string)
	}
	if awsAccountID, ok := d.GetOk("aws_account_id"); ok {
		awsGovCloud.AwsAccountId = awsAccountID.(string)
	}
	if m, ok := d.GetOk("metric_collection_mode"); ok {
		awsGovCloud.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	}
	if name, ok := d.GetOk("name"); ok {
		awsGovCloud.Name = name.(string)
	}
	if secretKeyID, ok := d.GetOk("secret_access_key"); ok {
		awsGovCloud.SecretAccessKey = secretKeyID.(cloud.SecureValue)
	}
	input := cloud.CloudLinkCloudAccountsInput{
		AwsGovcloud: []cloud.CloudAwsGovCloudLinkAccountInput{awsGovCloud},
	}
	return input
}
func resourceNewRelicAwsGovCloudLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccountPayload, err := client.Cloud.GetLinkedAccount(accountID, linkedAccountID)

	if err != nil {
		return diag.FromErr(err)
	}
	readAwsGovCloudLinkAccount(d, linkedAccountPayload)
	return nil
}

func readAwsGovCloudLinkAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
	_ = d.Set("aws_account_id", result.ID)
	_ = d.Set("account_id", result.NrAccountId)
}

func resourceNewRelicAwsGovCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	//Setting up the linked account id to null after destroying the resource.
	d.SetId("")

	return nil
}
