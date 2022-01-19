//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAwsLinkAccount_Basic(t *testing.T) {
	randName := acctest.RandString(5)
	resourceName := "newrelic_cloud_aws_link_account.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsLinkAccountDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicCloudAwsLinkAccountConfig(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsLinkAccountExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicCloudAwsLinkAccountConfigUpdated(randName),
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

func testAccNewRelicCloudAwsLinkAccountConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_link_account" "foo" {
  arn = "arn:aws:iam::837935485030:role/NewRelicInfrastructure-Integrations-DTK"
  metric_collection_mode = "PULL"
  name = "%s"
}
`, name)
}

func testAccNewRelicCloudAwsLinkAccountConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_link_account" "foo" {
  arn = "arn:aws:iam::837935485030:role/NewRelicInfrastructure-Integrations-DTK"
  metric_collection_mode = "PULL"
  name = "%s"
}
`, name)
}
