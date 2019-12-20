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
		Text: "dmFyIGFzc2VydCA9IHJlcXVpcmUoJ2Fzc2VydCcpOw0KYXNzZXJ0LmVxdWFsKCcxJywgJzEnKTs=",
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
		APIKey: apiKey,
	})

	// Setup
	rand := nr.RandSeq(5)
	testIntegrationScriptedMonitor.Name = fmt.Sprintf("test-synthetics-monitor-%s", rand)
	monitorID, err := synthetics.CreateMonitor(testIntegrationScriptedMonitor)

	require.NoError(t, err)

	// Test: Update
	err = synthetics.UpdateMonitorScript(monitorID, testIntegrationMonitorScript)

	require.NoError(t, err)

	// Test: Get
	script, err := synthetics.GetMonitorScript(monitorID)

	require.NoError(t, err)
	require.Equal(t, testIntegrationMonitorScript.Text, script.Text)

	// Teardown
	err = synthetics.DeleteMonitor(monitorID)
}
