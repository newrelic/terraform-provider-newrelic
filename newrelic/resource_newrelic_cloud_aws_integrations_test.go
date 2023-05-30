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
		arn                    = "%[1]s"
		metric_collection_mode = "PULL"
		name                   = "integration test account"
	  }

	resource "newrelic_cloud_aws_integrations" "bar" {
		linked_account_id = newrelic_cloud_aws_link_account.foo.id
	  
		billing {
		  metrics_polling_interval = 3600
		}
		cloudtrail {
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		}
		health {
		  metrics_polling_interval = 3600
		}
		trusted_advisor {
		  metrics_polling_interval = 3600
		}
		vpc {
		  aws_regions              = ["us-east-1"]
		  fetch_nat_gateway        = true
		  fetch_vpn                = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		x_ray {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		s3 {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		doc_db {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		sqs {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  queue_prefixes           = ["test prefix"]
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ebs {
		  metrics_polling_interval = 3600
		  fetch_extended_inventory = true
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		alb {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  load_balancer_prefixes   = ["test prefix"]
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticache {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		api_gateway {
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		  stage_prefixes           = ["test prefix"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		auto_scaling {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_app_sync {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_athena {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_cognito {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_direct_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_fsx {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_glue {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_kinesis_analytics {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_media_convert {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_media_package_vod {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_mq {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_msk {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_neptune {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_qldb {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_route53resolver {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_states {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_transit_gateway {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_waf {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_wafv2 {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		cloudfront {
		  fetch_lambdas_at_edge    = true
		  fetch_tags               = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
	  }
`, arn)
}

func testAccNewRelicAwsIntegrationsConfigUpdated(arn string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_aws_link_account" "foo" {
		arn                    = "%[1]s"
		metric_collection_mode = "PULL"
		name                   = "integration test account - updated"
	  }
	resource "newrelic_cloud_aws_integrations" "bar" {
		linked_account_id = newrelic_cloud_aws_link_account.foo.id
		billing {
		  metrics_polling_interval = 21600
		}
		cloudtrail {
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		}
		health {
		  metrics_polling_interval = 3600
		}
		trusted_advisor {
		  metrics_polling_interval = 3600
		}
		vpc {
		  aws_regions              = ["us-east-1"]
		  fetch_nat_gateway        = true
		  fetch_vpn                = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		x_ray {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		s3 {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		doc_db {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		sqs {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  queue_prefixes           = ["test prefix"]
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ebs {
		  metrics_polling_interval = 3600
		  fetch_extended_inventory = true
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		alb {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  load_balancer_prefixes   = ["test prefix"]
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticache {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		api_gateway {
		  metrics_polling_interval = 3600
		  aws_regions              = ["us-east-1"]
		  stage_prefixes           = ["test prefix"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		auto_scaling {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_app_sync {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_athena {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_cognito {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_direct_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_fsx {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_glue {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_kinesis_analytics {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_media_convert {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_media_package_vod {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_mq {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_msk {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_neptune {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_qldb {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_route53resolver {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_states {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_transit_gateway {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_waf {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		aws_wafv2 {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 3600
		}
		cloudfront {
		  fetch_lambdas_at_edge    = true
		  fetch_tags               = true
		  metrics_polling_interval = 3600
		  tag_key                  = "test"
		  tag_value                = "test"
		}
	  }
`, arn)
}
