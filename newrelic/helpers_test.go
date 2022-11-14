//go:build integration || unit
// +build integration unit

//
// Test Helpers
//

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

var (
	nrInternalAccount = os.Getenv("NR_ACC_TESTING") != ""
)

func newIntegrationTestClient() (*newrelic.NewRelic, error) {
	return newrelic.New(newrelic.ConfigPersonalAPIKey(testAccAPIKey))
}

func testAccDeleteNewRelicAlertPolicy(name string) func() {
	return func() {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		params := alerts.AlertsPoliciesSearchCriteriaInput{}
		policies, _ := client.Alerts.QueryPolicySearch(providerConfig.AccountID, params)

		for _, p := range policies {
			if p.Name == name {
				_, _ = client.Alerts.DeletePolicyMutation(providerConfig.AccountID, p.ID)
				break
			}
		}
	}
}

// A custom check function to log the internal state during a test run.
// nolint:deadcode,unused
func logState(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		t.Logf("State: %s\n", s)

		return nil
	}
}

func testAccImportStateIDFunc(resourceName string, metadata string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		idWithMetadata := fmt.Sprintf("%s:%s", rs.Primary.ID, metadata)

		return idWithMetadata, nil
	}
}

func avoidEmptyAccountID() {
	if os.Getenv("NEW_RELIC_ACCOUNT_ID") == "" {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", "11111")
	}
}

// retrieveIdsFromEnvOrSkip retrieves certain variables from the environment needed by the api access key tests.
func retrieveIdsFromEnvOrSkip(t *testing.T, envKey string) (string, int) {
	envValue := os.Getenv(envKey)
	var id int
	if envValue == "" {
		t.Skip(fmt.Sprintf("Skipping test: config %s not set", envKey))
	} else {
		id, _ = strconv.Atoi(envValue)
	}

	return envValue, id
}
