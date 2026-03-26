//go:build integration || NGEP

package newrelic

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/newrelic/newrelic-client-go/v2/pkg/pipelinecontrol"
)

func TestAccNewRelicFederatedLogsSetup_Basic(t *testing.T) {
	t.Parallel()

	var (
		resourceName = "newrelic_federated_logs_setup.foo"

		name                       = fmt.Sprintf("test-federated-log-setup-%s", acctest.RandString(8))
		cloudProvider              = "AWS"
		cloudProviderRegion        = "us-east-1"
		dataLocationBucket         = "my-test-bucket"
		dataProcessingConnectionId = "conn-abc123"
		nrAccountId                = os.Getenv("NEW_RELIC_ACCOUNT_ID")
		nrRegion                   = "US01"
		queryConnectionId          = "qconn-xyz456"
		status                     = "ACTIVE"

		nameUpdated        = fmt.Sprintf("test-federated-log-setup-updated-%s", acctest.RandString(8))
		descriptionUpdated = "Updated description for federated log setup acceptance test."
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if nrAccountId == "" {
				t.Skip("NEW_RELIC_ACCOUNT_ID must be set for acceptance tests")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFederatedLogsSetupDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccFederatedLogsSetupConfig(
					name, "", cloudProvider, cloudProviderRegion,
					dataLocationBucket, dataProcessingConnectionId,
					nrAccountId, nrRegion, queryConnectionId, status,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogsSetupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider_region", cloudProviderRegion),
					resource.TestCheckResourceAttr(resourceName, "data_location_bucket", dataLocationBucket),
					resource.TestCheckResourceAttr(resourceName, "data_processing_connection_id", dataProcessingConnectionId),
					resource.TestCheckResourceAttr(resourceName, "nr_account_id", nrAccountId),
					resource.TestCheckResourceAttr(resourceName, "nr_region", nrRegion),
					resource.TestCheckResourceAttr(resourceName, "query_connection_id", queryConnectionId),
					resource.TestCheckResourceAttr(resourceName, "status", status),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
				),
			},
			// Test: Update (name and description)
			{
				Config: testAccFederatedLogsSetupConfig(
					nameUpdated, descriptionUpdated, cloudProvider, cloudProviderRegion,
					dataLocationBucket, dataProcessingConnectionId,
					nrAccountId, nrRegion, queryConnectionId, status,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogsSetupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "status", status),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
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

func testAccCheckNewRelicFederatedLogsSetupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_federated_logs_setup" {
			continue
		}

		_, err := client.Pipelinecontrol.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return nil // Success: resource has been destroyed.
		}

		return fmt.Errorf("federated log setup '%s' still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckNewRelicFederatedLogsSetupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no federated log setup ID is set in state")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.Pipelinecontrol.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching federated log setup with ID %s: %w", rs.Primary.ID, err)
		}

		if resp == nil {
			return fmt.Errorf("entity with ID %s returned nil response", rs.Primary.ID)
		}

		if _, ok := (*resp).(*pipelinecontrol.EntityManagementFederatedLogsSetupEntity); !ok {
			return fmt.Errorf("entity %s is not of type EntityManagementFederatedLogsSetupEntity", rs.Primary.ID)
		}

		return nil
	}
}

func testAccFederatedLogsSetupConfig(
	name, description, cloudProvider, cloudProviderRegion,
	dataLocationBucket, dataProcessingConnectionId,
	nrAccountId, nrRegion, queryConnectionId, status string,
) string {
	descriptionAttr := ""
	if description != "" {
		descriptionAttr = fmt.Sprintf(`description = "%s"`, description)
	}

	return fmt.Sprintf(`
resource "newrelic_federated_logs_setup" "foo" {
  name                          = "%s"
  %s
  cloud_provider                = "%s"
  cloud_provider_region         = "%s"
  data_location_bucket          = "%s"
  data_processing_connection_id = "%s"
  nr_account_id                 = "%s"
  nr_region                     = "%s"
  query_connection_id           = "%s"
  status                        = "%s"
}
`, name, descriptionAttr, cloudProvider, cloudProviderRegion,
		dataLocationBucket, dataProcessingConnectionId,
		nrAccountId, nrRegion, queryConnectionId, status)
}
