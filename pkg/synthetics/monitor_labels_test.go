// +build unit

package synthetics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGetMonitorLabelsJson = `
	{
	"labels": [{
		"type": "Testing",
		"value": "Mbnhl",
		"href": "https://synthetics.newrelic.com/synthetics/api/v4/monitors/labels/Testing:Mbnhl"
	}],
	"count": 1
}
	`
)

func TestGetMonitorLabel(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testGetMonitorLabelsJson, http.StatusOK)

	r, err := synthetics.GetMonitorLabels(testMonitorID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r))
	assert.Equal(t, r[0].Type, "Testing")
	assert.Equal(t, r[0].Value, "Mbnhl")
}

func TestAddMonitorLabel(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, "", http.StatusOK)

	err := synthetics.AddMonitorLabel(testMonitorID, "test", "test")
	assert.NoError(t, err)
}

func TestDeleteMonitorLabel(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, "", http.StatusOK)

	err := synthetics.DeleteMonitorLabel(testMonitorID, "test", "test")
	assert.NoError(t, err)
}
