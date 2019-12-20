package alerts

import (
	"fmt"
	"os"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
	"github.com/stretchr/testify/require"
)

var (
	testRandStr                      = nr.RandSeq(5)
	testIntegrationSyntheticsMonitor = synthetics.Monitor{
		Name:         fmt.Sprintf("test-synthetics-alert-conditions-monitor-%s", testRandStr),
		Type:         synthetics.MonitorTypes.Simple,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       synthetics.MonitorStatus.Enabled,
		SLAThreshold: 7,
		APIVersion:   "LATEST",
	}
	testIntegrationAlertPolicy = AlertPolicy{
		Name:               fmt.Sprintf("test-synthetics-alert-conditions-policy-%s", testRandStr),
		IncidentPreference: "PER_POLICY",
	}
	testIntegrationSyntheticsCondition = SyntheticsCondition{
		Name: fmt.Sprintf("test-synthetics-alert-condition-%s", testRandStr),
	}
)

func TestIntegrationSyntheticsConditions(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	alerts := New(config.Config{
		APIKey: apiKey,
	})

	synth := synthetics.New(config.Config{
		APIKey: apiKey,
	})

	// Setup
	monitorID, err := synth.CreateMonitor(testIntegrationSyntheticsMonitor)

	require.NoError(t, err)

	policy, err := alerts.CreateAlertPolicy(testIntegrationAlertPolicy)

	require.NoError(t, err)

	// Deferred Teardown
	defer func() {
		// Teardown
		_, err = alerts.DeleteAlertPolicy(policy.ID)
		if err != nil {
			t.Logf("Error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}

		err = synth.DeleteMonitor(monitorID)
		if err != nil {
			t.Logf("Error cleaning up synthetics monitor %s (%s): %s",
				monitorID, testIntegrationSyntheticsMonitor.Name, err)
		}
	}()

	// Test: Create
	testIntegrationSyntheticsCondition.MonitorID = monitorID
	created, err := alerts.CreateSyntheticsCondition(policy.ID, testIntegrationSyntheticsCondition)

	require.NoError(t, err)
	require.NotZero(t, created)

	// Test: List
	conditions, err := alerts.ListSyntheticsConditions(policy.ID)

	require.NoError(t, err)
	require.NotZero(t, conditions)

	// Test: Update
	created.Name = fmt.Sprintf("test-synthetics-alert-condition-updated-%s", testRandStr)
	updated, err := alerts.UpdateSyntheticsCondition(*created)

	require.NoError(t, err)
	require.NotZero(t, updated)

	// Test: Delete
	deleted, err := alerts.DeleteSyntheticsCondition(updated.ID)

	require.NoError(t, err)
	require.NotZero(t, deleted)
}
