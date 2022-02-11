package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
	"log"
	"strconv"
	"strings"
)

func dataSourceNewRelicCloudAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAwsLinkAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the account in New Relic.",
			},
			"provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cloud provider of the account, e.g. aws, gcp, azure",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account",
			},
		},
	}
}

func dataSourceNewRelicAwsLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*ProviderConfig)

	client := cfg.NewClient

	log.Printf("[INFO] Reading New Relic Cloud Accounts")

	name := d.Get("name").(string)
	provider := d.Get("provider").(string)
	accountID := selectAccountID(cfg, d)

	accounts, err := client.Cloud.GetLinkedAccountsWithContext(ctx, provider)

	if err != nil {
		return diag.FromErr(err)
	}

	var account *cloud.CloudLinkedAccount

	for _, a := range *accounts {
		if strings.EqualFold(a.Name, name) && a.NrAccountId == accountID {
			account = &a
			break
		}
	}

	if account == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any account for provider '%s", name, provider))
	}

	d.SetId(strconv.Itoa(account.ID))

	return diag.FromErr(flattenCloudAccount(account, d, accountID))
}

func flattenCloudAccount(account *cloud.CloudLinkedAccount, d *schema.ResourceData, accountID int) error {
	var err error

	err = d.Set("name", account.Name)
	if err != nil {
		return err
	}

	err = d.Set("provider", account.Provider)
	if err != nil {
		return err
	}

	err = d.Set("account_id", accountID)
	if err != nil {
		return err
	}

	return nil
}
