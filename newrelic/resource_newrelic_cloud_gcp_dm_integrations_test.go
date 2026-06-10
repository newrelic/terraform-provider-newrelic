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

func TestAccNewRelicCloudGcpDmIntegrations_Basic(t *testing.T) {
	testProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	testWifCredential := os.Getenv("INTEGRATION_TESTING_GCP_WIF_CREDENTIAL")

	if testProjectID == "" || testWifCredential == "" {
		t.Skip("skipping: INTEGRATION_TESTING_GCP_PROJECT_ID and INTEGRATION_TESTING_GCP_WIF_CREDENTIAL must be set")
	}
	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skip("skipping: NEW_RELIC_SUBACCOUNT_ID must be set")
	}

	resourceName := "newrelic_cloud_gcp_dm_integrations.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "gcp") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudGcpDmIntegrationsDestroyed,
		Steps: []resource.TestStep{
			// Create: link account + big_query + api_gateway
			{
				Config: testAccNewRelicCloudGcpDmIntegrationsConfig(testProjectID, testWifCredential),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudGcpDmIntegrationsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "linked_account_id"),
					resource.TestCheckResourceAttr(resourceName, "big_query.0.metrics_polling_interval", "400"),
					resource.TestCheckResourceAttr(resourceName, "big_query.0.fetch_tags", "true"),
					resource.TestCheckResourceAttr(resourceName, "big_query.0.fetch_table_metrics", "true"),
					resource.TestCheckResourceAttr(resourceName, "api_gateway.0.metrics_polling_interval", "400"),
				),
			},
			// Update: add firebase_auth, remove api_gateway
			{
				Config: testAccNewRelicCloudGcpDmIntegrationsConfigUpdated(testProjectID, testWifCredential),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudGcpDmIntegrationsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "big_query.0.metrics_polling_interval", "400"),
					resource.TestCheckResourceAttr(resourceName, "firebase_auth.0.metrics_polling_interval", "400"),
					resource.TestCheckNoResourceAttr(resourceName, "api_gateway.0.metrics_polling_interval"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNewRelicCloudGcpDmIntegrationsExists(n string) resource.TestCheckFunc {
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
		if err != nil {
			return err
		}
		if linkedAccount == nil {
			return fmt.Errorf("GCP Dimensional Metrics linked account not found: %d", linkedAccountID)
		}
		if len(linkedAccount.Integrations) == 0 {
			return fmt.Errorf("GCP Dimensional Metrics integrations not configured for linked account %d", linkedAccountID)
		}
		return nil
	}
}

func testAccCheckNewRelicCloudGcpDmIntegrationsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_gcp_dm_integrations" && r.Type != "newrelic_cloud_gcp_dm_link_account" {
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

// testAccNewRelicCloudGcpDmIntegrationsConfig creates a link account plus integrations
// with big_query (fetch_tags + fetch_table_metrics) and api_gateway enabled.
func testAccNewRelicCloudGcpDmIntegrationsConfig(projectID, wifCredential string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "%d"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_gcp_dm_link_account" "test" {
  provider       = newrelic.cloud-integration-provider
  account_id     = %d
  name           = "tf-test-gcp-dm-integrations"
  project_id     = %q
  wif_credential = %q
}

resource "newrelic_cloud_gcp_dm_integrations" "test" {
  provider          = newrelic.cloud-integration-provider
  account_id        = %d
  linked_account_id = newrelic_cloud_gcp_dm_link_account.test.id

  big_query {
    metrics_polling_interval = 400
    fetch_tags               = true
    fetch_table_metrics      = true
  }

  api_gateway {
    metrics_polling_interval = 400
  }
}
`, testSubAccountID, testSubAccountID, projectID, wifCredential, testSubAccountID)
}

// testAccNewRelicCloudGcpDmIntegrationsConfigUpdated keeps big_query, adds firebase_auth,
// and removes api_gateway.
func testAccNewRelicCloudGcpDmIntegrationsConfigUpdated(projectID, wifCredential string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "%d"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_gcp_dm_link_account" "test" {
  provider       = newrelic.cloud-integration-provider
  account_id     = %d
  name           = "tf-test-gcp-dm-integrations"
  project_id     = %q
  wif_credential = %q
}

resource "newrelic_cloud_gcp_dm_integrations" "test" {
  provider          = newrelic.cloud-integration-provider
  account_id        = %d
  linked_account_id = newrelic_cloud_gcp_dm_link_account.test.id

  big_query {
    metrics_polling_interval = 400
    fetch_tags               = true
    fetch_table_metrics      = true
  }

  firebase_auth {
    metrics_polling_interval = 400
  }
}
`, testSubAccountID, testSubAccountID, projectID, wifCredential, testSubAccountID)
}
