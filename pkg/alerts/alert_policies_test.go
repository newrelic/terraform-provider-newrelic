package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func NewTestAlerts(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestListAlertPolicies(t *testing.T) {
	t.Parallel()
	alerts := NewTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"policies": [
				{
					"id": 579506,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-1",
					"created_at": 1575438237690,
					"updated_at": 1575438237690
				},
				{
					"id": 579509,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-2",
					"created_at": 1575438240632,
					"updated_at": 1575438240632
				}
			]
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	expected := []AlertPolicy{
		{
			ID:                 579506,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-1",
			CreatedAt:          1575438237690,
			UpdatedAt:          1575438237690,
		},
		{
			ID:                 579509,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-2",
			CreatedAt:          1575438240632,
			UpdatedAt:          1575438240632,
		},
	}

	actual, err := alerts.ListAlertPolicies()

	if err != nil {
		t.Fatalf("ListApplications error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListApplications response is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListApplications response differs from expected: %s", diff)
	}
}
