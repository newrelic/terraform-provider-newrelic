package synthetics

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

func NewTestSynthetics(handler http.Handler) Synthetics {
	ts := httptest.NewServer(handler)

	c := New(newrelic.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestListApplications(t *testing.T) {
	synthetics := NewTestSynthetics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"monitors": [
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
			]
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	monitorOptions := MonitorOptions{
		ValidationString:       "",
		VerifySSL:              false,
		BypassHEADRequest:      false,
		TreatRedirectAsFailure: false,
	}

	timestamp, _ := time.Parse(time.RFC3339, "2019-11-27T19:11:05.076+0000")

	expected := []Monitor{
		{
			ID:           "72733a02-9701-4279-8ac3-8f6281a5a1a9",
			Name:         "test-synthetics-monitor",
			Type:         "SIMPLE",
			Frequency:    15,
			URI:          "https://google.com",
			Locations:    []string{"AWS_US_EAST_1"},
			Status:       "DISABLED",
			SLAThreshold: 7,
			UserID:       0,
			APIVersion:   "LATEST",
			ModifiedAt:   timestamp,
			CreatedAt:    timestamp,
			Options:      monitorOptions,
		},
	}

	actual, err := synthetics.ListMonitors()

	if err != nil {
		t.Fatalf("ListMonitors error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListMonitors response is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListMonitors response differs from expected: %s", diff)
	}
}

func TestAccListApplications(t *testing.T) {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	synthetics := New(newrelic.Config{
		APIKey: apiKey,
	})

	_, err := synthetics.ListMonitors()

	if err != nil {
		t.Fatalf("ListMonitors error: %s", err)
	}
}
