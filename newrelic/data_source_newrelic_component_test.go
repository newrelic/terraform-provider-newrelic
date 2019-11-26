package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicComponent_Basic(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicComponentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicComponentDataSource("data.newrelic_component.foo"),
				),
			},
		},
	})
}

func testAccCheckNewRelicComponentDataSource(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a component from New Relic")
		}

		return nil
	}
}

func testAccNewRelicComponentConfig() string {
	return `
data "newrelic_plugin" "foo" {
	guid = "net.kenjij.newrelic_redis_plugin"
}
data "newrelic_component" "foo" {
	plugin_id = "${data.newrelic_plugin.foo.id}"
	name = "MyRedisServer"
}
`
}
