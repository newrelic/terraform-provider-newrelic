//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAwsLinkAccount_Basic(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicAwsLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsLinkAccountConfig(rName),
			},
			//Test: Update
			//TODO
			{
				Config: testAccNewRelicAwsLinkAccountConfigUpdated(rName),
			},
		},
	})
}

func testAccNewRelicAwsLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_aws_link_account" {
			continue
		}
		resourceId, err := strconv.Atoi(r.Primary.ID)
		if err != nil {
			fmt.Errorf("unable to convert string to int")
		}
		_, err = client.Cloud.GetLinkedAccount(testAccountID, resourceId)
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccNewRelicAwsLinkAccountConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_aws_link_account" "account" {
		arn = ""
		metric_collection_mode = "push"
		name = "%[1]s"
	}
	`, name)
}

func testAccNewRelicAwsLinkAccountConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_aws_link_account" "account" {
		arn = ""
		metric_collection_mode = "push"
		name = "%[1]s-updated"
	}
	`, name)
}
