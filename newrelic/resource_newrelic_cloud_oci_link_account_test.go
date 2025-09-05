//go:build integration || CLOUD
// +build integration CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicCloudOciLinkAccount_Basic(t *testing.T) {
	testOciLinkAccountName := fmt.Sprintf("tf_cloud_link_account_test_oci_%s", acctest.RandString(5))
	resourceName := "newrelic_cloud_oci_link_account.foo"

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testOciTenantID := os.Getenv("INTEGRATION_TESTING_OCI_TENANT_ID")
	if testOciTenantID == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_TENANT_ID must be set for this acceptance test")
	}

	OciLinkAccountTestConfig := map[string]string{
		"name":       testOciLinkAccountName,
		"account_id": strconv.Itoa(testSubAccountID),
		"tenant_id":  testOciTenantID,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "oci") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudOciLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicOciLinkAccountConfig(OciLinkAccountTestConfig, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudOciLinkAccountExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicOciLinkAccountConfig(OciLinkAccountTestConfig, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudOciLinkAccountExists(resourceName),
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

func testAccCheckNewRelicCloudOciLinkAccountExists(n string) resource.TestCheckFunc {
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

func testAccCheckNewRelicCloudOciLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_oci_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked oci account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicOciLinkAccountConfig(OciLinkAccountTestConfig map[string]string, updated bool) string {
	if updated == true {
		OciLinkAccountTestConfig["name"] += "_updated"
	}

	return fmt.Sprintf(`
	provider "newrelic" {
		account_id = "` + OciLinkAccountTestConfig["account_id"] + `"
		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_oci_link_account" "foo" {
		provider        	   = newrelic.cloud-integration-provider
		tenant_id              = "` + OciLinkAccountTestConfig["tenant_id"] + `"
		name                   = "` + OciLinkAccountTestConfig["name"] + `"
		account_id			   = "` + OciLinkAccountTestConfig["account_id"] + `"
	}`)
}
