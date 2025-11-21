package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

func dataSourceNewRelicAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAccountRead,
		Schema: map[string]*schema.Schema{
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
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region code of the account (e.g., us01, eu01).",
			},
		},
	}
}

func dataSourceNewRelicAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic accounts")

	var diags diag.Diagnostics

	// Get organization ID
	organization, err := client.Organization.GetOrganization()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to fetch organization information: %v", err),
		})
		return diags
	}

	organizationID := organization.ID

	id, idOk := d.GetOk("account_id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		// Default to the provider's AccountID if no lookup attributes are provided.
		id, idOk = selectAccountID(providerConfig, d), true
	}

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf(`exactly one of "name" or "account_id" is required to locate a New Relic account`))
	}

	var account *customeradministration.OrganizationAccount

	// Build filter input based on whether we're searching by ID or name
	filterInput := customeradministration.OrganizationAccountFilterInput{
		OrganizationId: customeradministration.OrganizationAccountOrganizationIdFilterInput{
			Eq: organizationID,
		},
	}

	if idOk {
		// If searching by ID, add ID filter
		filterInput.ID = customeradministration.OrganizationAccountIdFilterInput{
			Eq: id.(int),
		}
	}

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		getAccountsResponse, err := client.CustomerAdministration.GetAccounts(
			"",
			filterInput,
			[]customeradministration.OrganizationAccountSortInput{},
		)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		accounts := getAccountsResponse.Items

		if len(accounts) == 0 {
			return resource.RetryableError(fmt.Errorf("no accounts found matching the criteria"))
		}

		// If searching by name, filter client-side
		if nameOk {
			var matchedAccount *customeradministration.OrganizationAccount
			for _, a := range accounts {
				if a.Name == name.(string) {
					matchedAccount = &a
					break
				}
			}

			if matchedAccount == nil {
				return resource.NonRetryableError(fmt.Errorf("the name '%s' does not match any New Relic accounts", name))
			}
			account = matchedAccount
		} else {
			// If searching by ID, we should have exactly one result
			if len(accounts) != 1 {
				return resource.RetryableError(fmt.Errorf("expected 1 account, found %d", len(accounts)))
			}
			account = &accounts[0]
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diag.FromErr(flattenAccountData(account, d))
}

func flattenAccountData(a *customeradministration.OrganizationAccount, d *schema.ResourceData) error {
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

	err = d.Set("region", a.RegionCode)
	if err != nil {
		return err
	}

	return nil
}
