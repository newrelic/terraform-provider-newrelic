// +build integration

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	testExpectedApplicationName string
)

func TestAccNewRelicApplicationSettings_Basic(t *testing.T) {
	resourceName := "newrelic_application_settings.app"
	testExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicApplicationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicApplicationExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicApplicationConfigUpdated(testExpectedApplicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicApplicationExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckNewRelicApplicationDestroy(s *terraform.State) error {
	// We expect the application to still exist
	return nil
}

// The test application for this data source is created in provider_test.go
func testAccNewRelicApplicationConfig() string {
	return fmt.Sprintf(`
resource "newrelic_application_settings" "app" {
	name = "%s"
	app_apdex_threshold = "0.9"
	end_user_apdex_threshold = "0.8"
	enable_real_user_monitoring = true
}
`, testExpectedApplicationName)
}

func testAccNewRelicApplicationConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_application_settings" "app" {
	name = "%s-updated"
	app_apdex_threshold = "0.8"
	end_user_apdex_threshold = "0.7"
	enable_real_user_monitoring = false
}
`, name)
}

func testAccCheckNewRelicApplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no application ID is set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return nil
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		_, err = client.APM.GetApplication(int(id))
		if err != nil {
			return err
		}

		return nil
	}
}

func testPreCheck(t *testing.T) {
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		t.Fatal("NEW_RELIC_API_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_LICENSE_KEY"); v == "" {
		t.Fatal("NEW_RELIC_LICENSE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_ADMIN_API_KEY"); v == "" {
		t.Log("[WARN] NEW_RELIC_ADMIN_API_KEY has not been set for acceptance tests")
	}

	testCreateApplication(t)

	time.Sleep(5 * time.Second)
}

func testCreateApplication(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(testExpectedApplicationName),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
	)

	if err != nil {
		t.Fatalf("Error setting up New Relic application: %s", err)
	}

	if err := app.WaitForConnection(30 * time.Second); err != nil {
		t.Fatalf("Unable to setup New Relic application connection: %s", err)
	}

	app.RecordCustomEvent("terraform test", nil)
	app.Shutdown(30 * time.Second)
}
