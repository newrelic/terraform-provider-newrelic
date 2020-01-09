// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPluginsConditions(t *testing.T) {
	t.Parallel()

	var (
		randomString = nr.RandSeq(5)
		alertPolicy  = Policy{
			Name:               fmt.Sprintf("test-integration-plugins-policy-%s", randomString),
			IncidentPreference: "PER_POLICY",
		}
		conditionName        = fmt.Sprintf("test-integration-plugins-condition-%s", randomString)
		conditionNameUpdated = fmt.Sprintf("test-integration-plugins-condition-updated-%s", randomString)
		condition            = PluginsCondition{
			Name:              conditionName,
			Enabled:           true,
			Entities:          []string{"212222915"},
			Metric:            "Component/Connection/Clients[connections]",
			MetricDescription: "Connected Clients",
			RunbookURL:        "https://example.com/runbook",
			Terms: []ConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    10,
					TimeFunction: "all",
				},
			},
			ValueFunction: "average",
			Plugin: AlertPlugin{
				ID:   "21709",
				GUID: "net.kenjij.newrelic_redis_plugin",
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	policy, err := client.CreatePolicy(alertPolicy)

	require.NoError(t, err)

	condition.PolicyID = policy.ID

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	createResult, err := client.CreatePluginsCondition(condition)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: List
	listResult, err := client.ListPluginsConditions(createResult.PolicyID)

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	readResult, err := client.GetPluginsCondition(createResult.PolicyID, createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = conditionNameUpdated
	updateResult, err := client.UpdatePluginsCondition(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, conditionNameUpdated, updateResult.Name)

	// Test: Delete
	result, err := client.DeletePluginsCondition(createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
}
