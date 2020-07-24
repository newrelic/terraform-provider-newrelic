// +build integration

package newrelic

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

var (
	nrInternalAccount = os.Getenv("NR_ACC_TESTING") != ""
)

func testAccDeleteNewRelicAlertPolicy(name string) func() {
	return func() {
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		params := alerts.ListPoliciesParams{
			Name: name,
		}
		policies, _ := client.Alerts.ListPolicies(&params)

		for _, p := range policies {
			if p.Name == name {
				_, _ = client.Alerts.DeletePolicy(p.ID)
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
