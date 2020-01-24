package newrelic

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccExpectedAlertPolicyName  string
	testAccAPIKey                   string
	testAccProviders                map[string]terraform.ResourceProvider
	testAccProvider                 *schema.Provider
	testAccCleanupComplete          = false
)

func init() {
	testAccExpectedAlertChannelName = fmt.Sprintf("%s tf-test@example.com", acctest.RandString(5))
	testAccExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccExpectedAlertPolicyName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"newrelic": testAccProvider,
	}
	testAccAPIKey = os.Getenv("NEWRELIC_API_KEY")
	if v := os.Getenv("NEWRELIC_API_KEY"); v == "" {
		testAccAPIKey = "foo"
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NEWRELIC_API_KEY"); v == "" {
		t.Fatal("NEWRELIC_API_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEWRELIC_LICENSE_KEY"); v == "" {
		t.Fatal("NEWRELIC_LICENSE_KEY must be set for acceptance tests")
	}

	testAccApplicationsCleanup(t)
	testAccCreateApplication(t)
}

func testAccCreateApplication(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(testAccExpectedApplicationName),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_LICENSE_KEY")),
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

func testAccApplicationsCleanup(t *testing.T) {
	// Only run cleanup once per test run
	if testAccCleanupComplete {
		return
	}

	client := apm.New(config.Config{
		APIKey: os.Getenv("NEWRELIC_API_KEY"),
	})

	params := apm.ListApplicationsParams{
		Name: "tf_test",
	}

	applications, err := client.ListApplications(&params)

	if err != nil {
		t.Logf("error fetching applications: %s", err)
	}

	deletedAppCount := 0

	for _, app := range applications {
		if !app.Reporting {
			_, err = client.DeleteApplication(app.ID)

			if err == nil {
				deletedAppCount++
			}
		}
	}

	t.Logf("testacc cleanup of %d applications complete", deletedAppCount)

	testAccCleanupComplete = true
}
