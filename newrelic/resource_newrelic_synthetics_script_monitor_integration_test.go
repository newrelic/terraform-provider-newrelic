//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
)

func TestAccNewRelicSyntheticsScriptAPIMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_script_monitor.foo"
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	monitorTypeStr := string(SyntheticsMonitorTypes.SCRIPT_API)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsScriptMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptAPIMonitorConfig(rName, monitorTypeStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicSyntheticsScriptAPIMonitorConfig(fmt.Sprintf("%s-updated", rName), monitorTypeStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"locations_public",
					"location_private",
					"status",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tag",
					"script",
					"enable_screenshot_on_failure_and_script",
				},
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptBrowserMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_script_monitor.bar"
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsScriptMonitorDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptBrowserMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsScriptBrowserMonitorConfig(fmt.Sprintf("%s-updated", rName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsScriptMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"locations_public",
					"location_private",
					"status",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tag",
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
			name	=	"%s"
			type	=	"%s"
			locations_public	=	["AP_SOUTH_1"]
			period	=	"EVERY_HOUR"
			status	=	"ENABLED"
			script	=	"console.log('terraform integration test')"
			script_language	=	"JAVASCRIPT"
			runtime_type	=	"NODE_API"
			runtime_type_version	=	"16.10"
			tag {
				key	=	"some_key"
				values	=	["some_value"]
			}
		}`, name, scriptMonitorType)
}

func testAccNewRelicSyntheticsScriptBrowserMonitorConfig(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "bar" {
			enable_screenshot_on_failure_and_script	=	true
			locations_public	=	["AP_SOUTH_1"]
			name	=	"%[1]s"
			period	=	"EVERY_HOUR"
			runtime_type_version	=	"100"
			runtime_type	=	"CHROME_BROWSER"
			script_language	=	"JAVASCRIPT"
			status	=	"ENABLED"
			type	=	"SCRIPT_BROWSER"
			script	=	"$browser.get('https://one.newrelic.com')"
			tag {
				key	= "Name"
				values	= ["scriptedMonitor"]
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

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(10 * time.Second)

		result, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}
		if string((*result).GetGUID()) != rs.Primary.ID {
			return fmt.Errorf("the monitor is not found %v - %v", (*result).GetGUID(), rs.Primary.ID)
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

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(10 * time.Second)

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
