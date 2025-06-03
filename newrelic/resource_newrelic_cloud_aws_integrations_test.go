//go:build integration || CLOUD
// +build integration CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAwsIntegrations_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_aws_integrations.bar"
	testAWSIntegrationName := fmt.Sprintf("tf_cloud_integrations_test_aws_%s", acctest.RandString(5))

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testAWSArn := os.Getenv("INTEGRATION_TESTING_AWS_ARN")
	if testAWSArn == "" {
		t.Skipf("INTEGRATION_TESTING_AWS_ARN must be set for this acceptance test")
	}

	AWSIntegrationsTestConfig := map[string]string{
		"name":       testAWSIntegrationName,
		"account_id": strconv.Itoa(testSubAccountID),
		"arn":        testAWSArn,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "aws") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsIntegrationsConfig(AWSIntegrationsTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsIntegrationsExist(resourceName),
				),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
			},
			//Test: Update
			{
				Config: testAccNewRelicAwsIntegrationsConfigUpdated(AWSIntegrationsTestConfig),
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

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

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
		if r.Type != "newrelic_cloud_aws_integrations" && r.Type != "newrelic_cloud_aws_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked aws account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicAwsIntegrationsConfig(AWSIntegrationsTestConfig map[string]string) string {
	return fmt.Sprintf(`
	provider "newrelic" {
  		account_id = "` + AWSIntegrationsTestConfig["account_id"] + `"
  		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_aws_link_account" "foo" {
        provider        	   = newrelic.cloud-integration-provider
		arn                    = "` + AWSIntegrationsTestConfig["arn"] + `"
		metric_collection_mode = "PULL"
		name                   = "` + AWSIntegrationsTestConfig["name"] + `"
		account_id			   = "` + AWSIntegrationsTestConfig["account_id"] + `"
	  }

	  resource "newrelic_cloud_aws_integrations" "bar" {
		provider        	   = newrelic.cloud-integration-provider
		linked_account_id 	   = newrelic_cloud_aws_link_account.foo.id
	    account_id             = "` + AWSIntegrationsTestConfig["account_id"] + `"
		billing {
		  metrics_polling_interval = 6000
		}
		cloudtrail {
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		}
		health {
		  metrics_polling_interval = 6000
		}
		trusted_advisor {
		  metrics_polling_interval = 6000
		}
		vpc {
		  aws_regions              = ["us-east-1"]
		  fetch_nat_gateway        = true
		  fetch_vpn                = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		x_ray {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		s3 {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		doc_db {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		sqs {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  queue_prefixes           = ["test prefix"]
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ebs {
		  metrics_polling_interval = 6000
		  fetch_extended_inventory = true
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		alb {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  load_balancer_prefixes   = ["test prefix"]
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticache {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		api_gateway {
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		  stage_prefixes           = ["test prefix"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		auto_scaling {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_app_sync {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_athena {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_cognito {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_direct_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_fsx {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_glue {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_kinesis_analytics {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_media_convert {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_media_package_vod {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_mq {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_msk {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_neptune {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_qldb {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_route53resolver {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_states {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_transit_gateway {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_waf {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_wafv2 {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		cloudfront {
		  fetch_lambdas_at_edge    = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		dynamodb {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ec2 {
		  aws_regions              = ["us-east-1"]
		  duplicate_ec2_tags       = true
		  fetch_ip_addresses       = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ecs {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		efs {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticbeanstalk {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticsearch {
		  aws_regions              = ["us-east-1"]
		  fetch_nodes              = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elb {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		}
		emr {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		iam {
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		iot {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		kinesis {
		  aws_regions              = ["us-east-1"]
		  fetch_shards             = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		kinesis_firehose {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		lambda {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		rds {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		redshift {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		route53 {
		  fetch_extended_inventory = true
		  metrics_polling_interval = 6000
		}
		ses {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		sns {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  metrics_polling_interval = 6000
		}
		security_hub {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
	  }
`)
}

func testAccNewRelicAwsIntegrationsConfigUpdated(AWSIntegrationsTestConfig map[string]string) string {
	return fmt.Sprintf(`
	provider "newrelic" {
  		account_id = "` + AWSIntegrationsTestConfig["account_id"] + `"
  		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_aws_link_account" "foo" {
		provider        	   = newrelic.cloud-integration-provider
		arn                    = "` + AWSIntegrationsTestConfig["arn"] + `"
		metric_collection_mode = "PULL"
		name                   = "` + AWSIntegrationsTestConfig["name"] + `_updated"
		account_id			   = "` + AWSIntegrationsTestConfig["account_id"] + `"
	  }
	  resource "newrelic_cloud_aws_integrations" "bar" {
        provider        	   = newrelic.cloud-integration-provider
		linked_account_id      = newrelic_cloud_aws_link_account.foo.id
		account_id			   = "` + AWSIntegrationsTestConfig["account_id"] + `"
		billing {
		  metrics_polling_interval = 10000
		}
		cloudtrail {
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		}
		health {
		  metrics_polling_interval = 6000
		}
		trusted_advisor {
		  metrics_polling_interval = 6000
		}
		vpc {
		  aws_regions              = ["us-east-1"]
		  fetch_nat_gateway        = true
		  fetch_vpn                = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		x_ray {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		s3 {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		doc_db {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		sqs {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  queue_prefixes           = ["test prefix"]
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ebs {
		  metrics_polling_interval = 6000
		  fetch_extended_inventory = true
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		alb {
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  load_balancer_prefixes   = ["test prefix"]
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticache {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		api_gateway {
		  metrics_polling_interval = 6000
		  aws_regions              = ["us-east-1"]
		  stage_prefixes           = ["test prefix"]
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		auto_scaling {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_app_sync {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_athena {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_cognito {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_direct_connect {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_fsx {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_glue {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_kinesis_analytics {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_media_convert {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_media_package_vod {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_mq {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_msk {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_neptune {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_qldb {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_route53resolver {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_states {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_transit_gateway {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_waf {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		aws_wafv2 {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		cloudfront {
		  fetch_lambdas_at_edge    = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		dynamodb {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ec2 {
		  aws_regions              = ["us-east-1"]
		  duplicate_ec2_tags       = true
		  fetch_ip_addresses       = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		ecs {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		efs {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticbeanstalk {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elasticsearch {
		  aws_regions              = ["us-east-1"]
		  fetch_nodes              = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		elb {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		}
		emr {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		iam {
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		iot {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		kinesis {
		  aws_regions              = ["us-east-1"]
		  fetch_shards             = true
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		kinesis_firehose {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		lambda {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		rds {
		  aws_regions              = ["us-east-1"]
		  fetch_tags               = true
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		redshift {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		  tag_key                  = "test"
		  tag_value                = "test"
		}
		route53 {
		  fetch_extended_inventory = true
		  metrics_polling_interval = 6000
		}
		ses {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 6000
		}
		sns {
		  aws_regions              = ["us-east-1"]
		  fetch_extended_inventory = true
		  metrics_polling_interval = 6000
		}
		aws_auto_discovery {
			aws_regions              = ["us-east-1"]
			metrics_polling_interval = 6000
		}
		security_hub {
		  aws_regions              = ["us-east-1"]
		  metrics_polling_interval = 86400
		}
	  }
`)
}
