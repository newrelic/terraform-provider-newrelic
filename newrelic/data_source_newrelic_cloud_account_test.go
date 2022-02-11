//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAccountDataSource_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_account.account"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicCloudAccountDataSourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAccountDataSource("data.newrelic_cloud_account.account"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("NEW-DTK-NAME"))),
			},
		},
	})

}

func testNewRelicCloudAccountDataSourceBasicConfig() string {
	return fmt.Sprintf(`
data "newrelic_cloud_account" "account" {
	account_id = 2508259
	name = "NEW-DTK-NAME"
	provider = "aws"
}
`)
}

func testNewRelicCloudAccountDataSourceErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_cloud_account" "account" {
	name = "NEW-DTK-NAME"
	provider = "aws"
}
`)
}

func testAccCheckNewRelicCloudAccountDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a cloud account from New Relic")
		}

		return nil
	}
}
