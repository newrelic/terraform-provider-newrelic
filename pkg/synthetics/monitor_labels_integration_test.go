// +build integration

package synthetics

import (
	"fmt"
	"os"
	"strings"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testIntegrationLabeledMonitor = Monitor{
		Type:         MonitorTypes.APITest,
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

func TestIntegrationGetMonitorLabels(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	synthetics := New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})

	// Setup
	rand := nr.RandSeq(5)
	testIntegrationLabeledMonitor.Name = fmt.Sprintf("test-synthetics-monitor-%s", rand)
	monitor, err := synthetics.CreateMonitor(testIntegrationLabeledMonitor)
	require.NoError(t, err)

	labels, err := synthetics.GetMonitorLabels(monitor.ID)
	require.NoError(t, err)

	assert.Equal(t, len(labels), 0)

	// Test: Add
	err = synthetics.AddMonitorLabel(monitor.ID, "testing", rand)
	require.NoError(t, err)

	labels, err = synthetics.GetMonitorLabels(monitor.ID)
	require.NoError(t, err)
	assert.Equal(t, len(labels), 1)
	assert.Equal(t, (*labels[0]).Type, "Testing")
	assert.Equal(t, (*labels[0]).Value, strings.Title(rand))

	// Test: Delete
	err = synthetics.DeleteMonitorLabel(monitor.ID, "testing", rand)
	require.NoError(t, err)

	// Verify Delete
	labels, err = synthetics.GetMonitorLabels(monitor.ID)
	require.NoError(t, err)
	assert.Equal(t, len(labels), 0)

	// Deferred teardown
	defer func() {
		err = synthetics.DeleteMonitor(monitor.ID)

		if err != nil {
			t.Logf("error cleaning up monitor %s (%s): %s", monitor.ID, monitor.Name, err)
		}
	}()
}
