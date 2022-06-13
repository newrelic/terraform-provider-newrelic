//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"testing"
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
				Config: testAccNewRelicSyntheticsSimpleMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsSimpleMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsSimpleMonitorConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "foo" {
	  custom_headers{
		name	= "Name"
		value	= "simpleMonitor"
		}
	  treat_redirect_as_failure	= true
	  validation_string			= "success"
	  bypass_head_request		= true
	  verify_ssl				= true
	  locations_public 			= ["AP_SOUTH_1"]
	  name      				= "%[1]s"
	  period 					= "EVERY_MINUTE"
	  status    				= "ENABLED"
	  type      				= "SIMPLE"
	  tags{
		key		= "monitor"
		values	= ["myMonitor"]
	  }
	  uri       				= "https://www.one.newrelic.com"
	}`, name)
}

func testAccNewRelicSyntheticsSimpleMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "foo" {
	  custom_headers{
		name	= "name"
		value	= "simpleMonitorUpdated"
	  }
	  treat_redirect_as_failure	= false
	  validation_string			= "succeeded"
	  bypass_head_request		= false
	  verify_ssl				= false
	  locations_public 			= ["AP_SOUTH_1","AP_EAST_1"]
	  name      				= "%[1]s-updated"
	  period 					= "EVERY_5_MINUTES"
	  status    				= "DISABLED"
	  type      				= "SIMPLE"
	  tags{
		key		= "monitor"
		values	=[ "myMonitor","simple_monitor"]
	  }
	  uri       				= "https://www.one.newrelic.com"
	}
`, name)
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
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfig(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  custom_headers{
			name	= "name"
			value	= "simple_browser"
		  }
		  enable_screenshot_on_failure_and_script	= true
		  validation_string							= "success"
		  verify_ssl								= true
		  locations_public 							= ["AP_SOUTH_1"]
		  name      								= "%[1]s"
		  period 									= "EVERY_MINUTE"
		  runtime_type_version						= "100"
		  runtime_type								= "CHROME_BROWSER"
		  script_language							= "JAVASCRIPT"
		  status    								= "ENABLED"
		  type      								= "BROWSER"
          uri										= "https://www.one.newrelic.com"
		  tags{
			key		= "name"
			values	= ["SimpleBrowserMonitor"]
		  }
		}
		`, name)
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  custom_headers{
			name  = "name"
			value = "simple_browser"
		  }
		  enable_screenshot_on_failure_and_script = false
		  validation_string                       = "success"
		  verify_ssl                              = false
		  locations_public                        = ["AP_SOUTH_1","AP_EAST_1"]
		  name                                    = "%[1]s-Updated"
		  period                                  = "EVERY_5_MINUTES"
		  runtime_type_version                    = "100"
		  runtime_type                            = "CHROME_BROWSER"
		  script_language                         = "JAVASCRIPT"
		  status                                  = "DISABLED"
		  type                                    = "BROWSER"
		  uri                                     = "https://www.one.newrelic.com"
		  tags{
			key     = "name"
			values  = ["SimpleBrowserMonitor","my_monitor"]
		  }
		}
		`, name)
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
