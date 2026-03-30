//go:build integration || FLEET

package newrelic

import (
t"os"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetConfigurationVersion_Basic(t *testing.T) {
	resourceName := "newrelic_fleet_configuration_version.foo"
	rName := fmt.Sprintf("tf-test-config-%s", acctest.RandString(5))
	configContent1 := `{"log_level": "info"}`
	configContent2 := `{"log_level": "debug"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigurationVersionConfig(rName, configContent1, configContent2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "blob_id"),
				),
			},
		},
	})
}

func testAccNewRelicFleetConfigurationVersionConfig(name, configContent1, configContent2 string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "base" {
  name                   = "%s"
  agent_type             = "NRInfra"
  managed_entity_type    = "HOST"
  configuration_content  = %q
}

resource "newrelic_fleet_configuration_version" "foo" {
  configuration_id      = newrelic_fleet_configuration.base.id
  configuration_content = %q
}
`, name, configContent1, configContent2)
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
