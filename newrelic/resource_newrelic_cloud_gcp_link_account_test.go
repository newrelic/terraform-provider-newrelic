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

func TestAccNewRelicCloudGcpLinkAccount(t *testing.T) {
	resourceName := "newrelic_cloud_gcp_link_account.foo"

	testGcpProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	if testGcpProjectID == "" {
		t.Skipf("INTEGRATION_TESTING_GCP_PROJECT_ID must be set for acceptance test")
	}

	testGcpAccountName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudGcpLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudGcpLinkAccountConfig(testGcpAccountName, testGcpProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpLinkAccountExists(resourceName),
				),
			},
			//Test: Update

			// NOTE: Skipping this step due to an API issue.

			//{
			//	Config: testAccNewRelicCloudGcpLinkAccountConfigUpdated(testGcpAccountName, testGcpProjectID),
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

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

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

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("Linked gcp account still exists: #{err}")
		}

	}
	return nil
}

func testAccNewRelicCloudGcpLinkAccountConfig(name string, projectId string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_gcp_link_account" "foo"{
			name = "%[1]s"
			account_id = %[3]d
			project_id="%[2]s"
	}
	`, name, projectId, testAccountID)
}

func testAccNewRelicCloudGcpLinkAccountConfigUpdated(name string, projectId string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_gcp_link_account" "foo"{
			name = "%[1]s-updated"
			account_id = %[3]d
			project_id="%[2]s"
	}
	`, name, projectId, testAccountID)
}
