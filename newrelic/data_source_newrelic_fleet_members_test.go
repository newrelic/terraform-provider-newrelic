//go:build integration || FLEET

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNewRelicFleetMembersDataSource_All reads all members of a fleet
// (no ring filter) and checks the data source populates the members list.
func TestAccNewRelicFleetMembersDataSource_All(t *testing.T) {
	dsName := "data.newrelic_fleet_members.all"
	fleetName := fmt.Sprintf("tf-test-ds-members-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)
	entityIDs := testAccFleetMembersEntityIDs(t, 1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetMembersDataSourceAll(fleetName, entityIDs[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "fleet_id"),
					resource.TestCheckResourceAttrSet(dsName, "members.#"),
				),
			},
		},
	})
}

// TestAccNewRelicFleetMembersDataSource_ByRing reads members filtered by ring.
func TestAccNewRelicFleetMembersDataSource_ByRing(t *testing.T) {
	dsName := "data.newrelic_fleet_members.by_ring"
	fleetName := fmt.Sprintf("tf-test-ds-members-ring-%s", acctest.RandString(5))

	setupFleetTestCredentials(t)
	entityIDs := testAccFleetMembersEntityIDs(t, 1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFleetEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicFleetMembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFleetMembersDataSourceByRing(fleetName, entityIDs[0]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "fleet_id"),
					resource.TestCheckResourceAttr(dsName, "ring", "default"),
				),
			},
		},
	})
}

// Config templates

func testAccFleetMembersDataSourceAll(fleetName, entityID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

resource "newrelic_fleet_members" "default" {
  fleet_id = newrelic_fleet.test.id
  ring {
    name       = "default"
    entity_ids = [%q]
  }
}

data "newrelic_fleet_members" "all" {
  fleet_id   = newrelic_fleet.test.id
  depends_on = [newrelic_fleet_members.default]
}
`, fleetName, entityID)
}

func testAccFleetMembersDataSourceByRing(fleetName, entityID string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet" "test" {
  name                = %q
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
}

resource "newrelic_fleet_members" "default" {
  fleet_id = newrelic_fleet.test.id
  ring {
    name       = "default"
    entity_ids = [%q]
  }
}

data "newrelic_fleet_members" "by_ring" {
  fleet_id   = newrelic_fleet.test.id
  ring       = "default"
  depends_on = [newrelic_fleet_members.default]
}
`, fleetName, entityID)
}
