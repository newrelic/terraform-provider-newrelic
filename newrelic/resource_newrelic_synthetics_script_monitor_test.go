//go:build integration || SYNTHETICS

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func TestAccNewRelicSyntheticsScriptAPIMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_script_monitor.foo"
	rName := generateNameForIntegrationTestResource()
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
			// Test: Update
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
					// Technical limitations with the API prevent us from setting the following attributes.
					"locations_public",
					"location_private",
					"tag",
					"script",
					"enable_screenshot_on_failure_and_script",
					SyntheticsUseLegacyRuntimeAttrLabel,
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
			// Test: Create
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
					// Technical limitations with the API prevent us from setting the following attributes.
					"locations_public",
					"location_private",
					"tag",
					"script",
					"enable_screenshot_on_failure_and_script",
					"device_orientation",
					"device_type",
					SyntheticsUseLegacyRuntimeAttrLabel,
					"browsers",
					"devices",
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
			locations_public	=	["AP_SOUTH_1", "US_EAST_1"]
			name	=	"%[1]s"
			period	=	"EVERY_HOUR"
			status	=	"ENABLED"
			type	=	"SCRIPT_BROWSER"

			runtime_type_version	=	"100"
			runtime_type	=	"CHROME_BROWSER"
			script_language	=	"JAVASCRIPT"
			script	=	"$browser.get('https://one.newrelic.com')"

			browsers = ["CHROME", "FIREFOX"]
			devices = ["DESKTOP", "TABLET_LANDSCAPE", "MOBILE_PORTRAIT", "MOBILE_LANDSCAPE", "TABLET_PORTRAIT"]
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

		if rs.Primary.Attributes["monitor_id"] != string((*result).(*entities.SyntheticMonitorEntity).MonitorId) {
			return fmt.Errorf("the monitor id doesnot match, expected: %v", rs.Primary.Attributes["monitor_id"])
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorScriptUpdate(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		accountId := providerConfig.AccountID

		result, err := client.Synthetics.GetScript(accountId, synthetics.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}
		if len(result.Text) == 0 {
			return fmt.Errorf("Synthetic Monitor Script update not successful !!")
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
