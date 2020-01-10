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
	testIntegrationScriptedMonitor = Monitor{
		Type:         MonitorTypes.ScriptedBrowser,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       MonitorStatus.Disabled,
		SLAThreshold: 7,
		UserID:       0,
		APIVersion:   "LATEST",
		Options:      MonitorOptions{},
	}
	testIntegrationMonitorScriptLocation = MonitorScriptLocation{
		Name: "AWS_US_EAST_1",
		HMAC: "MjhiNGE4MjVlMDE1N2M4NDQ4MjNjNDFkZDEyYTRjMmUzZDE3NGJlNjU0MWFmOTJlMzNiODExOGU2ZjhkZTY4ZQ",
	}
	testIntegrationMonitorScript = MonitorScript{
		Text: "asdf",
		Locations: []MonitorScriptLocation{
			testIntegrationMonitorScriptLocation,
		},
	}
)

func TestIntegrationMonitorScripts(t *testing.T) {
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
	testIntegrationScriptedMonitor.Name = fmt.Sprintf("test-synthetics-monitor-%s", rand)
	monitor, err := synthetics.CreateMonitor(testIntegrationScriptedMonitor)

	require.NoError(t, err)

	// Test: Update
	updated, err := synthetics.UpdateMonitorScript(monitor.ID, testIntegrationMonitorScript)

	require.NoError(t, err)
	require.NotNil(t, updated)

	// Test: Get
	script, err := synthetics.GetMonitorScript(monitor.ID)

	require.NoError(t, err)
	require.Equal(t, testIntegrationMonitorScript.Text, script.Text)

	// Deferred teardown
	defer func() {
		err = synthetics.DeleteMonitor(monitor.ID)

		if err != nil {
			t.Logf("error cleaning up monitor %s (%s): %s", monitor.ID, monitor.Name, err)
		}
	}()
}
