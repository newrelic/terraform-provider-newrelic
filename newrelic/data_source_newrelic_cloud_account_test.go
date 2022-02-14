//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAccountDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicCloudAccountDataSourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAccountDataSourceExists("data.newrelic_cloud_account.account")),
			},
		},
	})
}

func TestAccNewRelicCloudAccountDataSource_Error(t *testing.T) {
	expectedErrorMsg := regexp.MustCompile("does not match any account for provider")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testNewRelicCloudAccountDataSourceErrorConfig(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testNewRelicCloudAccountDataSourceBasicConfig() string {
	return fmt.Sprintf(`
data "newrelic_cloud_account" "account" {
	account_id = 2508259
	name = "NEW-DTK-NAME"
	cloud_provider = "aws"
}
`)
}

func testNewRelicCloudAccountDataSourceErrorConfig() string {
	return fmt.Sprintf(`
data "newrelic_cloud_account" "account" {
	name = "NEW-DTK-NAME"
	cloud_provider = "aws"
}
`)
}

func testAccCheckNewRelicCloudAccountDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		id := r.Primary.ID
		a := r.Primary.Attributes

		if id == "" {
			return fmt.Errorf("expected to get an account from New Relic")
		}

		if a["name"] == "" {
			return fmt.Errorf("expected to get a name from New Relic")
		}

		return nil
	}
}
