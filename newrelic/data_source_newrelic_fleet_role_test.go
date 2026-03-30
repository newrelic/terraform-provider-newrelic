//go:build integration

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetRoleDataSource_Basic(t *testing.T) {
	resourceName := "data.newrelic_fleet_role.manager"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetRoleDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Fleet Manager"),
					resource.TestCheckResourceAttr(resourceName, "type", "STANDARD"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "scope", "fleet"),
				),
			},
		},
	})
}

func testAccNewRelicFleetRoleDataSourceConfig() string {
	return `
data "newrelic_fleet_role" "manager" {
  name = "Fleet Manager"
  type = "STANDARD"
}
`
}
