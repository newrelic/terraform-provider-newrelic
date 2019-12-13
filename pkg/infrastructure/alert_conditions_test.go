// +build unit

package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func NewTestInfrastructure(handler http.Handler) Infrastructure {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestListAlertConditions(t *testing.T) {
	t.Parallel()
	infra := NewTestInfrastructure(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	expected := []AlertCondition{
		{},
	}

	actual, err := infra.ListAlertConditions()

	if err != nil {
		t.Fatalf("ListAlertConditions error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListAlertConditions response is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListAlertConditions response differs from expected: %s", diff)
	}
}
