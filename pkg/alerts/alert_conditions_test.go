package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListConditionsResponseJSON = `{
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

	testConditionJSON = `{
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

	testConditionUpdateJSON = `{
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
					"duration": "10",
					"operator": "below",
					"priority": "warning",
					"threshold": "0.5",
					"time_function": "all"
				}
			]
		}
	}`
)

func TestListConditions(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListConditionsResponseJSON, http.StatusOK)

	expected := []*Condition{
		{
			PolicyID:   333,
			ID:         123,
			Type:       "apm_app_metric",
			Name:       "Apdex (High)",
			Enabled:    true,
			Entities:   []string{"321"},
			Metric:     "apdex",
			RunbookURL: "",
			Terms: []ConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    0.9,
					TimeFunction: "all",
				},
			},
			UserDefined: ConditionUserDefined{
				Metric:        "",
				ValueFunction: "",
			},
			Scope:               "application",
			GCMetric:            "",
			ViolationCloseTimer: 0,
		},
	}

	actual, err := alerts.ListConditions(333)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListConditionsResponseJSON, http.StatusOK)

	expected := &Condition{
		PolicyID:   333,
		ID:         123,
		Type:       "apm_app_metric",
		Name:       "Apdex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: ConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	actual, err := alerts.GetCondition(333, 123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testConditionJSON, http.StatusCreated)

	condition := Condition{
		PolicyID:   333,
		Type:       "apm_app_metric",
		Name:       "Adpex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: ConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	expected := &Condition{
		PolicyID:   333,
		ID:         123,
		Type:       "apm_app_metric",
		Name:       "Apdex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: ConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	actual, err := alerts.CreateCondition(condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testConditionUpdateJSON, http.StatusCreated)

	condition := Condition{
		PolicyID:   333,
		Type:       "apm_app_metric",
		Name:       "Adpex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: ConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	expected := &Condition{
		PolicyID:   333,
		ID:         123,
		Type:       "apm_app_metric",
		Name:       "Apdex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []ConditionTerm{
			{
				Duration:     10,
				Operator:     "below",
				Priority:     "warning",
				Threshold:    0.5,
				TimeFunction: "all",
			},
		},
		UserDefined: ConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	actual, err := alerts.UpdateCondition(condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testConditionJSON, http.StatusOK)

	expected := &Condition{
		ID:         123,
		Type:       "apm_app_metric",
		Name:       "Apdex (High)",
		Enabled:    true,
		Entities:   []string{"321"},
		Metric:     "apdex",
		RunbookURL: "",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    0.9,
				TimeFunction: "all",
			},
		},
		UserDefined: ConditionUserDefined{
			Metric:        "",
			ValueFunction: "",
		},
		Scope:               "application",
		GCMetric:            "",
		ViolationCloseTimer: 0,
	}

	actual, err := alerts.DeleteCondition(123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
