// +build integration

package synthetics

import (
	"fmt"
	"os"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	testIntegrationMonitor = Monitor{
		Type:         MonitorTypes.Simple,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       MonitorStatus.Disabled,
		SLAThreshold: 7,
		UserID:       0,
		APIVersion:   "LATEST",
		Options:      MonitorOptions{},
	}
)

func TestIntegrationMonitors(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	synthetics := New(config.Config{
		APIKey: apiKey,
	})

	rand := nr.RandSeq(5)
	testIntegrationMonitor.Name = fmt.Sprintf("test-synthetics-monitor-%s", rand)

	// Test: Create
	monitorID, err := synthetics.CreateMonitor(testIntegrationMonitor)

	require.NoError(t, err)
	require.NotNil(t, monitorID)

	// Test: List
	monitors, err := synthetics.ListMonitors()

	require.NoError(t, err)
	require.NotNil(t, monitors)
	require.Greater(t, len(monitors), 0)

	// Test: Get
	monitor, err := synthetics.GetMonitor(monitorID)

	require.NoError(t, err)
	require.NotNil(t, monitor)

	// Test: Update
	updatedName := fmt.Sprintf("test-synthetics-monitor-updated-%s", rand)
	monitor.Name = updatedName
	err = synthetics.UpdateMonitor(*monitor)

	require.NoError(t, err)

	monitor, err = synthetics.GetMonitor(monitorID)

	require.NoError(t, err)
	require.Equal(t, updatedName, monitor.Name)

	// Test: Delete
	err = synthetics.DeleteMonitor(monitorID)

	require.NoError(t, err)
}
