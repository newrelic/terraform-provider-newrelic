package newrelic

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestProviderConfig(t *testing.T) {
	c := ProviderConfig{
		PersonalAPIKey: "abc123",
		AccountID:      123,
	}

	hasNerdGraphCreds := c.hasNerdGraphCredentials()

	if !hasNerdGraphCreds {
		t.Error("hasNerdGraphCreds should be true")
	}
}

func TestAccNewRelicProvider_Region(t *testing.T) {
	// This error message will occur when configuring
	// US region with EU API URLs when using the TF test account.
	expectedErrorMsg := "403 response returned"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Region "US"
			{
				Config: testAccNewRelicProviderConfig("US", "", rName),
			},
			// Test: Region "EU"
			{
				Config:      testAccNewRelicProviderConfig("EU", "", rName),
				ExpectError: regexp.MustCompile(expectedErrorMsg),
			},
			// Test: Override US region URLs with EU region URLs (will result in an auth error)
			{
				Config:      testAccNewRelicProviderConfig("US", `nerdgraph_api_url = "https://api.eu.newrelic.com/graphql"`, rName),
				ExpectError: regexp.MustCompile(expectedErrorMsg),
			},
			// Test: Override EU region URLs with US region URLs (should work since the TF acct is US-based)
			{
				Config: testAccNewRelicProviderConfig("EU", `nerdgraph_api_url = "https://api.newrelic.com/graphql"`, rName),
			},
			// Test: Case insensitivity
			{
				Config: testAccNewRelicProviderConfig("us", "", rName),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		t.Log("[WARN] NEW_RELIC_API_KEY has not been set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_LICENSE_KEY"); v == "" {
		t.Fatal("NEW_RELIC_LICENSE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_ADMIN_API_KEY"); v == "" {
		t.Fatal("NEW_RELIC_ADMIN_API_KEY must be set for acceptance tests")
	}

	//testAccApplicationsCleanup(t)
	testAccCreateApplication(t)
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
