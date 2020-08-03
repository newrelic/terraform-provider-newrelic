// +build integration

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"
)

func TestExpandAlertConditionEntities(t *testing.T) {
	flattened := []interface{}{123, 456}
	expected := []string{"123", "456"}

	expanded := expandAlertConditionEntities(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenAlertConditionEntities(t *testing.T) {
	expanded := []string{"123", "456"}
	expected := []int{123, 456}

	flattened, err := flattenAlertConditionEntities(&expanded)

	require.NoError(t, err)
	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}

func TestExpandAlertConditionTerms(t *testing.T) {
	flattened := []interface{}{
		map[string]interface{}{
			"duration":      123,
			"operator":      "above",
			"priority":      "critical",
			"threshold":     123.456,
			"time_function": "all",
		},
	}

	expected := []alerts.ConditionTerm{
		{
			Duration:     123,
			Operator:     alerts.OperatorTypes.Above,
			Priority:     alerts.PriorityTypes.Critical,
			Threshold:    123.456,
			TimeFunction: alerts.TimeFunctionTypes.All,
		},
	}

	expanded := expandAlertConditionTerms(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenAlertConditionTerms(t *testing.T) {
	expanded := []alerts.ConditionTerm{
		{
			Duration:     123,
			Operator:     "above",
			Priority:     "critical",
			Threshold:    123.456,
			TimeFunction: alerts.TimeFunctionTypes.All,
		},
		{
			Duration:     123,
			Operator:     "equal",
			Priority:     "warning",
			Threshold:    123.456,
			TimeFunction: alerts.TimeFunctionTypes.Any,
		},
	}

	expected := []map[string]interface{}{
		{
			"duration":      123,
			"operator":      alerts.OperatorTypes.Above,
			"priority":      alerts.PriorityTypes.Critical,
			"threshold":     123.456,
			"time_function": alerts.TimeFunctionTypes.All,
		},
		{
			"duration":      123,
			"operator":      alerts.OperatorTypes.Equal,
			"priority":      alerts.PriorityTypes.Warning,
			"threshold":     123.456,
			"time_function": alerts.TimeFunctionTypes.Any,
		},
	}

	flattened := flattenAlertConditionTerms(&expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}
