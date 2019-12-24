package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListAlertConditionsResponseJSON = `{
		"conditions": [
			{
				"id": 123,
				"type": "apm_app_metric",
				"name": "Apdex (High)",
				"enabled": true,
				"entities": [
					"321"
				],
				"metric": "apdex",
				"condition_scope": "application",
				"terms": [
					{
						"duration": "5",
						"operator": "above",
						"priority": "critical",
						"threshold": "0.9",
						"time_function": "all"
					}
				]
			}
		]
	}`

	testCreateAlertConditionJSON = `{
		"condition": {
			"id": 123,
			"type": "apm_app_metric",
			"name": "Apdex (High)",
			"enabled": true,
			"entities": [
				"321"
			],
			"metric": "apdex",
			"condition_scope": "application",
			"violation_close_timer": 0,
			"terms": [
				{
					"duration": "5",
					"operator": "above",
					"priority": "critical",
					"threshold": "0.9",
					"time_function": "all"
				}
			]
		}
	}`
)

func TestListAlertConditions(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListAlertConditionsResponseJSON, http.StatusOK)

	expected := []*AlertCondition{
		{
			PolicyID:   0,
			ID:         123,
			Type:       "apm_app_metric",
			Name:       "Apdex (High)",
			Enabled:    true,
			Entities:   []string{"321"},
			Metric:     "apdex",
			RunbookURL: "",
			Terms: []AlertConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    0.9,
					TimeFunction: "all",
				},
			},
			UserDefined: AlertConditionUserDefined{
				Metric:        "",
				ValueFunction: "",
			},
			Scope:               "application",
			GCMetric:            "",
			ViolationCloseTimer: 0,
		},
	}

	actual, err := alerts.ListAlertConditions(2233)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetAlertCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListAlertConditionsResponseJSON, http.StatusOK)

	expected := &AlertCondition{
		PolicyID:   0,
		ID:         123,
		Type:       "apm_app_metric",
		Name:       "Apdex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: AlertConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	actual, err := alerts.GetAlertCondition(0, 123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateAlertCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testCreateAlertConditionJSON, http.StatusCreated)

	condition := AlertCondition{
		PolicyID:   0,
		Type:       "apm_app_metric",
		Name:       "Adpex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: AlertConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	expected := &AlertCondition{
		PolicyID:   0,
		ID:         123,
		Type:       "apm_app_metric",
		Name:       "Apdex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []AlertConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: AlertConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	actual, err := alerts.CreateAlertCondition(condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
