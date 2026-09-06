//go:build integration || SYNTHETICS

package newrelic

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

var tv bool = true

func TestAccNewRelicSyntheticsBrowserMonitor_DeviceEmulationError(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config:      testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulationError(rName, string(SyntheticsMonitorTypes.BROWSER)),
				ExpectError: regexp.MustCompile("all of `device_orientation,device_type` must be\nspecified"),
			},
		},
	})
}

func TestAccNewRelicSyntheticsBrowserMonitor_DeviceEmulationErrorUpdate(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	resourceName := "newrelic_synthetics_monitor.foo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulation(rName, string(SyntheticsMonitorTypes.BROWSER)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
			// Test: Update ; By removing device_type field
			{
				Config:      testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulationError(rName, string(SyntheticsMonitorTypes.BROWSER)),
				ExpectError: regexp.MustCompile("all of `device_orientation,device_type` must be\nspecified"),
			},
			// Test: Update ; Added back removed device_type field
			{
				Config: testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulation(rName, string(SyntheticsMonitorTypes.BROWSER)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulationError(name string, monitorType string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "foo" {
		status           = "ENABLED"
		name             = "%s"
		period           = "EVERY_MINUTE"
		uri              = "https://www.one.newrelic.com"
		type             = "%s"
		locations_public = ["AP_SOUTH_1"]
	  
		custom_header {
		  name  = "Name"
		  value = "browserMonitor"
		}
	  
		enable_screenshot_on_failure_and_script = true
		validation_string                       = "success"
		verify_ssl                              = true
		runtime_type_version                    = "100"
		runtime_type                            = "CHROME_BROWSER"
		script_language                         = "JAVASCRIPT"
		device_orientation                      = "LANDSCAPE"
	  
		tag {
			key    = "butterscotch"
			values = ["cake"]
		}
}`, name, monitorType)
}

func testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulation(name string, monitorType string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "foo" {
		status           = "ENABLED"
		name             = "%s"
		period           = "EVERY_MINUTE"
		uri              = "https://www.one.newrelic.com"
		type             = "%s"
		locations_public = ["AP_SOUTH_1"]
	  
		custom_header {
		  name  = "Name"
		  value = "browserMonitor"
		}
	  
		enable_screenshot_on_failure_and_script = true
		validation_string                       = "success"
		verify_ssl                              = true
		runtime_type_version                    = "100"
		runtime_type                            = "CHROME_BROWSER"
		script_language                         = "JAVASCRIPT"
		device_orientation                      = "LANDSCAPE"
		device_type								= "MOBILE"
	  
		tag {
			key    = "butterscotch"
			values = ["cake"]
		}
}`, name, monitorType)
}

func testAccNewRelicSyntheticsBrowserMonitorConfig_DeviceEmulationLegacyRuntimeError(name string, monitorType string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "foo" {
		status           = "ENABLED"
		name             = "%s"
		period           = "EVERY_MINUTE"
		uri              = "https://www.one.newrelic.com"
		type             = "%s"
		locations_public = ["AP_SOUTH_1"]
	  
		custom_header {
		  name  = "Name"
		  value = "browserMonitor"
		}
	  
		enable_screenshot_on_failure_and_script = true
		validation_string                       = "success"
		verify_ssl                              = true
		device_orientation                      = "LANDSCAPE"
		device_type								= "MOBILE"
	  
		tag {
			key    = "butterscotch"
			values = ["cake"]
		}
}`, name, monitorType)
}
func TestAccNewRelicSyntheticsSimpleMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
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
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Technical limitations with the API prevent us from setting the following attributes.
					"locations_public",
					"locations_private",
					"bypass_head_request",
					"treat_redirect_as_failure",
					"runtime_type",
					"runtime_type_version",
					"script_language",
					"tag",
					"enable_screenshot_on_failure_and_script",
					"custom_header",
					"device_orientation",
					"device_type",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsSimpleMonitorConfig(name string, monitorType string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	custom_header {
		name  = "Name"
		value = "simpleMonitor"
	}
	treat_redirect_as_failure = false
	validation_string         = "success"
	bypass_head_request       = false
	verify_ssl                = false
	locations_public          = ["AP_SOUTH_1"]
	name                      = "%s"
	period                    = "EVERY_MINUTE"
	status                    = "ENABLED"
	type                      = "%s"
	tag {
		key    = "pineapple"
		values = ["pizza"]
	}
	uri = "https://www.one.newrelic.com"
}`, name, monitorType)
}

func testAccNewRelicSyntheticsSimpleMonitorConfigUpdated(name string, monitorType string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	custom_header {
		name  = "name"
		value = "simpleMonitorUpdated"
	}
	treat_redirect_as_failure = true
	validation_string         = "succeeded"
	bypass_head_request       = false
	verify_ssl                = true
	locations_public          = ["AP_SOUTH_1", "AP_EAST_1"]
	name                      = "%s-updated"
	period                    = "EVERY_5_MINUTES"
	status                    = "DISABLED"
	type                      = "%s"
	tag {
		key    = "pineapple"
		values = ["pizza", "cake"]
	}
	uri = "https://www.one.newrelic.com"
}`, name, monitorType)
}

func TestAccNewRelicSyntheticsSimpleBrowserMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_monitor.bar"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
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
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Technical limitations with the API prevent us from setting the following attributes.
					"locations_public",
					"locations_private",
					"bypass_head_request",
					"treat_redirect_as_failure",
					"tag",
					"enable_screenshot_on_failure_and_script",
					"custom_header",
					"runtime_type_version",
					"runtime_type",
					"script_language",
					"device_orientation",
					"device_type",
					"browsers",
					"devices",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfig(name string, monitorType string) string {
	return fmt.Sprintf(`
	resource "newrelic_synthetics_monitor" "bar" {
		custom_header {
			name	= "custom-header-1"
			value	= "header-value-1"
		}
		custom_header {
			name	= "customer-header-2"
			value	= "header-value-2"
		}
		browsers = ["CHROME", "FIREFOX"]
		devices = ["DESKTOP","MOBILE_PORTRAIT", "TABLET_LANDSCAPE", "MOBILE_LANDSCAPE", "TABLET_PORTRAIT"]
		enable_screenshot_on_failure_and_script	=	true
		validation_string	=	"success"
		verify_ssl	=	true
		locations_public	=	["AP_SOUTH_1"]
		name	=	"%s"
		period	=	"EVERY_MINUTE"
		runtime_type_version	=	"100"
		runtime_type	=	"CHROME_BROWSER"
		script_language	=	"JAVASCRIPT"
		status	=	"ENABLED"
		type	=	"%s"
		uri	=	"https://www.one.newrelic.com"
	}`, name, monitorType)
}

func testAccNewRelicSyntheticsSimpleBrowserMonitorConfigUpdated(name string, monitorType string) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "bar" {
			custom_header {
				name  = "name"
				value = "simple_browser"
			}
			browsers = ["CHROME"]
			devices = ["DESKTOP","MOBILE_PORTRAIT", "TABLET_LANDSCAPE"]
			enable_screenshot_on_failure_and_script	=	false
			validation_string	=	"success"
			verify_ssl	=	false
			locations_public	=	["AP_SOUTH_1","AP_EAST_1"]
			name	=	"%s-Updated"
			period	=	"EVERY_5_MINUTES"
			runtime_type_version	=	"100"
			runtime_type	=	"CHROME_BROWSER"
			script_language	=	"JAVASCRIPT"
			status	=	"DISABLED"
			type	=	"%s"
			uri	=	"https://www.one.newrelic.com"
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

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(30 * time.Second)

		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if string((*found).GetGUID()) != rs.Primary.ID {
			return fmt.Errorf("the monitor is not found %v - %v", (*found).GetGUID(), rs.Primary.ID)
		}

		if rs.Primary.Attributes["monitor_id"] != string((*found).(*entities.SyntheticMonitorEntity).MonitorId) {
			return fmt.Errorf("the monitor id doesnot not match expected: %v", rs.Primary.Attributes["monitor_id"])
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

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(30 * time.Second)

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}

	return nil
}

func TestSyntheticsSimpleBrowserMonitor_PeriodInMinutes(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := createIntegrationTestClient(t)

	monitorName := generateNameForIntegrationTestResource()

	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{
		Locations: synthetics.SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: synthetics.SyntheticsMonitorStatus(synthetics.SyntheticsMonitorStatusTypes.ENABLED),
		Uri:    "https://www.one.newrelic.com",

		// Simple Browser Monitors also seem to have the API return an error when no runtime
		// is specified; updated the test accordingly, after the Legacy Runtime EOL.
		Runtime: &synthetics.SyntheticsRuntimeInput{
			RuntimeType:        "CHROME_BROWSER",
			RuntimeTypeVersion: "100",
			ScriptLanguage:     "JAVASCRIPT",
		},
		AdvancedOptions: synthetics.SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
			ResponseValidationText:             "SUCCESS",
			CustomHeaders: &[]synthetics.SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "Synthetics",
				},
			},
			UseTlsValidation: &tv,
		},
	}

	createSimpleBrowserMonitor, err := a.SyntheticsCreateSimpleBrowserMonitor(testAccountID, simpleBrowserMonitorInput)
	var periodInMinutes int = syntheticsMonitorPeriodInMinutesValueMap[createSimpleBrowserMonitor.Monitor.Period]

	require.NoError(t, err)
	require.NotNil(t, createSimpleBrowserMonitor)
	require.Equal(t, 0, len(createSimpleBrowserMonitor.Errors))
	require.Equal(t, 5, periodInMinutes)
}

func createIntegrationTestClient(t *testing.T) synthetics.Synthetics {
	tc := mock.NewIntegrationTestConfig(t)

	return synthetics.New(tc)
}
