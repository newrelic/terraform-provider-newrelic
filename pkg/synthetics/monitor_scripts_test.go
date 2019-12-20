// +build unit

package synthetics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testMonitorScriptLocation = MonitorScriptLocation{
		Name: "my_vse_enabled_location",
		HMAC: "MjhiNGE4MjVlMDE1N2M4NDQ4MjNjNDFkZDEyYTRjMmUzZDE3NGJlNjU0MWFmOTJlMzNiODExOGU2ZjhkZTY4ZQ",
	}
	testMonitorScript = MonitorScript{
		Text: "dmFyIGFzc2VydCA9IHJlcXVpcmUoJ2Fzc2VydCcpOw0KYXNzZXJ0LmVxdWFsKCcxJywgJzEnKTs",
		Locations: []MonitorScriptLocation{
			testMonitorScriptLocation,
		},
	}
	testMonitorScriptJson = `
	{
		"scriptText": "dmFyIGFzc2VydCA9IHJlcXVpcmUoJ2Fzc2VydCcpOw0KYXNzZXJ0LmVxdWFsKCcxJywgJzEnKTs",
		"scriptLocations": [
			{
				"name": "my_vse_enabled_location",
				"hmac": "MjhiNGE4MjVlMDE1N2M4NDQ4MjNjNDFkZDEyYTRjMmUzZDE3NGJlNjU0MWFmOTJlMzNiODExOGU2ZjhkZTY4ZQ"
			}
		]
	}
	`
)

func TestGetMonitorScript(t *testing.T) {
	t.Parallel()
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testMonitorScriptJson))

		if err != nil {
			t.Fatal(err)
		}
	}))

	script, err := synthetics.GetMonitorScript(testMonitorID)

	assert.NoError(t, err)
	assert.NotNil(t, script)
}

func TestUpdateMonitorScript(t *testing.T) {
	t.Parallel()
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testMonitorScriptJson))

		if err != nil {
			t.Fatal(err)
		}
	}))

	err := synthetics.UpdateMonitorScript(testMonitorID, testMonitorScript)

	assert.NoError(t, err)
}
