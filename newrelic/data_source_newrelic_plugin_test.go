package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicPlugin_Basic(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicPluginConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPluginDataSource("data.newrelic_plugin.foo"),
				),
			},
		},
	})
}

func testAccCheckNewRelicPluginDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a plugin from New Relic")
		}

		return nil
	}
}

func testAccNewRelicPluginConfig() string {
	return `
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
`
}
