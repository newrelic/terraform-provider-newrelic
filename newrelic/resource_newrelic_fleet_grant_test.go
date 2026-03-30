//go:build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetGrant_Basic(t *testing.T) {
	// Skip: Requires a valid group_id which varies per organization
	// In real usage, users would provide their actual group IDs from their organization
	t.Skip("Skipping: requires organization-specific group_id")

	resourceName := "newrelic_fleet_grant.foo"
	fleetName := fmt.Sprintf("tf-test-fleet-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetGrantConfig(fleetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "fleet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
				),
			},
		},
	})
}

func testAccNewRelicFleetGrantConfig(fleetName string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = "%s"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

data "newrelic_fleet_role" "manager" {
  name = "Fleet Manager"
  type = "STANDARD"
}

resource "newrelic_fleet_grant" "foo" {
  fleet_id = newrelic_fleet.test.id

  # Note: This test requires a valid group_id
  # In a real test, you would get this from a data source or environment variable
  grant {
    group_id = "test-group-id"
    role_id  = data.newrelic_fleet_role.manager.id
  }
}
`, fleetName)
}
