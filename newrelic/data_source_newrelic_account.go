package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/accounts"
)

func dataSourceNewRelicAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAccountRead,
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(accounts.RegionScopeTypes.IN_REGION),
				Description:  `The scope of the account in New Relic.  Valid values are "global" and "in_region".  Defaults to "in_region".`,
				ValidateFunc: validation.StringInSlice([]string{"global", "in_region"}, true),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the account in New Relic.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the account in New Relic.",
			},
		},
	}
}

func dataSourceNewRelicAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic accounts")

	scope := accounts.RegionScope(strings.ToUpper(d.Get("scope").(string)))

	id, idOk := d.GetOk("account_id")
	name, nameOk := d.GetOk("name")

	params := accounts.ListAccountsParams{
		Scope: &scope,
	}

	accts, err := client.Accounts.ListAccountsWithContext(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	var account *accounts.AccountOutline

	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf(`one of "name" or "account_id" is required to locate a New Relic account`))
	}

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf(`exactly one of "name" or "account_id" is required to locate a New Relic account`))
	}

	if nameOk {
		for _, a := range accts {
			if a.Name == name.(string) {
				account = &a
				break
			}
		}

		if account == nil {
			return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic accounts", name))
		}
	}

	if idOk {
		for _, a := range accts {
			if a.ID == id.(int) {
				account = &a
				break
			}
		}

		if account == nil {
			return diag.FromErr(fmt.Errorf("the id '%d' does not match any New Relic accounts", id))
		}
	}

	return diag.FromErr(flattenAccountData(account, d))
}

func flattenAccountData(a *accounts.AccountOutline, d *schema.ResourceData) error {
	d.SetId(strconv.Itoa(a.ID))
	var err error

	err = d.Set("name", a.Name)
	if err != nil {
		return err
	}

	err = d.Set("account_id", a.ID)
	if err != nil {
		return err
	}

	return nil
}
