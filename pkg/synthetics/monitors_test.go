// +build unit

package synthetics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func NewTestSynthetics(handler http.Handler) Synthetics {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

var (
	testMonitorID    = "72733a02-9701-4279-8ac3-8f6281a5a1a9"
	testTimestamp, _ = time.Parse(time.RFC3339, "2019-11-27T19:11:05.076+0000")

	testMonitorOptions = MonitorOptions{
		ValidationString:       "",
		VerifySSL:              false,
		BypassHEADRequest:      false,
		TreatRedirectAsFailure: false,
	}

	testMonitor = Monitor{
		ID:           testMonitorID,
		Name:         "test-synthetics-monitor",
		Type:         MonitorTypes.Simple,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       MonitorStatus.Disabled,
		SLAThreshold: 7,
		UserID:       0,
		APIVersion:   "LATEST",
		ModifiedAt:   testTimestamp,
		CreatedAt:    testTimestamp,
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
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(fmt.Sprintf(`
		{
			"monitors": [%s]
		}
		`, testMonitorJson)))

		if err != nil {
			t.Fatal(err)
		}
	}))

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
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testMonitorJson))

		if err != nil {
			t.Fatal(err)
		}
	}))

	actual, err := synthetics.GetMonitor(testMonitorID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testMonitor, actual)
}

func TestCreateMonitor(t *testing.T) {
	t.Parallel()
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testMonitorJson))

		if err != nil {
			t.Fatal(err)
		}
	}))

	id, err := synthetics.CreateMonitor(testMonitor)

	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestUpdateMonitor(t *testing.T) {
	t.Parallel()
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testMonitorJson))

		if err != nil {
			t.Fatal(err)
		}
	}))

	err := synthetics.UpdateMonitor(testMonitor)

	assert.NoError(t, err)
}

func TestDeleteMonitor(t *testing.T) {
	t.Parallel()
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testMonitorJson))

		if err != nil {
			t.Fatal(err)
		}
	}))

	err := synthetics.DeleteMonitor(testMonitor.ID)

	assert.NoError(t, err)
}
