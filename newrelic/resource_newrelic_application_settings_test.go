//go:build integration || APM

package newrelic

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

var (
	testExpectedApplicationName string
	testApplicationGUID         = "Mzk1NzUyNHxBUE18QVBQTElDQVRJT058NTc4ODU1MzYx"
)

func TestAccNewRelicApplicationSettings_Basic(t *testing.T) {
	resourceName := "newrelic_application_settings.app"
	testExpectedApplicationName = fmt.Sprintf("dummy_app_pro_test_%s", acctest.RandString(10))
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
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           testApplicationGUID,
				ImportStateVerifyIgnore: []string{"is_imported"},
			},
		},
	})
}

func TestAccNewRelicApplicationSettings_UserMonitoringValidation(t *testing.T) {
	expectedMsg, _ := regexp.Compile("use_server_side_config must be set to true when transaction_tracer is configured")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicApplicationConfigUserMonitoringValidation(testExpectedApplicationName),
				ExpectError: expectedMsg,
			},
		},
	})
}

func TestAccNewRelicApplicationSettings_TransactionTracerValidation(t *testing.T) {
	expectedMsg, _ := regexp.Compile("`transaction_threshold_value` must be set when `transaction_threshold_type` is 'VALUE'")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicApplicationConfigTransactionTracerValidation(testExpectedApplicationName),
				ExpectError: expectedMsg,
			},
		},
	})
}

func testAccCheckNewRelicApplicationDestroy(s *terraform.State) error {
	// We expect the application to still exist
	return nil
}

func testAccNewRelicApplicationConfig() string {
	return fmt.Sprintf(`
		resource "newrelic_application_settings" "app" {
			guid = "%[2]s"
			name = "%[1]s"
			app_apdex_threshold = "0.5"
			use_server_side_config = true
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
			enable_slow_sql = false
			enable_thread_profiler = false
		}`, testExpectedApplicationName, testApplicationGUID)
}

func testAccNewRelicApplicationConfigUpdated(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_application_settings" "app" {
			guid = "%[2]s"
			name = "%[1]s-updated"
			app_apdex_threshold = "0.5"
			use_server_side_config = true
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
			enable_slow_sql = true	
			tracer_type = "OPT_OUT"
			enable_thread_profiler = false
		}`, name, testApplicationGUID)
}

func testAccNewRelicApplicationConfigUserMonitoringValidation(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_application_settings" "app" {
			guid = "%[2]s"
			name = "%[1]s-updated"
			app_apdex_threshold = "0.5"
			use_server_side_config = false
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
		}`, name, testApplicationGUID)
}

func testAccNewRelicApplicationConfigTransactionTracerValidation(name string) string {
	return fmt.Sprintf(`
		resource "newrelic_application_settings" "app" {
			guid = "%[2]s"
			name = "%[1]s-updated"
			app_apdex_threshold = "0.5"
			use_server_side_config = true
			transaction_tracer{
			   explain_query_plans{
				 query_plan_threshold_value = "0.5"
				 query_plan_threshold_type = "VALUE"
			   }
			   stack_trace_threshold_value = "0.5"
			   transaction_threshold_type = "VALUE"
			   sql{
				 record_sql = "RAW"
			   }
			}
		}`, name, testApplicationGUID)
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

		time.Sleep(2 * time.Second)
		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
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
}
