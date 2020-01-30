package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testMultiLocationSyntheticsConditionPolicyID = 111111
	testMultiLocationSyntheticsConditionsJSON    = `
		{
			"location_failure_conditions": [
				{
					"id": 11425367,
					"name": "Zach is testing",
					"enabled": true,
					"entities": [
						"0d3d7d23-9d7e-44ba-ac74-242ce325161c",
						"767958fe-f47a-49ac-88e3-9fbd0d85b2a0",
						"5d968daa-a226-45b9-b877-1d6601e599d8"
					],
					"terms": [
						{
							"priority": "warning",
							"threshold": 1
						},
						{
							"priority": "critical",
							"threshold": 2
						}
					],
					"violation_time_limit_seconds": 3600
				}
			]
		}
	`

	testMultiLocationSyntheticsConditionJSON = `
		{
			"location_failure_condition": {
				"id": 11425367,
				"name": "Zach is testing",
				"enabled": true,
				"entities": [
					"0d3d7d23-9d7e-44ba-ac74-242ce325161c",
					"767958fe-f47a-49ac-88e3-9fbd0d85b2a0",
					"5d968daa-a226-45b9-b877-1d6601e599d8"
				],
				"terms": [
					{
						"priority": "warning",
						"threshold": 1
					},
					{
						"priority": "critical",
						"threshold": 2
					}
				],
				"violation_time_limit_seconds": 3600
			}
		}
	`

	testMultiLocationSyntheticsCondition = MultiLocationSyntheticsCondition{
		ID:      11425367,
		Name:    "Zach is testing",
		Enabled: true,
		Entities: []string{
			"0d3d7d23-9d7e-44ba-ac74-242ce325161c",
			"767958fe-f47a-49ac-88e3-9fbd0d85b2a0",
			"5d968daa-a226-45b9-b877-1d6601e599d8",
		},
		Terms: []MultiLocationSyntheticsConditionTerm{
			{"warning", 1},
			{"critical", 2},
		},
		ViolationTimeLimitSeconds: 3600,
	}

	testCreateMultiLocationSyntheticsCondition = MultiLocationSyntheticsCondition{
		Name:    "Zach is testing",
		Enabled: true,
		Entities: []string{
			"0d3d7d23-9d7e-44ba-ac74-242ce325161c",
			"767958fe-f47a-49ac-88e3-9fbd0d85b2a0",
			"5d968daa-a226-45b9-b877-1d6601e599d8",
		},
		Terms: []MultiLocationSyntheticsConditionTerm{
			{"warning", 1},
			{"critical", 2},
		},
		ViolationTimeLimitSeconds: 3600,
	}
)

func TestListMultiLocationSyntheticsConditions(t *testing.T) {
	t.Parallel()
	respJSON := testMultiLocationSyntheticsConditionsJSON
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := []*MultiLocationSyntheticsCondition{&testMultiLocationSyntheticsCondition}

	actual, err := alerts.ListMultiLocationSyntheticsConditions(testMultiLocationSyntheticsConditionPolicyID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestCreateMultiLocationSyntheticsCondition(t *testing.T) {
	t.Parallel()
	respJSON := testMultiLocationSyntheticsConditionJSON
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &testMultiLocationSyntheticsCondition

	actual, err := alerts.CreateMultiLocationSyntheticsCondition(testCreateMultiLocationSyntheticsCondition, testMultiLocationSyntheticsConditionPolicyID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestUpdateMultiLocationSyntheticsCondition(t *testing.T) {
	t.Parallel()
	respJSON := testMultiLocationSyntheticsConditionJSON
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &testMultiLocationSyntheticsCondition

	actual, err := alerts.UpdateMultiLocationSyntheticsCondition(testMultiLocationSyntheticsCondition)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestDeleteMultiLocationSyntheticsCondition(t *testing.T) {
	t.Parallel()
	respJSON := testMultiLocationSyntheticsConditionJSON
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &testMultiLocationSyntheticsCondition

	actual, err := alerts.DeleteMultiLocationSyntheticsCondition(testMultiLocationSyntheticsCondition.ID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}
