// +build unit

package synthetics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testMonitorScriptLocation = MonitorScriptLocation{
		Name: "AWS_US_EAST_1",
		HMAC: "MjhiNGE4MjVlMDE1N2M4NDQ4MjNjNDFkZDEyYTRjMmUzZDE3NGJlNjU0MWFmOTJlMzNiODExOGU2ZjhkZTY4ZQ",
	}
	testMonitorScript = MonitorScript{
		Text: "asdf",
		Locations: []MonitorScriptLocation{
			testMonitorScriptLocation,
		},
	}
	testMonitorScriptJson = `
	{
		"scriptText": "asdf"
	}
	`
)

func TestGetMonitorScript(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testMonitorScriptJson, http.StatusOK)

	script, err := synthetics.GetMonitorScript(testMonitorID)

	assert.NoError(t, err)
	assert.NotNil(t, script)
}

func TestUpdateMonitorScript(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testMonitorScriptJson, http.StatusOK)

	script, err := synthetics.UpdateMonitorScript(testMonitorID, testMonitorScript)

	assert.NoError(t, err)
	assert.NotNil(t, script)
}
