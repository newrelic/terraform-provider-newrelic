//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAzureLinkAccount_Basic(t *testing.T) {
	testAzureLinkAccountName := fmt.Sprintf("tf_cloud_link_account_test_azure_%s", acctest.RandString(5))
	resourceName := "newrelic_cloud_azure_link_account.foo"

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testAzureApplicationID := os.Getenv("INTEGRATION_TESTING_AZURE_APPLICATION_ID")
	if testAzureApplicationID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_APPLICATION_ID must be set for acceptance test")
	}

	testAzureClientSecretID := os.Getenv("INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID")
	if testAzureClientSecretID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID must be set for acceptance test")
	}

	testAzureSubscriptionID := os.Getenv("INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID")
	if testAzureSubscriptionID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID must be set for acceptance test")
	}

	testAzureTenantID := os.Getenv("INTEGRATION_TESTING_AZURE_TENANT_ID")
	if testAzureTenantID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_TENANT_ID must be set for acceptance test")
	}

	azureLinkAccountTestConfig := map[string]string{
		"name":            testAzureLinkAccountName,
		"account_id":      strconv.Itoa(testSubAccountID),
		"application_id":  testAzureApplicationID,
		"client_secret":   testAzureClientSecretID,
		"subscription_id": testAzureSubscriptionID,
		"tenant_id":       testAzureTenantID,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "azure") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAzureLinkAccountDestroy,
		Steps: []resource.TestStep{

			// Test: Create
			{
				Config: testAccNewRelicAzureLinkAccountConfig(azureLinkAccountTestConfig, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAzureLinkAccountExists(resourceName),
				),
			},

			// Test: Update
			{
				Config: testAccNewRelicAzureLinkAccountConfig(azureLinkAccountTestConfig, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAzureLinkAccountExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_id", "client_secret", "subscription_id", "tenant_id"},
			},
		},
	})
}

func testAccCheckNewRelicAzureLinkAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)

		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")

		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if err != nil && linkedAccount == nil {
			return err
		}

		return nil

	}

}

func testAccCheckNewRelicCloudAzureLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {

		if r.Type != "newrelic_cloud_azure_link_account" {
			continue

		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked azure account still exists: #{err}")
		}
	}
	return nil
}

func testAccNewRelicAzureLinkAccountConfig(azureLinkAccountTestConfig map[string]string, updated bool) string {
	if updated == true {
		azureLinkAccountTestConfig["name"] += "-updated"
	}

	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "` + azureLinkAccountTestConfig["account_id"] + `"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_azure_link_account" "foo" {
  provider        = newrelic.cloud-integration-provider
  application_id  = "` + azureLinkAccountTestConfig["application_id"] + `"
  client_secret   = "` + azureLinkAccountTestConfig["client_secret"] + `"
  subscription_id = "` + azureLinkAccountTestConfig["subscription_id"] + `"
  tenant_id       = "` + azureLinkAccountTestConfig["tenant_id"] + `"
  name            = "` + azureLinkAccountTestConfig["name"] + `"
  account_id      = "` + azureLinkAccountTestConfig["account_id"] + `"
}
`)
}
