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

	// Create
	createResult := testCreateAlertPolicy(t, client, policy)

	// Read
	readResult := testReadAlertPolicy(t, client, createResult)

	// Update
	updateResult := testUpdateAlertPolicy(t, client, readResult)

	// Delete
	testDeleteAlertPolicy(t, client, updateResult)
}

func testCreateAlertPolicy(t *testing.T, client Alerts, policy AlertPolicy) *AlertPolicy {
	result, err := client.CreateAlertPolicy(policy)

	if err != nil {
		t.Fatalf("CreateAlertPolicy error: %s", err)
	}

	if result == nil {
		t.Fatalf("CreateAlertPolicy result is nil")
	}

	return result
}

func testReadAlertPolicy(t *testing.T, client Alerts, policy *AlertPolicy) *AlertPolicy {
	result, err := client.GetAlertPolicy(policy.ID)

	if err != nil {
		t.Fatalf("GetAlertPolicy error: %s", err)
	}

	if result == nil {
		t.Fatalf("GetAlertPolicy result is nil")
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
		t.Fatalf("UpdateAlertPolicy error: %s", err)
	}

	if result == nil {
		t.Fatalf("UpdateAlertPolicy result is nil")
	}

	return result
}

func testDeleteAlertPolicy(t *testing.T, client Alerts, policy *AlertPolicy) {
	p := *policy
	err := client.DeleteAlertPolicy(p.ID)

	if err != nil {
		t.Fatalf("DeleteAlertPolicy error: %s", err)
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
