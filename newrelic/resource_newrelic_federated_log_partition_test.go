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

func TestAccNewRelicFederatedLogPartition_Basic(t *testing.T) {
	t.Parallel()

	var (
		partitionResourceName = "newrelic_federated_log_partition.foo"

		nrAccountId = os.Getenv("NEW_RELIC_ACCOUNT_ID")

		setupName         = fmt.Sprintf("test-fl-setup-%s", acctest.RandString(8))
		partitionName     = fmt.Sprintf("test-fl-partition-%s", acctest.RandString(8))
		dataLocationUri   = "s3://my-test-bucket/logs/"
		partitionDatabase = "my_database"
		partitionTable    = "my_table"
		status            = "ACTIVE"

		partitionNameUpdated     = fmt.Sprintf("test-fl-partition-updated-%s", acctest.RandString(8))
		descriptionUpdated       = "Updated description for federated log partition acceptance test."
		retentionDurationUpdated = 30
		retentionUnitUpdated     = "DAYS"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if nrAccountId == "" {
				t.Skip("NEW_RELIC_ACCOUNT_ID must be set for acceptance tests")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFederatedLogPartitionDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccFederatedLogPartitionConfig(
					setupName, nrAccountId,
					partitionName, "", dataLocationUri,
					partitionDatabase, partitionTable,
					status, 0, "",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogPartitionExists(partitionResourceName),
					resource.TestCheckResourceAttr(partitionResourceName, "name", partitionName),
					resource.TestCheckResourceAttr(partitionResourceName, "data_location_uri", dataLocationUri),
					resource.TestCheckResourceAttr(partitionResourceName, "partition_database", partitionDatabase),
					resource.TestCheckResourceAttr(partitionResourceName, "partition_table", partitionTable),
					resource.TestCheckResourceAttr(partitionResourceName, "nr_account_id", nrAccountId),
					resource.TestCheckResourceAttr(partitionResourceName, "status", status),
					resource.TestCheckResourceAttr(partitionResourceName, "is_default", "false"),
					resource.TestCheckResourceAttrSet(partitionResourceName, "account_id"),
					resource.TestCheckResourceAttrSet(partitionResourceName, "setup_id"),
				),
			},
			// Test: Update (name, description, and retention policy)
			{
				Config: testAccFederatedLogPartitionConfig(
					setupName, nrAccountId,
					partitionNameUpdated, descriptionUpdated, dataLocationUri,
					partitionDatabase, partitionTable,
					status, retentionDurationUpdated, retentionUnitUpdated,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogPartitionExists(partitionResourceName),
					resource.TestCheckResourceAttr(partitionResourceName, "name", partitionNameUpdated),
					resource.TestCheckResourceAttr(partitionResourceName, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(partitionResourceName, "retention_duration", fmt.Sprintf("%d", retentionDurationUpdated)),
					resource.TestCheckResourceAttr(partitionResourceName, "retention_unit", retentionUnitUpdated),
					resource.TestCheckResourceAttrSet(partitionResourceName, "account_id"),
				),
			},
			// Test: Import
			{
				ResourceName:      partitionResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNewRelicFederatedLogPartitionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_federated_log_partition" {
			continue
		}

		_, err := client.Pipelinecontrol.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return nil // Success: resource has been destroyed.
		}

		return fmt.Errorf("federated log partition '%s' still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckNewRelicFederatedLogPartitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no federated log partition ID is set in state")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.Pipelinecontrol.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching federated log partition with ID %s: %w", rs.Primary.ID, err)
		}

		if resp == nil {
			return fmt.Errorf("entity with ID %s returned nil response", rs.Primary.ID)
		}

		if _, ok := (*resp).(*pipelinecontrol.EntityManagementFederatedLogPartitionEntity); !ok {
			return fmt.Errorf("entity %s is not of type EntityManagementFederatedLogPartitionEntity", rs.Primary.ID)
		}

		return nil
	}
}

// testAccFederatedLogPartitionConfig generates HCL for a federated log partition along
// with the setup resource it depends on. Pass retentionDuration=0 to omit retention policy.
func testAccFederatedLogPartitionConfig(
	setupName, nrAccountId,
	partitionName, description, dataLocationUri,
	partitionDatabase, partitionTable,
	status string,
	retentionDuration int, retentionUnit string,
) string {
	descriptionAttr := ""
	if description != "" {
		descriptionAttr = fmt.Sprintf(`  description       = "%s"`, description)
	}

	retentionAttrs := ""
	if retentionDuration > 0 {
		retentionAttrs = fmt.Sprintf(`
  retention_duration = %d
  retention_unit     = "%s"`, retentionDuration, retentionUnit)
	}

	return fmt.Sprintf(`
resource "newrelic_federated_log_setup" "setup" {
  name                          = "%s"
  cloud_provider                = "AWS"
  cloud_provider_region         = "us-east-1"
  data_location_bucket          = "my-test-bucket"
  data_processing_connection_id = "conn-abc123"
  nr_account_id                 = "%s"
  nr_region                     = "US01"
  query_connection_id           = "qconn-xyz456"
  status                        = "ACTIVE"
}

resource "newrelic_federated_log_partition" "foo" {
  setup_id           = newrelic_federated_log_setup.setup.id
  name               = "%s"
%s
  data_location_uri  = "%s"
  is_default         = false
  nr_account_id      = "%s"
  partition_database = "%s"
  partition_table    = "%s"
  status             = "%s"
%s
}
`, setupName, nrAccountId,
		partitionName, descriptionAttr, dataLocationUri,
		nrAccountId, partitionDatabase, partitionTable, status,
		retentionAttrs)
}
