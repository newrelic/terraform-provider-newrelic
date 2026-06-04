//go:build integration || CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudGcpDmLinkAccount_Basic(t *testing.T) {
	testProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	testWifCredential := os.Getenv("INTEGRATION_TESTING_GCP_WIF_CREDENTIAL")

	if testProjectID == "" || testWifCredential == "" {
		t.Skip("skipping: INTEGRATION_TESTING_GCP_PROJECT_ID and INTEGRATION_TESTING_GCP_WIF_CREDENTIAL must be set")
	}
	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skip("skipping: NEW_RELIC_SUBACCOUNT_ID must be set")
	}

	resourceName := "newrelic_cloud_gcp_dm_link_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "gcp") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudGcpDmLinkAccountDestroyed,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicCloudGcpDmLinkAccountConfig(testProjectID, testWifCredential, "tf-test-gcp-v2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudGcpDmLinkAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-gcp-v2"),
					resource.TestCheckResourceAttr(resourceName, "project_id", testProjectID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Rename (Update name only — all other fields are ForceNew)
			{
				Config: testAccNewRelicCloudGcpDmLinkAccountConfig(testProjectID, testWifCredential, "tf-test-gcp-v2-renamed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-gcp-v2-renamed"),
				),
			},
			// Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wif_credential"}, // write-only; not returned by API
			},
		},
	})
}

func testAccCheckNewRelicCloudGcpDmLinkAccountExists(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		linkedAccountID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error converting resource ID to int: %w", err)
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, linkedAccountID)
		if err != nil || linkedAccount == nil {
			return fmt.Errorf("GCP v2 linked account not found: %w", err)
		}
		return nil
	}
}

func testAccCheckNewRelicCloudGcpDmLinkAccountDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_gcp_dm_link_account" {
			continue
		}
		linkedAccountID, err := strconv.Atoi(r.Primary.ID)
		if err != nil {
			return fmt.Errorf("error converting resource ID to int: %w", err)
		}
		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, linkedAccountID)
		if linkedAccount != nil && err == nil {
			return fmt.Errorf("GCP v2 linked account still exists: %d", linkedAccountID)
		}
	}
	return nil
}

func testAccNewRelicCloudGcpDmLinkAccountConfig(projectID, wifCredential, name string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "%d"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_gcp_dm_link_account" "test" {
  provider       = newrelic.cloud-integration-provider
  account_id     = %d
  name           = %q
  project_id     = %q
  wif_credential = %q
}
`, testSubAccountID, testSubAccountID, name, projectID, wifCredential)
}
