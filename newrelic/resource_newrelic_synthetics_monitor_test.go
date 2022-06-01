//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		name="Name"
		value="simpleMonitor"
		}
	  treat_redirect_as_failure=true
	  validation_string="success"
	  bypass_head_request=true
	  verify_ssl=true
	  locations = ["AP_SOUTH_1"]
	  name      = "%[1]s"
	  frequency = 5
	  status    = "ENABLED"
	  type      = "SIMPLE"
	  tags{
		key="monitor"
		values=["myMonitor"]
	  }
	  uri       = "https://www.one.newrelic.com"
	}`, name)
}

func testAccNewRelicSyntheticsSimpleMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "foo" {
	  custom_headers{
		name="name"
		value="simpleMonitorUpdated"
	  }
	  treat_redirect_as_failure=false
	  validation_string="succeeded"
	  bypass_head_request=false
	  verify_ssl=false
	  locations = ["AP_SOUTH_1","AP_EAST_1"]
	  name      = "%[1]s-updated"
	  frequency = 10
	  status    = "DISABLED"
	  type      = "SIMPLE"
	  tags{
		key="monitor"
		values=["myMonitor","simple_monitor"]
	  }
	  uri       = "https://www.one.newrelic.com"
	}
`, name)
}

///////////////////////////////
//simple browser monitor test//
//////////////////////////////

func TestAccNewRelicSyntheticsSimpleBrowserMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
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
			name="name"
			value="simple_browser"
		  }
		  enable_screenshot_on_failure_and_script=true
		  validation_string="success"
		  verify_ssl=true
		  locations = ["AP_SOUTH_1"]
		  name      = "%[1]s"
		  frequency = 5
		  runtime_type_version="100"
		  runtime_type="CHROME_BROWSER"
		  script_language="JAVASCRIPT"
		  status    = "ENABLED"
		  type      = "BROWSER"
		  tags{
			key="name"
			values=["SimpleBrowserMonitor"]
		  }
		  uri="https://www.one.newrelic.com"
		}
		`, name)
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  custom_headers{
			name="name"
			value="simple_browser"
		  }
		  enable_screenshot_on_failure_and_script=false
		  validation_string="success"
		  verify_ssl=false
		  locations = ["AP_SOUTH_1","AP_EAST_1"]
		  name      = "%[1]s-Updated"
		  frequency = 10
		  runtime_type_version="100"
		  runtime_type="CHROME_BROWSER"
		  script_language="JAVASCRIPT"
		  status    = "DISABLED"
		  type      = "BROWSER"
		  tags{
			key="name"
			values=["SimpleBrowserMonitor","my_monitor"]
		  }
		  uri="https://www.one.newrelic.com"
		}
		`, name)
}

////////////////////////
//scripted-api-monitor//
////////////////////////

func TestAccNewRelicSyntheticsScriptedAPIMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptedAPIMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsScriptedAPIMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsScriptedAPIMonitorConfig(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  locations = ["AP_SOUTH_1"]
		  name      = "%[1]s"
		  frequency = 5
		  runtime_type="NODE_API"
		  script_language="JAVASCRIPT"
		  runtime_type_version="16.10"
		  script=<<EOF
		  var assert = require('assert');
		  $http.post('http://httpbin.org/post',
		  // Post data
		  {
		  json: {
		  widgetType: 'gear',
		  widgetCount: 10
		  }
		  },
		  // Callback
		  function (err, response, body) {
		  assert.equal(response.statusCode, 200, 'Expected a 200 OK response');
		  console.log('Response:', body.json);
		  assert.equal(body.json.widgetType, 'gear', 'Expected a gear widget type');
		  assert.equal(body.json.widgetCount, 10, 'Expected 10 widgets');
		  }
		  );
		  EOF
		  status    = "ENABLED"
		  type      = "SCRIPT_API"
		  tags{
			key="Name"
			values=["SimpleMonitor"]
		  }
		}
		`, name)
}

func testAccNewRelicSyntheticsScriptedAPIMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  locations = ["AP_SOUTH_1","AP_EAST_1"]
		  name      = "%[1]s-updated"
		  frequency = 10
		  runtime_type="NODE_API"
		  script_language="JAVASCRIPT"
		  runtime_type_version="16.10"
		  script=<<EOF
		  var assert = require('assert');
		  $http.post('http://httpbin.org/post',
		  // Post data
		  {
		  json: {
		  widgetType: 'gear',
		  widgetCount: 10
		  }
		  },
		  // Callback
		  function (err, response, body) {
		  assert.equal(response.statusCode, 200, 'Expected a 200 OK response');
		  console.log('Response:', body.json);
		  assert.equal(body.json.widgetType, 'gear', 'Expected a gear widget type');
		  assert.equal(body.json.widgetCount, 10, 'Expected 10 widgets');
		  }
		  );
		  EOF
		  status    = "DISABLED"
		  type      = "SCRIPT_API"
		  tags{
			key="Name"
			values=["SimpleMonitor","hello"]
		  }
		}
		`, name)
}

////////////////////////////
//Scripted-Browser-monitor//
////////////////////////////

func TestAccNewRelicSyntheticsScriptedBrowserMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicSyntheticsScriptedBrowserMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicSyntheticsScriptedBrowserMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsScriptedBrowserMonitorConfig(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  enable_screenshot_on_failure_and_script=true
		  locations = ["AP_SOUTH_1"]
		  name      = "%[1]s"
		  frequency = 10
		  runtime_type_version="100"
		  runtime_type="CHROME_BROWSER"
		  script_language="JAVASCRIPT"
		  status    = "ENABLED"
		  type      = "SCRIPT_BROWSER"
		  script=<<EOF
		  var assert = require('assert');
		  $browser.get('https://one.newrelic.com')
		  EOF
		  tags{
			key="Name"
			values=["scriptedMonitor"]
		  }
		}
		`, name)
}

func testAccNewRelicSyntheticsScriptedBrowserMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
		  enable_screenshot_on_failure_and_script=false
		  locations = ["AP_SOUTH_1","AP_EAST_1"]
		  name      = "%[1]s_updated"
		  frequency = 10
		  runtime_type_version="100"
		  runtime_type="CHROME_BROWSER"
		  script_language="JAVASCRIPT"
		  status    = "DISABLED"
		  type      = "SCRIPT_BROWSER"
		  script=<<EOF
		  var assert = require('assert');
		  $browser.get('https://one.newrelic.com')
		  EOF
		  tags{
			key="Name"
			values=["scriptedMonitor","hello"]
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

		_, err := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if err == nil {
			return fmt.Errorf("synthetics monitor still exists")
		}

	}
	return nil
}
