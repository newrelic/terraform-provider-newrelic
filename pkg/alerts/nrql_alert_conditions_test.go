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
