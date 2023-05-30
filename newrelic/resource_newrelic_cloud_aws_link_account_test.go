//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicCloudAwsLinkAccount_Basic(t *testing.T) {
	testAWSLinkAccountName := fmt.Sprintf("tf_cloud_link_account_test_aws_%s", acctest.RandString(5))
	resourceName := "newrelic_cloud_aws_link_account.foo"

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testAWSArn := os.Getenv("INTEGRATION_TESTING_AWS_ARN")
	if testAWSArn == "" {
		t.Skipf("INTEGRATION_TESTING_AWS_ARN must be set for this acceptance test")
	}

	AWSLinkAccountTestConfig := map[string]string{
		"name":       testAWSLinkAccountName,
		"account_id": strconv.Itoa(testSubAccountID),
		"arn":        testAWSArn,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "aws") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsLinkAccountConfig(AWSLinkAccountTestConfig, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsLinkAccountExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicAwsLinkAccountConfig(AWSLinkAccountTestConfig, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsLinkAccountExists(resourceName),
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

func testAccCheckNewRelicCloudAwsLinkAccountExists(n string) resource.TestCheckFunc {
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

		if err != nil && linkedAccount == nil {
			return err
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAwsLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_aws_link_account" {
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

func testAccNewRelicAwsLinkAccountConfig(AWSLinkAccountTestConfig map[string]string, updated bool) string {
	if updated == true {
		AWSLinkAccountTestConfig["name"] += "_updated"
	}

	return fmt.Sprintf(`
	provider "newrelic" {
		account_id = "` + AWSLinkAccountTestConfig["account_id"] + `"
		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_aws_link_account" "foo" {
		provider        	   = newrelic.cloud-integration-provider
		arn                    = "` + AWSLinkAccountTestConfig["arn"] + `"
		metric_collection_mode = "PULL"
		name                   = "` + AWSLinkAccountTestConfig["name"] + `"
		account_id			   = "` + AWSLinkAccountTestConfig["account_id"] + `"
	}`)
}
