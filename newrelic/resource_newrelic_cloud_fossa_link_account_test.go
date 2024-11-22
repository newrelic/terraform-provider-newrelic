//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudFossaLinkAccount_Basic(t *testing.T) {
	testFossaLinkAccountName := fmt.Sprintf("tf_cloud_link_account_test_fossa_%s", acctest.RandString(5))
	resourceName := "newrelic_cloud_fossa_link_account.foo"

	// if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
	// 	t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	// }
	// testFossaApiKey := os.Getenv("INTEGRATION_TESTING_FOSSA_API_KEY")
	// if testFossaApiKey == "" {
	// 	t.Skip("INTEGRATION_TESTING_FOSSA_API_KEY must be set for acceptance test")
	// }
	// testFossaExternalKey := os.Getenv("INTEGRATION_TESTING_FOSSA_EXTERNAL_KEY")
	// if testFossaExternalKey == "" {
	// 	t.Skip("INTEGRATION_TESTING_FOSSA_EXTERNAL_KEY must be set for acceptance test")
	// }

	testFossaApiKey := "f3ab70a3d29029c94ab8057f1b30838b"
	testSubAccountID := 3806526
	testFossaExternalKey := "test-vka-external-key"

	fossaLinkAccountTestConfig := map[string]string{
		"name":         testFossaLinkAccountName,
		"api_key":      testFossaApiKey,
		"external_key": testFossaExternalKey,
		"account_id":   strconv.Itoa(testSubAccountID),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "fossa") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudfossaLinkAccountDestroy,
		Steps: []resource.TestStep{

			// Test: Create
			{
				Config: testAccNewRelicFossaLinkAccountConfig(fossaLinkAccountTestConfig, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFossaLinkAccountExists(resourceName),
				),
			},

			// Test: Update
			{
				Config: testAccNewRelicFossaLinkAccountConfig(fossaLinkAccountTestConfig, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFossaLinkAccountExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckNewRelicFossaLinkAccountExists(n string) resource.TestCheckFunc {
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

func testAccCheckNewRelicCloudFossaLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {

		if r.Type != "newrelic_cloud_fossa_link_account" {
			continue

		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked fossa account still exists: #{err}")
		}
	}
	return nil
}

func testAccNewRelicFossaLinkAccountConfig(fossaLinkAccountTestConfig map[string]string, updated bool) string {
	if updated == true {
		fossaLinkAccountTestConfig["name"] += "-updated"
	}

	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "` + fossaLinkAccountTestConfig["account_id"] + `"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_fossa_link_account" "foo" {
	name         = "` + fossaLinkAccountTestConfig["name"] + `"
	api_key      = "` + fossaLinkAccountTestConfig["api_key"] + `"
	external_key = "` + fossaLinkAccountTestConfig["external_key"] + `"
}
`)
}

func testAccCheckNewRelicCloudfossaLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {

		if r.Type != "newrelic_cloud_fossa_link_account" {
			continue

		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked fossa account still exists: #{err}")
		}
	}
	return nil
}
