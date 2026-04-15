//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/stretchr/testify/require"
)

func TestFlattenAlertPolicy_Basic(t *testing.T) {
	t.Parallel()

	mockPolicy := &alerts.AlertsPolicy{
		ID:                  "123456",
		Name:                "Test Policy",
		IncidentPreference:  alerts.AlertsIncidentPreference("PER_POLICY"),
		EntityGuid:          "test-entity-guid-123",
	}

	r := resourceNewRelicAlertPolicy()
	d := r.TestResourceData()
	d.SetId("123456")

	accountID := 12345

	err := flattenAlertPolicy(mockPolicy, d, accountID)
	require.NoError(t, err)

	require.Equal(t, "Test Policy", d.Get("name"))
	require.Equal(t, "PER_POLICY", d.Get("incident_preference"))
	require.Equal(t, accountID, d.Get("account_id"))
	require.Equal(t, "test-entity-guid-123", d.Get("entity_guid"))
}

func TestFlattenAlertPolicy_EmptyEntityGuid(t *testing.T) {
	t.Parallel()

	mockPolicy := &alerts.AlertsPolicy{
		ID:                  "123456",
		Name:                "Test Policy",
		IncidentPreference:  alerts.AlertsIncidentPreference("PER_CONDITION"),
		EntityGuid:          "",
	}

	r := resourceNewRelicAlertPolicy()
	d := r.TestResourceData()
	d.SetId("123456")

	accountID := 12345

	err := flattenAlertPolicy(mockPolicy, d, accountID)
	require.NoError(t, err)

	require.Equal(t, "Test Policy", d.Get("name"))
	require.Equal(t, "PER_CONDITION", d.Get("incident_preference"))
	require.Equal(t, accountID, d.Get("account_id"))
	require.Equal(t, "", d.Get("entity_guid"))
}

func TestFlattenAlertPolicy_AllIncidentPreferences(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		incidentPreference string
		entityGuid         string
	}{
		{
			name:               "PER_POLICY",
			incidentPreference: "PER_POLICY",
			entityGuid:         "entity-guid-per-policy",
		},
		{
			name:               "PER_CONDITION",
			incidentPreference: "PER_CONDITION",
			entityGuid:         "entity-guid-per-condition",
		},
		{
			name:               "PER_CONDITION_AND_TARGET",
			incidentPreference: "PER_CONDITION_AND_TARGET",
			entityGuid:         "entity-guid-per-condition-and-target",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPolicy := &alerts.AlertsPolicy{
				ID:                  "123456",
				Name:                "Test Policy",
				IncidentPreference:  alerts.AlertsIncidentPreference(tc.incidentPreference),
				EntityGuid:          tc.entityGuid,
			}

			r := resourceNewRelicAlertPolicy()
			d := r.TestResourceData()
			d.SetId("123456")

			accountID := 12345

			err := flattenAlertPolicy(mockPolicy, d, accountID)
			require.NoError(t, err)

			require.Equal(t, tc.incidentPreference, d.Get("incident_preference"))
			require.Equal(t, tc.entityGuid, d.Get("entity_guid"))
		})
	}
}
