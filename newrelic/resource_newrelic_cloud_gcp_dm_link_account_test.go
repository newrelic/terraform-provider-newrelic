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
	testAudience := os.Getenv("INTEGRATION_TESTING_GCP_WIF_AUDIENCE")
	testSAEmail := os.Getenv("INTEGRATION_TESTING_GCP_WIF_SA_EMAIL")

	if testProjectID == "" || testAudience == "" || testSAEmail == "" {
		t.Skip("skipping: INTEGRATION_TESTING_GCP_PROJECT_ID, INTEGRATION_TESTING_GCP_WIF_AUDIENCE and INTEGRATION_TESTING_GCP_WIF_SA_EMAIL must be set")
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
				Config: testAccNewRelicCloudGcpDmLinkAccountConfig(testProjectID, testAudience, testSAEmail, "tf-test-gcp-dm"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudGcpDmLinkAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-gcp-dm"),
					resource.TestCheckResourceAttr(resourceName, "project_id", testProjectID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Rename (Update name only — all other fields are ForceNew)
			{
				Config: testAccNewRelicCloudGcpDmLinkAccountConfig(testProjectID, testAudience, testSAEmail, "tf-test-gcp-dm-renamed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-gcp-dm-renamed"),
				),
			},
			// Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"audience", "service_account_email"}, // write-only; not returned by API
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
			return fmt.Errorf("GCP Dimensional Metrics linked account not found: %w", err)
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
			return fmt.Errorf("GCP Dimensional Metrics linked account still exists: %d", linkedAccountID)
		}
	}
	return nil
}

func testAccNewRelicCloudGcpDmLinkAccountConfig(projectID, audience, serviceAccountEmail, name string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "%d"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_gcp_dm_link_account" "test" {
  provider              = newrelic.cloud-integration-provider
  account_id            = %d
  name                  = %q
  project_id            = %q
  audience              = %q
  service_account_email = %q
}
`, testSubAccountID, testSubAccountID, name, projectID, audience, serviceAccountEmail)
}
