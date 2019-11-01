package newrelic

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	newrelicAgent "github.com/newrelic/go-agent"
	newrelicSDK "github.com/paultyng/go-newrelic/v4/api"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccExpectedAlertPolicyName  string
	testAccProviders                map[string]terraform.ResourceProvider
	testAccProvider                 *schema.Provider
)

func init() {
	testAccExpectedAlertChannelName = fmt.Sprintf("%s tf-test@example.com", acctest.RandString(5))
	testAccExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccExpectedAlertPolicyName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"newrelic": testAccProvider,
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

func testAccCheckDestroy() {
	// delete fake application used during tests
	config := newrelicSDK.Config{
		APIKey: os.Getenv(("NEWRELIC_API_KEY")),
	}
	client := newrelicSDK.New(config)

	apps, _ := client.ListApplications()

	for _, app := range apps {
		if app.Name == testAccExpectedApplicationName {
			_ = client.DeleteApplication(app.ID)
			break
		}
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NEWRELIC_API_KEY"); v == "" {
		t.Log(v)
		t.Fatal("NEWRELIC_API_KEY must be set for acceptance tests")
	}

	// setup fake application by logging some metrics
	if v := os.Getenv("NEWRELIC_LICENSE_KEY"); len(v) > 0 {
		config := newrelicAgent.NewConfig(testAccExpectedApplicationName, v)
		app, err := newrelicAgent.NewApplication(config)
		if err != nil {
			t.Log(err)
			t.Fatal("Error setting up New Relic application")
		}

		if err := app.WaitForConnection(30 * time.Second); err != nil {
			t.Log(err)
			t.Fatal("Unable to setup New Relic application connection")
		}

		if err := app.RecordCustomEvent("terraform test", nil); err != nil {
			t.Log(err)
			t.Fatal("Unable to record custom event in New Relic")
		}

		app.Shutdown(30 * time.Second)
	} else {
		t.Log(v)
		t.Fatal("NEWRELIC_LICENSE_KEY must be set for acceptance tests")
	}
}
