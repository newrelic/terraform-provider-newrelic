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

func TestGetPolicy(t *testing.T) {
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

	// GetPolicy returns a pointer *Policy
	expected := &Policy{
		ID:                 579506,
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.GetPolicy(579506)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestListPolicies(t *testing.T) {
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

	expected := []Policy{
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

	actual, err := alerts.ListPolicies(nil)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestListPoliciesWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "test-alert-policy-1"

	alerts := newTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		if name != expectedName {
			t.Errorf(`expected name filter "%s", received: "%s"`, expectedName, name)
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

	expected := []Policy{
		{
			ID:                 579506,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-1",
			CreatedAt:          &testTimestamp,
			UpdatedAt:          &testTimestamp,
		},
	}

	params := ListPoliciesParams{
		Name: expectedName,
	}

	actual, err := alerts.ListPolicies(&params)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestCreatePolicy(t *testing.T) {
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

	policy := Policy{
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
	}

	expected := &Policy{
		ID:                 123,
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.CreatePolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestUpdatePolicy(t *testing.T) {
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
	policy := Policy{
		ID:                 123,
		IncidentPreference: "PER_POLICY",
		Name:               "name",
	}

	// Updated policy expectation
	expected := &Policy{
		ID:                 123,
		IncidentPreference: "PER_CONDITION",
		Name:               "name-updated",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.UpdatePolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestDeletePolicy(t *testing.T) {
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

	expected := &Policy{
		ID:                 123,
		IncidentPreference: "PER_CONDITION",
		Name:               "name-updated",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.DeletePolicy(123)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}
