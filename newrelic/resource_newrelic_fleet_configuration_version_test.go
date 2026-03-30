//go:build integration

package newrelic

import (
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
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
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
