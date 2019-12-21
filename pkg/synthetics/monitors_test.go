// +build unit

package synthetics

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testMonitorID = "72733a02-9701-4279-8ac3-8f6281a5a1a9"
	testTime, _   = time.Parse(syntheticsTimeFormat, "2019-11-27T19:11:05.076+0000")
	testTimestamp = Time(testTime)

	testMonitorOptions = MonitorOptions{
		ValidationString:       "",
		VerifySSL:              false,
		BypassHEADRequest:      false,
		TreatRedirectAsFailure: false,
	}

	testMonitor = Monitor{
		ID:           testMonitorID,
		Name:         "test-synthetics-monitor",
		Type:         MonitorTypes.Ping,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       MonitorStatus.Disabled,
		SLAThreshold: 7,
		UserID:       0,
		APIVersion:   "LATEST",
		ModifiedAt:   &testTimestamp,
		CreatedAt:    &testTimestamp,
		Options:      testMonitorOptions,
	}

	testMonitorJson = `
		{
			"id": "72733a02-9701-4279-8ac3-8f6281a5a1a9",
			"name": "test-synthetics-monitor",
			"type": "SIMPLE",
			"frequency": 15,
			"uri": "https://google.com",
			"locations": [
				"AWS_US_EAST_1"
			],
			"status": "DISABLED",
			"slaThreshold": 7,
			"options": {

			},
			"modifiedAt": "2019-11-27T19:11:05.076+0000",
			"createdAt": "2019-11-27T19:11:05.076+0000",
			"userId": 0,
			"apiVersion": "LATEST"
		}
	`
)

func TestListMonitors(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "monitors": [%s] }`, testMonitorJson)
	synthetics := newMockResponse(t, respJSON, http.StatusOK)

	expected := []Monitor{
		testMonitor,
	}

	actual, err := synthetics.ListMonitors()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetMonitor(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testMonitorJson, http.StatusOK)

	actual, err := synthetics.GetMonitor(testMonitorID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testMonitor, actual)
}

func TestCreateMonitor(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testMonitorJson, http.StatusOK)

	id, err := synthetics.CreateMonitor(testMonitor)

	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestUpdateMonitor(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testMonitorJson, http.StatusOK)

	err := synthetics.UpdateMonitor(testMonitor)

	assert.NoError(t, err)
}

func TestDeleteMonitor(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testMonitorJson, http.StatusOK)

	err := synthetics.DeleteMonitor(testMonitor.ID)

	assert.NoError(t, err)
}
