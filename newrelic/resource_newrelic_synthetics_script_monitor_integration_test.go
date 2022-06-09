//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
)

func TestAccNewRelicSyntheticsScriptAPIMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_script_monitor.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		// PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptAPIMonitorConfig(rName, string(SyntheticsMonitorTypes.SCRIPT_API)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
				Destroy: false,
			},
		},
		// CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
	})
}

func testAccNewRelicSyntheticsScriptAPIMonitorConfig(name string, scriptMonitorType string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_script_monitor" "foo" {
	name                 = "%s"
	type                 = "%s"
	locations_public     = ["AP_SOUTH_1"]
	period               = "EVERY_HOUR"
	status               = "ENABLED"
	script               = "console.log('terraform integration test')"
	#script_language      = "javascript"
	#runtime_type         = "NODE_API"
	#runtime_type_version = "16.10.0"

	tags {
		key    = "someKey"
		values = ["somevalue"]
	}
}`, name, scriptMonitorType)
}

func testAccCheckNewRelicSyntheticsScriptMonitorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		fmt.Print("\n\n **************************** \n")
		fmt.Printf("\n rs.Primary.ID:  %+v \n", rs.Primary.ID)

		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			fmt.Printf(rs.Primary.ID)
			return err
		}

		fmt.Printf("\n FOUND:  %+v \n", *found)
		fmt.Print("\n **************************** \n\n")

		// if string((*found).GetGUID()) != rs.Primary.ID {
		// 	fmt.Errorf("the monitor is not found %v - %v", (*found).GetGUID(), rs.Primary.ID)
		// }

		return nil
	}
}