package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accountmanagement"
)

func resourceNewRelicWorkloadAccountManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAccountCreate,
		ReadContext:   resourceNewRelicAccountRead,
		UpdateContext: resourceNewRelicAccountUpdate,
		DeleteContext: resourceNewRelicAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the account to be created",
				Required:    true,
			},
			"region": {
				Type:         schema.TypeString,
				Description:  "A description of what this parsing rule represents.",
				ValidateFunc: validation.StringInSlice([]string{"us01", "eu01"}, false),
				Required:     true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
		},
	}

}

func resourceNewRelicAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		account, err := getCreatedAccountByID(ctx, client, d.Id())
		//		fmt.Println("read", account.ID, err)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if account == nil {
			return resource.RetryableError(fmt.Errorf("account not found"))
		}
		d.Set("region", account.RegionCode)
		d.Set("name", account.Name)

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}
	return nil
}

func resourceNewRelicAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createAccountInput := accountmanagement.AccountManagementCreateInput{
		Name:       d.Get("name").(string),
		RegionCode: d.Get("region").(string),
	}
	created, err := client.AccountManagement.AccountManagementCreateAccount(createAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: Account not created. Please check the input details")
	}
	accountId := created.ManagedAccount.ID

	d.SetId(strconv.Itoa(accountId))
	return resourceNewRelicAccountRead(ctx, d, meta)
}
func resourceNewRelicAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	updateAccountInput := accountmanagement.AccountManagementUpdateInput{
		Name: d.Get("name").(string),
		ID:   accountId,
	}
	updated, err := client.AccountManagement.AccountManagementUpdateAccount(updateAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		return diag.Errorf("err: Account not Updated. Please check the input details")
	}

	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Account cannot be deleted via Terraform. Please contact support",
	})
	return diags
}

func getCreatedAccountByID(ctx context.Context, client *newrelic.NewRelic, ruleID string) (*accountmanagement.AccountManagementManagedAccount, error) {

	accountId, err := strconv.Atoi(ruleID)
	if err != nil {
		return nil, err
	}
	accounts, err := client.AccountManagement.GetManagedAccounts()
	if err != nil && accounts == nil {
		return nil, err
	}
	for _, account := range *accounts {
		if account.ID == accountId {
			return &account, nil
		}
	}
	return nil, nil

}
