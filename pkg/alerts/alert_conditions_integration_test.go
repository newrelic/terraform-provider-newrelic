// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationConditions(t *testing.T) {
	t.Parallel()

	var (
		testConditionRandStr = nr.RandSeq(5)
		testConditionPolicy  = Policy{
			Name: fmt.Sprintf("test-integration-alert-conditions-%s",
				testConditionRandStr),
			IncidentPreference: "PER_POLICY",
		}
		testCondition = Condition{
			Type:       "apm_app_metric",
			Name:       "Adpex (High)",
			Enabled:    true,
			Entities:   []string{},
			Metric:     "apdex",
			RunbookURL: "",
			Terms: []ConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    0.9,
					TimeFunction: "all",
				},
			},
			UserDefined: ConditionUserDefined{
				Metric:        "",
				ValueFunction: "",
			},
			Scope:               "application",
			GCMetric:            "",
			ViolationCloseTimer: 0,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	policy, err := client.CreatePolicy(testConditionPolicy)

	require.NoError(t, err)

	testCondition.PolicyID = policy.ID

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	createResult, err := client.CreateCondition(testCondition)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get
	listResult, err := client.ListConditions(createResult.PolicyID)

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	readResult, err := client.GetCondition(createResult.PolicyID, createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = "Apdex Update Test"
	updateResult, err := client.UpdateCondition(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, "Apdex Update Test", updateResult.Name)

	// Test: Delete
	result, err := client.DeleteCondition(updateResult.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
}
