//go:build integration || FLEET

package newrelic

import (
t"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetRoleDataSource_Basic(t *testing.T) {
	resourceName := "data.newrelic_fleet_role.manager"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
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
