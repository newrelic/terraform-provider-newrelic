// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAlertConditions(t *testing.T) {
	t.Parallel()

	var (
		testAlertConditionRandStr = nr.RandSeq(5)
		testAlertConditionPolicy  = AlertPolicy{
			Name: fmt.Sprintf("test-integration-alert-conditions-%s",
				testAlertConditionRandStr),
			IncidentPreference: "PER_POLICY",
		}
		testAlertCondition = AlertCondition{
			Type:       "apm_app_metric",
			Name:       "Adpex (High)",
			Enabled:    true,
			Entities:   []string{},
			Metric:     "apdex",
			RunbookURL: "",
			Terms: []AlertConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    0.9,
					TimeFunction: "all",
				},
			},
			UserDefined: AlertConditionUserDefined{
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
	policy, err := client.CreateAlertPolicy(testAlertConditionPolicy)

	require.NoError(t, err)

	testAlertCondition.PolicyID = policy.ID

	// Deferred teardown
	defer func() {
		_, err := client.DeleteAlertPolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	createResult, err := client.CreateAlertCondition(testAlertCondition)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get
	listResult, err := client.ListAlertConditions(createResult.PolicyID)

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	readResult, err := client.GetAlertCondition(createResult.PolicyID, createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = "Apdex Update Test"
	updateResult, err := client.UpdateAlertCondition(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, "Apdex Update Test", updateResult.Name)

	// Test: Delete
	result, err := client.DeleteAlertCondition(updateResult.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
}
