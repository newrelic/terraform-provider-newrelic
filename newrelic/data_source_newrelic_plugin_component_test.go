// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicPluginComponent_Basic(t *testing.T) {
	t.Skip()

	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicPluginComponentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPluginComponentDataSource("data.newrelic_plugin_component.foo"),
				),
			},
		},
	})
}

func testAccCheckNewRelicPluginComponentDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a plugin component from New Relic")
		}

		return nil
	}
}

func testAccNewRelicPluginComponentConfig() string {
	return `
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_plugin_component" "foo" {
	plugin_id = data.newrelic_plugin.foo.id
	name = "MyRedisServer"
}
`
}
