// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicSyntheticsMonitor_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMonitorConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMonitorConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicSyntheticsMonitor_OptionalArgs(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMonitorConfigOmitOptionalArgs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMonitorConfigUpdateIncludeOptionalArgs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicSyntheticsMonitor_Browser(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMonitorConfigBrowser(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMonitorConfigBrowserUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicSyntheticsMonitor_ScriptBrowser(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMonitorConfigScriptBrowser(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMonitorConfigScriptBrowserUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicSyntheticsMonitor_ScriptAPI(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsMonitorConfigScriptAPI(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsMonitorConfigScriptAPIUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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

		found, err := client.Synthetics.GetMonitor(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("synthetics monitor not found: %v - %v", rs.Primary.ID, found)
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

		_, err := client.Synthetics.GetMonitor(r.Primary.ID)
		if err == nil {
			return fmt.Errorf("synthetics monitor still exists")
		}

	}
	return nil
}

func testAccNewRelicSyntheticsMonitorConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name       = "%[1]s"
	type       = "SIMPLE"
	frequency  = 1
	status     = "DISABLED"
	locations  = ["AWS_US_EAST_1"]

	uri                       = "https://example.com"
	validation_string         = "add example validation check here"
	verify_ssl                = false
	bypass_head_request       = false
	treat_redirect_as_failure = false
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-updated"
	type      = "SIMPLE"
	frequency = 5
	status    = "ENABLED"
	locations = ["AWS_US_EAST_1", "AWS_US_WEST_1"]

	uri                       = "https://example-updated.com"
	validation_string         = "add example validation check here updated"
	verify_ssl                = true
	bypass_head_request       = true
	treat_redirect_as_failure = true
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigOmitOptionalArgs(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s"
	type      = "SIMPLE"
	frequency = 1
	status    = "ENABLED"
	locations = ["AWS_US_EAST_1", "AWS_US_WEST_1"]

	uri = "https://example.com"
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigUpdateIncludeOptionalArgs(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-updated"
	type      = "SIMPLE"
	frequency = 5
	status    = "ENABLED"
	locations = ["AWS_US_EAST_1", "AWS_US_WEST_1"]

	uri                       = "https://example-updated.com"
	validation_string         = "this should should exist in the config now"
	verify_ssl                = false
	bypass_head_request       = false
	treat_redirect_as_failure = false
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigBrowser(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-browser-test"
	type      = "BROWSER"
	frequency = 1
	status    = "DISABLED"
	locations = ["AWS_US_EAST_1"]

	uri               = "https://example.com"
	validation_string = "this text should exist in the response"
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigBrowserUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-browser-test-updated"
	type      = "BROWSER"
	frequency = 5
	status    = "ENABLED"
	locations = ["AWS_US_EAST_1", "AWS_US_WEST_1"]

	uri               = "https://example-updated.com"
	validation_string = "this text should exist in the response updated"
	verify_ssl        = false
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigScriptBrowser(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-script-browser-test"
	type      = "SCRIPT_BROWSER"
	frequency = 1
	status    = "DISABLED"
	locations = ["AWS_US_EAST_1"]
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigScriptBrowserUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-script-browser-test-updated"
	type      = "SCRIPT_BROWSER"
	frequency = 5
	status    = "ENABLED"
	locations = ["AWS_US_EAST_2"]
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigScriptAPI(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-script-api-test"
	type      = "SCRIPT_API"
	frequency = 1
	status    = "DISABLED"
	locations = ["AWS_US_EAST_1"]
}
`, name)
}

func testAccNewRelicSyntheticsMonitorConfigScriptAPIUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name      = "%[1]s-script-api-test-updated"
	type      = "SCRIPT_API"
	frequency = 5
	status    = "ENABLED"
	locations = ["AWS_US_EAST_2"]
}
`, name)
}
