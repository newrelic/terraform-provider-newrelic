// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationMultiLocationSyntheticsConditions(t *testing.T) {
	t.Parallel()

	var (
		testIntegrationInfrastructureConditionRandStr = nr.RandSeq(5)
		testIntegrationInfrastructureConditionPolicy  = Policy{
			Name: fmt.Sprintf("test-integration-location-failure-condition-%s",
				testIntegrationInfrastructureConditionRandStr),
			IncidentPreference: "PER_POLICY",
		}

		testIntegrationMultiLocationSyntheticsCondition = MultiLocationSyntheticsCondition{
			Name:    fmt.Sprintf("test-integration-location-failure-condition-%s", testIntegrationInfrastructureConditionRandStr),
			Enabled: false,
			Terms: []MultiLocationSyntheticsConditionTerm{
				{"warning", 10},
				{"critical", 11},
			},
		}
	)

	alerts := newIntegrationTestClient(t)

	// Setup
	policy, err := alerts.CreatePolicy(testIntegrationInfrastructureConditionPolicy)
	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := alerts.DeletePolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	created, err := alerts.CreateMultiLocationSyntheticsCondition(testIntegrationMultiLocationSyntheticsCondition, policy.ID)
	require.NoError(t, err)
	require.NotZero(t, created)

	defer func() {
		_, err := alerts.DeleteMultiLocationSyntheticsCondition(created.ID)
		if err != nil {
			t.Logf("error cleaning up location failure condition %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// // Test: List
	conditions, err := alerts.ListMultiLocationSyntheticsConditions(policy.ID)

	require.NoError(t, err)
	require.Greater(t, len(conditions), 0)

	// Test: Update
	created.Name = "Updated"
	created.Enabled = true
	updated, err := alerts.UpdateMultiLocationSyntheticsCondition(*created)

	require.NoError(t, err)
	require.NotZero(t, updated)

}
