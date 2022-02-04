//go:build integration
// +build integration

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
	// resourceName := "newrelic_cloud_awsGovCloud_link_account.foo"
	
	}

	resource.ParallelTest(t, resource.TestCase {
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicawsGovCloudLinkAccountDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccCheckNewRelicAwsGovCloudLinkAccountConfig(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsGovCloudLinkAccountExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccCheckNewRelicAwsGovCloudLinkAccountConfigUpdated(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsGovCloudLinkAccountExists(resourceName),
				),
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
			fmt.Errorf("error converting string id to int")
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
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked awsGovcloud account still exists: #{err}")
		}
	}

	return nil
}

func testAccCheckNewRelicawsGovCloudLinkAccountConfig(name string) string {
	return fmt.Sprintf(`
    resource "newrelic_cloud_awsGovcloud_link_account" "account" {
    access_key_id =""
	aws_account_id=""
	metric_collection_mode = "PULL"
    name = "%s"
	secret_access_key = ""
}
`, name)
}

func testAccCheckNewRelicawsGovCloudLinkAccountConfigUpdated(name string) string {
	return fmt.Sprintf(`
    resource ""newrelic_cloud_awsGovcloud_link_account" "account" {
    access_key_id =""
	aws_account_id=""
	metric_collection_mode = "PULL"
    name = "%s"
	secret_access_key = ""
}
`, name)
}
