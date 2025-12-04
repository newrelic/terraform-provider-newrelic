package newrelic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accounts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

func dataSourceNewRelicAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAccountRead,
		Schema: map[string]*schema.Schema{
			NewRelicAccountManagementSchemaName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the account in New Relic.",
				ConflictsWith: []string{"account_id"},
			},
			NewRelicAccountManagementSchemaRegion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region code of the account (e.g., us01, eu01).",
			},
			"account_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The ID of the account in New Relic.",
				ConflictsWith: []string{NewRelicAccountManagementSchemaName},
			},
			// deprecated and no longer used by the data source; just adding this here for feature parity
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(accounts.RegionScopeTypes.IN_REGION),
				Description:  `The scope of the account in New Relic.  Valid values are "global" and "in_region".  Defaults to "in_region".`,
				ValidateFunc: validation.StringInSlice([]string{"global", "in_region"}, true),
			},
		},
	}
}

func dataSourceNewRelicAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var err error

	organization, err := client.Organization.GetOrganization()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch organization information: %v", err))
	}

	organizationID := organization.ID

	accountID, accountIDOk := d.GetOk("account_id")
	name, nameOk := d.GetOk(NewRelicAccountManagementSchemaName)

	if !accountIDOk && !nameOk {
		// Default to the provider's AccountID if no lookup attributes are provided.
		accountID, accountIDOk = selectAccountID(providerConfig, d), true
	}

	var account *customeradministration.OrganizationAccount

	// Build filter input based on whether we're searching by ID or name
	filterInput := customeradministration.OrganizationAccountFilterInput{
		OrganizationId: customeradministration.OrganizationAccountOrganizationIdFilterInput{
			Eq: organizationID,
		},
	}

	if accountIDOk {
		filterInput.ID = customeradministration.OrganizationAccountIdFilterInput{
			Eq: accountID.(int),
		}
	}

	if nameOk {
		filterInput.Name = customeradministration.OrganizationAccountNameFilterInput{
			Contains: name.(string),
		}
		// Note: Using a "Contains" filter may return multiple results, as it performs a partial match.
		// To ensure accuracy, additional steps are required to identify the exact match from the results.
	}

	getAccountsResponse, err := client.CustomerAdministration.GetAccounts(
		"",
		filterInput,
		[]customeradministration.OrganizationAccountSortInput{},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	accounts := getAccountsResponse.Items

	if len(accounts) == 0 {
		return diag.FromErr(fmt.Errorf("no accounts found matching the criteria"))
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

		if matchedAccount == nil || matchedAccount.Name == "" {
			return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic accounts", name))
		}
		account = matchedAccount
	} else {
		// If searching by ID, we should have exactly one result
		if len(accounts) != 1 {
			return diag.FromErr(fmt.Errorf("expected 1 account, found %d", len(accounts)))
		}
		account = &accounts[0]
	}

	d.SetId(strconv.Itoa(account.ID))

	err = d.Set(NewRelicAccountManagementSchemaName, account.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("account_id", account.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(NewRelicAccountManagementSchemaRegion, account.RegionCode)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
