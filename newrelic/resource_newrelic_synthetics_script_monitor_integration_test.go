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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsScriptMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptAPIMonitorConfig(rName, string(SyntheticsMonitorTypes.SCRIPT_API)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicSyntheticsScriptAPIMonitorConfigUpdated(rName, string(SyntheticsMonitorTypes.SCRIPT_API)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true, //name,type
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"locations_public",
					"locations_private",
					"status",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tags",
					"script",
					"enable_screenshot_on_failure_and_script",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsScriptAPIMonitorConfig(name string, scriptMonitorType string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "foo" {
			name					=	"%s"
			type					=	"%s"
			locations_public		=	["AP_SOUTH_1"]
			period					=	"EVERY_HOUR"
			status					=	"ENABLED"
			script					=	"console.log('terraform integration test')"
			script_language			=	"JAVASCRIPT"
			runtime_type			=	"NODE_API"
			runtime_type_version	=	"16.10"
			tags {
				key		=	"some_key"
				values	=	["some_value"]
			}
		}`, name, scriptMonitorType)
}

func testAccNewRelicSyntheticsScriptAPIMonitorConfigUpdated(name string, scriptMonitorType string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "foo" {
			name					=	"%s-updated"
			type					=	"%s"
			locations_public		=	["AP_SOUTH_1","AP_EAST_1"]
			period					=	"EVERY_6_HOURS"
			status					=	"DISABLED"
			script					=	"console.log('terraform integration test updated')"
			script_language			=	"JAVASCRIPT"
			runtime_type			=	"NODE_API"
			runtime_type_version	=	"16.10"
			tags {
				key		=	"some_key"
				values	=	["some_value","some_other_value"]
			}
		}`, name, scriptMonitorType)
}

func TestAccNewRelicSyntheticsScriptedBrowserMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_script_monitor.bar"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsScriptMonitorDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptedBrowserMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsScriptBrowserMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true, //name,type
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"locations_public",
					"locations_private",
					"status",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tags",
					"script",
					"enable_screenshot_on_failure_and_script",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsScriptedBrowserMonitorConfig(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "bar" {
			enable_screenshot_on_failure_and_script	=	true
			locations_public						=	["AP_SOUTH_1"]
			name									=	"%[1]s"
			period									=	"EVERY_HOUR"
			runtime_type_version					=	"100"
			runtime_type							=	"CHROME_BROWSER"
			script_language							=	"JAVASCRIPT"
			status									=	"ENABLED"
			type									=	"SCRIPT_BROWSER"
			script									=	"$browser.get('https://one.newrelic.com')"
			tags {
				key		= "Name"
				values	= ["scriptedMonitor"]
			}
		}`, name)
}

func testAccNewRelicSyntheticsScriptBrowserMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "bar" {
			enable_screenshot_on_failure_and_script	=	false
			locations_public						=	["AP_SOUTH_1","AP_EAST_1"]
			name									=	"%[1]s_updated"
			period									=	"EVERY_HOUR"
			runtime_type_version					=	"100"
			runtime_type							=	"CHROME_BROWSER"
			script_language							=	"JAVASCRIPT"
			status									=	"DISABLED"
			type									=	"SCRIPT_BROWSER"
			script									=	"$browser.get('https://one.newrelic.com')"
			tags {
				key		=	"Name"
				values	=	["scriptedMonitor","hello"]
			}
		}`, name)
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

		result, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if string((*result).GetGUID()) != rs.Primary.ID {
			fmt.Errorf("the monitor is not found %v - %v", (*result).GetGUID(), rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsScriptMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_script_monitor" {
			continue
		}

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))

		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
