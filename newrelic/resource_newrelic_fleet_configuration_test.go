//go:build integration || FLEET

package newrelic

import (
t"os"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetConfiguration_Basic(t *testing.T) {
	resourceName := "newrelic_fleet_configuration.foo"
	rName := fmt.Sprintf("tf-test-config-%s", acctest.RandString(5))
	configContent := `{"log_level": "info"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckFleetEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicFleetConfigurationConfig(rName, configContent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "agent_type", "NRInfra"),
					resource.TestCheckResourceAttr(resourceName, "managed_entity_type", "HOST"),
					resource.TestCheckResourceAttrSet(resourceName, "entity_guid"),
				),
			},
		},
	})
}

func testAccNewRelicFleetConfigurationConfig(name, content string) string {
	return fmt.Sprintf(`
resource "newrelic_fleet_configuration" "foo" {
  name                   = "%s"
  agent_type             = "NRInfra"
  managed_entity_type    = "HOST"
  configuration_content  = %q
}
`, name, content)
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
