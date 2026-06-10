//go:build LOGGING_INTEGRATIONS

// NOTE: Removed the "integration" tag so these tests no longer run under
// the standard `make test-integration-all` (-tags=integration) CI job.
// The federated-logs API gateway returns ACCESS_DENIED for the shared
// Terraform provider test account because the account lacks the
// federated_logs entitlement on productLine=generic. Re-add the
// "integration" tag once the entitlement is granted.
// Mirrors newrelic-client-go#1425.

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/newrelic/newrelic-client-go/v2/pkg/federatedlogs"
)

// TestAccNewRelicFederatedLogsPartition_Basic exercises the full chain:
// AWS connections → setup → partition. The partition is a non-default
// partition; the default partition is created by the setup itself.
//
// Required env vars: same as the setup test (NEW_RELIC_API_KEY,
// NEW_RELIC_ACCOUNT_ID) plus the federated logs feature flag on the account.
func TestAccNewRelicFederatedLogsPartition_Basic(t *testing.T) {
	resourceName := "newrelic_federated_logs_partition.foo"
	rName := generateNameForIntegrationTestResource()
	roleArn := "arn:aws:iam::123456789012:role/tf-test-role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFederatedLogsPartitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFederatedLogsPartitionConfig(rName, roleArn, "Initial partition", true, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogsPartitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-partition"),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial partition"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage.0.table", "log_transactions_secondary"),
					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_retention_policy.0.duration", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "setup_id"),
					resource.TestCheckResourceAttrSet(resourceName, "lifecycle_status.0.status"),
				),
			},
			// Update description + active + retention duration (all mutable per FederatedLogsUpdatePartitionInput).
			{
				Config: testAccNewRelicFederatedLogsPartitionConfig(rName, roleArn, "Updated partition", false, 60),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogsPartitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated partition"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_retention_policy.0.duration", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"created_at",
					"updated_at",
					"lifecycle_status",
					"health_check",
				},
			},
		},
	})
}

func testAccNewRelicFederatedLogsPartitionConfig(name, roleArn, description string, active bool, retentionDuration int) string {
	return fmt.Sprintf(`
resource "newrelic_aws_connection" "ingest" {
  name       = "%[1]s-ingest"
  enabled    = true
  region     = "us-east-1"
  role_arn   = "%[2]s"
  scope_type = "ORGANIZATION"
  scope_id   = "fb33fea3-4d7e-4736-9701-acb59a634fdf"
}

resource "newrelic_aws_connection" "query" {
  name       = "%[1]s-query"
  enabled    = true
  region     = "us-east-1"
  role_arn   = "%[2]s"
  scope_type = "ORGANIZATION"
  scope_id   = "fb33fea3-4d7e-4736-9701-acb59a634fdf"
}

resource "newrelic_federated_logs_setup" "parent" {
  name        = "%[1]s-setup"
  description = "Parent setup for partition acceptance test"

  storage {
    data_location_bucket      = "tf-test-fed-logs-bucket"
    database                  = "tf_test_fed_logs_db"
    data_ingest_connection_id = newrelic_aws_connection.ingest.id
    query_connection_id       = newrelic_aws_connection.query.id

    cloud_provider_configuration {
      provider = "AWS"
      region   = "us-east-1"
    }
  }

  default_partition {
    storage {
      table             = "log_transactions_default"
      data_location_uri = "s3://tf-test-fed-logs-bucket/log_transactions_default"
    }
  }
}

resource "newrelic_federated_logs_partition" "foo" {
  setup_id    = newrelic_federated_logs_setup.parent.id
  name        = "%[1]s-partition"
  description = "%[3]s"
  active      = %[4]t

  storage {
    table             = "log_transactions_secondary"
    data_location_uri = "s3://tf-test-fed-logs-bucket/log_transactions_secondary"
  }

  data_retention_policy {
    duration = %[5]d
    unit     = "DAYS"
  }
}
`, name, roleArn, description, active, retentionDuration)
}

func testAccCheckNewRelicFederatedLogsPartitionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no partition ID set in state")
		}
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient
		resp, err := client.Federatedlogs.GetPartitionWithContext(context.Background(), providerConfig.AccountID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching federated logs partition %s: %w", rs.Primary.ID, err)
		}
		if resp == nil {
			return fmt.Errorf("federated logs partition %s returned nil from API", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicFederatedLogsPartitionDestroy(s *terraform.State) error {
	// Partition deletion is a soft-delete via lifecycleStatus DELETING. The
	// entity stays queryable state.
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_federated_logs_partition" {
			continue
		}
		resp, err := client.Federatedlogs.GetPartitionWithContext(context.Background(), providerConfig.AccountID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("expected partition %s to be queryable in DELETING state, got error: %w", rs.Primary.ID, err)
		}
		if resp == nil {
			return fmt.Errorf("expected partition %s to be queryable in DELETING state, got nil response", rs.Primary.ID)
		}
		if string(resp.LifecycleStatus.Status) != string(federatedlogs.FederatedLogsLifecycleStateTypes.DELETING) {
			return fmt.Errorf("federated logs partition %s expected lifecycleStatus DELETING, got %q",
				rs.Primary.ID, resp.LifecycleStatus.Status)
		}
	}
	return nil
}

// TestAccNewRelicFederatedLogsPartition_ImmutableFieldsError pins the contract
// enforced by validateFederatedLogsPartitionDiff: editing a create-only field
// on an existing partition must error at plan time. The first step creates
// the partition; the second tries to change storage.table and asserts the
// plan fails with the expected error.
func TestAccNewRelicFederatedLogsPartition_ImmutableFieldsError(t *testing.T) {
	resourceName := "newrelic_federated_logs_partition.foo"
	rName := generateNameForIntegrationTestResource()
	roleArn := "arn:aws:iam::123456789012:role/tf-test-role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFederatedLogsPartitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFederatedLogsPartitionConfig(rName, roleArn, "Initial partition", true, 30),
				Check:  resource.ComposeTestCheckFunc(testAccCheckNewRelicFederatedLogsPartitionExists(resourceName)),
			},
			{
				Config:      testAccNewRelicFederatedLogsPartitionConfigDifferentTable(rName, roleArn),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`storage\.table cannot be updated after creation`),
			},
		},
	})
}

// testAccNewRelicFederatedLogsPartitionConfigDifferentTable mirrors the basic
// partition config but with a different storage.table — used to drive the
// ImmutableFieldsError test.
func testAccNewRelicFederatedLogsPartitionConfigDifferentTable(name, roleArn string) string {
	return fmt.Sprintf(`
resource "newrelic_aws_connection" "ingest" {
  name       = "%[1]s-ingest"
  enabled    = true
  region     = "us-east-1"
  role_arn   = "%[2]s"
  scope_type = "ORGANIZATION"
  scope_id   = "fb33fea3-4d7e-4736-9701-acb59a634fdf"
}

resource "newrelic_aws_connection" "query" {
  name       = "%[1]s-query"
  enabled    = true
  region     = "us-east-1"
  role_arn   = "%[2]s"
  scope_type = "ORGANIZATION"
  scope_id   = "fb33fea3-4d7e-4736-9701-acb59a634fdf"
}

resource "newrelic_federated_logs_setup" "parent" {
  name        = "%[1]s-setup"
  description = "Parent setup for partition acceptance test"

  storage {
    data_location_bucket      = "tf-test-fed-logs-bucket"
    database                  = "tf_test_fed_logs_db"
    data_ingest_connection_id = newrelic_aws_connection.ingest.id
    query_connection_id       = newrelic_aws_connection.query.id

    cloud_provider_configuration {
      provider = "AWS"
      region   = "us-east-1"
    }
  }

  default_partition {
    storage {
      table             = "log_transactions_default"
      data_location_uri = "s3://tf-test-fed-logs-bucket/log_transactions_default"
    }
  }
}

resource "newrelic_federated_logs_partition" "foo" {
  setup_id    = newrelic_federated_logs_setup.parent.id
  name        = "%[1]s-partition"
  description = "Initial partition"

  storage {
    table             = "log_transactions_changed"
    data_location_uri = "s3://tf-test-fed-logs-bucket/log_transactions_changed"
  }

  data_retention_policy {
    duration = 30
    unit     = "DAYS"
  }
}
`, name, roleArn)
}
