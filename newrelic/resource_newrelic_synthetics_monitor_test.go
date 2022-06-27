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

///////////////////////
//simple monitor test//
//////////////////////

func TestAccNewRelicSyntheticsSimpleMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsSimpleMonitorConfig(rName, string(SyntheticsMonitorTypes.SIMPLE)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsSimpleMonitorConfigUpdated(rName, string(SyntheticsMonitorTypes.SIMPLE)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true, //name,type,uri
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"location_public",
					"location_private",
					"status",
					"validation_string",
					"verify_ssl",
					"bypass_head_request",
					"treat_redirect_as_failure",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tag",
					"enable_screenshot_on_failure_and_script",
					"custom_header",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsSimpleMonitorConfig(name string, monitorType string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "foo" {
			custom_header {
				name	=	"Name"
				value	=	"simpleMonitor"
				}
			treat_redirect_as_failure	=	false
			validation_string	=	"success"
			bypass_head_request	=	false
			verify_ssl	=	false
			location_public	=	["AP_SOUTH_1"]
			name	=	"%s"
			period	=	"EVERY_MINUTE"
			status	=	"ENABLED"
			type	=	"%s"
			tag {
				key	=	"Name"
				values	=	["apple"]
			}
			uri	=	"https://www.one.newrelic.com"
		}`, name, monitorType)
}

func testAccNewRelicSyntheticsSimpleMonitorConfigUpdated(name string, monitorType string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "foo" {
			custom_header{
				name	=	"name"
				value	=	"simpleMonitorUpdated"
			}
			treat_redirect_as_failure	=	true
			validation_string	=	"succeeded"
			bypass_head_request	=	true
			verify_ssl	=	true
			location_public	=	["AP_SOUTH_1","AP_EAST_1"]
			name	=	"%s-updated"
			period	=	"EVERY_5_MINUTES"
			status	=	"DISABLED"
			type	=	"%s"
			tag {
				key	=	"Name"
				values	=	[ "pineApple","fruit"]
			}
			uri	=	"https://www.one.newrelic.com"
		}`, name, monitorType)
}

///////////////////////////////
//simple browser monitor test//
//////////////////////////////

func TestAccNewRelicSyntheticsSimpleBrowserMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.bar"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitorConfig(rName, string(SyntheticsMonitorTypes.BROWSER)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitorConfigUpdated(rName, string(SyntheticsMonitorTypes.BROWSER)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true, //name,type,uri
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"period",
					"location_public",
					"location_private",
					"status",
					"validation_string",
					"verify_ssl",
					"bypass_head_request",
					"treat_redirect_as_failure",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tag",
					"enable_screenshot_on_failure_and_script",
					"custom_header",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfig(name string, monitorType string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "bar" {
		custom_header{
			name	= "name"
			value	= "simple_browser"
		}
		enable_screenshot_on_failure_and_script	=	true
		validation_string	=	"success"
		verify_ssl	=	true
		location_public	=	["AP_SOUTH_1"]
		name	=	"%s"
		period	=	"EVERY_MINUTE"
		runtime_type_version	=	"100"
		runtime_type	=	"CHROME_BROWSER"
		script_language	=	"JAVASCRIPT"
		status	=	"ENABLED"
		type	=	"%s"
		uri	=	"https://www.one.newrelic.com"
		tag {
			key	=	"name"
			values	=	["SimpleBrowserMonitor"]
		}
	}`, name, monitorType)
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfigUpdated(name string, monitorType string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
			custom_header{
				name  = "name"
				value = "simple_browser"
			}
			enable_screenshot_on_failure_and_script	=	false
			validation_string	=	"success"
			verify_ssl	=	false
			location_public	=	["AP_SOUTH_1","AP_EAST_1"]
			name	=	"%s-Updated"
			period	=	"EVERY_5_MINUTES"
			runtime_type_version	=	"100"
			runtime_type	=	"CHROME_BROWSER"
			script_language	=	"JAVASCRIPT"
			status	=	"DISABLED"
			type	=	"%s"
			uri	=	"https://www.one.newrelic.com"
			tag {
				key	=	"name"
				values	=	["SimpleBrowserMonitor","my_monitor"]
		  	}
		}`, name, monitorType)
}

func testAccCheckNewRelicSyntheticsMonitorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if string((*found).GetGUID()) != rs.Primary.ID {
			fmt.Errorf("the monitor is not found %v - %v", (*found).GetGUID(), rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_monitor" {
			continue
		}

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
