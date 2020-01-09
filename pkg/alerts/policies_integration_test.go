// +build integration

package alerts

import (
	"fmt"
	"os"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var (
	testIntegrationPolicyNameRandStr = nr.RandSeq(5)
)

func TestIntegrationPolicy(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	policy := Policy{
		IncidentPreference: "PER_POLICY",
		Name:               fmt.Sprintf("test-alert-policy-%s", testIntegrationPolicyNameRandStr),
	}

	// Test: Create
	createResult := testCreatePolicy(t, client, policy)

	// Test: Read
	readResult := testReadPolicy(t, client, createResult)

	// Test: Update
	updateResult := testUpdatePolicy(t, client, readResult)

	// Test: Delete
	testDeletePolicy(t, client, updateResult)
}

func testCreatePolicy(t *testing.T, client Alerts, policy Policy) *Policy {
	result, err := client.CreatePolicy(policy)

	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testReadPolicy(t *testing.T, client Alerts, policy *Policy) *Policy {
	result, err := client.GetPolicy(policy.ID)

	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testUpdatePolicy(t *testing.T, client Alerts, policy *Policy) *Policy {
	policyUpdated := Policy{
		ID:                 policy.ID,
		IncidentPreference: "PER_CONDITION",
		Name:               fmt.Sprintf("test-alert-policy-updated-%s", testIntegrationPolicyNameRandStr),
	}

	result, err := client.UpdatePolicy(policyUpdated)

	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testDeletePolicy(t *testing.T, client Alerts, policy *Policy) {
	p := *policy
	_, err := client.DeletePolicy(p.ID)

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
