//go:build integration

package newrelic

import (
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
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
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
