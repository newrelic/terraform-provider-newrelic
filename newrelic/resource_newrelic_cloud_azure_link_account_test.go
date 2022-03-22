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
	randName := acctest.RandString(5)
	resourceName := "newrelic_cloud_azure_link_account.foo"

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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAzureLinkAccountDestroy,
		Steps: []resource.TestStep{

			// Test: Create
			{
				Config: testAccNewRelicAzureLinkAccountConfig(testAzureApplicationID, testAzureClientSecretID, testAzureSubscriptionID, testAzureTenantID, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAzureLinkAccountExists(resourceName),
				),
			},

			// Test: Update
			{
				Config: testAccNewRelicAzureLinkAccountConfigUpdated(testAzureApplicationID, testAzureClientSecretID, testAzureSubscriptionID, testAzureSubscriptionID, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAzureLinkAccountExists(resourceName),
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
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

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
			fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked azure account still exists: #{err}")
		}
	}
	return nil
}

func testAccNewRelicAzureLinkAccountConfig(applicationID string, clientSecretID string, subscriptionID string, tenantID string, name string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_azure_link_account" "foo"{
	application_id = "%[1]s"
	client_secret = "%[2]s"
	subscription_id = "%[3]s"
	tenant_id = "%[4]s"
	name  = "%[5]s"
}
	`, applicationID, clientSecretID, subscriptionID, tenantID, name)
}

func testAccNewRelicAzureLinkAccountConfigUpdated(applicationID string, clientSecretID string, subscriptionID string, tenantID string, name string) string {
	return fmt.Sprintf(`
   resource "newrelic_cloud_azure_link_account" "foo"{
      application_id = "%[1]s"
       client_secret = "%[2]s"
       subscription_id = "%[3]s"
       tenant_id = "%[4]s"
       name = "%[5]s-updated"
   }
   `, applicationID, clientSecretID, subscriptionID, tenantID, name)
}
