//go:build integration || AUTH
// +build integration AUTH

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

func TestAccNewRelicAccountDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataSourceConfigByID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataSourceExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountDataSource_ByName(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataSourceExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountDataSource_MissingAttributes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataSourceConfigMissingAttributes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataSourceExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountDataSource_ConflictingAttributes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAccountDataSourceConfigConflictingAttributes(),
				ExpectError: regexp.MustCompile("exactly one of"),
			},
		},
	})
}

func testAccCheckNewRelicAccountDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		id := r.Primary.ID
		a := r.Primary.Attributes

		if id == "" {
			return fmt.Errorf("expected to get an account from New Relic")
		}

		if a["name"] == "" {
			return fmt.Errorf("expected to get an account name from New Relic")
		}

		if a["account_id"] == "" {
			return fmt.Errorf("expected to get an account ID from New Relic")
		}

		// Verify the account exists using the customeradministration package
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		// Convert account ID from string to int
		accountID, err := strconv.Atoi(a["account_id"])
		if err != nil {
			return fmt.Errorf("failed to convert account ID to integer: %v", err)
		}

		// Get organization ID
		organization, err := client.Organization.GetOrganization()
		if err != nil {
			return fmt.Errorf("failed to fetch organization information: %v", err)
		}

		// Fetch account using customeradministration package
		ctx := context.Background()
		getAccountsResponse, err := client.CustomerAdministration.GetAccounts(
			"",
			customeradministration.OrganizationAccountFilterInput{
				OrganizationId: customeradministration.OrganizationAccountOrganizationIdFilterInput{
					Eq: organization.ID,
				},
				ID: customeradministration.OrganizationAccountIdFilterInput{
					Eq: accountID,
				},
			},
			[]customeradministration.OrganizationAccountSortInput{},
		)

		if err != nil {
			return fmt.Errorf("failed to fetch account details: %v", err)
		}

		accounts := getAccountsResponse.Items
		if len(accounts) == 0 {
			return fmt.Errorf("account not found: %d", accountID)
		}

		if len(accounts) != 1 {
			return fmt.Errorf("expected 1 account, found %d", len(accounts))
		}

		if accounts[0].ID != accountID {
			return fmt.Errorf("expected account ID %d, got %d", accountID, accounts[0].ID)
		}

		return nil
	}
}

func testAccNewRelicAccountDataSourceConfigByID() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	account_id = "%d"
}
`, testAccountID)
}

func testAccNewRelicAccountDataSourceConfigByName() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	name = "%s"
}
`, testAccountName)
}

func testAccNewRelicAccountDataSourceConfigMissingAttributes() string {
	return `data "newrelic_account" "acc" {}`
}

func testAccNewRelicAccountDataSourceConfigConflictingAttributes() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	name = "%s"
	account_id = "%d"
}
`, testAccountName, testAccountID)
}
