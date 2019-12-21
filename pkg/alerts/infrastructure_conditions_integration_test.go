// +build integration

package alerts

import (
	"fmt"
	"os"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	testIntegrationInfrastructureConditionRandStr = nr.RandSeq(5)
	testIntegrationInfrastructureConditionPolicy  = AlertPolicy{
		Name: fmt.Sprintf("test-integration-infrastructure-conditions-%s",
			testIntegrationInfrastructureConditionRandStr),
		IncidentPreference: "PER_POLICY",
	}
	testIntegrationInfrastructureConditionThreshold = InfrastructureConditionThreshold{
		Duration: 6,
		Value:    0,
	}

	testIntegrationInfrastructureCondition = InfrastructureCondition{
		Comparison:   "equal",
		Critical:     &testIntegrationInfrastructureConditionThreshold,
		Enabled:      true,
		Name:         "Java is running",
		ProcessWhere: "(commandName = 'java')",
		Type:         "infra_process_running",
		Where:        "(hostname LIKE '%cassandra%')",
	}
)

func TestIntegrationListInfrastructureConditions(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	alerts := New(config.Config{
		APIKey: apiKey,
	})

	// Setup
	policy, err := alerts.CreateAlertPolicy(testIntegrationInfrastructureConditionPolicy)

	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := alerts.DeleteAlertPolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	testIntegrationInfrastructureCondition.PolicyID = policy.ID
	created, err := alerts.CreateInfrastructureCondition(testIntegrationInfrastructureCondition)

	require.NoError(t, err)
	require.NotZero(t, created)

	// Test: List
	conditions, err := alerts.ListInfrastructureConditions(policy.ID)

	require.NoError(t, err)
	require.Greater(t, len(conditions), 0)

	// Test: Get
	condition, err := alerts.GetInfrastructureCondition(created.ID)

	require.NoError(t, err)
	require.NotZero(t, condition)

	// Test: Update
	testIntegrationInfrastructureCondition.Name = "Updated"
	updated, err := alerts.UpdateInfrastructureCondition(testIntegrationInfrastructureCondition)

	require.NoError(t, err)
	require.NotZero(t, updated)

	// Test: Delete
	err = alerts.DeleteInfrastructureCondition(created.ID)

	require.NoError(t, err)
}
