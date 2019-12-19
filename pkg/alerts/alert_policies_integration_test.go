// +build integration

package alerts

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestIntegrationAlertPolicyCRUD(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	policy := AlertPolicy{
		IncidentPreference: "PER_POLICY",
		Name:               "integration-test-alert-policy",
	}

	// Test: Create
	createResult := testCreateAlertPolicy(t, client, policy)

	// Test: Read
	readResult := testReadAlertPolicy(t, client, createResult)

	// Test: Update
	updateResult := testUpdateAlertPolicy(t, client, readResult)

	// Test: Delete
	testDeleteAlertPolicy(t, client, updateResult)
}

func testCreateAlertPolicy(t *testing.T, client Alerts, policy AlertPolicy) *AlertPolicy {
	result, err := client.CreateAlertPolicy(policy)

	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testReadAlertPolicy(t *testing.T, client Alerts, policy *AlertPolicy) *AlertPolicy {
	result, err := client.GetAlertPolicy(policy.ID)

	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testUpdateAlertPolicy(t *testing.T, client Alerts, policy *AlertPolicy) *AlertPolicy {
	policyUpdated := AlertPolicy{
		ID:                 policy.ID,
		IncidentPreference: "PER_CONDITION",
		Name:               "integration-test-alert-policy-updated",
	}

	result, err := client.UpdateAlertPolicy(policyUpdated)

	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testDeleteAlertPolicy(t *testing.T, client Alerts, policy *AlertPolicy) {
	p := *policy
	_, err := client.DeleteAlertPolicy(p.ID)

	if err != nil {
		t.Fatal(err)
	}
}

func newClient(t *testing.T) Alerts {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	return New(config.Config{
		APIKey: apiKey,
	})
}
