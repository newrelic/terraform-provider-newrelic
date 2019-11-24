package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicPlugin_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicPluginConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicPlugin("data.newrelic_plugin.guid"),
				),
			},
		},
	})
}

func testAccNewRelicPlugin(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a plugin from New Relic")
		}

		return nil
	}
}

// The test plugin for this data source is created in provider_test.go
func testAccNewRelicPluginConfig() string {
	return fmt.Sprintf(`
data "newrelic_plugin" "foo" {
	guid = "%s"
}
`, testAccExpectedPluginName)
}
