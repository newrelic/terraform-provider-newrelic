//go:build integration || FLEET

package newrelic

import (
t"os"
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
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
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

func setupFleetTestCredentials(t *testing.T) {
	t.Helper()

	// Set fleet credentials for this test
	originalAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	originalAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	t.Cleanup(func() {
		os.Setenv("NEW_RELIC_API_KEY", originalAPIKey)
		os.Setenv("NEW_RELIC_ACCOUNT_ID", originalAccountID)
	})

	fleetAPIKey := os.Getenv("NEW_RELIC_FLEET_TEST_API_KEY")
	fleetAccountID := os.Getenv("NEW_RELIC_FLEET_TEST_ACCOUNT_ID")
	if fleetAPIKey != "" {
		os.Setenv("NEW_RELIC_API_KEY", fleetAPIKey)
	}
	if fleetAccountID != "" {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", fleetAccountID)
	}
}
