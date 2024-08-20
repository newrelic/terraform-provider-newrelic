//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCurrentUserDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicCurrentUserDataSourceConfiguration(),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckCurrentUserDataSourceExists(t, "data.newrelic_current_user.foo"),
				),
			},
		},
	})
}

func testAccNewRelicCheckCurrentUserDataSourceExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get an ID of the matching user")
		}

		return nil
	}
}

func testAccNewRelicCurrentUserDataSourceConfiguration() string {
	return fmt.Sprintf(`
data "newrelic_current_user" "foo" {
}`)
}
