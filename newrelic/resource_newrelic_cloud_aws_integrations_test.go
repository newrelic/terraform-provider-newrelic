//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAwsIntegrations_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_aws_integrations.bar"

	testAwsArn := os.Getenv("INTEGRATION_TESTING_AWS_ARN")
	if testAwsArn == "" {
		t.Skipf("INTEGRATION_TESTING_AWS_ARN must be set for this acceptance test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsIntegrationsConfig(testAwsArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsIntegrationsExist(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicAwsIntegrationsConfigUpdated(testAwsArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsIntegrationsExist(resourceName),
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

func testAccCheckNewRelicCloudAwsIntegrationsExist(n string) resource.TestCheckFunc {
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

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			return fmt.Errorf("An error occurred creating AWS integrations")
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAwsIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_aws_integrations" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked aws account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicAwsIntegrationsConfig(arn string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_aws_link_account" "foo" {
		arn = "%[1]s"
		metric_collection_mode = "PULL"
		name = "integration test account"
	}

	resource "newrelic_cloud_aws_integrations" "bar" {
		linked_account_id = newrelic_cloud_aws_link_account.foo.id

		billing {
			metrics_polling_interval = 6000
		}
		cloudtrail {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
		}
		health {
			metrics_polling_interval = 6000
		}
		trusted_advisor {
			metrics_polling_interval = 6000
		}
		vpc {
			aws_regions = ["us-east-1"]
			fetch_nat_gateway = true
			fetch_vpn = true
			metrics_polling_interval = 6000
			tag_key = "test"
			tag_value = "test"
		}
		x_ray {
			aws_regions = ["us-east-1"]
			metrics_polling_interval = 6000
		}
		s3 {
			metrics_polling_interval = 6000
		}
		doc_db {
			metrics_polling_interval = 6000
		}
		sqs {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
			tag_key = "test"
			tag_value = "test"
		}
		ebs {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
			tag_key = "test"
			tag_value = "test"
		}
		alb {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
		}
		elasticache {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
		}
	}
`, arn)
}

func testAccNewRelicAwsIntegrationsConfigUpdated(arn string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_aws_link_account" "foo" {
		arn = "%[1]s"
		metric_collection_mode = "PULL"
		name = "integration test account - updated"
	}
	resource "newrelic_cloud_aws_integrations" "bar" {
		linked_account_id = newrelic_cloud_aws_link_account.foo.id
		billing {
			metrics_polling_interval = 10000
		}
		cloudtrail {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
		}
		health {
			metrics_polling_interval = 6000
		}
		trusted_advisor {
			metrics_polling_interval = 6000
		}
		vpc {
			aws_regions = ["us-east-1"]
			fetch_nat_gateway = true
			fetch_vpn = true
			metrics_polling_interval = 6000
			tag_key = "test"
			tag_value = "test"
		}
		x_ray {
			aws_regions = ["us-east-1"]
			metrics_polling_interval = 6000
		}
		s3 {
			metrics_polling_interval = 6000
		}
		doc_db {
			metrics_polling_interval = 6000
		}
		sqs {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
			tag_key = "test"
			tag_value = "test"
		}
		ebs {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
			tag_key = "test"
			tag_value = "test"
		}
		alb {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
		}
		elasticache {
			metrics_polling_interval = 6000
			aws_regions = ["us-east-1"]
		}
	}
`, arn)
}
