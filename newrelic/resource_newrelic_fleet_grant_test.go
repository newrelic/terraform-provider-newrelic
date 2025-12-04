//go:build integration || FLEET
// +build integration FLEET

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetGrant(t *testing.T) {
	t.Parallel()

	// Note: This test assumes the existence of a fleet, a group, and a role.
	// Replace these with actual values from your test environment.
	// It is recommended to create these dependencies as resources in the test configuration.
	const fleetID = "your-fleet-id"
	const groupID1 = "your-group-id-1"
	const roleID1 = 12345
	const groupID2 = "your-group-id-2"
	const roleID2 = 54321

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetGrantConfig(fleetID, groupID1, roleID1, groupID2, roleID2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_fleet_grant.foo", "fleet_id", fleetID),
					resource.TestCheckResourceAttr("newrelic_fleet_grant.foo", "grant.#", "2"),
				),
			},
		},
	})
}

func testAccNewRelicFleetGrantConfig(fleetID, groupID1 string, roleID1 int, groupID2 string, roleID2 int) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_grant" "foo" {
  fleet_id = "%s"

  grant {
    group_id = "%s"
    role_id  = %d
  }

  grant {
    group_id = "%s"
    role_id  = %d
  }
}
`, fleetID, groupID1, roleID1, groupID2, roleID2)
}
