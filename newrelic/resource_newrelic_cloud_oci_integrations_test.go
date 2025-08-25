//go:build integration || CLOUD
// +build integration CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudOciIntegrations_Basic(t *testing.T) {
	t.Skipf("Skipping test until OCI Environment Variables are fixed")
	resourceName := "newrelic_cloud_oci_integrations.foo1"
	testOciIntegrationName := fmt.Sprintf("tf_cloud_integrations_test_oci_%s", acctest.RandString(5))

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testOciTenantID := os.Getenv("INTEGRATION_TESTING_OCI_TENANT_ID")
	if testOciTenantID == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_TENANT_ID must be set for acceptance test")
	}

	OciIntegrationTestConfig := map[string]string{
		"name":       testOciIntegrationName,
		"account_id": strconv.Itoa(testSubAccountID),
		"tenant_id":  testOciTenantID,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "oci") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudOciIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudOciIntegrationsConfig(OciIntegrationTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudOciIntegrationsExists(resourceName),
				),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
			},
			//Test: Update
			{
				Config: testAccNewRelicCloudOciIntegrationsConfigUpdated(OciIntegrationTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudOciIntegrationsExists(resourceName),
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

func testAccNewRelicCloudOciIntegrationsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)
		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			return fmt.Errorf("An error occurred creating GCP integrations")
		}

		return nil
	}
}

func testAccNewRelicCloudOciIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_oci_integrations" && r.Type != "newrelic_cloud_oci_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("GCP Linked account is not unlinked: #{err}")
		}
	}
	return nil
}

func testAccNewRelicCloudOciIntegrationsConfig(OciIntegrationTestConfig map[string]string) string {
	return fmt.Sprintf(`
		provider "newrelic" {
			account_id = "` + OciIntegrationTestConfig["account_id"] + `"
			alias      = "cloud-integration-provider"
		}

		resource "newrelic_cloud_oci_link_account" "foo" {
			provider        = newrelic.cloud-integration-provider
			name       = "` + OciIntegrationTestConfig["name"] + `"
			account_id = "` + OciIntegrationTestConfig["account_id"] + `"
			tenant_id = "` + OciIntegrationTestConfig["tenant_id"] + `"
		}

		resource "newrelic_cloud_oci_integrations" "foo1" {
		provider        = newrelic.cloud-integration-provider
		account_id 		= "` + OciIntegrationTestConfig["account_id"] + `"
		linked_account_id = newrelic_cloud_oci_link_account.foo.id
		oci_metadata_and_tags {}
		}
	`)
}

func testAccNewRelicCloudOciIntegrationsConfigUpdated(OciIntegrationTestConfig map[string]string) string {
	return fmt.Sprintf(`
		provider "newrelic" {
			account_id = "` + OciIntegrationTestConfig["account_id"] + `"
			alias      = "cloud-integration-provider"
		}

		resource "newrelic_cloud_oci_link_account" "foo" {
			provider        = newrelic.cloud-integration-provider
			name       = "` + OciIntegrationTestConfig["name"] + `"
			account_id = "` + OciIntegrationTestConfig["account_id"] + `"
			tenant_id = "` + OciIntegrationTestConfig["tenant_id"] + `"
		}

		resource "newrelic_cloud_oci_integrations" "foo1" {
			provider        = newrelic.cloud-integration-provider
			account_id 		= "` + OciIntegrationTestConfig["account_id"] + `"
			linked_account_id = newrelic_cloud_oci_link_account.foo.id
			oci_metadata_and_tags {}
		}
	`)
}
