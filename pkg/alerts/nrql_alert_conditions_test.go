// +build unit

package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListNrqlAlertConditionsResponseJSON = `{
		"nrql_conditions": [
			{
				"type": "static",
				"id": 12345,
				"name": "NRQL Test Alert",
				"enabled": true,
				"value_function": "single_value",
				"violation_time_limit_seconds": 3600,
				"terms": [
					{
						"duration": "5",
						"operator": "above",
						"priority": "critical",
						"threshold": "1",
						"time_function": "all"
					}
				],
				"nrql": {
					"query": "SELECT count(*) FROM Transactions",
					"since_value": "3"
				}
			}
		]
	}`

	testNrqlAlertConditionJSON = `{
		"nrql_condition": {
			"type": "static",
			"id": 12345,
			"name": "NRQL Test Alert",
			"enabled": true,
			"value_function": "single_value",
			"violation_time_limit_seconds": 3600,
			"terms": [
				{
					"duration": "5",
					"operator": "above",
					"priority": "critical",
					"threshold": "1",
					"time_function": "all"
				}
			],
			"nrql": {
				"query": "SELECT count(*) FROM Transactions",
				"since_value": "3"
			}
		}
	}`

	testNrqlConditionUpdatedJSON = `{
		"nrql_condition": {
			"type": "static",
			"id": 12345,
			"name": "NRQL Test Alert Updated",
			"enabled": false,
			"value_function": "single_value",
			"violation_time_limit_seconds": 3600,
			"terms": [
				{
					"duration": "5",
					"operator": "below",
					"priority": "critical",
					"threshold": "1",
					"time_function": "all"
				}
			],
			"nrql": {
				"query": "SELECT count(*) FROM Transactions",
				"since_value": "3"
			},
			"runbook_url": "https://www.example.com/docs"
		}
	}`
)

func TestListNrqlAlertConditions(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListNrqlAlertConditionsResponseJSON, http.StatusOK)

	expected := []*NrqlCondition{
		{
			Nrql: NrqlQuery{
				Query:      "SELECT count(*) FROM Transactions",
				SinceValue: "3",
			},
			Terms: []AlertConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    1,
					TimeFunction: "all",
				},
			},
			Type:                "static",
			Name:                "NRQL Test Alert",
			RunbookURL:          "",
			ValueFunction:       "single_value",
			PolicyID:            123,
			ID:                  12345,
			ViolationCloseTimer: 3600,
			Enabled:             true,
		},
	}

	actual, err := alerts.ListNrqlAlertConditions(123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetNrqlAlertCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListNrqlAlertConditionsResponseJSON, http.StatusOK)

	expected := &NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ValueFunction:       "single_value",
		PolicyID:            123,
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
	}

	actual, err := alerts.GetNrqlAlertCondition(123, 12345)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateNrqlAlertCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testNrqlAlertConditionJSON, http.StatusCreated)

	condition := NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ValueFunction:       "single_value",
		PolicyID:            123,
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
	}

	expected := &condition

	actual, err := alerts.CreateNrqlAlertCondition(condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateNrqlAlertCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testNrqlConditionUpdatedJSON, http.StatusCreated)

	condition := NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ValueFunction:       "single_value",
		PolicyID:            123,
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
	}

	expected := &NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "below",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert Updated",
		RunbookURL:          "https://www.example.com/docs",
		ValueFunction:       "single_value",
		PolicyID:            123,
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             false,
	}

	actual, err := alerts.UpdateNrqlAlertCondition(condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteNrqlAlertCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testNrqlAlertConditionJSON, http.StatusOK)

	expected := &NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ValueFunction:       "single_value",
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
	}

	actual, err := alerts.DeleteNrqlAlertCondition(12345)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
