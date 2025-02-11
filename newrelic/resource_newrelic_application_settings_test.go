//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

var (
	testExpectedApplicationName string
)

func TestAccNewRelicApplicationSettings_Basic(t *testing.T) {
	resourceName := "newrelic_application_settings.app"
	testExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
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
			guid = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTY3MjMyMjY0"
			name = "%[1]s"
			app_apdex_threshold = "0.5"
			enable_real_user_monitoring = true
			transaction_tracer{
			   explain_query_plans{
				 query_plan_threshold_value = "0.5"
				 query_plan_threshold_type = "VALUE"
			   }
			   stack_trace_threshold_value = "0.5"
			   transaction_threshold_value = "0.5"
			   transaction_threshold_type = "VALUE"
			   sql{
				 record_sql = "RAW"
			   }
			}
			error_collector{
			  expected_error_classes = []
			  expected_error_codes = []
			  ignored_error_classes = []
			  ignored_error_codes = []
			}
			tracer_type = "OPT_OUT"
			enable_thread_profiler = false
		}`, testExpectedApplicationName)
}

func testAccNewRelicApplicationConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_application_settings" "app" {
			guid = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTY3MjMyMjY0"
			name = "%[1]s-updated"
			app_apdex_threshold = "0.5"
			enable_real_user_monitoring = true
			transaction_tracer{
			   explain_query_plans{
				 query_plan_threshold_value = "0.5"
				 query_plan_threshold_type = "VALUE"
			   }
			   stack_trace_threshold_value = "0.5"
			   transaction_threshold_value = "0.5"
			   transaction_threshold_type = "VALUE"
			   sql{
				 record_sql = "RAW"
			   }
			}
			error_collector{
			  expected_error_classes = []
			  expected_error_codes = []
			  ignored_error_classes = []
			  ignored_error_codes = []
			}
			tracer_type = "OPT_OUT"
			enable_thread_profiler = false
		}`, name)
}

func testAccCheckNewRelicApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no application ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		time.Sleep(5 * time.Second)
		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		res, foundOk := (*found).(*entities.ApmApplicationEntity)
		if !foundOk {
			return fmt.Errorf("no application found")
		}
		if res.GUID != common.EntityGUID(rs.Primary.ID) {
			return fmt.Errorf("no application found")
		}

		return nil
	}
}

func testPreCheck(t *testing.T) {
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		t.Skipf("NEW_RELIC_API_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_LICENSE_KEY"); v == "" {
		t.Skipf("NEW_RELIC_LICENSE_KEY must be set for acceptance tests")
	}

	//testCreateApplication(t)

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
