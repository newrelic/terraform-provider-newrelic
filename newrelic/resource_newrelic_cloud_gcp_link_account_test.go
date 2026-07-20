//go:build integration || CLOUD

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

func TestAccNewRelicCloudGcpLinkAccount(t *testing.T) {
	t.Skipf("Skipping test until GCP Environment Variables are fixed")
	resourceName := "newrelic_cloud_gcp_link_account.foo"
	testGCPLinkAccountName := fmt.Sprintf("tf_cloud_link_account_test_gcp_%s", acctest.RandString(5))

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testGCPProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	if testGCPProjectID == "" {
		t.Skipf("INTEGRATION_TESTING_GCP_PROJECT_ID must be set for acceptance test")
	}

	GCPLinkAccountTestConfig := map[string]string{
		"name":       testGCPLinkAccountName,
		"account_id": strconv.Itoa(testSubAccountID),
		"project_id": testGCPProjectID,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "gcp") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudGcpLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudGcpLinkAccountConfig(GCPLinkAccountTestConfig, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpLinkAccountExists(resourceName),
				),
			},
			//Test: Update

			// NOTE: Skipping this step due to an API issue.

			//{
			//	Config: testAccNewRelicCloudGcpLinkAccountConfig(GCPLinkAccountTestConfig, true),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccNewRelicCloudGcpLinkAccountExists(resourceName),
			//	),
			//},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicCloudGcpLinkAccount_WIF exercises the GCP Dimensional Metrics
// (keyless) linking mode of newrelic_cloud_gcp_link_account, where audience and
// service_account_email are supplied and the account is linked via Workload
// Identity Federation instead of a service-account key.
func TestAccNewRelicCloudGcpLinkAccount_WIF(t *testing.T) {
	testProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	testAudience := os.Getenv("INTEGRATION_TESTING_GCP_WIF_AUDIENCE")
	testSAEmail := os.Getenv("INTEGRATION_TESTING_GCP_WIF_SA_EMAIL")

	if testProjectID == "" || testAudience == "" || testSAEmail == "" {
		t.Skip("skipping: INTEGRATION_TESTING_GCP_PROJECT_ID, INTEGRATION_TESTING_GCP_WIF_AUDIENCE and INTEGRATION_TESTING_GCP_WIF_SA_EMAIL must be set")
	}
	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skip("skipping: NEW_RELIC_SUBACCOUNT_ID must be set")
	}

	resourceName := "newrelic_cloud_gcp_link_account.wif"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "gcp") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudGcpLinkAccountDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicCloudGcpLinkAccountWIFConfig(testProjectID, testAudience, testSAEmail, "tf-test-gcp-dm"),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpLinkAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-gcp-dm"),
					resource.TestCheckResourceAttr(resourceName, "project_id", testProjectID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Rename (only name is mutable; all other fields are ForceNew)
			{
				Config: testAccNewRelicCloudGcpLinkAccountWIFConfig(testProjectID, testAudience, testSAEmail, "tf-test-gcp-dm-renamed"),
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

func testAccNewRelicCloudGcpLinkAccountWIFConfig(projectID, audience, serviceAccountEmail, name string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "%d"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_gcp_link_account" "wif" {
  provider                         = newrelic.cloud-integration-provider
  account_id                       = %d
  name                             = %q
  project_id                       = %q
  use_workload_identity_federation = true
  audience                         = %q
  service_account_email            = %q
}
`, testSubAccountID, testSubAccountID, name, projectID, audience, serviceAccountEmail)
}

func testAccNewRelicCloudGcpLinkAccountExists(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if err != nil && linkedAccount == nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicCloudGcpLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_gcp_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)
		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("Linked gcp account still exists: #{err}")
		}

	}
	return nil
}

func testAccNewRelicCloudGcpLinkAccountConfig(GCPLinkAccountTestConfig map[string]string, updated bool) string {
	if updated == true {
		GCPLinkAccountTestConfig["name"] += "_updated"
	}

	return `
	provider "newrelic" {
  		account_id = "` + GCPLinkAccountTestConfig["account_id"] + `"
  		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_gcp_link_account" "foo"{
            provider        = newrelic.cloud-integration-provider
			name 		= "` + GCPLinkAccountTestConfig["name"] + `"
            account_id  = "` + GCPLinkAccountTestConfig["account_id"] + `"
			project_id  = "` + GCPLinkAccountTestConfig["project_id"] + `"
	}
	`
}
