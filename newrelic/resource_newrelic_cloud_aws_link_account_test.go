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
	rName := acctest.RandString(5)
	resourceName := "newrelic_cloud_aws_link_account.foo"
	testAwsArn := os.Getenv("INTEGRATION_TESTING_AWS_ARN")

	if testAwsArn == "" {
		t.Skipf("INTEGRATION_TESTING_AWS_ARN must be set for this acceptance test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsLinkAccountConfig(rName, testAwsArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsLinkAccountExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicAwsLinkAccountConfigUpdated(rName, testAwsArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsLinkAccountExists(resourceName),
				),
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
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

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
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked aws account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicAwsLinkAccountConfig(name string, arn string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_aws_link_account" "foo" {
		arn = "%[2]s"
		metric_collection_mode = "PUSH"
		name = "%[1]s"
	}
	`, name, arn)
}

func testAccNewRelicAwsLinkAccountConfigUpdated(name string, arn string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_aws_link_account" "foo" {
		arn = "%[2]s"
		metric_collection_mode = "PUSH"
		name = "%[1]s-updated"
	}
	`, name, arn)
}
