//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccNewRelicCloudAwsGovCloudIntegrations_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_aws_govcloud_integrations.foo"

	randName := acctest.RandString(5)

	testAwsGovCloudAccessKeyId := os.Getenv("INTEGRATION_TESTING_AWSGOVCLOUD_ACCESS_KEY_ID")
	if testAwsGovCloudAccessKeyId == "" {
		t.Skipf("INTEGRATION_TESTING_AWSGOVCLOUD_ACCESS_KEY_ID must be set for acceptance test")
	}

	testAwsGovCloudAwsAccountId := os.Getenv("INTEGRATION_TESTING_AWSGOVCLOUD_AWS_ACCOUNT_ID")
	if testAwsGovCloudAwsAccountId == "" {
		t.Skipf("INTEGRATION_TESTING_AWSGOVCLOUD_AWS_ACCOUNT_ID must be set for acceptance test")
	}

	testAwsGovCloudSecretAccessKey := os.Getenv("INTEGRATION_TESTING_AWSGOVCLOUD_SECRET_ACCESS_KEY")
	if testAwsGovCloudSecretAccessKey == "" {
		t.Skipf("INTEGRATION_TESTING_AWSGOVCLOUD_SECRET_ACCESS_KEY must be set for acceptance test")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudAwsGovCloudIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudAwsGovCloudIntegrationsConfig(testAwsGovCloudAccessKeyId, testAwsGovCloudAwsAccountId, testAwsGovCloudSecretAccessKey, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudAwsGovCloudIntegrationsExist(resourceName)),
			},
			//Test: Update
			{
				Config: testAccNewRelicCloudAwsGovCloudIntegrationsConfig(testAwsGovCloudAccessKeyId, testAwsGovCloudAwsAccountId, testAwsGovCloudSecretAccessKey, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudAwsGovCloudIntegrationsExist(resourceName)),
			},
			//Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNewRelicCloudAwsGovCloudIntegrationsExist(n string) resource.TestCheckFunc {
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
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			fmt.Errorf("An error occured creating awsGovCloud integrations")
		}

		return nil
	}
}

func testAccNewRelicCloudAwsGovCloudIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_aws_govcloud_integrations" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked awsGovCloud account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicCloudAwsGovCloudIntegrationsConfig(access_key_id string, aws_account_id string, secret_access_key string, name string) string {
	return fmt.Sprintf(`
    resource "newrelic_cloud_awsGovcloud_link_account" "account" {
    access_key_id ="%[1]s"
	aws_account_id="%[2]s"
	metric_collection_mode = "PULL"
    name = "%[4]s"
	secret_access_key = "%[3]s"
}
	resource "newrelic_cloud_aws_govcloud_integrations" "foo" {
	  account_id=2520528
	  linked_account_id=newrelic_cloud_awsGovcloud_link_account.account.id
	  alb{
		metrics_polling_interval=1000
		aws_regions=["us-east-1"]
		fetch_extended_inventory=true
		fetch_tags=true
		load_balancer_prefixes=[""]
		tag_key=""
		tag_value=""
	  }
	  api_gateway{
		metrics_polling_interval=1000
		aws_regions=[""]
		stage_prefixes=[""]
		tag_key=""
		tag_value=""
	  }
	  auto_scaling{
		metrics_polling_interval=1000
		aws_regions=[""]
	  }
	  aws_direct_connect{
		metrics_polling_interval=1000
		aws_regions=[""]
	  }
	  aws_states{
		metrics_polling_interval=1000
		aws_regions=[""]
	  }
	  cloudtrail{
		metrics_polling_interval=1000
		aws_regions=[""]
	  }
	  dynamo_db{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_extended_inventory=true
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  ebs{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_extended_inventory=true
		tag_key=""
		tag_value=""
	  }
	  ec2{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_ip_addresses=true
		tag_key=""
		tag_value=""
	  }
	  elastic_search{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_nodes=true
		tag_key=""
		tag_value=""
	  }
	  elb{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_extended_inventory=true
		fetch_tags=true
	  }
	  emr{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  iam{
		metrics_polling_interval=1000
		tag_key=""
		tag_value=""
	  }
	  lambda{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  rds{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  red_shift{
		metrics_polling_interval=1000
		aws_regions=[""]
		tag_key=""
		tag_value=""
	  }
	  route53{
		metrics_polling_interval=1000
		fetch_extended_inventory=true
	  }
	  s3{
		metrics_polling_interval=1000
		fetch_extended_inventory=true
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  sns{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_extended_inventory=true
	  }
	  sqs{
		metrics_polling_interval=1000
		aws_regions=[""]
		fetch_extended_inventory=true
		fetch_tags=true
		queue_prefixes=[""]
		tag_key=""
		tag_value=""
	  }
}
`, access_key_id, aws_account_id, secret_access_key, name)
}

func testAccNewRelicCloudAwsGovCloudIntegrationsConfigUpdated(access_key_id string, aws_account_id string, secret_access_key string, name string) string {
	return fmt.Sprintf(`
    resource ""newrelic_cloud_awsGovcloud_link_account" "account" {
    access_key_id ="%[1]s"
	aws_account_id="%[2]s"
	metric_collection_mode = "PULL"
    name = "%[4]s-Updated"
	secret_access_key = "%[3]s"
}
	resource "newrelic_cloud_aws_govcloud_integrations" "foo" {
	  account_id=2520528
	  linked_account_id=newrelic_cloud_awsGovcloud_link_account.account.id
	  alb{
		metrics_polling_interval=2000
		aws_regions=["us-east-1"]
		fetch_extended_inventory=true
		fetch_tags=true
		load_balancer_prefixes=[""]
		tag_key=""
		tag_value=""
	  }
	  api_gateway{
		metrics_polling_interval=2000
		aws_regions=[""]
		stage_prefixes=[""]
		tag_key=""
		tag_value=""
	  }
	  auto_scaling{
		metrics_polling_interval=2000
		aws_regions=[""]
	  }
	  aws_direct_connect{
		metrics_polling_interval=2000
		aws_regions=[""]
	  }
	  aws_states{
		metrics_polling_interval=2000
		aws_regions=[""]
	  }
	  cloudtrail{
		metrics_polling_interval=2000
		aws_regions=[""]
	  }
	  dynamo_db{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_extended_inventory=true
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  ebs{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_extended_inventory=true
		tag_key=""
		tag_value=""
	  }
	  ec2{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_ip_addresses=true
		tag_key=""
		tag_value=""
	  }
	  elastic_search{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_nodes=true
		tag_key=""
		tag_value=""
	  }
	  elb{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_extended_inventory=true
		fetch_tags=true
	  }
	  emr{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  iam{
		metrics_polling_interval=2000
		tag_key=""
		tag_value=""
	  }
	  lambda{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  rds{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  red_shift{
		metrics_polling_interval=2000
		aws_regions=[""]
		tag_key=""
		tag_value=""
	  }
	  route53{
		metrics_polling_interval=2000
		fetch_extended_inventory=true
	  }
	  s3{
		metrics_polling_interval=2000
		fetch_extended_inventory=true
		fetch_tags=true
		tag_key=""
		tag_value=""
	  }
	  sns{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_extended_inventory=true
	  }
	  sqs{
		metrics_polling_interval=2000
		aws_regions=[""]
		fetch_extended_inventory=true
		fetch_tags=true
		queue_prefixes=[""]
		tag_key=""
		tag_value=""
	  }
}
`, access_key_id, aws_account_id, secret_access_key, name)
}
