package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAwsLinkAccount_Basic(t *testing.T) {
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsAccountLinkAccountConfig(rName),
			},
			//Test: Update
			//TODO
			{
				Config: testAccNewRelicAwsAccountLinkAccountConfigUpdated(rName),
			},
		},
	})
}

func testAccNewRelicAwsAccountLinkAccountConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_aws_link_account" "account" {
		arn = ""
		metric_collection_mode = "push"
		name = "%[1]s"
	}
	`, name)
}

func testAccNewRelicAwsAccountLinkAccountConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_aws_link_account" "account" {
		arn = ""
		metric_collection_mode = "push"
		name = "%[1]s-updated"
	}
	`, name)
}
