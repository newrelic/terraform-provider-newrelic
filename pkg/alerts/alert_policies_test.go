// +build unit

package alerts

import (
	"net/http"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
	"github.com/stretchr/testify/require"
)

var (
	testTimestamp = serialization.EpochTime(time.Unix(1575438237690, 0))
)

func TestGetAlertPolicy(t *testing.T) {
	t.Parallel()
	respJSON := `
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
				"created_at": 1575438237690,
				"updated_at": 1575438237690
			}
		]
	}`
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	// GetAlertPolicy returns a pointer *AlertPolicy
	expected := &AlertPolicy{
		ID:                 579506,
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.GetAlertPolicy(579506)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestListAlertPolicies(t *testing.T) {
	t.Parallel()
	respJSON := `
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
				"created_at": 1575438237690,
				"updated_at": 1575438237690
			}
		]
	}`
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := []AlertPolicy{
		{
			ID:                 579506,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-1",
			CreatedAt:          &testTimestamp,
			UpdatedAt:          &testTimestamp,
		},
		{
			ID:                 579509,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-2",
			CreatedAt:          &testTimestamp,
			UpdatedAt:          &testTimestamp,
		},
	}

	actual, err := alerts.ListAlertPolicies(nil)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestListAlertPoliciesWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "test-alert-policy-1"

	alerts := newTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		if name != expectedName {
			t.Errorf(`expected name filter "%s", recieved: "%s"`, expectedName, name)
		}

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
			CreatedAt:          &testTimestamp,
			UpdatedAt:          &testTimestamp,
		},
	}

	params := ListAlertPoliciesParams{
		Name: expectedName,
	}

	actual, err := alerts.ListAlertPolicies(&params)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestCreateAlertPolicy(t *testing.T) {
	t.Parallel()
	respJSON := `
	{
		"policy": {
			"id": 123,
			"incident_preference": "PER_POLICY",
			"name": "test-alert-policy-1",
			"created_at": 1575438237690,
			"updated_at": 1575438237690
		}
	}
	`
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	policy := AlertPolicy{
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
	}

	expected := &AlertPolicy{
		ID:                 123,
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.CreateAlertPolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestUpdateAlertPolicy(t *testing.T) {
	t.Parallel()
	respJSON := `
	{
		"policy": {
			"id": 123,
			"incident_preference": "PER_CONDITION",
			"name": "name-updated",
			"created_at": 1575438237690,
			"updated_at": 1575438237690
		}
	}`

	alerts := newMockResponse(t, respJSON, http.StatusOK)

	// Original policy
	policy := AlertPolicy{
		ID:                 123,
		IncidentPreference: "PER_POLICY",
		Name:               "name",
	}

	// Updated policy expectation
	expected := &AlertPolicy{
		ID:                 123,
		IncidentPreference: "PER_CONDITION",
		Name:               "name-updated",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.UpdateAlertPolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestDeleteAlertPolicy(t *testing.T) {
	t.Parallel()
	respJSON := `
	{
		"policy": {
			"id": 123,
			"incident_preference": "PER_CONDITION",
			"name": "name-updated",
			"created_at": 1575438237690,
			"updated_at": 1575438237690
		}
	}`
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &AlertPolicy{
		ID:                 123,
		IncidentPreference: "PER_CONDITION",
		Name:               "name-updated",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.DeleteAlertPolicy(123)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}
