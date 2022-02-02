//go:build integration || unit
// +build integration unit

//
// Test helpers
//

package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccExpectedAlertPolicyName  string
	testAccAPIKey                   string
	testAccProviders                map[string]*schema.Provider
	testAccProvider                 *schema.Provider
	testAccountID                   int
	testSubaccountID                int
	testAccountName                 string
	//testAccCleanupComplete          = false
)

func init() {
	testAccExpectedAlertChannelName = fmt.Sprintf("%s tf-test@example.com", acctest.RandString(5))
	testAccExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccExpectedAlertPolicyName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"newrelic": testAccProvider,
	}
	testAccAPIKey = os.Getenv("NEW_RELIC_API_KEY")
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		testAccAPIKey = "foo"
	}

	if v, _ := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID")); v != 0 {
		testAccountID = v
	}

	// Used for cross-account scenarios if needed, such as dashboard widgets.
	if v, _ := strconv.Atoi(os.Getenv("NEW_RELIC_SUBACCOUNT_ID")); v != 0 {
		testSubaccountID = v
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
		t.Skipf("[WARN] NEW_RELIC_API_KEY has not been set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_LICENSE_KEY"); v == "" {
		t.Skipf("NEW_RELIC_LICENSE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_ACCOUNT_ID"); v == "" {
		t.Skipf("NEW_RELIC_ACCOUNT_ID must be set for acceptance tests")
	}

	//testAccApplicationsCleanup(t)
	testAccCreateApplication(t)

	// We need to give the entity search engine time to index the app so
	// we try to get the entity, and retry if it fails for a certain amount
	// of time
	client := entities.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})
	params := entities.EntitySearchQueryBuilder{
		Name:   testAccExpectedApplicationName,
		Type:   "APPLICATION",
		Domain: "APM",
	}

	retryErr := resource.RetryContext(context.Background(), 15*time.Second, func() *resource.RetryError {
		entityResults, err := client.GetEntitySearchWithContext(context.Background(), entities.EntitySearchOptions{}, "", params, []entities.EntitySearchSortCriteria{})
		if err != nil {
			return resource.RetryableError(err)
		}

		if entityResults.Count != 1 {
			return resource.RetryableError(fmt.Errorf("Entity not found, or found more than one"))
		}

		return nil
	})

	if retryErr != nil {
		t.Fatalf("Unable to find application entity: %s", retryErr)
	}
}

func testAccCreateApplication(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
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
