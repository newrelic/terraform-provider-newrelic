//go:build integration || CLOUD

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

func TestNewrelicAwsGovCloudLinkAccount_Basic(t *testing.T) {
	randName := acctest.RandString(5)

	resourceName := "resourceNewRelicAwsGovCloudLinkAccount.foo"

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
		CheckDestroy: testAccCheckNewRelicawsGovCloudLinkAccountDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicAwsGovCloudLinkAccountConfig(testAwsGovCloudAccessKeyId, testAwsGovCloudAwsAccountId, testAwsGovCloudSecretAccessKey, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsGovCloudLinkAccountExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccCheckNewRelicAwsGovCloudLinkAccountConfigUpdated(testAwsGovCloudAccessKeyId, testAwsGovCloudAwsAccountId, testAwsGovCloudSecretAccessKey, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsGovCloudLinkAccountExists(resourceName),
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

func testAccCheckNewRelicAwsGovCloudLinkAccountExists(n string) resource.TestCheckFunc {
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

		if err != nil && linkedAccount == nil {
			return err
		}

		return nil
	}
}

func testAccCheckNewRelicawsGovCloudLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_awsGovcloud_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked awsGovcloud account still exists: #{err}")
		}
	}

	return nil
}

func testAccCheckNewRelicAwsGovCloudLinkAccountConfig(access_key_id string, aws_account_id string, secret_access_key string, name string) string {
	return fmt.Sprintf(`
    resource "newrelic_cloud_awsGovcloud_link_account" "account" {
    access_key_id ="%[1]s"
	aws_account_id="%[2]s"
	metric_collection_mode = "PULL"
    name = "%[4]s"
	secret_access_key = "%[3]s"
}
`, access_key_id, aws_account_id, secret_access_key, name)
}

func testAccCheckNewRelicAwsGovCloudLinkAccountConfigUpdated(access_key_id string, aws_account_id string, secret_access_key string, name string) string {
	return fmt.Sprintf(`
    resource ""newrelic_cloud_awsGovcloud_link_account" "account" {
    access_key_id ="%[1]s"
	aws_account_id="%[2]s"
	metric_collection_mode = "PULL"
    name = "%[4]s-Updated"
	secret_access_key = "%[3]s"
}
`, access_key_id, aws_account_id, secret_access_key, name)
}
