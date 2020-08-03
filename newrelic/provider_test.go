// +build integration unit
//
// Test helpers
//

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccExpectedAlertPolicyName  string
	testAccAPIKey                   string
	testAccProviders                map[string]terraform.ResourceProvider
	testAccProvider                 *schema.Provider
	testAccountID                   int
	testAccountName                 string
	//testAccCleanupComplete          = false
)

func init() {
	testAccExpectedAlertChannelName = fmt.Sprintf("%s tf-test@example.com", acctest.RandString(5))
	testAccExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccExpectedAlertPolicyName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"newrelic": testAccProvider,
	}
	testAccAPIKey = os.Getenv("NEW_RELIC_API_KEY")
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		testAccAPIKey = "foo"
	}

	if v, _ := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID")); v != 0 {
		testAccountID = v
	}

	if v := os.Getenv("NEW_RELIC_ACCOUNT_NAME"); v != "" {
		testAccountName = v
	} else {
		testAccountName = "New Relic Terraform Provider Acceptance Testing"
	}
}

func testAccNewRelicProviderConfig(region string, baseURL string, resourceName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	region = "%[1]s"

	%[2]s
}

resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[3]s"
}
`, region, baseURL, resourceName)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		t.Fatal("[WARN] NEW_RELIC_API_KEY has not been set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_LICENSE_KEY"); v == "" {
		t.Fatal("NEW_RELIC_LICENSE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_ADMIN_API_KEY"); v == "" {
		t.Fatal("NEW_RELIC_ADMIN_API_KEY must be set for acceptance tests")
	}

	//testAccApplicationsCleanup(t)
	testAccCreateApplication(t)

	// We need to give the entity search engine time to index the app
	time.Sleep(5 * time.Second)
}

func testAccCreateApplication(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(testAccExpectedApplicationName),
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

// func testAccApplicationsCleanup(t *testing.T) {
// 	// Only run cleanup once per test run
// 	if testAccCleanupComplete {
// 		return
// 	}

// 	client := apm.New(config.Config{
// 		APIKey: os.Getenv("NEW_RELIC_API_KEY"),
// 	})

// 	params := apm.ListApplicationsParams{
// 		Name: "tf_test",
// 	}

// 	applications, err := client.ListApplications(&params)

// 	if err != nil {
// 		t.Logf("error fetching applications: %s", err)
// 	}

// 	deletedAppCount := 0

// 	for _, app := range applications {
// 		if !app.Reporting {
// 			_, err = client.DeleteApplication(app.ID)

// 			if err == nil {
// 				deletedAppCount++
// 			}
// 		}
// 	}

// 	t.Logf("testacc cleanup of %d applications complete", deletedAppCount)

// 	testAccCleanupComplete = true
// }
