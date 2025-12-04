//go:build AUTH
// +build AUTH

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

func TestAccNewRelicAccountManagement_Basic(t *testing.T) {
	resourceName := "newrelic_account_management.foo"
	rName := acctest.RandString(7)
	rNameUpdated := acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAccountCreateConfig("Test " + rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "Test "+rName),
					resource.TestCheckResourceAttr(resourceName, "region", "us01"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAccountUpdateConfig("Updated " + rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "Updated "+rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "region", "us01"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAccountManagement_Import(t *testing.T) {
	t.Skipf("Skipping import test: Account import is resulting in import/destroy deadlock inconsistencies")

	resourceName := "newrelic_account_management.foo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Import
			{
				ImportState:  true,
				Config:       testAccNewRelicAccountImportConfig(),
				ResourceName: resourceName,
				// do not change this
				ImportStateId:      "7400957",
				ImportStateCheck:   testAccCheckNewRelicAccountImportCheck(resourceName),
				ImportStatePersist: true,
			},
		},
	})
}
func TestAccNewRelicAccountManagementInvalidRegion(t *testing.T) {
	rName := acctest.RandString(7)
	expectedErrorMsg := regexp.MustCompile(`expected region to be one of \[us01 eu01\], got abcd01`)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicAccountCreateInvalidRegionConfig("Test " + rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}
func TestAccNewRelicAccountManagementInCorrectRegion(t *testing.T) {
	rName := acctest.RandString(7)
	expectedErrorMsg := regexp.MustCompile(`An error occurred resolving this field|cannot create account -- no configured parent account`)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicAccountCreateConfigInCorrectRegion("Test " + rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccNewRelicAccountImportConfig() string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo" {
  name   = ""
  region = "us01"
}
`)
}

func testAccNewRelicAccountCreateConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[1]s"
	region= "us01"
}
`, name)
}

func testAccNewRelicAccountUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[1]s"
	region= "us01"
}
`, name)
}

func testAccCheckNewRelicAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		// Convert account ID from string to int
		accountID, err := strconv.Atoi(rs.Primary.ID)
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

		// Suppress unused variable warning
		_ = ctx

		return nil
	}
}

func testAccNewRelicAccountCreateInvalidRegionConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[1]s"
	region= "abcd01"
}
`, name)
}

func testAccNewRelicAccountCreateConfigInCorrectRegion(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[1]s"
	region= "eu01"
}
`, name)
}

func testAccCheckNewRelicAccountImportCheck(resourceName string) resource.ImportStateCheckFunc {
	return func(state []*terraform.InstanceState) error {
		expectedRegionCode := "us01"
		region := state[0].Attributes["region"]
		if region != expectedRegionCode {
			return fmt.Errorf(
				"%s: Attribute '%s' expected %#v got nil",
				resourceName,
				"region.#",
				expectedRegionCode,
			)
		}

		return nil
	}
}
