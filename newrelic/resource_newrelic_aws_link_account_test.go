package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	}
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
