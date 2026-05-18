//go:build integration || LOGGING_INTEGRATIONS

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccNewRelicFederatedLogsSetup_Basic exercises create / read / update /
// import / destroy for the newrelic_federated_logs_setup resource.
func TestAccNewRelicFederatedLogsSetup_Basic(t *testing.T) {
	resourceName := "newrelic_federated_logs_setup.foo"
	rName := generateNameForIntegrationTestResource()
	roleArn := "arn:aws:iam::123456789012:role/tf-test-role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFederatedLogsSetupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFederatedLogsSetupConfig(rName, roleArn, "Initial setup"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogsSetupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial setup"),
					resource.TestCheckResourceAttr(resourceName, "storage.0.data_location_bucket", "tf-test-fed-logs-bucket"),
					resource.TestCheckResourceAttr(resourceName, "storage.0.database", "tf_test_fed_logs_db"),
					resource.TestCheckResourceAttr(resourceName, "storage.0.cloud_provider_configuration.0.provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "storage.0.cloud_provider_configuration.0.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "default_partition.0.storage.0.table", "log_transactions_default"),
					resource.TestCheckResourceAttrSet(resourceName, "default_partition_id"),
					resource.TestCheckResourceAttrSet(resourceName, "lifecycle_status.0.status"),
				),
			},
			// Update name and description (mutable via FederatedLogsUpdateSetup).
			{
				Config: testAccNewRelicFederatedLogsSetupConfig(rName+"-updated", roleArn, "Updated setup"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicFederatedLogsSetupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated setup"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Computed lifecycle/health-check timestamps drift between create
				// response and the post-import read; skip strict comparison there.
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

func testAccNewRelicFederatedLogsSetupConfig(name, roleArn, description string) string {
	return fmt.Sprintf(`
resource "newrelic_aws_connection" "ingest" {
  name     = "%[1]s-ingest"
  enabled  = true
  region   = "us-east-1"
  role_arn = "%[2]s"
}

resource "newrelic_aws_connection" "query" {
  name     = "%[1]s-query"
  enabled  = true
  region   = "us-east-1"
  role_arn = "%[2]s"
}

resource "newrelic_federated_logs_setup" "foo" {
  name        = "%[1]s"
  description = "%[3]s"

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

    data_retention_policy {
      duration = 30
      unit     = "DAYS"
    }
  }
}
`, name, roleArn, description)
}

func testAccCheckNewRelicFederatedLogsSetupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no setup ID set in state")
		}
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.Federatedlogs.GetSetupWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching federated logs setup %s: %w", rs.Primary.ID, err)
		}
		if resp == nil {
			return fmt.Errorf("federated logs setup %s returned nil from API", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckNewRelicFederatedLogsSetupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_federated_logs_setup" {
			continue
		}
		resp, err := client.Federatedlogs.GetSetupWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			// API returns an error for deleted IDs — that's the destroy success path.
			return nil
		}
		if resp != nil {
			return fmt.Errorf("federated logs setup %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}

// TestAccNewRelicFederatedLogsSetup_ImmutableFieldsError pins the contract
// enforced by validateFederatedLogsSetupDiff: editing a create-only field on
// an existing setup must error at plan time rather than silently dropping the
// change or destructively recreating the resource. The first step creates the
// setup; the second step tries to change the bucket and asserts the plan
// fails with the expected error.
func TestAccNewRelicFederatedLogsSetup_ImmutableFieldsError(t *testing.T) {
	resourceName := "newrelic_federated_logs_setup.foo"
	rName := generateNameForIntegrationTestResource()
	roleArn := "arn:aws:iam::123456789012:role/tf-test-role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFederatedLogsSetupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFederatedLogsSetupConfig(rName, roleArn, "Initial setup"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckNewRelicFederatedLogsSetupExists(resourceName)),
			},
			{
				Config:      testAccNewRelicFederatedLogsSetupConfigDifferentBucket(rName, roleArn),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`storage\.data_location_bucket cannot be updated after creation`),
			},
		},
	})
}

// testAccNewRelicFederatedLogsSetupConfigDifferentBucket is identical to
// testAccNewRelicFederatedLogsSetupConfig except the bucket is changed; used
// to drive the ImmutableFieldsError test.
func testAccNewRelicFederatedLogsSetupConfigDifferentBucket(name, roleArn string) string {
	return fmt.Sprintf(`
resource "newrelic_aws_connection" "ingest" {
  name     = "%[1]s-ingest"
  enabled  = true
  region   = "us-east-1"
  role_arn = "%[2]s"
}

resource "newrelic_aws_connection" "query" {
  name     = "%[1]s-query"
  enabled  = true
  region   = "us-east-1"
  role_arn = "%[2]s"
}

resource "newrelic_federated_logs_setup" "foo" {
  name        = "%[1]s"
  description = "Initial setup"

  storage {
    data_location_bucket      = "tf-test-fed-logs-bucket-changed"
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
      data_location_uri = "s3://tf-test-fed-logs-bucket-changed/log_transactions_default"
    }

    data_retention_policy {
      duration = 30
      unit     = "DAYS"
    }
  }
}
`, name, roleArn)
}
